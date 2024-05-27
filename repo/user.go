package repo

import (
	"beli-mang/model"
	"context"

	"github.com/jmoiron/sqlx"
)

type StaffRepo interface {
	InsertStaff(ctx context.Context, staff model.Staff, hashPassword string) error
	GetStaffByEmail(ctx context.Context, email, role string) (model.Staff, error)
	GetStaffByUsername(ctx context.Context, username string) (model.Staff, error)
}

type staffRepo struct {
	db *sqlx.DB
}

func NewStaffRepo(db *sqlx.DB) StaffRepo {
	return &staffRepo{db: db}
}

func (r *staffRepo) InsertStaff(ctx context.Context, staff model.Staff, hashPassword string) error {

	query := `INSERT INTO "user" (id, username, role, password, email, "createdAt") VALUES ($1, $2, $3, $4, $5, NOW())`
	_, err := r.db.ExecContext(ctx, query, staff.ID, staff.Username, staff.Role, hashPassword, staff.Email)
	if err != nil {
		return err
	}
	return nil
}

func (r *staffRepo) GetStaffByEmail(ctx context.Context, email, role string) (model.Staff, error) {
	var staff model.Staff
	query := `SELECT id, username, role, email, "createdAt" FROM "user" WHERE email = $1 AND role = $2 LIMIT 1`
	err := r.db.GetContext(ctx, &staff, query, email, role)
	if err != nil {
		return model.Staff{}, err
	}
	return staff, nil

}

func (r *staffRepo) GetStaffByUsername(ctx context.Context, username string) (model.Staff, error) {
	var staff model.Staff
	query := `SELECT id, username, role, email, password, "createdAt" FROM "user" WHERE username = $1`
	err := r.db.Get(&staff, query, username)
	return staff, err
}
