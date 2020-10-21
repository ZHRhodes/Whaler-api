package models

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/heroku/whaler-api/utils"
	"github.com/heroku/whaler-api/graph/model"
)

type Organization struct {
	DBModel
	Name  string `json:"name"`
	Users []User `json:"users"`
}

func (org *Organization) Create() map[string]interface{} {
	DB().Create(org)

	if len(org.ID) == 0 {
		fmt.Print(fmt.Sprint("the org id was not set"))
		return utils.Message(5001, "Failed to create organization, connection error", true, map[string]interface{}{})
	}

	data := map[string]interface{}{"organization": org}
	response := utils.Message(2000, "Organization has been created", false, data)
	return response
}

func CreateOrganization(db *gorm.DB, newOrganization model.NewOrganization) (*Organization, error) {
	organization := &Organization{
		Name: newOrganization.Name,
	}

	err := db.Create(organization).Error

	return organization, err
}

func FetchOrganization(db *gorm.DB, preloads []string, orgID string) (*Organization, error) {
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

	for idx := range org.Users {
		org.Users[idx].Password = ""
	}

	if err != nil {
		fmt.Print(fmt.Sprint("Failed to fetch Organization", err))
		return nil, err
	}

	return org, nil
}
