package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

// TODO make maps pointing with email bool || pointer to node for constant lookups

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", loginHandler)
	router.HandleFunc("/cancel", cancelHandler)
	router.HandleFunc("/goodbye", goodbyeHandler)
	router.HandleFunc("/signup", signupHandler)
	router.HandleFunc("/home/{id}", homeHandler)
	router.HandleFunc("/home/{id}/tweet", createTweetHandler)
	router.HandleFunc("/home/{id}/follow", followFriendHandler)
	http.ListenAndServe(":8000", router)
}

func goodbyeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("goodbye.html"))
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}
}

func cancelHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("cancel.html"))
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}
	// send rest api req to cancel the account
	jsonData := map[string]string{"email": r.FormValue("email"), "password": r.FormValue("password")}
	jsonValue, _ := json.Marshal(jsonData)
	response, err := http.Post("http://localhost:9000/cancel", "application/json", bytes.NewBuffer(jsonValue))
	if err == nil {
		body, e := ioutil.ReadAll(response.Body)
		if e == nil {
			var data map[string]bool
			json.Unmarshal(body, &data)
			if data["Success"] == true {
				http.Redirect(w, r, "/goodbye", 302)
				return
			} else {
				tmpl.Execute(w, struct{ Authfail bool }{true})
			}
		}
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("login.html"))
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	// Call rest api to check authentication. send user email and password and get response
	jsonData := map[string]string{"email": r.FormValue("email"), "password": r.FormValue("password")}
	jsonValue, _ := json.Marshal(jsonData)
	response, err := http.Post("http://localhost:9000/login", "application/json", bytes.NewBuffer(jsonValue))
	if err == nil {
		body, e := ioutil.ReadAll(response.Body)
		if e == nil {
			var data map[string]bool
			json.Unmarshal(body, &data)
			if data["Success"] == true {
				homeURL := "/home/" + r.FormValue("email")
				http.Redirect(w, r, homeURL, 302)
				return
			} else {
				tmpl.Execute(w, struct{ Authfail bool }{true})
			}
		}
	}

}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("signUp.html"))
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	jsonData := map[string]string{"Email": r.FormValue("new_email"), "FirstName": r.FormValue("new_fname"),
		"LastName": r.FormValue("new_lname"), "Password": r.FormValue("new_password")}
	jsonValue, _ := json.Marshal(jsonData)
	response, err := http.Post("http://localhost:9000/signup", "application/json", bytes.NewBuffer(jsonValue))
	if err == nil {
		body, e := ioutil.ReadAll(response.Body)
		if e == nil {
			var data map[string]bool
			json.Unmarshal(body, &data)
			fmt.Print(data)
			tmpl.Execute(w, struct {
				Success          bool
				DuplicateAccount bool
			}{Success: data["Success"], DuplicateAccount: data["DuplicateAccount"]})
			return
		}
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("home.html"))
	if r.Method != http.MethodPost {
		//get email from URL so that we can request rest to give tweets for all other users except this one
		vars := mux.Vars(r)
		userSignedIn := vars["id"]
		//send post req to rest API to get a list of tweets to show
		jsonData := map[string]string{"userId": userSignedIn}
		jsonValue, _ := json.Marshal(jsonData)
		response, err := http.Post("http://localhost:9000/showTweets", "application/json", bytes.NewBuffer(jsonValue))
		//response from API
		if err == nil {
			body, e := ioutil.ReadAll(response.Body)
			if e == nil {
				//define tweet struct to get jsonData
				type tweets struct {
					Email string `json:"userId,omitempty"`
					Tweet string `json:"userTweet,omitempty"`
				}

				var data []tweets
				json.Unmarshal(body, &data)
				type tweetContainer struct {
					ContainerEmail string `json:"userId,omitempty"`
					Containerdata  []tweets
				}
				containerObject := tweetContainer{ContainerEmail: userSignedIn, Containerdata: data}
				fmt.Print(containerObject)
				tmpl.Execute(w, containerObject)
			}
		}
		return
	}
}

func createTweetHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("tweet.html"))
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	vars := mux.Vars(r)
	userSignedIn := vars["id"]
	//send post req to rest API to get a list of tweets to show
	jsonData := map[string]string{"userId": userSignedIn, "userTweet": r.FormValue("tweet")}
	jsonValue, _ := json.Marshal(jsonData)
	response, err := http.Post("http://localhost:9000/createTweet", "application/json", bytes.NewBuffer(jsonValue))
	//response from API
	if err == nil {
		body, e := ioutil.ReadAll(response.Body)
		if e == nil {
			var data map[string]bool
			json.Unmarshal(body, &data)
			if data["Success"] {
				homeURL := "/home/" + userSignedIn
				http.Redirect(w, r, homeURL, 302)
				return
			} else {
				tmpl.Execute(w, struct{ Tweetfail bool }{true})
			}
		}
	}
}

func followFriendHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("follow.html"))
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	vars := mux.Vars(r)
	userSignedIn := vars["id"]
	//send post req to rest API to get a list of tweets to show
	jsonData := map[string]string{"userId": userSignedIn, "userToFollow": r.FormValue("user_follow")}
	jsonValue, _ := json.Marshal(jsonData)
	response, err := http.Post("http://localhost:9000/followUser", "application/json", bytes.NewBuffer(jsonValue))
	//response from API
	if err == nil {
		body, e := ioutil.ReadAll(response.Body)
		if e == nil {
			var data map[string]bool
			json.Unmarshal(body, &data)
			if data["Success"] {
				homeURL := "/home/" + userSignedIn
				http.Redirect(w, r, homeURL, 302)
				return
			} else {
				tmpl.Execute(w, struct{ Followfail bool }{true})
			}
		}
	}
}
