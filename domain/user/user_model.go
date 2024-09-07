package user

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/podossaem/podoroot/domain/exception"
	"golang.org/x/crypto/argon2"
)

type (
	User struct {
		ID              string    `json:"id"`
		JoinMethod      string    `json:"joinMethod"`
		AccountID       string    `json:"accountId"`
		AccountPassword string    `json:"accountPassword"`
		Nickname        string    `json:"nickname"`
		Role            int       `json:"role"`
		CreatedAt       time.Time `json:"createdAt"`
	}

	UserInfo struct {
		JoinMethod string    `json:"joinMethod"`
		AccountID  string    `json:"accountId"`
		Nickname   string    `json:"string"`
		CreatedAt  time.Time `json:"createdAt"`
	}

	hashConfig struct {
		memory      uint32
		iterations  uint32
		parallelism uint8
		saltLength  uint32
		keyLength   uint32
	}
)

func (e *User) ComparePassword(password string) error {
	hc, salt, hash, err := e.decodeHash(e.AccountPassword)
	if err != nil {
		return err
	}

	otherHash := argon2.IDKey([]byte(password), salt, hc.iterations, hc.memory, hc.parallelism, hc.keyLength)

	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return nil
	}

	return exception.ErrUnauthorized
}

func (e *User) HashPassword() error {
	// 추후 password 암호 버저닝을 위해
	hc := &hashConfig{
		memory:      64 * 1024,
		iterations:  3,
		parallelism: 2,
		saltLength:  16,
		keyLength:   32,
	}
	salt := make([]byte, hc.saltLength)
	if _, err := rand.Read(salt); err != nil {
		return err
	}

	hash := argon2.IDKey(
		[]byte(e.AccountPassword),
		salt,
		hc.iterations,
		hc.memory,
		hc.parallelism,
		hc.keyLength,
	)
	base64Salt := base64.RawStdEncoding.EncodeToString(salt)
	base64Hash := base64.RawStdEncoding.EncodeToString(hash)
	encodedHash := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		hc.memory,
		hc.iterations,
		hc.parallelism,
		base64Salt,
		base64Hash,
	)

	e.AccountPassword = encodedHash

	return nil
}

func (e *User) decodeHash(password string) (hc *hashConfig, salt, hash []byte, err error) {
	vals := strings.Split(password, "$")
	if len(vals) != 6 {
		return nil, nil, nil, errors.New("invalid hash")
	}

	var version int
	if _, err := fmt.Sscanf(vals[2], "v=%d", &version); err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, errors.New("incompatible version")
	}

	hc = &hashConfig{}
	if _, err := fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &hc.memory, &hc.iterations, &hc.parallelism); err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	hc.saltLength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	hc.keyLength = uint32(len(hash))

	return hc, salt, hash, nil
}

const (
	JoinMethod_Email = "email"
	Role_User        = 100
	Role_Master      = 10000
)
