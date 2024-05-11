package model

import "time"

type SolutionInput struct {
	TaskId int            `json:"taskId"`
	Files  []SolutionFile `json:"files"`
}
type Solution struct {
	Id        int            `json:"id"`
	UserId    int            `json:"userId"`
	TaskId    int            `json:"taskId"`
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

type Solutions []Solution
type SolutionsMin []SolutionMin

func (solutions Solutions) Min() SolutionsMin {
	solutionsMin := make(SolutionsMin, len(solutions))
	for i, s := range solutions {
		solutionsMin[i] = *s.Min()
	}
	return solutionsMin
}

type SolutionFile struct {
	Id         int
	SolutionId int
	Name       string
	Code       string
}