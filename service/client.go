package service

import (
	dbcall "my-homepage/dbcall"
	model "my-homepage/struct"
)

// 회원가입
func AddSignup(i model.Signup) error {
	err := dbcall.InsertSignup(i)
	if err != nil {
		return err
	}

	return nil
}

// 로그인
func Login(i model.Login) (bool, error) {
	ok, err := dbcall.Login(i)
	if err != nil {
		return false, err
	}

	return ok, nil
}
