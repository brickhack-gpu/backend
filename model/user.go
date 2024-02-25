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

	Deposits      []*Deposit      `bun:"rel:has-many,join:id=user_id"`
	Purchases     []*Purchase     `bun:"rel:has-many,join:id=user_id"`
	Notifications []*Notification `bun:"rel:has-many,join:id=user_id"`
	Products      []*Product      `bun:"rel:has-many,join:id=user_id"`
}

type Notification struct {
	bun.BaseModel `bun:"table:notifications"`

	ID        int64     `bun:"id,pk,autoincrement" json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Read      bool      `bun:"default:false" json:"read"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"createdAt"`

	UserID int64 `bun:",notnull"`
	User   *User `bun:"rel:belongs-to,join:user_id=id"`
}
