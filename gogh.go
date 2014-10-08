package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/google/go-github/github"
)

var (
	user     *string
	upstream *string
	debug    *bool
	outdated *bool
)

func init() {
	user = flag.String("user", os.Getenv("USER"), "Github username")
	upstream = flag.String("upstream", "upstream", "Github user or org for upstream fork")
	debug = flag.Bool("debug", false, "Enable debug logging")
	outdated = flag.Bool("outdated", false, "Show outdated comments")
	flag.Parse()
}

func main() {
	remotes, _ := getRemotes()
	shas, _ := getSHA()
	if *debug {
		for _, sha := range shas {
			log.Println("SHA:", sha)
		}
	}

	var remote *Remote
	for _, r := range remotes {
		if r.LocalName == *upstream {
			remote = r
			break
		}
	}
	if remote == nil {
		log.Fatalln("Couldn't find upstream remote named", *upstream)
	}
	client := github.NewClient(nil)
	data, _, _ := client.PullRequests.List(remote.User, remote.Repo, nil)
	for _, pr := range data {
		for _, sha := range shas {
			if *pr.Head.SHA == sha {
				if *debug {
					log.Println("Found the PR:", *pr.URL)
				}
				getCommentsForPullRequest(client, remote, pr)
				os.Exit(0)
			}
		}
	}
	fmt.Println("Couldn't find Pull Request matching this branch")
	os.Exit(1)
}

func getCommentsForPullRequest(client *github.Client, remote *Remote, pr github.PullRequest) {
	comments, _, _ := client.PullRequests.ListComments(remote.User, remote.Repo, *pr.Number, nil)
	diff := getDiffFromURL(*pr.URL)
	diffMap := processDiffIntoDiffMap(diff)
	basePath, _ := getTopLevelPath()
	for _, comment := range comments {
		p := basePath
		if comment.Path != nil {
			p = path.Join(basePath, *comment.Path)
		}
		pos := 0
		if comment.Position != nil && comment.Path != nil {
			diffLine := diffMap[*comment.Path][*comment.Position]
			fmt.Println(diffLine.Line)
			pos = diffLine.RightIndex
		} else if !*outdated {
			continue
		}
		fmt.Printf("%s:%d:%s (@%s)\n", p, pos, *comment.Body, *comment.User.Login)
	}
}
