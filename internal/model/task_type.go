package model

import "slices"

type TaskType struct {
	Handle     string `json:"handle"`
	Incomplete string `json:"incomplete"`
	Category   string `json:"category"`
	Complete   string `json:"complete"`
	User       bool   `json:"user"`
	EcmTask    bool   `json:"ecmTask"`
	TaskCount  int
}

func (tt TaskType) IsSelected(selectedTaskTypes []string) bool {
	return slices.Contains(selectedTaskTypes, tt.Handle)
}
