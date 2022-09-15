package cmd

import (
	"bytes"
	"github.com/gatecheckdev/gatecheck/pkg/exporter/defectDojo"
	"io"
	"os"
	"path"
	"testing"
)

func Test_ConfigInitCmd(t *testing.T) {
	// Provoke an error with improper file name

	actual := new(bytes.Buffer)
	command := NewRootCmd(defectDojo.Exporter{})
	command.SetOut(actual)
	command.SetErr(actual)

	t.Run("bad filename", func(t *testing.T) {
		tempDir := "\000x"
		command.SetArgs([]string{"config", "init", tempDir})
		err := command.Execute()
		if err == nil {
			t.Fatal("Expected file access error")
		}
	})

	t.Run("bad permissions", func(t *testing.T) {
		tempDir := t.TempDir()
		if err := os.Chmod(tempDir, 0000); err != nil {
			t.Fatal(err)
		}

		command.SetArgs([]string{"config", "init", tempDir})
		err := command.Execute()

		if err == nil {
			t.Fatal("Expected file access error")
		}
	})

	t.Run("directory", func(t *testing.T) {
		tempDir := t.TempDir()

		command.SetArgs([]string{"config", "init", tempDir})
		t.Log(tempDir)
		err := command.Execute()
		if err != nil {
			t.Fatal(err)
		}

		fileInfo, err := os.Stat(path.Join(tempDir, "gatecheck.yaml"))
		if err != nil {
			t.Error(fileInfo)
			t.Fatal(err)
		}

		if fileInfo.Size() < 70 {
			t.Fatal("File size is unexpectedly small")
		}
	})

	t.Run("bad file permission", func(t *testing.T) {
		tempDirPath := t.TempDir()
		tempDir := path.Join(tempDirPath, "custom.yaml")
		if err := os.Chmod(tempDirPath, 0000); err != nil {
			t.Fatal(err)
		}

		command.SetArgs([]string{"config", "init", tempDir})
		err := command.Execute()

		if err == nil {
			t.Fatal("Expected file access error")
		}
	})

	t.Run("file already exists", func(t *testing.T) {
		tempFile := path.Join(t.TempDir(), "custom.yaml")
		f, _ := os.Create(tempFile)
		_, _ = io.Copy(f, bytes.NewBufferString("Sample Content"))

		command.SetArgs([]string{"config", "init", tempFile})
		if err := command.Execute(); err == nil {
			t.Fatal("expected error for pre-existing file")
		}
	})

	t.Run("file", func(t *testing.T) {
		tempDir := path.Join(t.TempDir(), "custom-name.yaml")
		command.SetArgs([]string{"config", "init", tempDir})
		err := command.Execute()
		if err != nil {
			t.Fatal(err)
		}

		fileInfo, err := os.Stat(tempDir)
		if err != nil {
			t.Fatal(err)
		}

		if fileInfo.Size() < 70 {
			t.Fatal("File size is unexpectedly small")
		}
	})

	t.Run("file with project name", func(t *testing.T) {
		tempDir := path.Join(t.TempDir(), "custom-name.yaml")
		command.SetArgs([]string{"config", "init", tempDir, "test project name"})
		err := command.Execute()
		if err != nil {
			t.Fatal(err)
		}

		fileInfo, err := os.Stat(tempDir)
		if err != nil {
			t.Fatal(err)
		}

		if fileInfo.Size() < 70 {
			t.Fatal("File size is unexpectedly small")
		}
	})
}
