package handlers

import (
	database "SMS-panel/database"
	"SMS-panel/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	merchantID      = "860C78FA-D6A9-48AE-805D-2B33B52309D2"
	callbackURL     = "http://localhost:8080/accounts/payment/verify"
	zarinpalRequest = "https://sandbox.banktest.ir/zarinpal/api.zarinpal.com/pg/v4/payment/request.json"
	zarinpalVerify  = "https://sandbox.banktest.ir/zarinpal/api.zarinpal.com/pg/v4/payment/verify.json"
	zarinpalGateURL = "https://sandbox.banktest.ir/zarinpal/www.zarinpal.com/pg/StartPay/"
)

type ZarinpalData struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Authority string `json:"authority"`
	FeeType   string `json:"fee_type"`
	Fee       int    `json:"fee"`
}

type ZarinpalResponse struct {
	Data   ZarinpalData  `json:"data"`
	Errors []interface{} `json:"errors"`
}

type ZarinpalVerify struct {
	Data   ZarinpalVerifyData `json:"data"`
	Errors ZarinpalErrors     `json:"errors"`
}

type ZarinpalVerifyData struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	CardHash string `json:"card_hash"`
	CardPan  string `json:"card_pan"`
	RefID    int    `json:"ref_id"`
	FeeType  string `json:"fee_type"`
	Fee      int    `json:"fee"`
}

type ZarinpalErrors struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type AmountFee struct {
	Fee int `json:"fee"`
}

type VerifyResponse struct {
	Status    string `json:"Status"`
	Authority string `json:"Authority"`
}

type RequestResponse struct {
	PaymentUrl string `json:"payment_url"`
}

// @Summary Add budget request
// @Description Zarinpal Payment to add budget to account
// @Tags payment
// @Accept json
// @Produce json
// @Param body body AmountFee true "Payment request details"
// @Success 200 {object} RequestResponse
// @Failure 400 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /accounts/payment/request [post]
func PaymentRequestHandler(c echo.Context) error {
	// Connect To The Datebase
	db, err := database.GetConnection()
	if err != nil {
		return c.JSON(http.StatusBadGateway, models.Response{ResponseCode: 502, Message: "Can't Connect To Database"})
	}

	// Read Request Body
	jsonBody := make(map[string]int)
	err = json.NewDecoder(c.Request().Body).Decode(&jsonBody)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Invalid JSON"})
	}

	if _, ok := jsonBody["fee"]; !ok {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Input Json doesn't include fee"})
	}

	account := c.Get("account")
	userID := account.(models.Account).UserID

	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		// Handle the error (e.g., user not found)
		return c.JSON(http.StatusNotFound, models.Response{ResponseCode: 404, Message: "User Not Founded"})
	}

	data := map[string]interface{}{
		"merchant_id":  merchantID,
		"amount":       jsonBody["fee"],
		"callback_url": callbackURL,
		"description":  "Add budget to account",
		"metadata": map[string]string{
			"mobile": user.Phone,
			"email":  user.Email,
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			ResponseCode: 500,
			Message:      "Failed to marshal JSON data",
		})
	}

	resp, err := http.Post(zarinpalRequest, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return c.JSON(http.StatusBadGateway, models.Response{
			ResponseCode: 502,
			Message:      "Failed to send POST request",
		})
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			ResponseCode: 500,
			Message:      "Failed to read body",
		})
	}

	var result ZarinpalResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			ResponseCode: 500,
			Message:      "Failed to parse response",
		})
	}

	var transaction models.Transaction
	transaction.AccountID = user.ID
	transaction.Amount = int64(jsonBody["fee"])
	transaction.Status = "Wait"
	transaction.Authority = result.Data.Authority
	transaction.CreatedAt = time.Now()

	// Insert Transaction Object Into Database
	createdTransaction := db.Create(&transaction)
	if createdTransaction.Error != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{ResponseCode: 500, Message: "Transaction Cration Failed"})
	}
	var response RequestResponse
	response.PaymentUrl = fmt.Sprintf("%s%s", zarinpalGateURL, result.Data.Authority)

	return c.JSON(http.StatusOK, response)
}

// @Summary Verify budget payment and add budget
// @Description Verify Zarinpal Payment to add budget to account
// @Tags payment
// @Accept json
// @Produce json
// @Param body body VerifyResponse true "Payment verify details"
// @Success 200 {string} string
// @Failure 400 {string} ErrorResponse
// @Failure 422 {string} ErrorResponse
// @Failure 500 {string} ErrorResponse
// @Router /accounts/payment/verify [get]
func PaymentVerifyHandler(c echo.Context) error {
	// Connect To The Datebase
	db, err := database.GetConnection()
	if err != nil {
		return c.JSON(http.StatusBadGateway, models.Response{ResponseCode: 502, Message: "Can't Connect To Database"})
	}

	authority := c.QueryParam("Authority")
	status := c.QueryParam("Status")

	var transaction models.Transaction
	if err := db.Where(&models.Transaction{Authority: authority}).First(&transaction).Error; err != nil {
		// Handle the error (e.g., transaction not found)
		return c.JSON(http.StatusNotFound, models.Response{ResponseCode: 404, Message: "Transaction Not Founded"})
	}

	if status == "NOK" {
		transaction.Status = "Failed"
		db.Save(&transaction)
		return c.JSON(http.StatusBadRequest, "Failed Payment")
	}

	data := map[string]interface{}{
		"merchant_id": merchantID,
		"amount":      transaction.Amount,
		"authority":   transaction.Authority,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			ResponseCode: 500,
			Message:      "Failed to marshal JSON data",
		})
	}

	resp, err := http.Post(zarinpalVerify, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return c.JSON(http.StatusBadGateway, models.Response{
			ResponseCode: 502,
			Message:      "Failed to send POST request",
		})
	}
	defer resp.Body.Close()

	// Read Request Body
	jsonBody := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&jsonBody)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, models.Response{ResponseCode: 422, Message: "Invalid JSON"})
	}

	if data, ok := jsonBody["data"]; ok {
		if dataMap, ok := data.(map[string]interface{}); ok {
			if code, ok := dataMap["code"]; ok {
				if code == float64(100) {
					transaction.Status = "Okay"
					db.Save(&transaction)

					accountID := transaction.AccountID
					var account models.Account
					if err := db.First(&account, accountID).Error; err != nil {
						// Handle the error (e.g., account not found)
						return c.JSON(http.StatusNotFound, models.Response{ResponseCode: 404, Message: "Account Not Founded"})
					}

					account.Budget += transaction.Amount
					db.Save(&account)

					return c.JSON(http.StatusOK, "Successful Payment")
				} else if code == float64(101) {
					return c.JSON(http.StatusOK, "Transaction had verified")
				}
			}
		}

	}

	transaction.Status = "Failed"
	db.Save(&transaction)
	return c.JSON(http.StatusBadRequest, "Failed Payment")
}
