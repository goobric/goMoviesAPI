package main

import (
	"encoding/json" //used to create a new movie
	"fmt"
	"log"
	"math/rand"
	"net/http" //create a server in go
	"strconv"  //convert id string to int

	"github.com/gorilla/mux" //used for data
)

//struct of type Movie
type Movie struct {
	ID string `json:"id"`
	Isbn string `json:"isbn"`
	Title string `json:"title"`
	Director *Director `json:"director"` //pointer to Director struct
}

//struct of type Director
type Director struct {
	Firstname string `json:"firstname"`
	Lastname string `json:"lastname"`
}

var movies []Movie //slice of type Movie

//get all movies function
func getMovies(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies) //pass the movies slice to json encoder
}
//passing a pointer of the request that is sent

//delete a movie
func deleteMovie(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range movies {
		if item.ID == params["id"]{
			movies = append(movies[:index], movies[index+1:]...) //from index+1 will now be appended to the start of the slice
			break
		}
	}
	//return all remaining movies
	json.NewEncoder(w).Encode(movies)
}

//get individual movie
func getMovie(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	//loop through all the movies, use blank identifier as 'index' is not used
	for _, item := range movies {
		if item.ID == params["id"]{
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

//create a new movie
func createMovie(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)
	movie.ID = strconv.Itoa(rand.Intn(100000000))
	movies = append(movies, movie)
	json.NewEncoder(w).Encode(movie)
}

//update a movie, id from params, from mux.Vars inside (r), range over the movies
func updateMovie(w http.ResponseWriter, r *http.Request){
	//set json content type
	w.Header().Set("Content-Type", "application/json")
	//id from params
	params := mux.Vars(r)
	//loop over the movies, range
	//delete the movie with the id sent
	//add a new movie - the movie that is sent in the body of Postman
	for index, item := range movies {
		if item.ID == params["id"]{
			movies = append(movies[:index], movies[index+1:]...)
			var movie Movie
			_ = json.NewDecoder(r.Body).Decode(&movie)
			movie.ID = params["id"]
			movies = append(movies, movie)
			json.NewEncoder(w).Encode(movie)
			return
		}
	}
}

func main() {
	r := mux.NewRouter()

	//slice of movies
	movies = append(movies, Movie{ID: "01", Isbn:"67890", Title:"Movie One", Director: &Director{Firstname:"John", Lastname:"Doe"}})
	movies = append(movies, Movie{ID:"02", Isbn:"123456", Title:"Movie Two", Director: &Director{Firstname:"Sally", Lastname:"Smith"}})
	movies = append(movies, Movie{ID:"03", Isbn:"123456", Title:"Movie Three", Director: &Director{Firstname:"Zoe", Lastname:"Ali"}})
	//&reference of the object address of type Director pointer
	// 5 function handlers and 5 routes
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Printf("Starting server at port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", r))
}
