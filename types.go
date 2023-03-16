package main

type Message struct {
	Fact string `json:"fact"`
}

type Circle struct {
	id    string
	name  string
	users int
}
