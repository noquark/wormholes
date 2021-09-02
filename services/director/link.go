package main

type Link struct {
	ID     string `json:"id" redis:"id"`
	Target string `json:"target" redis:"target"`
	Tag    string `json:"tag" redis:"tag"`
}
