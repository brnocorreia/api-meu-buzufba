package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/brnocorreia/api-meu-buzufba/internal/common/dto"
	"github.com/brnocorreia/api-meu-buzufba/internal/infra/http/middleware"
	"github.com/brnocorreia/api-meu-buzufba/internal/infra/mail"
	"github.com/brnocorreia/api-meu-buzufba/internal/modules/session"
	"github.com/brnocorreia/api-meu-buzufba/internal/modules/user"
	"github.com/brnocorreia/api-meu-buzufba/pkg/cache"
	"github.com/brnocorreia/api-meu-buzufba/pkg/crypto"
	"github.com/brnocorreia/api-meu-buzufba/pkg/fault"
	"github.com/brnocorreia/api-meu-buzufba/pkg/logging"
	"github.com/brnocorreia/api-meu-buzufba/pkg/token"
)

const (
	accessTokenDuration  = time.Minute * 15    // 15 minutes
	refreshTokenDuration = time.Hour * 24 * 30 // 30 days
)

// TODO: Remove userService dependency and user only the userRepo
// TODO: Remove sessionService dependency and user only the sessionRepo
type service struct {
	log            logging.Logger
	userService    user.Service
	userRepo       user.Repository
	sessionService session.Service
	sessionRepo    session.Repository
	mailer         *mail.Mail
	cache          *cache.Cache
	secretKey      string
}

func NewService(
	log logging.Logger,
	userService user.Service,
	userRepo user.Repository,
	sessionService session.Service,
	sessionRepo session.Repository,
	mailer *mail.Mail,
	cache *cache.Cache,
	secretKey string,
) Service {
	return &service{
		log:            log,
		userService:    userService,
		userRepo:       userRepo,
		sessionService: sessionService,
		sessionRepo:    sessionRepo,
		mailer:         mailer,
		cache:          cache,
		secretKey:      secretKey,
	}
}

func (s service) Logout(ctx context.Context) error {
	c, ok := ctx.Value(middleware.AuthKey{}).(*token.Claims)
	if !ok {
		s.log.Error(ctx, "context does not contain auth key")
		return fault.NewUnauthorized("access token no provided")
	}

	sessRecord, err := s.sessionRepo.GetActiveByUserID(ctx, c.UserID)
	if err != nil {
		return fault.NewBadRequest("failed to retrieve active session")
	} else if sessRecord == nil {
		return fault.NewNotFound("active session not found")
	}

	session := session.NewFromModel(*sessRecord)
	session.Deactivate()

	err = s.sessionRepo.Update(ctx, session.Model())
	if err != nil {
		return fault.NewBadRequest("failed to deactivate session")
	}

	// Delete the session from the cache when the user logs out
	go func() {
		ctx := context.Background()

		cacheKey := fmt.Sprintf("sess:%s", c.UserID)
		has, err := s.cache.Has(ctx, cacheKey)
		if err != nil {
			s.log.Errorw(ctx, "failed to check if session is in cache", logging.Err(err))
		} else if has {
			err = s.cache.Delete(ctx, c.UserID)
			if err != nil {
				s.log.Errorw(ctx, "failed to delete session from cache", logging.Err(err))
			}
		}
	}()

	return nil
}

func (s service) Activate(ctx context.Context, userId string) error {
	userRecord, err := s.userRepo.GetByID(ctx, userId)
	if err != nil {
		return fault.NewBadRequest("failed to get user by id")
	} else if userRecord == nil {
		return fault.NewNotFound("user not found")
	}

	if userRecord.Activated {
		return fault.New(
			"expired activation link",
			fault.WithHTTPCode(http.StatusBadRequest),
			fault.WithTag(fault.EXPIRED),
		)
	}

	user := user.NewFromModel(*userRecord)
	user.Activate()

	err = s.userRepo.Update(ctx, user.Model())
	if err != nil {
		s.log.Errorw(ctx, "failed to update user", logging.Err(err))
		return fault.NewBadRequest("failed to update user")
	}

	return nil
}

func (s service) GetSignedUser(ctx context.Context) (*dto.UserResponse, error) {
	c, ok := ctx.Value(middleware.AuthKey{}).(*token.Claims)
	if !ok {
		s.log.Error(ctx, "context does not contain auth key")
		return nil, fault.NewUnauthorized("access token no provided")
	}

	userRecord, err := s.userRepo.GetByID(ctx, c.UserID)
	if err != nil {
		return nil, fault.NewBadRequest("failed to retrieve user")
	} else if userRecord == nil {
		return nil, fault.NewNotFound("user not found")
	}

	user := &dto.UserResponse{
		ID:          userRecord.ID,
		Name:        userRecord.Name,
		Username:    userRecord.Username,
		Email:       userRecord.Email,
		IsUfba:      userRecord.IsUfba,
		Activated:   userRecord.Activated,
		ActivatedAt: userRecord.ActivatedAt,
		CreatedAt:   userRecord.CreatedAt,
		UpdatedAt:   userRecord.UpdatedAt,
	}

	return user, nil
}

func (s service) Register(ctx context.Context, input dto.CreateUser) error {
	_, err := s.userService.CreateUser(ctx, input)
	if err != nil {
		s.log.Errorw(ctx, "failed to create user", logging.Err(err))
		return err // The error is already being handled in the user service
	}

	// TODO: Send a welcome email here in the future

	return nil
}

func (s service) Login(ctx context.Context, email, password, ip, agent string) (*dto.LoginResponse, error) {
	userRecord, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fault.NewBadRequest("failed to get user by id")
	} else if userRecord == nil {
		return nil, fault.NewNotFound("user not found")
	}
	userID := userRecord.ID

	if !crypto.PasswordMatches(password, userRecord.Password) {
		return nil, fault.NewUnauthorized("invalid credentials")
	}

	err = s.sessionRepo.DeactivateAll(ctx, userID)
	if err != nil {
		return nil, fault.NewBadRequest("failed to deactivate user sessions")
	}

	accessToken, _, err := token.Gen(s.secretKey, userID, accessTokenDuration)
	if err != nil {
		s.log.Errorw(ctx, "failed to generate access token", logging.Err(err))
		return nil, fault.NewUnauthorized(err.Error())
	}

	refreshToken, _, err := token.Gen(s.secretKey, userID, refreshTokenDuration)
	if err != nil {
		s.log.Errorw(ctx, "failed to generate refresh token", logging.Err(err))
		return nil, fault.NewUnauthorized(err.Error())
	}

	params := dto.CreateSession{

		IP:           ip,
		Agent:        agent,
		UserID:       userID,
		RefreshToken: refreshToken,
	}
	sess, err := s.sessionService.CreateSession(ctx, params)
	if err != nil {
		return nil, err // The error is already being handled in the user service
	}

	response := dto.LoginResponse{
		SessionID:    sess.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return &response, nil
}
