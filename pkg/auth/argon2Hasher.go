package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"github.com/HunterGooD/voice_friend/user_service/internal/domain/entity"
	"regexp"
	"strconv"

	"github.com/pkg/errors"
	"golang.org/x/crypto/argon2"
)

const templateArgonString = "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"

var parseRegExp = regexp.MustCompile(`\$argon2id\$v=(\d+)\$m=(\d+),t=(\d+),p=(\d+)\$([^$]+)\$([^$]+)`)

type Argon2Hasher struct {
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
	saltLen uint32
}

func NewArgon2Hasher(time, memory, keyLen, saltLen uint32, threads uint8) *Argon2Hasher {
	return &Argon2Hasher{time, memory, threads, keyLen, saltLen}
}

func (ah *Argon2Hasher) HashPassword(password string) (string, error) {
	salt, err := ah.generateSalt()
	if err != nil {
		return "", errors.Wrap(err, "generate salt error")
	}
	hash := argon2.IDKey([]byte(password), salt, ah.time, ah.memory, ah.threads, ah.keyLen)

	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)
	encodedHash := base64.RawStdEncoding.EncodeToString(hash)
	hashedPassword := fmt.Sprintf(templateArgonString, argon2.Version, ah.memory, ah.time, ah.threads, encodedSalt, encodedHash)

	return hashedPassword, nil
}

func (ah *Argon2Hasher) CheckPassword(password, hashedPassword string) (bool, error) {
	var parallelism uint8
	var memory, iterations uint32
	var salt, hash string

	matches := parseRegExp.FindStringSubmatch(hashedPassword)
	if matches == nil || len(matches) != 7 {
		return false, errors.Wrap(entity.ErrDataNotValid, "Error parsing password params")
	}

	// skip 0 index because its source string
	var err error
	memory64, err := strconv.ParseUint(matches[2], 10, 32)
	if err != nil {
		return false, errors.Wrap(err, "Error parsing memory params")
	}
	memory = uint32(memory64)

	iterations64, err := strconv.ParseUint(matches[3], 10, 32)
	if err != nil {
		return false, errors.Wrap(err, "Error parsing iterations params")
	}
	iterations = uint32(iterations64)

	parallelism64, err := strconv.ParseUint(matches[4], 10, 8)
	if err != nil {
		return false, errors.Wrap(err, "Error parsing parallelism params")
	}
	parallelism = uint8(parallelism64)

	salt = matches[5]
	hash = matches[6]

	saltBytes, err := base64.RawStdEncoding.DecodeString(salt)
	if err != nil {
		return false, errors.Wrap(err, "error decode base64 salt")
	}

	hashBytes, err := base64.RawStdEncoding.DecodeString(hash)
	if err != nil {
		return false, errors.Wrap(err, "error decode base64 hash password")
	}

	newHash := argon2.IDKey([]byte(password), saltBytes, iterations, memory, parallelism, ah.keyLen)

	return subtle.ConstantTimeCompare(hashBytes, newHash) == 1, nil
}

func (ah *Argon2Hasher) generateSalt() ([]byte, error) {
	salt := make([]byte, ah.saltLen)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}
