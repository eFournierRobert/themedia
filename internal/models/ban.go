package models

import "time"

type Ban struct {
	EndDatetime time.Time `json:"end_datetime"`
}
