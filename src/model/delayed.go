package model

import "time"

type DelayedTask struct {
	Id         int       `json:"id"`
	UserId     int       `json:"userId"`
	TargetId   *int      `json:"targetId"`
	Action     string    `json:"action"`
	TargetName string    `json:"targetName"`
	Status     string    `json:"status"`
	Details    string    `json:"details"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type DelayedTasks []DelayedTask