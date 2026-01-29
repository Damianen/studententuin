package domain

import "time"

type User struct {
	id string
	email string
	name string
	avatarUrl string
	githubId string
	githubUsername string
	githubAccessToken string
	createdAt time.Time
	updatedA time.Time
}
