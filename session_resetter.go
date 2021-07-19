package dbhook

import (
	"context"
	"database/sql/driver"
)

type SessionResetter struct {
	*Conn
}

func (s *SessionResetter) ResetSession(ctx context.Context) error {
	c, ok := s.Conn.Conn.(driver.SessionResetter)
	if !ok {
		return ErrNonSessionResetter
	}

	return c.ResetSession(ctx) // nolint:wrapcheck // need clear error
}
