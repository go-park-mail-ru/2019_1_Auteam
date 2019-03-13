package server

import (
	"crypto/md5"
	"fmt"
)

func HashPassword(password string) string {
	salt := "lewkrjhnljkdfsgbhkfgjdscwf"
	hash1 := fmt.Sprintf("%x", md5.Sum([]byte(salt+password)))
	hash2 := fmt.Sprint("%x", md5.Sum([]byte(hash1)))
	return fmt.Sprint("%x", md5.Sum([]byte(hash2)))
}
