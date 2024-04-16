package service

import (
	"database/sql"
	"log"
)

// Recipe is a lightweight recipe type w/o ingredients
type Recipe struct {
	Name string   `json:"name"`
	ID   int      `json:"id"`
	Tags []string `json:"tags"`
}

// GetAllRecipes returns all recipes in the recipe table
func GetAllRecipes(db *sql.DB, userID string) ([]Recipe, error) {
	accountID, err := GetAccountID(db, userID)

	if err != nil {
		log.Println("Error getting account ID")
		return nil, err
	}

	recipesQuery := `
		SELECT recipe.id, name, tag_name FROM recipe
			LEFT JOIN recipe_tag on recipe.id = recipe_tag.recipe_id
			WHERE account_id = ?
			ORDER BY lower(recipe.name);
	`
	results, err := db.Query(recipesQuery, accountID)

	if err != nil {
		log.Println("Error querying recipes")
		return nil, err
	}

	recipes := []Recipe{}

	for results.Next() {
		r := Recipe{Tags: []string{}}
		var tag sql.NullString
		err = results.Scan(&r.ID, &r.Name, &tag)
		if err != nil {
			return nil, err
		}

		if len(recipes) > 0 {
			lastRecipe := recipes[len(recipes)-1]
			log.Println(r.ID, lastRecipe.ID)
			if r.ID == lastRecipe.ID {
				lastRecipe.Tags = append(lastRecipe.Tags, tag.String)
				log.Println(lastRecipe.Tags)
				continue
			}
		}

		if tag.Valid {
			r.Tags = []string{tag.String}
		}
		recipes = append(recipes, r)
	}
	return recipes, nil
}
