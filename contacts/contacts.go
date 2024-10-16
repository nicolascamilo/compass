package contacts

import (
	"github.com/nicolascamilo/compass/domain"
	"github.com/resilva87/stringmetric"
	"regexp"
)

const (
	HIGH_ACCURACY   = "HIGH"
	MEDIUM_ACCURACY = "MEDIUM"
	LOW_ACCURACY    = "LOW"
	NO_MATCH        = "NONE"
)

func CompareContacts(source domain.Contact, potentialDuplicate domain.Contact) Accuracy {
	firstNameScore := stringmetric.RatcliffObershelpMetric(source.FirstName, potentialDuplicate.FirstName)
	lastNameScore := stringmetric.RatcliffObershelpMetric(source.LastName, potentialDuplicate.LastName)
	emailScore := stringmetric.RatcliffObershelpMetric(getEmailUsername(source.Email), getEmailUsername(potentialDuplicate.Email))
	// I'm assuming that zip code can't be pretty much similar? Might be if zip code 1 is 1234 and 2 is 123
	// So, it's either the same or not
	zipcodeScore := float64(0)
	if source.ZipCode == potentialDuplicate.ZipCode {
		zipcodeScore = 1
	}
	addressScore := stringmetric.RatcliffObershelpMetric(source.Address, potentialDuplicate.Address)

	return parseScore(firstNameScore + lastNameScore + emailScore + zipcodeScore + addressScore)
}

var emailRegex = regexp.MustCompile(`^[^@]+`)

// Emails might have different domains, but I only care about usernames
func getEmailUsername(email string) string {
	return emailRegex.FindString(email)
}

func parseScore(score float64) Accuracy {
	// Since the max value is 5,
	// Anything above 3.75 should be a High
	// I would assume that a score >= 2.50 would be Medium
	// Between 2.50 and 1.25, Low
	// Below 1.25, probably not a match?
	if score > 3.75 {
		return HIGH_ACCURACY
	} else if score < 3.75 && score >= 2.50 {
		return MEDIUM_ACCURACY
	} else if score < 2.50 && score > 1.25 {
		return LOW_ACCURACY
	}
	return NO_MATCH
}

func (a Accuracy) IsMatch() bool {
	// I would probably not do this? Just because is not clear but this is more simple than writing more ifs
	return a != NO_MATCH
}

type Accuracy string

type Match struct {
	IDSource  int64
	IDMatched int64
	Accuracy  Accuracy
}

func FindPossibleMatches(contacts []domain.Contact) []Match {
	matches := make([]Match, 0, len(contacts))

	// This is not optimal, but since there's a possibility of having more than one match per user then we might need to iterate
	// There might a solution with hashes though? But to be crystal clear, I don't think I have enough time to think about it
	for i := 0; i < len(contacts); i++ {
		for j := i + 1; j < len(contacts); j++ {
			contactI := contacts[i]
			contactJ := contacts[j]
			result := CompareContacts(contactI, contactJ)
			if result.IsMatch() {
				matches = append(matches, Match{
					IDSource:  contactI.ID,
					IDMatched: contactJ.ID,
					Accuracy:  result,
				},
				)
			}
		}
	}

	return matches
}
