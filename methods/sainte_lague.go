package methods

type SainteLague struct {
	
}

func (this SainteLague) CalculateQuotients(votes int, seats_awarded int) float64 {
	return (float64(votes) / float64((seats_awarded + 1)))
}