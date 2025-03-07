package main

import (
	"bufio"
	"flag"
	"fmt"
	todo "github.com/dwrcx/tut-todo-cli"
	"io"
	"os"
	"strings"
)

var todoFileName = ".todo.json"

func main() {
	add := flag.Bool("a", false, "Add a new task")
	list := flag.Bool("l", false, "List tasks")
	verbose := flag.Bool("v", false, "Use with -l for verbose list output")
	filterComplete := flag.Bool("fc", false, "Use with -l to filter complete tasks")
	done := flag.Int("d", 0, "Mark task as done")
	undone := flag.Int("ud", 0, "Revert (undo) a completed task")
	remove := flag.Int("rm", 0, "Remove task")
	clear := flag.Bool("clear", false, "Remove all tasks")
	flag.Parse()

	if os.Getenv("TODO_FILENAME") != "" {
		todoFileName = os.Getenv("TODO_FILENAME")
	}

	l := &todo.List{}
	loadTasks(l, todoFileName)

	switch {
	case *list:
		printTasks(*l, *verbose, *filterComplete)

	case *done > 0:
		modifyTask(l.Complete, *done, "Completed", l)

	case *undone > 0:
		modifyTask(l.UndoComplete, *undone, "Reverted", l)

	case *add:
		t, err := getTask(os.Stdin, flag.Args()...)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		l.Add(t)

		saveTasks(l, todoFileName)
		fmt.Printf("Added Task [%s]\n", t)

	case *remove > 0:
		modifyTask(l.Delete, *remove, "Removed", l)

	case *clear:
		l.DeleteAll()
		fmt.Printf("All tasks removed.\n")
		saveTasks(l, todoFileName)

	default:
		fmt.Fprintf(os.Stderr, "Invalid option. Use -h for help.\n\n")
		showHelp()
	}
}

func loadTasks(l *todo.List, filename string) {
	if err := l.Get(filename); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func saveTasks(l *todo.List, filename string) {
	if err := l.Save(filename); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func showHelp() {
	fmt.Println(`Usage: todo [command] [options]

Commands:
  a     "task"    Add a new task
  rm    <number>  Remove a task
  clear           Remove all tasks
  d     <number>  Mark a task as done
  ud    <number>  Revert (undo) a completed task
  l               List tasks
  l -v            List tasks verbose
  l -fc           List tasks and filter complete items

Options:
  -h, --help      Show this help message`)
}

func getTask(r io.Reader, args ...string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}

	s := bufio.NewScanner(r)
	s.Scan()

	if err := s.Err(); err != nil {
		return "", err
	}

	if len(s.Text()) == 0 {
		return "", fmt.Errorf("task cannot be blank")
	}

	return s.Text(), nil
}

func modifyTask(action func(int) error, idx int, actionLabel string, l *todo.List) {
	if idx > len(*l) || idx <= 0 {
		fmt.Fprintf(os.Stderr,
			"Invalid task number. Enter a task number between 1 and %d\n", len(*l))
		return
	}

	taskName := (*l)[idx-1].Task
	if err := action(idx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	saveTasks(l, todoFileName)
	fmt.Printf("%s Task [%s]\n", actionLabel, taskName)
}

func printTasks(l todo.List, verbose bool, filterComplete bool) {
	if len(l) == 0 {
		fmt.Println("No tasks found.")
		return
	}

	for i, item := range l {
		if filterComplete && item.Done {
			continue
		}

		status := " "
		if item.Done {
			status = "X"
		}

		if verbose {
			date := item.CreatedAt.Format("Mon 02/01")
			hrMin := item.CreatedAt.Format("15:04")
			if !item.CompletedAt.IsZero() {
				date = item.CompletedAt.Format("Mon 02/01")
				hrMin = item.CompletedAt.Format("15:04")
			}
			fmt.Printf("%s %d: %s %s %s\n", status, i+1, date, hrMin, item.Task)
		} else {
			fmt.Printf("%s %d: %s\n", status, i+1, item.Task)
		}
	}
}
