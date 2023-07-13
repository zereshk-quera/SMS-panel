package handlers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strings"
	"time"

	database "SMS-panel/database"
	"SMS-panel/models"
	"SMS-panel/utils"

	"gorm.io/gorm"
)

// Declare errors
type PhoneBooksNotFoundError struct {
	Message string
}

func (e PhoneBooksNotFoundError) Error() string {
	return fmt.Sprintf(e.Message)
}

type AcountDoesNotHaveBudgetError struct {
	Message string
}

func (e AcountDoesNotHaveBudgetError) Error() string {
	return fmt.Sprintf(e.Message)
}

type PhoneBooksNumbersAreEmptyError struct {
	Message string
}

func (e PhoneBooksNumbersAreEmptyError) Error() string {
	return fmt.Sprintf(e.Message)
}

type SenderNumberNotFoundError struct {
	Message string
}

func (e SenderNumberNotFoundError) Error() string {
	return fmt.Sprintf(e.Message)
}

type SendMessageStatus struct {
	ID     int
	Status bool
}

type Message struct {
	Text        string `json:"text"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
}

func MockSendMessage(message *Message) (string, error) {
	// Query the database to retrieve bad words
	var badWords []models.Bad_Word
	db, err := database.GetConnection()
	if err != nil {
		return "failed to connect db", err
	}
	db.Find(&badWords)

	for _, badWord := range badWords {
		regex := regexp.MustCompile(badWord.Regex)
		if regex.MatchString(message.Text) {
			return "message not sent", errors.New("Message contains a bad word")
		}
	}

	log.Printf("Message sent - Text: %s, Source: %s, Destination: %s", message.Text, message.Source, message.Destination)

	return "Message sent successfully", nil
}

func CreateSMSTemplate(template string, phoneNumber models.PhoneBookNumber) string {
	template = strings.ReplaceAll(template, "%name", phoneNumber.Name)
	template = strings.ReplaceAll(template, "%prefix", phoneNumber.Prefix)

	currentTime := time.Now().Format("2006-01-02 15:04:05")
	template = strings.ReplaceAll(template, "%date", currentTime)
	template = strings.ReplaceAll(template, "%username", phoneNumber.Username)

	return template
}

func SendMessageToPhoneBooks(
	ctx context.Context,
	body SendSMessageToPhoneBooksBody,
	db *gorm.DB,
) error {
	// Check if all phone books exist
	phoneBooksExist, err := CheckPhoneBooksExist(ctx, db, body.PhoneBooks)
	if err != nil {
		return err
	}
	if !phoneBooksExist {
		return PhoneBooksNotFoundError{Message: "Phone book not found!"}
	}

	// Check if sender number is available
	senderNumberExisted := utils.IsSenderNumberExist(
		ctx, db, body.SenderNumber, body.Account.UserID,
	)
	if !senderNumberExisted {
		return SenderNumberNotFoundError{Message: "Sender number not found!"}
	}

	// Fetch all phone book numbers
	var phoneBookNumbers []models.PhoneBookNumber
	err = db.WithContext(ctx).Joins(
		"JOIN phone_books ON phone_books.id = phone_book_numbers.phone_book_id").
		Where("phone_books.account_id = ? AND phone_books.name IN ?", body.Account.ID, body.PhoneBooks).
		Find(&phoneBookNumbers).Error
	if err != nil {
		fmt.Println(err, "========")
		return PhoneBooksNumbersAreEmptyError{Message: "Phone books are empty!"}
	} else if len(phoneBookNumbers) == 0 {
		return PhoneBooksNumbersAreEmptyError{Message: "Phone books are empty!"}
	}

	fmt.Println(phoneBookNumbers)

	// Check that if user have enough budget
	var smsCost int
	err = db.WithContext(ctx).Table("configuration").
		Where("name = (?)", "group sms").
		Pluck("value", &smsCost).Error

	if err != nil {
		log.Printf("Failed to retrieve SMS costs: %s", err.Error())
		return err
	}

	cost := int64(smsCost * len(phoneBookNumbers))
	haveAccountBudget := utils.DoesAcountHaveBudget(
		body.Account.Budget, cost,
	)
	if !haveAccountBudget {
		return AcountDoesNotHaveBudgetError{Message: "You don't have enough budget!"}
	}

	// send message
	statusOfMessages := make(chan SendMessageStatus, len(phoneBookNumbers))
	for messageID, phoneNumber := range phoneBookNumbers {
		cost := smsCost
		if !utils.DoesAcountHaveBudget(body.Account.Budget, int64(cost)) {
			return AcountDoesNotHaveBudgetError{Message: "You don't have enough budget!"}
		}
		message := CreateSMSTemplate(body.Message, phoneNumber)
		go SendGroupMessage(
			statusOfMessages, message, messageID, body.Account, phoneNumber,
		)
	}

	// check status of sms
	for i := 0; i < len(phoneBookNumbers); i++ {
		messageStatus := <-statusOfMessages
		phoneNumber := phoneBookNumbers[messageStatus.ID]
		message := CreateSMSTemplate(body.Message, phoneNumber)
		sms := models.SMSMessage{
			Sender:    body.SenderNumber,
			Recipient: phoneNumber.Phone,
			Message:   message,
			Schedule:  nil,
			CreatedAt: time.Now(),
			AccountID: body.Account.ID,
		}

		if messageStatus.Status {
			sms.DeliveryReport = "Message sent successfully"

			body.Account.Budget -= int64(smsCost)
			if err := db.Save(&body.Account).Error; err != nil {
				log.Printf("Failed to update account's budget: %s", err.Error())
				return err
			}
		} else {
			sms.DeliveryReport = "Message sent field"
		}

		if err := db.Create(&sms).Error; err != nil {
			log.Println("Message field to save in database.")
		}

	}

	return nil
}

// Check if all phone books are exist.
func CheckPhoneBooksExist(ctx context.Context, db *gorm.DB, phoneBooks []string) (bool, error) {
	var exists bool

	for _, pb := range phoneBooks {
		err := db.WithContext(ctx).Model(&models.PhoneBook{}).
			Select("count(*) > 0").
			Where("name = ?", pb).
			Find(&exists).
			Error
		if err != nil {
			log.Println(err)
			return false, err
		}

		if !exists {
			return false, nil
		}
	}

	return true, nil
}

func sendMessageApiWithError(message string, phoneNumber string) error {
	n := rand.Intn(10)
	if n%2 == 0 {
		return errors.New("field to send message!")
	} else {
		return nil
	}
}

func sendMessageApiWithSuccess(message string, phoneNumber string) error {
	return nil
}

func sendMessageApi(message string, phoneNumber string) error {
	// return sendMessageApiWithError(message, phoneNumber)
	return sendMessageApiWithSuccess(message, phoneNumber)
}

func SendGroupMessage(
	ch chan<- SendMessageStatus,
	message string,
	messageID int,
	account models.Account,
	phoneNumber models.PhoneBookNumber,
) {
	err := sendMessageApi(message, phoneNumber.Phone)
	if err != nil {
		ch <- SendMessageStatus{ID: messageID, Status: false}
		log.Printf(
			"Field to sent message - Text: %s, Source: %s, Destination: %s",
			message,
			account.Username,
			phoneNumber.Phone,
		)
		return
	}
	// Log the message details
	log.Printf(
		"Message sent - Text: %s, Source: %s, Destination: %s",
		message,
		account.Username,
		phoneNumber.Phone,
	)
	ch <- SendMessageStatus{ID: messageID, Status: true}
}

func DoesAcountHaveBudget(smsCost int, smsCounts int, budget int64) bool {
	return budget > int64(smsCost*smsCounts)
}
