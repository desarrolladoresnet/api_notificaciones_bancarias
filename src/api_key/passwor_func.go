package api_key

import (
	"crypto/rand"
	"crypto/subtle"

	"encoding/base64"
	"log"

	"golang.org/x/crypto/argon2"
)

////////////////////////////////////////////
////////////////////////////////////////////
////////////////////////////////////////////

func HashPassword(password string) (string, error) {
	const saltLength = 16
	const keyLength = 32

	// Generar salt
	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// Generar hash usando argon2 y el salt generado
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, keyLength)

	// Combinar salt y hash, y codificarlos en base64
	hashSalt := append(salt, hash...)
	encoded := base64.RawStdEncoding.EncodeToString(hashSalt)

	log.Printf("Pass: %v\nHash generado: %v", password, encoded) // Log para depuración
	return encoded, nil
}

////////////////////////////////////////////
////////////////////////////////////////////
////////////////////////////////////////////

// CheckPasswordHash compara la contraseña proporcionada con el hash almacenado
func CheckPasswordHash(password, encodedHash string) bool {
	const saltLength = 16
	const keyLength = 32

	// fmt.Printf("Password: %v\n", password)
	// fmt.Printf("Encoded Hash: %v\n", encodedHash)

	// Decodificar el string base64 para obtener el salt y el hash
	hashSalt, err := base64.RawStdEncoding.DecodeString(encodedHash)
	if err != nil {
		log.Printf("Error al decodificar el hash: %v", err)
		return false
	}

	// Extraer el salt y el hash almacenados
	salt := hashSalt[:saltLength]
	storedHash := hashSalt[saltLength:]

	// Calcular el hash con la contraseña proporcionada y el mismo salt
	newHash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, keyLength)

	// log.Printf("Hash calculado: %v", newHash)     // Log para depuración
	// log.Printf("Hash almacenado: %v", storedHash) // Log para depuración

	// Comparar el hash almacenado y el hash calculado
	return subtle.ConstantTimeCompare(storedHash, newHash) == 1
}
