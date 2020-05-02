package models

import (
	"github.com/heroku/whaler-api/utils"
	"github.com/jinzhu/gorm"
)

type Account struct {
	DBModel
	Name                string `json:"name"`
	Industry            string `json:"industry"`
	Description         string `json:"description"`
	Tier                int    `json:"tier"`
	URL                 string `json:"url"`
	Location            string `json:"location"`
	HeadcountUpperBound int    `json:"headcountUpperBound"`
	HeadcountLowerBound int    `json:"headcountLowerBound"`
	RevenueUpperBound   int    `json:"revenueUpperBound"`
	RevenueLowerBound   int    `json:"RevenueLowerBound"`
	AssignedTo          []User `json:"assignedTo"`
	//notes
	//contacts
}

func (account *Account) Create() map[string]interface{} {
	DB().Create(account)

	if account.ID <= 0 {
		return utils.Message(5001, "Failed to create account, connection error", true, map[string]interface{}{})
	}

	data := map[string]interface{}{"account": account}
	response := utils.Message(2000, "Account has been created", false, data)
	return response
}

func FetchAccounts(workspaceID string) map[string]interface{} {
	workspace := &Workspace{}
	err := DB().Table("workspaces").Where("id = ?", workspaceID).First(workspace).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.Message(5001, "Workspace with the given id not found", true, map[string]interface{}{})
		} else {
			return utils.Message(5001, "Unable to fetch workspace, connection error", true, map[string]interface{}{})
		}
	}
	accounts := []Account{}
	DB().Model(&workspace).Related(&accounts).Find(&workspace.Accounts)

	return utils.Message(2000, "Workspace accounts fetched successfully", false, accounts)
}
