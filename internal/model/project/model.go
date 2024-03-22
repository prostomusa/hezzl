package project

import (
	"time"
)

type Project struct {
	Id        int       `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
}
