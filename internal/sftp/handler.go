package sftp

import (
	"context"
	"io"
	"os"
	"path/filepath"

	pkg "github.com/pkg/sftp"
)

type Handlers struct{ VFS *VirtualHandler }

func (h *Handlers) Fileread(r *pkg.Request) (io.ReaderAt, error) {
	p, err := h.VFS.AllowedRead(context.Background(), r.Filepath)
	if err != nil {
		return nil, err
	}
	return os.Open(p)
}
func (h *Handlers) Filewrite(r *pkg.Request) (io.WriterAt, error) {
	p, err := h.VFS.AllowedWrite(context.Background(), r.Filepath)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return nil, err
	}
	return os.OpenFile(p, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
}
func (h *Handlers) Filelist(r *pkg.Request) (pkg.ListerAt, error) {
	p, err := h.VFS.AllowedRead(context.Background(), r.Filepath)
	if err != nil {
		return nil, err
	}
	return pkg.ReadDir(p)
}
func (h *Handlers) Filecmd(r *pkg.Request) error {
	_, err := h.VFS.AllowedWrite(context.Background(), r.Filepath)
	return err
}
func (h *Handlers) Filestat(r *pkg.Request) (pkg.ListerAt, error) {
	p, err := h.VFS.AllowedRead(context.Background(), r.Filepath)
	if err != nil {
		return nil, err
	}
	return pkg.ReadDir(p)
}
