package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"
)

const ExitFail = 1
const BumpPrefix = "// bump: "
const ChangelogPrefix = "## [UNRELEASED]"

func main() {
	rootFlag := flag.String("r", ".", "root directory for gatecheck")
	tagFlag := flag.String("t", "", "tag name (version v.X.X.X) for release")

	flag.Parse()
	if *tagFlag == "" {
		check(errors.New("missing tag flag -t"))
	}

	root := *rootFlag
	version := *tagFlag

	fmt.Fprintln(os.Stderr, "Bump to Version:", version)
	f, err := os.OpenFile(path.Join(root, "cmd/gatecheck/main.go"), os.O_RDWR, 0664)
	check(err)

	tempBuf := new(bytes.Buffer)
	check(bumpReplace(f, tempBuf, version))
	f.Seek(0, 0)
	_, err = tempBuf.WriteTo(f)
	check(err)

	f, err = os.OpenFile(path.Join(root, "CHANGELOG.md"), os.O_RDWR, 0664)
	tempBuf = new(bytes.Buffer)
	check(err)
	check(bumpChangelog(f, tempBuf, version))
	f.Seek(0, 0)
	b := bytes.TrimSuffix(tempBuf.Bytes(), []byte("\n"))
	_, err = bytes.NewReader(b).WriteTo(f)
	check(err)
}

func bumpChangelog(src io.Reader, dest io.Writer, newVersion string) error {
	content, err := io.ReadAll(src)
	if err != nil {
		return err
	}

	newHeader := fmt.Sprintf("%s\n\n## [%s] - %s",
		ChangelogPrefix,
		strings.Replace(newVersion, "v", "", 1),
		time.Now().Format("2006-01-02"))

	_, err = strings.NewReader(strings.Replace(string(content), ChangelogPrefix, newHeader, 1)).WriteTo(dest)
	return err
}
func bumpReplace(src io.Reader, dest io.Writer, replaceValue string) error {
	scanner := bufio.NewScanner(src)
	nextLineContent := ""
	for scanner.Scan() {
		if nextLineContent != "" {
			_, err := dest.Write([]byte(nextLineContent + "\n"))
			if err != nil {
				return err
			}
			nextLineContent = ""
			continue
		}
		if strings.Contains(scanner.Text(), BumpPrefix) {
			nextLineContent = strings.ReplaceAll(scanner.Text(), BumpPrefix, "")
			nextLineContent = strings.ReplaceAll(nextLineContent, "{0}", replaceValue)
		}
		_, err := dest.Write(append(scanner.Bytes(), []byte("\n")...))
		if err != nil {
			return err
		}
	}
	return nil
}

func check(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(ExitFail)
	}
}
