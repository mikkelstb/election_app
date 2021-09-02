package main

import (
	"encoding/json"
	"os"
)

func main() {

	//fmt.Println("V2 starting")

	district_file, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	//d := new(District)
	var js interface{}
	err = json.Unmarshal(district_file, &js)
	if err != nil {
		panic(err)
	}

	country := buildDistrict(js)

	country.initParties([]string{
		"Andre*",
		"Arbeiderpartiet",
		"Fremskrittspartiet",
		"Høyre",
		"Kristelig Folkeparti",
		"Miljøpartiet De Grønne",
		"Rødt",
		"Senterpartiet",
		"Sosialistisk Venstreparti",
		"Venstre",
	})

	votes := readVoteFile("resources/vote_files/stortinget_norstat.csv")

	for _, dis_vote := range votes.districts {
		district := country.findDistrict(dis_vote.name)
		//fmt.Println("Found: " + district.Name)
		for party, vote := range dis_vote.parties {
			district.SetVotes(party, vote)
		}
	}

	for index := range country.SubDistricts {
		country.SubDistricts[index].addAllSeats()
		country.SubDistricts[index].printVotes()
	}

	country.addAllAdditionalSeats()

	country.printVotes()
}

func buildDistrict(js interface{}) *District {
	data := js.(map[string]interface{})

	name, exists := data["Name"].(string)
	if !exists {
		panic("Error Name is missing from District")
	}
	seats := checkInt(data["Seats"], 0)
	additional_seats := checkInt(data["AdditionalSeats"], 0)
	threshold := checkFloat32(data["Threshold"], 0.0)
	first_divisor := checkFloat32(data["FirstDivisor"], 1)

	var qf QuotientFunc
	qf_string, exists := data["QuotientFunc"]
	if exists {
		switch qf_string.(string) {
		case "dhont":
			qf = dhont{}
		case "sainte_lague":
			qf = sainteLague{}
		}
	} else {
		qf = none{}
	}

	d := NewDistrict(name, seats, additional_seats, qf, threshold, first_divisor, nil)

	sub_districts, ok := data["SubDistricts"]
	if ok {
		for _, sd := range sub_districts.([]interface{}) {
			d.addSubdistrict(buildDistrict(sd))
		}
	}
	return d
}

func checkInt(i interface{}, deafult int) int {
	if i == nil {
		return deafult
	} else {
		return int(i.(float64))
	}
}

func checkFloat32(i interface{}, deafult float32) float32 {
	if i == nil {
		return deafult
	} else {
		return float32(i.(float64))
	}
}
