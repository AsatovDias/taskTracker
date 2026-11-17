package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

type Status string

const (
	StatusTodo       Status = "todo"
	StatusInProgress Status = "in-progress"
	StatusDone       Status = "done"
)

type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type TaskStore struct {
	FilePath string
}

func NewTaskStore(filePath string) *TaskStore {
	return &TaskStore{FilePath: filePath}
}

func (s *TaskStore) Load() ([]Task, error) {
	data, err := os.ReadFile(s.FilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []Task{}, nil
		}
		return nil, fmt.Errorf("failed to read tasks file: %w", err)
	}
	if len(data) == 0 {
		return []Task{}, nil
	}

	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("failed to parse tasks JSON: %w", err)
	}
	return tasks, nil
}

func (s *TaskStore) Save(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize tasks: %w", err)
	}
	if err := os.WriteFile(s.FilePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write tasks file: %w", err)
	}
	return nil
}

func (s *TaskStore) Add(description string) (Task, error) {
	if description == "" {
		return Task{}, errors.New("task description cannot be empty")
	}

	tasks, err := s.Load()
	if err != nil {
		return Task{}, err
	}

	nextID := 1
	for _, t := range tasks {
		if t.ID >= nextID {
			nextID = t.ID + 1
		}
	}

	now := time.Now().Truncate(time.Second)
	newTask := Task{
		ID:          nextID,
		Description: description,
		Status:      StatusTodo,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	tasks = append(tasks, newTask)
	if err := s.Save(tasks); err != nil {
		return Task{}, err
	}
	return newTask, nil
}

func (s *TaskStore) Update(id int, description string) error {
	if description == "" {
		return errors.New("task description cannot be empty")
	}

	tasks, err := s.Load()
	if err != nil {
		return err
	}

	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Description = description
			tasks[i].UpdatedAt = time.Now().Truncate(time.Second)
			return s.Save(tasks)
		}
	}
	return fmt.Errorf("task with ID %d not found", id)
}

func (s *TaskStore) Delete(id int) error {
	tasks, err := s.Load()
	if err != nil {
		return err
	}

	for i, t := range tasks {
		if t.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return s.Save(tasks)
		}
	}
	return fmt.Errorf("task with ID %d not found", id)
}

func (s *TaskStore) UpdateStatus(id int, status Status) error {
	tasks, err := s.Load()
	if err != nil {
		return err
	}

	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Status = status
			tasks[i].UpdatedAt = time.Now().Truncate(time.Second)
			return s.Save(tasks)
		}
	}
	return fmt.Errorf("task with ID %d not found", id)
}

func (s *TaskStore) List(filter Status) ([]Task, error) {
	tasks, err := s.Load()
	if err != nil {
		return nil, err
	}

	if filter == "" {
		return tasks, nil
	}

	var filtered []Task
	for _, t := range tasks {
		if t.Status == filter {
			filtered = append(filtered, t)
		}
	}
	return filtered, nil
}
