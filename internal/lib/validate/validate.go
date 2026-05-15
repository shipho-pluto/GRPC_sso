package validate

import (
	ssov1 "github.com/shipho-pluto/grpc_proto/gen/go/sso"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	emptyValueAppID = 0
)

func ValidateRegister(req *ssov1.RegisterRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	return nil
}

func ValidateLogin(req *ssov1.LoginRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	if req.GetAppId() == emptyValueAppID {
		return status.Error(codes.InvalidArgument, "app_id is required")
	}

	return nil
}

func ValidateIsAdmin(req *ssov1.IsAdminRequest) error {
	if req.GetUserId() <= emptyValueAppID {
		return status.Error(codes.InvalidArgument, "user_id is required")
	}

	return nil
}

func ValidateLogout(req *ssov1.LogoutRequest) error {
	if req.GetUserId() <= emptyValueAppID {
		return status.Error(codes.InvalidArgument, "user_id is required")
	}

	return nil
}
