package methods

type Result struct {
	Parties []Party
	quotients map[string]float64
	Threshold float64
	Total_votes int
}


type Party struct {
	Name string
	Votes int
	Percentage float64
	Seats int
}


func (this *Result) LoadVotes(votes map[string]int) {

	this.quotients = make(map[string]float64)
	this.Parties = make([]Party, 0)

	for party_name, party_votes := range votes {
		party := new(Party)
		party.Name = party_name
		party.Votes = party_votes
		this.Total_votes = this.Total_votes + party_votes
		this.Parties = append(this.Parties, *party)
	}

	for x:= range this.Parties {
		this.Parties[x].Percentage = float64(this.Parties[x].Votes)*100/float64(this.Total_votes)
	}
}


func (this *Result) SetThreshold (threshold float64) {
	this.Threshold = threshold
}

func (this *Result) GetNextMandate() string {
	this.calculateQuotients()
	largest := this.getLargestParty()
	largest.Seats++
	return largest.Name
}

func (this *Result) GetAllMandates(numberOfSeats int) []Party {
	for x := 0; x < numberOfSeats; x++ {
		this.GetNextMandate()
	}
	return this.Parties
}

// Calculate from the formula: Votes/(Seats Awarded so far + 1)
func (this *Result) calculateQuotients() {
	for x := range this.Parties {
		if this.Parties[x].Percentage > this.Threshold {
			this.quotients[this.Parties[x].Name] = float64(this.Parties[x].Votes) / float64((this.Parties[x].Seats + 1))
		}
	}
}

func (this *Result) getLargestParty() *Party {
	var largest *Party

	for x := range this.Parties {
		if largest != nil {
			if this.quotients[this.Parties[x].Name] > this.quotients[largest.Name] {
				largest = &this.Parties[x]
			}
		} else {
			largest = &this.Parties[x]
		}
	}
	return largest
}
