package title

import (
	"fmt"
	"github.com/michaelvmata/path/session"
	"github.com/michaelvmata/path/world"
)

func ListCharacters(s *session.Session, players map[string]*world.Character) {
	for rank, player := range players {
		s.Outgoing <- fmt.Sprintf("[%s] %s", rank, player.Name)
	}
}
