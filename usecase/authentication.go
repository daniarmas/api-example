package usecase

import (
	"errors"

	"github.com/daniarmas/api-example/dto"
	"github.com/daniarmas/api-example/models"
	"github.com/daniarmas/api-example/repository"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
	"gorm.io/gorm"
)

type AuthenticationService interface {
	SignIn(signInRequest *dto.SignInRequest, metadata *metadata.MD) (*dto.SignInResponse, error)
	SignOut(all *bool, authorizationTokenFk *string, metadata *metadata.MD) error
	RefreshToken(refreshToken *string, metadata *metadata.MD) (*dto.RefreshToken, error)
}

type authenticationService struct {
	dao repository.DAO
}

func NewAuthenticationService(dao repository.DAO) AuthenticationService {
	return &authenticationService{dao: dao}
}

func (v *authenticationService) SignIn(signInRequest *dto.SignInRequest, metadata *metadata.MD) (*dto.SignInResponse, error) {
	var userRes *models.User
	var userErr, refreshTokenErr, authorizationTokenErr, jwtRefreshTokenErr, jwtAuthorizationTokenErr error
	var refreshTokenRes *models.RefreshToken
	var authorizationTokenRes *models.AuthorizationToken
	var jwtAuthorizationTokenRes, jwtRefreshTokenRes *string
	err := repository.DB.Transaction(func(tx *gorm.DB) error {
		userRes, userErr = v.dao.NewUserQuery().GetUser(tx, &models.User{Email: signInRequest.Email}, nil)
		if userErr != nil {
			return userErr
		} else if *userRes == (models.User{}) {
			return errors.New("user not found")
		}
		password := v.dao.NewHashPasswordQuery().CheckPasswordHash(signInRequest.Password, userRes.Password)
		if !password {
			return errors.New("password incorrect")
		}
		deleteRefreshTokenRes, deleteRefreshTokenErr := v.dao.NewRefreshTokenQuery().DeleteRefreshToken(tx, &models.RefreshToken{UserFk: userRes.ID}, &[]string{"id"})
		if deleteRefreshTokenErr != nil {
			return deleteRefreshTokenErr
		}
		if len(*deleteRefreshTokenRes) != 0 {
			_, deleteAuthorizationTokenErr := v.dao.NewAuthorizationTokenQuery().DeleteAuthorizationToken(tx, &models.AuthorizationToken{RefreshTokenFk: (*deleteRefreshTokenRes)[0].ID})
			if deleteAuthorizationTokenErr != nil {
				return deleteAuthorizationTokenErr
			}
		}
		refreshTokenRes, refreshTokenErr = v.dao.NewRefreshTokenQuery().CreateRefreshToken(tx, &models.RefreshToken{UserFk: userRes.ID})
		if refreshTokenErr != nil {
			return refreshTokenErr
		}
		authorizationTokenRes, authorizationTokenErr = v.dao.NewAuthorizationTokenQuery().CreateAuthorizationToken(tx, &models.AuthorizationToken{RefreshTokenFk: refreshTokenRes.ID, UserFk: userRes.ID})
		if authorizationTokenErr != nil {
			return authorizationTokenErr
		}
		authorizationTokenId := authorizationTokenRes.ID.String()
		refreshTokenId := refreshTokenRes.ID.String()
		jwtRefreshTokenRes, jwtRefreshTokenErr = v.dao.NewTokenQuery().CreateJwtRefreshToken(&refreshTokenId)
		if jwtRefreshTokenErr != nil {
			return jwtRefreshTokenErr
		}
		jwtAuthorizationTokenRes, jwtAuthorizationTokenErr = v.dao.NewTokenQuery().CreateJwtAuthorizationToken(&authorizationTokenId)
		if jwtAuthorizationTokenErr != nil {
			return jwtAuthorizationTokenErr
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &dto.SignInResponse{AuthorizationToken: *jwtAuthorizationTokenRes, RefreshToken: *jwtRefreshTokenRes, User: *userRes}, nil
}

func (v *authenticationService) SignOut(all *bool, authorizationTokenFk *string, metadata *metadata.MD) error {
	err := repository.DB.Transaction(func(tx *gorm.DB) error {
		if *all {
			authorizationTokenParseRes, authorizationTokenParseErr := v.dao.NewTokenQuery().ParseJwtAuthorizationToken(&metadata.Get("authorization")[0])
			if authorizationTokenParseErr != nil {
				switch authorizationTokenParseErr.Error() {
				case "Token is expired":
					return errors.New("authorizationtoken expired")
				case "signature is invalid":
					return errors.New("signature is invalid")
				case "token contains an invalid number of segments":
					return errors.New("token contains an invalid number of segments")
				default:
					return authorizationTokenParseErr
				}
			}
			authorizationTokenRes, authorizationTokenErr := v.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: uuid.MustParse(*authorizationTokenParseRes)}, &[]string{"id", "user_fk"})
			if authorizationTokenErr != nil {
				return authorizationTokenErr
			} else if *authorizationTokenRes == (models.AuthorizationToken{}) {
				return errors.New("unauthenticated")
			}
			var refreshTokenIds []string
			deleteRefreshTokenRes, deleteRefreshTokenErr := v.dao.NewRefreshTokenQuery().DeleteRefreshToken(tx, &models.RefreshToken{UserFk: authorizationTokenRes.UserFk}, &[]string{"id"})
			if deleteRefreshTokenErr != nil {
				return deleteRefreshTokenErr
			}
			for _, e := range *deleteRefreshTokenRes {
				refreshTokenIds = append(refreshTokenIds, e.ID.String())
			}
			_, deleteAuthorizationTokenErr := v.dao.NewAuthorizationTokenQuery().DeleteAuthorizationTokenIn(tx, "refresh_token_fk IN ?", &refreshTokenIds)
			if deleteAuthorizationTokenErr != nil {
				return deleteAuthorizationTokenErr
			}
			return nil
		} else if *authorizationTokenFk != "" {
			authorizationTokenParseRes, authorizationTokenParseErr := v.dao.NewTokenQuery().ParseJwtAuthorizationToken(&metadata.Get("authorization")[0])
			if authorizationTokenParseErr != nil {
				switch authorizationTokenParseErr.Error() {
				case "Token is expired":
					return errors.New("authorizationtoken expired")
				case "signature is invalid":
					return errors.New("signature is invalid")
				case "token contains an invalid number of segments":
					return errors.New("token contains an invalid number of segments")
				default:
					return authorizationTokenParseErr
				}
			}
			authorizationTokenByReqRes, authorizationTokenByReqErr := v.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: uuid.MustParse(*authorizationTokenFk)}, &[]string{"id", "user_fk", "device_fk"})
			if authorizationTokenByReqErr != nil {
				return authorizationTokenByReqErr
			}
			authorizationTokenRes, authorizationTokenErr := v.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: uuid.MustParse(*authorizationTokenParseRes)}, &[]string{"id", "user_fk"})
			if authorizationTokenErr != nil {
				return authorizationTokenErr
			} else if *authorizationTokenRes == (models.AuthorizationToken{}) {
				return errors.New("unauthenticated")
			} else if authorizationTokenRes.UserFk != authorizationTokenByReqRes.UserFk {
				return errors.New("permission denied")
			}
			deleteRefreshTokenRes, deleteRefreshTokenErr := v.dao.NewRefreshTokenQuery().DeleteRefreshToken(tx, &models.RefreshToken{UserFk: authorizationTokenByReqRes.UserFk}, &[]string{"id"})
			if deleteRefreshTokenErr != nil {
				return deleteRefreshTokenErr
			}
			_, deleteAuthorizationTokenErr := v.dao.NewAuthorizationTokenQuery().DeleteAuthorizationToken(tx, &models.AuthorizationToken{RefreshTokenFk: (*deleteRefreshTokenRes)[0].ID})
			if deleteAuthorizationTokenErr != nil {
				return deleteAuthorizationTokenErr
			}
			return nil
		} else {
			authorizationTokenParseRes, authorizationTokenParseErr := v.dao.NewTokenQuery().ParseJwtAuthorizationToken(&metadata.Get("authorization")[0])
			if authorizationTokenParseErr != nil {
				switch authorizationTokenParseErr.Error() {
				case "Token is expired":
					return errors.New("authorizationtoken expired")
				case "signature is invalid":
					return errors.New("signature is invalid")
				case "token contains an invalid number of segments":
					return errors.New("token contains an invalid number of segments")
				default:
					return authorizationTokenParseErr
				}
			}
			authorizationTokenRes, authorizationTokenErr := v.dao.NewAuthorizationTokenQuery().GetAuthorizationToken(tx, &models.AuthorizationToken{ID: uuid.MustParse(*authorizationTokenParseRes)}, &[]string{"id", "user_fk", "device_fk"})
			if authorizationTokenErr != nil {
				return authorizationTokenErr
			} else if *authorizationTokenRes == (models.AuthorizationToken{}) {
				return errors.New("unauthenticated")
			}
			deleteRefreshTokenRes, deleteRefreshTokenErr := v.dao.NewRefreshTokenQuery().DeleteRefreshToken(tx, &models.RefreshToken{UserFk: authorizationTokenRes.UserFk}, &[]string{"id"})
			if deleteRefreshTokenErr != nil {
				return deleteRefreshTokenErr
			}
			_, deleteAuthorizationTokenErr := v.dao.NewAuthorizationTokenQuery().DeleteAuthorizationToken(tx, &models.AuthorizationToken{RefreshTokenFk: (*deleteRefreshTokenRes)[0].ID})
			if deleteAuthorizationTokenErr != nil {
				return deleteAuthorizationTokenErr
			}
			return nil
		}
	})
	if err != nil {
		return err
	}
	return nil
}

func (v *authenticationService) RefreshToken(refreshToken *string, metadata *metadata.MD) (*dto.RefreshToken, error) {
	var jwtAuthorizationTokenRes, jwtRefreshTokenRes *string
	var jwtAuthorizationTokenErr, jwtRefreshTokenErr error
	err := repository.DB.Transaction(func(tx *gorm.DB) error {
		refreshTokenParseRes, refreshTokenParseErr := v.dao.NewTokenQuery().ParseJwtRefreshToken(refreshToken)
		if refreshTokenParseErr != nil {
			switch refreshTokenParseErr.Error() {
			case "Token is expired":
				return errors.New("refreshtoken expired")
			case "signature is invalid":
				return errors.New("signature is invalid")
			case "token contains an invalid number of segments":
				return errors.New("token contains an invalid number of segments")
			default:
				return refreshTokenParseErr
			}
		}
		refreshTokenRes, refreshTokenErr := v.dao.NewRefreshTokenQuery().GetRefreshToken(tx, &models.RefreshToken{ID: uuid.MustParse(*refreshTokenParseRes)}, &[]string{"id", "user_fk", "device_fk"})
		if refreshTokenErr != nil {
			return refreshTokenErr
		} else if *refreshTokenRes == (models.RefreshToken{}) {
			return errors.New("unauthenticated")
		}
		userRes, userErr := v.dao.NewUserQuery().GetUser(tx, &models.User{ID: refreshTokenRes.UserFk}, &[]string{"id"})
		if userErr != nil {
			return userErr
		} else if *userRes == (models.User{}) {
			return errors.New("user not found")
		}
		deleteRefreshTokenRes, deleteRefreshTokenErr := v.dao.NewRefreshTokenQuery().DeleteRefreshToken(tx, &models.RefreshToken{ID: refreshTokenRes.ID}, &[]string{"id"})
		if deleteRefreshTokenErr != nil {
			return deleteRefreshTokenErr
		}
		_, deleteAuthorizationTokenErr := v.dao.NewAuthorizationTokenQuery().DeleteAuthorizationToken(tx, &models.AuthorizationToken{RefreshTokenFk: (*deleteRefreshTokenRes)[0].ID})
		if deleteAuthorizationTokenErr != nil {
			return deleteAuthorizationTokenErr
		}
		refreshTokenRes, refreshTokenErr = v.dao.NewRefreshTokenQuery().CreateRefreshToken(tx, &models.RefreshToken{UserFk: userRes.ID})
		if refreshTokenErr != nil {
			return refreshTokenErr
		}
		authorizationTokenRes, authorizationTokenErr := v.dao.NewAuthorizationTokenQuery().CreateAuthorizationToken(tx, &models.AuthorizationToken{RefreshTokenFk: refreshTokenRes.ID, UserFk: userRes.ID})
		if authorizationTokenErr != nil {
			return authorizationTokenErr
		}
		authorizationTokenId := authorizationTokenRes.ID.String()
		refreshTokenId := refreshTokenRes.ID.String()
		jwtRefreshTokenRes, jwtRefreshTokenErr = v.dao.NewTokenQuery().CreateJwtRefreshToken(&refreshTokenId)
		if jwtRefreshTokenErr != nil {
			return jwtRefreshTokenErr
		}
		jwtAuthorizationTokenRes, jwtAuthorizationTokenErr = v.dao.NewTokenQuery().CreateJwtAuthorizationToken(&authorizationTokenId)
		if jwtAuthorizationTokenErr != nil {
			return jwtAuthorizationTokenErr
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &dto.RefreshToken{RefreshToken: *jwtRefreshTokenRes, AuthorizationToken: *jwtAuthorizationTokenRes}, nil
}
