package trainer

import(
	"fmt"
	"time"
	"os"
	"bufio"
	"strings"
	"strconv"
	xrand "golang.org/x/exp/rand"
	"github.com/kaushikc92/chess"
	"github.com/kaushikc92/chess/uci"
	"github.com/kaushikc92/cagliostro/pkg/lichess"
	"github.com/kaushikc92/cagliostro/pkg/db"
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
			case "":
				continue
			case "draw":
				fmt.Printf(game.Position().Board().Draw())
			case "check":
				err := checkPosition(game.Position().String())
				if err != nil {
					fmt.Printf("%v\n", err)
				}
			case "update":
				depth, err := strconv.Atoi(inputs[1])
				if err != nil {
					fmt.Printf("%v\n", err)
				} else {
					err = updatePosition(game.Position(), depth, inputs[2])
					if err != nil {
						fmt.Printf("%v\n", err)
					}
				}
			default:
				notation := chess.UCINotation{}
				move, err := notation.Decode(game.Position(), inputs[0])
				if err != nil {
					fmt.Printf("%v\n", err)
					continue
				}
				err = game.Move(move)
				if err != nil {
					fmt.Printf("%v\n", err)
				} else {
					playerTurn = false
				}
			}
		} else {
			posData, _ := lichess.PositionData(game.FEN())
			moveStr := castleCorrect(chooseMove(posData))
			notation := chess.UCINotation{}
			move, err := notation.Decode(game.Position(), moveStr)
			if err != nil {
				panic(err)
			}
			fmt.Printf("%v\n", moveStr)
			err = game.Move(move)
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
	return posData.Moves[int(categorical.Rand())].Uci
}

func checkPosition(fenString string) error {
	pos, err := db.GetPosition(fenString)
	if err != nil {
		switch err.(type) {
		case *db.ErrPositionDoesntExist :
			fmt.Print("This position does not exist in the database\n")
		default:
			return err
		}
	} else {
		fmt.Printf("%+v\n", pos)
	}
	return nil
}

func updatePosition(position *chess.Position, depth int, myMove string) error {
	fenString := position.String()
	dbpos, err := db.GetPosition(fenString)
	if err != nil {
		switch err.(type) {
		case *db.ErrPositionDoesntExist :
			bestMove, err := getMove(position, depth)
			if err != nil {
				return err
			}
			newDbpos := db.Position{
				Fen: fenString,
				BestMove: bestMove,
				Depth: depth,
				MyMove: myMove,
			}
			db.UpsertPosition(newDbpos)
		default:
			return err
		}
	} else {
		if depth > dbpos.Depth {
			bestMove, err := getMove(position, depth)
			if err != nil {
				return err
			}
			newDbpos := db.Position{
				Fen: fenString,
				BestMove: bestMove,
				Depth: depth,
				MyMove: myMove,
			}
			db.UpsertPosition(newDbpos)
		} else {
			newDbpos := db.Position{
				Fen: fenString,
				BestMove: dbpos.BestMove,
				Depth: dbpos.Depth,
				MyMove: myMove,
			}
			db.UpsertPosition(newDbpos)
		}
	}
	return nil
}

func getMove(position *chess.Position, depth int) (string, error) {
	eng, err := uci.New("stockfish")
	if err != nil {
		return "", err
	}
	defer eng.Close()
	if err != nil {
		return "", err
	}
	setPos := uci.CmdPosition{Position: position}
	setGo := uci.CmdGo{Depth: depth}
	if err := eng.Run(uci.CmdUCINewGame, setPos, setGo, uci.CmdEval); err != nil {
		return "", err
	}
	bestMove := eng.SearchResults().BestMove
	moveStr := bestMove.String()

	return moveStr, nil
}

func castleCorrect(moveStr string) string {
	switch moveStr{
	case "e1h1":
		return "e1g1"
	case "e8h8":
		return "e8g8"
	case "e1a1":
		return "e1c1"
	case "e8a8":
		return "e8c8"
	default:
		return moveStr
	}
}
