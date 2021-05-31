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
	SalesforceID      *string                  `json:"salesforceID" gorm:"uniqueIndex"`
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
	Notes             *string                  `json:"notes" gorm:"-"`
	AssignmentEntries []AccountAssignmentEntry `json:"assignmentEntries" gorm:"foreignKey:AccountID;references:ID"`
	Trackers          []*User                  `json:"trackers" gorm:"many2many:account_trackers;"`
	Collaborators     []User                   `json:"collaborators" gorm:"many2many:account_collaborators;"`
	Contacts          []*Contact               `json:"contacts"`
	AssignedTo        *string                  `json:"assignedTo"`
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

func SaveAccounts(senderID *string, newAccounts []*model.NewAccount, userID string) ([]*Account, error) {
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

	//Same as for contacts:
	//A rare short term solution: this is to idenfiy single saves
	//aka a FE changed something specifically. Will not send notifications for
	//Salesforce sync saves. After moving SF to BE, this can be refined more easily.
	if len(newAccounts) == 1 {
		go SendAccountChangeMessage(senderID, userID)
	}

	//will only return latest error
	return savedAccounts, error
}

func SendAccountChangeMessage(senderID *string, userID string) {
	user := FetchUser(userID)
	Consumer.ModelChanged(user.OrganizationID, senderID)
}

func SaveAccount(account *Account) (*Account, error) {
	err := DB().Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "salesforce_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"updated_at", "name", "industry",
			"salesforce_id", "description", "number_of_employees", "annual_revenue",
			"billing_city", "billing_state", "phone", "website", "type", "state", "assigned_to"}),
	}).Create(&account).Error

	return account, err
}

func FetchAccounts(userID string) ([]*Account, error) {
	var accounts = []*Account{}
	user := FetchUser(userID)
	err := db.Model(&user).Association("TrackedAccounts").Find(&accounts)
	return accounts, err
}

func ApplyAccountTrackingChanges(trackingChanges []*model.AccountTrackingChange, userID string) ([]*Account, error) {
	user := FetchUser(userID)
	var error error
	for _, change := range trackingChanges {
		account := createAccountFromNewAccount(*change.Account)
		var existingAccount Account
		if account.SalesforceID != nil {
			db.First(&existingAccount, "salesforce_id = ?", account.SalesforceID)
		}

		account.ID = existingAccount.ID

		if change.NewState == "tracked" {
			db.Model(&user).Association("TrackedAccounts").Append(account)
		} else if change.NewState == "untracked" {
			db.Model(&user).Association("TrackedAccounts").Delete(account)
		}
	}

	var trackedAccounts []*Account
	db.Model(&user).Association("TrackedAccounts").Find(&trackedAccounts)

	//TODO: return more informative response, and more than just latest error
	return trackedAccounts, error
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
		AssignedTo:        newAccount.AssignedTo,
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
