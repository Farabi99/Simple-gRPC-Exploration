package controllers

import (
	"context"

	"github.com/Farabi99/Simple-gRPC-Exploration/models"
	pb "github.com/Farabi99/Simple-gRPC-Exploration/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserController implements the generated gRPC server interface
type UserController struct {
	pb.UnimplementedUserServiceServer
	repo models.UserRepository // Inject the model
}

func NewUserController(repo models.UserRepository) *UserController {
	return &UserController{repo: repo}
}

func (c *UserController) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	user, err := c.repo.Create(req.GetName(), req.GetEmail())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	return &pb.UserResponse{
		User: &pb.User{
			Id:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	}, nil
}

func (c *UserController) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	user, err := c.repo.GetByID(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}

	return &pb.UserResponse{
		User: &pb.User{
			Id:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	}, nil
}

func (c *UserController) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	user, err := c.repo.Update(req.GetId(), req.GetName(), req.GetEmail())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}

	return &pb.UserResponse{
		User: &pb.User{
			Id:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	}, nil
}

func (c *UserController) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	err := c.repo.Delete(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}

	return &pb.DeleteUserResponse{Success: true}, nil
}

func (c *UserController) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	// Provide a default limit if the client doesn't send one
	limit := int(req.GetLimit())
	if limit <= 0 || limit > 100 {
		limit = 10 // default chunk size
	}

	users, nextCursor, err := c.repo.List(limit, req.GetCursor())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	// Map the Model Users to Protobuf Users
	var pbUsers []*pb.User
	for _, u := range users {
		pbUsers = append(pbUsers, &pb.User{
			Id:    u.ID,
			Name:  u.Name,
			Email: u.Email,
		})
	}

	return &pb.ListUsersResponse{
		Users:      pbUsers,
		NextCursor: nextCursor,
	}, nil
}
