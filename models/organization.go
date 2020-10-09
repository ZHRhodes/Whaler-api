package models

import (
	"fmt"

	"github.com/jinzhu/gorm"

	"github.com/heroku/whaler-api/utils"
)

type Organization struct {
	DBModel
	Name  string `json:"name"`
	Users []User `json:"users"`
}

func (org *Organization) Create() map[string]interface{} {
	DB().Create(org)

	if org.ID <= 0 {
		fmt.Print(fmt.Sprint("the org id was less than zero"))
		return utils.Message(5001, "Failed to create organization, connection error", true, map[string]interface{}{})
	}

	data := map[string]interface{}{"organization": org}
	response := utils.Message(2000, "Organization has been created", false, data)
	return response
}

func FetchOrganization(db *gorm.DB, preloads []string, orgID int) (*Organization, error) {
	shouldFetchUsers := false
	for _, value := range preloads {
		if value == "users" || value == "organization.users" {
			shouldFetchUsers = true
		}
	}

	fmt.Printf(fmt.Sprint("Should fetch users: %t", shouldFetchUsers))

	org := &Organization{}
	res := db.Table("organizations").Where("id = ?", orgID)

	if shouldFetchUsers {
		res = res.Preload("Users")
	}

	err := res.First(org).Error

	if err != nil {
		fmt.Print(fmt.Sprint("Failed to fetch Organization", err))
		return nil, err
	}
	return org, nil
}
