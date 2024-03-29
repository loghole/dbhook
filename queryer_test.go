package dbhook

import (
	"context"
	"database/sql/driver"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/loghole/dbhook/mocks"
)

func TestQueryerContext_QueryContext(t *testing.T) {
	t.Parallel()

	var (
		ctx  = context.Background()
		ctrl = gomock.NewController(t)
	)

	type fields struct {
		makeConn  func() driver.Conn
		makeHooks func() Hook
	}

	type args struct {
		query string
		args  []driver.NamedValue
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    driver.Rows
		wantErr bool
	}{
		{
			name: "pass with QueryContext",
			fields: fields{
				makeConn: func() driver.Conn {
					conn := mocks.NewMockAdvancedConn(ctrl)
					conn.EXPECT().
						QueryContext(ctx, "SELECT", nil).
						Times(1).
						DoAndReturn(func(ctx context.Context, query string, arg []driver.NamedValue) (driver.Rows, error) {
							rows := mocks.NewMockRows(ctrl)

							return rows, nil
						})

					return conn
				},
				makeHooks: func() Hook {
					hook := NewHooks(WithHook(&testHook{}))

					return hook
				},
			},
			args: args{
				query: "SELECT",
			},
			want: mocks.NewMockRows(ctrl),
		},
		{
			name: "pass with Query",
			fields: fields{
				makeConn: func() driver.Conn {
					conn := mocks.NewMockQueryerConn(ctrl)
					conn.EXPECT().
						Query("SELECT", []driver.Value{}).
						Times(1).
						DoAndReturn(func(query string, arg []driver.Value) (driver.Rows, error) {
							rows := mocks.NewMockRows(ctrl)

							return rows, nil
						})

					return conn
				},
				makeHooks: func() Hook {
					return nil
				},
			},
			args: args{
				query: "SELECT",
			},
			want: mocks.NewMockRows(ctrl),
		},
		{
			name: "error with Query named",
			fields: fields{
				makeConn: func() driver.Conn {
					conn := mocks.NewMockQueryerConn(ctrl)

					return conn
				},
				makeHooks: func() Hook {
					return nil
				},
			},
			args: args{
				query: "SELECT",
				args: []driver.NamedValue{
					{
						Ordinal: 1,
						Value:   1,
					},
					{
						Name:    "name",
						Ordinal: 2,
						Value:   2,
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "without Query",
			fields: fields{
				makeConn: func() driver.Conn {
					conn := mocks.NewMockConn(ctrl)

					return conn
				},
				makeHooks: func() Hook {
					return nil
				},
			},
			args: args{
				query: "SELECT",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "with before error",
			fields: fields{
				makeConn: func() driver.Conn {
					conn := mocks.NewMockAdvancedConn(ctrl)

					return conn
				},
				makeHooks: func() Hook {
					hook := NewHooks(WithHooksBefore(&testBeforeHookWithError{}))

					return hook
				},
			},
			args:    args{query: "SELECT"},
			want:    nil,
			wantErr: true,
		},
		{
			name: "with after error",
			fields: fields{
				makeConn: func() driver.Conn {
					conn := mocks.NewMockAdvancedConn(ctrl)
					conn.EXPECT().
						QueryContext(ctx, "SELECT", nil).
						Times(1).
						DoAndReturn(func(ctx context.Context, query string, arg []driver.NamedValue) (driver.Rows, error) {
							rows := mocks.NewMockRows(ctrl)

							return rows, nil
						})

					return conn
				},
				makeHooks: func() Hook {
					hook := NewHooks(WithHooksAfter(&testAfterHookWithError{}))

					return hook
				},
			},
			args:    args{query: "SELECT"},
			want:    nil,
			wantErr: true,
		},
		{
			name: "with error",
			fields: fields{
				makeConn: func() driver.Conn {
					conn := mocks.NewMockAdvancedConn(ctrl)
					conn.EXPECT().
						QueryContext(ctx, "SELECT", nil).
						Times(1).
						DoAndReturn(func(ctx context.Context, query string, arg []driver.NamedValue) (driver.Rows, error) {
							return nil, errors.New("some error")
						})

					return conn
				},
				makeHooks: func() Hook {
					return nil
				},
			},
			args:    args{query: "SELECT"},
			want:    nil,
			wantErr: true,
		},
		{
			name: "with error with hook",
			fields: fields{
				makeConn: func() driver.Conn {
					conn := mocks.NewMockAdvancedConn(ctrl)
					conn.EXPECT().
						QueryContext(ctx, "SELECT", nil).
						Times(1).
						DoAndReturn(func(ctx context.Context, query string, arg []driver.NamedValue) (driver.Rows, error) {
							return nil, errors.New("some error")
						})

					return conn
				},
				makeHooks: func() Hook {
					hook := NewHooks(WithHooksError(&testErrorHookWithError{}))

					return hook
				},
			},
			args:    args{query: "SELECT"},
			want:    nil,
			wantErr: true,
		},
		{
			name: "with fix error",
			fields: fields{
				makeConn: func() driver.Conn {
					conn := mocks.NewMockAdvancedConn(ctrl)
					conn.EXPECT().
						QueryContext(ctx, "SELECT", nil).
						Times(1).
						DoAndReturn(func(ctx context.Context, query string, arg []driver.NamedValue) (driver.Rows, error) {
							return nil, errors.New("some error")
						})

					return conn
				},
				makeHooks: func() Hook {
					hook := NewHooks(WithHooksError(&testFixErrorHookWithError{}))

					return hook
				},
			},
			args: args{query: "SELECT"},
			want: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			conn := &QueryerContext{
				Conn: &Conn{
					Conn:  tt.fields.makeConn(),
					hooks: tt.fields.makeHooks(),
				},
			}

			got, err := conn.QueryContext(ctx, tt.args.query, tt.args.args)
			switch tt.wantErr {
			case true:
				assert.Error(t, err)
			case false:
				assert.NoError(t, err)
			}

			assert.Equal(t, got, tt.want)
		})
	}
}
