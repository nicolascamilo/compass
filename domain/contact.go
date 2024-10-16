package domain

// I did this to avoid stuttering (contacts.Contact)
type Contact struct {
	ID        int64
	FirstName string
	LastName  string
	Email     string
	ZipCode   string
	Address   string
}
