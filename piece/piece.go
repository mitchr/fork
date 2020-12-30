package piece

type Piece struct {
	Color Color
	Type  pieceType
}

func New(c Color, t pieceType) *Piece {
	return &Piece{c, t}
}

func (p Piece) String() string {
	return p.Color.String() + " " + p.Type.String()
}

type pieceType byte

func (p pieceType) String() string {
	return [...]string{"Null", "Pawn", "Bishop", "Knight", "Rook", "Queen", "King"}[p]
}

const (
	Null pieceType = iota // no type (no piece at this square)
	Pawn
	Bishop
	Knight
	Rook
	Queen
	King
)

type Color byte

const (
	Blank Color = iota // no color (no piece at this square)
	White
	Black
)

func (c Color) String() string {
	if c == 0 {
		return "Blank"
	} else if c == 1 {
		return "White"
	} else {
		return "Black"
	}
}
