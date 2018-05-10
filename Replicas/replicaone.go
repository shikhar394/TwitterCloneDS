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

type UserTweetReplicate struct {
	UserID       string
	UserTweet    string
	OperationLog []OperationDetails
	CommitLog    []int
}

type FollowUserReplicate struct {
	UserID       string
	UserToFollow string
	OperationLog []OperationDetails
	CommitLog    []int
}

type CancelUserReplicate struct {
	UserID       string
	Password     string
	OperationLog []OperationDetails
	CommitLog    []int
}

type CreateUserReplicate struct {
	Email        string
	FirstName    string
	LastName     string
	Password     string
	OperationLog []OperationDetails
	CommitLog    []int
}

var UserTweetsMap = make(map[string][]UserTweet)

var UserFollower = make(map[string][]string)

// TODO make maps pointing with email bool || pointer to node for constant lookups

var UserMap = make(map[string]*User)

var OperationLog []OperationDetails

var CommitLog []int

var PRIMARYPORT = 9000

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/userReplicate", userReplicateHandler).Methods("POST")
	router.HandleFunc("/tweetReplicate", tweetReplicateHandler).Methods("POST")
	router.HandleFunc("/followerReplicate", followerReplicateHandler).Methods("POST")
	router.HandleFunc("/deleteUserReplicate", deleteUserReplicateHandler).Methods("POST")
	router.HandleFunc("/CommitIndexHandler", CommitIndexHandler).Methods("POST")
	http.ListenAndServe(":9001", router)
}

//The Master will call this endpoint to replicate the latest user map it has.
func deleteUserReplicateHandler(w http.ResponseWriter, r *http.Request) {
	body, e := ioutil.ReadAll(r.Body)
	var result map[string]bool
	result = make(map[string]bool)
	if e == nil {
		var params map[string]CancelUserReplicate
		json.Unmarshal(body, &params)
		CancelUserReplicate := params["user"]
		user := CancelUserReplicate.UserID
		//password := CancelUserReplicate.Password
		PrimaryCommitLog := CancelUserReplicate.CommitLog
		PrimaryOperationLog := CancelUserReplicate.OperationLog
		OperationLog = append(OperationLog, OperationDetails{"deleteUserReplicateHandler", body, user})
		delete(UserMap, user)
		if len(OperationLog) < len(PrimaryOperationLog) {
			go fixLogs(PrimaryOperationLog, PrimaryCommitLog)
		} else {
			fmt.Print("Logs perfect \n\n")
		}
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
		var params map[string]FollowUserReplicate
		json.Unmarshal(body, &params)
		FollowUserDetails := params["user"]
		user := FollowUserDetails.UserID
		userToFollow := FollowUserDetails.UserToFollow
		PrimaryCommitLog := FollowUserDetails.CommitLog
		PrimaryOperationLog := FollowUserDetails.OperationLog
		OperationLog = append(OperationLog, OperationDetails{"followerReplicateHandler", body, user})
		UserFollower[user] = append(UserFollower[user], userToFollow)
		if len(OperationLog) < len(PrimaryOperationLog) {
			go fixLogs(PrimaryOperationLog, PrimaryCommitLog)
		} else {
			fmt.Print("Logs perfect \n\n")
		}
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
		var params map[string]UserTweetReplicate
		json.Unmarshal(body, &params)
		fmt.Printf("Userssss %v\n", params)
		TweetDetails := params["user"]
		PrimaryCommitLog := TweetDetails.CommitLog
		PrimaryOperationLog := TweetDetails.OperationLog
		user := TweetDetails.UserID
		userTweet := TweetDetails.UserTweet
		// if !ok || !ok1 {
		// 	fmt.Printf("Error converting %v %v %v %v\n", ok, ok1, ok0, ok2)
		// }

		OperationLog = append(OperationLog, OperationDetails{"tweetReplicateHandler", body, user})
		fmt.Printf("Log: %v \n PrimaryLog %v\n CommitLog %v", len(OperationLog), PrimaryOperationLog, PrimaryCommitLog)
		UserTweetsMap[user] = append(UserTweetsMap[user], UserTweet{Email: user, Tweet: userTweet})
		if len(OperationLog) < len(PrimaryOperationLog) {
			go fixLogs(PrimaryOperationLog, PrimaryCommitLog)
		} else {
			fmt.Print("Logs perfect \n\n")
		}
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
		var params map[string]CreateUserReplicate
		json.Unmarshal(body, &params)
		NewUserReplicate := params["user"]
		newUser.Email = NewUserReplicate.Email
		newUser.FirstName = NewUserReplicate.FirstName
		newUser.LastName = NewUserReplicate.LastName
		newUser.Password = NewUserReplicate.Password
		PrimaryOperationLog := NewUserReplicate.OperationLog
		PrimaryCommitLog := NewUserReplicate.CommitLog
		UserMap[newUser.Email] = newUser
		OperationLog = append(OperationLog, OperationDetails{"userReplicateHandler", body, newUser.Email})
		fmt.Printf("User values %v", UserMap)
		if len(OperationLog) < len(PrimaryOperationLog) {
			go fixLogs(PrimaryOperationLog, PrimaryCommitLog)
		} else {
			fmt.Print("Logs perfect \n\n")
		}
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

func fixLogs(PrimaryOperationLog []OperationDetails, PrimaryCommitLog []int) {
	OperationLog = PrimaryOperationLog[:]
	CommitLog = PrimaryCommitLog[:]
	fmt.Printf("Fixing Logs %v\n\n", OperationLog)
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
