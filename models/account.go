package models

import (
	"github.com/heroku/whaler-api/graph/model"
	"github.com/heroku/whaler-api/utils"
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
	RevenueLowerBound   int    `json:"revenueLowerBound"`
	AssignedTo          []User `json:"assignedTo"`
	//notes
	//contacts
}

//DEPRECATED -- REST
func (account *Account) Create() map[string]interface{} {
	DB().Create(account)

	if account.ID <= 0 {
		return utils.Message(5001, "Failed to create account, connection error", true, map[string]interface{}{})
	}

	data := map[string]interface{}{"account": account}
	response := utils.Message(2000, "Account has been created", false, data)
	return response
}

func CreateAccount(newAccount model.NewAccount) (*Account, error) {
	account := &Account{
		Name: newAccount.Name,
		Industry: newAccount.Industry,
		Description: newAccount.Description,
		Tier: newAccount.Tier,
		URL: newAccount.URL,
		Location: newAccount.Location,
		HeadcountUpperBound: newAccount.HeadcountUpperBound,
		HeadcountLowerBound: newAccount.HeadcountLowerBound,
		RevenueUpperBound: newAccount.RevenueLowerBound,
		RevenueLowerBound: newAccount.RevenueLowerBound
	}

	err := DB().Create(account).Error

	if account.ID <= 0 {
		return nil, err
	}

	return account, nil
}
