package stop

import (
	"net/http"
	"sync"

	"github.com/brnocorreia/api-meu-buzufba/pkg/fault"
	"github.com/brnocorreia/api-meu-buzufba/pkg/logging"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	httputil "github.com/brnocorreia/api-meu-buzufba/pkg/http_util"
)

const (
	stopHandlerJourney = "stop handler"
)

var (
	instance *handler
	Once     sync.Once
)

type handler struct {
	stopService Service
}

func NewHandler(stopService Service) *handler {
	Once.Do(func() {
		instance = &handler{stopService: stopService}
	})
	return instance
}

func (h handler) Register(r *chi.Mux) {

	r.Route("/api/v1/stops", func(r chi.Router) {
		r.Get("/{slug}", h.handleGetStopBySlug)
		r.Get("/{id}", h.handleGetStopByID)
	})
}

func (h handler) handleGetStopBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	stop, err := h.stopService.GetStopBySlug(r.Context(), slug)
	if err != nil {
		logging.Error("failed to get stop by slug", err, zap.String("journey", stopHandlerJourney))
		fault.NewHTTPError(w, err)
		return
	}
	httputil.WriteJSON(w, http.StatusOK, stop)
}

func (h handler) handleGetStopByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	stop, err := h.stopService.GetStopByID(r.Context(), id)
	if err != nil {
		logging.Error("failed to get stop by id", err, zap.String("journey", stopHandlerJourney))
		fault.NewHTTPError(w, err)
		return
	}
	httputil.WriteJSON(w, http.StatusOK, stop)
}
