<h1 align="center">git-utils</h1>

A command-line tool that provides various utilities for working with Git repositories.
<br><br>

## Features

- `pull`: Pull branches from multiple repositories.
- `fetch`: Fetch branches from multiple repositories.
- `grep`: Search for a pattern in file contents across multiple repositories.
- `checkout`: Checkout a branch in multiple repositories.

# Installation

## Using Go

```
go install github.com/arzkar/git-utils@latest
```

## From Github

1. Make sure you have Git installed on your system.
2. Download the latest release of Git Utils from the [Releases](https://github.com/arzkar/git-utils/releases) page.
3. Extract the downloaded archive to a location of your choice.
4. Add the extracted directory to your system's PATH.

# Usage

```
git-utils v0.2.0
Copyright (c) Arbaaz Laskar <arzkar.dev@gmail.com>

A CLI for performing various operations on git repositories

Usage:
  git-utils [command]

Available Commands:
  checkout    Checkout a branch in all repositories
  completion  Generate the autocompletion script for the specified shell
  fetch       Fetch all or specified branches
  grep        Search for a pattern in files
  help        Help about any command
  pull        Pull all or specified branches

Flags:
  -h, --help   help for git-utils

Use "git-utils [command] --help" for more information about a command.
```

# Example

## Note

You can specify the optional `--dir` flag to search within a specific directory. By default, the current directory is used.

### Pull

The `pull` command allows you to update your local branch with the latest changes from the remote repository.

Command:
git-utils pull <branch>

### Fetch

The `fetch` command fetches the latest changes from the remote repository. You can specify either a single branch or "all" to fetch all branches.

Command:
git-utils fetch <branch>

Example:
git-utils fetch feature-branch

### Checkout

The `checkout` command allows you to switch between branches in a Git repository.

Command:
git-utils checkout <branch>

Example:
git-utils checkout develop

### Grep

The `grep` command searches for a specified pattern in the files of a Git repository.

Command:
git-utils grep <pattern> [--dir=<directory>]

Example:
git-utils grep "TODO"
