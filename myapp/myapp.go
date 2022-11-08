package myapp

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// user struct
type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreateAt  time.Time `json:"create_at"`
}

var userMap map[int]*User
var idSeq int

// index handler
func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("url : ", r.URL)

	fmt.Fprintf(w, "hello world!")
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("url : ", r.URL)

	if len(userMap) == 0 {
		// there is no users
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "no users")
	}

	users := []*User{}
	for _, u := range userMap {
		//users = append(users,u)
		users = append(users, u)
	}
	data, _ := json.Marshal(users)

	w.Header().Add("Content-Type", "application/json") // before WriteHeader
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(data))
}

// create user
func createUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("url : ", r.URL)

	user := new(User)
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	idSeq++
	// created user
	user.ID = idSeq
	user.CreateAt = time.Now()
	userMap[user.ID] = user

	w.Header().Add("Content-Type", "application/json") // before WriteHeader
	w.WriteHeader(http.StatusCreated)
	data, _ := json.Marshal(user)
	fmt.Fprint(w, string(data))
}

// get user
func getUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Print("vars : ", vars)
	//fmt.Fprintf(w, "User Id: %v", vars["id"])

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}
	user, ok := userMap[id]
	if !ok {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "No User ID : ", id)
		return
	}

	w.Header().Add("Content-Type", "application/json") // before WriteHeader
	w.WriteHeader(http.StatusOK)

	data, _ := json.Marshal(user)
	fmt.Fprint(w, string(data))
}

// delete user
func deleteUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	_, ok := userMap[id]
	if !ok {
		// not found
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "No User ID : ", id)
		return
	}
	delete(userMap, id)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User ID[%d] Deleted\n", id)
}

// update user
func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	reqUser := new(User)
	err := json.NewDecoder(r.Body).Decode(reqUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}
	user, ok := userMap[reqUser.ID]
	if !ok {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "not found User ID : ", reqUser.ID)
		return
	}
	if reqUser.FirstName != "" {
		user.FirstName = reqUser.FirstName
	}
	if reqUser.LastName != "" {
		user.LastName = reqUser.LastName
	}
	if reqUser.Email != "" {
		user.Email = reqUser.Email
	}

	user.CreateAt = time.Now()

	w.Header().Add("Content-Type", "application/json") // before WriteHeader
	w.WriteHeader(http.StatusOK)

	data, _ := json.Marshal(user)
	fmt.Fprint(w, string(data))
}

// NewHandler make a new myapp handler
func NewHandler() http.Handler {
	userMap = make(map[int]*User)
	idSeq = 0

	//mux := http.NewServeMux()
	mux := mux.NewRouter() // use gorilla mux
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/users", usersHandler).Methods("GET")
	mux.HandleFunc("/users", createUserHandler).Methods("POST")
	mux.HandleFunc("/users", updateUserHandler).Methods("PUT")
	mux.HandleFunc("/users/{id:[0-9]+}", getUserInfoHandler).Methods("GET")
	mux.HandleFunc("/users/{id:[0-9]+}", deleteUserInfoHandler).Methods("DELETE")
	return mux
}
