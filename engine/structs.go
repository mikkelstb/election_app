package engine

import (
	"fmt"
	"sort"
)


type Election struct {
	Allocated_seats int
	Districts []District
	Parties []Party
	total_votes int
	method QuotientAlg
	Threshold_percent float32
	Threshold_seats int
}

type District struct {
	Name string
	Allocated_seats int
	Votes map[string]int
	DistrictParties []DistrictParty
	method QuotientAlg
	total_votes int
}

type DistrictParty struct {
	Name string
	party *Party
	district_votes int
	district_seats int
}


type Party struct {
	Name string
	votes int
	percentage float32
	district_seats int
	additional_seats int
}


func (p *Party) getTotalSeats() int {
	return p.district_seats + p.additional_seats
}


func (e *Election) getPartyRef(partyname string) *Party {
	for index := range e.Parties {
		if e.Parties[index].Name == partyname {
			return &e.Parties[index]
		}
	}
	return nil
}


// func (e *Election) addVotesToParty(p *Party, votes int) {
// 	p.votes += votes
// }


func (d *District) initParties(parties *[]Party) {
	d.DistrictParties = make([]DistrictParty, len(*parties))
	for index := range *parties {
		d.DistrictParties[index].Name = (*parties)[index].Name
		d.DistrictParties[index].party = &(*parties)[index]
		d.DistrictParties[index].district_votes = d.Votes[d.DistrictParties[index].Name]
		(*parties)[index].votes += d.Votes[d.DistrictParties[index].Name]
		d.total_votes += d.DistrictParties[index].district_votes
	}
}


func (d *District) GetAllSeats() {
	
	for i := 0; i < d.Allocated_seats; i++ {
		d.getNextSeat()
	}
}


func (d *District) getNextSeat() {
	largest_party := new(DistrictParty)

	for i := range d.DistrictParties {
		next_party := &(d.DistrictParties)[i]
		q1 := d.method.calculateQuotient(largest_party.district_votes, largest_party.district_seats)
		q2 := d.method.calculateQuotient(next_party.district_votes, next_party.district_seats)
		if q2 > q1 {
			largest_party = next_party
		}
	}
	largest_party.district_seats++
	largest_party.party.district_seats++
}


func (d *District) printResult() {

	fmt.Println(d.Name)
	sort.Slice(d.DistrictParties, func(i, j int) bool { return d.DistrictParties[i].district_votes > d.DistrictParties[j].district_votes})
	for _, party := range d.DistrictParties {
		percentage := float32(party.district_votes)/float32(d.total_votes)*100
		fmt.Printf("%-40s \t %7d %5.1f %% %3d \n", party.Name, party.district_votes, percentage, party.district_seats)
	}

	fmt.Printf("Total Votes \t\t\t\t\t %7d\n", d.total_votes)
	fmt.Println()
}


// 	sort.Slice(d.Parties, func(i int, j int) bool { return d.Parties[i].Votes > d.Parties[j].Votes})
// 	for _, p := range d.Parties {
// 		fmt.Printf("%-40s \t %7d  %5.1f %% %3d \n", p.Name, p.Votes, float32(p.Votes)*100/float32(d.Total_votes), p.seats_won)
// 	}


func (e *Election) printResult() {

}



type QuotientAlg interface {
	calculateQuotient(votes int, seats_won int) float32
}


type Dhont struct {}

func (dhont Dhont) calculateQuotient(votes int, seats_won int) float32 {
	return float32(votes)/float32(seats_won+1)
}


type SainteLague struct {}

func (saint_l SainteLague) calculateQuotient(votes int, seats_won int) float32{
	return float32(votes)/float32(seats_won*2 +1)
}
