package session

import (
	"context"
	"fmt"
	"time"

	"github.com/brnocorreia/api-meu-buzufba/internal/common/dto"
	"github.com/brnocorreia/api-meu-buzufba/internal/infra/database/model"
	"github.com/brnocorreia/api-meu-buzufba/internal/infra/http/middleware"
	"github.com/brnocorreia/api-meu-buzufba/internal/infra/http/token"
	"github.com/brnocorreia/api-meu-buzufba/internal/modules/user"
	"github.com/brnocorreia/api-meu-buzufba/pkg/cache"
	"github.com/brnocorreia/api-meu-buzufba/pkg/fault"
	"github.com/brnocorreia/api-meu-buzufba/pkg/logging"
	"go.uber.org/zap"

	"github.com/medama-io/go-useragent"
)

const sessionServiceJourney = "session service"

type ServiceConfig struct {
	SessionRepo Repository
	UserRepo    user.Repository

	Cache *cache.Cache

	SecretKey string
}

type service struct {
	sessionRepo Repository
	userRepo    user.Repository
	cache       *cache.Cache

	secretKey string
}

func NewService(c ServiceConfig) Service {
	return &service{
		sessionRepo: c.SessionRepo,
		userRepo:    c.UserRepo,
		cache:       c.Cache,
		secretKey:   c.SecretKey,
	}
}

func (s service) RenewAccessToken(ctx context.Context, refreshToken string) (*dto.RenewAccessToken, error) {
	claims, err := token.Verify(s.secretKey, refreshToken)
	if err != nil {
		logging.Error("invalid refresh token", err,
			zap.String("journey", sessionServiceJourney))
		return nil, fault.NewUnauthorized("invalid refresh token")
	}

	sessRecord, err := s.sessionRepo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		logging.Error("failed to retrieve session", err,
			zap.String("journey", sessionServiceJourney))
		return nil, fault.NewBadRequest("failed to retrieve session")
	}
	session := NewFromModel(*sessRecord)

	if session.IsExpired() {
		logging.Error("session has expired", err,
			zap.String("journey", sessionServiceJourney))
		return nil, fault.NewBadRequest("session has expired")
	}

	if session.UserID() != claims.UserID {
		logging.Info("unauthorized user",
			zap.String("journey", sessionServiceJourney),
			zap.String("sessionUserID", session.UserID()),
			zap.String("claimsUserID", claims.UserID))
		return nil, fault.NewUnauthorized("unauthorized user")
	}

	newAccessToken, _, err := token.Gen(s.secretKey, claims.UserID, time.Minute*15)
	if err != nil {
		logging.Error("failed to generate access token", err,
			zap.String("journey", sessionServiceJourney))
		return nil, fault.NewBadRequest("failed to generate access token")
	}

	return &dto.RenewAccessToken{
		AccessToken:        newAccessToken,
		AccessTokenExpires: time.Now().Add(time.Minute * 15),
	}, nil
}

func (s service) GetSessionByUserID(ctx context.Context, userID string) (*dto.SessionResponse, error) {
	var cachedSession *model.Session
	err := s.cache.GetStruct(ctx, fmt.Sprintf("sess:%s", userID), &cachedSession)
	if err != nil {
		switch {
		case fault.GetTag(err) == fault.CACHE_MISS:
			logging.Info("cache: miss session not found",
				zap.String("journey", sessionServiceJourney))
		default:
			logging.Error("failed to query session from cache", err,
				zap.String("journey", sessionServiceJourney))
		}
	}

	if cachedSession != nil {
		if time.Now().After(cachedSession.Expires) {
			logging.Info("session has expired",
				zap.String("journey", sessionServiceJourney),
				zap.String("userID", userID))
			return nil, fault.NewBadRequest("session has expired")
		}

		return &dto.SessionResponse{
			ID:        cachedSession.ID,
			Agent:     cachedSession.Agent,
			IP:        cachedSession.IP,
			Active:    cachedSession.Active,
			CreatedAt: cachedSession.CreatedAt,
			UpdatedAt: cachedSession.UpdatedAt,
		}, nil
	}

	sessionRecord, err := s.sessionRepo.GetActiveByUserID(ctx, userID)
	if err != nil {
		logging.Error("failed to retrieve active session", err,
			zap.String("journey", sessionServiceJourney))
		return nil, fault.NewBadRequest("failed to retrieve session")
	} else if sessionRecord == nil {
		logging.Info("active session not found",
			zap.String("journey", sessionServiceJourney),
			zap.String("userID", userID))
		return nil, fault.NewNotFound("session not found")
	}

	res := &dto.SessionResponse{
		ID:        sessionRecord.ID,
		Agent:     sessionRecord.Agent,
		IP:        sessionRecord.IP,
		Active:    sessionRecord.Active,
		CreatedAt: sessionRecord.CreatedAt,
		UpdatedAt: sessionRecord.UpdatedAt,
	}

	cacheKey := fmt.Sprintf("sess:%s", userID)
	err = s.cache.SetStruct(ctx, cacheKey, sessionRecord, time.Minute*30)
	if err != nil {
		logging.Error("failed to cache session", err,
			zap.String("journey", sessionServiceJourney))
	}

	return res, nil
}

func (s service) GetAllSessions(ctx context.Context) ([]dto.SessionResponse, error) {
	c, ok := ctx.Value(middleware.AuthKey{}).(*token.Claims)
	if !ok {
		logging.Error("context does not contain auth key",
			fmt.Errorf("access token no provided"),
			zap.String("journey", sessionServiceJourney))
		return nil, fault.NewUnauthorized("access token no provided")
	}

	records, err := s.sessionRepo.GetAllByUserID(ctx, c.UserID)
	if err != nil {
		logging.Error("failed to retrieve sessions", err,
			zap.String("journey", sessionServiceJourney))
		return nil, fault.NewBadRequest("failed to retrieve sessions")
	}

	if len(records) == 0 {
		return make([]dto.SessionResponse, 0), nil
	}

	// Pre-allocate the slice to avoid reallocations
	// This is more efficient than appending to the slice
	sessions := make([]dto.SessionResponse, len(records))
	for i, s := range records {
		sessions[i] = dto.SessionResponse{
			ID:        s.ID,
			Agent:     s.Agent,
			IP:        s.IP,
			Active:    s.Active,
			CreatedAt: s.CreatedAt,
			UpdatedAt: s.UpdatedAt,
		}
	}

	return sessions, nil
}

func (s service) CreateSession(ctx context.Context, input dto.CreateSession) (*dto.SessionResponse, error) {
	userRecord, err := s.userRepo.GetByID(ctx, input.UserID)
	if err != nil {
		logging.Error("failed to retrieve user", err,
			zap.String("journey", sessionServiceJourney))
		return nil, fault.NewBadRequest("failed to retrieve user")
	} else if userRecord == nil {
		logging.Info("user not found",
			zap.String("journey", sessionServiceJourney),
			zap.String("userID", input.UserID))
		return nil, fault.NewNotFound("user not found with ID: " + input.UserID)
	}
	userID := userRecord.ID

	rawAgent := useragent.NewParser().Parse(input.Agent)
	agent := fmt.Sprintf("%s em %s", rawAgent.Browser(), rawAgent.OS())
	// In case the user agent is not a browser, we will use "unknown agent"
	// Likely to happen in mobile devices or in CLI/Postman and similar tools
	// Output: "<browser> em <os>"
	if rawAgent.Browser() == "" || rawAgent.OS() == "" {
		agent = "unknown agent"
	}

	sess, err := New(userID, input.IP, agent, input.RefreshToken)
	if err != nil {
		logging.Error("failed to create session entity", err,
			zap.String("journey", sessionServiceJourney))
		return nil, fault.NewUnprocessableEntity("failed to create session entity")
	}

	err = s.sessionRepo.Insert(ctx, sess.Model())
	if err != nil {
		logging.Error("failed to insert session entity", err,
			zap.String("journey", sessionServiceJourney))
		return nil, fault.NewBadRequest("failed to insert session entity")
	}

	res := dto.SessionResponse{
		ID:        sess.ID(),
		Agent:     sess.Agent(),
		IP:        sess.IP(),
		Active:    sess.Active(),
		CreatedAt: sess.CreatedAt(),
		UpdatedAt: sess.UpdatedAt(),
	}

	return &res, nil
}
