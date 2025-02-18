package models

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type POIEntity struct {
	DB *sql.DB
}

type POI struct {
	Name    string
	Tourism sql.NullString
	Natural sql.NullString
}

func (p *POIEntity) Get(lat float64, lon float64, radius int) ([]POI, error) {
	qry := `
		SELECT 
			name, 
			tourism, 
			"natural" 
		FROM tourism_poi
		WHERE name IS NOT NULL 
		AND ST_DWithin(
			ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography, 
			geom::geography, 
			$3
		)
		ORDER BY 
			ST_Distance(
				ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography, 
				geom::geography
			) ASC 
		LIMIT 10;
	`

	rows, err := p.DB.Query(qry, lon, lat, radius)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pois []POI

	for rows.Next() {
		var poi POI
		err := rows.Scan(&poi.Name, &poi.Tourism, &poi.Natural)
		if err != nil {
			return nil, err
		}
		pois = append(pois, poi)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return pois, nil
}
