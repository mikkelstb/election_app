package methods

type Method interface{
	LoadVotes(map[string]int)
	getNextMandate() string
	getAllMandates(numberOfSeats int) map[string]int
}

