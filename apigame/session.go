package apigame

import (
	"time"
	"user-apiv2/apiuser"
)

type Key int

const MyKey Key = 0

/*Session has the database handle the services can reference them. By making the
GameService a non-pointer field we reduce the allocations required when creating
a new session.*/
type Session struct {
	db             interface{}
	Gameservice    GameService
	GameGuiService *GameGuiService
	//Pageservice    apipage.PageService
	//PageGuiService *apipage.PageGuiService
	now time.Time
	*apiuser.Session
}

func NewSession(uDB GDBType, us *apiuser.Session, gg *GameGuiService) *Session {

	s := &Session{
		db:      uDB,
		Session: us,
		//PageGuiService: gs,
		GameGuiService: gg,
	}
	s.Gameservice.session = s
	//s.Pageservice.Session = s
	return s
}

// // GameService returns a game service associated with this session.
// func (s *Session) GameService() GameService {
// 	return &s.Gameservice
// }
