package puzzles

import(
	"fmt"
	"os"
	"bufio"
	"strings"
	"github.com/kaushikc92/cagliostro/pkg/db"
)

func Interactive(minRating int, maxRating int) error {
	for {
		// Retrieve puzzle in range and display
		puzzle, err := db.GetPuzzle(minRating, maxRating)
		if err != nil {
			panic(err)
		} else {
			fmt.Printf("Fen: %v", puzzle.Fen)
		}

		reader := bufio.NewReader(os.Stdin)
		inp, _ := reader.ReadString('\n')
		inp = inp[:len(inp) - 1]
		inputs := strings.Split(inp, " ")

		if inputs[0] == "n" {
			fmt.Println("New Puzzle")
		}
	}
	return nil
}
