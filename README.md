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

## Usage

### Config File

Your editor command and directory aliases will be stored in
`~/.config/gopen/gopen.json`, which can be initially created with the --init
option, or its shorthand `-i`. Both the file and the directory will be created
if they don't exist.

```bash
gopen -i
# Creating a new config file...

gopen -i
# Found config file - exiting...
```

### Editor Command

The `--editor-cmd` option, or its shorthand `-e`, allows you to get or set your
editor command. Using it with no additional command-line arguments will get the
current editor command. Adding the command (or the path to an executable
binary) as an argument will set it in the config.

```bash
gopen -e vi

gopen -e
# vi
```

