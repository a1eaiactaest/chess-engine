package engine

import (
	"fmt"
	"os"
	"sort"

	"github.com/notnil/chess"
)

const MaxVal = 100000

type Game struct {
	game              *chess.Game
	position          *chess.Position
	leavesExplored    int
	pieceValues       map[chess.PieceType]int
	pieceSquareTables map[chess.PieceType][]int
	positionStack     []*chess.Position
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
		positionStack:     make([]*chess.Position, 0, 100),
	}, nil
}

func (g *Game) pushPosition(pos *chess.Position) {
	g.positionStack = append(g.positionStack, pos)
}

func (g *Game) popPosition() {
	stackLen := len(g.positionStack)
	if stackLen > 0 {
		g.position = g.positionStack[stackLen-1]
		g.positionStack = g.positionStack[:stackLen-1]
	}
}

func (g *Game) Leaves() int {
	leaves := g.leavesExplored
	g.leavesExplored = 0
	return leaves
}

func (g *Game) GetLeaves() int {
	return g.leavesExplored
}

func (g *Game) evaluateMobility() int {
	// Count number of legal moves for each side
	whiteMoves := len(g.game.ValidMoves())
	if whiteMoves == 0 {
		return 0
	}
	// Use first valid move instead of empty move
	firstMove := g.game.ValidMoves()[0]
	g.game.Move(firstMove)
	blackMoves := len(g.game.ValidMoves())
	g.game.Position() // Reset position to previous state
	return (whiteMoves - blackMoves) * 10
}

func (g *Game) evaluateKingSafety() int {
	val := 0
	board := g.game.Position().Board()
	turn := g.game.Position().Turn()

	// Find king positions first
	var whiteKing, blackKing chess.Square
	for sq := 0; sq < 64; sq++ {
		p := board.Piece(chess.Square(sq))
		if p == chess.NoPiece {
			continue
		}
		if p.Type() == chess.King {
			if p.Color() == chess.White {
				whiteKing = chess.Square(sq)
			} else {
				blackKing = chess.Square(sq)
			}
		}
	}

	// Count attacking pieces near the king
	for sq := 0; sq < 64; sq++ {
		p := board.Piece(chess.Square(sq))
		if p == chess.NoPiece {
			continue
		}

		// Only evaluate pieces attacking the king
		if p.Color() != turn {
			kingSquare := whiteKing
			if p.Color() == chess.White {
				kingSquare = blackKing
			}

			// Calculate Manhattan distance to king
			kingFile := int(kingSquare) % 8
			kingRank := int(kingSquare) / 8
			pieceFile := sq % 8
			pieceRank := sq / 8
			distance := abs(kingFile-pieceFile) + abs(kingRank-pieceRank)

			if distance <= 2 {
				val -= 50 * int(p.Type())
			}
		}
	}

	return val
}

func (g *Game) evaluatePawnStructure() int {
	val := 0
	board := g.game.Position().Board()

	// Create a map to track pawns by file
	pawnsByFile := make(map[int][]chess.Square)
	for sq := 0; sq < 64; sq++ {
		p := board.Piece(chess.Square(sq))
		if p == chess.NoPiece || p.Type() != chess.Pawn {
			continue
		}
		file := sq % 8
		pawnsByFile[file] = append(pawnsByFile[file], chess.Square(sq))
	}

	// Evaluate doubled pawns
	for _, pawns := range pawnsByFile {
		if len(pawns) > 1 {
			val -= 20 * (len(pawns) - 1)
		}
	}

	// Evaluate isolated pawns
	for file, pawns := range pawnsByFile {
		hasNeighbor := false
		for offset := -1; offset <= 1; offset += 2 {
			neighborFile := file + offset
			if neighborFile >= 0 && neighborFile < 8 {
				if len(pawnsByFile[neighborFile]) > 0 {
					hasNeighbor = true
					break
				}
			}
		}
		if !hasNeighbor {
			val -= 30 * len(pawns)
		}
	}

	return val
}

func (g *Game) Evaluate() int {
	val := 0
	board := g.game.Position().Board()

	// Material and piece-square tables
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

	// Add mobility evaluation
	val += g.evaluateMobility()

	// Add king safety evaluation
	val += g.evaluateKingSafety()

	// Add pawn structure evaluation
	val += g.evaluatePawnStructure()

	return val
}

type MoveScore struct {
	moves []chess.Move
	score int
}

func (g *Game) moveScore(move *chess.Move) int {
	score := 0
	board := g.game.Position().Board()
	targetSq := move.S2()
	targetPiece := board.Piece(targetSq)

	// Captures
	if targetPiece != chess.NoPiece {
		score += 10 * int(g.pieceValues[targetPiece.Type()])
	}

	// Center control
	centerSquares := []chess.Square{chess.E4, chess.E5, chess.D4, chess.D5}
	for _, sq := range centerSquares {
		if targetSq == sq {
			score += 30
		}
	}

	// Development
	if g.game.Position().Turn() == chess.White {
		if move.S1() == chess.E2 && move.S2() == chess.E4 {
			score += 50
		}
		if move.S1() == chess.G1 && move.S2() == chess.F3 {
			score += 30
		}
		if move.S1() == chess.F1 && (move.S2() == chess.C4 || move.S2() == chess.B5) {
			score += 30
		}
	} else {
		if move.S1() == chess.E7 && move.S2() == chess.E5 {
			score += 50
		}
		if move.S1() == chess.B8 && move.S2() == chess.C6 {
			score += 30
		}
		if move.S1() == chess.F8 && (move.S2() == chess.C5 || move.S2() == chess.B4) {
			score += 30
		}
	}

	return score
}

func (g *Game) Minmax(
	fromBot int,
	depth int,
	lastMove *chess.Move,
	alpha int,
	beta int,
	isMax bool) MoveScore {

	moveTree := make([]chess.Move, 0, 10)

	if lastMove != nil {
		moveTree = append(moveTree, *lastMove)
	}

	if fromBot == 0 {
		return MoveScore{moveTree, g.Evaluate()}
	}

	moves := g.game.ValidMoves()

	// Sort moves for better alpha-beta pruning
	sort.Slice(moves, func(i, j int) bool {
		return g.moveScore(moves[i]) > g.moveScore(moves[j])
	})

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

			// Save current position
			oldFen := g.game.Position().String()

			// Make move
			g.game.Move(move)
			g.position = g.game.Position()

			result := g.Minmax(fromBot-1, depth+1, move, alpha, beta, false)

			// Restore position
			oldPos, _ := chess.FEN(oldFen)
			g.game = chess.NewGame(oldPos)
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

			// Save current position
			oldFen := g.game.Position().String()

			// Make move
			g.game.Move(move)
			g.position = g.game.Position()

			result := g.Minmax(fromBot-1, depth+1, move, alpha, beta, true)

			// Restore position
			oldPos, _ := chess.FEN(oldFen)
			g.game = chess.NewGame(oldPos)
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

func (g *Game) IDS(depth int, isMax bool) string {
	bestMove := ""
	_ = -MaxVal // Initialize but don't use bestScore since we only need the move

	// Start with a smaller depth and gradually increase
	maxDepth := min(depth, 4) // Limit maximum depth to prevent hanging
	for d := 1; d <= maxDepth; d++ {
		result := g.Minmax(d, 0, nil, -MaxVal, MaxVal, isMax)
		if len(result.moves) > 0 {
			// Validate that the move is legal before returning it
			move := result.moves[0]
			validMoves := g.game.ValidMoves()
			isValid := false
			for _, validMove := range validMoves {
				if validMove.S1() == move.S1() && validMove.S2() == move.S2() {
					isValid = true
					break
				}
			}
			if isValid {
				bestMove = move.String()
			} else {
				// If the move is invalid, try to find any legal move
				if len(validMoves) > 0 {
					bestMove = validMoves[0].String()
				}
			}
		}
	}

	// If we still don't have a valid move, try to find any legal move
	if bestMove == "" {
		validMoves := g.game.ValidMoves()
		if len(validMoves) > 0 {
			bestMove = validMoves[0].String()
		}
	}

	return bestMove
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func FeedbackEngine() {
	defaultFen := "r1b1k3/ppp1nQ2/4P1pN/2q5/8/6P1/5PBP/R3R1K1 b - - 2 28"
	fen := defaultFen

	if len(os.Args) > 1 {
		fen = os.Args[1]
	}

	engine, err := NewGame(fen)
	if err != nil {
		fmt.Printf("Error initializing engine: %v\n", err)
		os.Exit(1)
	}

	move := engine.IDS(5, false)
	fmt.Printf("Best move: %s for %s\n", move, engine.game.Position().Turn())
	fmt.Printf("Evaluation : %d\n", engine.Evaluate())
}
