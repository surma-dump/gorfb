package main

import (
	"image"
	"strings"
)

func PerformClick(c Client, pos image.Point) {
	c.SendMessage(&PointerEventMessage{
		Position:   pos,
		MouseState: MouseState{}.Set(MouseButtonLeft),
	})
	c.SendMessage(&PointerEventMessage{
		Position:   pos,
		MouseState: MouseState{},
	})
}

func PerformDoubleClick(c Client, pos image.Point) {
	PerformClick(c, pos)
	PerformClick(c, pos)
}

type Direction int

const (
	DirectionUp Direction = -1 + iota
	DirectionNone
	DirectionDown
)

func Scroll(c Client, d Direction) {
	msg := &PointerEventMessage{
		Position:   c.LastMousePosition(),
		MouseState: MouseState{},
	}
	if d == DirectionUp {
		msg.MouseState = msg.MouseState.Set(MouseButtonWheelUp)
	} else if d == DirectionDown {
		msg.MouseState = msg.MouseState.Set(MouseButtonWheelDown)
	}
	c.SendMessage(msg)
	msg.MouseState = MouseState{}
	c.SendMessage(msg)
}

type state func(s string) (keys []int, remainder string, newState state)

func directState(s string) ([]int, string, state) {
	if len(s) == 0 {
		return []int{0}, "", nil
	}

	if s[0] == '\\' && len(s) >= 2 {
		return []int{int(s[1])}, s[2:], directState
	} else if s[0] == '[' {
		return []int{}, s[1:], composeState
	}
	return []int{int(s[0])}, s[1:], directState
}

var (
	KeyAliases = map[string]int{
		"CTRL":   xkControlL,
		"SHIFT":  xkShiftL,
		"ALT":    xkAltL,
		"SUPER":  xkSuperL,
		"META":   xkMetaL,
		"UP":     xkUp,
		"DOWN":   xkDown,
		"LEFT":   xkLeft,
		"RIGHT":  xkRight,
		"RETURN": xkReturn,
	}
)

func composeState(s string) ([]int, string, state) {
	keys := make([]int, 0, 4)

	if len(s) == 0 {
		s = "]"
	}

	token := ""
	for len(s) > 0 && s[0] != ']' {
		if s[0] == '\\' {
			s = s[1:]
		}
		token += s[0:1]
		s = s[1:]

		if s[0] == '+' || s[0] == ']' {
			var key int
			if keyAlias, ok := KeyAliases[strings.ToUpper(token)]; ok {
				key = keyAlias
			}
			if len(token) == 1 {
				key = int(token[0])
			}
			if key != 0 {
				keys = append(keys, key)
			}
			if s[0] == '+' {
				s = s[1:]
			}
			token = ""
		}
	}
	if len(s) > 0 {
		s = s[1:]
	}
	return keys, s, directState
}

func TypeString(c Client, s string) {
	for keys, r, state := directState(s); state != nil; keys, r, state = state(r) {
		msg := &KeyEventMessage{
			Pressed: true,
		}
		for _, k := range keys {
			msg.Key = k
			c.SendMessage(msg)
		}

		msg.Pressed = false
		for _, k := range keys {
			msg.Key = k
			c.SendMessage(msg)
		}
	}
}
