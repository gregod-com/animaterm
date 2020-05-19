package animaterm

// Position 2D position
type Position struct {
	x      int
	y      int
	offset int
}

// CreatePos takes x and y coodinates in percent
// and normalizes input to values between 0 and 100
func CreatePos(x int, y int) IRelativePosition {
	pos := &Position{offset: 0}
	pos.SetXandY(x, y)
	return pos
}

// GetX  see IRelativePosition
func (p *Position) GetX() int {
	return p.x
}

// GetY see IRelativePosition
func (p *Position) GetY() int {
	return p.y
}

// GetXandY see IRelativePosition
func (p *Position) GetXandY() (int, int) {
	return p.x, p.y
}

// SetX see IRelativePosition
func (p *Position) SetX(percentx int) error {
	if percentx > 100 {
		percentx = 100
	}
	if percentx < -100 {
		percentx = -100
	}
	p.x = percentx
	return nil
}

// SetY see IRelativePosition
func (p *Position) SetY(percenty int) error {
	if percenty > 100 {
		percenty = 100
	}
	if percenty < -100 {
		percenty = -100
	}
	p.y = percenty
	return nil
}

// SetXandY see IRelativePosition
func (p *Position) SetXandY(percentx int, percenty int) error {
	p.SetX(percentx)
	p.SetY(percenty)
	return nil
}

// IncrementOffset see IRelativePosition
func (p *Position) IncrementOffset() IRelativePosition {
	p.offset++
	return p
}

// ResetOffset see IRelativePosition
func (p *Position) ResetOffset() IRelativePosition {
	p.offset = 0
	return p
}

// GetOffset see IRelativePosition
func (p *Position) GetOffset() int {
	return p.offset
}

// SetOffset see IRelativePosition
func (p *Position) SetOffset(offset int) IRelativePosition {
	p.offset = offset
	return p
}

// DistanceTo see IRelativePosition
func (p *Position) DistanceTo(p2 IRelativePosition) IRelativePosition {
	distance := CreatePos(p2.GetX()-p.GetX(), p2.GetY()-p.GetY())
	return distance
}

// AddDistance see IRelativePosition
func (p *Position) AddDistance(p2 IRelativePosition) IRelativePosition {
	newpos := CreatePos(p.GetX()+p2.GetX(), p.GetY()+p2.GetY())
	newpos.SetOffset(p.GetOffset())
	return newpos
}

// MultiplyWith see IRelativePosition
func (p *Position) MultiplyWith(factor float32) IRelativePosition {
	newpos := CreatePos(int(float32(p.GetX())*factor), int(float32(p.GetY())*factor))
	newpos.SetOffset(p.GetOffset())
	return newpos
}
