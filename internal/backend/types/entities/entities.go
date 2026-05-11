package entities

type World struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	KeyFacts    []string `json:"key_facts"`
	References  []string `json:"references"`
}

type Scenario struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	KeyDetails  []string `json:"key_details"`
	References  []string `json:"references"`
}

type Character struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Personality   string   `json:"personality"`
	Backstory     string   `json:"backstory"`
	Goals         []string `json:"goals"`
	Relationships []string `json:"relationships"`
	References    []string `json:"references"`
}
