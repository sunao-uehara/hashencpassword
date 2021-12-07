package storages

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"log"
)

// global variable to store passwords
var PasswordList = make(map[int]*Password)

type Password struct {
	ID       int    `json:"id"`
	Password string `json:"password"`
}

// CreatePasswrod creates new id and store into global variables
func CreatePasswordID() int {
	log.Println("CreatePasswordID")

	newID := len(PasswordList) + 1 // increment by 1
	p := &Password{
		ID: newID,
	}
	PasswordList[newID] = p

	return newID
}

// UpdatePassword updates password by id
func UpdatePassword(id int, password string) error {
	log.Println("UpdatePassword for id:", id)

	p, ok := PasswordList[id]
	if !ok {
		return fmt.Errorf("id %d not found", id)
	}

	var err error
	p.Password, err = hashPassword(password)
	if err != nil {
		return err
	}
	PasswordList[id] = p

	return nil
}

// GetPassword gets Password obj by id
func GetPassword(id int) (*Password, error) {
	log.Println("GetPassword for id:", id)

	p, ok := PasswordList[id]
	if !ok {
		return nil, fmt.Errorf("id %d not found", id)
	}

	return p, nil
}

// hashPassword generates base64 encoded strings of sha512 hash of the provided string
func hashPassword(password string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("password is empty")
	}

	// hash and encode password
	sha := sha512.New()
	sha.Write([]byte(password))
	enc := base64.URLEncoding.EncodeToString(sha.Sum(nil))

	// log.Println(enc)
	return enc, nil
}
