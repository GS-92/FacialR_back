package faceR_user

import "time"

type User struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Entries  int       `json:"entries"`
	Joined   time.Time `json:"joined"`
}
