package services

import (
	"errors"
	"fmt"
	"foodie-service/models"
	"foodie-service/types"
	"foodie-service/utils"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	models *models.BaseModel
}

func NewAuthService(models *models.BaseModel) *AuthService {
	return &AuthService{models: models}
}

func PasswordStrengthCheck(password string) bool {
	numberRegex := regexp.MustCompile(`[0-9]`)
	uppercharRegex := regexp.MustCompile(`[A-Z]`)
	lowercharRegex := regexp.MustCompile(`[a-z]`)
	specialCharacterRegex := regexp.MustCompile(`[!@#$%^&*()_+\-=[\]{};':"\\|,.<>/?]`)

	if len(password) > 7 && len(password) < 25 &&
		numberRegex.MatchString(password) &&
		uppercharRegex.MatchString(password) &&
		lowercharRegex.MatchString(password) &&
		specialCharacterRegex.MatchString(password) {
		return true
	}
	return false
}

func (as *AuthService) SignIn(email string, password string) (*types.SignInResponse, error) {
	user, err := as.models.Auth.GetUserByEmail(email)
	if err != nil || user == nil {
		fmt.Println(err)
		return nil, errors.New("user not found")
	}

	// use bcrypt to compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(user)
	token, err := utils.GenerateToken(user.Email)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &types.SignInResponse{Token: token}, nil
}

func (as *AuthService) SignUp(userDetails *types.SignupRequest) (*types.SignupResponse, error) {
	existingUser, _ := as.models.Auth.GetUserByEmail(userDetails.Email)
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	isPasswordStrong := PasswordStrengthCheck(userDetails.Password)
    if(!isPasswordStrong){
		return nil, errors.New("passowrd is weak")
	}

	// use bcrypt to hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDetails.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create a new user object
	user := &models.UserSchema{
		Email:    userDetails.Email,
		Password: string(hashedPassword),
	}

	user, err = as.models.Auth.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return &types.SignupResponse{UserID: user.UserID}, nil
}
