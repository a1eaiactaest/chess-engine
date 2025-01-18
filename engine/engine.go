package engine

import (
	"github.com/notnil/chess"
)

const MaxVal = 100000

type Game struct {
	game              *chess.Game
	position          *chess.Position
	leavesExplored    int
	pieceValues       map[chess.PieceType]int
	pieceSquareTables map[chess.PieceType][]int
}

func NewGame(fen string) (*Game, error) {
	fenPos, err := chess.FEN(fen)
	if err != nil {
		return nil, err
	}

	game := chess.NewGame(fenPos)

	pieceValues := map[chess.PieceType]int{
		chess.Pawn:   100,
		chess.Bishop: 300,
		chess.Knight: 300,
		chess.Rook:   500,
		chess.Queen:  900,
		chess.King:   0,
	}

	// ideally, this would be dynamic
	pieceSquareTables := map[chess.PieceType][]int{
		chess.Pawn: {
			0, 0, 0, 0, 0, 0, 0, 0, // 8
			50, 50, 50, 50, 50, 50, 50, 50, // 7
			10, 10, 20, 30, 30, 20, 10, 10, // 6
			5, 5, 10, 25, 25, 10, 5, 5, // 5
			0, 0, 0, 20, 20, 0, 0, 0, // 4
			5, -5, -10, 0, 0, -10, -5, 5, // 3
			5, 10, 10, -20, -20, 10, 10, 5, // 2
			0, 0, 0, 0, 0, 0, 0, 0, // 1
			//     a  b  c  d  e  f  g  h
		},
		chess.Knight: {
			-50, -40, -30, -30, -30, -30, -40, -50,
			-40, -20, 0, 0, 0, 0, -20, -40,
			-30, 0, 10, 15, 15, 10, 0, -30,
			-30, 5, 15, 20, 20, 15, 5, -30,
			-30, 0, 15, 20, 20, 15, 0, -30,
			-30, 5, 10, 15, 15, 10, 5, -30,
			-40, -20, 0, 5, 5, 0, -20, -40,
			-50, -40, -30, -30, -30, -30, -40, -50,
		},
		chess.Bishop: {
			-20, -10, -10, -10, -10, -10, -10, -20,
			-10, 0, 0, 0, 0, 0, 0, -10,
			-10, 0, 5, 10, 10, 5, 0, -10,
			-10, 5, 5, 10, 10, 5, 5, -10,
			-10, 0, 10, 10, 10, 10, 0, -10,
			-10, 10, 10, 10, 10, 10, 10, -10,
			-10, 5, 0, 0, 0, 0, 5, -10,
			-20, -10, -10, -10, -10, -10, -10, -20,
		},
		chess.Rook: {
			0, 0, 0, 0, 0, 0, 0, 0,
			5, 10, 10, 10, 10, 10, 10, 5,
			-5, 0, 0, 0, 0, 0, 0, -5,
			-5, 0, 0, 0, 0, 0, 0, -5,
			-5, 0, 0, 0, 0, 0, 0, -5,
			-5, 0, 0, 0, 0, 0, 0, -5,
			-5, 0, 0, 0, 0, 0, 0, -5,
			0, 0, 0, 5, 5, 0, 0, 0,
		},
		chess.Queen: {
			-20, -10, -10, -5, -5, -10, -10, -20,
			-10, 0, 0, 0, 0, 0, 0, -10,
			-10, 0, 5, 5, 5, 5, 0, -10,
			-5, 0, 5, 5, 5, 5, 0, -5,
			0, 0, 5, 5, 5, 5, 0, -5,
			-10, 5, 5, 5, 5, 5, 0, -10,
			-10, 0, 5, 0, 0, 0, 0, -10,
			-20, -10, -10, -5, -5, -10, -10, -20,
		},
		chess.King: {
			-30, -40, -40, -50, -50, -40, -40, -30,
			-30, -40, -40, -50, -50, -40, -40, -30,
			-30, -40, -40, -50, -50, -40, -40, -30,
			-30, -40, -40, -50, -50, -40, -40, -30,
			-20, -30, -30, -40, -40, -30, -30, -20,
			-10, -20, -20, -20, -20, -20, -20, -10,
			20, 20, 0, 0, 0, 0, 20, 20,
			20, 30, 10, 0, 0, 10, 30, 20,
		},
	}

	return &Game{
		game:              game,
		position:          game.Position(),
		leavesExplored:    0,
		pieceValues:       pieceValues,
		pieceSquareTables: pieceSquareTables,
	}, nil
}

func (g *Game) Leaves() int {
	leaves := g.leavesExplored
	g.leavesExplored = 0
	return leaves
}

func (g *Game) Evaluate() int {
	val := 0
	board := g.game.Position().Board()

	for sq := 0; sq < 64; sq++ {
		p := board.Piece(chess.Square(sq))
		if p == chess.NoPiece {
			continue
		}

		value := g.pieceValues[p.Type()]
		if p.Color() == chess.White {
			val += value
			if table, exists := g.pieceSquareTables[p.Type()]; exists {
				val += table[63-sq] // flip board
			}
		} else {
			val -= value
			if table, exists := g.pieceSquareTables[p.Type()]; exists {
				val -= table[sq]
			}
		}
	}
	return val
}

type MoveScore struct {
	moves []chess.Move
	score int
}

func (g *Game) Minmax(
	fromBot int,
	depth int,
	lastMove *chess.Move,
	alpha int,
	beta int,
	isMax bool) MoveScore {

	moveTree := make([]chess.Move, 0)

	if lastMove != nil {
		moveTree = append(moveTree, *lastMove)
	}

	if fromBot == 0 {
		return MoveScore{moveTree, g.Evaluate()}
	}

	moves := g.game.ValidMoves()

	if len(moves) == 0 {
		if g.game.Outcome() == chess.WhiteWon {
			return MoveScore{moveTree, MaxVal}
		} else if g.game.Outcome() == chess.BlackWon {
			return MoveScore{moveTree, -MaxVal}
		}
		return MoveScore{moveTree, 0}
	}

	var bestScore int
	var bestMove chess.Move

	if isMax {
		bestScore = -MaxVal
		for _, move := range moves {
			g.leavesExplored += 1

			oldFen, _ := chess.FEN(g.position.String()) // this is dumb but there's no other way
			// make move, eval, undo move
			g.game.Move(move) // make
			g.position = g.game.Position() // save state

			// eval state, NOTE: move maybe should be a & pointer
			result := g.Minmax(fromBot-1, depth+1, move, alpha, beta, false)

			g.game = chess.NewGame(oldFen)
			g.position = g.game.Position()

			if result.score > bestScore {
				bestScore = result.score
				bestMove = *move
				moveTree = result.moves
			}

			if result.score >= beta {
				moveTree = append(moveTree, bestMove)
				return MoveScore{moveTree, bestScore}
			}

			if result.score > alpha {
				alpha = result.score
			}
		}

	} else {
		bestScore = MaxVal
		for _, move := range moves {
			g.leavesExplored += 1

			oldFen, _ := chess.FEN(g.position.String())
			g.game.Move(move)
			g.position = g.game.Position()

			result := g.Minmax(fromBot-1, depth+1, move, alpha, beta, true)

			g.game = chess.NewGame(oldFen)
			g.position = g.game.Position()

			if result.score < bestScore {
				bestScore = result.score
				bestMove = *move
				moveTree = result.moves
			}
			if result.score <= alpha {
				moveTree = append(moveTree, bestMove)
				return MoveScore{moveTree, bestScore}
			}
			if result.score < beta {
				beta = result.score
			}
		}
	}
	moveTree = append(moveTree, bestMove)
	return MoveScore{moveTree, bestScore}
}

