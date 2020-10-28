package models

import (
	"github.com/heroku/whaler-api/graph/model"
	"github.com/heroku/whaler-api/utils"
	"gorm.io/gorm/clause"
)

type Account struct {
	DBModel
	Name              string                   `json:"name"`
	OwnerID           string                   `json:"ownerID"`
	SalesforceOwnerID *string                  `json:"salesforceOwnerID"`
	SalesforceID      *string                  `json:"salesforceID"`
	Industry          *string                  `json:"industry"`
	Description       *string                  `json:"description"`
	NumberOfEmployees *string                  `json:"numberOfEmployees"`
	AnnualRevenue     *string                  `json:"annualRevenue"`
	BillingCity       *string                  `json:"billingCity"`
	BillingState      *string                  `json:"billingState"`
	Phone             *string                  `json:"phone"`
	Website           *string                  `json:"website"`
	Type              *string                  `json:"type"`
	State             *string                  `json:"state"`
	Notes             *string                  `json:"notes"`
	AssignmentEntries []AccountAssignmentEntry `json:"assignmentEntries" gorm:"foreignKey:AccountID;references:ID"`
	Collaborators     []User                   `json:"collaborators" gorm:"many2many:account_collaborators;"`
	// AssignedTo          []User `json:"assignedTo"`
	//contacts
}

//DEPRECATED -- REST
func (account *Account) Create() map[string]interface{} {
	DB().Create(account)

	if len(account.ID) == 0 {
		return utils.Message(5001, "Failed to create account, connection error", true, map[string]interface{}{})
	}

	data := map[string]interface{}{"account": account}
	response := utils.Message(2000, "Account has been created", false, data)
	return response
}

func CreateAccount(newAccount model.NewAccount) (*Account, error) {
	account := createAccountFromNewAccount(newAccount)

	err := DB().Create(account).Error

	if len(account.ID) == 0 {
		return nil, err
	}

	return account, nil
}

//When saving, we need to 1. set user.OwnedAccount = newAccount and 2. set user.CollaboratingAccount = newAccounts
//In the future, we can't just assume that the accounts they're saving are all their owned accounts
func SaveAccounts(newAccounts []*model.NewAccount, userID string) ([]*Account, error) {
	var accounts = []*Account{}
	for _, newAccount := range newAccounts {
		account := createAccountFromNewAccount(*newAccount)
		if account.OwnerID == "" {
			account.OwnerID = userID
		}
		accounts = append(accounts, account)
	}

	err := DB().Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "owner_id", "industry",
			"salesforce_id", "description", "number_of_employees", "annual_revenue",
			"billing_city", "billing_state", "phone", "website", "type", "state", "notes"}),
	}).Create(&accounts).Error

	user := FetchUser(userID)

	user.OwnedAccounts = accounts
	user.CollaboratingAccounts = accounts

	db.Model(&user).Association("OwnedAccounts").Replace(accounts)
	db.Model(&user).Association("CollaboratingAccounts").Replace(accounts)

	return accounts, err
}

func FetchAccounts(userID string) ([]*Account, error) {
	var accounts = []*Account{}
	err := db.Model(&User{DBModel: DBModel{ID: userID}}).Association("CollaboratingAccounts").Find(&accounts)
	return accounts, err
}

func createAccountFromNewAccount(newAccount model.NewAccount) *Account {
	id := SafelyUnwrap(newAccount.ID)
	ownerID := SafelyUnwrap(newAccount.OwnerID)

	return &Account{
		DBModel:           DBModel{ID: id},
		Name:              newAccount.Name,
		OwnerID:           ownerID,
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

//Move this
func SafelyUnwrap(value *string) string {
	if value != nil {
		return *value
	} else {
		return ""
	}
}
