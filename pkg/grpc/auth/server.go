package auth

import (
	"context"
	"errors"
	"github.com/Pinkman-77/grpc-auth/pkg/service/auth"
	"github.com/Pinkman-77/grpc-auth/pkg/storage"

	ssov1 "github.com/Pinkman-77/grpc-auth/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context, email string, password string, appID int) (token string, err error)
	Register(ctx context.Context, email string, password string) (userID int64, err error)
	Admin(ctx context.Context, userID int64) (bool, error)
}

type Server struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth *auth.Auth) {
	ssov1.RegisterAuthServer(gRPC, &Server{auth: auth})
}

const emptyV = 0

func (s *Server) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	if err := validateLogin(req); err != nil {
		return nil, err
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &ssov1.LoginResponse{
		Token: token,
	}, nil

}

func (s *Server) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	if err := validateRegister(req); err != nil {
		return nil, err
	}

	userID, err := s.auth.Register(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &ssov1.RegisterResponse{UserId: userID}, nil

}

func (s *Server) Admin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	if err := validateAdmin(req); err != nil {
		return nil, err
	}

	isAdmin, err := s.auth.Admin(ctx, req.GetUserId())

	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &ssov1.IsAdminResponse{IsAdmin: isAdmin}, nil
}
func validateLogin(req *ssov1.LoginRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is not provided")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is not provided")
	}

	if req.GetAppId() == emptyV {
		return status.Error(codes.InvalidArgument, "app id is not provided")
	}

	return nil
}
func validateRegister(req *ssov1.RegisterRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is not provided")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is not provided")
	}

	return nil
}

func validateAdmin(req *ssov1.IsAdminRequest) error {
	if req.GetUserId() == emptyV {
		return status.Error(codes.InvalidArgument, "user id is not provided")
	}

	return nil
}
