package board

import (
	"fmt"
	"testing"
)

func TestKnightMoves(t *testing.T) {
	b := New()
	bCopy := *b
	fmt.Printf("%T, %T\n", b, bCopy)

	b.movePiece('b', 1, 'd', 5)
	moves := b.possibleMoves('d', 5)
	fmt.Println(moves)
	fmt.Println(b)
	moves[0](b)
	fmt.Println(b)

	fmt.Println(bCopy)
	bCopy.movePiece('b', 1, 'd', 5)
	fmt.Println(bCopy)
	// fmt.Println(bCopy)
	// moves[5](bCopy)
	// fmt.Println(bCopy)
}

func TestBishopMoves(t *testing.T) {
	b := New()
	b.movePiece('c', 1, 'd', 5)
	moves := b.possibleMoves('d', 5)
	fmt.Println(b)
	fmt.Println(moves, len(moves))
}

func TestRookMoves(t *testing.T) {
	b := New()
	b.movePiece('a', 1, 'd', 5)
	moves := b.possibleMoves('d', 5)
	fmt.Println(b)
	fmt.Println(moves, len(moves))
}

func TestQueenMoves(t *testing.T) {
	b := New()
	b.movePiece('d', 1, 'd', 5)
	moves := b.possibleMoves('d', 5)
	fmt.Println(b)
	fmt.Println(moves, len(moves))
}
