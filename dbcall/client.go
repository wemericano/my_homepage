package dbcall

import (
	"log"
	db "my-homepage/database"
	model "my-homepage/struct"
)

func InsertSignup(i model.Signup) error {
	q := `INSERT INTO [user] (user_id, user_pw, user_name) VALUES (@p1, @p2, @p3)`
	_, err := db.DB.Exec(q, i.ID, i.Password, i.Username)
	if err != nil {
		log.Println("InsertSignup error: ", err)
		return err
	}

	return nil
}

func Login(i model.Login) (bool, error) {
	var v int
	q := `SELECT COUNT(*) FROM [user] WHERE user_id = @p1 AND user_pw = @p2`
	err := db.DB.QueryRow(q, i.ID, i.Password).Scan(&v)
	if err != nil {
		log.Println("Login error: ", err)
		return false, err
	}

	if v != 1 {
		return false, nil
	}

	return true, nil
}
