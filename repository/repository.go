package repository

import (
	"encoding/csv"
	"github.com/nicolascamilo/compass/domain"
	"log"
	"os"
	"strconv"
)

type repository struct {
	path string
}

type Repository interface {
	GetContacts() []domain.Contact
}

func New(path string) Repository {
	return &repository{path: path}
}

func (r *repository) GetContacts() []domain.Contact {
	file, err := os.Open(r.path)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	lines, err := reader.ReadAll()
	if err != nil {
		log.Panic(err)
	}

	contacts := make([]domain.Contact, 0, len(lines)-1)

	for _, line := range lines[1:] {
		parsedID, err := strconv.ParseInt(line[0], 10, 64)
		if err != nil {
			log.Panicf("file is bad formed, error: %s", err.Error())
		}
		contacts = append(contacts, domain.Contact{
			ID:        parsedID,
			FirstName: line[1],
			LastName:  line[2],
			Email:     line[3],
			ZipCode:   line[4],
			Address:   line[5],
		})
	}

	return contacts
}
