```
┌──────────────────────────────────────────────────────┐
│                                                      │
│                                                      │
│      ██████╗  ██████╗ ██████╗ ███████╗███╗   ██╗     │
│     ██╔════╝ ██╔═══██╗██╔══██╗██╔════╝████╗  ██║     │
│     ██║  ███╗██║   ██║██████╔╝█████╗  ██╔██╗ ██║     │
│     ██║   ██║██║   ██║██╔═══╝ ██╔══╝  ██║╚██╗██║     │
│     ╚██████╔╝╚██████╔╝██║     ███████╗██║ ╚████║     │
│      ╚═════╝  ╚═════╝ ╚═╝     ╚══════╝╚═╝  ╚═══╝     │
│                                                      │
│                                                      │
└──────────────────────────────────────────────────────┘
```

**A simple CLI to quick-start coding projects**

The premise of this command-line utility is to save an editor of choice and a
list of aliases for your local development projects instead of "polluting" your
system-level configs (e.g., `.bashrc`). Then, Gopen command will `cd` into that
folder and open your editor of choice.

## Installation

- Check the [releases](https://github.com/wipdev-tech/gopen/releases) for the
  binary (only Linux for now).
- To install from source, you need to have Go installed, clone this repo and
  run this command:

```bash
go install
```

This should build the the `gopen` binary and install it in the directory
specified by your `GOBIN` environment variable.

## Usage

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
