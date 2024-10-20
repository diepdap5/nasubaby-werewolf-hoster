package repository

import (
	"database/sql"
	"fmt"
	"nasubaby-werewolf-hoster/cmd/model"
)

func GetGameBases(db *sql.DB) ([]*model.GameBase, error) {
	// Prepare the query
	query := "SELECT id, name, description FROM game_base"

	// Execute the query
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	// Check if there are no rows returned
	if !rows.Next() {
		return nil, fmt.Errorf("no data found")
	}

	defer rows.Close()

	// Iterate over the rows and scan the results into a slice of GameBase structs
	var gameBases []*model.GameBase
	for rows.Next() {
		gb := &model.GameBase{}
		if err := rows.Scan(&gb.ID, &gb.RoleCount, &gb.RolesList); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		gameBases = append(gameBases, gb)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	return gameBases, nil
}
