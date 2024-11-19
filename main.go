package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

type User struct {
	Name	string	`json:"name"`
}

var UserCache = make(map[int]User)

var cacheMutex sync.RWMutex

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", HandleHome)
	mux.HandleFunc("POST /createUser", HandleCreateUser)
	mux.HandleFunc("GET /users/{id}", HandleGetUser)
	mux.HandleFunc("DELETE /users/{id}", HandleDeleteUser)

	println("Server listening at port 3000")

	http.ListenAndServe(":3000", mux)
}

func HandleHome(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Hello from home")
}

func HandleCreateUser(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","application/json")
	w.Header().Set("Access-Control-Allow-Origin","*")
	w.Header().Set("Access-Control-Allow-Methods","POST")
	w.Header().Set("Access-Control-Allow-Headers","Content-Type")
	
	var user User;
	err := json.NewDecoder(r.Body).Decode(&user)
	
	if err!=nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if user.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return 
	}

	cacheMutex.Lock()
	UserCache[len(UserCache) + 1] = user
	cacheMutex.Unlock()

	fmt.Println(UserCache)
	json.NewEncoder(w).Encode(user)
}

func HandleGetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","application/json")
	w.Header().Set("Access-Control-Allow-Origin","*")
	w.Header().Set("Access-Control-Allow-Methods","GET")
	
	id, err := strconv.Atoi(r.PathValue("id"))
	if err!=nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	cacheMutex.RLock()
	user, ok := UserCache[id]
	cacheMutex.RUnlock()

	if !ok {
		http.Error(w, fmt.Sprintf("User with id: %d not found", id), http.StatusNotFound)
		return
	}

	j, err := json.Marshal(user)
	if err!=nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err!=nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}	

	if _, ok := UserCache[id]; !ok {
		http.Error(w, fmt.Sprintf("User with id: %d not found", id), http.StatusNotFound)
		return
	}

	cacheMutex.Lock()
	delete(UserCache, id)
	cacheMutex.Unlock()

	fmt.Println(UserCache)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "UserId %d deleted successfully \n", id)
}
