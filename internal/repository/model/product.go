package model

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	Id          uuid.UUID `db:"id"`
	DateTime    time.Time `db:"datetime"`
	Type        string    `db:"type"`
	ReceptionId uuid.UUID `db:"receptionid"`
}
