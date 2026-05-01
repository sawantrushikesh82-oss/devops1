package sftp

import (
	"devops1/internal/db"
	"github.com/rs/zerolog"
)

type Session struct {
	Username string
	UserType string
	BasePath string
	Queries  *db.Queries
	Logger   zerolog.Logger
}
