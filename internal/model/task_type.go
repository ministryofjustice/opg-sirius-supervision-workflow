package model

import "slices"

type TaskType struct {
	Handle     string
	Incomplete string
	Category   string
	Complete   string
	User       bool
	EcmTask    bool
	TaskCount  int
}

func (tt TaskType) IsSelected(selectedTaskTypes []string) bool {
	return slices.Contains(selectedTaskTypes, tt.Handle)
}
