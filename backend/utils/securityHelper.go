package utils

import "golang.org/x/crypto/bcrypt"

func Hash(password string) ([]byte, error) {
	//qtde de hasheamento dinamico conforme tamanho de password, mínimo: 4 e máximo: 31, conforme bcrypt(07/12/2020)
	quantity := len(password)
	if len(password) <= 10 {
		quantity = 10
	} else if len(password) >= 16 { //16 pra não demorar muito
		quantity = 16
	}

	return bcrypt.GenerateFromPassword([]byte(password), quantity)
}

func VerifyPassword(PasswordWithHash, passwordString string) error {
	return bcrypt.CompareHashAndPassword([]byte(PasswordWithHash), []byte(passwordString))
}
