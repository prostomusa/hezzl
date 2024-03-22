package clickhouse

import (
	"database/sql"
	"time"
)

type ClickHouseLog struct {
	Id          int            `db:"Id" json:"id"`
	ProjectId   int            `db:"ProjectId" json:"project_id"`
	Name        string         `db:"Name" json:"name"`
	Description sql.NullString `db:"Description" json:"description"`
	Priority    int            `db:"Priority" json:"priority"`
	Removed     bool           `db:"Removed" json:"removed"`
	EventTime   time.Time      `db:"EventTime" json:"event_time"`
}
