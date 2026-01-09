package model

import "time"

type Task struct {
	ID          string
	Title       string
	Description string
	Created_At  time.Time
}
