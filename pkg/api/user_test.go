package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var theApp App

func TestMain(m *testing.M) {
	// setup
	theApp = App{}
	theApp.Initialize("microgram", "yourpassword", "microgram")
	ensureTableExists()

	code := m.Run()
	// shutdown

	os.Exit(code)
}

const tableCreationQuery = `
create table if not exists photographers
(
	email VARCHAR(50) UNIQUE PRIMARY KEY,
	username VARCHAR(50) UNIQUE NOT NULL,
	phone VARCHAR(15) UNIQUE NOT NULL,
	firstname VARCHAR(50) NULL,
	lastname VARCHAR(50) NULL,
	city	VARCHAR(50) NULL,
	gender	ENUM('m', 'f'),
	birthdate DATE NULL
)
`

func ensureTableExists() {
	if _, err := theApp.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	theApp.DB.Exec("DELETE FROM photographers")
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	resRecorder := httptest.NewRecorder()
	theApp.Router.ServeHTTP(resRecorder, req)

	return resRecorder
}

func checkResponseCode(t *testing.T, expectedCode, actualCode int) {
	if expectedCode != actualCode {
		t.Errorf("Expected response code %d. Got %d\n", expectedCode, actualCode)
	}
}

func TestCreateUser(t *testing.T) {
	clearTable()
	payload := []byte(`{
		"email":"jojorancu@gmail.com",
		"phone":"0812345678",
		"username":"jojorancu",
		"gender":"m"
	}`)

	req, _ := http.NewRequest("POST", "/photographer", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var thePhotographer map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &thePhotographer)

	if thePhotographer["error"] != nil {
		t.Errorf("Caused error by: '%v'", thePhotographer["error"])
	} else {
		if thePhotographer["email"] != "jojorancu@gmail.com" {
			t.Errorf("Expected user email to be 'jojorancu@gmail.com'. Got '%v'", thePhotographer["email"])
		}

		if thePhotographer["phone"] != "0812345678" {
			t.Errorf("Expected user phone to be '0812345678'. Got '%v'", thePhotographer["phone"])
		}

		if thePhotographer["username"] != "jojorancu" {
			t.Errorf("Expected user phone to be 'jojorancu'. Got '%v'", thePhotographer["username"])
		}

		if thePhotographer["gender"] != "m" {
			t.Errorf("Expected user gender to be 'm'. Got '%v'", thePhotographer["gender"])
		}
	}
}

func TestGetUser(t *testing.T) {
	req, _ := http.NewRequest("GET", "/photographer/jojorancu", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var thePhotographer map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &thePhotographer)

	if thePhotographer["error"] != nil {
		t.Errorf("Caused error by: '%v'", thePhotographer["error"])
	} else {
		if thePhotographer["email"] != "jojorancu@gmail.com" {
			t.Errorf("Expected user email to be 'jojorancu@gmail.com'. Got '%v'", thePhotographer["email"])
		}

		if thePhotographer["phone"] != "0812345678" {
			t.Errorf("Expected user phone to be '0812345678'. Got '%v'", thePhotographer["phone"])
		}

		if thePhotographer["username"] != "jojorancu" {
			t.Errorf("Expected user phone to be 'jojorancu'. Got '%v'", thePhotographer["username"])
		}

		if thePhotographer["gender"] != "m" {
			t.Errorf("Expected user gender to be 'm'. Got '%v'", thePhotographer["gender"])
		}
	}
}

func TestGetNonExistentUser(t *testing.T) {
	req, _ := http.NewRequest("GET", "/photographer/jojorancuu", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var thePhotographer map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &thePhotographer)

	if thePhotographer["error"] != "User not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'User not found'. Got '%s'", thePhotographer["error"])
	}
}

func TestGetNonExistentUser_Regex(t *testing.T) {
	req, _ := http.NewRequest("GET", "/photographer/jojorancu!!!", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)

	var thePhotographer map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &thePhotographer)

	if thePhotographer["error"] != "User should not be exist" {
		t.Errorf("Expected the 'error' key of the response to be set to 'User should not be exist. Got '%s'", thePhotographer["error"])
	}
}

func TestUpdateUser(t *testing.T) {
	req, _ := http.NewRequest("GET", "/photographer/jojorancu", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var recordPhotographer map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &recordPhotographer)

	payload := []byte(`{
		"email":"secretmethod94@gmail.com",
		"phone":"08123456789",
		"gender":"m"
	}`)
	req, _ = http.NewRequest("PUT", "/photographer/jojorancu", bytes.NewBuffer(payload))
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var thePhotographer map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &thePhotographer)

	if thePhotographer["error"] != nil {
		t.Errorf("Caused error by: '%v'", thePhotographer["error"])
	} else {
		if thePhotographer["username"] != recordPhotographer["username"] {
			t.Errorf("Expected the username to remain the same '%v'. Got '%v'", recordPhotographer["username"], thePhotographer["username"])
		}
		if thePhotographer["email"] == recordPhotographer["email"] {
			t.Errorf("Expected the email to change from '%v' to '%v'. Got '%v'", recordPhotographer["email"], thePhotographer["email"], thePhotographer["email"])
		}
		if thePhotographer["phone"] == recordPhotographer["phone"] {
			t.Errorf("Expected the phone to change from '%v' to '%v'. Got '%v'", recordPhotographer["phone"], thePhotographer["phone"], thePhotographer["phone"])
		}
		if thePhotographer["gender"] != recordPhotographer["gender"] {
			t.Errorf("Expected the gender to remain the same '%v'. Got '%v'", recordPhotographer["gender"], thePhotographer["gender"])
		}
	}
}

func TestDeleteUser(t *testing.T) {
	req, _ := http.NewRequest("GET", "/photographer/jojorancu", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/photographer/jojorancu", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/photographer/jojorancu", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}
