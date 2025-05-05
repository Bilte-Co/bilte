package domain

type Projects map[string][]Project

type Project struct {
	Company      string   `json:"company"`
	ClientName   string   `json:"client_name"`
	Technologies []string `json:"technologies"`
	Description  string   `json:"description"`
}

type Link struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	External bool   `json:"external"`
}

type Experience struct {
	Name        string   `json:"name"`
	Slug        string   `json:"slug"`
	Role        string   `json:"role"`
	Description string   `json:"description"`
	Start       string   `json:"start"`
	End         string   `json:"end"`
	Bullets     []string `json:"bullets"`
}

type OtherExperience struct {
	Name  string `json:"name"`
	Role  string `json:"role"`
	Start string `json:"start"`
	End   string `json:"end"`
}

type Education struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Degree      string `json:"degree"`
	Graduation  string `json:"graduation"`
}

type Competancy struct {
	Name  string   `json:"name"`
	Items []string `json:"items"`
}

type Resume struct {
	Name            string            `json:"name"`
	Title           string            `json:"title"`
	Email           string            `json:"email"`
	Phone           string            `json:"phone"`
	Location        string            `json:"location"`
	Links           []Link            `json:"links"`
	Intro           string            `json:"intro"`
	Competancies    []Competancy      `json:"competancies"`
	Experience      []Experience      `json:"experience"`
	OtherExperience []OtherExperience `json:"other_experience"`
	Education       []Education       `jsone:"education"`
	Interests       []string          `json:"interests"`
}
