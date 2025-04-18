package stop

import (
	"time"

	"github.com/brnocorreia/api-meu-buzufba/internal/infra/database/model"
	"github.com/brnocorreia/api-meu-buzufba/pkg/fault"
	"github.com/brnocorreia/api-meu-buzufba/pkg/slug"
	"github.com/brnocorreia/api-meu-buzufba/pkg/uid"
)

type stop struct {
	id              string
	slug            string
	name            string
	latitude        float64
	longitude       float64
	security_rating int
	is_active       bool
	created_at      time.Time
	updated_at      time.Time
}

func NewFromModel(m model.Stop) *stop {
	return &stop{
		id:              m.ID,
		slug:            m.Slug,
		name:            m.Name,
		latitude:        m.Latitude,
		longitude:       m.Longitude,
		security_rating: m.SecurityRating,
		is_active:       m.IsActive,
		created_at:      m.CreatedAt,
		updated_at:      m.UpdatedAt,
	}
}

func New(name string, latitude float64, longitude float64, security_rating int) (*stop, error) {

	s := stop{
		id:              uid.New("stop"),
		slug:            slug.New(name),
		name:            name,
		latitude:        latitude,
		longitude:       longitude,
		security_rating: security_rating,
		is_active:       true,
		created_at:      time.Now(),
		updated_at:      time.Now(),
	}

	if err := s.validate(); err != nil {
		return nil, fault.New("failed to create stop entity",
			fault.WithTag(fault.INVALID_ENTITY),
			fault.WithError(err),
		)
	}

	return &s, nil
}

func (s *stop) Model() model.Stop {
	return model.Stop{
		ID:             s.id,
		Slug:           s.slug,
		Name:           s.name,
		Latitude:       s.latitude,
		Longitude:      s.longitude,
		SecurityRating: s.security_rating,
		IsActive:       s.is_active,
		CreatedAt:      s.created_at,
		UpdatedAt:      s.updated_at,
	}
}

func (s *stop) Inactivate() {
	s.is_active = false
}

func (s *stop) validate() error {
	if s.name == "" {
		return fault.New("name is required")
	}

	if s.latitude > 90 || s.latitude < -90 {
		return fault.New("latitude must be between -90 and 90")
	}

	if s.longitude > 180 || s.longitude < -180 {
		return fault.New("longitude must be between -180 and 180")
	}

	if s.security_rating < 0 || s.security_rating > 5 {
		return fault.New("security rating must be between 0 and 5")
	}

	return nil
}
