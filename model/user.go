package model

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users"`

	ID           int64     `bun:"id,pk,autoincrement"`
	Username     string    `bun:",notnull"`
	PasswordHash string    `bun:",notnull"`
	Active       bool      `bun:"default:true"`
	Admin        bool      `bun:"default:false"`
	SSHKey       string    `bun:"ssh_key"`
	CreatedAt    time.Time `bun:",nullzero,notnull,default:current_timestamp"`

	Deposits  []*Deposit  `bun:"rel:has-many,join:id=user_id"`
	Purchases []*Purchase `bun:"rel:has-many,join:id=user_id"`
	Products  []*Product  `bun:"rel:has-many,join:id=user_id"`
}
