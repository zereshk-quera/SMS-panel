package handlers

import (
	"errors"
	"net/http"

	"SMS-panel/models"
	"SMS-panel/utils"

	echo "github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UpdatePhoneBookNumberRequest struct {
	Prefix string `json:"prefix"`
	Name   string `json:"name"`
	Phone  string `json:"phone"`
}
type CreatePhoneBookNumberRequest struct {
	Name        string `json:"name"`
	Phone       string `json:"phone"`
	PhoneBookID uint   `json:"phoneBookID"`
	Prefix      string `json:"prefix"`
	Username    string `json:"username"`
}

// CreatePhoneBookNumber creates a new phone book number
// @Summary Create a new phone book number
// @Description Create a new phone book number
// @Tags PhoneBookNumbers
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization Token"
// @Param phoneBookNumber body CreatePhoneBookNumberRequest true "Phone book number object"
// @Success 201 {object} models.PhoneBookNumber
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /account/phone-books/phone-book-numbers [post]
func (p *PhonebookHandler) CreatePhoneBookNumber(c echo.Context) error {
	phoneBookNumber := models.PhoneBookNumber{}
	account := c.Get("account").(models.Account)

	if err := c.Bind(&phoneBookNumber); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if phoneBookNumber.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Name is required"})
	}

	if phoneBookNumber.Phone == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Phone is required"})
	}

	if phoneBookNumber.PhoneBookID == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Phone book is required"})
	}

	// Retrieve the existing PhoneBook with the given PhoneBookID
	phoneBook := models.PhoneBook{}
	result := p.db.First(&phoneBook, phoneBookNumber.PhoneBookID)
	if result.Error != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Phone book not found"})
	}

	// Set the AccountID and assign the existing PhoneBook to PhoneBookNumber
	phoneBookNumber.PhoneBook = phoneBook
	phoneBookNumber.PhoneBook.AccountID = account.ID

	// Check Phone Number Validation
	if !utils.ValidatePhone(phoneBookNumber.Phone) {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "Invalid Phone Number"})
	}

	// Is Input Phone Number Unique or Not
	var existingPhoneBookNumber models.PhoneBookNumber
	p.db.Where("phone = ? AND prefix = ?", phoneBookNumber.Phone, phoneBookNumber.Prefix).First(&existingPhoneBookNumber)
	if existingPhoneBookNumber.ID != 0 {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "Input Phone Number has already been registered"})
	}

	p.db.Where("username = ?", phoneBookNumber.Username).First(&existingPhoneBookNumber)
	if existingPhoneBookNumber.ID != 0 {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "Input Username has already been registered"})
	}

	result = p.db.Create(&phoneBookNumber)
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
// @Param Authorization header string true "Authorization Token"
// @Success 200 {array} models.PhoneBookNumber
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /account/phone-books/{phoneBookID}/phone-book-numbers [get]
func (p *PhonebookHandler) ListPhoneBookNumbers(c echo.Context) error {
	phoneBookID := c.Param("phoneBookID")

	account := c.Get("account").(models.Account)

	var phoneBookNumber []models.PhoneBookNumber
	result := p.db.
		Joins("JOIN phone_books ON phone_books.id = phone_book_numbers.phone_book_id AND phone_books.account_id = ?", account.ID).
		Where("phone_book_id = ?", phoneBookID).
		Find(&phoneBookNumber)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, "Phonebook not found")
		}
		return c.JSON(http.StatusInternalServerError, result.Error.Error())
	}

	return c.JSON(http.StatusOK, phoneBookNumber)
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
func (p *PhonebookHandler) ReadPhoneBookNumber(c echo.Context) error {
	phoneBookNumberID := c.Param("phoneBookNumberID")

	var phoneBookNumber models.PhoneBookNumber
	result := p.db.First(&phoneBookNumber, phoneBookNumberID)
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
func (p *PhonebookHandler) UpdatePhoneBookNumber(c echo.Context) error {
	phoneBookNumberID := c.Param("phoneBookNumberID")

	var existingPhoneBookNumber models.PhoneBookNumber
	result := p.db.First(&existingPhoneBookNumber, phoneBookNumberID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, "Phone book number not found")
		}
		return c.JSON(http.StatusInternalServerError, result.Error.Error())
	}

	var updatedPhoneBookNumber UpdatePhoneBookNumberRequest
	if err := c.Bind(&updatedPhoneBookNumber); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Update the fields of the existing phone book number
	if updatedPhoneBookNumber.Prefix != "" {
		existingPhoneBookNumber.Prefix = updatedPhoneBookNumber.Prefix
	}
	if updatedPhoneBookNumber.Name != "" {
		existingPhoneBookNumber.Name = updatedPhoneBookNumber.Name
	}
	if updatedPhoneBookNumber.Phone != "" {
		existingPhoneBookNumber.Phone = updatedPhoneBookNumber.Phone
	}

	// Use the `clause.OnConflict` to avoid updating the primary key
	result = p.db.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Save(&existingPhoneBookNumber)
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
func (p *PhonebookHandler) DeletePhoneBookNumber(c echo.Context) error {
	phoneBookNumberID := c.Param("phoneBookNumberID")

	var phoneBookNumber models.PhoneBookNumber
	result := p.db.First(&phoneBookNumber, phoneBookNumberID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, "Phone book number not found")
		}
		return c.JSON(http.StatusInternalServerError, result.Error.Error())
	}

	result = p.db.Delete(&phoneBookNumber)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, result.Error.Error())
	}

	return c.JSON(http.StatusOK, "Phone book number deleted")
}
