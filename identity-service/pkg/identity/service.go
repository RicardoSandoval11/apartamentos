package identity

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"time"

	"github.com/RicardoSandoval11/apartamentos/identity-service/conf"
	"github.com/RicardoSandoval11/apartamentos/identity-service/entities"
	"github.com/RicardoSandoval11/apartamentos/identity-service/pkg/verification"
	"github.com/RicardoSandoval11/apartamentos/identity-service/utils"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

// Custom errors
var ErrEmailInUse = errors.New("email already exists")
var ErrVerificationCodeExpired = errors.New("code expired")
var ErrFailedAttemptsExceed = errors.New("number of attempts exceeded")
var ErrInvalidVerificationCode = errors.New("invalid verification code")
var ErrIdentityBlocked = errors.New("identity blocked")
var ErrIdentityNotVerified = errors.New("identity is not verified")
var ErrInvalidCredentials = errors.New("invalid credentials")

const maxNumberOfAttempts = 5

// Service types
type CustomTokenClaims struct {
	Uid   uuid.UUID
	Email string
}

type AccessTokenClaims struct {
	jwt.StandardClaims
	CustomTokenClaims
}

// Services inputs
type RegisterInput struct {
	Email    string
	Password string
}

type RegisterOutput struct {
	Success bool
}

type VerifyAccountInput struct {
	Email string
	Code  string
}

type VerifyAccountOutput struct {
	Success bool
}

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	Token        string
	RefreshToken string
}

type RefreshTokenInput struct {
	RefreshToken string
}

type RefreshTokenOutput struct {
	Token        *string
	RefreshToken *string
}

type Service interface {
	Register(ctx context.Context, input RegisterInput) (*RegisterOutput, error)
	VerifyAccount(ctx context.Context, input VerifyAccountInput) (*VerifyAccountOutput, error)
	Login(ctx context.Context, input LoginInput) (*LoginOutput, error)
	RefreshToken(ctx context.Context, input RefreshTokenInput) (*RefreshTokenOutput, error)
}

type identityService struct {
	identityRepository     Repository
	verificationRepository verification.Repository
	signKey                *rsa.PrivateKey
	verifyKey              *rsa.PublicKey
	authTokenDuration      time.Duration
	refreshTokenDuration   time.Duration
}

func NewIdentityService(
	identityRepository Repository,
	verificationRepository verification.Repository,
	config *conf.Config,
) Service {

	privateKeyBytes := []byte(config.PrivateKeyAuthToken)
	publicKeyBytes := []byte(config.PublicKeyAuthToken)

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)

	if err != nil {
		panic("invalid private key")
	}

	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)

	if err != nil {
		panic("invalid public key")
	}

	return &identityService{
		identityRepository:     identityRepository,
		verificationRepository: verificationRepository,
		signKey:                signKey,
		verifyKey:              verifyKey,
		authTokenDuration:      time.Duration(config.AuthTokenDurationMinutes * int(time.Minute)),
		refreshTokenDuration:   time.Duration(config.RefreshTokenDurationMinutes * int(time.Minute)),
	}
}

func (s *identityService) Register(ctx context.Context, input RegisterInput) (*RegisterOutput, error) {

	identity, err := s.identityRepository.GetByEmail(ctx, input.Email)

	if err != nil {
		return nil, err
	}

	if identity != nil {
		return nil, ErrEmailInUse
	}

	newUser := entities.Identity{
		Email:        input.Email,
		PasswordHash: utils.GenerateHash(input.Password),
		PublicId:     uuid.New(),
	}

	identityId, err := s.identityRepository.Create(ctx, newUser)

	if err != nil {
		return nil, err
	}

	code, err := s.verificationRepository.SaveVerification(ctx, *identityId)

	if err != nil {
		return nil, err
	}

	fmt.Printf("Sending verification code: %s\n", code)

	result := RegisterOutput{
		Success: true,
	}

	return &result, nil
}

func (s *identityService) VerifyAccount(ctx context.Context, input VerifyAccountInput) (*VerifyAccountOutput, error) {

	identity, err := s.identityRepository.GetByEmail(ctx, input.Email)

	if err != nil {
		return nil, err
	}

	verification, err := s.verificationRepository.GetVerificationByIdentityId(ctx, identity.Id)

	if err != nil {
		return nil, err
	}

	if verification.Verified {
		return &VerifyAccountOutput{
			Success: true,
		}, nil
	}

	currDate := time.Now()

	if currDate.After(verification.Expiration) {
		return nil, ErrVerificationCodeExpired
	}

	attempts, err := s.verificationRepository.GetAttemptsByVerificationId(ctx, verification.Id)

	if err != nil {
		return nil, err
	}

	if attempts.Times >= maxNumberOfAttempts {
		return nil, ErrFailedAttemptsExceed
	}

	if verification.Code != input.Code {
		s.verificationRepository.IncreaseVerificationAttempts(ctx, verification.Id)
		return nil, ErrInvalidVerificationCode
	}

	identity.Status = entities.IdentityVerified

	s.verificationRepository.MarkVerificationAsVerified(ctx, verification.Id)
	s.identityRepository.Update(ctx, "status", *identity)

	return &VerifyAccountOutput{
		Success: true,
	}, nil
}

func (s *identityService) Login(ctx context.Context, input LoginInput) (*LoginOutput, error) {

	identity, err := s.identityRepository.GetByEmail(ctx, input.Email)

	if err != nil {
		return nil, err
	}

	if identity.Status == entities.IdentityBlocked {
		return nil, ErrIdentityBlocked
	}

	if identity.Status != entities.IdentityVerified {
		return nil, ErrIdentityNotVerified
	}

	receivedPassword := utils.GenerateHash(input.Password)

	if receivedPassword != identity.PasswordHash {
		// TODO: Control number of failed attempts or block account
		return nil, ErrInvalidCredentials
	}

	authToken, err := s.generateAuthToken(identity.PublicId, identity.Email)

	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateRefreshToken(identity.PublicId)

	if err != nil {
		return nil, err
	}

	return &LoginOutput{
		Token:        authToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *identityService) RefreshToken(ctx context.Context, input RefreshTokenInput) (*RefreshTokenOutput, error) {
	// 1. Decodificar refresh token

	// 2. Obtener identity de la base de datos
	identityId, err := s.identityRepository.GetByPublicId()

	// 3. Validar si la cuenta esta bloqueada

	// 4. Validar si la cuenta esta verificada

}

func (s *identityService) generateAuthToken(publicId uuid.UUID, email string) (string, error) {
	jwtParser := jwt.New(jwt.GetSigningMethod("RS256"))

	expiresAt := time.Now().Add(s.authTokenDuration)

	jwtParser.Claims = &AccessTokenClaims{
		jwt.StandardClaims{
			ExpiresAt: expiresAt.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		CustomTokenClaims{
			Uid:   publicId,
			Email: email,
		},
	}

	jwtParser.SignedString(s.signKey)

	return jwtParser.SignedString(s.signKey)
}

func (s *identityService) generateRefreshToken(publicId uuid.UUID) (string, error) {
	jwtParser := jwt.New(jwt.GetSigningMethod("RS256"))

	expiresAt := time.Now().Add(s.refreshTokenDuration)

	jwtParser.Claims = &AccessTokenClaims{
		jwt.StandardClaims{
			ExpiresAt: expiresAt.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		CustomTokenClaims{
			Uid: publicId,
		},
	}

	jwtParser.SignedString(s.signKey)

	return jwtParser.SignedString(s.signKey)
}
