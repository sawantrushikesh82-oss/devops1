package sftp

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"devops1/internal/models"
	pkg "github.com/pkg/sftp"
)

type VirtualHandler struct{ Session *Session }

func (v *VirtualHandler) allowedReadDir() string {
	if v.Session.UserType == string(models.UserTypeExternal) {
		return filepath.Join(v.Session.BasePath, v.Session.Username, "out")
	}
	return filepath.Join(v.Session.BasePath, v.Session.Username, "in")
}
func (v *VirtualHandler) allowedWriteDir() string {
	if v.Session.UserType == string(models.UserTypeExternal) {
		return filepath.Join(v.Session.BasePath, v.Session.Username, "in")
	}
	return filepath.Join(v.Session.BasePath, v.Session.Username, "out")
}

func cleanAllowed(base, p string) (string, error) {
	joined := filepath.Clean(filepath.Join(base, p))
	if !strings.HasPrefix(joined, filepath.Clean(base)+string(os.PathSeparator)) && joined != filepath.Clean(base) {
		return "", errors.New("traversal")
	}
	return joined, nil
}

func (v *VirtualHandler) AllowedRead(ctx context.Context, path string) (string, error) {
	target, err := cleanAllowed(v.allowedReadDir(), path)
	if err != nil {
		return "", pkg.ErrSSHFxPermissionDenied
	}
	if !strings.HasPrefix(target, filepath.Clean(v.allowedReadDir())) {
		_ = v.Session.Queries.WriteLog(ctx, models.Log{Username: v.Session.Username, Filename: path, Status: models.LogStatusDenied, Component: "sftp", Error: "read denied"})
		return "", pkg.ErrSSHFxPermissionDenied
	}
	return target, nil
}

func (v *VirtualHandler) AllowedWrite(ctx context.Context, path string) (string, error) {
	target, err := cleanAllowed(v.allowedWriteDir(), path)
	if err != nil {
		return "", pkg.ErrSSHFxPermissionDenied
	}
	if !strings.HasPrefix(target, filepath.Clean(v.allowedWriteDir())) {
		_ = v.Session.Queries.WriteLog(ctx, models.Log{Username: v.Session.Username, Filename: path, Status: models.LogStatusDenied, Component: "sftp", Error: "write denied"})
		return "", pkg.ErrSSHFxPermissionDenied
	}
	return target, nil
}
