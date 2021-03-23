package models

import (
	"github.com/heroku/whaler-api/graph/model"
	"github.com/heroku/whaler-api/utils"
	"gorm.io/gorm/clause"
)

type Account struct {
	DBModel
	Name              string                   `json:"name"`
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
	Trackers          []*User                  `json:"trackers" gorm:"many2many:account_trackers;"`
	Collaborators     []User                   `json:"collaborators" gorm:"many2many:account_collaborators;"`
	Contacts          []*Contact               `json:"contacts"`
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

func SaveAccounts(newAccounts []*model.NewAccount, userID string) ([]*Account, error) {
	var savedAccounts = []*Account{}
	var error error
	for _, newAccount := range newAccounts {
		account := createAccountFromNewAccount(*newAccount)
		savedAccount, err := SaveAccount(account)
		savedAccounts = append(savedAccounts, savedAccount)
		if error == nil {
			error = err
		}
	}
	//will only return latest error
	return savedAccounts, error
}

func SaveAccount(account *Account) (*Account, error) {
	err := DB().Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "salesforce_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"updated_at", "name", "owner_id", "industry",
			"salesforce_id", "description", "number_of_employees", "annual_revenue",
			"billing_city", "billing_state", "phone", "website", "type", "state"}),
	}).Create(&account).Error

	return account, err
}

func FetchAccounts(userID string) ([]*Account, error) {
	var accounts = []*Account{}
	// err := db.Model(&User{DBModel: DBModel{ID: userID}}).Association("CollaboratingAccounts").Find(&accounts)
	//Temporarily opening access to all accounts for all users

	err := db.Find(&accounts).Error
	return accounts, err
}

func ApplyAccountTrackingChanges(trackingChanges []*model.AccountTrackingChange, userID string) (string, error) {
	//if new state is 'tracked'
	//  if account does not exist in DB, save it
	//  add account to current user's tracked accounts
	//else if new state is 'untracked'
	//  remove account from user's tracked accounts
	user := FetchUser(userID)
	var error error
	for _, change := range trackingChanges {
		account := createAccountFromNewAccount(*change.Account)
		if change.NewState == "tracked" {
			account, err := SaveAccount(account)
			if error == nil {
				error = err
			}
			//TODO: This might be re-appending existing tracked accounts, but need to check to confirm
			db.Model(&user).Association("TrackedAccounts").Append(account)
		} else if change.NewState == "untracked" {
			db.Model(&user).Association("TrackedAccounts").Delete(account)
		}
	}

	//TODO: return more informative success, and more than just latest error
	return "done", error
}

func createAccountFromNewAccount(newAccount model.NewAccount) *Account {
	id := SafelyUnwrap(newAccount.ID)

	return &Account{
		DBModel:           DBModel{ID: id},
		Name:              newAccount.Name,
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
