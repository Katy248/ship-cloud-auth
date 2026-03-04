package data

import "time"

type TimestampsModel struct {
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}

func (m *TimestampsModel) UpdateTimestamp() {
	m.UpdatedAt = time.Now()
}
