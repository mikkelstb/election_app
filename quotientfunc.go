package main

/*
	Quotientfunt calculates a quotient to determine the order of which
	party or candidate should be awarded the next seat
	The function takes 3 parameters:
	int: number of votes
	int: number of seats already awarded
	float62: The number to manipulate number of seats if seats is 0
*/
type QuotientFunc interface {
	calculate(int, int, float32) float32
}

type dhont struct{}

func (d dhont) calculate(party_votes, seats int, first_divisor float32) float32 {
	if first_divisor > 0 && first_divisor < 2 && seats == 0 {
		return float32(party_votes) / float32(first_divisor)
	}
	return float32(party_votes) / float32(seats+1)
}

type sainteLague struct{}

func (s sainteLague) calculate(party_votes, seats int, first_divisor float32) float32 {
	if first_divisor > 0 && first_divisor < 2 && seats == 0 {
		return float32(party_votes) / float32(first_divisor)
	}
	return float32(party_votes) / (2*float32(seats) + 1)
}

type none struct{}

func (n none) calculate(party_votes, seats int, first_divisor float32) float32 {
	return 0
}
