package trainer

import(
	"fmt"
	"time"
	"os"
	"bufio"
	"strings"
	xrand "golang.org/x/exp/rand"
	"github.com/kaushikc92/chess"
	"github.com/kaushikc92/cagliostro/pkg/lichess"
	"gonum.org/v1/gonum/stat/distuv"

)

func Interactive(fenString string) error {
	fen, err := chess.FEN(fenString)
	if err != nil {
		panic(err)
	}
	game := chess.NewGame(fen)
	playerTurn := true
	for {
		if playerTurn {
			fmt.Print("Enter your move:\n")
			reader := bufio.NewReader(os.Stdin)
			inp, _ := reader.ReadString('\n')
			inp = inp[:len(inp) - 1]
			inputs := strings.Split(inp, " ")
			switch inputs[0] {
			case "draw":
				fmt.Printf(game.Position().Board().Draw())
			case "check":
			case "update":
			case "eval":
			case "force":
			default:
				err := game.MoveStr(inputs[0])
				if err != nil {
					fmt.Printf("%v", err)
				} else {
					playerTurn = false
				}
			}
		} else {
			posData, _ := lichess.PositionData(game.FEN())
			moveStr := chooseMove(posData)
			fmt.Printf("%v\n", moveStr)
			err := game.MoveStr(moveStr)
			if err != nil {
				panic(err)
			}
			playerTurn = true
		}
	}

	return nil
}

func chooseMove(posData lichess.PositionDataResults) string {
	n := len(posData.Moves)
	weights := make([]float64, n)
	for i:= 0; i<n; i++ {
		move := posData.Moves[i]
		weights[i] = float64(move.White + move.Draws + move.Black)
	}
	source := xrand.NewSource(uint64(time.Now().UnixNano()))
	categorical := distuv.NewCategorical(weights, source)
	return posData.Moves[int(categorical.Rand())].San
}
