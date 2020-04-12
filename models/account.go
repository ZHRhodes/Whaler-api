package models

import (
	"time"

	"github.com/heroku/whaler-api/utils"
)

type DBModel struct {
	ID        uint       `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}

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
