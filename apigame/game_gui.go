package apigame

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"user-apiv2/data"
)

type GameGuiService struct {
	URL *url.URL
	//gameService GameService
}

func NewGameGuiService() *GameGuiService {
	a := &GameGuiService{URL: &url.URL{}}
	return a
}

func (g *GameGuiService) SaveAGame(gm *Game) {

}

type request struct {
	Letter string
	ID     GameID
}

func (g *GameGuiService) GuessALetter(letter string, gmID GameID) (*Game, error) {
	//Validate arguments.
	if gmID == 0 {
		return nil, data.ErrGameRequired
	}
	req := request{
		Letter: letter,
		ID:     gmID,
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(&req)
	resp, err := http.Post("http://localhost:8006/games/guessletter", "application/json; charset=utf-8", b)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Decode response into JSON.
	var respBody Game
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, err
	}

	// Copy returned dial.

	return &respBody, nil

}

func (g *GameGuiService) MustName(Name string) string {
	if Name == "" {
		return "Guess A Letter Game"
	}
	return Name
}

var gct int

func (g *GameGuiService) Count(m int) bool {
	fmt.Println("\ncout = ", gct)
	gct++
	if math.Mod(float64(m), 4) == 0 {
		return true
	}
	return false
}
