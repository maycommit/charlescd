package models

type Workspace struct {
	Base
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Clusters    []Cluster `json:"clusters"`
}
