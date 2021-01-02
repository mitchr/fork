package board

import (
	"fmt"
	"sync"

	"github.com/mitchr/fork/piece"
)

type board struct {
	positions   [8][8]*piece.Piece
	currentTurn piece.Color
	turnCount   int       // max of 50 possibly
	lastMove    move      // string representing last move made (used for en passant)
	moveList    [][2]move // a list of 2-tuples, where the first element is white's move and second is black's
	// possibly need some kind of field for both the white and black king and rooks indicating if they have moved or not, which would disallow them from castling; but we might also consider keeping a list of all the moves made by each player and just searching this list for any king/rook movement if that player tries to castle
}

func New() *board {
	b := &board{}
	b.setup()
	return b
}

// two boards are equal if they have exactly the same piece positions
func (b1 *board) Equals(b2 *board) bool {
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if b1.positions[i][j].Type != b2.positions[i][j].Type || b1.positions[i][j].Color != b2.positions[i][j].Color {
				return false
			}
		}
	}
	return true
}

// execute a series of moves
func (b *board) Execute(move ...move) {
	for _, v := range move {
		// need to verify if move is legal
		// what defines a legal move?
		// 	movement is possible for the designated piece (i.e. rook cannot move diagonally)
		//	if a capture happens, that capture is able to take place
		//	this move does not place the current player's kind in check

		v(b)
	}
}

// returns board i,j coordinates from file,rank
func (b *board) fileAndRankToMatrix(file rune, rank int) (int, int) {
	return 8 - rank, int(file - 'a')
}

// finds the piece at file i and rank j (i.e. findPiece('c', 4))
// returns nil if no piece was found at that file and rank, or if file and rank are illformed
func (b *board) findPiece(file rune, rank int) *piece.Piece {
	// make sure we aren't accessing out of bounds
	if file > 'h' || file < 'a' || rank > 8 || rank < 1 {
		return nil
	}
	return b.positions[8-rank][file-'a']
}

// move the piece at f1,r1 to f2,r2
// if an opponent piece is at f2,r2, it will be captured
// does not do any move verification, so be careful
func (b *board) movePiece(f1 rune, r1 int, f2 rune, r2 int) {
	p := b.findPiece(f1, r1)
	if p == nil || p.Type == piece.Null { // out of bounds or no piece at position
		return
	}

	// leave p's position empty
	i, j := b.fileAndRankToMatrix(f1, r1)
	b.positions[i][j] = &piece.Piece{}
	// move p to new position
	i, j = b.fileAndRankToMatrix(f2, r2)
	b.positions[i][j] = p
}

// returns a list of possible moves for a piece at file f and rank r
func (b *board) possibleMoves(f rune, r int) []move {
	p := b.findPiece(f, r)
	if p == nil || p.Type == piece.Null { // no piece found for this position
		return nil
	}

	moves := []move{} // TODO: consider using LL instead of array here

	// TODO: before any piece can move, we have to make sure that this move will not put the current player's king in check
	switch p.Type {
	case piece.Pawn:
		switch p.Color {
		case piece.White:
			if o := b.findPiece(f, r+1); o != nil && o.Type == piece.Null { // can move 1 forward if unobstructed
				moves = append(moves, func(b *board) { b.movePiece(f, r, f, r+1) })
			}
			if r == 2 && b.findPiece(f, r+2).Type == piece.Null { // if a white pawn is still on the 2 rank, it can also move 2 ranks forward
				// TODO: should also enable some kind of en-passant flag here!
				moves = append(moves, func(b *board) { b.movePiece(f, r, f, r+2) })
			}
			if o := b.findPiece(f+1, r+1); o != nil && o.Color == piece.Black { // can capture black piece to the right
				moves = append(moves, func(b *board) { b.movePiece(f, r, f+1, r+1) })
			}
			if o := b.findPiece(f+1, r+1); o != nil && o.Color == piece.Black { // can capture black piece to the left
				moves = append(moves, func(b *board) { b.movePiece(f, r, f-1, r+1) })
			}
		case piece.Black:
			if o := b.findPiece(f, r-1); o != nil && o.Type == piece.Null {
				moves = append(moves, func(b *board) { b.movePiece(f, r, f, r-1) })
			}
			if r == 7 && b.findPiece(f, r-2).Type == piece.Null {
				// TODO: should also enable some kind of en-passant flag here!
				moves = append(moves, func(b *board) { b.movePiece(f, r, f, r-2) })
			}
			if o := b.findPiece(f+1, r-1); o != nil && o.Color == piece.White {
				moves = append(moves, func(b *board) { b.movePiece(f, r, f+1, r-1) })
			}
			if o := b.findPiece(f+1, r-1); o != nil && o.Color == piece.White {
				moves = append(moves, func(b *board) { b.movePiece(f, r, f-1, r-1) })
			}
		}

	case piece.Bishop:
		var wg sync.WaitGroup
		wg.Add(4)
		c := make(chan move)

		// right-up
		go func(c chan<- move) {
			defer wg.Done()
			for file, rank := f+1, r+1; file <= 'h' && rank <= 8; file, rank = file+1, rank+1 {
				o := b.findPiece(file, rank)
				if o.Type == piece.Null {
					c <- func(b *board) { b.movePiece(f, r, file, rank) }
				} else if o.Color != b.currentTurn {
					c <- func(b *board) { b.movePiece(f, r, file, rank) }
					break
				} else if o.Color == b.currentTurn {
					break
				}
			}
		}(c)

		// right-down
		go func(c chan<- move) {
			defer wg.Done()
			for file, rank := f+1, r-1; file <= 'h' && rank >= 1; file, rank = file+1, rank-1 {
				o := b.findPiece(file, rank)
				if o.Type == piece.Null {
					c <- func(b *board) { b.movePiece(f, r, file, rank) }
				} else if o.Color != b.currentTurn {
					c <- func(b *board) { b.movePiece(f, r, file, rank) }
					break
				} else if o.Color == b.currentTurn {
					break
				}
			}
		}(c)

		// left-up
		go func(c chan<- move) {
			defer wg.Done()
			for file, rank := f-1, r+1; file >= 'a' && rank <= 8; file, rank = file-1, rank+1 {
				o := b.findPiece(file, rank)
				if o.Type == piece.Null {
					c <- func(b *board) { b.movePiece(f, r, file, rank) }
				} else if o.Color != b.currentTurn {
					c <- func(b *board) { b.movePiece(f, r, file, rank) }
					break
				} else if o.Color == b.currentTurn {
					break
				}
			}
		}(c)

		// left-down
		go func(c chan<- move) {
			defer wg.Done()
			for file, rank := f-1, r-1; file >= 'a' && rank >= 1; file, rank = file-1, rank-1 {
				o := b.findPiece(file, rank)
				if o.Type == piece.Null {
					c <- func(b *board) { b.movePiece(f, r, file, rank) }
				} else if o.Color != b.currentTurn {
					c <- func(b *board) { b.movePiece(f, r, file, rank) }
					break
				} else if o.Color == b.currentTurn {
					break
				}
			}
		}(c)

		// reader that consume all moves from channel
		go func(c <-chan move) {
			for m := range c {
				moves = append(moves, m)
			}
		}(c)
		// wait for all goroutine to finish
		wg.Wait()

		// by closing c here, we can remove the reader goroutine from the call space
		// If we neglect to close c, there will be a continually blocking goroutine
		// in the background which leaks memory depending on how many times this
		// particular branch runs
		close(c)

	case piece.Knight:
		// the notations in the comments here are referring to White, but it doesn't matter because every Knight can always make the same 8 moves
		// if the spot is empty or a color opposite of the one whose current turn it is, that piece can move there
		locations := [][]interface{}{
			{f + 1, r + 2}, //^^->
			{f - 1, r + 2}, //^^<-
			{f - 2, r + 1}, //<-<-^
			{f - 2, r - 1}, //<-<-v
			{f + 2, r + 1}, //->->^
			{f + 2, r - 1}, //->->v
			{f - 1, r - 2}, //vv<-
			{f + 1, r - 2}, //vv->-
		}
		for _, v := range locations {
			if o := b.findPiece(v[0].(rune), v[1].(int)); o != nil && (o.Type == piece.Null || o.Color != b.currentTurn) {
				moves = append(moves, func(b *board) { b.movePiece(f, r, v[0].(rune), v[1].(int)) })
			}
		}

	case piece.Rook: // adopts same goroutine pattern as Bishop
		var wg sync.WaitGroup
		wg.Add(4)
		c := make(chan move)

		// check file above
		go func(c chan<- move) {
			defer wg.Done()
			for i := r + 1; i <= 8; i++ {
				o := b.findPiece(f, i)
				// if there's an empty space here, then we can move there
				if o.Type == piece.Null {
					c <- func(b *board) { b.movePiece(f, r, f, i) }
				} else if o.Color != b.currentTurn {
					// if we encounter an opponent piece, then we can move there, but we
					// should not look for other positions in this direction
					c <- func(b *board) { b.movePiece(f, r, f, i) }
					break
				} else if o.Color == b.currentTurn {
					// if we encounter the same colored piece, then we can stop looking
					break
				}
			}
		}(c)

		// check file below
		go func(c chan<- move) {
			defer wg.Done()
			for i := r - 1; i >= 1; i-- {
				o := b.findPiece(f, i)
				if o.Type == piece.Null {
					c <- func(b *board) { b.movePiece(f, r, f, i) }
				} else if o.Color != b.currentTurn {
					c <- func(b *board) { b.movePiece(f, r, f, i) }
					break
				} else if o.Color == b.currentTurn {
					break
				}
			}
		}(c)

		// check rank left
		go func(c chan<- move) {
			defer wg.Done()
			for j := f - 1; j >= 'a'; j-- {
				o := b.findPiece(j, r)
				if o.Type == piece.Null {
					c <- func(b *board) { b.movePiece(f, r, j, r) }
				} else if o.Color != b.currentTurn {
					c <- func(b *board) { b.movePiece(f, r, j, r) }
					break
				} else if o.Color == b.currentTurn {
					break
				}
			}
		}(c)

		// check rank right
		go func(c chan<- move) {
			defer wg.Done()
			for j := f + 1; j <= 'h'; j++ {
				o := b.findPiece(j, r)
				if o.Type == piece.Null {
					c <- func(b *board) { b.movePiece(f, r, j, r) }
				} else if o.Color != b.currentTurn {
					c <- func(b *board) { b.movePiece(f, r, j, r) }
					break
				} else if o.Color == b.currentTurn {
					break
				}
			}
		}(c)

		go func(c <-chan move) {
			for m := range c {
				moves = append(moves, m)
			}
		}(c)

		wg.Wait()

		close(c)

	case piece.Queen: // a queen is a combination of a bishop and a rook
		// mutate p to look like a bishop
		p.Type = piece.Bishop
		moves = b.possibleMoves(f, r)
		// mutate p to look like a rook
		p.Type = piece.Rook
		moves = append(moves, b.possibleMoves(f, r)...)
		// mutate back
		p.Type = piece.Queen

	case piece.King:
		locations := [][]interface{}{
			{f - 1, r + 1}, // <-^
			{f - 1, r},     // <-
			{f - 1, r - 1}, // <-v
			{f, r - 1},     // v
			{f + 1, r - 1}, // ->v
			{f + 1, r},     // ->
			{f + 1, r + 1}, // ->^
			{f, r + 1},     // ^
		}
		for _, v := range locations {
			if o := b.findPiece(v[0].(rune), v[1].(int)); o != nil && (o.Type == piece.Null || o.Color != b.currentTurn) {
				moves = append(moves, func(b *board) { b.movePiece(f, r, v[0].(rune), v[1].(int)) })
			}
		}
	}

	return moves
}

// depending on whose turn it is,
func (b *board) isKinginCheck() bool {
	return false
}

// setup for standard piece placement
func (b *board) setup() {
	// white goes first
	b.currentTurn = piece.White

	b.positions[0] = [8]*piece.Piece{piece.New(piece.Black, piece.Rook),
		piece.New(piece.Black, piece.Knight),
		piece.New(piece.Black, piece.Bishop),
		piece.New(piece.Black, piece.Queen),
		piece.New(piece.Black, piece.King),
		piece.New(piece.Black, piece.Bishop),
		piece.New(piece.Black, piece.Knight),
		piece.New(piece.Black, piece.Rook)}
	b.positions[1] = [8]*piece.Piece{piece.New(piece.Black, piece.Pawn),
		piece.New(piece.Black, piece.Pawn),
		piece.New(piece.Black, piece.Pawn),
		piece.New(piece.Black, piece.Pawn),
		piece.New(piece.Black, piece.Pawn),
		piece.New(piece.Black, piece.Pawn),
		piece.New(piece.Black, piece.Pawn),
		piece.New(piece.Black, piece.Pawn)}

	// fill all blank spaces with a null piece
	for i := 2; i < 6; i++ {
		for j := 0; j < 8; j++ {
			b.positions[i][j] = &piece.Piece{}
		}
	}

	b.positions[6] = [8]*piece.Piece{piece.New(piece.White, piece.Pawn),
		piece.New(piece.White, piece.Pawn),
		piece.New(piece.White, piece.Pawn),
		piece.New(piece.White, piece.Pawn),
		piece.New(piece.White, piece.Pawn),
		piece.New(piece.White, piece.Pawn),
		piece.New(piece.White, piece.Pawn),
		piece.New(piece.White, piece.Pawn)}
	b.positions[7] = [8]*piece.Piece{piece.New(piece.White, piece.Rook),
		piece.New(piece.White, piece.Knight),
		piece.New(piece.White, piece.Bishop),
		piece.New(piece.White, piece.Queen),
		piece.New(piece.White, piece.King),
		piece.New(piece.White, piece.Bishop),
		piece.New(piece.White, piece.Knight),
		piece.New(piece.White, piece.Rook)}
}

// sets all spaces to null
func (b *board) blank() {
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			b.positions[i][j] = &piece.Piece{}
		}
	}
}

func (b board) String() string {
	str := ""
	for i := 0; i < 8; i++ {
		str += fmt.Sprintf("%v |", 8-i)
		for j := 0; j < 8; j++ {
			p := b.positions[i][j]
			if p.Type == piece.Null {
				str += "* "
			} else if p.Type == piece.Knight { // K is already taken :/
				str += "N "
			} else {
				str += fmt.Sprintf("%s ", string(p.Type.String()[0]))
			}
		}
		str += "\n"
	}
	str += "   _______________\n"
	str += "   a b c d e f g h"
	return str
}
