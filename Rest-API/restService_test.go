package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	server *httptest.Server
	reader io.Reader
	//signUpURL string
)

func TestSignUp(t *testing.T) {
	server := httptest.NewServer(TestHandlers())
	//signUpURL := server.URL + "/signup"
	userJson := `{"Email": "a@gmail.com", "FirstName": "a", "LastName": "b", "Password": "a"}`

	reader = strings.NewReader(userJson)

	request, err := http.NewRequest("POST", server.URL+"/signup", reader)

	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}
	body, _ := ioutil.ReadAll(res.Body)
	var data map[string]bool
	json.Unmarshal(body, &data)
	if data["Success"] == true {
		fmt.Printf("Test Passed.")
	}
}

func TestFollowUser(t *testing.T) {
	server := httptest.NewServer(TestHandlers())
	//signUpURL := server.URL + "/signup"
	userJson := `{"userID": "abc@gmail.com", "userToFollow": "bcd@gmail.com"}`

	reader = strings.NewReader(userJson)

	request, err := http.NewRequest("POST", server.URL+"/followUser", reader)

	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}
	body, _ := ioutil.ReadAll(res.Body)
	var data map[string]bool
	json.Unmarshal(body, &data)
	if data["Success"] == true {
		fmt.Printf("Test Passed. \n")
	}
}

func TestCreateTweet(t *testing.T) {
	server := httptest.NewServer(TestHandlers())
	//signUpURL := server.URL + "/signup"
	userJson := `{"userID": "abc@gmail.com", "userTweet": "Testing This Good"}`

	reader = strings.NewReader(userJson)

	request, err := http.NewRequest("POST", server.URL+"/createTweet", reader)

	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}
	body, _ := ioutil.ReadAll(res.Body)
	var data map[string]bool
	json.Unmarshal(body, &data)
	if data["Success"] == true {
		fmt.Printf("Test Passed. \n")
	}
}

func TestShowTweets(t *testing.T) {
	server := httptest.NewServer(TestHandlers())
	//signUpURL := server.URL + "/signup"
	userJson := `{"userID": "abc@gmail.com"}`

	reader = strings.NewReader(userJson)

	request, err := http.NewRequest("POST", server.URL+"/showTweets", reader)

	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}
	body, _ := ioutil.ReadAll(res.Body)
	var result []UserTweet
	json.Unmarshal(body, &result)
	if len(result) > 0 {
		fmt.Print("Test Passed. \n")
	}
}

func TestCancelHandler(t *testing.T) {
	server := httptest.NewServer(TestHandlers())
	fmt.Printf("Checking when user exists. Success should be true. \n")

	//signUpURL := server.URL + "/signup"
	userJson := `{"email": "abc@gmail.com", "password":"abc"}`

	reader = strings.NewReader(userJson)

	request, err := http.NewRequest("POST", server.URL+"/cancel", reader)

	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}
	body, _ := ioutil.ReadAll(res.Body)
	var data map[string]bool
	json.Unmarshal(body, &data)
	if data["Success"] == true {
		fmt.Printf("")
	}
	fmt.Printf("Checking when user doesn't exist. Success should be false. \n")
	userJson = `{"email": "abcddd@gmail.com", "password":"abc"}`

	reader = strings.NewReader(userJson)

	request, err = http.NewRequest("POST", server.URL+"/cancel", reader)

	res, err = http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}
	body, _ = ioutil.ReadAll(res.Body)
	var data1 map[string]bool
	json.Unmarshal(body, &data1)
	if data1["Success"] == false {
		fmt.Printf("Test Passed. \n")
	}
}

func TestLoginHandler(t *testing.T) {
	server := httptest.NewServer(TestHandlers())
	fmt.Printf("Checking when user exists. Success should be true. \n")

	//signUpURL := server.URL + "/signup"
	userJson := `{"email": "abc@gmail.com", "password":"abc"}`

	reader = strings.NewReader(userJson)

	request, err := http.NewRequest("POST", server.URL+"/cancel", reader)

	res, err := http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}
	body, _ := ioutil.ReadAll(res.Body)
	var data map[string]bool
	json.Unmarshal(body, &data)
	if data["Success"] == true {
		fmt.Printf("")
	}
	fmt.Printf("Checking when user doesn't exist. Success should be false. \n")
	userJson = `{"email": "abcddd@gmail.com", "password":"abc"}`

	reader = strings.NewReader(userJson)

	request, err = http.NewRequest("POST", server.URL+"/cancel", reader)

	res, err = http.DefaultClient.Do(request)

	if err != nil {
		t.Error(err)
	}
	body, _ = ioutil.ReadAll(res.Body)
	var data1 map[string]bool
	json.Unmarshal(body, &data1)
	if data1["Success"] == false {
		fmt.Printf("Test Passed. \n")
	}
}
