package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type User struct {
	Email     string `json:"Email,omitempty"`
	FirstName string `json:"FirstName,omitempty"`
	LastName  string `json:"LastName,omitempty"`
	Password  string `json:"Password,omitempty"`
	PosInList *list.Element
}

type UserTweet struct {
	Email string `json:"userId,omitempty"`
	Tweet string `json:"userTweet,omitempty"`
}

// TODO make maps pointing with email bool || pointer to node for constant lookups

var Email string

var UserMap = make(map[string]*User)

var UserList = list.New() // Only for distributing among servers while making distributed.
//Users := make(map[Emai])
// user email to tweet map TODO: add structs
var tweets []UserTweet

func main() {
	//add dummy tweet data
	tweets = append(tweets, UserTweet{Email: "b@b", Tweet: "bb tweet"})
	tweets = append(tweets, UserTweet{Email: "c@c", Tweet: "bb2 tweet"})
	tweets = append(tweets, UserTweet{Email: "a@a", Tweet: "ba tweet"})
	router := mux.NewRouter()
	router.HandleFunc("/login", loginHandler).Methods("POST")
	router.HandleFunc("/cancel", cancelHandler).Methods("POST")
	router.HandleFunc("/signup", signupHandler).Methods("POST")
	router.HandleFunc("/showTweets", showTweetsHandler).Methods("POST")
	router.HandleFunc("/createTweet", createTweets).Methods("POST")
	http.ListenAndServe(":9000", router)
}

func createTweets(w http.ResponseWriter, r *http.Request) {
	var result map[string]bool
	result = make(map[string]bool)
	body, e := ioutil.ReadAll(r.Body)
	if e == nil {
		var params map[string]string
		json.Unmarshal(body, &params)
		userID := params["userID"]
		userTweet := params["userTweet"]
		//make new  Tweet and store in slice
		tweets = append(tweets, UserTweet{Email: userID, Tweet: userTweet})
		result["Success"] = true
	} else {
		result["Success"] = false
	}

	//send data back
	jData, err := json.Marshal(result)
	if err != nil {
		panic(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
	return
}

func showTweetsHandler(w http.ResponseWriter, r *http.Request) {
	// a map of email and tweet to be sent to client
	var result []UserTweet
	body, e := ioutil.ReadAll(r.Body)
	if e == nil {
		var params map[string]string
		json.Unmarshal(body, &params)
		userID := params["userId"]
		for _, tweet := range tweets {
			if tweet.Email != userID {
				//add to results
				result = append(result, UserTweet{Email: tweet.Email, Tweet: tweet.Tweet})
			}
		}
		fmt.Print(result)
		//send data back
		jData, err := json.Marshal(result)
		if err != nil {
			panic(err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jData)
		return
	}
}

func cancelHandler(w http.ResponseWriter, r *http.Request) {
	body, e := ioutil.ReadAll(r.Body)
	var result map[string]bool
	result = make(map[string]bool)
	if e == nil {
		var params map[string]string
		json.Unmarshal(body, &params)
		userEmail := params["email"]
		userPassword := params["password"]
		User, ok := UserMap[userEmail]
		if ok {
			if User.Password == userPassword {
				//remove user
				UserList.Remove(User.PosInList)
				delete(UserMap, userEmail)
				result["Success"] = true
			} else {
				//fail - do not do anything
				result["Success"] = false
			}
		} else {
			result["Success"] = false
		}
	}
	//send data back
	jData, err := json.Marshal(result)
	if err != nil {
		panic(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
	return
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	body, e := ioutil.ReadAll(r.Body)
	var result map[string]bool
	result = make(map[string]bool)
	if e == nil {
		var params map[string]string
		json.Unmarshal(body, &params)
		userEmail := params["email"]
		userPassword := params["password"]
		User, ok := UserMap[userEmail]
		if ok {
			if User.Password == userPassword {
				result["Success"] = true
			} else {
				result["Success"] = false
			}
		} else {
			result["Success"] = false
		}
	}
	//send data back
	jData, err := json.Marshal(result)
	if err != nil {
		panic(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
	return
}

func sortedInsert(newUser *User) *list.Element {
	//fmt.Printf("\n\nSorted insert on email %v\n\n", newUser.Email)
	if UserList.Front() != nil {
		for CurrentUser := UserList.Front(); CurrentUser != nil; CurrentUser = CurrentUser.Next() {
			if strings.Compare(CurrentUser.Value.(*User).Email, (*newUser).Email) == 1 {
				return UserList.InsertBefore(newUser, CurrentUser)
			}
		}
	} else {
		return UserList.PushFront(newUser)
	}
	return UserList.PushBack(newUser)
	//fmt.Printf("\n\nDone inserting\n\n")
}

func readUsers() {
	for CurrentUser := UserList.Front(); CurrentUser != nil; CurrentUser = CurrentUser.Next() {
		fmt.Printf("User %v, Password %v\n", CurrentUser.Value.(*User).Email, CurrentUser.Value.(*User).Password)
	}
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	body, e := ioutil.ReadAll(r.Body)
	newUser := new(User)
	if e == nil {
		var params map[string]string
		json.Unmarshal(body, &params)
		newUser.Email = params["Email"]
		newUser.FirstName = params["FirstName"]
		newUser.LastName = params["LastName"]
		newUser.Password = params["Password"]
	}
	var result map[string]bool
	result = make(map[string]bool)
	//Users.PushBack(newUser) // TODO Make sorted instead and store in map for constant time lookup
	_, ok := UserMap[newUser.Email]
	if !ok {
		UserElementList := sortedInsert(newUser)
		newUser.PosInList = UserElementList
		UserMap[(newUser).Email] = newUser
		readUsers()
		//fmt.Print(UserList.Front().Value.(*User).FirstName)map
		result["Success"] = true
		result["DuplicateAccount"] = false
	} else {
		//tmpl.Execute(w, struct{ Success bool }{false})
		result["Success"] = false
		result["DuplicateAccount"] = true
	}
	jData, err := json.Marshal(result)
	if err != nil {
		panic(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
	return

}