package engine

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)


func ReadElectionFile(filename string) *Election {

	json_string, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading: " + err.Error())
	}

	election := Election{}
	json.Unmarshal(json_string, &election)

	return &election
}



func (e *Election) CalculateDistrictSeats() {
	
	for index := range e.Districts {
		e.Districts[index].initParties(&e.Parties)
		e.Districts[index].method = Dhont{}
		e.Districts[index].GetAllSeats()
		e.total_votes += e.Districts[index].total_votes
	}
	for index := range e.Parties {
		e.Parties[index].percentage = float32(e.Parties[index].votes)*100/float32(e.total_votes)
	}
}


func (e *Election) calculateAdditionalSeats() {

	e.method = Dhont{}
	for i := 0; i < e.Allocated_seats; i++ {
		e.getNextSeat()
	}
}


func (e *Election) getNextSeat() {
	largest_party := new(Party)
	for i := range e.Parties {
		next_party := &e.Parties[i]
		if(next_party.percentage > e.Threshold_percent) || (next_party.district_seats >= e.Threshold_seats) {
			q1 := e.method.calculateQuotient(largest_party.votes, largest_party.getTotalSeats())
			q2 := e.method.calculateQuotient(next_party.votes, next_party.getTotalSeats())
			if q2 > q1 {
				largest_party = next_party
			}
		}
	}
	fmt.Println(largest_party.Name)
	largest_party.additional_seats++
}


func (e *Election) PrintDistrictResults() {
	
	for index := range e.Districts {
		e.Districts[index].printResult()
	}

	sort.Slice(e.Parties, func(i, j int) bool {return e.Parties[i].votes > e.Parties[j].votes})

	for _, party := range e.Parties {
		fmt.Printf("%-40s \t %7d %5.1f %% %3d %3d %3d\n", party.Name, party.votes, party.percentage, party.district_seats, party.additional_seats, party.getTotalSeats())
	}
	fmt.Printf("Total Votes \t\t\t\t\t %7d\n", e.total_votes)
}



