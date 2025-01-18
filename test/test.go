package test

import (
  "fmt"
  "github.com/notnil/chess"
)

func TestGochess() {
  game := chess.NewGame()
  moves := []string{"e4", "e5", "Nf3", "Nc6", "Bb5"}

  for _, move := range moves {
    err := game.MoveStr(move)
    if err != nil {
      fmt.Printf("error making move %s: %v\n", move, err)
      return
    }
  }

  fmt.Printf("current fen: %s\n", game.Position().String())
  fmt.Printf("\nPGN:\n%s\n", game.String())

  if game.Outcome() != chess.NoOutcome {
    fmt.Printf("\nGame is over. result :%s\n", game.Outcome())
  }

  fmt.Printf("\nCurrent board position:\n%s\n", game.Position().Board().Draw())

  validMoves := game.ValidMoves()
  fmt.Printf("\nn of valid moves: %d\n", len(validMoves))
  for i, move := range validMoves {
    if i < 5 {
      fmt.Printf("- %s\n", move.String())
    }
  }
}
