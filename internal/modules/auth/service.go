package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/brnocorreia/api-meu-buzufba/internal/common/dto"
	"github.com/brnocorreia/api-meu-buzufba/internal/infra/http/middleware"
	"github.com/brnocorreia/api-meu-buzufba/internal/infra/http/token"
	"github.com/brnocorreia/api-meu-buzufba/internal/infra/mail"
	"github.com/brnocorreia/api-meu-buzufba/internal/modules/session"
	"github.com/brnocorreia/api-meu-buzufba/internal/modules/user"
	"github.com/brnocorreia/api-meu-buzufba/pkg/cache"
	"github.com/brnocorreia/api-meu-buzufba/pkg/crypto"
	"github.com/brnocorreia/api-meu-buzufba/pkg/dbutil"
	"github.com/brnocorreia/api-meu-buzufba/pkg/fault"
	"github.com/brnocorreia/api-meu-buzufba/pkg/logging"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

const (
	accessTokenDuration  = time.Minute * 15    // 15 minutes
	refreshTokenDuration = time.Hour * 24 * 30 // 30 days
	authServiceJourney   = "auth service"
)

type ServiceConfig struct {
	SecretKey      string
	UserService    user.Service
	UserRepo       user.Repository
	SessionService session.Service
	SessionRepo    session.Repository
	Mailer         *mail.Mail
	Cache          *cache.Cache
}

type service struct {
	userRepo       user.Repository
	sessionService session.Service
	sessionRepo    session.Repository
	mailer         *mail.Mail
	cache          *cache.Cache
	secretKey      string
}

func NewService(c ServiceConfig) Service {
	return &service{
		userRepo:       c.UserRepo,
		sessionService: c.SessionService,
		sessionRepo:    c.SessionRepo,
		mailer:         c.Mailer,
		cache:          c.Cache,
		secretKey:      c.SecretKey,
	}
}

func (s service) Logout(ctx context.Context) error {
	c, ok := ctx.Value(middleware.AuthKey{}).(*token.Claims)
	if !ok {
		logging.Error("context does not contain auth key",
			fmt.Errorf("access token not provided"),
			zap.String("journey", authServiceJourney))
		return fault.NewUnauthorized("access token not provided")
	}

	sessRecord, err := s.sessionRepo.GetActiveByUserID(ctx, c.UserID)
	if err != nil {
		logging.Error("failed to retrieve active session", err,
			zap.String("journey", authServiceJourney))
		return fault.NewBadRequest("failed to retrieve active session")
	} else if sessRecord == nil {
		logging.Error("active session not found", err,
			zap.String("journey", authServiceJourney))
		return fault.NewNotFound("active session not found")
	}

	sess := session.NewFromModel(*sessRecord)
	sess.Deactivate()

	err = s.sessionRepo.Update(ctx, sess.Model())
	if err != nil {
		logging.Error("failed to deactivate session", err,
			zap.String("journey", authServiceJourney))
		return fault.NewBadRequest("failed to deactivate session")
	}

	err = s.cache.Delete(ctx, fmt.Sprintf("sess:%s", c.UserID))
	if err != nil {
		logging.Error("failed to delete session from cache", err,
			zap.String("journey", authServiceJourney))
	}

	// Delete the session from the cache when the user logs out
	go func() {
		ctx := context.Background()

		cacheKey := fmt.Sprintf("sess:%s", c.UserID)
		has, err := s.cache.Has(ctx, cacheKey)
		if err != nil {
			logging.Error("failed to check if session is in cache", err,
				zap.String("journey", authServiceJourney))
		} else if has {
			err = s.cache.Delete(ctx, c.UserID)
			if err != nil {
				logging.Error("failed to delete session from cache", err,
					zap.String("journey", authServiceJourney))
			}
		}
	}()

	return nil
}

func (s service) Activate(ctx context.Context, userId string) error {
	userRecord, err := s.userRepo.GetByID(ctx, userId)
	if err != nil {
		logging.Error("failed to retrieve user", err,
			zap.String("journey", authServiceJourney))
		return fault.NewBadRequest("failed to get user by id")
	} else if userRecord == nil {
		logging.Error("user not found", err,
			zap.String("journey", authServiceJourney))
		return fault.NewNotFound("user not found")
	}

	if userRecord.Activated {
		logging.Error("user already activated", err,
			zap.String("journey", authServiceJourney))
		return fault.New(
			"expired activation link",
			fault.WithHTTPCode(http.StatusBadRequest),
			fault.WithTag(fault.EXPIRED),
		)
	}

	u := user.NewFromModel(*userRecord)
	u.Activate()

	err = s.userRepo.Update(ctx, u.Model())
	if err != nil {
		logging.Error("failed to update user", err,
			zap.String("journey", authServiceJourney))
		return fault.NewBadRequest("failed to update user")
	}

	return nil
}

func (s service) GetSignedUser(ctx context.Context) (*dto.UserResponse, error) {
	c, ok := ctx.Value(middleware.AuthKey{}).(*token.Claims)
	if !ok {
		logging.Error("context does not contain auth key",
			fmt.Errorf("access token not provided"),
			zap.String("journey", authServiceJourney))
		return nil, fault.NewUnauthorized("access token not provided")
	}

	userRecord, err := s.userRepo.GetByID(ctx, c.UserID)
	if err != nil {
		logging.Error("failed to retrieve user", err,
			zap.String("journey", authServiceJourney))
		return nil, fault.NewBadRequest("failed to retrieve user")
	} else if userRecord == nil {
		logging.Error("user not found", err,
			zap.String("journey", authServiceJourney))
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
	userRecord, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		logging.Error("failed to retrieve user", err,
			zap.String("journey", authServiceJourney))
		return fault.NewBadRequest("failed to get user by email")
	} else if userRecord != nil {
		logging.Error("failed to create user: e-mail already taken", err,
			zap.String("journey", authServiceJourney))
		return fault.NewConflict("e-mail already taken")
	}

	isUfba := checkIfUserEmailIsUfba(input.Email)

	newUser, err := user.New(input.Name, input.Username, input.Email, input.Password, isUfba)
	if err != nil {
		logging.Error("failed to create user", err,
			zap.String("journey", authServiceJourney))
		return fault.NewUnprocessableEntity("failed to create user entity")
	}
	model := newUser.Model()

	if err = s.userRepo.Insert(ctx, model); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" { // 23505 is the code for unique constraint violation
			field := dbutil.ExtractFieldFromDetail(pqErr.Detail)
			return fault.NewConflict(fmt.Sprintf("%s already taken", field))
		}
		logging.Error("failed to insert user", err,
			zap.String("journey", authServiceJourney))
		return fault.NewBadRequest("failed to insert user")
	}

	// TODO: Send a welcome email here in the future

	return nil
}

func (s service) Login(ctx context.Context, email, password, ip, agent string) (*dto.LoginResponse, error) {
	userRecord, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		logging.Error("failed to retrieve user", err,
			zap.String("journey", authServiceJourney))
		return nil, fault.NewBadRequest("failed to get user by email")
	} else if userRecord == nil {
		logging.Error("user not found", err,
			zap.String("journey", authServiceJourney))
		return nil, fault.NewNotFound("user not found")
	}
	userID := userRecord.ID

	if !crypto.PasswordMatches(password, userRecord.Password) {
		return nil, fault.NewUnauthorized("invalid credentials")
	}

	err = s.sessionRepo.DeactivateAll(ctx, userID)
	if err != nil {
		logging.Error("failed to deactivate user sessions", err,
			zap.String("journey", authServiceJourney))
		return nil, fault.NewBadRequest("failed to deactivate user sessions")
	}

	accessToken, _, err := token.Gen(s.secretKey, userID, accessTokenDuration)
	if err != nil {
		logging.Error("failed to generate access token", err,
			zap.String("journey", authServiceJourney))
		return nil, fault.NewUnauthorized(err.Error())
	}

	refreshToken, _, err := token.Gen(s.secretKey, userID, refreshTokenDuration)
	if err != nil {
		logging.Error("failed to generate refresh token", err,
			zap.String("journey", authServiceJourney))
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

func checkIfUserEmailIsUfba(email string) bool {
	return strings.HasSuffix(email, "@ufba.br")
}
