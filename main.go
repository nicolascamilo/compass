package main

import (
	"github.com/nicolascamilo/compass/contacts"
	"github.com/nicolascamilo/compass/repository"
)

func main() {
	repo := repository.New("./test/Sample.csv")
	contacts.FindPossibleMatches(repo.GetContacts())
}
