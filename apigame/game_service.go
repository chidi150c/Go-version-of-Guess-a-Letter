package apigame

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"user-apiv2/data"
)

var nextID GameID

func init() {
	nextID = 1
}

const AddGameKey Key = 1

type GameService struct {
	session *Session
}

func (g *GameService) ServerGuessALetter(letter string, gameID GameID) (*Game, error) {

	gm, err := g.GetGame(gameID)
	if err != nil {
		return nil, err
	}
	gm.GuessedLetterR = letter
	gl := strings.ToUpper(gm.GuessedLetterR)
	//cw := strings.ToUpper(gm.correctWord)
	//result := gm
	fmt.Printf("\n\n\n********************* In ServerGuessALetter0 *************************  %vn\n\n\n", gm)
	if gm.GuessedLetterR == "" {
		fmt.Printf("\n\n\n********************* In ServerGuessALetter1 *************************  %vn\n\n\n", gm)
		gm.JustStarted = "false"
		if err := g.UpdateGame(gm); err != nil {
			return nil, err
		}
		return gm, nil
	}
	if !strings.Contains("ABCDEFGHIJKLMNOPQRSTUVWXYZ", gl) {
		gm.Invalid = "true"
		gm.JustStarted = "false"
		if err := g.UpdateGame(gm); err != nil {
			return gm, err
		}
		return gm, nil
	} else {
		gm.Invalid = "false"
	}
	fmt.Printf("\n\n\n********************* In ServerGuessALetter2 *************************  %vn\n\n\n", gm)

	gm.guessedLetter = gm.GuessedLetterR[0]
	if !strings.Contains(gm.correctWord, gl) {
		fmt.Printf("\n\n\n********************* In ServerGuessALetter3 *************************  %vn\n\n\n", gm)
		gm.WrongGuessesR = append(gm.WrongGuessesR, gm.guessedLetter)
	} else {
		for {
			fmt.Printf("\n\n\n********************* In ServerGuessALetter4 *************************  %vn\n\n\n", gm)
			index := strings.Index(gm.wrongGuesses, gl)
			if index == -1 {
				fmt.Printf("\n\n\n********************* In ServerGuessALetter5 *************************  %vn\n\n\n", gm)
				break
			}
			gm.wordSoFar[index] = gm.guessedLetter
			gm.wrongGuesses = strings.Replace(gm.wrongGuesses, gl, "-", 1)
			gm.WordSoFarR = string(gm.wordSoFar)
			fmt.Printf("\n\n\n********************* In ServerGuessALetter6 *************************  %vn\n\n\n", gm)
		}
	}

	fmt.Printf("\n\n\n********************* In ServerGuessALetter7 *************************  %vn\n\n\n", gm)
	gm.JustStarted = "false"
	if err := g.UpdateGame(gm); err != nil {
		return nil, err
	}
	if gm.Count == 6 {
		gm.Winorloss = g.WinOrLoss(gm)
	}
	gm.Count++
	return gm, nil
}

func (g *GameService) ServeRender(letter string, gameID GameID) (string, error) {
	fmt.Println("ServeRender **********", letter, gameID)

	//return data.Response{Hash: "#"}, nil
	return "#", nil
}

func (g *GameService) Start(ctx context.Context) (string, GameID, error) {
	// Retrieve current session user.
	u, ok := ctx.Value(AddGameKey).(*data.User)
	if !ok {
		return "", 0, data.ErrUnauthorized
	}
	var bB []byte
	resp, _ := http.PostForm("http://watchout4snakes.com/wo4snakes/Random/RandomWord", url.Values{})
	if resp != nil {
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			defer resp.Body.Close()
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return "", 0, err
			}
			bB = bodyBytes
		}
	} else {
		bB = []byte("correction")
	}
	gm := &Game{
		correctWord:    string(bB),
		PlayerID:       u.Username,
		JustStarted:    "true",
		GuessedLetterR: "",
		guessedLetter:  0,
		Count:          0,
	}
	length := len(gm.correctWord)
	gm.wordSoFar = make([]byte, length)
	gm.correctWord = strings.ToUpper(gm.correctWord)
	gm.wrongGuesses = gm.correctWord
	for i := 0; i < length; i++ {
		gm.wordSoFar[i] = '-'
	}
	gm.WordSoFarR = string(gm.wordSoFar)
	fmt.Println(gm)
	id, err := g.AddGame(ctx, gm)
	if err != nil {
		log.Fatal(err)
	}
	return gm.GuessedLetterR, id, nil
}

func (g *GameService) WinOrLoss(gm *Game) bool {
	gm.Word = gm.correctWord
	a := strings.ToLower(string(gm.wordSoFar))
	b := strings.ToLower(string(gm.correctWord))
	return reflect.DeepEqual(a, b)
}

func (g *GameService) AddGame(ctx context.Context, gam *Game) (GameID, error) {
	// Retrieve current session user.
	u, ok := ctx.Value(AddGameKey).(*data.User)
	if !ok {
		return 0, data.ErrUnauthorized
	}
	db, ok := g.session.db.(GDBType)
	if !ok {
		return 0, data.ErrGamDbUnreachable
	}
	if gam.PlayerID != u.Username {
		return gam.ID, data.ErrUnauthorized
	}
	db[nextID] = gam
	gam.ID = nextID
	nextID++
	return gam.ID, nil
}

func (g *GameService) GetGame(gamid GameID) (*Game, error) {
	db, ok := g.session.db.(GDBType)
	if !ok {
		return nil, data.ErrGamDbUnreachable
	}
	if gamid != 0 {
		gam := db[gamid]
		if gam == nil {
			return nil, data.ErrGameNotFound
		}
		return gam, nil
	}
	return nil, data.ErrGameIDRequired
}

func (g *GameService) DeleteGame(ctx context.Context, gamid GameID) error {
	u, ok := ctx.Value(AddGameKey).(data.User)
	if !ok {
		return data.ErrUnauthorized
	}
	db, ok := g.session.db.(GDBType)
	if !ok {
		return data.ErrGamDbUnreachable
	}
	if gamid != 0 {
		gam := db[gamid]
		if gam == nil {
			return data.ErrGameNotFound
		}
		if gam.ID != gamid && gam.PlayerID == u.Username {
			return data.ErrUnauthorized
		}
		delete(db, gamid)
	}
	return nil
}

func (g *GameService) ListGames() ([]*Game, error) {
	mdb, ok := g.session.db.(GDBType)
	if !ok {
		return nil, data.ErrGamDbUnreachable
	}
	var db []*Game
	for _, b := range mdb {
		db = append(db, b)
	}
	return db, nil
}

func (m *GameService) UpdateGame(gam *Game) error {
	// u, ok := ctx.Value(AddGameKey).(data.User)
	// if !ok {
	// 	return data.ErrUnauthorized
	// }
	db, ok := m.session.db.(GDBType)
	if !ok {
		return data.ErrGamDbUnreachable
	}
	// Only allow player to update Product.
	gameInDB, ok := db[gam.ID]
	if !ok {
		return fmt.Errorf("memory db: product not found with ID %v", gam.ID)
	} else if gameInDB.ID != gam.ID { /* && gam.PlayerID == u.Username*/
		return data.ErrUnauthorized
		//return fmt.Errorf("memory db: Non player not allowed to update Product %v", b.ID)
	}
	if gam.ID == 0 {
		return errors.New("memory db: product with unassigned ID passed into updateProduct")
	}
	db[gam.ID] = gam
	return nil
}
