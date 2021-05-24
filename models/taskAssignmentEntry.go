package models

type TaskAssignmentEntry struct {
	DBModel
	TaskID     string `json:"taskId"`
	AssignedBy string `json:"assignedBy"`
	AssignedTo string `json:"assignedTo"`
}
