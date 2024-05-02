package model

type ProjectPublic struct {
	ProjectBase
	Id    int `json:"id"`
	Tasks []Task `json:"tasks"`
}