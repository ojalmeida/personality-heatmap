package database

import (
	"database/sql"
	"fmt"
	"personality-heatmap/models"
	"strings"
)

func Insert(profile models.Profile, db *sql.DB) error {

	_, err := db.Exec(fmt.Sprintf(`INSERT INTO profile (profile_id, name, birth_date, distance, interests, bio) VALUES ("%s", "%s", "%s", "%d", "%s", "%s")`,

		profile.User.ID,
		profile.User.Name,
		profile.User.BirthDate,
		profile.Distance,
		strings.Join(profile.Interests, ","),
		profile.User.Bio,
	))

	return err

}
