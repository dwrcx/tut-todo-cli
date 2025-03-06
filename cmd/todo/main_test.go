package main_test

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

var (
	binName  = "todo"
	fileName = ".todo.json"
)

func TestMain(m *testing.M) {
	fmt.Println("Building tool...")

	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	build := exec.Command("go", "build", "-o", binName)

	if err := build.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "cannot build tool %s: %s", binName, err)
		os.Exit(1)
	}

	fmt.Println("running tests...")
	result := m.Run()

	fmt.Println("cleaning up...")
	os.Remove(binName)
	os.Remove(fileName)

	os.Exit(result)
}

func TestTodoCLI(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	cmdPath := filepath.Join(dir, binName)

	t.Run("AddNewTaskFromArguments", func(t *testing.T) {
		task := "test task number 1"
		cmd := exec.Command(cmdPath, "-add", task)
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}

		t.Cleanup(func() { os.Remove(fileName) })
	})

	t.Run("AddNewTaskFromSTDIN", func(t *testing.T) {
		task2 := "test task number 2"
		cmd := exec.Command(cmdPath, "-add", task2)
		cmdStdIn, err := cmd.StdinPipe()
		if err != nil {
			t.Fatal(err)
		}
		io.WriteString(cmdStdIn, task2)
		cmdStdIn.Close()

		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}

		t.Cleanup(func() { os.Remove(fileName) })
	})

	t.Run("ListTasks", func(t *testing.T) {
		task1 := "list task 1"
		cmd := exec.Command(cmdPath, "-add", task1)
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}

		task2 := "list task 2"
		cmd = exec.Command(cmdPath, "-add", task2)
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}

		cmd = exec.Command(cmdPath, "-list")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		expected := fmt.Sprintf("  1: %s\n  2: %s\n", task1, task2)
		if expected != string(out) {
			t.Errorf("expected %q, got %q instead\n", expected, string(out))
		}
		t.Cleanup(func() { os.Remove(fileName) })
	})
}
