package models

import (
	"fmt"
	"time"
)

type Task struct {
	DBModel
	AssociatedTo      *string               `json:"associatedTo"`
	Description       string                `json:"description"`
	Done              bool                  `json:"done"`
	Type              *string               `json:"type"`
	DueDate           *time.Time            `json:"dueDate"`
	AssignmentEntries []TaskAssignmentEntry `json:"assignmentEntries"`
	AssignedTo        *string               `json:"assignedTo"`
}

func CreateTask(newTask Task) (*Task, error) {
	err := DB().Create(&newTask).Error

	if err != nil {
		fmt.Println("Failed to create task.", err)
	}

	return &newTask, nil
}

func SaveTask(saveTask Task) (*Task, error) {
	err := DB().Save(&saveTask).Error

	if err != nil {
		fmt.Println("Failed to save task.", err)
	}

	return &saveTask, err
}

func FetchTasks(associatedTo string) ([]*Task, error) {
	var tasks = []*Task{}
	err := db.Where(&Task{AssociatedTo: &associatedTo}).Find(&tasks).Error

	if err != nil {
		fmt.Println("Failed to fetch tasks.", err)
	}

	return tasks, err
}
