package api

import (
	"13spr/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}
