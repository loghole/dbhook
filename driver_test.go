package dbhook

import (
	"database/sql/driver"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	sqlite "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"

	"github.com/loghole/dbhook/mocks"
)

func TestDriver_Open(t *testing.T) {
	t.Parallel()

	var (
		ctrl = gomock.NewController(t)
		name = ":memory:"
	)

	type fields struct {
		makeDriver func() driver.Driver
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "pass",
			fields: fields{
				makeDriver: func() driver.Driver {
					return &sqlite.SQLiteDriver{}
				},
			},
		},
		{
			name: "failed open",
			fields: fields{
				makeDriver: func() driver.Driver {
					driver := mocks.NewMockDriver(ctrl)
					driver.EXPECT().Open(name).Times(1).Return(nil, errors.New("some error"))

					return driver
				},
			},
			wantErr: true,
		},
		{
			name: "open ExecerQueryerSessionResetter",
			fields: fields{
				makeDriver: func() driver.Driver {
					conn := mocks.NewMockAdvancedConn(ctrl)

					driver := mocks.NewMockDriver(ctrl)
					driver.EXPECT().Open(name).Times(1).Return(conn, nil)

					return driver
				},
			},
		},
		{
			name: "open Execer",
			fields: fields{
				makeDriver: func() driver.Driver {
					conn := mocks.NewMockEConn(ctrl)

					driver := mocks.NewMockDriver(ctrl)
					driver.EXPECT().Open(name).Times(1).Return(conn, nil)

					return driver
				},
			},
		},
		{
			name: "open Queryer",
			fields: fields{
				makeDriver: func() driver.Driver {
					conn := mocks.NewMockQConn(ctrl)

					driver := mocks.NewMockDriver(ctrl)
					driver.EXPECT().Open(name).Times(1).Return(conn, nil)

					return driver
				},
			},
		},
		{
			name: "open Query",
			fields: fields{
				makeDriver: func() driver.Driver {
					conn := mocks.NewMockQueryerConn(ctrl)

					driver := mocks.NewMockDriver(ctrl)
					driver.EXPECT().Open(name).Times(1).Return(conn, nil)

					return driver
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			drv := &Driver{Driver: tt.fields.makeDriver()}

			got, err := drv.Open(name)
			switch tt.wantErr {
			case true:
				assert.Error(t, err)
				assert.Nil(t, got)
			case false:
				assert.NoError(t, err)
				assert.NotNil(t, got)
			}
		})
	}
}
