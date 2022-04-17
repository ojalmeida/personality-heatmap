package data

import (
	"gopkg.in/yaml.v2"
	"os"
	"personality-heatmap/models"
)

var Data models.ProcessingData

func LoadData(path string) error {

	in, err := os.ReadFile(path)

	if err != nil {
		return err
	}

	err = yaml.Unmarshal(in, &Data)

	if err != nil {
		return err
	}

	return nil

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
