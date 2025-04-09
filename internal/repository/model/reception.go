package model

import (
	"time"

	"github.com/google/uuid"
)

type Reception struct {
	Id       uuid.UUID
	DateTime time.Time
	PvzId    uuid.UUID
	Status   string
}
