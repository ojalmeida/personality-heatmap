package geo

import (
	"personality-heatmap/models"
)

func GetNodesCoordinates(targetCityCoordinates models.Coordinate) [3]models.Coordinate {

	// Builds a triangle with vertices apart 10km from center, with central angles of 120 degrees
	return [3]models.Coordinate{

		// Rounding 10km to a delta of 0.1 degrees, in South America
		{
			Lat: targetCityCoordinates.Lat - 0.1,
			Lng: targetCityCoordinates.Lng,
		},

		{

			Lat: targetCityCoordinates.Lat + 0.05,
			Lng: targetCityCoordinates.Lng - 0.0866,
		},

		{

			Lat: targetCityCoordinates.Lat + 0.05,
			Lng: targetCityCoordinates.Lng + 0.0866,
		},
	}

}
