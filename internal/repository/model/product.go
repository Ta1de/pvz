package model

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	Id          uuid.UUID
	DateTime    time.Time
	Type        string
	ReceptionId uuid.UUID
}
