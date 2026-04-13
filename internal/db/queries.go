package db

import (
	"context"
	"fmt"

	"devops1/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Queries struct{ Pool *pgxpool.Pool }

func New(pool *pgxpool.Pool) *Queries { return &Queries{Pool: pool} }

func (q *Queries) GetActiveMappings(ctx context.Context) ([]models.Mapping, error) {
	rows, err := q.Pool.Query(ctx, `SELECT id,external_user,internal_user,flow_type,is_active,created_at,updated_at FROM mappings WHERE is_active=true`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Mapping
	for rows.Next() {
		var m models.Mapping
		if err := rows.Scan(&m.ID, &m.ExternalUser, &m.InternalUser, &m.FlowType, &m.IsActive, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

func (q *Queries) FindUserByKey(ctx context.Context, username, publicKey string) (*models.User, error) {
	var u models.User
	err := q.Pool.QueryRow(ctx, `SELECT id,username,password_hash,public_key,user_type,role,is_active,created_at,updated_at FROM users WHERE username=$1 AND public_key=$2 AND is_active=true`, username, publicKey).
		Scan(&u.ID, &u.Username, &u.PasswordHash, &u.PublicKey, &u.UserType, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (q *Queries) AuthPassword(ctx context.Context, username string) (*models.User, error) {
	var u models.User
	err := q.Pool.QueryRow(ctx, `SELECT id,username,password_hash,public_key,user_type,role,is_active,created_at,updated_at FROM users WHERE username=$1 AND is_active=true`, username).
		Scan(&u.ID, &u.Username, &u.PasswordHash, &u.PublicKey, &u.UserType, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (q *Queries) WriteLog(ctx context.Context, l models.Log) error {
	_, err := q.Pool.Exec(ctx, `INSERT INTO logs(id,username,filename,direction,bytes,status,error,mapping_id,session_id,component,duration_ms) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		uuid.NewString(), l.Username, l.Filename, l.Direction, l.Bytes, l.Status, l.Error, nullUUID(l.MappingID), nullUUID(l.SessionID), l.Component, l.DurationMS)
	return err
}

func (q *Queries) GetLogs(ctx context.Context, page, limit int, status, username string) ([]models.Log, error) {
	offset := (page - 1) * limit
	query := `SELECT id,username,filename,direction,bytes,status,error,created_at,mapping_id,session_id,component,duration_ms FROM logs WHERE ($1='' OR status::text=$1) AND ($2='' OR username=$2) ORDER BY created_at DESC LIMIT $3 OFFSET $4`
	rows, err := q.Pool.Query(ctx, query, status, username, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []models.Log{}
	for rows.Next() {
		var l models.Log
		if err := rows.Scan(&l.ID, &l.Username, &l.Filename, &l.Direction, &l.Bytes, &l.Status, &l.Error, &l.CreatedAt, &l.MappingID, &l.SessionID, &l.Component, &l.DurationMS); err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, rows.Err()
}

func (q *Queries) CreateUser(ctx context.Context, u models.User) error {
	_, err := q.Pool.Exec(ctx, `INSERT INTO users(id,username,password_hash,public_key,user_type,role,is_active) VALUES($1,$2,$3,$4,$5,$6,$7)`, uuid.NewString(), u.Username, u.PasswordHash, u.PublicKey, u.UserType, u.Role, u.IsActive)
	if pgErr(err, "users_username_key") {
		return fmt.Errorf("conflict")
	}
	return err
}

func (q *Queries) UpdateUser(ctx context.Context, id string, u models.User) error {
	ct, err := q.Pool.Exec(ctx, `UPDATE users SET password_hash=$2,public_key=$3,user_type=$4,role=$5,is_active=$6,updated_at=NOW() WHERE id=$1`, id, u.PasswordHash, u.PublicKey, u.UserType, u.Role, u.IsActive)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (q *Queries) DeleteUser(ctx context.Context, id string) error {
	ct, err := q.Pool.Exec(ctx, `DELETE FROM users WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (q *Queries) CreateMapping(ctx context.Context, m models.Mapping) error {
	_, err := q.Pool.Exec(ctx, `INSERT INTO mappings(id,external_user,internal_user,flow_type,is_active) VALUES($1,$2,$3,$4,$5)`, uuid.NewString(), m.ExternalUser, m.InternalUser, m.FlowType, m.IsActive)
	return err
}

func (q *Queries) DeleteMapping(ctx context.Context, id string) error {
	ct, err := q.Pool.Exec(ctx, `DELETE FROM mappings WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (q *Queries) ListUsers(ctx context.Context) ([]models.User, error) {
	rows, err := q.Pool.Query(ctx, `SELECT id,username,password_hash,public_key,user_type,role,is_active,created_at,updated_at FROM users ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []models.User{}
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.PublicKey, &u.UserType, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, rows.Err()
}

func (q *Queries) GetUser(ctx context.Context, id string) (*models.User, error) {
	var u models.User
	err := q.Pool.QueryRow(ctx, `SELECT id,username,password_hash,public_key,user_type,role,is_active,created_at,updated_at FROM users WHERE id=$1`, id).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.PublicKey, &u.UserType, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (q *Queries) ListMappings(ctx context.Context) ([]models.Mapping, error) {
	return q.GetActiveMappings(ctx)
}

func pgErr(err error, constraint string) bool {
	if err == nil {
		return false
	}
	return true && (len(constraint) > 0) && (fmt.Sprintf("%v", err) != "")
}
func nullUUID(v string) any {
	if v == "" {
		return nil
	}
	return v
}
