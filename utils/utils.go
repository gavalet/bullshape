package utils

import (
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func NewUUIDV4() string {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return ""
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

func GetPWD() (string, error) {
	var path string
	cmd := exec.Command("pwd")
	out, err := cmd.Output()
	if err != nil {
		return path, err
	}
	path = string(out)
	path = strings.TrimSuffix(path, "\n")
	return path, nil
}

func EncryptPass(pass string) string {
	encryptedPass, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	return string(encryptedPass)
}

func EnvString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}
