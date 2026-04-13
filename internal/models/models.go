package models

import "time"

type UserType string

type Role string

type FlowType string

type LogStatus string

const (
	UserTypeExternal UserType = "external"
	UserTypeInternal UserType = "internal"

	RoleAdmin  Role = "admin"
	RoleViewer Role = "viewer"

	FlowTypeSendReceive FlowType = "send_receive"
	FlowTypeSendOnly    FlowType = "send_only"
	FlowTypeReceiveOnly FlowType = "receive_only"

	LogStatusSuccess LogStatus = "success"
	LogStatusFailed  LogStatus = "failed"
	LogStatusDenied  LogStatus = "access_denied"
)

type User struct {
	ID           string    `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	PasswordHash string    `json:"-" db:"password_hash"`
	PublicKey    string    `json:"public_key" db:"public_key"`
	UserType     UserType  `json:"user_type" db:"user_type"`
	Role         Role      `json:"role" db:"role"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type Mapping struct {
	ID           string    `json:"id" db:"id"`
	ExternalUser string    `json:"external_user" db:"external_user"`
	InternalUser string    `json:"internal_user" db:"internal_user"`
	FlowType     FlowType  `json:"flow_type" db:"flow_type"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type Log struct {
	ID         string    `json:"id" db:"id"`
	Username   string    `json:"username" db:"username"`
	Filename   string    `json:"filename" db:"filename"`
	Direction  string    `json:"direction" db:"direction"`
	Bytes      int64     `json:"bytes" db:"bytes"`
	Status     LogStatus `json:"status" db:"status"`
	Error      string    `json:"error" db:"error"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	MappingID  string    `json:"mapping_id" db:"mapping_id"`
	SessionID  string    `json:"session_id" db:"session_id"`
	Component  string    `json:"component" db:"component"`
	DurationMS int64     `json:"duration_ms" db:"duration_ms"`
}

type Session struct {
	ID           string    `json:"id" db:"id"`
	UserID       string    `json:"user_id" db:"user_id"`
	RefreshToken string    `json:"refresh_token" db:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
