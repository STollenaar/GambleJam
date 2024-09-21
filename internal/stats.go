package internal

import (
	"fmt"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/stollenaar/gamblingjam/util"
)

var (
	winnerTicketAsset  *ebiten.Image
	scratchTicketAsset *ebiten.Image
)

const (
	layoutTime = "15:04 PM"
)

type Item interface {
	Interact() interface{}
}

type Stats struct {
	drawWinnner bool

	money int
	day   int
	time  time.Time

	inventory []Item
	home      *Home
	store     *Store
}

func init() {
	igbm, _, err := ebitenutil.NewImageFromFile("./assets/winner.png")
	if err != nil {
		log.Fatal(err)
	}
	winnerTicketAsset = igbm

	sigbm, _, err := ebitenutil.NewImageFromFile("./assets/scratchTicket.png")
	if err != nil {
		log.Fatal(err)
	}
	scratchTicketAsset = sigbm
}

func (s *Stats) CheckTicket(slot int) (t *Ticket) {
	ticket := s.inventory[slot]
	if ticket == nil {
		return nil
	}
	s.advanceTime(time.Minute * 10)
	if game := ticket.Interact(); game != nil {
		s.money += game.(*TicketGame).Win
		// Play winning sounds/graphics
		fmt.Printf("You won: %s, $%d\n", game.(*TicketGame).Prize, game.(*TicketGame).Win)
		go util.PlayMP3("assets/wahoo.mp3")
		s.drawWinnner = true
		t = ticket.(*Ticket)
	} else {
		fmt.Println("Better luck next time")
	}
	s.inventory[slot] = nil
	return t
}

func (s *Stats) HandleButtons(place Place) bool {
	if place == STORE {
		ticket := s.store.FindTicket()
		if s.time.Hour() >= 21 {
			return false
		}

		if ticket != nil && s.money >= ticket.Cost && s.findEmptySlot(ticket) {
			s.money -= ticket.Cost
			s.advanceTime(time.Minute * 5)
			return true
		}
	}
	return false
}

func (s *Stats) findEmptySlot(input Item) bool {
	for i, item := range s.inventory {
		if item == nil {
			s.inventory[i] = input
			return true
		}
	}
	return false
}

func (s *Stats) advanceTime(d time.Duration) {
	if s.time.Add(d).Day() == s.time.Day() {
		s.time = s.time.Add(d)
	}
}
