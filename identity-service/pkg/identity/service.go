package identity

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/RicardoSandoval11/apartamentos/identity-service/conf"
	"github.com/RicardoSandoval11/apartamentos/identity-service/entities"
	"github.com/RicardoSandoval11/apartamentos/identity-service/pkg/verification"
	"github.com/RicardoSandoval11/apartamentos/identity-service/utils"
	"github.com/golang-jwt/jwt"
)

// Custom errors
var ErrEmailInUse = errors.New("email already exists")
var ErrVerificationCodeExpired = errors.New("code expired")
var ErrFailedAttemptsExceed = errors.New("number of attempts exceeded")
var ErrInvalidVerificationCode = errors.New("invalid verification code")
var ErrIdentityBlocked = errors.New("identity blocked")
var ErrIdentityNotVerified = errors.New("identity is not verified")

const maxNumberOfAttempts = 5

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
	Token        *string
	RefreshToken *string
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
	config                 *conf.Config
}

func NewIdentityService(
		identityRepository Repository, 
		verificationRepository verification.Repository,
		config &conf.Config // SEGUIR AQUI
	) Service {
	return &identityService{
		identityRepository:     identityRepository,
		verificationRepository: verificationRepository,
	}
}

func (s *identityService) Register(ctx context.Context, input RegisterInput) (*RegisterOutput, error) {

	// 1. Validar si el correo ya existe
	identity, err := s.identityRepository.GetByEmail(ctx, input.Email)

	if err != nil {
		return nil, err
	}

	// 2. Si el correo ya existe, retornar error
	if identity != nil {
		return nil, ErrEmailInUse
	}

	// 3. Si el correo no existe, crear registro
	newUser := entities.Identity{
		Email:        input.Email,
		PasswordHash: utils.GenerateHash(input.Password),
	}

	identityId, err := s.identityRepository.Create(ctx, newUser)

	if err != nil {
		return nil, err
	}

	// 4. Generate and save verification code
	code, err := s.verificationRepository.SaveVerification(ctx, *identityId)

	if err != nil {
		return nil, err
	}

	// 5. Enviar correo de verificacion de cuenta
	fmt.Printf("Sending verification code: %s\n", code)

	// 6. Retornar respuesta de exito
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
	// 1. Verificar si la cuenta esta verificada o bloqueada
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

	jwtParser = jwt.New(jwt.GetSigningMethod("RS256"))
}

func (s *identityService) RefreshToken(ctx context.Context, input RefreshTokenInput) (*RefreshTokenOutput, error) {

}
