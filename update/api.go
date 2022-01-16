package update

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
)

// IngredientName holds the name of an ingredient (required for JSON unmarshalling)
type IngredientName struct {
	Name string `json:"strIngredient1"`
}

// CocktailID holds the ID of a cocktail (required for JSON unmarshalling)
type CocktailID struct {
	ID string `json:"idDrink"`
}

type ingredientsNameResult struct {
	Ingredients []IngredientName `json:"drinks"`
}

type ingredientsResult struct {
	Ingredients []Ingredient `json:"ingredients"`
}

type cocktailIDResult struct {
	Cocktails []CocktailID `json:"drinks"`
}

type cocktailResult struct {
	Cocktails []Cocktail `json:"drinks"`
}

// Ingredient holds the details of an ingredient (name, type, image url, and whether it's alcohol or not)
type Ingredient struct {
	Name      string `json:"strIngredient"`
	Type      string `json:"strType"`
	IsAlcohol string `json:"strAlcohol"`
	Image     string
}

// Cocktail holds the details of a cocktail (name, category, IBA, glass type, ingredients, image url, and instructions)
type Cocktail struct {
	Name         string
	Category     string
	IBA          string
	Glass        string
	Instructions string
	Ingredients  string
	Measurements string
	Image        string
}

// API holds the Resty client (for now)
type API struct {
	HTTP *resty.Client
}

// GetAllIngredients method fetches all the ingredients from the API
/*
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
	var ingredientsList ingredientsNameResult
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
	var ingredientResult ingredientsResult

	resp, err := api.HTTP.R().Get("search.php?i=" + ingredientName.Name)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, errors.New("unexpected error in the API")
	}

	json.Unmarshal(resp.Body(), &ingredientResult)

	if len(ingredientResult.Ingredients) == 0 {
		return nil, errors.New("unexpected error in the API")
	}

	ingredient := ingredientResult.Ingredients[0]

	ingredient.Image = fmt.Sprintf("www.thecocktaildb.com/images/ingredients/%s-Medium.png", ingredient.Name)

	return &ingredient, nil
}

// GetAllCocktails method fetches all the cocktails from the API
/*
	API response example:

	{
		"drinks": [
			{
				"strDrink": "1-900-FUK-MEUP",
				"strDrinkThumb": "https://www.thecocktaildb.com/images/media/drink/uxywyw1468877224.jpg",
				"idDrink": "15395"
			},
			{
				"strDrink": "110 in the shade",
				"strDrinkThumb": "https://www.thecocktaildb.com/images/media/drink/xxyywq1454511117.jpg",
				"idDrink": "15423"
			},
			{
				"strDrink": "151 Florida Bushwacker",
				"strDrinkThumb": "https://www.thecocktaildb.com/images/media/drink/rvwrvv1468877323.jpg",
				"idDrink": "14588"
			},
		]
	}

	Returns an array of cocktail IDs ([]CocktailID)
*/
func (api API) GetAllCocktails() ([]CocktailID, error) {
	var cocktailsList cocktailIDResult
	resp, err := api.HTTP.R().Get("filter.php?a=Alcoholic")

	if err != nil {
		return nil, err
	}

	json.Unmarshal(resp.Body(), &cocktailsList)

	return cocktailsList.Cocktails, nil
}

// GetCocktailDetails method fetches the details of a cocktail
/*
	API response example:

	{
		"drinks": [
			{
			"idDrink": "11007",
			"strDrink": "Margarita",
			"strDrinkAlternate": null,
			"strTags": "IBA,ContemporaryClassic",
			"strVideo": null,
			"strCategory": "Ordinary Drink",
			"strIBA": "Contemporary Classics",
			"strAlcoholic": "Alcoholic",
			"strGlass": "Cocktail glass",
			"strInstructions": "Rub the rim of the glass with the lime slice to make the salt stick to it. Take care to moisten only the outer rim and sprinkle the salt on it. The salt should present to the lips of the imbiber and never mix into the cocktail. Shake the other ingredients with ice, then carefully pour into the glass.",
			"strIngredient1": "Tequila",
			"strIngredient2": "Triple sec",
			"strIngredient3": "Lime juice",
			"strIngredient4": "Salt",
			"strMeasure1": "1 1/2 oz ",
			"strMeasure2": "1/2 oz ",
			"strMeasure3": "1 oz ",
			}
		]
	}
	Returns the cocktail details (Cocktail struct)
*/
func (api API) GetCocktailDetails(cocktailID CocktailID) (*Cocktail, error) {

	resp, err := api.HTTP.R().Get("lookup.php?i=" + cocktailID.ID)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, errors.New("unexpected error in the API")
	}

	var result map[string]interface{}

	json.Unmarshal(resp.Body(), &result)

	if result["drinks"] == nil {
		return nil, errors.New("unexpected error in the API")
	}

	details := result["drinks"].([]interface{})[0].(map[string]interface{})

	var cocktail Cocktail

	var ingredients []string
	var measurements []string

	for key, value := range details {

		if value == nil {
			continue
		}

		val := value.(string)

		if strings.Contains(key, "strIngredient") {
			ingredients = append(ingredients, val)
			continue
		}

		if strings.Contains(key, "strMeasure") {
			measurements = append(measurements, val)
			continue
		}

		switch key {
		case "strDrink":
			cocktail.Name = val
		case "strCategory":
			cocktail.Category = val
		case "strIBA":
			cocktail.IBA = val
		case "strGlass":
			cocktail.Glass = val
		case "strInstructions":
			cocktail.Instructions = val
		case "strDrinkThumb":
			cocktail.Image = val
		}
	}

	cocktail.Ingredients = strings.Join(ingredients, ",")
	cocktail.Measurements = strings.Join(measurements, ",")

	return &cocktail, nil
}

// NewAPIClient method creates a new API clients
func NewAPIClient(baseURL string) *API {
	HTTP := resty.New()
	HTTP.SetHostURL(baseURL)
	return &API{HTTP}
}
