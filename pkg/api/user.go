package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/mail"
)

type photographer struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Address   string `json:"address"`
	City      string `json:"city"`
	Gender    string `json:"gender"`
	Birthdate string `json:"birthdate"`
}

func (p *photographer) createPhotographer(db *sql.DB) error {
	_, err := mail.ParseAddress(p.Email)
	if err != nil {
		return err
	}

	statement := fmt.Sprintf("INSERT INTO photographers(username, email, phone, gender) VALUES('%s','%s','%s','%s')", p.Username, p.Email, p.Phone, p.Gender)
	_, err = db.Exec(statement)
	if err != nil {
		return err
	}

	return nil
}

func (p *photographer) getPhotographer(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT email, phone, gender FROM photographers where username='%s'", p.Username)
	return db.QueryRow(statement).Scan(&p.Email, &p.Phone, &p.Gender)
}

func (p *photographer) updatePhotographer(db *sql.DB) error {
	_, err := mail.ParseAddress(p.Email)
	if err != nil {
		return err
	}

	statement := fmt.Sprintf("UPDATE photographers SET email='%s', phone='%s', gender='%s' WHERE username='%s'", p.Email, p.Phone, p.Gender, p.Username)
	_, err = db.Exec(statement)
	if err != nil {
		return err
	}

	return nil
}

func (p *photographer) deletePhotographer(db *sql.DB) error {
	statement := fmt.Sprintf("DELETE FROM photographers WHERE username='%s'", p.Username)
	_, err := db.Exec(statement)
	if err != nil {
		return err
	}
	return nil
}

func getPhotographers(db *sql.DB, start, count int) ([]photographer, error) {
	return nil, errors.New("Not implemented yet!")
}
