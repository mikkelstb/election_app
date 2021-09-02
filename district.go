package main

import (
	"fmt"
	"sort"
	"strings"
)

type District struct {
	Name            string
	Seats           int
	AdditionalSeats int
	Threshold       float32
	PotentialVoters int
	SeatCalculator  QuotientFunc
	FirstDivisor    float32
	Parties         map[string]*DistrictVote
	SubDistricts    []District
	VoteChannel     chan VoteEvent
	//SeatChannel     chan SeatEvent
	TotalVotes int
}

func NewDistrict(name string, seats, additional_seats int, qf QuotientFunc, threshold float32, first_divisor float32, parties []string) *District {

	d := District{
		Name:            name,
		Seats:           seats,
		AdditionalSeats: additional_seats,
		Threshold:       threshold,
		SeatCalculator:  qf,
		FirstDivisor:    first_divisor,
		SubDistricts:    []District{},
		Parties:         make(map[string]*DistrictVote),
	}
	return &d
}

func (d *District) findDistrict(name string) *District {

	if name == d.Name {
		return d
	}
	for index := range d.SubDistricts {
		fd := d.SubDistricts[index].findDistrict(name)
		if fd != nil {
			return fd
		}
	}
	return nil
}

func (d *District) addSubdistrict(sd *District) {
	sd.VoteChannel = make(chan VoteEvent)
	go d.updateVotes(sd.VoteChannel)
	d.SubDistricts = append(d.SubDistricts, *sd)
}

func (d *District) SetVotes(party string, votes int) {
	difference := votes - d.Parties[party].votes
	//fmt.Println(difference)
	d.Parties[party].votes = votes
	d.VoteChannel <- VoteEvent{party: party, vote_difference: difference}
	d.TotalVotes += difference
}

func (d *District) updateVotes(data <-chan VoteEvent) {
	for {
		val := <-data
		d.Parties[val.party].votes += val.vote_difference
		d.Parties[val.party].seats += val.seat_difference
		//fmt.Println(val.seat_difference)
		d.TotalVotes += val.vote_difference
	}
}

func (d *District) List() {
	fmt.Printf("%-30s %3d %3d\n", d.Name, d.Seats, d.AdditionalSeats)
	for _, sub_district := range d.SubDistricts {
		sub_district.listAsSubdir(1)
	}
	s, as := d.getTotalSeats()
	fmt.Printf("%-30s %3d %3d\n", "Totalt:", s, as)
}

func (d *District) listAsSubdir(debth int) {
	fmt.Printf("%-30s %3d %3d\n", (strings.Repeat("-", debth) + d.Name), d.Seats, d.AdditionalSeats)
	debth++
	for _, sub_district := range d.SubDistricts {
		sub_district.listAsSubdir(debth)
	}
}

func (d *District) getTotalSeats() (int, int) {
	var seats, additional_seats int
	seats = d.Seats
	additional_seats = d.AdditionalSeats
	for _, sub_dist := range d.SubDistricts {
		s, as := sub_dist.getTotalSeats()
		seats += s
		additional_seats += as
	}
	return seats, additional_seats
}

func (d *District) initParties(new_parties []string) {
	for index := range new_parties {
		d.Parties[new_parties[index]] = new(DistrictVote)
		for _, sd := range d.SubDistricts {
			sd.initParties(new_parties)
		}
	}
}

func (d *District) SetSeatCalculator(qf QuotientFunc) {
	d.SeatCalculator = qf
}

func (d *District) calculateQuotient(party string) float32 {
	return d.SeatCalculator.calculate(d.Parties[party].votes, d.Parties[party].seats+d.Parties[party].additional_seats, d.FirstDivisor)
}

/*
	addAllAdditionalSeats()
	For all AdditionalSeats in the district, these are added taking into account how many
	District seats each party has won within the subdistricts
*/
func (d *District) addAllAdditionalSeats() {

	for i := range d.Parties {
		if float32(d.Parties[i].votes)/float32(d.TotalVotes)*100 < d.Threshold {
			d.Parties[i].passes_threshold = false
		} else {
			d.Parties[i].passes_threshold = true
		}
		if strings.HasSuffix(i, "*") {
			d.Parties[i].passes_threshold = false
		}
	}

	for i := 0; i < d.AdditionalSeats; i++ {
		d.addAdditionalSeat()
	}
}

func (d *District) addAdditionalSeat() {
	var largest_party string
	for party := range d.Parties {
		if _, ok := d.Parties[largest_party]; ok {
			if d.Parties[party].passes_threshold {
				if d.calculateQuotient(largest_party) < d.calculateQuotient(party) {
					largest_party = party
				}
			}
		} else {
			if d.Parties[party].passes_threshold {
				largest_party = party
			}
		}
	}
	d.Parties[largest_party].additional_seats++
}

/*
	addSeat() calculates the quota for each party, and adds 1 to the seat count of the largest
*/
func (d *District) addSeat() {
	var largest_party string
	for party := range d.Parties {
		if _, ok := d.Parties[largest_party]; ok {
			if d.Parties[party].passes_threshold {
				if d.calculateQuotient(largest_party) < d.calculateQuotient(party) {
					largest_party = party
				}
			}
		} else {
			if d.Parties[party].passes_threshold {
				largest_party = party
			}
		}
	}
	d.Parties[largest_party].seats++
	d.VoteChannel <- VoteEvent{party: largest_party, seat_difference: 1}
}

/* addAllSeats sets all seats to 0, and calculates the seats according to partyvotes and district rules*/
func (d *District) addAllSeats() {

	for i := range d.Parties {
		if float32(d.Parties[i].votes)/float32(d.TotalVotes)*100 < d.Threshold {
			d.Parties[i].passes_threshold = false
		} else {
			d.Parties[i].passes_threshold = true
		}
		if strings.HasSuffix(i, "*") {
			d.Parties[i].passes_threshold = false
		}
	}

	for i := range d.SubDistricts {
		d.SubDistricts[i].addAllSeats()
	}

	for i := range d.Parties {
		d.Parties[i].seats = 0
	}
	for i := 0; i < d.Seats; i++ {
		d.addSeat()
	}
}

func (d *District) printVotes() {
	fmt.Println("Votes for: " + d.Name)

	party_list := make([]string, 0, len(d.Parties))
	for k := range d.Parties {
		party_list = append(party_list, k)
	}
	sort.Strings(party_list)

	for _, party := range party_list {
		stat := d.Parties[party]
		if d.AdditionalSeats == 0 {
			fmt.Printf("%-30v %7d %5.1f%% %3d\n", party, stat.votes, float32(stat.votes)/float32(d.TotalVotes)*100, stat.seats)
		} else {
			fmt.Printf("%-30v %7d %5.1f%% %3d %3d %3d\n", party, stat.votes, float32(stat.votes)/float32(d.TotalVotes)*100, stat.seats, stat.additional_seats, stat.seats+stat.additional_seats)
		}
	}
	fmt.Printf("Total votes: %d\n", d.TotalVotes)
	fmt.Println()
}
