package models

import (
	"fmt"

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
		return utils.Message(5001, "Failed to create organization, connection error.", true, map[string]interface{}{})
	}

	data := map[string]interface{}{"organization": org}
	response := utils.Message(2000, "Organization has been created", false, data)
	return response
}
