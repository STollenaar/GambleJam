package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"sort"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Ticket struct {
	Name  string        `json:"name"`
	Cost  int           `json:"cost"`
	Games []*TicketGame `json:"games"`
}

type TicketGame struct {
	Prize string  `json:"prize"`
	Win   int     `json:"win"`
	Odds  float64 `json:"odds"`
}

var (
	Tickets      []*Ticket
	TicketAssets map[string]*ebiten.Image
	TicketNames  []string
)

func init() {
	data, err := os.ReadFile("./configs/tickets.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(data, &Tickets)
	if err != nil {
		log.Fatal(err)
	}

	// Sorting the slice by the 'name' field
	sort.Slice(Tickets, func(i, j int) bool {
		return Tickets[i].Name < Tickets[j].Name
	})

	TicketAssets = make(map[string]*ebiten.Image)
	assets, err := os.ReadDir("./assets/tickets")
	if err != nil {
		log.Fatal(err)
	}
	for _, asset := range assets {
		igbm, _, err := ebitenutil.NewImageFromFile(fmt.Sprintf("./assets/tickets/%s", asset.Name()))
		if err != nil {
			log.Fatal(err)
		}
		nameNorm := strings.Title(strings.ReplaceAll(strings.ReplaceAll(asset.Name(), "_", " "), ".png", ""))
		TicketAssets[nameNorm] = igbm
		TicketNames = append(TicketNames, nameNorm)
	}
	sort.Strings(TicketNames)
}

func (t *Ticket) Interact() interface{} {
	// Generate a random number between 0 and 1
	random := rand.Float64()

	var totalCumulative float64

	for _, game := range t.Games {
		totalCumulative += game.Odds
		if random <= totalCumulative {
			return game
		}
	}
	return nil
}
