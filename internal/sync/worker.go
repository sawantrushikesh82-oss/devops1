package sync

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"devops1/internal/config"
	"devops1/internal/db"
	"devops1/internal/models"
	"github.com/google/uuid"
)

type Worker struct {
	Cfg     *config.Config
	Queries *db.Queries
}

func (w *Worker) Sync(ctx context.Context, m models.Mapping) error {
	switch m.FlowType {
	case models.FlowTypeSendReceive:
		if err := w.syncDir(ctx, m.ExternalUser, "in", m.InternalUser, "in", m); err != nil {
			return err
		}
		return w.syncDir(ctx, m.InternalUser, "out", m.ExternalUser, "out", m)
	case models.FlowTypeSendOnly:
		return w.syncDir(ctx, m.InternalUser, "out", m.ExternalUser, "out", m)
	case models.FlowTypeReceiveOnly:
		return w.syncDir(ctx, m.ExternalUser, "in", m.InternalUser, "in", m)
	default:
		return fmt.Errorf("unknown flow type: %s", m.FlowType)
	}
}

func (w *Worker) syncDir(ctx context.Context, srcUser, srcSub, dstUser, dstSub string, m models.Mapping) error {
	src := filepath.Join(w.Cfg.SFTP.DataRoot, srcUser, srcSub)
	dst := filepath.Join(w.Cfg.SFTP.DataRoot, dstUser, dstSub)
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		rel, _ := filepath.Rel(src, path)
		dstPath := filepath.Join(dst, rel)
		if err := os.MkdirAll(filepath.Dir(dstPath), 0o755); err != nil {
			return err
		}
		srcRec := FileRecord{Path: path, ModTime: info.ModTime(), Size: info.Size()}
		dstInfo, statErr := os.Stat(dstPath)
		if statErr == nil {
			dstRec := FileRecord{Path: dstPath, ModTime: dstInfo.ModTime(), Size: dstInfo.Size()}
			if !NeedsSync(srcRec, dstRec) {
				return nil
			}
			if w.Cfg.Sync.ChecksumEnabled && ChecksumMatch(path, dstPath) {
				return nil
			}
		}
		return WithRetry(3, 5*time.Second, func() error {
			if err := copyAtomic(path, dstPath); err != nil {
				return err
			}
			return w.Queries.WriteLog(ctx, models.Log{Username: srcUser, Filename: rel, Direction: srcSub + "->" + dstSub, Bytes: info.Size(), Status: models.LogStatusSuccess, MappingID: m.ID, Component: "sync"})
		})
	})
}

func copyAtomic(src, dst string) error {
	sf, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sf.Close()
	tmp := filepath.Join(filepath.Dir(dst), ".tmp_"+uuid.NewString())
	tf, err := os.OpenFile(tmp, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	if _, err = io.Copy(tf, sf); err != nil {
		_ = tf.Close()
		return err
	}
	if err = tf.Close(); err != nil {
		return err
	}
	return os.Rename(tmp, dst)
}
