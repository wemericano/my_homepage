package dbcall

import (
	db "my-homepage/database"
	model "my-homepage/struct"
)

func InsertSignup(i model.Signup) error {
	q := `INSERT INTO client (client_id, password, name) VALUES (?, ?, ?)`
	_, err := db.DB.Exec(q, i.ID, i.Password, i.Username)
	if err != nil {
		return err
	}

	return nil
}
