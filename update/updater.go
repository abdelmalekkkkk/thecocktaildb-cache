package update

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
	RequestsChannel chan IngredientName
	ResultsChannel  chan Ingredient
}


// Run method fetches data from the API and stores it in the Redis Server
func (updater Updater) Run() {

	log.Printf("Started database update")

	go updater.updateIngredients()

}

func (updater Updater) updateIngredients() {

	log.Printf("Setting ingredientsID to 0")
	updater.Redis.Set(*updater.Ctx, "ingredientID", 0, 0)
	updater.Redis.Set(*updater.Ctx, "alcoholID", 0, 0)

	log.Printf("Fetching ingredients from API")
	ingredients, err := updater.API.GetAllIngredients()

	if err != nil {
		log.Println(err)
		log.Println("Couldn't retrieve the ingredients. Aborting the update.")
	}

	log.Printf("%d ingredients retrieved", len(ingredients))

	updater.RequestsChannel = make(chan IngredientName, 10)
	updater.ResultsChannel = make(chan Ingredient, 10)

	var wg sync.WaitGroup

	go func() {
		for _, ingredient := range ingredients {
			updater.RequestsChannel <- ingredient
		}
		close(updater.RequestsChannel)
	}()

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go updater.getIngredientWorker(&wg)
		go updater.addIngredientWorker(&wg)
	}

	wg.Wait()

}

func (updater Updater) getIngredientWorker(wg *sync.WaitGroup) {
	for ingredient := range updater.RequestsChannel {
		ingredientDetails, err := updater.API.GetIngredientDetails(ingredient)
		if err != nil {
			log.Printf("Error with %v", ingredient)
			continue
		}
		log.Printf("%v", ingredientDetails)
		go func() {
			updater.ResultsChannel <- *ingredientDetails
		}()
	}
}

func (updater Updater) addIngredientWorker(wg *sync.WaitGroup) {
	for ingredient := range updater.ResultsChannel {
		go updater.addIngredient(ingredient)
	}
}

func (updater Updater) addIngredient(ingredient Ingredient) {
	client := updater.Redis
	entry := make(map[string]interface{})
	entry["name"] = ingredient.Name
	var id int64
	var key string
	if ingredient.IsAlcohol == "Yes" {
		key = "alcohol"
	} else {
		key = "ingredient"
		entry["type"] = ingredient.Type
	}
	id = client.Incr(*updater.Ctx, key+"ID").Val()
	log.Printf("ID %d, Key %s,The entry is %+v", id, key, entry)
	client.HMSet(*updater.Ctx, key+":"+strconv.FormatInt(id, 10), entry)
}
