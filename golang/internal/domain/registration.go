package domain

import "time"

type Registration struct {
	ID        string
	Email     string
	FullName  string
	Phone     string
	CreatedAt time.Time
}

func (r Registration) CreatedAtFormatted() string {
	return r.CreatedAt.UTC().Format("2006-01-02 15:04:05") + " UTC"
}
