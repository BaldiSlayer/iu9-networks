package main

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"math/rand"
	"net/smtp"
	"time"
)

func sendEmail(toEmail string, subject string, messageBody template.HTML) (bool, error) {
	smtpHost := "mail.nic.ru"
	smtpPort := 465
	username := "dts21@dactyl.su"
	password := "12345678990DactylSUDTS"

	// Формирование сообщения
	// message := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", toEmail, subject, messageBody)
	message := fmt.Sprintf("To: %s\r\nSubject: %s\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s", toEmail, subject, messageBody)

	// Настройка подключения с поддержкой SSL
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // Настройте это правильно для безопасного использования в боевом режиме
		ServerName:         smtpHost,
	}

	// Подключение к SMTP серверу
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", smtpHost, smtpPort), tlsConfig)
	if err != nil {
		return false, fmt.Errorf("error connecting to SMTP server: %s", err)
	}
	defer conn.Close()

	// Авторизация на SMTP сервере
	auth := smtp.PlainAuth("", username, password, smtpHost)
	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		return false, fmt.Errorf("error creating SMTP client: %s", err)
	}
	defer client.Close()

	if err := client.Auth(auth); err != nil {
		return false, fmt.Errorf("error authenticating: %s", err)
	}

	// Отправка сообщения
	if err := client.Mail(username); err != nil {
		return false, fmt.Errorf("error setting sender: %s", err)
	}

	if err := client.Rcpt(string(toEmail)); err != nil {
		return false, fmt.Errorf("error setting recipient: %s", err)
	}

	w, err := client.Data()
	if err != nil {
		return false, fmt.Errorf("error opening data connection: %s", err)
	}
	defer w.Close()

	_, err = w.Write([]byte(message))
	if err != nil {
		return false, fmt.Errorf("error writing message: %s", err)
	}

	// Завершение соединения
	if err := client.Quit(); err != nil {
		return false, fmt.Errorf("error quitting: %s", err)
	}

	return true, nil
}

func task1() {
	success, err := sendEmail("lisov.a2005@yandex.ru", "Помогите пожалуйста", "Хочу сдохнуть")
	if err != nil {
		fmt.Println("Email sending failed:", err)
	} else if success {
		fmt.Println("Email sent successfully!")
	} else {
		fmt.Println("Email sending failed for unknown reasons.")
	}
}

type EmailData struct {
	Username string
	Email    string
	Message  string
}

func insertRecord(db *sql.DB, username string, email string, message string) {
	query := `
        INSERT INTO email_distribution_lisov (username, email, message)
        SELECT * FROM (SELECT ?, ?, ?) AS tmp
        WHERE NOT EXISTS (
            SELECT id FROM email_distribution_lisov WHERE username = ? AND email = ?
        ) LIMIT 1
    `

	result, err := db.Exec(query, username, email, message, username, email)
	if err != nil {
		log.Fatal(err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected > 0 {
		log.Println("Record inserted successfully")
	} else {
		log.Println("Record already exists")
	}
}

func task2() {
	db, err := sql.Open("mysql", "iu9networkslabs:Je2dTYr6@tcp(students.yss.su:3306)/iu9networkslabs")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("DROP TABLE IF EXISTS email_distribution_lisov")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS email_distribution_lisov (
			id INT AUTO_INCREMENT PRIMARY KEY,
    		username VARCHAR(255),
    		email VARCHAR(255),
    		message TEXT
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		ALTER TABLE email_distribution_lisov CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
	`)
	if err != nil {
		log.Fatal(err)
	}

	result, err := db.Exec(`
    INSERT INTO email_distribution_lisov (username, email, message)
    SELECT * FROM (SELECT 'Alisov', 'lisov.a2005@yandex.ru', '<h1>Hello</h1><p style="font-style: italic;">This is a test message</p>') AS tmp
    WHERE NOT EXISTS (
        SELECT id FROM email_distribution_lisov WHERE username = 'Alisov' AND email = 'lisov.a2005@yandex.ru'
    ) LIMIT 1
  `)
	if err != nil {
		log.Fatal(err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected > 0 {
		log.Println("Record inserted successfully")
	} else {
		log.Println("Record already exists")
	}

	result, err = db.Exec(`
    INSERT INTO email_distribution_lisov (username, email, message)
    SELECT * FROM (SELECT 'Alisov', 'lisov.a2005@yandex.ru', '<h1>Hello</h1><p style="font-style: italic;">This is a test message</p>') AS tmp
    WHERE NOT EXISTS (
        SELECT id FROM email_distribution_lisov WHERE username = 'Alisov' AND email = 'lisov.a2005@yandex.ru'
    ) LIMIT 1
  `)
	if err != nil {
		log.Fatal(err)
	}
	rowsAffected, _ = result.RowsAffected()
	if rowsAffected > 0 {
		log.Println("Record inserted successfully")
	} else {
		log.Println("Record already exists")
	}

	// insertRecord(db, "Alisov", "lisov.a2005@yandex.ru", "<h1>Hello</h1><p style=\"font-style: italic;\">This is a test message</p>")

	rows, err := db.Query("SELECT username, email, message FROM email_distribution_lisov")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var emailDataList []EmailData

	for rows.Next() {
		var emailData EmailData
		if err := rows.Scan(&emailData.Username, &emailData.Email, &emailData.Message); err != nil {
			log.Fatal(err)
		}
		emailDataList = append(emailDataList, emailData)
	}

	for _, emailData := range emailDataList {
		for i := 0; i < 1; i++ {
			_, err := sendEmail(emailData.Email, emailData.Username, template.HTML("Dear, <b>"+emailData.Username+"</b>"+emailData.Message))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Email sent to", emailData.Email)
			randomDelay := time.Duration(rand.Intn(5)+1) * time.Second
			time.Sleep(randomDelay)
		}
	}

}

func main() {
	// task1()
	task2()
}
