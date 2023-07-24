package model

type RefData struct {
	Handle string `json:"handle"`
	Label  string `json:"label"`
}

func (r RefData) Is(handle string) bool {
	return r.Handle == handle
}

func (r RefData) IsIn(handles []string) bool {
	for _, handle := range handles {
		if r.Is(handle) {
			return true
		}
	}
	return false
}
