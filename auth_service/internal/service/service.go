package service

import (
	"context"
	"fmt"
	"github.com/paw1a/grpc-media-converter/auth_service/config"
	"github.com/paw1a/grpc-media-converter/auth_service/internal/domain"
	"github.com/paw1a/grpc-media-converter/auth_service/internal/repository"
	"github.com/paw1a/grpc-media-converter/auth_service/pb"
	"github.com/paw1a/grpc-media-converter/auth_service/pkg/utils"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcService struct {
	pb.UnimplementedAuthServiceServer
	cfg      *config.Config
	userRepo repository.Users
	jwt      *utils.JwtProvider
}

func NewGrpcService(cfg *config.Config, userRepo repository.Users, jwt *utils.JwtProvider) *GrpcService {
	return &GrpcService{cfg: cfg, userRepo: userRepo, jwt: jwt}
}

func (s *GrpcService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	var user domain.User

	_, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err == nil {
		return nil, s.errResponse(codes.AlreadyExists,
			fmt.Errorf("user with email %v already exists", req.Email))
	}

	user.Email = req.Email
	user.Password = utils.HashPassword(req.Password)

	_, err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, s.errResponse(codes.Internal, err)
	}

	return &pb.RegisterResponse{}, nil
}

func (s *GrpcService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, s.errResponse(codes.NotFound,
			fmt.Errorf("user with email %v not found", req.Email))
	}

	match := utils.CheckPasswordHash(req.Password, user.Password)
	if !match {
		return nil, s.errResponse(codes.InvalidArgument, errors.New("invalid password"))
	}

	token, err := s.jwt.GenerateToken(user)
	if err != nil {
		return nil, s.errResponse(codes.Internal, err)
	}

	return &pb.LoginResponse{Token: token}, nil
}

func (s *GrpcService) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	claims, err := s.jwt.ValidateToken(req.Token)
	if err != nil {
		return nil, s.errResponse(codes.Internal, err)
	}

	user, err := s.userRepo.FindByID(ctx, claims.Id)
	if err != nil {
		return nil, s.errResponse(codes.Internal, err)
	}

	return &pb.ValidateResponse{UserId: user.Id}, nil
}

func (s *GrpcService) errResponse(c codes.Code, err error) error {
	return status.Error(c, err.Error())
}
