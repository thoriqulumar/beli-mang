package service

import (
	"beli-mang/config"
	"beli-mang/model"
	"beli-mang/pkg/crypto"
	cerr "beli-mang/pkg/customErr"
	"beli-mang/repo"
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type StaffService interface {
	RegisterAdmin(ctx context.Context, newStaff model.Staff) (model.StaffWithToken, error)
	LoginAdmin(ctx context.Context, staff model.Staff) (model.StaffWithToken, error)
	RegisterUser(ctx context.Context, newStaff model.Staff) (model.StaffWithToken, error)
	LoginUser(ctx context.Context, staff model.Staff) (model.StaffWithToken, error)
}

type staffSvc struct {
	cfg  *config.Config
	repo repo.StaffRepo
}

func NewStaffService(cfg *config.Config, r repo.StaffRepo) StaffService {
	return &staffSvc{
		cfg:  cfg,
		repo: r,
	}
}

func (s *staffSvc) RegisterAdmin(ctx context.Context, newStaff model.Staff) (model.StaffWithToken, error) {
	usernameStaff, err := s.repo.GetStaffByUsername(ctx, newStaff.Username)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return model.StaffWithToken{}, cerr.New(http.StatusInternalServerError, err.Error())
		}
	}
	if usernameStaff.Username != "" {
		return model.StaffWithToken{}, cerr.New(http.StatusConflict, "Username already in use")
	}

	staff, err := s.repo.GetStaffByEmail(ctx, newStaff.Email, "admin")
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return model.StaffWithToken{}, cerr.New(http.StatusInternalServerError, err.Error())

		}
	}

	if staff.Email != "" {
		return model.StaffWithToken{}, cerr.New(http.StatusConflict, "Email conflict with another admin")
	}

	hashedPassword, err := crypto.GenerateHashedPassword(newStaff.Password, s.cfg.BcryptSalt)
	if err != nil {
		return model.StaffWithToken{}, err
	}
	newStaff.Password = hashedPassword

	id := uuid.New()
	newStaff.ID = id
	newStaff.CreatedAt = time.Now()
	newStaff.Role = "admin"

	//save to database
	err = s.repo.InsertStaff(ctx, newStaff, hashedPassword)
	if err != nil {
		return model.StaffWithToken{}, err
	}
	// Generate token
	token, err := crypto.GenerateToken(newStaff, s.cfg.JWTSecret)
	if err != nil {
		return model.StaffWithToken{}, err

	}
	serviceResponse := model.StaffWithToken{
		AccessToken: token,
	}
	return serviceResponse, nil
}

func (s *staffSvc) LoginAdmin(ctx context.Context, staff model.Staff) (model.StaffWithToken, error) {
	staffAdmin, err := s.repo.GetStaffByUsername(ctx, staff.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.StaffWithToken{}, cerr.New(http.StatusBadRequest, "user not found")
		}
		return model.StaffWithToken{}, cerr.New(http.StatusInternalServerError, "database error: "+err.Error())
	}
	err = crypto.VerifyPassword(staff.Password, staffAdmin.Password)
	if err != nil {
		return model.StaffWithToken{}, cerr.New(http.StatusBadRequest, "Invalid password")
	}
	token, err := crypto.GenerateToken(staffAdmin, s.cfg.JWTSecret)
	if err != nil {
		return model.StaffWithToken{}, cerr.New(http.StatusBadRequest, err.Error())
	}
	serviceResponse := model.StaffWithToken{
		AccessToken: token,
	}
	return serviceResponse, nil
}

func (s *staffSvc) RegisterUser(ctx context.Context, newStaff model.Staff) (model.StaffWithToken, error) {
	usernameStaff, err := s.repo.GetStaffByUsername(ctx, newStaff.Username)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return model.StaffWithToken{}, cerr.New(http.StatusInternalServerError, err.Error())
		}
	}
	if usernameStaff.Username != "" {
		return model.StaffWithToken{}, cerr.New(http.StatusConflict, "Username already in use")
	}

	staff, err := s.repo.GetStaffByEmail(ctx, newStaff.Email, "user")
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return model.StaffWithToken{}, cerr.New(http.StatusInternalServerError, err.Error())

		}
	}

	if staff.Email != "" {
		return model.StaffWithToken{}, cerr.New(http.StatusConflict, "Email conflict with another user")
	}

	hashedPassword, err := crypto.GenerateHashedPassword(newStaff.Password, s.cfg.BcryptSalt)
	if err != nil {
		return model.StaffWithToken{}, err
	}
	newStaff.Password = hashedPassword

	id := uuid.New()
	newStaff.ID = id
	newStaff.CreatedAt = time.Now()
	newStaff.Role = "user"

	//save to database
	err = s.repo.InsertStaff(ctx, newStaff, hashedPassword)
	if err != nil {
		return model.StaffWithToken{}, err
	}
	// Generate token
	token, err := crypto.GenerateToken(newStaff, s.cfg.JWTSecret)
	if err != nil {
		return model.StaffWithToken{}, err

	}
	serviceResponse := model.StaffWithToken{
		AccessToken: token,
	}
	return serviceResponse, nil
}

func (s *staffSvc) LoginUser(ctx context.Context, staff model.Staff) (model.StaffWithToken, error) {
	staffAdmin, err := s.repo.GetStaffByUsername(ctx, staff.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.StaffWithToken{}, cerr.New(http.StatusBadRequest, "user not found")
		}
		return model.StaffWithToken{}, cerr.New(http.StatusInternalServerError, "database error: "+err.Error())
	}
	err = crypto.VerifyPassword(staff.Password, staffAdmin.Password)
	if err != nil {
		return model.StaffWithToken{}, cerr.New(http.StatusBadRequest, "Invalid password")
	}
	token, err := crypto.GenerateToken(staffAdmin, s.cfg.JWTSecret)
	if err != nil {
		return model.StaffWithToken{}, cerr.New(http.StatusBadRequest, err.Error())
	}
	serviceResponse := model.StaffWithToken{
		AccessToken: token,
	}
	return serviceResponse, nil
}
