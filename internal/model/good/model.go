package good

import (
	"database/sql"
	"time"
)

type GoodEntity struct {
	Id          int            `db:"id" json:"id"`
	ProjectId   int            `db:"project_id" json:"project_id"`
	Name        string         `db:"name" json:"name"`
	Description sql.NullString `db:"description" json:"description"`
	Priority    int            `db:"priority" json:"priority"`
	Removed     bool           `db:"removed" json:"removed"`
	CreatedAt   *time.Time     `db:"created_at" json:"created_at"`
}
