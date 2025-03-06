package main_test

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
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

// ========== Helpers ========

func setupTest(t *testing.T) {
	t.Helper()
	t.Cleanup(func() { os.Remove(fileName) })
}

func addTask(t *testing.T, name string) {
	t.Helper()
	_, err := runCLI("-add", name)
	if err != nil {
		t.Fatalf("failed to add task: %v", err)
	}
}

func listTasks(t *testing.T) string {
	t.Helper()
	out, err := runCLI("-list")
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
	cmd := exec.Command(filepath.Join(dir, binName), "-add")
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

	out, err := runCLI("-delete", "1")
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(out, "Task to delete was deleted") {
		t.Errorf("Expected delete message, got:\n%s", out)
	}

	output := listTasks(t)
	expected := "  1: Other task\n"

	if output != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, output)
	}
}
