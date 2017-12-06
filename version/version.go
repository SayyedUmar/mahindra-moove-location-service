package version

import "fmt"

var Commit, Branch, State, TimeStamp string

type Version struct {
	Commit    string
	Branch    string
	State     string
	TimeStamp string
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
