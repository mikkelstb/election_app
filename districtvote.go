package main

type DistrictVote struct {
	votes            int
	seats            int
	additional_seats int
	passes_threshold bool
}

type VoteEvent struct {
	party           string
	vote_difference int
}
