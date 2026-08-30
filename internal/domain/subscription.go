package domain

type Subscription struct {
	Id          uint
	Description string
	Value       float64
	StartDate   string
	EndDate     string
	PersonId    *uint
	CardId      *uint
}
