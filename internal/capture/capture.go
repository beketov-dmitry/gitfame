package capture

import (
	"log"
	"os/exec"
	"strings"
)

func MakeListOfFilesAndDirectories(revision, path string) (string, error) {
	cmd := exec.Command("git", "ls-tree", revision)
	cmd.Dir = path
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() git ls-tree with %s %s failed with %s\n", path, revision, err)
	}
	return string(out), err
}

func MakeListOfFileNames(revision, path string) (string, error) {
	cmd := exec.Command("git", "ls-tree", "--name-only", revision)
	cmd.Dir = path
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() git ls-tree --name-only %s failed with %s\n", path, err)
	}
	return string(out), err
}

func MakeListOfLastCommitsOfFile(revision, path string, dir string) (string, error) {
	cmd := exec.Command("git", "blame", "--porcelain", revision, path)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() git blame with %s %s failed with %s\n", path, dir, err)
	}
	return string(out), err
}

func MakeListOfEmptyFileChangers(revision, path string, dir string) (string, error) {
	cmd := exec.Command("git", "log", `--pretty=format:"%cn %H"`, revision, "--", path)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() git blame with %s %s failed with %s\n", path, dir, err)
	}
	lines := strings.Split(string(out), "\n")
	return lines[0], nil
}

//
//go build ./gitfame/cmd/gitfame/
//go install ./gitfame/cmd/gitfame/
