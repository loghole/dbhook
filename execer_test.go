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

func TestExecerContext_ExecContext(t *testing.T) {
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
		want    driver.Result
		wantErr bool
	}{
		{
			name: "pass with ExecContext",
			fields: fields{
				makeConn: func() driver.Conn {
					conn := mocks.NewMockAdvancedConn(ctrl)
					conn.EXPECT().
						ExecContext(ctx, "INSERT", nil).
						Times(1).
						DoAndReturn(func(ctx context.Context, query string, arg []driver.NamedValue) (driver.Result, error) {
							rows := mocks.NewMockResult(ctrl)

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
				query: "INSERT",
			},
			want: mocks.NewMockResult(ctrl),
		},
		{
			name: "pass with Exec",
			fields: fields{
				makeConn: func() driver.Conn {
					conn := mocks.NewMockQueryerConn(ctrl)
					conn.EXPECT().
						Exec("INSERT", []driver.Value{}).
						Times(1).
						DoAndReturn(func(query string, arg []driver.Value) (driver.Result, error) {
							rows := mocks.NewMockResult(ctrl)

							return rows, nil
						})

					return conn
				},
				makeHooks: func() Hook {
					return nil
				},
			},
			args: args{
				query: "INSERT",
			},
			want: mocks.NewMockResult(ctrl),
		},
		{
			name: "error with Exec named",
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
				query: "INSERT",
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
			name: "without Exec",
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
				query: "INSERT",
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
					hook := NewHooks(WithHook(&testHook{}), WithHooksBefore(&testBeforeHookWithError{}))

					return hook
				},
			},
			args:    args{query: "INSERT"},
			want:    nil,
			wantErr: true,
		},
		{
			name: "with after error",
			fields: fields{
				makeConn: func() driver.Conn {
					conn := mocks.NewMockAdvancedConn(ctrl)
					conn.EXPECT().
						ExecContext(ctx, "INSERT", nil).
						Times(1).
						DoAndReturn(func(ctx context.Context, query string, arg []driver.NamedValue) (driver.Result, error) {
							rows := mocks.NewMockResult(ctrl)

							return rows, nil
						})

					return conn
				},
				makeHooks: func() Hook {
					hook := NewHooks(WithHook(&testHook{}), WithHooksAfter(&testAfterHookWithError{}))

					return hook
				},
			},
			args:    args{query: "INSERT"},
			want:    nil,
			wantErr: true,
		},
		{
			name: "with error",
			fields: fields{
				makeConn: func() driver.Conn {
					conn := mocks.NewMockAdvancedConn(ctrl)
					conn.EXPECT().
						ExecContext(ctx, "INSERT", nil).
						Times(1).
						DoAndReturn(func(ctx context.Context, query string, arg []driver.NamedValue) (driver.Result, error) {
							return nil, errors.New("some error")
						})

					return conn
				},
				makeHooks: func() Hook {
					return nil
				},
			},
			args:    args{query: "INSERT"},
			want:    nil,
			wantErr: true,
		},
		{
			name: "with error with hook",
			fields: fields{
				makeConn: func() driver.Conn {
					conn := mocks.NewMockAdvancedConn(ctrl)
					conn.EXPECT().
						ExecContext(ctx, "INSERT", nil).
						Times(1).
						DoAndReturn(func(ctx context.Context, query string, arg []driver.NamedValue) (driver.Result, error) {
							return nil, errors.New("some error")
						})

					return conn
				},
				makeHooks: func() Hook {
					hook := NewHooks(WithHook(&testHook{}), WithHooksError(&testErrorHookWithError{}))

					return hook
				},
			},
			args:    args{query: "INSERT"},
			want:    nil,
			wantErr: true,
		},
		{
			name: "with fix error",
			fields: fields{
				makeConn: func() driver.Conn {
					conn := mocks.NewMockAdvancedConn(ctrl)
					conn.EXPECT().
						ExecContext(ctx, "INSERT", nil).
						Times(1).
						DoAndReturn(func(ctx context.Context, query string, arg []driver.NamedValue) (driver.Result, error) {
							return nil, errors.New("some error")
						})

					return conn
				},
				makeHooks: func() Hook {
					hook := NewHooks(WithHooksError(&testFixErrorHookWithError{}))

					return hook
				},
			},
			args: args{query: "INSERT"},
			want: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			conn := &ExecerContext{
				Conn: &Conn{
					Conn:  tt.fields.makeConn(),
					hooks: tt.fields.makeHooks(),
				},
			}

			got, err := conn.ExecContext(ctx, tt.args.query, tt.args.args)
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
