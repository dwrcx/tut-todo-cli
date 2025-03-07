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
	add := flag.Bool("add", false, "Add a new task")
	a := flag.Bool("a", false, "Alias for -add")

	remove := flag.Int("remove", 0, "Remove task")
	rm := flag.Int("rm", 0, "Alias for -remove")

	done := flag.Int("done", 0, "Mark task as done")
	d := flag.Int("d", 0, "Alias for -done")

	undo := flag.Int("undo", 0, "Revert (undo) a completed task")
	ud := flag.Int("ud", 0, "Alias for -undo")

	list := flag.Bool("list", false, "List tasks")
	ls := flag.Bool("l", false, "Alias for -list")

	verbose := flag.Bool("verbose", false, "Use with -list for verbose list output")
	v := flag.Bool("v", false, "Alias for -verbose")

	filter := flag.Bool("filter", false, "Use with -list to filter complete tasks")
	f := flag.Bool("f", false, "Alias for -filter")

	clear := flag.Bool("clear", false, "Remove all tasks")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"Usage: todo [command] [options]\n\n"+
				"Options:\n"+
				"  -a,  -add        Add a new task\n"+
				"  -rm, -remove N   Remove task N\n"+
				"  -d,  -done N     Mark task N as done\n"+
				"  -ud, -undo N     Revert completed task N\n"+
				"  -l,  -list       List tasks\n"+
				"  -v,  -verbose    Show detailed task info\n"+
				"  -f,  -filter     Filter out completed tasks\n"+
				"  -clear           Remove all tasks\n"+
				"  -h   --help      Show this help message\n\n"+
				"Examples:\n"+
				"  todo -a Water plants\n"+
				"  todo -done 1\n"+
				"  todo -l -v -f (list verbose and filter completed items)")
	}

	flag.Parse()

	if os.Getenv("TODO_FILENAME") != "" {
		todoFileName = os.Getenv("TODO_FILENAME")
	}

	l := &todo.List{}
	loadTasks(l, todoFileName)

	switch {
	case *list || *ls:
		printTasks(*l, *verbose || *v, *filter || *f)

	case *done > 0 || *d > 0:
		modifyTask(l.Complete, getHighVal(*done, *d), "Completed", l)

	case *undo > 0 || *ud > 0:
		modifyTask(l.UndoComplete, getHighVal(*undo, *ud), "Reverted", l)

	case *add || *a:
		t, err := getTask(os.Stdin, flag.Args()...)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		l.Add(t)

		saveTasks(l, todoFileName)
		fmt.Printf("Added Task [%s]\n", t)

	case *remove > 0 || *rm > 0:
		modifyTask(l.Delete, getHighVal(*remove, *rm), "Removed", l)

	case *clear:
		l.DeleteAll()
		fmt.Printf("All tasks removed.\n")
		saveTasks(l, todoFileName)

	default:
		fmt.Fprintf(os.Stderr, "Invalid option. Use -h for help.\n\n")
		flag.Usage()
	}
}

func getHighVal(num1, num2 int) int {
	if num1 != 0 {
		return num1
	}
	return num2
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
