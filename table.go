package posmigrator

import "time"

type Migration struct {
	ID          int64     `sql:"id"`
	Description string    `sql:"description"`
	MigratedAt  time.Time `sql:"migrated_at"`
}
