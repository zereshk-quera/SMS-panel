# SMS Panel Project

The SMS Panel Project is a simple and efficient platform for managing and sending SMS messages, built with Golang using GORM and the Echo web framework. This project allows users to register, log in, manage phone books and phone numbers, and send SMS messages to individuals or groups.

## Features

- User registration and login
- Account budget management
- Payment gateway creation and verification
- Phone book and phone number CRUD operations
- Single, Periodic, and Group SMS sending
- Sender number management and purchase
- SMS Sending with Templates: Personalize your SMS messages by including customizable templates with variables. Supported variables include:
  - `%name`: Replaced with the name of the phone book number.
  - `%date`: Replaced with the current date and time.
  - `%prefix`: Replaced with the prefix of the phone book number object.
  - `%username`: Replaced with the username of the phone book number.
- Admin Panel for managing users, configurations, and SMS messages

## Installation

1. Clone the repository to your local machine:

```bash
git clone https://github.com/zereshk-quera/SMS-panel.git
```

2. Navigate into the project directory:

```bash
cd SMS-panel
```

3. Install the required dependencies:

```bash
go mod download
```

4. Run the project:

```bash
go run .
```

The server will start running at `localhost:8080`.

5. Access the Swagger API documentation:

Swagger URL: [http://localhost:8080/swagger/](http://localhost:8080/swagger/)

The Swagger URL provides access to the Swagger API documentation for the SMS Panel Project.

## Usage

### Authentication

1. Register a new user:

```
POST accounts/register
```

2. Log in with the registered user:

```
POST accounts/login
```

### Account Management

1. Get user budget:

```
GET accounts/budget
```

### Payment

1. Create a payment gateway link:

```
POST accounts/payment/request
```

2. Verify the payment:

```
GET accounts/payment/verify
```

### Phone Book Management

1. CRUD operations for phone books:

```
GET /phone-books
POST /phone-books
PUT /phone-books/:id
DELETE /phone-books/:id
```

2. CRUD operations for phone book numbers:

```
GET accounts/phone-books/:id/numbers
POST accounts/phone-books/:id/numbers
PUT accounts/phone-books/:id/numbers/:numberId
DELETE accounts/phone-books/:id/numbers/:numberId
```

### SMS Sending

1. Get all available sender numbers:

```
GET accounts/sender-numbers
```

2. Get all sender number for purchase:

```
GET accounts/sender-numbers/sale
```

3. Send a single SMS:

```
POST /sms/single
```

4. Send a periodic SMS:

```
POST /sms/periodic
```

5. Send a group SMS:

```
POST /sms/phonebooks
```

## Admin Panel

The Admin Panel is a web-based interface for managing users, configurations, and SMS messages.

### Authentication

1. Log in as an admin:

```
POST admin/login
POST admin/register
```

### User Management

1. Activate a user:

```
PATCH admin/activate/:id
```

2. Deactivate a user:

```
PATCH admin/deactivate/:id
```

### Configuration Management

1. Create the configuration:

```
POST admin/add-config
```

### SMS Search and Reporting

1. Search for SMS messages containing a specific word:

```
GET admin/search/:word
```

2. Get a report of all SMS messages sent:

```
GET admin/sms-report
```

### Bad Word Management

1. Add a new bad word:

```
POST admin/bad-words/:word
```

### SMS Sending with Templates

The SMS Panel Project allows you to send SMS messages with customizable templates that can be dynamically populated with variables. This feature enables you to personalize the messages by inserting specific values such as names, dates, prefixes, and usernames into the content of the SMS.

To use this feature, you can include special placeholder variables in your SMS message templates. When sending the SMS, these variables will be substituted with the corresponding values based on the context. Here are the supported variables and their usage:

- `%name`: Replaced with the name of the phone book number.
- `%date`: Replaced with the current date and time.
- `%prefix`: Replaced with the prefix of the phone book number object.
- `%username`: Replaced with the username of the phone book number.

Here's an example of how you can use these variables in your SMS message template:

```plaintext
Hello %name! This is a reminder for your appointment on %date. Please arrive on time. Your prefix is %prefix. Thank you, %username.
```

In the above example, when sending an SMS using this template, the `%name` variable will be replaced with the actual name of the phone book number, `%date` will be replaced with the current date and time, `%prefix` will be replaced with the prefix of the phone book number, and `%username` will be replaced with the username of the phone book number.
