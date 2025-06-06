package main

type Camera struct {
	X, Y   float64
	Zoom   float64
	Width  int
	Height int
}

func (c *Camera) CenterOn(x, y float64) {
	c.X = x - float64(c.Width)/(2*c.Zoom)
	c.Y = y - float64(c.Height)/(2*c.Zoom)
}
