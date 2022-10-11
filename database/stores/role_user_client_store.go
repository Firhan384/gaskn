package stores

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

/*
*
Table model
*/
type RoleUserClient struct {
	gorm.Model
	ID         uuid.UUID `gorm:"type:char(36);primary_key"`
	ClientId   uuid.UUID `gorm:"type:char(36):index"`
	RoleUserId uuid.UUID `gorm:"type:char(36):index"`
	Client     Client    `gorm:"references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	RoleUser   RoleUser  `gorm:"references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	IsActive   bool
}

/*
*
This function is a feature that gorm has for making hooks,
this hook function is used to generate uuid when the user
performs the create action
*/
func (*RoleUserClient) BeforeCreate(tx *gorm.DB) error {
	tx.Statement.SetColumn("ID", uuid.New())
	return nil
}
