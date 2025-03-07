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
		if *complete > len(*l) || *complete <= 0 {
			fmt.Fprintf(os.Stderr, "Invalid task number. Select a number between 1 and %d\n\n", len(*l))
			return
		}

		taskName := (*l)[*complete-1].Task
		if err := l.Complete(*complete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Printf("Completed Task [%s]\n", taskName)

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

		fmt.Printf("Added Task [%s]\n", t)

	case *del > 0:
		if *del > len(*l) || *del <= 0 {
			fmt.Fprintf(os.Stderr, "Invalid task number. Select a number between 1 and %d\n\n", len(*l))
			return
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

		fmt.Printf("Deleted Task [%s]\n", taskName)

	default:
		fmt.Fprintf(os.Stderr, "Invalid option. Use -h for help.\n\n")
		showHelp()
	}
}

func showHelp() {
	fmt.Println(`Usage: todo [command] [options]

Commands:
  add      "task"   Add a new task
  delete   <number> Delete a task
  complete <number> Complete a task
  list              List tasks
  list -v           List tasks verbose

Options:
  -h, --help        Show this help message`)
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
