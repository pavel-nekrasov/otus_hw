package contracts

type Event struct {
	ID           string
	Title        string
	StartTime    int64
	EndTime      int64
	Description  string
	NotifyBefore string
	OwnerEmail   string
}
