package main

import (
	"encoding/json" // for encoding the data into json when sending it to postman
	"fmt"           // for printing
	"log"           // for logging out data or error
	"math/rand"     // for creating random 'id' for new movies which will be added by the user
	"net/http"      // for creating server
	"strconv"       // for converting the 'id'(i.e integer) generated by 'math/rand' into 'string'

	"github.com/gorilla/mux" // for routing
)

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"` // unique number assigned to the film
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var movies []Movie

func getMovies(w http.ResponseWriter, r *http.Request) {
	// here `r` is a pointer of request that we'll send from our postman to this function and
	// `w` is the response writer which gives back the response from the server back to function or frontend
	w.Header().Set("Content Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	// set json content type
	w.Header().Set("Content Type", "application/json")
	params := mux.Vars(r) // here params is the `ID` that we pass from Postman will go as params to our function
	// loop over the movies, range
	// delete the movie with the ID that you've sent
	for index, item := range movies {

		if item.ID == params["id"] {
			// below we will see how we can delete a movie using `append()`
			movies = append(movies[:index], movies[index+1:]...)
			// above we are appending rest of the `movies` in place of the given `movie` in this way we are removing the given movie from the list
			break
		}
	}
	json.NewEncoder(w).Encode(movies)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	// set json content type
	w.Header().Set("Content Type", "application/json")
	params := mux.Vars(r)
	for _, item := range movies {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content Type", "application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)
	movie.ID = strconv.Itoa(rand.Intn(10000000))
	movies = append(movies, movie)
	json.NewEncoder(w).Encode(movie)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	// set json content type
	w.Header().Set("Content Type", "application/json")
	// params
	params := mux.Vars(r)
	// loop over the movies, range
	// delete the movie with the ID that you've sent
	// add a new movie(i.e the movie that we sent in the body of the postman)
	for index, item := range movies {
		if item.ID == params["id"] {
			// deleting
			movies = append(movies[:index], movies[index+1:]...)
			var movie Movie
			_ = json.NewDecoder(r.Body).Decode(&movie)
			// adding
			movie.ID = params["id"]
			movies = append(movies, movie)
			json.NewEncoder(w).Encode(movie)
		}
	}
}

func main() {
	r := mux.NewRouter()
	// movies slices
	movies = append(movies, Movie{ID: "1", Isbn: "348738", Title: "Movie One", Director: &Director{Firstname: "John", Lastname: "Doe"}}) // `&` is to get the address and `*` is used to access that address or the pointer
	movies = append(movies, Movie{ID: "2", Isbn: "328746", Title: "Movie Two", Director: &Director{Firstname: "Steve", Lastname: "Smith"}})
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Printf("Starting server at port 8000\n")
	log.Fatal(http.ListenAndServe(":8000", r))
}
