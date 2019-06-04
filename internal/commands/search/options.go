package search

type options struct {
	clientID     string
	clientSecret string

	forceSearchAndSave bool

	workersCount int
}

var opts options
