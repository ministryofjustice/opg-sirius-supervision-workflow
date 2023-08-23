package model

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
	for _, selectedTaskType := range selectedTaskTypes {
		if tt.Handle == selectedTaskType {
			return true
		}
	}
	return false
}
