# Caching Server for TheCocktailDB

TheCocktailDB is a free API that has a large collection of cocktails. Unfortunately, it has some problems.

* It's a bit slow than what you would want
* The resources are not well structured
* The endpoints and their parameters are very confusing

This project tries to solve those issues by caching the entire API periodically using Redis. It does that by:

* Scheduling an Updater that fetches the API and stores its data on a Redis server periodically
* Exposing a few endpoints that are easier to understand to serve the data from the Redis server

# Running the project

1. Clone the repo

```
git clone https://github.com/Loukay/thecokctaildb-cache && cd thecokctaildb-cache
```

2. Create a .env file with the following variables:

```
REDIS_HOST=localhost
REDIS_PORT=6379
API_URL=https://url/api/json/v2/KEY
```

3. Install the packages and dependencies 

```console
go install
```

4. Run the project

```console
go run .
```
