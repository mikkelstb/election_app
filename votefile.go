package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type VoteData struct {
	districts []DistrictResult
}

type DistrictResult struct {
	name    string
	parties map[string]int
}

func readVoteFile(filename string) VoteData {

	var vote_data VoteData

	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	filescanner := bufio.NewScanner(file)

	//Read the first line and create the districts
	filescanner.Scan()
	for _, district := range strings.Split(filescanner.Text(), ";")[1:] {
		vote_data.districts = append(vote_data.districts, DistrictResult{name: district, parties: make(map[string]int)})
	}

	//Read the following lines of votes
	for filescanner.Scan() {

		party_data := strings.Split(filescanner.Text(), ";")
		for index, vote := range party_data[1:] {
			votes, err := strconv.Atoi(vote)
			if err != nil {
				panic(err)
			}
			vote_data.districts[index].parties[party_data[0]] = votes
		}
	}
	return vote_data
}
