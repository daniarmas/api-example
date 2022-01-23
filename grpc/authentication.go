package app

import (
	"context"

	"github.com/daniarmas/api-example/dto"
	pb "github.com/daniarmas/api-example/pkg"
	"github.com/daniarmas/api-example/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	gp "google.golang.org/protobuf/types/known/emptypb"
)

func (m *AuthenticationServer) SignIn(ctx context.Context, req *pb.SignInRequest) (*pb.SignInResponse, error) {
	var st *status.Status
	result, err := m.authenticationService.SignIn(&dto.SignInRequest{Password: req.Password, Email: req.Email})
	if err != nil {
		switch err.Error() {
		case "user not found":
			st = status.New(codes.NotFound, "User not found")
		case "password incorrect":
			st = status.New(codes.PermissionDenied, "Credentials incorrect")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return &pb.SignInResponse{RefreshToken: result.RefreshToken, AuthorizationToken: result.AuthorizationToken, User: &pb.User{Id: result.User.ID.String(), Email: result.User.Email, CreateTime: result.User.CreateTime.String(), UpdateTime: result.User.UpdateTime.String()}}, nil
}

func (m *AuthenticationServer) SignOut(ctx context.Context, req *gp.Empty) (*gp.Empty, error) {
	var st *status.Status
	md, _ := metadata.FromIncomingContext(ctx)
	metadata := utils.Metadata{
		Authorization:            md.Get("authorization")[0],
		Platform:                 md.Get("platform")[0],
		SystemVersion:            md.Get("systemversion")[0],
		SystemLanguage:           md.Get("systemlanguage")[0],
		NetworkType:              md.Get("networktype")[0],
		AppVersion:               md.Get("appversion")[0],
		App:                      md.Get("app")[0],
		DeviceId:                 md.Get("deviceid")[0],
		Ipv4:                     md.Get("ipv4")[0],
		Ipv6:                     md.Get("ipv6")[0],
		Model:                    md.Get("model")[0],
		FirebaseCloudMessagingId: md.Get("firebasecloudmessagingid")[0],
	}
	err := m.authenticationService.SignOut(&metadata)
	if err != nil {
		switch err.Error() {
		case "unauthenticated":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "permission denied":
			st = status.New(codes.PermissionDenied, "Permission denied")
		case "user not found":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "authorizationtoken expired":
			st = status.New(codes.Unauthenticated, "AuthorizationToken expired")
		case "signature is invalid":
			st = status.New(codes.Unauthenticated, "AuthorizationToken invalid")
		case "token contains an invalid number of segments":
			st = status.New(codes.Unauthenticated, "AuthorizationToken invalid")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return &gp.Empty{}, nil
}

func (m *AuthenticationServer) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	var st *status.Status
	md, _ := metadata.FromIncomingContext(ctx)
	metadata := utils.Metadata{
		Authorization:            md.Get("authorization")[0],
		Platform:                 md.Get("platform")[0],
		SystemVersion:            md.Get("systemversion")[0],
		SystemLanguage:           md.Get("systemlanguage")[0],
		NetworkType:              md.Get("networktype")[0],
		AppVersion:               md.Get("appversion")[0],
		App:                      md.Get("app")[0],
		DeviceId:                 md.Get("deviceid")[0],
		Ipv4:                     md.Get("ipv4")[0],
		Ipv6:                     md.Get("ipv6")[0],
		Model:                    md.Get("model")[0],
		FirebaseCloudMessagingId: md.Get("firebasecloudmessagingid")[0],
	}
	result, err := m.authenticationService.RefreshToken(&req.RefreshToken, &metadata)
	if err != nil {
		switch err.Error() {
		case "unauthenticated":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "permission denied":
			st = status.New(codes.PermissionDenied, "Permission denied")
		case "user not found":
			st = status.New(codes.Unauthenticated, "Unauthenticated")
		case "refreshtoken expired":
			st = status.New(codes.Unauthenticated, "RefreshToken expired")
		case "signature is invalid":
			st = status.New(codes.Unauthenticated, "RefreshToken invalid")
		case "token contains an invalid number of segments":
			st = status.New(codes.Unauthenticated, "RefreshToken invalid")
		default:
			st = status.New(codes.Internal, "Internal server error")
		}
		return nil, st.Err()
	}
	return &pb.RefreshTokenResponse{RefreshToken: result.RefreshToken, AuthorizationToken: result.AuthorizationToken}, nil
}
