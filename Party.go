package main

type Party struct {
	Name        string
	Abbrevation string
	Id          int
}

func (p *Party) String() string {
	return p.Name
}
