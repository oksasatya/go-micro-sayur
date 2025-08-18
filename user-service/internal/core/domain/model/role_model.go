package model

import "time"

type Roles struct {
	ID   int64 `gorm:"primaryKey;"`
	Name string
	// Mirror the join column mapping to match user_role(user_id, role_id)
	Users     []User `gorm:"many2many:user_role;joinForeignKey:RoleID;joinReferences:UserID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
