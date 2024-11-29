package contracts

type Event struct {
	ID           string
	Title        string
	StartTime    string
	EndTime      string
	Description  string
	NotifyBefore string
	OwnerEmail   string
}
