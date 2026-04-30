package models

import "time"

type AuditLog struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	AdminID   string    `gorm:"type:uuid;not null;index" json:"admin_id"`
	Action    string    `gorm:"type:varchar(100);not null" json:"action"`
	Target    string    `gorm:"type:varchar(100)" json:"target"`
	Details   string    `gorm:"type:text" json:"details"`
	IPAddress string    `gorm:"type:varchar(45)" json:"ip_address"`
	CreatedAt time.Time `json:"created_at"`

	Admin User `gorm:"foreignKey:AdminID" json:"admin"`
}

type SystemConfig struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Key       string    `gorm:"uniqueIndex;not null" json:"key"`
	Value     string    `gorm:"type:text;not null" json:"value"`
	Category  string    `gorm:"type:varchar(50);index" json:"category"`
	UpdatedAt time.Time `json:"updated_at"`
}
