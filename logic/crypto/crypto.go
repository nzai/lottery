package crypto

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"log"
)

func Encode(args ...string) string {
	t := sha1.New()

	for index := 0; index < len(args); index++ {
		_, err := t.Write([]byte(args[index]))
		if err != nil {
			log.Printf("加密时出错:%s\n", err.Error())
			return ""
		}
	}

	return hex.EncodeToString(t.Sum(nil))
}

func getUnique(l int) ([]byte, error) {
	buffer := make([]byte, l)
	_, err := rand.Read(buffer)
	if err != nil {
		return nil, err
	}

	return buffer, nil
}

//  生成UUID(32位)
func GetUUID() string {

	buffer := make([]byte, 16)
	_, err := rand.Read(buffer)
	if err != nil {
		return ""
	}

	return hex.EncodeToString(buffer)
}

//  生成UUID对应的Int64
func GetUniqueInt64() int64 {

	buffer, err := getUnique(16)
	if err != nil {
		return 0
	}

	result := int64(binary.BigEndian.Uint64(buffer))
	if result < 0 {
		return -result
	} else {
		return result
	}
}
