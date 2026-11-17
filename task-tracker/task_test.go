package main

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

func TestTaskLifecycle(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "tasks.json")
	store := NewTaskStore(tmpFile)

	// 1. Add tasks
	t1, err := store.Add("Buy groceries")
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}
	if t1.ID != 1 || t1.Status != StatusTodo {
		t.Errorf("unexpected task 1: %+v", t1)
	}

	t2, err := store.Add("Write Go CLI")
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}
	if t2.ID != 2 {
		t.Errorf("expected ID 2, got %d", t2.ID)
	}

	// 2. Update task
	if err := store.Update(1, "Buy groceries and cook dinner"); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// 3. Mark status
	if err := store.UpdateStatus(1, StatusInProgress); err != nil {
		t.Fatalf("UpdateStatus failed: %v", err)
	}
	if err := store.UpdateStatus(2, StatusDone); err != nil {
		t.Fatalf("UpdateStatus failed: %v", err)
	}

	// 4. List filters
	inProg, err := store.List(StatusInProgress)
	if err != nil || len(inProg) != 1 || inProg[0].Description != "Buy groceries and cook dinner" {
		t.Errorf("unexpected in-progress list: %+v, err: %v", inProg, err)
	}

	done, err := store.List(StatusDone)
	if err != nil || len(done) != 1 || done[0].ID != 2 {
		t.Errorf("unexpected done list: %+v, err: %v", done, err)
	}

	// 5. Delete task
	if err := store.Delete(1); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	all, err := store.List("")
	if err != nil || len(all) != 1 {
		t.Errorf("expected 1 remaining task, got %d", len(all))
	}
}

func TestRunCLI(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "tasks.json")
	store := NewTaskStore(tmpFile)
	var out bytes.Buffer

	if err := runCLI([]string{"add", "Learn Go"}, store, &out); err != nil {
		t.Fatalf("CLI add error: %v", err)
	}
	if !strings.Contains(out.String(), "Task added successfully (ID: 1)") {
		t.Errorf("unexpected output: %s", out.String())
	}

	out.Reset()
	if err := runCLI([]string{"list"}, store, &out); err != nil {
		t.Fatalf("CLI list error: %v", err)
	}
	if !strings.Contains(out.String(), "Learn Go") {
		t.Errorf("expected list output to contain 'Learn Go', got: %s", out.String())
	}
}
