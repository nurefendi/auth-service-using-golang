package dao

import (
	"github.com/google/uuid"
	"time"
)

type AuditorDAO struct {
	ID         uuid.UUID  `gorm:"primaryKey;type:uuid"`
	CreatedBy  string     `gorm:"'created_by'"`
	CreatedAt  time.Time  `gorm:"autoCreateTime:true'created_at'"`
	ModifiedBy *string    `gorm:"'modified_by'"`
	ModifiedAt *time.Time `gorm:"autoUpdateTime:true'modified_at'"`
}

// GenderLang represents the gender_lang table
type GenderLang struct {
	Gender int    `gorm:"type:tinyint(1);not null"`
	Lang   string `gorm:"type:varchar(20);not null"`
	Name   string `gorm:"type:varchar(100);not null"`
	AuditorDAO
}

// AuthUser represents the auth_user table
type AuthUser struct {
	FullName   string       `gorm:"type:varchar(255);not null"`
	Email      string       `gorm:"type:varchar(100);not null"`
	Username   string       `gorm:"type:varchar(100);not null"`
	Password   string       `gorm:"type:varchar(255);not null"`
	Gender     int          `gorm:"type:tinyint(1);not null"`
	Telephone  *string      `gorm:"type:varchar(15)"`
	HasDeleted bool         `gorm:"not null;default:false"`
	Picture    *string      `gorm:"type:varchar(255)"`
	GenderLang []GenderLang `gorm:"foreignKey:Gender;references:Gender"`
	AuditorDAO
}

// AuthGroup represents the auth_group table
type AuthGroup struct {
	Name        string  `gorm:"type:varchar(255);not null"`
	Description *string `gorm:"type:text"`
	AuditorDAO
}

// AuthUserGroup represents the auth_user_group table
type AuthUserGroup struct {
	UserID  uuid.UUID `gorm:"type:uuid;not null"`
	GroupID uuid.UUID `gorm:"type:uuid;not null"`
	User    AuthUser  `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	Group   AuthGroup `gorm:"foreignKey:GroupID;references:ID;constraint:OnDelete:CASCADE"`
	AuditorDAO
}

// AuthPortal represents the auth_portal table
type AuthPortal struct {
	Order    int              `gorm:"type:int(3);not null"`
	Path     string           `gorm:"type:varchar(255);not null"`
	Icon     *string          `gorm:"type:varchar(255)"`
	FontIcon *string          `gorm:"type:varchar(50)"`
	Lang     []AuthPortalLang `gorm:"foreignKey:PortalID;constraint:OnDelete:CASCADE"`
	AuditorDAO
}

// AuthPortalLang represents the auth_portal_lang table
type AuthPortalLang struct {
	PortalID    uuid.UUID `gorm:"type:uuid;not null"`
	Name        string    `gorm:"type:varchar(255);not null"`
	Description *string   `gorm:"type:text"`
	Lang        string    `gorm:"type:varchar(20);not null"`
	AuditorDAO
}

// AuthFunction represents the auth_function table
type AuthFunction struct {
	PortalID    uuid.UUID          `gorm:"type:uuid;not null"`
	Position    string             `gorm:"type:varchar(20);not null"`
	Icon        *string            `gorm:"type:varchar(255)"`
	FontIcon    *string            `gorm:"type:varchar(50)"`
	IsShow      bool               `gorm:"type:tinyint(1);not null;default:true"`
	ShortcutKey *string            `gorm:"type:varchar(100)"`
	Path        string             `gorm:"type:varchar(255);not null"`
	Lang        []AuthFunctionLang `gorm:"foreignKey:FunctionID;constraint:OnDelete:CASCADE"`
	AuditorDAO
}

// AuthFunctionLang represents the auth_function_lang table
type AuthFunctionLang struct {
	FunctionID  uuid.UUID `gorm:"type:uuid;not null"`
	Name        string    `gorm:"type:varchar(255);not null"`
	Lang        string    `gorm:"type:varchar(20);not null"`
	Description *string   `gorm:"type:text"`
	AuditorDAO
}

// AuthPermission represents the auth_permission table
type AuthPermission struct {
	GroupID     uuid.UUID    `gorm:"type:uuid;not null"`
	FunctionID  uuid.UUID    `gorm:"type:uuid;not null"`
	GrandCreate int          `gorm:"type:tinyint(1);not null;default:0"`
	GrandRead   int          `gorm:"type:tinyint(1);not null;default:0"`
	GrandUpdate int          `gorm:"type:tinyint(1);not null;default:0"`
	GrandDelete int          `gorm:"type:tinyint(1);not null;default:0"`
	Function    AuthFunction `gorm:"foreignKey:FunctionID;references:ID;constraint:OnDelete:CASCADE"`
	AuditorDAO
}
