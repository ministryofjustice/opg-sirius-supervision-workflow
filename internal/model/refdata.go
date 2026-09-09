package model

import "slices"

type RefData struct {
	Handle string `json:"handle"`
	Label  string `json:"label"`
}

func (r RefData) Is(handle string) bool {
	return r.Handle == handle
}

func (r RefData) IsIn(handles []string) bool {
	return slices.ContainsFunc(handles, r.Is)
}
