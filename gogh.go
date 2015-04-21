package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var (
	user          *string
	upstream      *string
	debug         *bool
	outdated      *bool
	issueComments *bool
)

func init() {
	user = flag.String("user", os.Getenv("USER"), "Github username")
	upstream = flag.String("upstream", "origin", "Github user or org for upstream fork")
	debug = flag.Bool("debug", false, "Enable debug logging")
	outdated = flag.Bool("outdated", false, "Show outdated comments")
	issueComments = flag.Bool("issue_comments", false, "Show issue comments")
	flag.Parse()
}

type tokenSource struct {
	token *oauth2.Token
}

func (t *tokenSource) Token() (*oauth2.Token, error) {
	return t.token, nil
}

func main() {
	remotes, _ := getRemotes()
	// shas, _ := getSHA()
	branch, _ := getBranch()
	if *debug {
		log.Println("branch", branch)
		//for _, sha := range shas {
		//log.Println("SHA:", sha)
		//}
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

	keyfile, err := os.Open(os.ExpandEnv("$HOME/.github.key"))
	if err != nil {
		log.Fatalln(err)
	}
	r := bufio.NewReader(keyfile)
	key, err := r.ReadString('\n')
	if err != nil {
		log.Fatalln(err)
	}

	ts := &tokenSource{
		&oauth2.Token{AccessToken: key},
	}
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)
	data, _, _ := client.PullRequests.List(remote.User, remote.Repo, nil)
	for _, pr := range data {
		if *pr.Head.Ref == branch {
			if *debug {
				log.Println("Found the PR:", *pr.URL)
			}
			getCommentsForPullRequest(client, remote, pr)
			os.Exit(0)
		}
	}
	fmt.Println("Couldn't find Pull Request matching this branch")
	os.Exit(1)
}

func getCommentsForPullRequest(client *github.Client, remote *Remote, pr github.PullRequest) {
	comments, _, err := client.PullRequests.ListComments(remote.User, remote.Repo, *pr.Number, nil)
	if err != nil {
		log.Fatalln(err)
	}
	diff := getDiffFromURL(*pr.URL)
	diffMap := processDiffIntoDiffMap(diff)
	basePath, _ := getTopLevelPath()
	if *debug {
		log.Println("# Comments:", len(comments))
	}
	for _, comment := range comments {
		p := basePath
		if comment.Path != nil {
			p = path.Join(basePath, *comment.Path)
		}
		pos := 1
		if comment.Position != nil && comment.Path != nil {
			if *debug {
				fmt.Println(*comment.Path, *comment.Position)
				fmt.Println(diffMap[*comment.Path])
			}
			if len(diffMap[*comment.Path]) == 0 {
				continue
			}
			diffLine := diffMap[*comment.Path][*comment.Position]
			fmt.Println(diffLine.Line)
			pos = diffLine.RightIndex
		} else if !*outdated {
			continue
		}
		fmt.Printf("%s:%d:%s (@%s)\n", p, pos, *comment.Body, *comment.User.Login)
	}
	if *issueComments {
		commit, _, err := client.Git.GetCommit(remote.User, remote.Repo, *pr.Head.SHA)
		if err != nil {
			log.Fatalln(err)
		}
		prComments, _, err := client.Issues.ListComments(remote.User, remote.Repo, *pr.Number, nil)
		if err != nil {
			log.Fatalln(err)
		}
		for _, comment := range prComments {
			if commit.Author.Date.Before(*comment.UpdatedAt) {
				body := strings.Replace(*comment.Body, "\r\n", "\n* ", -1)
				fmt.Printf(":: (@%s) %s\n\n", *comment.User.Login, body)
			}
		}
	}
}
