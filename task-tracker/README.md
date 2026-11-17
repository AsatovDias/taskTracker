# Task Tracker CLI (Go)

A fast, zero-dependency command-line interface (CLI) built in Go to track and manage your tasks. Tasks are persisted locally in a `tasks.json` file in the current working directory.

## Build & Run

```powershell
# Build the binary
go build -o task-cli.exe .

# Add tasks
.\task-cli.exe add "Buy groceries"
.\task-cli.exe add "Finish Go project"

# Update a task
.\task-cli.exe update 1 "Buy groceries and cook dinner"

# Mark task statuses
.\task-cli.exe mark-in-progress 1
.\task-cli.exe mark-done 2

# List tasks
.\task-cli.exe list
.\task-cli.exe list done
.\task-cli.exe list todo
.\task-cli.exe list in-progress

# Delete a task
.\task-cli.exe delete 1
```

## Running Tests

```powershell
go test -v ./...
```
