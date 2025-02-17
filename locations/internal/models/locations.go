package models

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type LocationEntity struct {
	DB *sql.DB
}

type Location struct {
	Name       string
	Place      string
	Population sql.NullInt64
}

func (l *LocationEntity) Get(lat float64, lon float64, radius int) ([]Location, error) {
	qry := `
		SELECT 
			name, 
			place, 
			population 
		FROM planet_osm_point 
		WHERE place IS NOT NULL 
		AND name IS NOT NULL 
		AND ST_DWithin(
			ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography, 
			ST_Transform(way, 4326)::geography, 
			$3
		) 
		ORDER BY 
			ST_Distance(
				ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography, 
				ST_Transform(way, 4326)::geography
			) ASC,
			CASE 
				WHEN place = 'town' THEN 1
				WHEN place = 'village' THEN 2
				WHEN place = 'suburb' THEN 3
				ELSE 4
			END,
			population DESC 
		LIMIT 10;
	`

	rows, err := l.DB.Query(qry, lon, lat, radius)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []Location

	for rows.Next() {
		var loc Location
		err := rows.Scan(&loc.Name, &loc.Place, &loc.Population)
		if err != nil {
			return nil, err
		}
		locations = append(locations, loc)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return locations, nil
}
