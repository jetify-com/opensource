package session

func (t *Token) Verify() error {
	if t.Keys != nil {
		return nil
	}

	return nil
}
