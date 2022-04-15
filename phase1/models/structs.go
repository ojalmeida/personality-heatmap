package models

type FakeLocation struct {
	Coordinate `json:"location"`

	Accuracy float32 `json:"accuracy"`
}

type Coordinate struct {
	Lat float64 `json:"lat" yaml:"lat"`
	Lng float64 `json:"lng" yaml:"lng"`
}

type PhaseData struct {
	Nodes [3]*Node `yaml:"nodes"`

	City City `yaml:"city"`
}

type Node struct {
	Name     string       `yaml:"name"`
	Location FakeLocation `yaml:"location"`
	APIToken string       `yaml:"APIToken"`
}

type City struct {
	Name        string     `yaml:"name"`
	Coordinates Coordinate `yaml:"coordinates"`
}
