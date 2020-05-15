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

func FetchOrganization(db *gorm.DB, orgID string) (*Organization, error) {
	org := &Organization{}
	//removed .Preload("Users") before .First.. how should i handle that with graphql, if at all?
	err := db.Table("organizations").Where("id = ?", orgID).First(org).Error
	if err != nil {
		return nil, err
	}
	return org, nil
}
