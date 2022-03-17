package trainer

import(
	"fmt"
	"github.com/kaushikc92/chess"
)

func Interactive(fenString string) error {
	fen, err := chess.FEN(fenString)
	if err != nil {
		panic(err)
	}
	game := chess.NewGame(fen)
	fmt.Printf(game.Position().Board().Draw())
	return nil
}
