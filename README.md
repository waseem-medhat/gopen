# Gopen

![preview](https://i.imgur.com/an39lik.gif)

**A simple CLI to quick-start coding projects**

The premise of this command-line utility is to save an editor of choice and a
list of aliases for your local development projects instead of "polluting" your
system-level configs (e.g., `.bashrc`). Then, Gopen command will `cd` into that
folder and open your editor of choice.

You can also save a git repo so that when you run the Gopen command and path
doesn't exist (usually happens after a fresh OS installation), it clones the
configured repo before opening the project.

Gopen can be used either via the interactive TUI or the command-line API.

## Installation

### Precompiled binaries

Check the [releases](https://github.com/waseem-medhat/gopen/releases) for Linux,
MacOS, or Windows binaries.

### From source

You need to have [Go](https://go.dev/doc/install) installed. The installation
is simply done by running this command:

```bash
go install github.com/waseem-medhat/gopen@latest
```

Alternatively, you can clone the repo, `cd` into its folder, and run:

```bash
go install
```

In any case, this should build the the `gopen` binary and install it in the
directory specified by your `GOBIN` environment variable (default is
`~/go/bin`). You might need to add that directory to your `PATH` variable if it
isn't there by default.

## Usage

For the interactive TUI, simply run `gopen` in your terminal.

### Config File

Your editor command and directory aliases will be stored in
`~/.config/gopen/gopen.json`, which can be initially created with the init
option, or its shorthand `i`. Both the file and the directory will be created
if they don't exist.

```bash
gopen i
# Creating a new config file...

gopen i
# Found config file - exiting...
```

### Editor Command

The `editor` option, or its shorthand `e`, allows you to get or set your editor
command. Using it with no additional command-line arguments will get the
current editor command. Adding the command (or the path to an executable
binary) as an argument will set it in the config.

```bash
gopen e vi

gopen e
# vi
```

### Directory Aliases

The `alias` option, or its shorthand `a`, allows you to list the aliases, get
the path assigned to a specific alias, or set a new one.

```bash
# list all aliases
gopen a

# get the path assigned to a specifc alias
gopen a myproj

# set a new alias
gopen a myproj path/to/my-proj
```

You can remove aliases using `remove` or its shorthand `r`.

```bash
gopen remove myproj
```

### Execution

Once you have your editor and aliases configured, simply provide the alias to
the `gopen` command. It will cd into the assigned path and open your editor.

```bash
gopen myproj
```

## Contributing

Any contributions are welcome! Feel free to raise issues for bug
reports/feature requests or open pull requests.

For PRs, please ensure the tests and linting rules pass before submitting the
PR, and add your own tests if applicable.
