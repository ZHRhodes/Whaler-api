package models

import (
	"fmt"

	"github.com/jinzhu/gorm"

	"github.com/heroku/whaler-api/utils"
)

type Organization struct {
	DBModel
	Name  string `json:"name"`
	Users []User `json:"users" gorm:"foreignkey:OrganizationID"`
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

func FetchOrg(orgID string) map[string]interface{} {
	org := &Organization{}
	err := DB().Table("organizations").Where("id = ?", orgID).Preload("Users").First(org).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.Message(5001, "Organization with the given id not found", true, map[string]interface{}{})
		} else {
			return utils.Message(5001, "Unable to fetch organization, connection error", true, map[string]interface{}{})
		}
	}

	return utils.Message(2000, "Organization fetched successfully", false, org)
}
