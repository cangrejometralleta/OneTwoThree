package main

// PersonBody is what a Client Sends.
// Weak on Purpose: the Wire Cannot be Trusted,
// so nothing here is a Business Type yet.
type PersonBody struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

// PersonView is what a Client Gets back.
type PersonView struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

// BuildPersonRecord Promotes untrusted Input into Business Types.
func (b PersonBody) BuildPersonRecord(id PersonID) Person {
	return Person{
		ID:    id,
		Name:  FullName(b.Name),
		Email: EmailAddress(b.Email),
		Phone: PhoneNumber(b.Phone),
	}
}

// RenderPersonView Demotes a Business Person back to the Wire.
func RenderPersonView(p Person) PersonView {
	return PersonView{
		ID:    uint(p.ID),
		Name:  string(p.Name),
		Email: string(p.Email),
		Phone: string(p.Phone),
	}
}

// RenderPeopleViews Repeats the Move for a whole Roster.
func RenderPeopleViews(people []Person) []PersonView {
	views := make([]PersonView, 0, len(people))
	for _, person := range people {
		views = append(views, RenderPersonView(person))
	}
	return views
}
