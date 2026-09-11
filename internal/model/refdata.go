package model

import "slices"

type RefData struct {
	Handle string
	Label  string
}

func (r RefData) Is(handle string) bool {
	return r.Handle == handle
}

func (r RefData) IsIn(handles []string) bool {
	return slices.ContainsFunc(handles, r.Is)
}
