package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
)

const defaultTasksFile = "tasks.json"

func main() {
	store := NewTaskStore(defaultTasksFile)
	if err := runCLI(os.Args[1:], store, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runCLI(args []string, store *TaskStore, out io.Writer) error {
	if len(args) == 0 {
		printUsage(out)
		return nil
	}

	cmd := args[0]
	switch cmd {
	case "add":
		if len(args) < 2 {
			return fmt.Errorf("usage: task-cli add <description>")
		}
		desc := strings.TrimSpace(strings.Join(args[1:], " "))
		task, err := store.Add(desc)
		if err != nil {
			return err
		}
		fmt.Fprintf(out, "Task added successfully (ID: %d)\n", task.ID)

	case "update":
		if len(args) < 3 {
			return fmt.Errorf("usage: task-cli update <id> <description>")
		}
		id, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid task ID %q: must be an integer", args[1])
		}
		desc := strings.TrimSpace(strings.Join(args[2:], " "))
		if err := store.Update(id, desc); err != nil {
			return err
		}
		fmt.Fprintf(out, "Task %d updated successfully\n", id)

	case "delete":
		if len(args) < 2 {
			return fmt.Errorf("usage: task-cli delete <id>")
		}
		id, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid task ID %q: must be an integer", args[1])
		}
		if err := store.Delete(id); err != nil {
			return err
		}
		fmt.Fprintf(out, "Task %d deleted successfully\n", id)

	case "mark-in-progress":
		if len(args) < 2 {
			return fmt.Errorf("usage: task-cli mark-in-progress <id>")
		}
		id, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid task ID %q: must be an integer", args[1])
		}
		if err := store.UpdateStatus(id, StatusInProgress); err != nil {
			return err
		}
		fmt.Fprintf(out, "Task %d marked as in-progress\n", id)

	case "mark-done":
		if len(args) < 2 {
			return fmt.Errorf("usage: task-cli mark-done <id>")
		}
		id, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid task ID %q: must be an integer", args[1])
		}
		if err := store.UpdateStatus(id, StatusDone); err != nil {
			return err
		}
		fmt.Fprintf(out, "Task %d marked as done\n", id)

	case "mark-todo":
		if len(args) < 2 {
			return fmt.Errorf("usage: task-cli mark-todo <id>")
		}
		id, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid task ID %q: must be an integer", args[1])
		}
		if err := store.UpdateStatus(id, StatusTodo); err != nil {
			return err
		}
		fmt.Fprintf(out, "Task %d marked as todo\n", id)

	case "list":
		var filter Status
		if len(args) > 1 {
			switch Status(args[1]) {
			case StatusTodo, StatusInProgress, StatusDone:
				filter = Status(args[1])
			default:
				return fmt.Errorf("invalid status filter %q (valid: todo, in-progress, done)", args[1])
			}
		}

		tasks, err := store.List(filter)
		if err != nil {
			return err
		}
		if len(tasks) == 0 {
			fmt.Fprintln(out, "No tasks found.")
			return nil
		}

		w := tabwriter.NewWriter(out, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "ID\tSTATUS\tDESCRIPTION\tCREATED\tUPDATED")
		for _, t := range tasks {
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n",
				t.ID,
				t.Status,
				t.Description,
				t.CreatedAt.Format("2006-01-02 15:04:05"),
				t.UpdatedAt.Format("2006-01-02 15:04:05"),
			)
		}
		return w.Flush()

	case "help", "-h", "--help":
		printUsage(out)

	default:
		return fmt.Errorf("unknown command %q (run 'task-cli help' for usage)", cmd)
	}

	return nil
}

func printUsage(out io.Writer) {
	fmt.Fprintln(out, `Task Tracker CLI - Manage your tasks from the terminal

Usage:
  task-cli <command> [arguments]

Commands:
  add <description>            Add a new task
  update <id> <description>    Update an existing task's description
  delete <id>                  Delete a task by ID
  mark-in-progress <id>        Mark a task as in-progress
  mark-done <id>               Mark a task as done
  mark-todo <id>               Mark a task as todo
  list [status]                List all tasks or filter by status (todo, in-progress, done)
  help                         Show this help message`)
}
