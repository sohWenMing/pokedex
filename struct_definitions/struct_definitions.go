package structdefinitions

type LocationAreaResult struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type Result struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
