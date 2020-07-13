package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"guru/models"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestUserHandler_Get(t *testing.T) {
	jsonStr := []byte(`{
 		"id": 1,
        "token": "sssss"
	}`)

	res, err := http.Post(
		fmt.Sprintf("%s/user/get", srv.URL),
		"application/json",
		bytes.NewBuffer(jsonStr),
	)

	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	var getUserResponse models.GetUserResponseModel
	if err = json.Unmarshal(resBytes, &getUserResponse); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, models.GetUserResponseModel{
		Id:           1,
		Balance:      75,
		DepositCount: 2,
		DepositSum:   200,
		BetCount:     1,
		BetSum:       50,
		WinCount:     1,
		WinSum:       25,
	}, getUserResponse)
}

func TestUserHandler_GetWrongToken(t *testing.T) {
	jsonStr := []byte(`{
 		"id": 1,
        "token": "ttttt"
	}`)

	res, err := http.Post(
		fmt.Sprintf("%s/user/get", srv.URL),
		"application/json",
		bytes.NewBuffer(jsonStr),
	)

	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	var errorResponse models.ErrorResponseModel
	if err = json.Unmarshal(resBytes, &errorResponse); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Equal(t, models.ErrorResponseModel{Error: "wrong token"}, errorResponse)
}

func TestUserHandler_GetNotFound(t *testing.T) {
	jsonStr := []byte(`{
 		"id": 5,
        "token": "sssss"
	}`)

	res, err := http.Post(
		fmt.Sprintf("%s/user/get", srv.URL),
		"application/json",
		bytes.NewBuffer(jsonStr),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	var errorResponse models.ErrorResponseModel
	if err = json.Unmarshal(resBytes, &errorResponse); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	assert.Equal(t, models.ErrorResponseModel{Error: "not found"}, errorResponse)
}

func TestUserHandler_Create(t *testing.T) {
	jsonStr := []byte(`{
		"id": 3,
		"balance": 75,
		"token": "string"
	}`)

	res, err := http.Post(
		fmt.Sprintf("%s/user/create", srv.URL),
		"application/json",
		bytes.NewBuffer(jsonStr),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	var errorResponse models.ErrorResponseModel
	if err = json.Unmarshal(resBytes, &errorResponse); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, models.ErrorResponseModel{Error: ""}, errorResponse)
}

func TestUserHandler_AddDeposit(t *testing.T) {
	jsonStr := []byte(`{
		"user_id": 3,
		"deposit_id": 1,
		"amount": 50,
		"token": "string"
	}`)

	res, err := http.Post(
		fmt.Sprintf("%s/user/deposit", srv.URL),
		"application/json",
		bytes.NewBuffer(jsonStr),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	var transactionResponse models.TransactionResponseModel
	if err = json.Unmarshal(resBytes, &transactionResponse); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, models.TransactionResponseModel{Error: "", Balance: 125}, transactionResponse)
}

func TestUserHandler_AddDepositWrongToken(t *testing.T) {
	jsonStr := []byte(`{
		"user_id": 3,
		"deposit_id": 1,
		"amount": 50,
		"token": "ttttt"
	}`)

	res, err := http.Post(
		fmt.Sprintf("%s/user/deposit", srv.URL),
		"application/json",
		bytes.NewBuffer(jsonStr),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	var errorResponse models.ErrorResponseModel
	if err = json.Unmarshal(resBytes, &errorResponse); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Equal(t, models.ErrorResponseModel{Error: "wrong token"}, errorResponse)
}

func TestUserHandler_AddDepositNotFound(t *testing.T) {
	jsonStr := []byte(`{
		"user_id": 5,
		"deposit_id": 1,
		"amount": 50,
		"token": "string"
	}`)

	res, err := http.Post(
		fmt.Sprintf("%s/user/deposit", srv.URL),
		"application/json",
		bytes.NewBuffer(jsonStr),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	var errorResponse models.ErrorResponseModel
	if err = json.Unmarshal(resBytes, &errorResponse); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	assert.Equal(t, models.ErrorResponseModel{Error: "not found"}, errorResponse)
}
