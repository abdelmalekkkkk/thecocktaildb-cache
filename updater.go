package main

import (
	"context"
	"log"
	"strconv"
	"sync"

	"github.com/go-redis/redis/v8"
)

// Updater handles the updating of the Redis database
type Updater struct {
	Ctx   *context.Context
	Redis *redis.Client
	*API
}

// Run method fetches data from the API and stores it in the Redis Server
func (updater Updater) Run() {

	log.Printf("Started database update")

	go updater.updateIngredients()
	go updater.updateCocktails()

}

func (updater Updater) updateIngredients() {

	log.Print("Setting ingredientID and alcoholID to 0")
	updater.Redis.Set(*updater.Ctx, "ingredientID", 0, 0)
	updater.Redis.Set(*updater.Ctx, "alcoholID", 0, 0)

	log.Printf("Fetching ingredients from API")
	ingredients, err := updater.API.GetAllIngredients()

	if err != nil {
		log.Println(err)
		log.Println("Couldn't retrieve the ingredients. Aborting the update.")
	}

	log.Printf("%d ingredients retrieved", len(ingredients))

	requestsChannel := make(chan IngredientName, 10)
	resultsChannel := make(chan Ingredient, 10)

	var wg sync.WaitGroup

	go func() {
		for _, ingredient := range ingredients {
			requestsChannel <- ingredient
		}
		close(requestsChannel)
	}()

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go updater.getIngredientWorker(&wg, requestsChannel, resultsChannel)
		go updater.addIngredientWorker(&wg, resultsChannel)
	}

	wg.Wait()

}

func (updater Updater) updateCocktails() {

	log.Print("Setting cocktailID to 0")
	updater.Redis.Set(*updater.Ctx, "cocktailID", 0, 0)

	log.Printf("Fetching cocktails from API")
	cocktails, err := updater.API.GetAllCocktails()

	if err != nil {
		log.Println(err)
		log.Println("Couldn't retrieve the cocktails. Aborting the update.")
	}

	log.Printf("%d alcohols retrieved", len(cocktails))

	requestsChannel := make(chan CocktailID, 10)
	resultsChannel := make(chan Cocktail, 10)

	var wg sync.WaitGroup

	go func() {
		for _, cocktail := range cocktails {
			requestsChannel <- cocktail
		}
		close(requestsChannel)
	}()

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go updater.getCocktailWorker(&wg, requestsChannel, resultsChannel)
		go updater.addCocktailWorker(&wg, resultsChannel)
	}

	wg.Wait()

}

func (updater Updater) getIngredientWorker(wg *sync.WaitGroup, requestsChannel <-chan IngredientName, resultsChannel chan<- Ingredient) {
	for ingredient := range requestsChannel {
		ingredientDetails, err := updater.API.GetIngredientDetails(ingredient)
		if err != nil {
			log.Printf("Error with %v", ingredient)
			continue
		}
		resultsChannel <- *ingredientDetails
	}
}

func (updater Updater) addIngredientWorker(wg *sync.WaitGroup, resultsChannel <-chan Ingredient) {
	for ingredient := range resultsChannel {
		go updater.addIngredient(ingredient)
	}
}

func (updater Updater) addIngredient(ingredient Ingredient) {
	client := updater.Redis

	entry := map[string]string{
		"id":    ingredient.ID,
		"name":  ingredient.Name,
		"image": "https://" + ingredient.Image,
		"type":  ingredient.Type,
	}

	key := "alcohol"

	if ingredient.IsAlcohol != "Yes" {
		key = "ingredient"
	}

	id := client.Incr(*updater.Ctx, key+"ID").Val()

	client.HMSet(*updater.Ctx, key+":"+strconv.FormatInt(id, 10), entry)
}

func (updater Updater) getCocktailWorker(wg *sync.WaitGroup, requestsChannel <-chan CocktailID, resultsChannel chan<- Cocktail) {
	for cocktail := range requestsChannel {
		cocktailDetails, err := updater.API.GetCocktailDetails(cocktail)
		if err != nil {
			log.Printf("Error with %v", cocktail)
			continue
		}
		resultsChannel <- *cocktailDetails
	}
}

func (updater Updater) addCocktailWorker(wg *sync.WaitGroup, resultsChannel <-chan Cocktail) {
	for cocktail := range resultsChannel {
		go updater.addCocktail(cocktail)
	}
}

func (updater Updater) addCocktail(cocktail Cocktail) {
	client := updater.Redis

	entry := map[string]string{
		"id":           cocktail.ID,
		"name":         cocktail.Name,
		"category":     cocktail.Category,
		"iba":          cocktail.IBA,
		"glass":        cocktail.Glass,
		"instructions": cocktail.Instructions,
		"ingredients":  cocktail.Ingredients,
		"measurements": cocktail.Measurements,
		"image":        cocktail.Image,
	}

	id := client.Incr(*updater.Ctx, "cocktailID").Val()

	client.HMSet(*updater.Ctx, "cocktail:"+strconv.FormatInt(id, 10), entry)
}
