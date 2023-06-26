package utils

import (
	"context"
	"errors"
	"net/mail"
	"strconv"
	"strings"
	"time"

	"SMS-panel/models"

	"gorm.io/gorm"
)

// This Function Validates Input Email.
func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// This Function Validates Input Phone Number.
func ValidatePhone(phone string) bool {
	hasCharacter := false
	for _, digit := range phone {
		if digit < 48 || digit > 57 {
			hasCharacter = true
			break
		}
	}
	if hasCharacter {
		return false
	}
	return strings.HasPrefix(phone, "09") && len(phone) == 11
}

// This Function Parse Input String to Integer on Input Base.
func ParseInt(s string, base int) int {
	n, err := strconv.ParseInt(s, base, 64)
	if err != nil {
		return 0
	}
	return int(n)
}

// This Function Validates Input National ID.
func ValidateNationalID(id string) bool {
	l := len(id)

	if l < 8 || ParseInt(id, 10) == 0 {
		return false
	}

	id = ("0000" + id)[l+4-10:]
	if ParseInt(id[3:9], 10) == 0 {
		return false
	}

	c := ParseInt(id[9:10], 10)
	s := 0
	for i := 0; i < 9; i++ {
		s += ParseInt(id[i:i+1], 10) * (10 - i)
	}
	s = s % 11

	return (s < 2 && c == s) || (s >= 2 && c == (11-s))
}

func ValidateTimeFormat(timeStr string) bool {
	_, err := time.Parse("15:04", timeStr)
	return err == nil
}

func ValidateJsonFormat(jsonBody map[string]interface{}, fields ...string) (string, error) {
	msg := "OK"
	for _, field := range fields {
		if _, ok := jsonBody[field]; !ok {
			msg = "Input Json doesn't include " + field
			break
		}
	}
	if msg != "OK" {
		return msg, errors.New("")
	}
	return msg, nil
}

// Check if sender number is available
func IsSenderNumberExist(
	ctx context.Context,
	db *gorm.DB,
	senderNumber string,
	userId uint,
) bool {
	var senderNumbersObjects models.SenderNumber

	err := db.WithContext(ctx).Model(&models.SenderNumber{}).
		Select("sender_numbers.number").
		Joins("LEFT JOIN user_numbers ON sender_numbers.id = user_numbers.number_id").
		Where(
			"(sender_numbers.is_default=true or (user_numbers.user_id = ? and user_numbers.is_available=true)) and sender_numbers.number = ?",
			userId, senderNumber).
		First(&senderNumbersObjects).Error

	return err == nil
}
