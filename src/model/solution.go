package model

import "time"

type SolutionInput struct {
	TaskId int           `json:"taskId"`
	Files  SolutionFiles `json:"files"`
}
type Solution struct {
	Id        int           `json:"id"`
	UserId    int           `json:"userId"`
	Status    string        `json:"status"`
	Result    string        `json:"result"`
	UpdatedAt time.Time     `json:"updatedAt"`
	Files     SolutionFiles `json:"files"`
	Task      TaskMin       `json:"task"`
}
type SolutionMin struct {
	Id        int       `json:"id"`
	UserId    int       `json:"userId"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updatedAt"`
	Task      TaskMin   `json:"task"`
}

func (s *Solution) Min() *SolutionMin {
	return &SolutionMin{
		Id:        s.Id,
		UserId:    s.UserId,
		Status:    s.Status,
		UpdatedAt: s.UpdatedAt,
		Task:      s.Task,
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
	Id         int    `json:"id"`
	SolutionId int    `json:"solutionId"`
	Name       string `json:"name"`
	Code       string `json:"code"`
}

type SolutionFiles []SolutionFile