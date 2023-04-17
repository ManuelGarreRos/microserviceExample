package main

type Message struct {
	Fact string `json:"fact"`
}

type Circle struct {
	ID    string
	Name  string
	Users int
}
