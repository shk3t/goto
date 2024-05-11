package model

import "time"

type DelayedTask struct {
	Id        int       `json:"id"`
	UserId    int       `json:"userId"`
	TargetId  *int      `json:"targetId"`
	Action    string    `json:"action"`
	Target    string    `json:"target"`
	Status    string    `json:"status"`
	Details   string    `json:"details"`
	UpdatedAt time.Time `json:"updatedAt"`
}