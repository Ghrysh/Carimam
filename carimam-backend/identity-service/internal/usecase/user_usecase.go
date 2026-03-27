package usecase

import (
	"errors"

	"github.com/ghrysh/carimam/identity-service/internal/models"
	"github.com/ghrysh/carimam/identity-service/internal/repository"
	"github.com/ghrysh/carimam/identity-service/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Name     string `json:"name" form:"name" binding:"required"`
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required,min=6"`
	Phone    string `json:"phone" form:"phone" binding:"required"`
	Role     string `json:"role" form:"role"` 
}

type LoginRequest struct {
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required"`
}

// DTO BARU UNTUK POTONG SALDO
type DeductBalanceRequest struct {
	UserID uint    `json:"user_id" binding:"required"`
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

type UserUseCase interface {
	Register(req RegisterRequest) error
	Login(req LoginRequest) (string, error)
	DeductBalance(req DeductBalanceRequest) error
	GetProfile(userID uint) (*models.User, error)
}

type userUseCase struct {
	repo repository.UserRepository
}

func NewUserUseCase(repo repository.UserRepository) UserUseCase {
	return &userUseCase{repo}
}

func (u *userUseCase) Register(req RegisterRequest) error {
	existingUser, _ := u.repo.FindByEmail(req.Email)
	if existingUser != nil {
		return errors.New("email sudah terdaftar, silakan gunakan email lain")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("gagal memproses password")
	}

	role := req.Role
	if role == "" {
		role = "eater"
	}

	newUser := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Phone:    req.Phone,
		Role:     role,
	}

	return u.repo.CreateUser(newUser)
}

func (u *userUseCase) Login(req LoginRequest) (string, error) {
	user, err := u.repo.FindByEmail(req.Email)
	if err != nil {
		return "", errors.New("email atau password salah") 
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", errors.New("email atau password salah")
	}

	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", errors.New("gagal membuat token akses")
	}

	return token, nil
}

// ===================================
// FUNGSI BARU: Logika Potong Saldo
// ===================================
func (u *userUseCase) DeductBalance(req DeductBalanceRequest) error {
	user, err := u.repo.FindByID(req.UserID)
	if err != nil {
		return errors.New("user tidak ditemukan")
	}

	if user.Balance < req.Amount {
		return errors.New("saldo tidak mencukupi untuk melakukan pesanan ini")
	}

	user.Balance -= req.Amount
	return u.repo.Update(user)
}

func (u *userUseCase) GetProfile(userID uint) (*models.User, error) {
	return u.repo.FindByID(userID)
}