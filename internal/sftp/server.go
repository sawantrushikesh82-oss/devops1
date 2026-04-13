package sftp

import (
	"context"
	"crypto/subtle"
	"fmt"
	"net"
	"os"

	"devops1/internal/config"
	"devops1/internal/db"
	"devops1/internal/metrics"
	"github.com/pkg/sftp"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh"
)

type Server struct {
	Cfg      *config.Config
	Queries  *db.Queries
	Logger   zerolog.Logger
	UserType string
}

func (s *Server) Start(ctx context.Context, address string) error {
	hostKey, err := os.ReadFile(s.Cfg.SFTP.HostKeyPath)
	if err != nil {
		return err
	}
	signer, err := ssh.ParsePrivateKey(hostKey)
	if err != nil {
		return err
	}
	scfg := &ssh.ServerConfig{
		PublicKeyCallback: s.publicKeyAuth,
		PasswordCallback:  s.passwordAuth,
	}
	scfg.AddHostKey(signer)
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	s.Logger.Info().Str("component", "sftp").Str("addr", address).Msg("listening")
	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		go s.handleConn(ctx, conn, scfg)
	}
}

func (s *Server) handleConn(ctx context.Context, conn net.Conn, scfg *ssh.ServerConfig) {
	defer conn.Close()
	metrics.SFTPConnectionsTotal.WithLabelValues(s.UserType, "accepted").Inc()
	sshConn, chans, reqs, err := ssh.NewServerConn(conn, scfg)
	if err != nil {
		metrics.SFTPConnectionsTotal.WithLabelValues(s.UserType, "rejected").Inc()
		return
	}
	defer sshConn.Close()
	go ssh.DiscardRequests(reqs)
	for ch := range chans {
		if ch.ChannelType() != "session" {
			_ = ch.Reject(ssh.UnknownChannelType, "unknown")
			continue
		}
		channel, requests, err := ch.Accept()
		if err != nil {
			continue
		}
		go func() {
			defer channel.Close()
			for req := range requests {
				if req.Type == "subsystem" && subtle.ConstantTimeCompare(req.Payload[4:], []byte("sftp")) == 1 {
					_ = req.Reply(true, nil)
					sess := &Session{Username: sshConn.User(), UserType: s.UserType, BasePath: s.Cfg.SFTP.DataRoot, Queries: s.Queries, Logger: s.Logger}
					h := &Handlers{VFS: &VirtualHandler{Session: sess}}
					server := sftp.NewRequestServer(channel, sftp.Handlers{FileGet: h, FilePut: h, FileCmd: h, FileList: h})
					_ = server.Serve()
				}
			}
		}()
	}
}

func (s *Server) publicKeyAuth(c ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
	u, err := s.Queries.FindUserByKey(context.Background(), c.User(), string(ssh.MarshalAuthorizedKey(key)))
	if err != nil || string(u.UserType) != s.UserType {
		return nil, fmt.Errorf("denied")
	}
	return nil, nil
}

func (s *Server) passwordAuth(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
	u, err := s.Queries.AuthPassword(context.Background(), c.User())
	if err != nil || string(u.UserType) != s.UserType {
		return nil, fmt.Errorf("denied")
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), pass) != nil {
		return nil, fmt.Errorf("denied")
	}
	return nil, nil
}
