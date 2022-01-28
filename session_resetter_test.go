package dbhook

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/loghole/dbhook/mocks"
)

func TestSessionResetter_ResetSession(t *testing.T) {
	t.Parallel()

	var (
		ctx  = context.Background()
		ctrl = gomock.NewController(t)
	)

	type fields struct {
		makeConn func() *Conn
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "pass",
			fields: fields{
				makeConn: func() *Conn {
					conn := mocks.NewMockAdvancedConn(ctrl)
					conn.EXPECT().
						ResetSession(ctx).
						DoAndReturn(func(ctx context.Context) error {
							return nil
						})

					return &Conn{Conn: conn}
				},
			},
		},
		{
			name: "error",
			fields: fields{
				makeConn: func() *Conn {
					conn := mocks.NewMockAdvancedConn(ctrl)
					conn.EXPECT().
						ResetSession(ctx).
						DoAndReturn(func(ctx context.Context) error {
							return errors.New("some error")
						})

					return &Conn{Conn: conn}
				},
			},
			wantErr: true,
		},
		{
			name: "without resetter",
			fields: fields{
				makeConn: func() *Conn {
					conn := mocks.NewMockConn(ctrl)

					return &Conn{Conn: conn}
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := &SessionResetter{Conn: tt.fields.makeConn()}

			err := s.ResetSession(ctx)
			switch tt.wantErr {
			case true:
				assert.Error(t, err)
			case false:
				assert.NoError(t, err)
			}
		})
	}
}
