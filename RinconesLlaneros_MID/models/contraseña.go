package models

import (
	"golang.org/x/crypto/bcrypt"
)

// Función para verificar si una contraseña coincide con el hash
func VerificarContraseña(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Función para hashear una contraseña
func HashearContraseña(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
