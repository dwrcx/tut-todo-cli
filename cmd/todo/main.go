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
	add := flag.Bool("add", false, "Add task to the todo list")
	list := flag.Bool("list", false, "List all tasks")
	verbose := flag.Bool("v", false, "Use with -list for verbose output")
	complete := flag.Int("complete", 0, "Item to be completed")
	del := flag.Int("delete", 0, "Item to be deleted")
	flag.Parse()

	if os.Getenv("TODO_FILENAME") != "" {
		todoFileName = os.Getenv("TODO_FILENAME")
	}

	l := &todo.List{}

	// Read the to-do items from file
	if err := l.Get(todoFileName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	switch {
	case *list:
		if *verbose {
			printTasksVerbose(*l)
		} else {
			printTasks(*l)
		}

	case *complete > 0:
		if err := l.Complete(*complete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	case *add:
		t, err := getTask(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		l.Add(t)

		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	case *del > 0:
		if *del > len(*l) || *del <= 0 {
			fmt.Fprintln(os.Stderr, "Invalid task number")
			os.Exit(1)
		}

		taskName := (*l)[*del-1].Task
		if err := l.Delete(*del); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Printf("%v was deleted.\n", taskName)
		printTasks(*l)

	default:
		fmt.Fprintln(os.Stderr, "Invalid option")
		os.Exit(1)
	}
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
