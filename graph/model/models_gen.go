// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type AccountID struct {
	ID string `json:"id"`
}

type NewAccount struct {
	ID                *string `json:"id"`
	SalesforceID      *string `json:"salesforceID"`
	Name              string  `json:"name"`
	Owner             string  `json:"owner"`
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
}

type NewContact struct {
	ID           *string `json:"id"`
	SalesforceID *string `json:"salesforceID"`
	FirstName    string  `json:"firstName"`
	LastName     string  `json:"lastName"`
	JobTitle     *string `json:"jobTitle"`
	State        *string `json:"state"`
	Email        *string `json:"email"`
	Phone        *string `json:"phone"`
	AccountID    *string `json:"accountID"`
}

type NewContactAssignmentEntry struct {
	ContactID  string  `json:"contactId"`
	AssignedBy string  `json:"assignedBy"`
	AssignedTo *string `json:"assignedTo"`
}

type NewOrganization struct {
	Name string `json:"name"`
}

type NewUser struct {
	Email          string `json:"email"`
	Password       string `json:"password"`
	OrganizationID string `json:"organizationID"`
}

type NewWorkspace struct {
	Name string `json:"name"`
}

type UserID struct {
	ID string `json:"id"`
}
