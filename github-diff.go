package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func getDiffFromURL(url string) []string {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/vnd.github.v3.diff")
	req.Header.Add("User-Agent", "github.com/barakmich/gogh")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	diff, _ := ioutil.ReadAll(resp.Body)
	return strings.Split(strings.TrimSpace(string(diff)), "\n")
}

type FileDiffMap map[string][]DiffLine

type DiffLine struct {
	Line       string
	LeftIndex  int
	RightIndex int
}

func processDiffIntoDiffMap(diffFile []string) FileDiffMap {
	out := make(FileDiffMap)
	var file string
	var dlines []DiffLine
	leftLine := 0
	rightLine := 0
	for i := 0; i < len(diffFile); {
		line := diffFile[i]
		if strings.HasPrefix(line, "diff") {
			if file != "" {
				out[file] = dlines
				dlines = nil
			}
			i += 3
			line = diffFile[i]
			file = line[6:]
		} else if strings.HasPrefix(line, "@@") {
			ls := 0
			rs := 0
			fmt.Sscanf(line, "@@ -%d,%d +%d,%d @@", &leftLine, &ls, &rightLine, &rs)
			dlines = append(dlines, DiffLine{line, ls, rs})
		} else if strings.HasPrefix(line, "+") {
			dlines = append(dlines, DiffLine{line, leftLine, rightLine})
			rightLine += 1
		} else if strings.HasPrefix(line, "-") {
			dlines = append(dlines, DiffLine{line, leftLine, rightLine})
			leftLine += 1
		} else {
			dlines = append(dlines, DiffLine{line, leftLine, rightLine})
			leftLine += 1
			rightLine += 1
		}
		i++
	}
	if file != "" {
		out[file] = dlines
		dlines = nil
	}
	return out
}
