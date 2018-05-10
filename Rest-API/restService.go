package main

import (
	"bytes"
	"container/list"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
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

var Email string

var UserMap = make(map[string]*User)

var UserList = list.New() // Only for distributing among servers while making distributed.
//Users := make(map[Emai])
// user email to tweet map TODO: add structs
var tweets []UserTweet

//var OperationLog = make([]OperationDetails)

var OperationLog []OperationDetails

var CommitLog []int //Access the op number from Operation Log

var totalServers = 3

func main() {
	//add dummy tweet data
	UserMap["abc@gmail.com"] = &User{"abc@gmail.com", "abc", "abc", "abc", nil}
	UserElementList := sortedInsert(UserMap["abc@gmail.com"])
	UserMap["abc@gmail.com"].PosInList = UserElementList

	UserMap["abcd@gmail.com"] = &User{"abcd@gmail.com", "abcd", "abcd", "abcd", nil}
	UserElementList = sortedInsert(UserMap["abcd@gmail.com"])
	UserMap["abcd@gmail.com"].PosInList = UserElementList

	UserMap["bcd@gmail.com"] = &User{"bcd@gmail.com", "bcd", "bcd", "bcd", nil}
	UserElementList = sortedInsert(UserMap["bcd@gmail.com"])
	UserMap["bcd@gmail.com"].PosInList = UserElementList

	UserTweetsMap["abc@gmail.com"] = append(UserTweetsMap["abc@gmail.com"], UserTweet{Email: "abc@gmail.com", Tweet: "first tweet"})
	UserTweetsMap["abcd@gmail.com"] = append(UserTweetsMap["abcd@gmail.com"], UserTweet{Email: "abcd@gmail.com", Tweet: "second tweet"})
	UserTweetsMap["abc@gmail.com"] = append(UserTweetsMap["abc@gmail.com"], UserTweet{Email: "abc@gmail.com", Tweet: "third tweet"})
	UserTweetsMap["bcd@gmail.com"] = append(UserTweetsMap["bcd@gmail.com"], UserTweet{Email: "bcd@gmail.com", Tweet: "fourth tweet"})

	UserFollower["abc@gmail.com"] = append(UserFollower["abc@gmail.com"], "abcd@gmail.com", "bcd@gmail.com")

	router := mux.NewRouter()
	go router.HandleFunc("/login", loginHandler).Methods("POST")
	go router.HandleFunc("/cancel", cancelHandler).Methods("POST")
	go router.HandleFunc("/signup", signupHandler).Methods("POST")
	go router.HandleFunc("/showTweets", showTweetsHandler).Methods("POST")
	go router.HandleFunc("/createTweet", createTweets).Methods("POST")
	go router.HandleFunc("/followUser", followUser).Methods("POST")
	http.ListenAndServe(":9000", router)
}

func TestHandlers() *mux.Router {
	UserMap["abc@gmail.com"] = &User{"abc@gmail.com", "abc", "abc", "abc", nil}
	UserElementList := sortedInsert(UserMap["abc@gmail.com"])
	UserMap["abc@gmail.com"].PosInList = UserElementList

	UserMap["abcd@gmail.com"] = &User{"abcd@gmail.com", "abcd", "abcd", "abcd", nil}
	UserElementList = sortedInsert(UserMap["abcd@gmail.com"])
	UserMap["abcd@gmail.com"].PosInList = UserElementList

	UserMap["bcd@gmail.com"] = &User{"bcd@gmail.com", "bcd", "bcd", "bcd", nil}
	UserElementList = sortedInsert(UserMap["bcd@gmail.com"])
	UserMap["bcd@gmail.com"].PosInList = UserElementList

	UserTweetsMap["abc@gmail.com"] = append(UserTweetsMap["abc@gmail.com"], UserTweet{Email: "abc@gmail.com", Tweet: "first tweet"})
	UserTweetsMap["abcd@gmail.com"] = append(UserTweetsMap["abcd@gmail.com"], UserTweet{Email: "abcd@gmail.com", Tweet: "second tweet"})
	UserTweetsMap["abc@gmail.com"] = append(UserTweetsMap["abc@gmail.com"], UserTweet{Email: "abc@gmail.com", Tweet: "third tweet"})
	UserTweetsMap["bcd@gmail.com"] = append(UserTweetsMap["bcd@gmail.com"], UserTweet{Email: "bcd@gmail.com", Tweet: "fourth tweet"})

	UserFollower["abc@gmail.com"] = append(UserFollower["abc@gmail.com"], "abcd@gmail.com", "bcd@gmail.com")
	router := mux.NewRouter()
	router.HandleFunc("/login", loginHandler).Methods("POST")
	router.HandleFunc("/cancel", cancelHandler).Methods("POST")
	router.HandleFunc("/signup", signupHandler).Methods("POST")
	router.HandleFunc("/showTweets", showTweetsHandler).Methods("POST")
	router.HandleFunc("/createTweet", createTweets).Methods("POST")
	router.HandleFunc("/followUser", followUser).Methods("POST")
	return router
}

func followUser(w http.ResponseWriter, r *http.Request) {
	var result map[string]bool
	result = make(map[string]bool)
	body, e := ioutil.ReadAll(r.Body)
	count := 0
	if e == nil {
		var params map[string]string
		json.Unmarshal(body, &params)
		user := params["userId"]
		userToFollow := params["userToFollow"]
		//Check if the user to follow actually exists
		_, ok := UserMap[userToFollow]
		if ok {
			UserFollower[user] = append(UserFollower[user], userToFollow)
			//Replicate on backend.
			jsonData := map[string]FollowUserReplicate{}
			jsonData["user"] = FollowUserReplicate{user, userToFollow, OperationLog, CommitLog}
			jsonValue, _ := json.Marshal(jsonData)
			OperationLog = append(OperationLog, OperationDetails{"followUser", body, user})
			fmt.Printf("\n FOLLOW LOG \n %v", OperationLog)
			go replicateData(jsonValue, &count, "/followerReplicate", len(OperationLog))
			//check if majority appended.
			result["Success"] = true
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

func createTweets(w http.ResponseWriter, r *http.Request) {
	var result map[string]bool
	result = make(map[string]bool)
	count := 0
	body, e := ioutil.ReadAll(r.Body)
	if e == nil {
		var params map[string]string
		json.Unmarshal(body, &params)
		user := params["userId"]
		userTweet := params["userTweet"]
		fmt.Print(user)
		//make new  Tweet and store in slice
		UserTweetsMap[user] = append(UserTweetsMap[user], UserTweet{Email: user, Tweet: userTweet})
		//Replicate on backend.
		OperationLog = append(OperationLog, OperationDetails{"createTweets", body, user})
		jsonData := map[string]UserTweetReplicate{}
		jsonData["user"] = UserTweetReplicate{user, userTweet, OperationLog, CommitLog}
		jsonValue, _ := json.Marshal(jsonData)
		go replicateData(jsonValue, &count, "/tweetReplicate", len(OperationLog))
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
	// var doneChannel chan bool
	// var done bool
	body, e := ioutil.ReadAll(r.Body)
	if e == nil {
		var params map[string]string
		json.Unmarshal(body, &params)
		userID := params["userId"]
		for _, follower := range UserFollower[userID] {
			if tweets, ok := UserTweetsMap[follower]; ok {
				for _, tweet := range tweets {
					result = append(result, tweet)
					//fmt.Printf("Tweets %v", result)
				}
			}
		}
		// for _, tweet := range tweets {
		// 	if tweet.Email != userID {
		// 		//add to results
		// 		result = append(result, UserTweet{Email: tweet.Email, Tweet: tweet.Tweet})
		// 	}
		// }
		//fmt.Print(result)
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
	count := 0
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
				//Replicate on backend.
				jsonData := map[string]CancelUserReplicate{}
				jsonData["user"] = CancelUserReplicate{userEmail, userPassword,
					OperationLog, CommitLog}
				jsonValue, _ := json.Marshal(jsonData)
				OperationLog = append(OperationLog, OperationDetails{"cancelHandler", body, userEmail})
				go replicateData(jsonValue, &count, "/deleteUserReplicate", len(OperationLog))
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
	count := 0

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
		//Replicate on backend.
		OperationLog = append(OperationLog, OperationDetails{"signupHandler", body, newUser.Email})
		jsonData := map[string]CreateUserReplicate{}
		jsonData["user"] = CreateUserReplicate{newUser.Email, newUser.FirstName, newUser.LastName, newUser.Password, OperationLog, CommitLog}
		jsonValue, _ := json.Marshal(jsonData)
		go replicateData(jsonValue, &count, "/userReplicate", len(OperationLog))

		// TODO if majority appends to the log send success to front end.
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

func replicateCommitIndex(commitIndex int, OperationNum int) {
	jsonData := map[string]int{"Commit": commitIndex, "Operation": OperationNum}
	jsonValue, _ := json.Marshal(jsonData)
	count := 0
	go replicateData(jsonValue, &count, "/CommitIndexHandler", -1)
	return
}

func replicateData(jsonValue []byte, count *int, handlerName string, opNumber int) {
	done := make(chan bool)
	for i := 1; i < totalServers; i++ {
		go func(i int, count *int) {
			response, err := http.Post("http://localhost:900"+strconv.Itoa(i)+handlerName, "application/json", bytes.NewBuffer(jsonValue))
			if err == nil {
				body, e := ioutil.ReadAll(response.Body)
				if e == nil {
					var data map[string]bool
					json.Unmarshal(body, &data)
					if data["Success"] == true {
						*count++
						done <- true
					}
				}
			}
		}(i, count)
	}
	for i := 0; i < totalServers/2; i++ {
		<-done
	}
	if opNumber > 0 {
		CommitLog = append(CommitLog, opNumber)
		replicateCommitIndex(len(CommitLog), opNumber)
		//fmt.Printf("CommitLog %v \n\n Operation Log %v", CommitLog, OperationLog)
	}
	//fmt.Printf("\n\n\nDONE WITH REPLICATIOn %v\n\n\n", *count)
}
