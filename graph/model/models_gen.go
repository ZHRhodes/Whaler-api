// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type AccountID struct {
	ID int `json:"id"`
}

type NewAccount struct {
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
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	JobTitle  *string `json:"jobTitle"`
	State     *string `json:"state"`
	Email     *string `json:"email"`
	Phone     *string `json:"phone"`
	AccountID *string `json:"accountID"`
}

type NewContactAssignmentEntry struct {
	ContactID  int     `json:"contactId"`
	AssignedBy string  `json:"assignedBy"`
	AssignedTo *string `json:"assignedTo"`
}

type NewUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type NewWorkspace struct {
	Name string `json:"name"`
}

type UserID struct {
	ID int `json:"id"`
}
