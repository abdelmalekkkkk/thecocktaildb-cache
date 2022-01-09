package update

import (
	"context"
	"log"
	"sync"

	"github.com/go-redis/redis/v8"
)

type Updater struct {
	Ctx   *context.Context
	Redis *redis.Client
	*API
	RequestsChannel chan IngredientName
	ResultsChannel  chan Ingredient
}

/*
	Fetches data from the API and stores it in the Redis Server
*/
func (updater Updater) Run() {

	log.Printf("Started database update")

	go updater.UpdateIngredients()

}

func (updater Updater) UpdateIngredients() {

	log.Printf("Setting ingredientsID to 0")
	updater.Redis.Set(*updater.Ctx, "ingredientsID", 0, 0)
	updater.Redis.Set(*updater.Ctx, "alcoholsID", 0, 0)

	log.Printf("Fetching ingredients from API")
	ingredients, err := updater.API.GetAllIngredients()

	if err != nil {
		log.Println(err)
		log.Println("Couldn't retrieve the ingredients. Aborting the update.")
	}

	log.Printf("%d ingredients retrieved", len(ingredients))

	updater.RequestsChannel = make(chan IngredientName, 10)

	var wg sync.WaitGroup

	go func() {
		for _, ingredient := range ingredients {
			updater.RequestsChannel <- ingredient
		}
	}()

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go updater.GetIngredientWorker(&wg)
		go updater.AddIngredientWorker(&wg)
	}

}

func (updater Updater) GetIngredientWorker(wg *sync.WaitGroup) {
	for ingredient := range updater.RequestsChannel {
		ingredientDetails, err := updater.API.GetIngredientDetails(ingredient)
		if err != nil {
			continue
		}
		updater.ResultsChannel <- *ingredientDetails
		log.Printf("%v", ingredientDetails)
	}
}

func (updater Updater) AddIngredientWorker(wg *sync.WaitGroup) {
	for ingredient := range updater.ResultsChannel {
		updater.AddIngredient(ingredient)
	}
}

func (updater Updater) AddIngredient(ingredient Ingredient) {
	client := updater.Redis
	if ingredient.IsAlcohol == "Yes" {
		id := client.Incr(*updater.Ctx, "alcoholsID")
		log.Printf("Adding drink to id %d", id)
	}
}
