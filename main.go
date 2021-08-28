package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"mikkelstb/election_app/methods"
	"os"
	"sort"
)


var finished bool
var js_results json_data

type json_data struct {
	Votes map[string]int
}


func main () {

	number_of_mandates := flag.Int("m", 0, "number of mandates")
	threshold := flag.Float64("t", 0, "Threshold for election")
	filename := flag.String("f", "", "filename with results")

	flag.Parse()

	if *filename != "" {
		file, err := os.ReadFile(*filename)
		if err != nil {
			fmt.Sprintln(err.Error())
		}
		json.Unmarshal(file, &js_results)
	}
	
	election_engine := new(methods.Result)
	election_engine.LoadVotes(js_results.Votes)
	election_engine.SetThreshold(*threshold)
	mandates := election_engine.GetAllMandates(*number_of_mandates)
	total_votes := election_engine.Total_votes
	printVotes(mandates, total_votes)
	printPage(election_engine)
}


func printVotes(parties []methods.Party, total_votes int) {


	sort.Slice(parties, func(i, j int) bool { return int(parties[i].Votes) > int(parties[j].Votes)} )

	fmt.Println("")

	for _, party := range parties {

		fmt.Printf("%-25s \t %7d   %5.2f %%   %3d \n", party.Name, party.Votes, party.Percentage, party.Seats)
	}
	fmt.Printf("total votes: %d\n", total_votes)
	fmt.Println("")
}


func printPage(results *methods.Result) {
	t, err := template.ParseFiles("./templates/results.html")
	if err != nil {
		fmt.Println(err.Error())
	}
	file, err := os.Create("./election.html")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer file.Close()

	err = t.Execute(file, results)
	if err != nil {
		fmt.Println(err.Error())
	}
}



	//reader := bufio.NewReader(os.Stdin)
	//number_of_parties := flag.Int("p", 0, "Number of parties")

/* 	else {
		for x := 0; x < *number_of_parties; x++ {
			fmt.Print("Enter party: ")
			party, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			party = strings.TrimSuffix(party, "\n")
			fmt.Print("Enter votes: ")
			party_votes, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			party_votes = strings.TrimSuffix(party_votes, "\n")
			party_votes_i, err := strconv.Atoi(party_votes)
			if err != nil {
				fmt.Println(err.Error())
				break
			}
			if party_votes_i < 0 {
				fmt.Println("Not a natural number")
				break
			}
			results.Votes[party] = party_votes_i
		}
	}
 */