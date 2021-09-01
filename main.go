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
	country.initParties([]string{"Soc.Dem.", "Liberale"})
	country.List()

	d := country.findDistrict("Danmark")
	//fmt.Println(d)
	d.SetVotes("Soc.Dem.", 50)
	d.SetVotes("Liberale", 100)

	d.addAllSeats()
	d.printVotes()

	//d.listAsSubdir(1)

	s := country.findDistrict("Sverige")
	s.SetVotes("Soc.Dem.", 25)
	s.SetVotes("Liberale", 50)

	s.SetVotes("Liberale", 0)
	d.SetVotes("Liberale", 0)

	d.addAllSeats()
	s.addAllSeats()

	s.printVotes()
	d.printVotes()
	country.printVotes()

	// parties := []string{"AP", "SP", "H", "KRF", "R", "SV", "FRP", "V", "MDG", "FNB"}
	// country.initParties(parties)

	// country.SetVotes("AP", 73122)
	// country.SetVotes("H", 92833)
	// country.SetVotes("MDG", 55772)
	// country.SetVotes("SV", 33258)
	// country.SetVotes("R", 26302)
	// country.SetVotes("FNB", 21346)
	// country.SetVotes("V", 21110)
	// country.SetVotes("FRP", 19272)
	// country.SetVotes("SP", 7980)
	// country.SetVotes("KRF", 6346)

	// country.addAllSeats()

	// country.printVotes()
}

func buildDistrict(js interface{}) District {
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
