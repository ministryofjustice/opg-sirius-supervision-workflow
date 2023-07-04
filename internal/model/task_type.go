package model

type TaskType struct {
	Handle     string `json:"handle"`
	Incomplete string `json:"incomplete"`
	Category   string `json:"category"`
	Complete   string `json:"complete"`
	User       bool   `json:"user"`
	EcmTask    bool   `json:"ecmTask"`
	IsSelected bool
	TaskCount  int
}
