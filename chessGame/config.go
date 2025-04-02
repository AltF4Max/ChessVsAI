package chessGame

import (
	"github.com/notnil/chess"
)

func PlayChess(game *chess.Game, move string) (bool, error) {
	err := game.MoveStr(move)
	if err != nil {
		return false, err
	}
	if game.Outcome() != chess.NoOutcome { //Game over
		return true, nil
	}
	return false, nil

}
