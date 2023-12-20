package tokenstore

import (
	"errors"
	"os"

	"go.jetpack.io/pkg/auth/session"
)

const storeDataVersion = "1"

type Store struct {
	rootDir string
}

type storeData struct {
	Version string           `json:"version"`
	Tokens  []*session.Token `json:"tokens"`
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

// ReadToken returns the top token for the given issuer and clientID.
// Top token was the last token written to the store. (last login/refresh)
func (s *Store) ReadToken(issuer string, clientID string) (*session.Token, error) {
	data, err := s.readData(issuer, clientID)
	if err != nil {
		return nil, err
	}
	if len(data.Tokens) == 0 {
		return nil, os.ErrNotExist
	}
	return data.Tokens[0], nil
}

func (s *Store) ReadTokens(issuer string, clientID string) ([]*session.Token, error) {
	data, err := s.readData(issuer, clientID)
	if err != nil {
		return nil, err
	}
	return data.Tokens, nil
}

// FindToken returns the first token for the given issuer and clientID where
// fn returns true.
func (s *Store) FindToken(
	issuer, clientID string,
	fn func(tok *session.Token) bool,
) (*session.Token, error) {
	tokens, err := s.ReadTokens(issuer, clientID)
	if err != nil {
		return nil, err
	}
	for _, tok := range tokens {
		if fn(tok) {
			return tok, nil
		}
	}
	return nil, os.ErrNotExist
}

func (s *Store) WriteToken(issuer string, clientID string, tok *session.Token) error {
	if tok == nil {
		// A nil token is the same as deleting the token.
		return s.DeleteToken(issuer, clientID)
	}
	data, err := s.readData(issuer, clientID)
	if err != nil {
		return err
	}
	data.addToken(tok)
	path := s.path(issuer, clientID)
	return writeJSONFile(
		path,
		data,
	)
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

func (s *Store) readData(issuer, clientID string) (storeData, error) {
	path := s.path(issuer, clientID)
	var data storeData
	err := readJSONFile(path, &data)
	if errors.Is(err, os.ErrNotExist) {
		return storeData{Version: storeDataVersion}, nil
	} else if err != nil {
		return storeData{Version: storeDataVersion}, err
	}
	return data, nil
}

func (sd *storeData) addToken(tok *session.Token) {
	tokens := []*session.Token{tok}
	for _, t := range sd.Tokens {
		if t.IDClaims().Subject != tok.IDClaims().Subject {
			tokens = append(tokens, t)
		}
	}
	sd.Tokens = tokens
}
