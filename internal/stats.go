package internal

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	winnerTicketAsset *ebiten.Image
)

type Item interface {
	Interact() interface{}
}

type Stats struct {
	drawWinnner bool

	money int
	day   int

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
}

func (s *Stats) CheckAllTickets() {
	var totalWinning int

	for index, ticket := range s.inventory {
		if game := ticket.Interact(); game != nil {
			totalWinning += game.(*TicketGame).Win
			// Play winning sounds/graphics
		}
		s.inventory[index] = nil
	}
	s.money += totalWinning
}

func (s *Stats) CheckTicket(slot int) (t *Ticket) {
	ticket := s.inventory[slot]
	if ticket == nil {
		return nil
	}
	if game := ticket.Interact(); game != nil {
		s.money += game.(*TicketGame).Win
		// Play winning sounds/graphics
		fmt.Printf("You won: %s, $%d\n", game.(*TicketGame).Prize, game.(*TicketGame).Win)
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

		if ticket != nil && s.money >= ticket.Cost && s.findEmptySlot(ticket) {
			s.money -= ticket.Cost
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
