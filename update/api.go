package update

import (
	"encoding/json"
	"errors"

	"github.com/go-resty/resty/v2"
)

// IngredientName holds the name of an ingredient (required for JSON unmarshalling)
type IngredientName struct {
	Name string `json:"strIngredient1"`
}

// Ingredient holds the details of an ingredient (name, type, image url, and whether it's alcohol or not)
type Ingredient struct {
	Name      string `json:"strIngredient"`
	Type      string `json:"strType"`
	IsAlcohol string `json:"strAlcohol"`
	Image     string
}

// IngredientResult holds a list of ingredients ([]Ingredient) (required for JSON unmarshalling)
type IngredientResult struct {
	Ingredients []Ingredient `json:"ingredients"`
}

// Result holds a lit of ingredient names ([]IngredientName) (required for JSON unmarshalling)
type Result struct {
	Ingredients []IngredientName `json:"drinks"`
}

// API holds the Resty client (for now)
type API struct {
	HTTP *resty.Client
}

// GetAllIngredients method fetches all the ingredients from the API
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
	resp, err := api.HTTP.R().Get("list.php?i=list")

	if err != nil {
		return nil, err
	}

	json.Unmarshal(resp.Body(), &ingredientsList)

	return ingredientsList.Ingredients, nil
}

// GetIngredientDetails method fetches the details of an ingredient
/*
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

	resp, err := api.HTTP.R().Get("search.php?i=" + ingredientName.Name)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, errors.New("unexpected error in the API")
	}

	json.Unmarshal(resp.Body(), &ingredientResult)

	return &ingredientResult.Ingredients[0], nil
}

// NewAPIClient method creates a new API clients
func NewAPIClient(baseURL string) *API {
	HTTP := resty.New()
	HTTP.SetHostURL(baseURL)
	return &API{HTTP}
}
