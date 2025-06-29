package controller

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/MrWhok/IMK-FP-BACKEND/model"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	//authentication
	tokenResponse := authenticationCreate()

	//service
	createProductRequest := model.ProductCreateModel{
		Name:     "Test Product",
		Price:    1000,
		Quantity: 1000,
	}
	requestBody, _ := json.Marshal(createProductRequest)

	request := httptest.NewRequest("POST", "/v1/api/product", bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Bearer "+tokenResponse["token"].(string))

	response, _ := appTest.Test(request)

	assert.Equal(t, 200, response.StatusCode)
	responseBody, _ := io.ReadAll(response.Body)

	webResponse := model.GeneralResponse{}
	_ = json.Unmarshal(responseBody, &webResponse)
	assert.Equal(t, 200, webResponse.Code)
	assert.Equal(t, "Success", webResponse.Message)

	jsonData, _ := json.Marshal(webResponse.Data)
	createProductResponse := model.ProductCreateModel{}
	_ = json.Unmarshal(jsonData, &createProductResponse)

	assert.Equal(t, createProductRequest.Name, createProductResponse.Name)
	assert.Equal(t, createProductRequest.Price, createProductResponse.Price)
	assert.Equal(t, createProductRequest.Quantity, createProductResponse.Quantity)
}
