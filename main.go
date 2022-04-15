package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	FirstName string `json:"first-name"`
	LastName  string `json:"last-name"`
}

var movies []Movie

func main() {
	//Define Data
	movies = append(movies, Movie{ID: "1", Isbn: "4323455", Title: "Movie One", Director: &Director{FirstName: "John", LastName: "Doe"}})
	movies = append(movies, Movie{ID: "2", Isbn: "55432344", Title: "Movie Two", Director: &Director{FirstName: "Steve", LastName: "Mcmaneman"}})

	r := mux.NewRouter()
	//Routes
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovies).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Printf("Start the server on port 8000")
	if err := http.ListenAndServe(":8000", r); err != nil {
		log.Fatal(err)
	}
}

//Delete a single movie
func deleteMovie(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)

	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
	//Return rest of movies
	json.NewEncoder(writer).Encode(movies)
}

//Update a movie
func updateMovie(writer http.ResponseWriter, request *http.Request) {
	//set json content type
	writer.Header().Set("Content-Type", "application/json")
	//params
	params := mux.Vars(request)
	//loop over the movies
	for index, item := range movies {
		if item.ID == params["id"] {
			//delete the movie with the id
			movies = append(movies[:index], movies[index+1:]...)

			//add new movie
			var movie Movie
			_ = json.NewDecoder(request.Body).Decode(&movie)
			movie.ID = params["id"]
			movies = append(movies, movie)
			json.NewEncoder(writer).Encode(movies)
		}
	}
}

//Create a movie
func createMovies(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	var movie Movie
	_ = json.NewDecoder(request.Body).Decode(&movie)

	//Define movie id
	movie.ID = strconv.Itoa(rand.Intn(1000))
	movies = append(movies, movie)

	json.NewEncoder(writer).Encode(movies)
}

//Get single movie
func getMovie(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)
	for _, item := range movies {
		if item.ID == params["id"] {
			json.NewEncoder(writer).Encode(item)
			return
		}
	}
}

//Get all movies
func getMovies(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(movies)
}
