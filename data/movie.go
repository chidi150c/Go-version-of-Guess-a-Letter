package data

import (
	"context"
	"fmt"
	"math"
	"time"
)

type MDBType map[MovieID]*Movie

type MovieID int64

//Movie holds info about the  Movies
type Movie struct {
	ID             MovieID `json:"id"`
	Name           string  `json:"name"`
	OwnerName      string  `json:"ownername"`
	OwnerID        string  `json:"ownerid"`
	OwnerCategory  string  `json:"ownercategory"`
	Token          string  `json:"token"`
	ProductionDate string
	Address        string    `json:"address"`
	Category       string    `json:"category"`
	ModTime        time.Time `json:"modTime"`
	Author         string
	Description    string
	ImageURL       string
	Phone          string
	EmailAddress   string
}

type MovieService interface {
	AddMovie(context.Context, *Movie) (MovieID, error)
	GetMovie(MovieID) (*Movie, error)
	DeleteMovie(context.Context, MovieID) error
	ListMovies() ([]*Movie, error)
	UpdateMovie(context.Context, *Movie) error
}

// SetCreatorAnonymous sets the OwnerNameID field to the "anonymous" ID.
func (b *Movie) SetCreatorAnonymous() {
	b.OwnerName = ""
	b.OwnerID = "anonymous"
}

// OwnerNameDisplayName returns a string appropriate for displaying the name of
// the user who created this movie object.
func (b *Movie) OwnerNameDisplayName() string {
	if b.OwnerID == "anonymous" {
		return "Anonymous"
	}
	return b.OwnerName
}

// type Counter interface {
// 	Count(m int) bool
// }

type Muu []*Movie

var ct int

func (b Muu) Count(m int) bool {
	fmt.Println("\ncout = ", ct)
	ct++
	if math.Mod(float64(m), 4) == 0 {
		return true
	}
	return false
}

func (b *Movie) MustName() string {
	if b.Name == "" {
		return "UnNamed"
	}
	return b.Name
}
