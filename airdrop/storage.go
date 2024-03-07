package airdrop

type Contributor struct {
	UserName     string
	TotalCommits uint64
}

type Storage interface {
	Store(repo string, contributors []Contributor) error
}
