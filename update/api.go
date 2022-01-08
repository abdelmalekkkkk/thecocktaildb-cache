package update

import (
	"encoding/json"

	"github.com/go-resty/resty/v2"
)

type Ingredient struct {
	Name string `json:"strIngredient1"`
}

type Result struct {
	Ingredients []Ingredient `json:"drinks"`
}

type API struct {
	Http *resty.Client
}

/*
	Get a list of all available ingredients

	API response example:

	{
		"drinks": [
			{
			"strIngredient1": "Light rum"
			},
			{
			"strIngredient1": "Applejack"
			},
			{
			"strIngredient1": "Gin"
			},
			{
			"strIngredient1": "Dark rum"
			},
			{
			"strIngredient1": "Sweet Vermouth"
			},
		]
	}

	Returns an array of Ingredient
*/
func (api API) GetAllIngredients() ([]Ingredient, error) {
	var ingredientsList Result

	resp, err := api.Http.R().Get("https://www.thecocktaildb.com/api/json/v1/1/list.php?i=list")

	if err != nil {
		return nil, err
	}

	json.Unmarshal(resp.Body(), &ingredientsList)
	
	return ingredientsList.Ingredients, nil
}

/*

	Creates a new API client

*/
func NewAPIClient() *API {
	Http := resty.New()
	return &API{Http}
}