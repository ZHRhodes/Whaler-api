package models

import (
	"github.com/heroku/whaler-api/graph/model"
	"github.com/heroku/whaler-api/utils"
	"gorm.io/gorm/clause"
)

type Account struct {
	DBModel
	Name              string  `json:"name"`
	Owner             string  `json:"owner"`
	SalesforceID      *string `json:"salesforceID"`
	Industry          *string `json:"industry"`
	Description       *string `json:"description"`
	NumberOfEmployees *string `json:"numberOfEmployees"`
	AnnualRevenue     *string `json:"annualRevenue"`
	BillingCity       *string `json:"billingCity"`
	BillingState      *string `json:"billingState"`
	Phone             *string `json:"phone"`
	Website           *string `json:"website"`
	Type              *string `json:"type"`
	State             *string `json:"state"`
	Notes             *string `json:"notes"`
	// AssignedTo          []User `json:"assignedTo"`
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
	account := createAccountFromNewAccount(newAccount)

	err := DB().Create(account).Error

	if account.ID <= 0 {
		return nil, err
	}

	return account, nil
}

func SaveAccounts(newAccounts []*model.NewAccount) ([]*Account, error) {
	var accounts = []*Account{}
	for _, newAccount := range newAccounts {
		accounts = append(accounts, createAccountFromNewAccount(*newAccount))
	}

	err := DB().Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "owner", "industry",
			"salesforce_id", "industry", "description", "number_of_employees", "annual_revenue",
			"billing_city", "billing_state", "phone", "website", "type", "state", "notes"}),
	}).Create(&accounts).Error

	return accounts, err
}

func createAccountFromNewAccount(newAccount model.NewAccount) *Account {
	var id int
	if newAccount.ID != nil {
		id = *newAccount.ID
	}
	return &Account{
		DBModel:           DBModel{ID: id},
		Name:              newAccount.Name,
		Owner:             newAccount.Owner,
		SalesforceID:      newAccount.SalesforceID,
		Industry:          newAccount.Industry,
		Description:       newAccount.Description,
		NumberOfEmployees: newAccount.NumberOfEmployees,
		AnnualRevenue:     newAccount.AnnualRevenue,
		BillingCity:       newAccount.BillingCity,
		BillingState:      newAccount.BillingState,
		Phone:             newAccount.Phone,
		Website:           newAccount.Website,
		Type:              newAccount.Type,
		State:             newAccount.State,
		Notes:             newAccount.Notes,
	}
}
