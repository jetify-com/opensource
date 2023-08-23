package tokenstore

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/gosimple/slug"
	"github.com/pkg/errors"
)

func (s *Store) path(issuer string, clientID string) string {
	return filepath.Join(s.rootDir, issuerSlug(issuer), slug.Make(clientID)+".json")
}

func issuerSlug(issuer string) string {
	issuer = strings.TrimPrefix(issuer, "https://")
	issuer = strings.TrimPrefix(issuer, "http://")
	issuer = strings.TrimSuffix(issuer, "/")
	return slug.Make(issuer)
}

func ensureDir(path string) error {
	dir := filepath.Dir(path)
	return os.MkdirAll(dir, 0700)
}

func writeJSONFile(path string, value any) error {
	err := ensureDir(path)
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return errors.WithStack(err)
	}
	return errors.WithStack(os.WriteFile(path, data, 0644))
}

func readJSONFile(path string, value any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return errors.WithStack(err)
	}
	return errors.WithStack(json.Unmarshal(data, value))
}
