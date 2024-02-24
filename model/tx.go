package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Deposit struct {
	bun.BaseModel `bun:"table:deposits"`

    ID        int64     `bun:"id,pk,autoincrement" json:"id"`
    Amount    float64   `bun:",notnull" json:"amount"`
    Status    string    `bun:",notnull" json:"status"`
    CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"createdAt"`

	UserID int64 `bun:",notnull"`
	User   *User `bun:"rel:belongs-to,join:user_id=id"`
}

type Purchase struct {
	bun.BaseModel `bun:"table:purchases"`

    ID        int64     `bun:"id,pk,autoincrement" json:"id"`
    Amount    float64   `bun:",notnull" json:"amount"`
    Status    string    `bun:",notnull" json:"status"`
    CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"createdAt"`

	UserID    int64    `bun:",notnull"`
	User      *User    `bun:"rel:belongs-to,join:user_id=id"`
	ProductID int64    `bun:",notnull"`
	Product   *Product `bun:"rel:belongs-to,join:product_id=id"`
}
