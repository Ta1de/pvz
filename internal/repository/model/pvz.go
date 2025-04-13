package model

import (
	"time"

	"github.com/google/uuid"
)

type Pvz struct {
	Id               uuid.UUID `db:"id"`
	RegistrationDate time.Time `db:"registrationdate"`
	City             string    `db:"city"`
}
