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

func SaveTask(saveTask Task) (*Task, error) {
	var err error
	if saveTask.ID == "" {
		err = DB().Create(&saveTask).Error
	} else if saveTask.DeletedAt != nil {
		err = DB().Delete(&saveTask).Error
	} else {
		err = DB().Save(&saveTask).Error
	}

	if err != nil {
		fmt.Println("Failed to save task.", err)
	}

	return &saveTask, err
}

func FetchTasks(associatedTo string) ([]*Task, error) {
	var tasks = []*Task{}
	err := db.Where("associated_to = ?", associatedTo).Find(&tasks).Error

	if err != nil {
		fmt.Println("Failed to fetch tasks.", err)
	}

	return tasks, err
}
