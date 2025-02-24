package dao

import (
    "time"
    "github.com/google/uuid"
)


type AuditorDAO struct {
	ID         uuid.UUID `gorm:"primaryKey;type:uuid"`
	CreatedBy  string    `gorm:"'created_by'"`
	CreatedAt  time.Time `gorm:"autoCreateTime:true'created_at'"`
	ModifiedBy *string    `gorm:"'modified_by'"`
	ModifiedAt *time.Time `gorm:"autoUpdateTime:true'modified_at'"`
}

// GenderLang represents the gender_lang table
type GenderLang struct {
    Gender     int       `gorm:"type:tinyint(1);not null"`
    Lang       string    `gorm:"type:varchar(20);not null"`
    Name       string    `gorm:"type:varchar(100);not null"`
	AuditorDAO
}

// AuthUser represents the auth_user table
type AuthUser struct {
    FullName   string    `gorm:"type:varchar(255);not null"`
    Email      string    `gorm:"type:varchar(100);not null"`
    Username   string    `gorm:"type:varchar(100);not null"`
    Password   string    `gorm:"type:varchar(255);not null"`
    Gender     int       `gorm:"type:tinyint(1);not null"`
    Telephone  *string   `gorm:"type:varchar(15)"`
    HasDeleted int       `gorm:"type:tinyint(1);not null;default:0"`
    Picture    *string   `gorm:"type:varchar(255)"`
	AuditorDAO
}

// AuthGroup represents the auth_group table
type AuthGroup struct {
    Name        string    `gorm:"type:varchar(255);not null"`
    Description *string   `gorm:"type:text"`
    AuditorDAO
}

// AuthUserGroup represents the auth_user_group table
type AuthUserGroup struct {
    UserID    uuid.UUID `gorm:"type:uuid;not null"`
    GroupID   uuid.UUID `gorm:"type:uuid;not null"`
    AuditorDAO
}

// AuthPortal represents the auth_portal table
type AuthPortal struct {
    Order     int       `gorm:"type:int(3);not null"`
    Path      string    `gorm:"type:varchar(255);not null"`
    Icon      *string   `gorm:"type:varchar(255)"`
    FontIcon  *string   `gorm:"type:varchar(50)"`
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
    PortalID    uuid.UUID `gorm:"type:uuid;not null"`
    Position    string    `gorm:"type:varchar(20);not null"`
    Icon        *string   `gorm:"type:varchar(255)"`
    FontIcon    *string   `gorm:"type:varchar(50)"`
    IsShow      int       `gorm:"type:tinyint(1);not null;default:1"`
    ShortcutKey *string   `gorm:"type:varchar(100)"`
    Path        string    `gorm:"type:varchar(255);not null"`
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
    GroupID     uuid.UUID `gorm:"type:uuid;not null"`
    FunctionID  uuid.UUID `gorm:"type:uuid;not null"`
    GrandCreate int       `gorm:"type:tinyint(1);not null;default:0"`
    GrandRead   int       `gorm:"type:tinyint(1);not null;default:0"`
    GrandUpdate int       `gorm:"type:tinyint(1);not null;default:0"`
    GrandDelete int       `gorm:"type:tinyint(1);not null;default:0"`
	AuditorDAO
}