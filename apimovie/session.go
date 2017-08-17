package apimovie

import (
	"time"
	"user-apiv2/apiuser"
	"user-apiv2/data"
)

type Key int

const MyKey Key = 0

/*Session has the database handle the services can reference them. By making the
MovieService a non-pointer field we reduce the allocations required when creating
a new session.*/
type Session struct {
	db           interface{}
	Movieservice MovieService
	now          time.Time
	*apiuser.Session
}

func NewSession(uDB data.MDBType, us *apiuser.Session) *Session {
	s := &Session{
		db:      uDB,
		Session: us,
	}
	s.Movieservice.session = s
	return s
}

// MovieService returns a movie service associated with this session.
func (s *Session) MovieService() data.MovieService {
	return &s.Movieservice
}
