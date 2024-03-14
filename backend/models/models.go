package models

import "time"

type ContactForm struct {
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"createdAt"`
}

type VideoWaitlistForm struct {
	Challenge   string    `json:"challenge"`
	Email       string    `json:"email"`
	Enhancement string    `json:"enhancement"`
	Features    string    `json:"features"`
	Feedback    string    `json:"feedback"`
	Tools       string    `json:"tools"`
	CreatedAt   time.Time `json:"createdAt"`
}
