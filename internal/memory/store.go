package memory

import (
	"errors"
	"sync"
)

type Store struct {
	mu      sync.RWMutex
	tasks   map[string]*Task
	results map[string]interface{} // Can store AgentResponse or ErrorResponse
}

func NewStore() *Store {
	return &Store{
		tasks:   make(map[string]*Task),
		results: make(map[string]interface{}),
	}
}

func (s *Store) SaveTask(task *Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if task.ID == "" {
		return errors.New("task ID cannot be empty")
	}
	s.tasks[task.ID] = task
	return nil
}

func (s *Store) GetTask(id string) (*Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, exists := s.tasks[id]
	if !exists {
		return nil, errors.New("task not found")
	}
	return task, nil
}

func (s *Store) SaveResult(taskID string, result interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.results[taskID] = result
	return nil
}

func (s *Store) GetResult(taskID string) (interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res, exists := s.results[taskID]
	if !exists {
		return nil, errors.New("result not found")
	}
	return res, nil
}
