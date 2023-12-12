package internal

type Pos struct {
	X, Y float32
}

func NewPos(x, y float32) Pos {
	return Pos{x, y}
}

func (p Pos) Subtract(p2 Pos) Pos {
	return Pos{X: p2.X - p.X, Y: p2.Y - p.Y}
}

func (p Pos) Add(p2 Pos) Pos {
	return Pos{X: p2.X + p.X, Y: p2.Y + p.Y}
}

type WInput struct {
	IsPressed      bool
	IsJustPressed  bool
	IsJustReleased bool
	Pos            Pos
	PressedPos     Pos
}

func (i *WInput) UpdatePosition(x, y int) {
	i.Pos.X = float32(x)
	i.Pos.Y = float32(y)
}
func (i *WInput) UpdatePressedPosition(x, y int) {
	i.PressedPos.X = float32(x)
	i.PressedPos.Y = float32(y)
}
