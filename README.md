<h1 align="center">git-utils</h1>


[![ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/arzkar)

A CLI for performing various operations on git repositories
<br>

## Features

- `pull`: Pull branches for all the repositories at once.
- `fetch`: Fetch branches for all the repositories at once.
- `grep`: Search for a pattern in file contents across multiple repositories.
- `checkout`: Checkout a branch for all the repositories at once.
- `tag`: Use custom tag messages for git repositories
- `bump`: Version bump the version

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
> git-utils
git-utils v0.5.1
Copyright (c) Arbaaz Laskar <arzkar.dev@gmail.com>

A CLI for performing various operations on git repositories

Usage:
  git-utils [command]

Available Commands:
  bump        Version bump the version
  checkout    Checkout a branch in all repositories
  completion  Generate the autocompletion script for the specified shell
  fetch       Fetch all or specified branches
  grep        Search for a pattern in files
  help        Help about any command
  pull        Pull all or specified branches
  tag         Create a new tag with a custom message for the repository

Flags:
      --config   Show app config
  -h, --help     help for git-utils

Use "git-utils [command] --help" for more information about a command.
```

# Example

## Note

- You can specify the optional `--dir` flag to search within a specific directory. By default, the current directory is used.

- Use the `--config` flag to show the location of the app directory and config file

- For pull & fetch commands, you can specify a single branch, a comma seperated list of branches or all.

### Pull

The `pull` command allows you to update your local branch with the latest changes for all the repositories at once.

Command:
`git-utils pull <branch> [--dir=<directory>]`

Example:
`git-utils pull main,devel`

### Fetch

The `fetch` command fetches the latest changes for all the repositories at once.

Command:
`git-utils fetch <branch> [--dir=<directory>]`

Example:
`git-utils fetch feature-branch`

### Checkout

The `checkout` command allows you to switch between branches for all the repositories at once.

Command:
`git-utils checkout <branch> [--dir=<directory>]`

Example:
`git-utils checkout develop`

### Grep

The `grep` command searches for a specified pattern in the files of all the repositories at once.

Command:
`git-utils grep <pattern> [--dir=<directory>]`

Example:
`git-utils grep "TODO"`

### Tag

The `tag` command reads the config file and uses custom message for tag

Command:
`git-utils tag -a <tag_name> -m <custom_message_keyword>`

Example:
`git-utils tag -a "v0.1.2" -m "changelog"`

Sample `config.json`:

```json
{
  "tags": {
    "messages": {
      "changelog": "Full changelog: https://github.com/{repo_owner}/{repo_name}/compare/{prevTag}...{newTag}"
    }
  }
}
```

Tag Templates variables available: `repo_owner, repo_name, prevTag, newTag`

### Bump

The `bump` command bumps the version set in the `.git-utils-bump.cfg` and search & replaces the version for the files set in the config file using the subcommands: `major`, `minor` & `patch`

Command:
`git-utils bump <command>`

Example:
`git-utils bump minor`

Sample `.git-utils.bump.cfg:

```yml
[bumpversion]
current_version = 0.4.0
commit          = True
tag             = True
tag_format      = v{tag}

[bumpversion:file:cmd/root.go]
search  = git-utils v{current_version}
replace = git-utils v{new_version}

[bumpversion:file:README.md]
search  = git-utils v{current_version}
replace = git-utils v{new_version}

[bumpversion:file:utils/update_checker.go]
search  = v{current_version}
replace = v{new_version}
```
