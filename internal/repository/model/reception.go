package model

import (
	"time"

	"github.com/google/uuid"
)

type Reception struct {
	Id       uuid.UUID `db:"id"`
	DateTime time.Time `db:"datetime"`
	PvzId    uuid.UUID `db:"pvzid"`
	Status   string    `db:"status"`
}
