package main

import (
	"fmt"
	_ "github.com/google/go-github/github"
	"os/exec"
)

func getSHA() (string, error) {
	commitCmd := exec.Command("git", "log --pretty='%H' HEAD -1")
	shaBytes, err := commitCmd.Output()
	if err != nil {
		return "", err
	}
	return string(shaBytes), nil
}

func main() {
	sha, _ := getSHA()
	fmt.Println(sha)
	//client := github.NewClient(nil)
	//client.PullRequests.List()
}
