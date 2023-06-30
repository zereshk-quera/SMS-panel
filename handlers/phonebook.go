package handlers

import (
	"errors"
	"net/http"

	"SMS-panel/models"

	echo "github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type PhoneBookRequest struct {
	AccountID uint   `json:"accountID" binding:"required"`
	Name      string `json:"name" binding:"required"`
}

type PhoneBookResponse struct {
	ID        uint   `json:"id"`
	AccountID uint   `json:"accountID"`
	Name      string `json:"name"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type PhonebookHandler struct {
	db *gorm.DB
}

func NewPhonebookHandler(db *gorm.DB) *PhonebookHandler {
	return &PhonebookHandler{db: db}
}

// @Summary Create a phone book entry
// @Description Create a new phone book entry
// @Tags phonebook
// @Accept json
// @Produce json
// @Param phoneBook body PhoneBookRequest true "Phone book entry data"
// @Success 201 {object} PhoneBookResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /account/phone-books/ [post]
func (p *PhonebookHandler) CreatePhoneBook(c echo.Context) error {
	var phoneBook models.PhoneBook
	account := c.Get("account").(models.Account)

	if err := c.Bind(&phoneBook); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if phoneBook.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Name is required"})
	}

	if phoneBook.AccountID != account.ID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Permision Denied"})
	}

	if err := p.db.Create(&phoneBook).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create phone book"})
	}

	return c.JSON(http.StatusCreated, phoneBook)
}

// @Summary Get all phone books
// @Description Get all phone books for a given account ID
// @Tags phonebook
// @Accept json
// @Produce json
// @Success 200 {array} PhoneBookResponse
// @Failure 500 {object} ErrorResponse
// @Router /account/phone-books/ [get]
func (p *PhonebookHandler) GetAllPhoneBooks(c echo.Context) error {
	account := c.Get("account").(models.Account)

	var phoneBooks []models.PhoneBook
	// Get all matched records
	result := p.db.Where("account_id = ?", account.ID).Find(&phoneBooks)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, result.Error.Error())
	}

	return c.JSON(http.StatusOK, phoneBooks)
}

// @Summary Get a phone book
// @Description Get a phone book by ID for a given account ID
// @Tags phonebook
// @Accept json
// @Produce json
// @Param phoneBookID path int true "Phone Book ID"
// @Success 200 {object} PhoneBookResponse
// @Failure 404 {string} string
// @Failure 500 {object} ErrorResponse
// @Router /account/phone-books/{phoneBookID} [get]
func (p *PhonebookHandler) ReadPhoneBook(c echo.Context) error {
	phoneBookID := c.Param("phoneBookID")
	account := c.Get("account").(models.Account)

	var phoneBook models.PhoneBook
	// Find the phone book with matching phoneBookID and accountID
	result := p.db.Where("id = ? AND account_id = ?", phoneBookID, account.ID).First(&phoneBook)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, "Phonebook not found")
		}
		return c.JSON(http.StatusInternalServerError, result.Error.Error())
	}

	return c.JSON(http.StatusOK, phoneBook)
}

// @Summary Update a phone book
// @Description Update a phone book by ID for a given account ID
// @Tags phonebook
// @Accept json
// @Produce json
// @Param phoneBookID path int true "Phone Book ID"
// @Param phoneBook body PhoneBookRequest true "Phone Book object"
// @Success 200 {object} PhoneBookResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {string} string
// @Failure 500 {object} ErrorResponse
// @Router /account/{accountID}/phone-books/{phoneBookID} [put]
func (p *PhonebookHandler) UpdatePhoneBook(c echo.Context) error {
	phoneBookID := c.Param("phoneBookID")
	account := c.Get("account").(models.Account)

	var phoneBook models.PhoneBook

	result := p.db.Where("id = ? AND account_id = ?", phoneBookID, account.ID).First(&phoneBook)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, "Phonebook not found")
		}
		return c.JSON(http.StatusInternalServerError, result.Error.Error())
	}

	if err := c.Bind(&phoneBook); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	result = p.db.Save(&phoneBook)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, result.Error.Error())
	}

	return c.JSON(http.StatusOK, phoneBook)
}

// @Summary Delete a phone book
// @Description Delete a phone book by ID for a given account ID
// @Tags phonebook
// @Accept json
// @Produce json
// @Param phoneBookID path int true "Phone Book ID"
// @Success 200 {string} string
// @Failure 404 {string} string
// @Failure 500 {object} ErrorResponse
// @Router /account/phone-books/{phoneBookID} [delete]
func (p *PhonebookHandler) DeletePhoneBook(c echo.Context) error {
	phoneBookID := c.Param("phoneBookID")
	account := c.Get("account").(models.Account)

	var phoneBook models.PhoneBook
	result := p.db.Where("id = ? AND account_id = ?", phoneBookID, account.ID).First(&phoneBook)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, "Phone book not found")
		}
		return c.JSON(http.StatusInternalServerError, result.Error.Error())
	}

	result = p.db.Delete(&phoneBook)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, result.Error.Error())
	}

	return c.JSON(http.StatusOK, "Phone book deleted")
}
