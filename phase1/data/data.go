package data

import (
	"github.com/ojalmeida/personality-heatmap/phase1/models"
	"gopkg.in/yaml.v2"
	"os"
)

var Data = &models.PhaseData{
	Nodes: [3]*models.Node{

		{
			Name:     "",
			Location: models.FakeLocation{},
			APIToken: "",
		},

		{
			Name:     "",
			Location: models.FakeLocation{},
			APIToken: "",
		},

		{
			Name:     "",
			Location: models.FakeLocation{},
			APIToken: "",
		},
	},

	City: models.City{
		Name: "",
		Coordinates: models.Coordinate{
			Lat: 0,
			Lng: 0,
		},
	},
}

func Save(path string) error {

	dataOut, err := yaml.Marshal(Data)

	if err != nil {
		return err
	}

	err = os.WriteFile(path, dataOut, 0660)

	if err != nil {
		return err
	}

	return nil

}
