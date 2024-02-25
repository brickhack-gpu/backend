package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Product struct {
	bun.BaseModel `bun:"table:products"`

	ID          int64     `bun:"id,pk,autoincrement" json:"id"`
	Price       float64   `bun:",notnull" json:"price"`
	Status      string    `bun:",notnull" json:"status"` // active, stopped, destroyed
	DNSLink     string    `bun:"dns_link" json:"dns_link"`
	GCPID       string    `bun:"gcp_id" json:"gcp_id"`
	Credentials string    `bun:"credentials" json:"credentials"`
	Storage     int       `bun:",notnull" json:"storage"`
	CreatedAt   time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"createdAt"`

	UserID         int64         `bun:",notnull"`
	User           *User         `bun:"rel:belongs-to,join:user_id=id"`
	ServerConfigID int64         `bun:",notnull"`
	ServerConfig   *ServerConfig `bun:"rel:belongs-to,join:server_config_id=id"`
	TemplateID     int64         `bun:",notnull"`
	Template       *Template     `bun:"rel:belongs-to,join:template_id=id"`
}

type ServerConfig struct {
	bun.BaseModel `bun:"table:server_types"`

	ID       int64   `bun:"id,pk,autoincrement" json:"id"`
	Region   string  `bun:",notnull" json:"region"`
	Zone     string  `bun:",notnull" json:"zone"`
	GPUType  string  `bun:"gpu_type" json:"gpu_type"`
	GPUCount int     `bun:"gpu_count" json:"gpu_count"`
	Price    float64 `bun:",notnull" json:"price"`
    MachineType string `bun:",notnull" json:"machine_type"`
    Active bool `bun:"default:true" json:"active"`
}

type Template struct {
	bun.BaseModel `bun:"table:templates"`

	ID          int64  `bun:"id,pk,autoincrement" json:"id"`
	Container   string `bun:"container" json:"container"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"` // image generation, text generation
    Active bool `bun:"default:true" json:"active"`
}
