package main

import (
	"image"
)

func PerformClick(c *Client, pos image.Point) {
	c.SetMouseState(pos, MouseState{}.Set(MouseButtonLeft))
	c.SetMouseState(pos, MouseState{})
}

func PerformDoubleClick(c *Client, pos image.Point) {
	PerformClick(c, pos)
	PerformClick(c, pos)
}

type Direction int

const (
	DirectionUp Direction = -1 + iota
	DirectionNone
	DirectionDown
)

func Scroll(c *Client, d Direction) {
	if d == DirectionUp {
		c.SetMouseState(c.MousePosition(), MouseState{}.Set(MouseButtonWheelUp))
	} else if d == DirectionDown {
		c.SetMouseState(c.MousePosition(), MouseState{}.Set(MouseButtonWheelDown))
	}
	c.SetMouseState(c.MousePosition(), MouseState{})
}

func TypeString(c *Client, s string) {
	for _, k := range s {
		c.PressKey(int(k))
		c.ReleaseKey(int(k))
	}
}
