package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

type User struct {
	Email     string `json:"Email,omitempty"`
	FirstName string `json:"FirstName,omitempty"`
	LastName  string `json:"LastName,omitempty"`
	Password  string `json:"Password,omitempty"`
}

type UserTweet struct {
	Email string `json:"userId,omitempty"`
	Tweet string `json:"userTweet,omitempty"`
}

type UserTweetList struct {
	Email     string `json:"userId,omitempty"`
	AllTweets []UserTweet
}
type OperationDetails struct {
	OperationName string
	Cmd           []byte
	ClientEmail   string //Gets email of the client as identifier
}

var UserTweetsMap = make(map[string][]UserTweet)

var UserFollower = make(map[string][]string)

// TODO make maps pointing with email bool || pointer to node for constant lookups

var UserMap = make(map[string]*User)

var OperationLog []OperationDetails

var CommitLog []int

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/userReplicate", userReplicateHandler).Methods("POST")
	router.HandleFunc("/tweetReplicate", tweetReplicateHandler).Methods("POST")
	router.HandleFunc("/followerReplicate", followerReplicateHandler).Methods("POST")
	router.HandleFunc("/deleteUserReplicate", deleteUserReplicateHandler).Methods("POST")
	router.HandleFunc("/CommitIndexHandler", CommitIndexHandler).Methods("POST")
	http.ListenAndServe(":9002", router)
}

//The Master will call this endpoint to replicate the latest user map it has.
func deleteUserReplicateHandler(w http.ResponseWriter, r *http.Request) {
	body, e := ioutil.ReadAll(r.Body)
	var result map[string]bool
	result = make(map[string]bool)
	if e == nil {
		var params map[string]string
		json.Unmarshal(body, &params)
		user := params["userEmail"]
		OperationLog = append(OperationLog, OperationDetails{"deleteUserReplicateHandler", body, user})
		delete(UserMap, user)
		result["Success"] = true
	} else {
		result["Success"] = false
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

//The Master will call this endpoint to replicate user-followers
func followerReplicateHandler(w http.ResponseWriter, r *http.Request) {
	body, e := ioutil.ReadAll(r.Body)
	var result map[string]bool
	result = make(map[string]bool)
	if e == nil {
		var params map[string]string
		json.Unmarshal(body, &params)
		user := params["userId"]
		userToFollow := params["userToFollow"]
		OperationLog = append(OperationLog, OperationDetails{"followerReplicateHandler", body, user})
		UserFollower[user] = append(UserFollower[user], userToFollow)
		result["Success"] = true
	} else {
		result["Success"] = false
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

//The Master will call this endpoint to replicate tweets
func tweetReplicateHandler(w http.ResponseWriter, r *http.Request) {
	body, e := ioutil.ReadAll(r.Body)
	var result map[string]bool
	result = make(map[string]bool)
	if e == nil {
		var params map[string]string
		json.Unmarshal(body, &params)
		user := params["userId"]
		userTweet := params["userTweet"]
		OperationLog = append(OperationLog, OperationDetails{"tweetReplicateHandler", body, user})
		fmt.Printf("Log: %v\n", OperationLog)
		UserTweetsMap[user] = append(UserTweetsMap[user], UserTweet{Email: user, Tweet: userTweet})
		result["Success"] = true
	} else {
		result["Success"] = false
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

//The master will call this endpoint to replicater a user.
func userReplicateHandler(w http.ResponseWriter, r *http.Request) {
	body, e := ioutil.ReadAll(r.Body)
	newUser := new(User)
	var result map[string]bool
	result = make(map[string]bool)
	if e == nil {
		var params map[string]string
		json.Unmarshal(body, &params)
		newUser.Email = params["Email"]
		newUser.FirstName = params["FirstName"]
		newUser.LastName = params["LastName"]
		newUser.Password = params["Password"]
		UserMap[newUser.Email] = newUser
		fmt.Printf("User values %v", UserMap)
		OperationLog = append(OperationLog, OperationDetails{"userReplicateHandler", body, newUser.Email})
		result["Success"] = true
	} else {
		result["Success"] = false
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

func CommitIndexHandler(w http.ResponseWriter, r *http.Request) {
	body, e := ioutil.ReadAll(r.Body)
	var result map[string]bool
	result = make(map[string]bool)
	if e == nil {
		var params map[string]int
		json.Unmarshal(body, &params)
		CommitIndex := params["Commit"]
		OperationIndex := params["Operation"]
		CommitLog = append(CommitLog, OperationIndex)
		fmt.Printf("CommitLog %v \n\n Operation Log %v", CommitLog, OperationLog)
		if len(CommitLog) == CommitIndex {
			result["Success"] = true
		} else {
			result["Success"] = true
		}
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
