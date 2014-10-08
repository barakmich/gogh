package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

func getTopLevelPath() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	bytes, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(bytes)), nil
}

func getSHA() ([]string, error) {
	cmd := exec.Command("git", "log", "--pretty=%H", "HEAD", "-10")
	bytes, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	shas := strings.Split(strings.TrimSpace(string(bytes)), "\n")
	return shas, nil
}

type Remote struct {
	LocalName string
	User      string
	Repo      string
}

func getRemotes() ([]*Remote, error) {
	var out []*Remote
	cmd := exec.Command("git", "remote", "-v")
	bytes, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(bytes)), "\n")
	re, err := regexp.Compile("github.com[:/](.*?)/(.*?).git")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	for _, line := range lines {
		record := strings.Split(line, "\t")
		lname := record[0]
		record = strings.Split(record[1], " ")
		path := record[0]
		fetchpush := record[1]
		if fetchpush != "(fetch)" {
			continue
		}
		submatch := re.FindStringSubmatch(path)
		if submatch == nil {
			continue
		}
		out = append(out, &Remote{
			LocalName: lname,
			User:      submatch[1],
			Repo:      submatch[2],
		})

	}
	return out, nil
}
