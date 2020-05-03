package session

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"hash"
	"sync"
	"time"

	"github.com/andrebq/authentic/internal/tcache"
)

const (
	randomLegSize   = 8
	tokenSize       = randomLegSize + 8 /* expiration time */
	signatureSize   = 32
	signedTokenSize = tokenSize + signatureSize
)

type (
	signedToken [tokenSize + signatureSize]byte
	// S handles session management using a TokenSet to persist session tokens
	//
	// Tokens don't provide extra information as they are intended to simply allow
	// the access to be protected.
	S struct {
		ts       TokenSet
		key      *Key
		bufPool  sync.Pool
		hmacPool sync.Pool
		duration time.Duration
	}

	// TokenSet is used to store session tokens. Each token have a TTL and the TokenSet
	// SHOULD retain the token for at least the provided TTL.
	//
	// Missing tokens are handled as expired.
	TokenSet interface {
		// Contains should return true if the token exists.
		Contains(string) (bool, error)
		// Add includes the token into the set of valid tokens
		Add(string, time.Time) error
	}

	// Key wraps a key which is used to compute an hmac from tokens
	Key struct {
		buf [32]byte
	}
)

var (
	errInvalid = errors.New("invalid token")
	errExpired = errors.New("expired")
)

// IsValid returns true if, and only if, err == nil.
//
// Check S.Verify for more details
func IsValid(err error) bool {
	return err == nil
}

// IsInvalid indicates a token that has an invalid format or has been tampered with
func IsInvalid(err error) bool {
	return errors.Is(err, errInvalid)
}

// IsExpired indicates a token that has a valid format, has not be tampered but reached its expiration
func IsExpired(err error) bool {
	return errors.Is(err, errExpired)
}

// New returns a new session holder, a nil TokenSet will result in a session which is kept only in memory
// and a nil *Key will result in a unique key which is discarded when the process finishes.
func New(ts TokenSet, hmacKey *Key) (*S, error) {
	var err error
	if ts == nil {
		ts = tcache.New()
	}
	if hmacKey == nil {
		hmacKey, err = RandomKey()
		if err != nil {
			return nil, err
		}
	}
	return &S{
		ts:  ts,
		key: hmacKey,
		bufPool: sync.Pool{
			New: func() interface{} {
				return new(signedToken)
			},
		},
		hmacPool: sync.Pool{
			New: func() interface{} {
				return hmac.New(sha256.New, hmacKey.buf[:])
			},
		},
		// approx one day
		duration: time.Hour * 24,
	}, nil
}

// Start returns a unique and small token to be included in a Cookie
func (s *S) Start(now time.Time) (string, time.Time, error) {
	expire := time.Now().Add(s.duration)
	buf := s.bufPool.Get().(*signedToken)
	defer s.bufPool.Put(buf)

	_, err := rand.Read(buf[:8])
	if err != nil {
		return "", time.Time{}, err
	}

	binary.BigEndian.PutUint64(buf[8:], uint64(expire.Unix()))

	s.sign(buf[tokenSize:], buf[:tokenSize])
	if err != nil {
		return "", time.Time{}, err
	}
	err = s.ts.Add(base64.URLEncoding.EncodeToString(buf[:tokenSize]), expire)
	if err != nil {
		return "", time.Time{}, err
	}
	return base64.URLEncoding.EncodeToString(buf[:]), expire, nil
}

// Verify if the given token is valid, which is indicated by a nil error
//
// Any non-nil error indicates that the session is either invalid/expired or could not be checked
//
// Use IsInvalid/IsExpired to extract more details. Do not expose the error directly to the client
// as it might contain details about server configuration.
//
// Use IsValid to convert this error into a boolean using the logic above
func (s *S) Verify(token string) error {
	enc := base64.URLEncoding

	sz := enc.DecodedLen(len(token))
	if sz != signedTokenSize {
		return errInvalid
	}

	input := s.bufPool.Get().(*signedToken)
	defer s.bufPool.Put(input)

	_, err := enc.Decode(input[:], []byte(token))
	if err != nil {
		return errInvalid
	}

	recomputed := s.bufPool.Get().(*signedToken)
	s.sign(recomputed[tokenSize:], input[:tokenSize])
	valid := hmac.Equal(recomputed[tokenSize:], input[tokenSize:])
	s.bufPool.Put(recomputed)

	if !valid {
		return errInvalid
	}

	expire := time.Unix(int64(binary.BigEndian.Uint64(input[randomLegSize:randomLegSize+8])), 0)
	valid = expire.After(time.Now())

	if !valid {
		return errExpired
	}

	valid, err = s.ts.Contains(enc.EncodeToString(input[:tokenSize]))
	if err != nil {
		return err
	}
	return nil
}

func (s *S) sign(out []byte, in []byte) {
	h := s.hmacPool.Get().(hash.Hash)
	h.Reset()
	h.Write(in)
	h.Sum(out[:0])
	s.hmacPool.Put(h)
}

// RandomKey returns a valid 32-byte key
func RandomKey() (*Key, error) {
	var k Key
	_, err := rand.Read(k.buf[:])
	if err != nil {
		return nil, err
	}
	return &k, nil
}
