package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/redis/go-redis/v9"
)

// Global context is required for all redis operations
var ctx = context.Background()
var redisAddress = os.Getenv("REDIS_ADDRESS")
var rdb *redis.Client

func buildRedis() {
	if redisAddress == "" {
		redisAddress = "localhost:6379"
	}
	rdb = redis.NewClient(&redis.Options{
		Addr:     redisAddress, // Redis address
		Password: "",           // No password set
		DB:       0,            // Use default DB
	})

	fmt.Printf("Connecting to Redis: PING\n")
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	fmt.Printf("Connected to Redis: %s\n", pong)
}

func main() {

	buildRedis()

	// Matches paths starting with /order/ (e.g. /order/123)
	http.HandleFunc("/order/", getByName)
	// Matches exactly /order
	http.HandleFunc("/order/keys", getAll)
	log.Println("Listening on port 3001 :)")
	log.Fatal(http.ListenAndServe(":3001", nil))
}

func getByName(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	// Extract the name parameter from the URL
	name := strings.TrimPrefix(r.URL.Path, "/order/")

	if name == "" {
		http.Error(w, "Name parameter is missing", http.StatusBadRequest)
		return
	}

	val, err := rdb.Get(ctx, name).Result()
	if err == redis.Nil {
		http.Error(w, "Key not found", http.StatusNotFound)
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		fmt.Fprint(w, val)
	}
}

func getAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	keys, err := rdb.Keys(ctx, "*").Result()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Keys: %v", keys)
}
