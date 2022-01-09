package update

import (
	"encoding/json"

	"github.com/go-resty/resty/v2"
)

type IngredientName struct {
	Name string `json:"strIngredient1"`
}

type Ingredient struct {
	Name      string `json:"strIngredient"`
	Type      string `json:"strType"`
	IsAlcohol string `json:"strAlcohol"`
}

type IngredientResult struct {
	Ingredients []Ingredient `json:"ingredients"`
}

type Result struct {
	Ingredients []IngredientName `json:"drinks"`
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
func (api API) GetAllIngredients() ([]IngredientName, error) {
	var ingredientsList Result

	resp, err := api.Http.R().Get("https://www.thecocktaildb.com/api/json/v1/1/list.php?i=list")

	if err != nil {
		return nil, err
	}

	json.Unmarshal(resp.Body(), &ingredientsList)

	return ingredientsList.Ingredients, nil
}

/*
	Get the the details of an ingredient

	API response example:

	{
	"ingredients": [
		{
		"idIngredient": "305",
		"strIngredient": "Light Rum",
		"strDescription": "Light rums, also referred to as \"silver\" or \"white\" rums, in general, have very little flavor aside from a general sweetness. Light rums are sometimes filtered after aging to remove any colour. The majority of light rums come from Puerto Rico. Their milder flavors make them popular for use in mixed drinks, as opposed to drinking them straight. Light rums are included in some of the most popular cocktails including the Mojito and the Daiquiri.",
		"strType": "Rum",
		"strAlcohol": "Yes",
		"strABV": null
		}
	]
	}

	Returns a struct of type Ingredient
*/
func (api API) GetIngredientDetails(ingredientName IngredientName) (*Ingredient, error) {
	var ingredientResult IngredientResult

	resp, err := api.Http.R().Get("https://www.thecocktaildb.com/api/json/v1/1/search.php?i=" + ingredientName.Name)

	if err != nil {
		return nil, err
	}

	json.Unmarshal(resp.Body(), &ingredientResult)

	return &ingredientResult.Ingredients[0], nil
}

/*

	Creates a new API client

*/
func NewAPIClient() *API {
	Http := resty.New()
	return &API{Http}
}
