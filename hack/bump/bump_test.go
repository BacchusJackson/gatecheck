package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestBumpReplace(t *testing.T) {
	inputString := "abc\ndef\nghi\njkl\nmno\n// bump: version {0}\nversion v1.0.0\npqr"
	t.Log(inputString)
	outputBuf := new(bytes.Buffer)
	err := bumpReplace(strings.NewReader(inputString), outputBuf, "v1.0.1")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(outputBuf.String())
}

func TestBumpChangelog(t *testing.T) {
	inputString := `# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [UNRELEASED]

### Fixed

- Bug where Gitleaks report with no secrets aren't properly decoded
	`
	outputBuf := new(bytes.Buffer)
	err := bumpChangelog(strings.NewReader(inputString), outputBuf, "v1.0.1")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(outputBuf.String())
}
