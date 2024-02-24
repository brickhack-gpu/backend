package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Product struct {
	bun.BaseModel `bun:"table:products"`

	ID          int64     `bun:"id,pk,autoincrement"`
	Price       float64   `bun:",notnull"`
	Status      string    `bun:",notnull"` // active, stopped, destroyed
	DNSLink     string    `bun:"dns_link"`
	GCPID       string    `bun:"gcp_id"`
	Credentials string    `bun:"credentials"`
	Storage     int       `bun:",notnull"`
	CreatedAt   time.Time `bun:",nullzero,notnull,default:current_timestamp"`

	UserID         int64         `bun:",notnull"`
	User           *User         `bun:"rel:belongs-to,join:user_id=id"`
	ParentID       int64         `bun:",notnull"`
	Product        *Product      `bun:"rel:belongs-to,join:parent_id=id"`
	ServerConfigId int64         `bun:",notnull"`
	ServerConfig   *ServerConfig `bun:"rel:belongs-to,join:server_config_id=id"`
}

type ServerConfig struct {
	bun.BaseModel `bun:"table:server_types"`

	ID       int64   `bun:"id,pk,autoincrement"`
	Region   string  `bun:",notnull"`
	Zone     string  `bun:",notnull"`
	GPUType  string  `bun:"gpu_type"`
	GPUCount int     `bun:"gpu_count"`
	Price    float64 `bun:",notnull"`
}
