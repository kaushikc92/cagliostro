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
	playTwoSides := false
	engineElo := 2200
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
					panic(err)
				}
			case "update":
				depth, err := strconv.Atoi(inputs[1])
				if err != nil {
					fmt.Printf("%v\n", err)
				} else {
					err = updatePosition(game.Position(), depth, inputs[2])
					if err != nil {
						panic(err)
					}
				}
			case "asyncupdate":
				depth, err := strconv.Atoi(inputs[1])
				if err != nil {
					fmt.Printf("%v\n", err)
				} else {
					err = asyncUpdatePosition(game.Position(), depth)
					if err != nil {
						panic(err)
					}
				}
			case "togglemode":
				playTwoSides = !playTwoSides
			case "elo":
				elo, err := strconv.Atoi(inputs[1])
				if err != nil {
					panic(err)
				}
				engineElo = elo
			case "switch":
				playerTurn = false
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
			if !playTwoSides {
				posData, _ := lichess.PositionData(game.FEN())
				moveStr := ""
				if posData.White + posData.Draws + posData.Black < 10 {
					moveStr, err = getMoveAtElo(game.Position(), engineElo, 10)
					if err != nil {
						panic(err)
					}
				} else {
					moveStr = chooseMove(posData)
				}
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
	pos, err := db.GetRepertoirePosition(fenString)
	if err != nil {
		switch err.(type) {
		case *db.ErrRecordDoesntExist :
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
	dbpos, err := db.GetRepertoirePosition(fenString)
	if err != nil {
		switch err.(type) {
		case *db.ErrRecordDoesntExist :
			bestMove, err := getMove(position, depth)
			if err != nil {
				return err
			}
			newDbpos := db.RepertoirePosition{
				Fen: fenString,
				BestMove: bestMove,
				Depth: depth,
				MyMove: myMove,
			}
			db.UpsertRepertoirePosition(newDbpos)
		default:
			return err
		}
	} else {
		if depth > dbpos.Depth {
			bestMove, err := getMove(position, depth)
			if err != nil {
				return err
			}
			newDbpos := db.RepertoirePosition{
				Fen: fenString,
				BestMove: bestMove,
				Depth: depth,
				MyMove: myMove,
			}
			db.UpsertRepertoirePosition(newDbpos)
		} else {
			newDbpos := db.RepertoirePosition{
				Fen: fenString,
				BestMove: dbpos.BestMove,
				Depth: dbpos.Depth,
				MyMove: myMove,
			}
			db.UpsertRepertoirePosition(newDbpos)
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
	setPos := uci.CmdPosition{Position: position}
	setGo := uci.CmdGo{Depth: depth}
	if err := eng.Run(uci.CmdUCINewGame, setPos, setGo); err != nil {
		return "", err
	}
	bestMove := eng.SearchResults().BestMove
	moveStr := bestMove.String()

	return moveStr, nil
}

func getMoveAtElo(position *chess.Position, elo int, moveTime int) (string, error) {
	eng, err := uci.New("stockfish")
	if err != nil {
		return "", err
	}
	defer eng.Close()
	setLimitStrengthOpt := uci.CmdSetOption{Name: "UCI_LimitStrength", Value: "true"}
	setEloOpt := uci.CmdSetOption{Name: "UCI_Elo", Value: strconv.Itoa(elo) }
	setPos := uci.CmdPosition{Position: position}
	setGo := uci.CmdGo{MoveTime: time.Second * time.Duration(moveTime)}
	if err := eng.Run(uci.CmdUCINewGame, setLimitStrengthOpt, setEloOpt, setPos, setGo); err != nil {
		return "", err
	}
	bestMove := eng.SearchResults().BestMove
	moveStr := bestMove.String()

	return moveStr, nil
}

func asyncUpdatePosition(position *chess.Position, depth int) error {
	upos := db.UpdatePosition {
		Fen: position.String(),
		Depth: depth,
	}
	return db.UpsertUpdatePosition(upos)
}
