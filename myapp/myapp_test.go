package myapp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndex(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	res, err := http.Get(ts.URL)
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)
	data, _ := ioutil.ReadAll(res.Body)
	assert.Equal("hello world!", string(data))
}

func TestUsers(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	res, err := http.Get(ts.URL + "/users")
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)
	data, _ := ioutil.ReadAll(res.Body)
	assert.Contains(string(data), "Get UserInfo")
}

func TestGetUserInfo(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	res, err := http.Get(ts.URL + "/users/89")
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)
	data, _ := ioutil.ReadAll(res.Body)
	assert.Contains(string(data), "User ID : ")
}

func TestCreateUserInfo(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	res, err := http.Post(ts.URL+"/users", "application/json", strings.NewReader(`{"first_name":"JongHan", "last_name":"Park", "email":"jhpark@sinsiawy.com"}`))

	assert.NoError(err)
	assert.Equal(http.StatusCreated, res.StatusCode) // code : StatusCreated

	user := new(User)
	err = json.NewDecoder(res.Body).Decode(user)

	assert.NoError(err)
	assert.NotEqual(0, user.ID)

	id := user.ID

	res, err = http.Get(ts.URL + "/users/" + strconv.Itoa(id))
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)

	//data, _ := ioutil.ReadAll(res.Body)
	//assert.Contains("User Id: "+strconv.Itoa(id), string(data))

	user2 := new(User)
	err = json.NewDecoder(res.Body).Decode(user2)
	assert.NoError(err)
	assert.Equal(user.ID, user2.ID)
	assert.Equal(user.FirstName, user2.FirstName)
}

func TestDeleteUser(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	req, _ := http.NewRequest("DELETE", ts.URL+"/users/1", nil)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)
	data, _ := ioutil.ReadAll(res.Body)
	log.Print(string(data))
}

func TestUpdateUser(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	// craete user
	res, err := http.Post(ts.URL+"/users", "application/json", strings.NewReader(`{"first_name":"JongHan", "last_name":"Park", "email":"jhpark@sinsiawy.com"}`))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, res.StatusCode) // code : StatusCreated

	data, _ := ioutil.ReadAll(res.Body)
	log.Print(string(data))

	user := new(User)
	err = json.Unmarshal(data, user)
	assert.NoError(err)

	updateStr := fmt.Sprintf(
		`{ 
			"id":%d,
			"first_name":"updated_fn",
			"last_name":"updated_ln",
			"email":"updated_email@sinsiway.com"
		}`,
		user.ID)

	// update user
	req, _ := http.NewRequest("PUT", ts.URL+"/users", strings.NewReader(updateStr))

	res, err = http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)

	data, _ = ioutil.ReadAll(res.Body)
	log.Print(string(data))
}
