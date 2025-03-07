package main

import (
	"bytes"
	"fmt"
	todo "github.com/dwrcx/tut-todo-cli"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

var (
	binName  = "todo"
	fileName = ".todo.json"
)

// ========== Utils ==========

func runCLI(args ...string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	cmd := exec.Command(filepath.Join(dir, binName), args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func captureOutput(f func()) string {
	var buf bytes.Buffer
	stdout := os.Stdout
	defer func() { os.Stdout = stdout }()

	r, w, _ := os.Pipe()
	os.Stdout = w
	f()
	w.Close()

	buf.ReadFrom(r)
	return buf.String()
}

// ========== Helpers ========

func setupTest(t *testing.T) {
	t.Helper()
	t.Cleanup(func() { os.Remove(fileName) })
}

func addTask(t *testing.T, name string) {
	t.Helper()
	_, err := runCLI("-a", name)
	if err != nil {
		t.Fatalf("failed to add task: %v", err)
	}
}

func listTasks(t *testing.T) string {
	t.Helper()
	out, err := runCLI("-l")
	if err != nil {
		t.Fatalf("failed to list tasks: %v", err)
	}
	return out
}

// ========== Tests ==========

func TestMain(m *testing.M) {
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	build := exec.Command("go", "build", "-o", binName)
	if err := build.Run(); err != nil {
		os.Stderr.WriteString("cannot build tool: " + err.Error())
		os.Exit(1)
	}

	result := m.Run()

	os.Remove(binName)
	os.Remove(fileName)
	os.Exit(result)
}

func TestAddTask(t *testing.T) {
	setupTest(t)
	addTask(t, "test task 1")
	output := listTasks(t)
	expected := "  1: test task 1\n"

	if output != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, output)
	}
}

func TestAddTaskFromSTDIN(t *testing.T) {
	setupTest(t)

	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	task := "task from STDIN"
	cmd := exec.Command(filepath.Join(dir, binName), "-a")
	cmdStdIn, err := cmd.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}

	io.WriteString(cmdStdIn, task)
	cmdStdIn.Close()

	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}

	output := listTasks(t)
	expected := "  1: task from STDIN\n"

	if output != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, output)
	}
}

func TestListTasks(t *testing.T) {
	setupTest(t)
	addTask(t, "Task 1")
	addTask(t, "Task 2")
	output := listTasks(t)
	expected := "  1: Task 1\n  2: Task 2\n"

	if output != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, output)
	}
}

func TestDeleteTask(t *testing.T) {
	setupTest(t)
	addTask(t, "Task to delete")
	addTask(t, "Other task")

	out, err := runCLI("-rm", "1")
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(out, "Removed Task [Task to delete]") {
		t.Errorf("Expected delete message, got:\n%s", out)
	}

	output := listTasks(t)
	expected := "  1: Other task\n"

	if output != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, output)
	}
}

func TestClearTasks(t *testing.T) {
	setupTest(t)
	addTask(t, "task A")
	addTask(t, "task B")
	addTask(t, "task C")
	addTask(t, "task D")

	out, err := runCLI("-clear")
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(out, "All tasks removed") {
		t.Errorf("Expected 'All tasks removed' got:\n%s", out)
	}

	output := listTasks(t)
	expected := "No tasks found.\n"
	if output != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, output)
	}
}

func TestPrintTasks(t *testing.T) {
	tasks := todo.List{
		{Task: "task A", Done: false},
		{Task: "task B", Done: true},
	}

	exp := "  1: task A\nX 2: task B\n"

	out := captureOutput(func() {
		printTasks(tasks)
	})

	if out != exp {
		t.Errorf("expected:\n%s\ngot:\n%s", exp, out)
	}
}

func TestPrintTasksVerbose(t *testing.T) {
	created := time.Date(2021, 2, 3, 4, 5, 0, 0, time.UTC)
	completed := created.Add(2 * time.Hour)

	tasks := todo.List{
		{Task: "task A", Done: false, CreatedAt: created},
		{Task: "task B", Done: true, CreatedAt: created, CompletedAt: completed},
	}

	exp := fmt.Sprintf(
		"  1: %s %s task A\n"+
			"X 2: %s %s task B\n",
		created.Format("Mon 02/01"), created.Format("15:04"),
		completed.Format("Mon 02/01"), completed.Format("15:04"),
	)

	out := captureOutput(func() {
		printTasksVerbose(tasks)
	})

	if out != exp {
		t.Errorf("expected:\n%s\ngot:\n%s", exp, out)
	}
}
