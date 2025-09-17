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

func Login(i model.Login) (bool, error) {
    var v int
    q := `SELECT COUNT(*) FROM client WHERE client_id = ? AND PASSWORD = ?`
    err := db.DB.QueryRow(q, i.ID, i.Password).Scan(&v)
    if err != nil {
        return false, err
    }

    if v != 1 {
        return false, nil
    }

    return true, nil
}
