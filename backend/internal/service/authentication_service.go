package service

import (
	"context"
	"github.com/aldian78/go-react-ecommerce/backend/internal/entity"
	"github.com/aldian78/go-react-ecommerce/backend/internal/repository"
	jwt2 "github.com/aldian78/go-react-ecommerce/common/jwt"
	baseutil "github.com/aldian78/go-react-ecommerce/common/utils"
	protoApi "github.com/aldian78/go-react-ecommerce/proto/pb/api"
	"github.com/aldian78/go-react-ecommerce/proto/pb/authentication"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	gocache "github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"go-micro.dev/v4/logger"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"os"
	"strings"
	"time"
)

type IAuthenticationService interface {
	RegisterService(ctx context.Context, request *authentication.RegisterRequest) (*authentication.RegisterResponse, error)
	LoginService(ctx context.Context, request *authentication.LoginRequest) (*authentication.LoginResponse, error)
	LogoutService(ctx context.Context, request *protoApi.APIREQ) (*authentication.LogoutResponse, error)
	ChangePasswordService(ctx context.Context, request *authentication.ChangePasswordRequest) (*authentication.ChangePasswordResponse, error)
	GetProfileService(ctx context.Context, request *authentication.GetProfileRequest) (*authentication.GetProfileResponse, error)
}

type authenticationService struct {
	authRepository repository.IAuthenticationRepository
	cacheService   *gocache.Cache
}

func NewAuthenticationService(authRepository repository.IAuthenticationRepository, cacheService *gocache.Cache) IAuthenticationService {
	return &authenticationService{
		authRepository: authRepository,
		cacheService:   cacheService,
	}
}

func (as *authenticationService) RegisterService(ctx context.Context, request *authentication.RegisterRequest) (*authentication.RegisterResponse, error) {
	if request.Password != request.PasswordConfirmation {
		return &authentication.RegisterResponse{
			Base: baseutil.BadRequestResponse("Password is not matched"),
		}, nil
	}

	user, err := as.authRepository.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}

	if user != nil {
		return &authentication.RegisterResponse{
			Base: baseutil.BadRequestResponse("User already exist"),
		}, nil
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), 10)
	if err != nil {
		return nil, err
	}

	// Insert ke db
	newUser := entity.User{
		Id:        uuid.NewString(),
		FullName:  request.FullName,
		Email:     request.Email,
		Password:  string(hashedPassword),
		RoleCode:  entity.UserRoleCustomer,
		CreatedAt: time.Now(),
		CreatedBy: &request.FullName,
	}
	err = as.authRepository.InsertUser(ctx, &newUser)
	if err != nil {
		return nil, err
	}

	return &authentication.RegisterResponse{
		Base: baseutil.SuccessResponse("User is registered"),
	}, nil
}

func (as *authenticationService) LoginService(ctx context.Context, request *authentication.LoginRequest) (*authentication.LoginResponse, error) {
	user, err := as.authRepository.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return &authentication.LoginResponse{
			Base: baseutil.BadRequestResponse("User is not registered"),
		}, nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated")
		}

		return nil, err
	}

	// generate jwt
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt2.JwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.Id,
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		Email:    user.Email,
		FullName: user.FullName,
		Role:     user.RoleCode,
	})
	secretKey := os.Getenv("JWT_SECRET")
	accessToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return nil, err
	}

	return &authentication.LoginResponse{
		//Base:        utils.SuccessResponse("Login successful"),
		AccessToken: accessToken,
	}, nil
}

func (as *authenticationService) LogoutService(ctx context.Context, request *protoApi.APIREQ) (*authentication.LogoutResponse, error) {
	//ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+request.Headers["Authorization"])
	//
	//jwtToken, err := jwtentity.ParseTokenFromContext(ctx)
	//if err != nil {
	//	panic(err.Error())
	//}

	authHeader := request.Headers["Authorization"]
	if authHeader == "" {
		return nil, baseutil.UnauthenticatedResponse()
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	tokenStr = strings.TrimSpace(tokenStr)

	logger.Info("check tokenStr : ", tokenStr)
	// Parse token
	jwtToken, err := jwt2.ParseToken(tokenStr)
	if err != nil {
		return nil, baseutil.UnauthenticatedResponse()
	}

	logger.Info("check jwtToken : ", jwtToken)

	tokenClaims, err := jwt2.GetClaimsFromContext(ctx)
	if err != nil {
		panic(err.Error())
	}

	as.cacheService.Set(jwtToken, "", time.Duration(tokenClaims.ExpiresAt.Time.Unix()-time.Now().Unix())*time.Second)

	return &authentication.LogoutResponse{
		Base: baseutil.SuccessResponse("Logout success"),
	}, nil
}

func (as *authenticationService) ChangePasswordService(ctx context.Context, request *authentication.ChangePasswordRequest) (*authentication.ChangePasswordResponse, error) {
	// Cek apakah new pass confirmation matched
	if request.NewPassword != request.NewPasswordConfirmation {
		return &authentication.ChangePasswordResponse{
			Base: baseutil.BadRequestResponse("New password is not matched"),
		}, nil
	}

	// Cek apakah old password sama
	jwtToken, err := jwt2.ParseTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}
	claims, err := jwt2.GetClaimsFromToken(jwtToken)
	if err != nil {
		return nil, err
	}

	user, err := as.authRepository.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return &authentication.ChangePasswordResponse{
			Base: baseutil.BadRequestResponse("User doesn't exist"),
		}, nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.OldPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return &authentication.ChangePasswordResponse{
				Base: baseutil.BadRequestResponse("Old password is not matched"),
			}, nil
		}

		return nil, err
	}

	// Update new password ke database
	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), 10)
	if err != nil {
		return nil, err
	}
	err = as.authRepository.UpdateUserPassword(ctx, user.Id, string(hashedNewPassword), user.FullName)
	if err != nil {
		return nil, err
	}

	// Kirim response
	return &authentication.ChangePasswordResponse{
		Base: baseutil.SuccessResponse("Change password success"),
	}, nil
}

func (as *authenticationService) GetProfileService(ctx context.Context, request *authentication.GetProfileRequest) (*authentication.GetProfileResponse, error) {
	claims, err := jwt2.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	user, err := as.authRepository.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return &authentication.GetProfileResponse{
			Base: baseutil.BadRequestResponse("User doesn't exist"),
		}, nil
	}

	return &authentication.GetProfileResponse{
		Base:        baseutil.SuccessResponse("Get profile success"),
		UserId:      claims.Subject,
		FullName:    claims.FullName,
		Email:       claims.Email,
		RoleCode:    claims.Role,
		MemberSince: timestamppb.New(user.CreatedAt),
	}, nil
}
