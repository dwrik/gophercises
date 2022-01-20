package cyoa

import (
	"encoding/json"
)

type Option struct {
	Text string
	Arc  string
}

type Arc struct {
	Title   string
	Story   []string
	Options []Option
}

type Story map[string]Arc

var s Story

func NextArc(arc string) Arc {
	if next, ok := s[arc]; ok {
		return next
	}
	return Arc{}
}

func ParseJSON(data []byte) error {
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	return nil
}
