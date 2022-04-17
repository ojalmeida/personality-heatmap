package database

import (
	"database/sql"
	"fmt"
	"personality-heatmap/data"
)

var Databases [3]*sql.DB

func CreateDatabases(dir string) {

	dbs := [3]*sql.DB{}

	for i := 0; i < 3; i++ {

		node := data.Data.Nodes[i]

		db, err := sql.Open("sqlite3", fmt.Sprintf("%s/%s.db", dir, node.Name))

		_, err = db.Exec(`CREATE TABLE IF NOT EXISTS profile (profile_id VARCHAR(50) NOT NULL PRIMARY KEY, name VARCHAR(50), birth_date VARCHAR(30), distance INT NOT NULL, interests VARCHAR(255), bio VARCHAR(1000))`)

		dbs[i] = db

		if err != nil {
			panic(err.Error())
		}

	}

	Databases = dbs

}
