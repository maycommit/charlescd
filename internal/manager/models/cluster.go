package models

import (
	"github.com/google/uuid"
)

type Cluster struct {
	Base
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Address     string    `json:"address"`
	AppKey      string    `json:"appKey"`
	WorkspaceID uuid.UUID `json:"workspaceId"`
}
