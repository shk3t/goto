package model

import "time"

type SolutionInput struct {
	TaskId int            `json:"taskId"`
	Files  []SolutionFile `json:"files"`
}
type Solution struct {
	Id        int            `json:"id"`
	TaskId    int            `json:"taskId"`
	UserId    int            `json:"userId"`
	Status    string         `json:"status"`
	Result    string         `json:"result"`
	UpdatedAt time.Time      `json:"updatedAt"`
	Files     []SolutionFile `json:"files"`
}
type SolutionMin struct {
	Id        int       `json:"id"`
	TaskId    int       `json:"taskId"`
	UserId    int       `json:"userId"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (s *Solution) Min() *SolutionMin {
	return &SolutionMin{
		Id:        s.Id,
		TaskId:    s.TaskId,
		UserId:    s.UserId,
		Status:    s.Status,
		UpdatedAt: s.UpdatedAt,
	}
}

type SolutionFile struct {
	Id         int
	SolutionId int
	Name       string
	Code       string
}