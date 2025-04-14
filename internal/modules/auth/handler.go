package auth

import (
	"net/http"
	"sync"

	"github.com/brnocorreia/api-meu-buzufba/internal/common/dto"
	"github.com/brnocorreia/api-meu-buzufba/internal/infra/http/middleware"
	"github.com/brnocorreia/api-meu-buzufba/pkg/fault"
	"github.com/brnocorreia/api-meu-buzufba/pkg/logging"
	"go.uber.org/zap"

	httputil "github.com/brnocorreia/api-meu-buzufba/pkg/http_util"
	"github.com/go-chi/chi/v5"
)

const (
	authHandlerJourney = "auth handler"
)

var (
	instance *handler
	Once     sync.Once
)

type handler struct {
	authService Service
	secretKey   string
}

func NewHandler(authService Service, secretKey string) *handler {
	Once.Do(func() {
		instance = &handler{
			authService: authService,
			secretKey:   secretKey,
		}
	})
	return instance
}

func (h handler) Register(r *chi.Mux) {
	m := middleware.NewWithAuth(h.secretKey)

	r.Route("/api/v1/auth", func(r chi.Router) {
		// Private
		r.With(m.WithAuth).Get("/me", h.handleGetSigned)
		r.With(m.WithAuth).Patch("/logout", h.handleLogout)
		// Public
		r.Get("/activate/{userId}", h.handleActivate)
		r.Post("/register", h.handleRegister)
		r.Post("/login", h.handleLogin)
	})
}

func (h handler) handleLogout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := h.authService.Logout(ctx)
	if err != nil {
		logging.Error("failed to logout", err, zap.String("journey", authHandlerJourney))
		fault.NewHTTPError(w, err)
		return
	}

	httputil.WriteSuccess(w, http.StatusOK)
}

func (h handler) handleActivate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := chi.URLParam(r, "userId")

	err := h.authService.Activate(ctx, userId)
	if err != nil {
		logging.Error("failed to activate user", err, zap.String("journey", authHandlerJourney))
		fault.NewHTTPError(w, err)
		return
	}

	httputil.WriteSuccess(w, http.StatusOK)
}

func (h handler) handleGetSigned(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := h.authService.GetSignedUser(ctx)
	if err != nil {
		logging.Error("failed to get signed user", err, zap.String("journey", authHandlerJourney))
		fault.NewHTTPError(w, err)
		return
	}

	httputil.WriteJSON(w, http.StatusOK, user)
}

func (h handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var body dto.CreateUser
	err := httputil.ReadRequestBody(w, r, &body)
	if err != nil {
		logErrorInReadRequestBody(err, r)
		fault.NewHTTPError(w, err)
		return
	}

	err = h.authService.Register(ctx, body)
	if err != nil {
		fault.NewHTTPError(w, err)
		return
	}

	httputil.WriteSuccess(w, http.StatusAccepted)
}

func (h handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := httputil.ReadRequestBody(w, r, &body)
	if err != nil {
		logErrorInReadRequestBody(err, r)
		fault.NewHTTPError(w, err)
		return
	}

	res, err := h.authService.Login(ctx, body.Email, body.Password, r.RemoteAddr, r.UserAgent())
	if err != nil {
		fault.NewHTTPError(w, err)
		return
	}

	httputil.WriteJSON(w, http.StatusOK, res)
}

func logErrorInReadRequestBody(err error, r *http.Request) {
	logging.Error("failed to read request body", err,
		zap.String("journey", authHandlerJourney),
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path))
}
