package apigame

type GameID int64

type GDBType map[GameID]*Game

type Game struct {
	ID             GameID
	correctWord    string
	guessedLetter  byte
	GuessedLetterR string
	wordSoFar      []byte
	WordSoFarR     string
	wrongGuesses   string
	WrongGuessesR  []byte
	PlayerID       string
	Invalid        string
	PlayerName     string
	Name           string
	Category       string
	Count          int
	Winorloss      bool
	Word           string
	ImageURL       string
	JustStarted    string
}

type Request struct {
	Letter string
	ID     GameID
}

type Response struct {
	Hash string
}

// type GameService interface {
// 	Start(context.Context) (GameID, error)
// 	//Guess(context.Context, *Game) (*Game, error)
// 	WinOrLoss(*Game) bool
// 	AddGame(context.Context, *Game) (GameID, error)
// 	GetGame(GameID) (*Game, error)
// 	DeleteGame(context.Context, GameID) error
// 	ListGames() ([]*Game, error)
// }

// type Session interface {
// 	GameService() GameService
// 	//UserService() UserService
// }

// SetCreatorAnonymous sets the PlayerID field to the "anonymous" ID.
// func (b *Game) SetCreatorAnonymous() {
// 	b.PlayerName = ""
// 	b.PlayerID = "anonymous"
// }

// PlayerNameDisplayName returns a string appropriate for displaying the name of
// the user who created this gamie object.
// func (b *Game) PlayerDisplayName() string {
// 	if b.PlayerID == "anonymous" {
// 		return "Anonymous"
// 	}
// 	return b.PlayerName
// }

// type Counter interface {
// 	Count(m int) bool
// }
