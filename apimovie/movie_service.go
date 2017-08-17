package apimovie

import (
	"context"
	"fmt"
	"user-apiv2/data"

	"github.com/pkg/errors"
)

var nextID data.MovieID

func init() {
	nextID = 1
}

const AddMovieKey Key = 1

type MovieService struct {
	session *Session
}

var _ data.MovieService = &MovieService{}

func (m *MovieService) AddMovie(ctx context.Context, mov *data.Movie) (data.MovieID, error) {
	// Retrieve current session user.
	u, ok := ctx.Value(AddMovieKey).(*data.User)
	if !ok {
		return 0, data.ErrUnauthorized
	}
	db, ok := m.session.db.(data.MDBType)
	if !ok {
		return 0, data.ErrMovDbUnreachable
	}
	mov.OwnerID = u.Username
	mov.ModTime = m.session.now
	db[nextID] = mov
	mov.ID = nextID
	nextID++
	return mov.ID, nil
}

func (m *MovieService) GetMovie(movid data.MovieID) (*data.Movie, error) {
	db, ok := m.session.db.(data.MDBType)
	if !ok {
		return nil, data.ErrMovDbUnreachable
	}
	if movid != 0 {
		mov := db[movid]
		if mov == nil {
			return nil, data.ErrMovieNotFound
		}
		return mov, nil
	}
	return nil, data.ErrMovieIDRequired
}

func (m *MovieService) DeleteMovie(ctx context.Context, movid data.MovieID) error {
	u, ok := ctx.Value(AddMovieKey).(data.User)
	if !ok {
		return data.ErrUnauthorized
	}
	db, ok := m.session.db.(data.MDBType)
	if !ok {
		return data.ErrMovDbUnreachable
	}
	if movid != 0 {
		mov := db[movid]
		if mov == nil {
			return data.ErrMovieNotFound
		}
		if mov.ID != movid && mov.OwnerID == u.Username {
			return data.ErrUnauthorized
		}
		delete(db, movid)
	}
	return nil
}

func (m *MovieService) ListMovies() ([]*data.Movie, error) {
	mdb, ok := m.session.db.(data.MDBType)
	if !ok {
		return nil, data.ErrMovDbUnreachable
	}
	var db []*data.Movie
	for _, b := range mdb {
		db = append(db, b)
	}
	return db, nil
}

func (m *MovieService) UpdateMovie(ctx context.Context, mov *data.Movie) error {
	u, ok := ctx.Value(AddMovieKey).(data.User)
	if !ok {
		return data.ErrUnauthorized
	}
	db, ok := m.session.db.(data.MDBType)
	if !ok {
		return data.ErrMovDbUnreachable
	}
	// Only allow owner to update Product.
	movieInDB, ok := db[mov.ID]
	if !ok {
		return fmt.Errorf("memory db: product not found with ID %v", mov.ID)
	} else if movieInDB.ID != mov.ID && mov.OwnerID == u.Username {
		return data.ErrUnauthorized
		//return fmt.Errorf("memory db: Non owner not allowed to update Product %v", b.ID)
	}
	if mov.ID == 0 {
		return errors.New("memory db: product with unassigned ID passed into updateProduct")
	}
	db[mov.ID] = mov
	return nil
}
