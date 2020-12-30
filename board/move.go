package board

// a move is a function that applies a single chess move to the board and mutates it accordingly
type move func(*board)

// given a move in algebraic notation, return a corresponding move function
// might need a simple parse here for the AN grammar
func NewMove(b *board, s string) move {
	if s == "0-0" { // king-side castle

	} else if s == "0-0-0" { // queen-side castle
	}

	// pieceSymbol := s[0]
	switch s[0] {
	case 'N':
	case 'B':
	case 'R':
	case 'Q':
	case 'K':
	}

	return nil
}

func parseAlgebra(s string) {

}
