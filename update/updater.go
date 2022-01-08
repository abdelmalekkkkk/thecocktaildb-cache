package update

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type Updater struct {
	Ctx *context.Context
	Redis *redis.Client
	*API
}



/*
	Fetches data from the API and stores it in the Redis Server
*/
func (updater Updater) Run()  {
	
	log.Printf("Started database update at %s\n", time.Now())

	ingredients, err := updater.API.GetAllIngredients()

	if err != nil {
		log.Println(err)
		log.Println("Couldn't retrieve the ingredients. Aborting the update.")
	}

	for _, v := range ingredients {
		fmt.Printf(v.Name)
	}


}