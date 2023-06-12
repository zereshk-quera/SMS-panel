package handlers

import (
	"errors"
	"net/http"

	database "SMS-panel/database"
	"SMS-panel/models"

	echo "github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type CreatePhoneBookNumberRequest struct {
	PhoneBookID uint   `json:"phoneBookID"`
	Prefix      string `json:"prefix"`
	Name        string `json:"name"`
	Phone       string `json:"phone"`
}
type UpdatePhoneBookNumberRequest struct {
	Prefix string `json:"prefix"`
	Name   string `json:"name"`
	Phone  string `json:"phone"`
}

// CreatePhoneBookNumber creates a new phone book number
// @Summary Create a new phone book number
// @Description Create a new phone book number
// @Tags PhoneBookNumbers
// @Accept json
// @Produce json
// @Param phoneBookNumber body CreatePhoneBookNumberRequest true "Phone book number object"
// @Success 201 {object} models.PhoneBookNumber
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /account/phone-books/phone-book-numbers [post]
func CreatePhoneBookNumber(c echo.Context) error {
	var request CreatePhoneBookNumberRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	phoneBookNumber := models.PhoneBookNumber{
		PhoneBookID: request.PhoneBookID,
		Prefix:      request.Prefix,
		Name:        request.Name,
		Phone:       request.Phone,
	}

	db, err := database.GetConnection()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	result := db.Create(&phoneBookNumber)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, result.Error.Error())
	}

	return c.JSON(http.StatusCreated, phoneBookNumber)
}

// ListPhoneBookNumbers retrieves all phone book numbers for a given PhoneBookID
// @Summary Get all phone book numbers for a given PhoneBookID
// @Description Get all phone book numbers for a given PhoneBookID
// @Tags PhoneBookNumbers
// @Accept json
// @Produce json
// @Param phoneBookID path string true "Phone book ID"
// @Success 200 {array} models.PhoneBookNumber
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /account/phone-books/{phoneBookID}/phone-book-numbers [get]
func ListPhoneBookNumbers(c echo.Context) error {
	phoneBookID := c.Param("phoneBookID")

	db, err := database.GetConnection()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	var phoneBookNumbers []models.PhoneBookNumber
	result := db.Where("phone_book_id = ?", phoneBookID).Find(&phoneBookNumbers)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, "Phone book not found")
		}
		return c.JSON(http.StatusInternalServerError, result.Error.Error())
	}

	return c.JSON(http.StatusOK, phoneBookNumbers)
}

// ReadPhoneBookNumber retrieves the data of a phone book number based on its ID
// @Summary Get phone book number by ID
// @Description Get phone book number by ID
// @Tags PhoneBookNumbers
// @Accept json
// @Produce json
// @Param phoneBookNumberID path string true "Phone book number ID"
// @Success 200 {object} models.PhoneBookNumber
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /account/phone-books/phone-book-numbers/{phoneBookNumberID} [get]
func ReadPhoneBookNumber(c echo.Context) error {
	phoneBookNumberID := c.Param("phoneBookNumberID")

	db, err := database.GetConnection()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	var phoneBookNumber models.PhoneBookNumber
	result := db.First(&phoneBookNumber, phoneBookNumberID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, "Phone book number not found")
		}
		return c.JSON(http.StatusInternalServerError, result.Error.Error())
	}

	return c.JSON(http.StatusOK, phoneBookNumber)
}

// UpdatePhoneBookNumber updates a phone book number based on its ID
// @Summary Update phone book number
// @Description Update phone book number
// @Tags PhoneBookNumbers
// @Accept json
// @Produce json
// @Param phoneBookNumberID path string true "Phone book number ID"
// @Param phoneBookNumber body UpdatePhoneBookNumberRequest true "Phone book number object"
// @Success 200 {object} models.PhoneBookNumber
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /account/phone-books/phone-book-numbers/{phoneBookNumberID} [put]
func UpdatePhoneBookNumber(c echo.Context) error {
	phoneBookNumberID := c.Param("phoneBookNumberID")

	db, err := database.GetConnection()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	var existingPhoneBookNumber models.PhoneBookNumber
	result := db.First(&existingPhoneBookNumber, phoneBookNumberID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, "Phone book number not found")
		}
		return c.JSON(http.StatusInternalServerError, result.Error.Error())
	}

	var updatedPhoneBookNumber models.PhoneBookNumber
	if err := c.Bind(&updatedPhoneBookNumber); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Update the fields of the existing phone book number
	existingPhoneBookNumber.Prefix = updatedPhoneBookNumber.Prefix
	existingPhoneBookNumber.Name = updatedPhoneBookNumber.Name
	existingPhoneBookNumber.Phone = updatedPhoneBookNumber.Phone

	result = db.Save(&existingPhoneBookNumber)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, result.Error.Error())
	}

	return c.JSON(http.StatusOK, existingPhoneBookNumber)
}

// DeletePhoneBookNumber deletes a phone book number based on its ID
// @Summary Delete phone book number
// @Description Delete phone book number
// @Tags PhoneBookNumbers
// @Accept json
// @Produce json
// @Param phoneBookNumberID path string true "Phone book number ID"
// @Success 200 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /account/phone-books/phone-book-numbers/{phoneBookNumberID} [delete]
func DeletePhoneBookNumber(c echo.Context) error {
	phoneBookNumberID := c.Param("phoneBookNumberID")

	db, err := database.GetConnection()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	var phoneBookNumber models.PhoneBookNumber
	result := db.First(&phoneBookNumber, phoneBookNumberID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, "Phone book number not found")
		}
		return c.JSON(http.StatusInternalServerError, result.Error.Error())
	}

	result = db.Delete(&phoneBookNumber)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, result.Error.Error())
	}

	return c.JSON(http.StatusOK, "Phone book number deleted")
}
