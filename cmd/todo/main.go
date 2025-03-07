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
	done := flag.Int("d", 0, "Mark task as done")
	remove := flag.Int("rm", 0, "Remove task")
	flag.Parse()

	if os.Getenv("TODO_FILENAME") != "" {
		todoFileName = os.Getenv("TODO_FILENAME")
	}

	l := &todo.List{}
	loadTasks(l, todoFileName)

	switch {
	case *list:
		if *verbose {
			printTasksVerbose(*l)
		} else {
			printTasks(*l)
		}

	case *done > 0:
		if err := validateTaskNumber(*done, l); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		taskName := (*l)[*done-1].Task
		if err := l.Complete(*done); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		saveTasks(l, todoFileName)
		fmt.Printf("Completed Task [%s]\n", taskName)

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
		if err := validateTaskNumber(*remove, l); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		taskName := (*l)[*remove-1].Task
		if err := l.Delete(*remove); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		saveTasks(l, todoFileName)
		fmt.Printf("Removed Task [%s]\n", taskName)

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

func validateTaskNumber(n int, l *todo.List) error {
	if n > len(*l) || n <= 0 {
		fmt.Fprintf(os.Stderr, "Invalid task number. Select a number between 1 and %d\n\n", len(*l))
	}
	return nil
}

func showHelp() {
	fmt.Println(`Usage: todo [command] [options]

Commands:
  a    "task"    Add a new task
  rm   <number>  Remove a task
  d    <number>  Mark a task as done
  l              List tasks
  l -v           List tasks verbose

Options:
  -h, --help     Show this help message`)
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

// printTasks displays the to-do list in a user-friendly way.
func printTasks(l todo.List) {
	if len(l) == 0 {
		fmt.Println("No tasks found.")
		return
	}

	for i, item := range l {
		status := " "
		if item.Done {
			status = "X"
		}
		fmt.Printf("%s %d: %s\n", status, i+1, item.Task)
	}
}

// printTasksVerbose displays the to-do list in verbose mode
func printTasksVerbose(l todo.List) {
	if len(l) == 0 {
		fmt.Println("No tasks found.")
		return
	}

	for i, item := range l {
		status := " "
		if item.Done {
			status = "X"
		}

		date := item.CreatedAt.Format("Mon 02/01")
		hrMin := item.CreatedAt.Format("15:04")
		if !item.CompletedAt.IsZero() {
			date = item.CompletedAt.Format("Mon 02/01")
			hrMin = item.CompletedAt.Format("15:04")
		}

		fmt.Printf("%s %d: %s %s %s\n", status, i+1, date, hrMin, item.Task)
	}
}
