package main

import (
	"encoding/json"
	"fmt"
	"golang-movie-server/utils"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Movie struct {
	ID        string    `json:"id"`
	Isbn      string    `json:"isbn"`
	Title     string    `json:"title"`
	Director  Director  `json:"director"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

type Director struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var Movies []Movie = []Movie{}

func init() {
	i := 0
	for i < 10 {
		id := uuid.New()
		now := time.Now()
		Movies = append(Movies, Movie{
			ID:    id.String(),
			Isbn:  id.String(),
			Title: utils.GenerateName(),
			Director: Director{
				FirstName: utils.GenerateName(),
				LastName:  utils.GenerateName(),
			},
			CreatedAt: now,
			UpdatedAt: now,
		})
		i++
	}
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/movies", getMovies).Methods(http.MethodGet)
	router.HandleFunc("/movies/{id}", getMovie).Methods(http.MethodGet)
	router.HandleFunc("/movies", createMovie).Methods(http.MethodPost)
	router.HandleFunc("/movies/{id}", updateMovie).Methods(http.MethodPut)
	router.HandleFunc("/movies/{id}", deleteMovie).Methods(http.MethodDelete)

	fmt.Println("Starting movies server at port 8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Movies)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, value := range Movies {
		if value.ID == params["id"] {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(value)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(Response{
		Code:    http.StatusNotFound,
		Message: fmt.Sprintf("movie with id `%v` not found", params["id"]),
	})
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	movie := Movie{}
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	defer r.Body.Close()
	id := uuid.New()
	now := time.Now()
	movie.ID = id.String()
	movie.CreatedAt = now
	movie.UpdatedAt = now
	Movies = append(Movies, movie)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Movies)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	movie := Movie{}
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	for index := range Movies {
		value := &Movies[index]
		if value.ID == params["id"] {
			if movie.Title != "" {
				value.Title = movie.Title
			}
			if movie.Director.FirstName != "" {
				value.Director.FirstName = movie.Director.FirstName
			}
			if movie.Director.LastName != "" {
				value.Director.LastName = movie.Director.LastName
			}
			value.UpdatedAt = time.Now()
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(Movies)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(Response{
		Code:    http.StatusNotFound,
		Message: fmt.Sprintf("movie with id `%v` not found", params["id"]),
	})
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index := range Movies {
		value := Movies[index]
		if value.ID == params["id"] {
			Movies = append(Movies[:index], Movies[index+1:]...)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(Movies)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(Response{
		Code:    http.StatusNotFound,
		Message: fmt.Sprintf("movie with id `%v` not found", params["id"]),
	})
}
