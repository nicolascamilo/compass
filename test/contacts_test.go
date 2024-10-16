package test

import (
	"github.com/nicolascamilo/compass/contacts"
	"github.com/nicolascamilo/compass/domain"
	"github.com/nicolascamilo/compass/repository"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFindPossibleMatches(t *testing.T) {
	testCases := []struct {
		name            string
		contacts        []domain.Contact
		expectedMatches []contacts.Match
	}{
		{
			name: "should be a high value because contacts are identical",
			contacts: []domain.Contact{
				{
					ID:        1,
					FirstName: "John",
					LastName:  "Doe",
					Email:     "johndoe@example.com",
					ZipCode:   "123",
					Address:   "Fake Street 123",
				},
				{
					ID:        2,
					FirstName: "John",
					LastName:  "Doe",
					Email:     "johndoe@test.com",
					ZipCode:   "123",
					Address:   "Fake Street 123",
				},
			},
			expectedMatches: []contacts.Match{
				{
					IDSource:  1,
					IDMatched: 2,
					Accuracy:  contacts.HIGH_ACCURACY,
				},
			},
		},
		{
			name: "should be a high value but similar values",
			contacts: []domain.Contact{
				{
					ID:        1,
					FirstName: "John",
					LastName:  "Doe",
					Email:     "johndoe@example.com",
					ZipCode:   "123",
					Address:   "Fake Street 123",
				},
				{
					ID:        2,
					FirstName: "Jones",
					LastName:  "Doe",
					Email:     "jonesdoe@example.com",
					ZipCode:   "123",
					Address:   "Fake Street 321",
				},
			},
			expectedMatches: []contacts.Match{
				{
					IDSource:  1,
					IDMatched: 2,
					Accuracy:  contacts.HIGH_ACCURACY,
				},
			},
		},
		{
			name: "medium accuracy",
			contacts: []domain.Contact{
				{
					ID:        1,
					FirstName: "John",
					LastName:  "Doe",
					Email:     "johndoe@example.com",
					ZipCode:   "123",
					Address:   "Fake Street 123",
				},
				{
					ID:        2,
					FirstName: "Johns",
					LastName:  "Does",
					Email:     "johns_does_fake_st@example.com",
					ZipCode:   "1234",
					Address:   "Fake St 13",
				},
			},
			expectedMatches: []contacts.Match{
				{
					IDSource:  1,
					IDMatched: 2,
					Accuracy:  contacts.MEDIUM_ACCURACY,
				},
			},
		},
		{
			name: "low accuracy",
			contacts: []domain.Contact{
				{
					ID:        1,
					FirstName: "John",
					LastName:  "Doe",
					Email:     "johndoe@example.com",
					ZipCode:   "123",
					Address:   "Fake Street 123",
				},
				{
					ID:        2,
					FirstName: "Jahn",
					LastName:  "Aldoe",
					Email:     "jahn_aldoe@example.com",
					ZipCode:   "1234",
					Address:   "Not Real 666",
				},
			},
			expectedMatches: []contacts.Match{
				{
					IDSource:  1,
					IDMatched: 2,
					Accuracy:  contacts.LOW_ACCURACY,
				},
			},
		},
		{
			name: "not a match",
			contacts: []domain.Contact{
				{
					ID:        1,
					FirstName: "John",
					LastName:  "Doe",
					Email:     "johndoe@example.com",
					ZipCode:   "123",
					Address:   "Fake Street 123",
				},
				{
					ID:        2,
					FirstName: "Jane",
					LastName:  "Clark",
					Email:     "jane_clark@example.com",
					ZipCode:   "1234",
					Address:   "Not Real 666",
				},
			},
			expectedMatches: []contacts.Match{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := contacts.FindPossibleMatches(testCase.contacts)
			require.Equal(t, testCase.expectedMatches, result)
		})
	}
}

/*
First run with *Contact in the compare func:
BenchmarkFindPossibleMatches-10    	       1	4413284250 ns/op	8582640368 B/op	93305499 allocs/op ----> 8 GB PER RUN, that's a lot

Then I realized that I could compile once the regex, this was the result:
BenchmarkFindPossibleMatches-10    	       1	3406769125 ns/op	7085265088 B/op	77315408 allocs/op ----> Still a lot but that's almost 1.5 GB less per run
*/
func BenchmarkFindPossibleMatches(b *testing.B) {
	file := "Sample.csv"
	repo := repository.New(file)
	contactList := repo.GetContacts()

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		b.StartTimer()
		contacts.FindPossibleMatches(contactList)
		b.StopTimer()
	}
}
