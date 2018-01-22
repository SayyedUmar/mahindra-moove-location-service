package version

import "fmt"

var Commit, Branch, State, TimeStamp string

type Version struct {
	Commit    string `json:"commit"`
	Branch    string `json:"branch"`
	State     string `json:"state"`
	TimeStamp string `json:"timestamp"`
}

func PrintVersion() {
	fmt.Printf("Commit : %s\n", Commit)
	fmt.Printf("Branch: %s\n", Branch)
	fmt.Printf("State : %s\n", State)
	fmt.Printf("TimeStamp: %s\n", TimeStamp)
}

func GetVersion() Version {
	return Version{Commit: Commit, Branch: Branch, State: State, TimeStamp: TimeStamp}
}
