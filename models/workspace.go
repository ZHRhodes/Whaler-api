package models

import "github.com/heroku/whaler-api/utils"

type Workspace struct {
	DBModel
	Name          string    `json:"name"`
	Accounts      []Account `json:"accounts"`
	Collaborators []User    `json:"collaborators"`
}

func (workspace *Workspace) Create() map[string]interface{} {
	//check if workspace already exists?
	DB().Create(workspace)

	if workspace.ID <= 0 {
		return utils.Message(5001, "Failed to create workspace, connection error", true, map[string]interface{}{})
	}

	data := map[string]interface{}{"workspace": workspace}
	response := utils.Message(2000, "Workspace has been created", false, data)
	return response
}
