package handlers

import (
	"errors"
	"net/http"

	database "SMS-panel/database"
	"SMS-panel/models"

	echo "github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// CreatePhoneBookNumber creates a new phone book number
func CreatePhoneBookNumber(c echo.Context) error {
	var phoneBookNumber models.PhoneBookNumber
	if err := c.Bind(&phoneBookNumber); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
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
func UpdatePhoneBookNumber(c echo.Context) error {
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

	if err := c.Bind(&phoneBookNumber); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	result = db.Save(&phoneBookNumber)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, result.Error.Error())
	}

	return c.JSON(http.StatusOK, phoneBookNumber)
}

// DeletePhoneBookNumber deletes a phone book number based on its ID
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
