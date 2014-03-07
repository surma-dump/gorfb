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

type state func(s string) (keys []int, remainder string, newState state)

func directState(s string) ([]int, string, state) {
	if len(s) == 0 {
		return []int{0}, "", nil
	}

	if s[0] == '\\' && len(s) >= 2 {
		return []int{int(s[1])}, s[2:], directState
	} else if s[0] == '[' {
		return []int{-1}, s[1:], composeState
	}
	return []int{int(s[0])}, s[1:], directState
}

var (
	KeyAliases = map[string]int{
		"Ctrl":   xkControlL,
		"Shift":  xkShiftL,
		"Alt":    xkAltL,
		"Super":  xkSuperL,
		"Meta":   xkMetaL,
		"Up":     xkUp,
		"Down":   xkDown,
		"Left":   xkLeft,
		"Right":  xkRight,
		"Return": xkReturn,
	}
)

func composeState(s string) ([]int, string, state) {
	if len(s) == 0 {
		return []int{0}, "", nil
	}
	token := ""
	for s[0] != ']' && s[0] != '+' {
		if s[0] == '\\' {
			s = s[1:]
		}
		token += s[0:1]
		s = s[1:]
	}
	key, ok := KeyAliases[token]
	if !ok {
		key = -1
	}
	if len(token) == 1 {
		key = int(token[0])
	}
	if s[0] == ']' {
		return []int{key}, s[1:], directState
	}
	additionalKeys, r, nextState := composeState(s[1:])
	return append([]int{key}, additionalKeys...), r, nextState
}

func TypeString(c *Client, s string) {
	for keys, r, state := directState(s); state != nil; keys, r, state = state(r) {
		for _, k := range keys {
			c.PressKey(k)
		}
		for _, k := range keys {
			c.ReleaseKey(k)
		}
	}
}
