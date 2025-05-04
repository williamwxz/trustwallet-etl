package internal

import (
	"fmt"
	"time"
)

type ProcessedUser struct {
	UUID           string    `json:"uuid"`
	FullName       string    `json:"full_name"`
	Email          string    `json:"email"`
	Gender         string    `json:"gender"`
	RegisteredDate time.Time `json:"registered_date"`
	Nationality    string    `json:"nationality"`
	Location       struct {
		City    string `json:"city"`
		Country string `json:"country"`
	} `json:"location"`
	ProcessedAt time.Time `json:"processed_at"`
}

func Transform(raw *RandomUserResponse) (*ProcessedUser, error) {
	if len(raw.Results) == 0 {
		return nil, fmt.Errorf("no results in random user response")
	}

	user := raw.Results[0]
	fullName := fmt.Sprintf("%s %s %s", user.Name.Title, user.Name.First, user.Name.Last)

	return &ProcessedUser{
		FullName:       fullName,
		Email:          user.Email,
		Gender:         user.Gender,
		RegisteredDate: user.Registered.Date.UTC(),
		ProcessedAt:    time.Now().UTC(),
		UUID:           user.Login.UUID,
		Nationality:    user.Nat,
		Location: struct {
			City    string `json:"city"`
			Country string `json:"country"`
		}{
			City:    user.Location.City,
			Country: user.Location.Country,
		},
	}, nil
}
