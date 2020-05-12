package models

import (
	"github.com/heroku/whaler-api/utils"
	"github.com/jinzhu/gorm"
)

type Workspace struct {
	DBModel
	Name          string    `json:"name"`
	Accounts      []Account `json:"accounts" gorm:"many2many:workspace_accounts;"`
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

//DEPRECATED -- REST
func FetchWorkspace(workspaceID string) map[string]interface{} {
	workspace := &Workspace{}
	err := DB().Table("workspaces").Where("id = ?", workspaceID).Preload("Accounts").First(workspace).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.Message(5001, "Workspace with the given id not found", true, map[string]interface{}{})
		} else {
			return utils.Message(5001, "Unable to fetch workspace, connection error", true, map[string]interface{}{})
		}
	}

	return utils.Message(2000, "Workspace fetched successfully", false, workspace)
}

func FetchWorkspaces(db *gorm.DB, userID int) ([]*Workspace, error) {
	user := User{}
	user.ID = userID
	workspaces := []*Workspace{}
	err := db.Model(&user).Related(&workspaces, "Workspaces").Error
	if err != nil {
		return []*Workspace{}, err
	}
	return workspaces, nil
}
