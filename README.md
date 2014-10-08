# gogh -- A Go Github Editor Tool
---

A simple, evolving tool meant to scratch a personal itch -- fetching Github comments for my current pull request based on my current branch, and pulling the locations into my editor.

## Installation
---
```
go get github.com/barakmich/gogh
```
As per usual

## Vim Configuration
---
I'm open to more plugins, but a simple enough

```
command Comments cexpr system('gogh')
```

Does the trick in your `.vimrc`, as long as your PWD is somewhere within the repository. Expects a remote named `upstream`, and you can either create that in git using

```
$ git remote add upstream <github path>
```

or you can set the flag for the command

```
command Comments cexpr system('gogh --upstream=my-upstream')
```

Then running `:Comments` or binding to that to the right keypress will load the valid comments for the PR for your branch into your QuickFix.

## License
---
Simplified BSD. Go nuts.
