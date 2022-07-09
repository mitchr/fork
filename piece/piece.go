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

func (p1 Piece) Equals(p2 Piece) bool {
	return p1.Color == p2.Color && p1.Type == p2.Type
}

type pieceType byte

func (p pieceType) String() string {
	return [...]string{"Pawn", "Bishop", "Knight", "Rook", "Queen", "King"}[p]
}

const (
	Pawn pieceType = iota
	Bishop
	Knight
	Rook
	Queen
	King
)

type Color byte

const (
	White Color = iota
	Black
)

func (c Color) String() string {
	if c == White {
		return "White"
	} else {
		return "Black"
	}
}
