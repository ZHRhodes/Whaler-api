package models

import (
	"github.com/heroku/whaler-api/utils"
)

type Account struct {
	DBModel
	Name                string `json:"name"`
	Industry            string `json:"industry"`
	Description         string `json:"description"`
	Tier                uint   `json:"tier"`
	URL                 string `json:"url"`
	Location            string `json:"location"`
	HeadcountUpperBound uint   `json:"headcountUpperBound"`
	HeadcountLowerBound uint   `json:"headcountLowerBound"`
	RevenueUpperBound   uint   `json:"revenueUpperBound"`
	RevenueLowerBound   uint   `json:"RevenueLowerBound"`
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
