package tokenstore

import (
	"errors"
	"os"

	"go.jetify.com/pkg/auth/session"
)

const storeDataVersion = "1"

type Store struct {
	rootDir string
}

type storeData struct {
	Version string `json:"version"`
	// Token order is significant. The first token is the default token.
	Tokens []*session.Token `json:"tokens"`
}

func New(rootDir string) (*Store, error) {
	// The store contains tokens that enable a particular user to authenticate.
	// It's important that the directory can only be read by that user.
	err := os.MkdirAll(rootDir, 0o700)
	if err != nil {
		return nil, err
	}
	return &Store{
		rootDir: rootDir,
	}, nil
}

func (s *Store) ReadTokens(issuer string, clientID string) ([]*session.Token, error) {
	data, err := s.readData(issuer, clientID)
	if err != nil {
		return nil, err
	}
	return data.Tokens, nil
}

func (s *Store) WriteToken(
	issuer, clientID string,
	tok *session.Token,
	makeDefault bool,
) error {
	if tok == nil {
		return errors.New("token is nil")
	}
	data, err := s.readData(issuer, clientID)
	if err != nil {
		return err
	}

	if makeDefault {
		data.addDefaultToken(tok)
	} else {
		data.replaceToken(tok)
	}

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
	data := storeData{Version: storeDataVersion}
	err := readJSONFile(path, &data)
	if errors.Is(err, os.ErrNotExist) {
		return data, nil
	} else if err != nil {
		return data, err
	}
	return data, nil
}

func (sd *storeData) addDefaultToken(tok *session.Token) {
	tokens := []*session.Token{tok}
	for _, t := range sd.Tokens {
		if t.IDClaims().Subject != tok.IDClaims().Subject {
			tokens = append(tokens, t)
		}
	}
	sd.Tokens = tokens
}

func (sd *storeData) replaceToken(tok *session.Token) {
	for idx, t := range sd.Tokens {
		if t.IDClaims().Subject == tok.IDClaims().Subject {
			sd.Tokens[idx] = tok
		}
	}
	// If we didn't find a match, should we add to end of list?
	// It likely means the token file was modified concurrently.
}
