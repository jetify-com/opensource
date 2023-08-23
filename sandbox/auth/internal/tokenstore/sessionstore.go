package tokenstore

import (
	"errors"
	"go/token"
	"os"
)

type Store struct {
	rootDir string
}

func New(rootDir string) (*Store, error) {
	// The store contains tokens that enable a particular user to authenticate.
	// It's important that the directory can only be read by that user.
	err := os.MkdirAll(rootDir, 0700)
	if err != nil {
		return nil, err
	}
	return &Store{
		rootDir: rootDir,
	}, nil
}

func (s *Store) ReadToken(issuer string, clientID string) *token.Token {
	var tok token.Token
	path := s.path(issuer, clientID)
	err := readJSONFile(path, &tok)
	if err != nil {
		return nil
	}
	return &tok
}

func (s *Store) WriteToken(issuer string, clientID string, tok *token.Token) error {
	if tok == nil {
		// A nil token is the same as deleting the token.
		return s.DeleteToken(issuer, clientID)
	}
	path := s.path(issuer, clientID)
	return writeJSONFile(path, tok)
}

func (s *Store) DeleteToken(issuer string, clientID string) error {
	path := s.path(issuer, clientID)

	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		// If the file doesn't exist, then we don't need to delete it. It's a no-op.
		return nil
	}

	return os.Remove(path)
}
