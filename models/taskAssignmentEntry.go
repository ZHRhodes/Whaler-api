package models

import (
	"fmt"

	"github.com/heroku/whaler-api/graph/model"
)

type TaskAssignmentEntry struct {
	DBModel
	TaskID     string  `json:"taskId"`
	AssignedBy string  `json:"assignedBy"`
	AssignedTo *string `json:"assignedTo"`
}

func CreateTaskAssignmentEntry(newEntry model.NewTaskAssignmentEntry) (*TaskAssignmentEntry, error) {
	fmt.Printf("\nCreating task assignment entry")
	var entry = &TaskAssignmentEntry{
		TaskID:     newEntry.TaskID,
		AssignedBy: newEntry.AssignedBy,
		AssignedTo: newEntry.AssignedTo,
	}

	var err = db.Create(entry).Error

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var task Task
	db.First(&task, "id = ?", newEntry.TaskID).Association("AssignmentEntries").Append(entry)
	db.Model(&task).Update("AssignedTo", entry.AssignedTo)
	fmt.Printf("\nUpdating assigned to field for taskId %s to assignedTo %s", task.ID, *entry.AssignedTo)

	return entry, nil
}

func FetchTaskAssignmentEntries(taskID string) ([]*TaskAssignmentEntry, error) {
	entries := []*TaskAssignmentEntry{}

	var task Task
	var err = DB().Debug().First(&task, "id = ?", taskID).Error
	association := DB().Model(&task).Association("AssignmentEntries")
	association.Find(&entries)

	if err != nil {
		fmt.Println("Something bad happened here...")
		fmt.Println(err)
	}

	return entries, nil
}
