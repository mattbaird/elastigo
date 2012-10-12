package main

import (
	"time"
)

// used in test suite, chosen to be similar to the documentation
type Tweet struct {
	User     string    `json:"user"`
	PostDate time.Time `json:"postDate"`
	Message  string    `json:"message"`
}

func NewTweet(user string, message string) Tweet {
	return Tweet{User: user, PostDate: time.Now(), Message: message}
}
