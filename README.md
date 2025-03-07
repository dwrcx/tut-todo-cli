# Todo CLI in Go

This project follows a tutorial from the book [Powerful Command Line Applications in Go](https://pragprog.com/titles/rggo/), written by Ricardo Gerardi.  
The tutorial portion ended at [this commit](https://github.com/dwrcx/tut-todo-cli/commit/8fc91c2842ab17308dbb3310227d835dab22faea)

Since then, I've extended the functionality as I realised I wanted to actually use it.

## Features
- Add, complete, undo and remove tasks
- List tasks with filtering and verbose output options
- Supports both long and short flag names (`-add` / `-a`)
- Persistent JSON storage (will explore SQLite integration soon)

## Usage

By default, the CLI creates a `.todo.json` file in the directory where it's run.  
This makes it easy to have directory-specific or project-specific todo lists.

If you prefer using a single shared todo file, set the following environment variable:
```sh
export TODO_FILENAME=path/to/file
```

```
Usage: todo [command] [options]

Options:
  -a,  -add        Add a new task
  -rm, -remove N   Remove task N
  -d,  -done N     Mark task N as done
  -ud, -undo N     Revert completed task N
  -l,  -list       List tasks
  -v,  -verbose    Show detailed task info
  -f,  -filter     Filter out completed tasks
  -clear           Remove all tasks
  -h   --help      Show this help message

Examples:
  todo -a "Water plants"
  todo -done 1
  todo -l -v -f # list verbose and filter completed tasks
```

## Learning Objectives
- Structuring a CLI in Go without a framework
- Handling user input and command-line arguments
- Managing file-based task storage (JSON)
- Exploring SQLite as an embedded database for CLI tools

## Roadmap
- Test embedding SQLite for persistent storage
- Investigate a migration path from JSON to SQLite

## Things I learned

### The Standard Go `flag` Library Has Some Limitations
- No built-in support for flag aliases (`-add`, `-a`)
- Doesn't allow "options" within a flag like `ss -tunlp`
- Combining options (`todo -lvf`) isn’t possible
- For fine-grained flag control, [cobra](https://github.com/spf13/cobra) seems like a strong alternative

### But You Can Still Build a Lot Without a Framework
- Writing an underlying API and treating the CLI as a "view" layer keeps things clean
- Manual flag handling is simple enough for smaller projects

### CLI's Make for Rewarding Projects
- I didn’t expect to enjoy building a todo CLI, but it turned out to be pretty satisfying
- The focused workflow keeps distractions minimal
- The develop/test/iterate loop is really fast
