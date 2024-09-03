package internal

import "fmt"

type Item interface {
	Interact() interface{}
}

type Stats struct {
	money int
	day   int

	inventory []Item
	home      *Home
	store     *Store
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

func (s *Stats) CheckTicket(slot int) {
	ticket := s.inventory[slot]
	if game := ticket.Interact(); game != nil {
		s.money += game.(*TicketGame).Win
		// Play winning sounds/graphics
		fmt.Printf("You won: %s, $%d\n", game.(*TicketGame).Prize, game.(*TicketGame).Win)
	} else {
		fmt.Println("Better luck next time")
	}
	s.inventory[slot] = nil

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
