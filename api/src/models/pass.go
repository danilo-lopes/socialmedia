package models

// Pass represents the password struc for password change
type Pass struct {
	New     string `json:"new"`
	Current string `json:"current"`
}
