package models

type Node struct {
	Address        string
	Port           string
	HtmlServerPort string
}

type Line struct {
	StartX int `json:"startX"`
	StartY int `json:"startY"`
	EndX   int `json:"endX"`
	EndY   int `json:"endY"`
}

type Message struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Message Line   `json:"message"`
}
