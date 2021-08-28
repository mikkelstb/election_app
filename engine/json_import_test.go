package engine

import (
	"testing"
)


func TestElection_PrintDistrictResults(t *testing.T) {
	type fields struct {
		Allocated_seats int
		Districts       []District
		Parties         []Party
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "Basic test",
			fields: fields{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := ReadElectionFile("../resources/folketinget_danmark_2021_1.json")
			e.CalculateDistrictSeats()
			e.calculateAdditionalSeats()
			e.PrintDistrictResults()
		})
	}
}
