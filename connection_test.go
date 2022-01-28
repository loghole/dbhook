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

func TestConn_PrepareContext(t *testing.T) {
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
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    driver.Stmt
		wantErr bool
	}{
		{
			name: "pass PrepareContext",
			fields: fields{
				makeConn: func() driver.Conn {
					conn := mocks.NewMockAdvancedConn(ctrl)
					conn.EXPECT().
						PrepareContext(ctx, "INSERT").
						Times(1).
						Return(mocks.NewMockStmtWithContext(ctrl), nil)

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
			want: &Stmt{
				Stmt:  mocks.NewMockStmtWithContext(ctrl),
				hooks: NewHooks(WithHook(&testHook{})),
				query: "INSERT",
			},
		},
		{
			name: "pass Prepare",
			fields: fields{
				makeConn: func() driver.Conn {
					conn := mocks.NewMockConn(ctrl)
					conn.EXPECT().
						Prepare("INSERT").
						Times(1).
						Return(mocks.NewMockStmt(ctrl), nil)

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
			want: &Stmt{
				Stmt:  mocks.NewMockStmt(ctrl),
				hooks: NewHooks(WithHook(&testHook{})),
				query: "INSERT",
			},
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
						PrepareContext(ctx, "INSERT").
						Times(1).
						Return(mocks.NewMockStmtWithContext(ctrl), nil)

					return conn
				},
				makeHooks: func() Hook {
					hook := NewHooks(WithHooksAfter(&testAfterHookWithError{}))

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
						PrepareContext(ctx, "INSERT").
						Times(1).
						Return(nil, errors.New("some error"))

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
						PrepareContext(ctx, "INSERT").
						Times(1).
						Return(nil, errors.New("some error"))

					return conn
				},
				makeHooks: func() Hook {
					hook := NewHooks(WithHooksError(&testErrorHookWithError{}))

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
						PrepareContext(ctx, "INSERT").
						Times(1).
						Return(nil, errors.New("some error"))

					return conn
				},
				makeHooks: func() Hook {
					hook := NewHooks(WithHooksError(&testFixErrorHookWithError{}))

					return hook
				},
			},
			args: args{query: "INSERT"},
			want: &Stmt{
				Stmt:  nil,
				hooks: NewHooks(WithHooksError(&testFixErrorHookWithError{})),
				query: "INSERT",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			conn := &Conn{
				Conn:  tt.fields.makeConn(),
				hooks: tt.fields.makeHooks(),
			}

			got, err := conn.PrepareContext(ctx, tt.args.query)
			switch tt.wantErr {
			case true:
				assert.Error(t, err)
			case false:
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestConn_BeginTx(t *testing.T) {
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
		opts driver.TxOptions
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    driver.Tx
		wantErr bool
	}{
		{
			name: "pass BeginTx",
			fields: fields{
				makeConn: func() driver.Conn {
					conn := mocks.NewMockAdvancedConn(ctrl)
					conn.EXPECT().
						BeginTx(ctx, driver.TxOptions{Isolation: 1}).
						Times(1).
						Return(mocks.NewMockTx(ctrl), nil)

					return conn
				},
				makeHooks: func() Hook {
					hook := NewHooks(WithHook(&testHook{}))

					return hook
				},
			},
			args: args{
				opts: driver.TxOptions{Isolation: 1},
			},
			want: &Tx{
				Tx:    mocks.NewMockTx(ctrl),
				hooks: NewHooks(WithHook(&testHook{})),
				ctx:   ctx,
			},
		},
		{
			name: "pass Begin",
			fields: fields{
				makeConn: func() driver.Conn {
					conn := mocks.NewMockConn(ctrl)
					conn.EXPECT().
						Begin().
						Times(1).
						Return(mocks.NewMockTx(ctrl), nil)

					return conn
				},
				makeHooks: func() Hook {
					hook := NewHooks(WithHook(&testHook{}))

					return hook
				},
			},
			args: args{
				opts: driver.TxOptions{Isolation: 1},
			},
			want: &Tx{
				Tx:    mocks.NewMockTx(ctrl),
				hooks: NewHooks(WithHook(&testHook{})),
				ctx:   ctx,
			},
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
			args: args{
				opts: driver.TxOptions{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "with after error",
			fields: fields{
				makeConn: func() driver.Conn {
					conn := mocks.NewMockAdvancedConn(ctrl)
					conn.EXPECT().
						BeginTx(ctx, driver.TxOptions{}).
						Times(1).
						Return(mocks.NewMockTx(ctrl), nil)

					return conn
				},
				makeHooks: func() Hook {
					hook := NewHooks(WithHooksAfter(&testAfterHookWithError{}))

					return hook
				},
			},
			args: args{
				opts: driver.TxOptions{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "with error",
			fields: fields{
				makeConn: func() driver.Conn {
					conn := mocks.NewMockAdvancedConn(ctrl)
					conn.EXPECT().
						BeginTx(ctx, driver.TxOptions{}).
						Times(1).
						Return(nil, errors.New("some error"))

					return conn
				},
				makeHooks: func() Hook {
					return nil
				},
			},
			args: args{
				opts: driver.TxOptions{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "with error with hook",
			fields: fields{
				makeConn: func() driver.Conn {
					conn := mocks.NewMockAdvancedConn(ctrl)
					conn.EXPECT().
						BeginTx(ctx, driver.TxOptions{}).
						Times(1).
						Return(nil, errors.New("some error"))

					return conn
				},
				makeHooks: func() Hook {
					hook := NewHooks(WithHooksError(&testErrorHookWithError{}))

					return hook
				},
			},
			args: args{
				opts: driver.TxOptions{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "with fix error",
			fields: fields{
				makeConn: func() driver.Conn {
					conn := mocks.NewMockAdvancedConn(ctrl)
					conn.EXPECT().
						BeginTx(ctx, driver.TxOptions{ReadOnly: true}).
						Times(1).
						Return(nil, errors.New("some error"))

					return conn
				},
				makeHooks: func() Hook {
					hook := NewHooks(WithHooksError(&testFixErrorHookWithError{}))

					return hook
				},
			},
			args: args{
				opts: driver.TxOptions{ReadOnly: true},
			},
			want: &Tx{
				Tx:    nil,
				hooks: NewHooks(WithHooksError(&testFixErrorHookWithError{})),
				ctx:   ctx,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			conn := &Conn{
				Conn:  tt.fields.makeConn(),
				hooks: tt.fields.makeHooks(),
			}

			got, err := conn.BeginTx(ctx, tt.args.opts)
			switch tt.wantErr {
			case true:
				assert.Error(t, err)
			case false:
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestConn_Close(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	type fields struct {
		makeConn func() driver.Conn
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "pass",
			fields: fields{
				makeConn: func() driver.Conn {
					conn := mocks.NewMockAdvancedConn(ctrl)
					conn.EXPECT().
						Close().
						Times(1).
						Return(nil)

					return conn
				},
			},
		},
		{
			name: "error",
			fields: fields{
				makeConn: func() driver.Conn {
					conn := mocks.NewMockAdvancedConn(ctrl)
					conn.EXPECT().
						Close().
						Times(1).
						Return(errors.New("some error"))

					return conn
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			conn := &Conn{
				Conn: tt.fields.makeConn(),
			}

			err := conn.Close()
			switch tt.wantErr {
			case true:
				assert.Error(t, err)
			case false:
				assert.NoError(t, err)
			}
		})
	}
}

func TestConn_CheckNamedValue(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	type fields struct {
		makeConn func() driver.Conn
	}

	type args struct {
		nv *driver.NamedValue
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "pass with CheckNamedValue",
			fields: fields{
				makeConn: func() driver.Conn {
					conn := mocks.NewMockAdvancedConn(ctrl)
					conn.EXPECT().
						CheckNamedValue(&driver.NamedValue{Name: "test"}).
						Times(1).
						Return(nil)

					return conn
				},
			},
			args: args{
				nv: &driver.NamedValue{Name: "test"},
			},
		},
		{
			name: "error with CheckNamedValue",
			fields: fields{
				makeConn: func() driver.Conn {
					conn := mocks.NewMockAdvancedConn(ctrl)
					conn.EXPECT().
						CheckNamedValue(&driver.NamedValue{Name: "test"}).
						Times(1).
						Return(errors.New("some error"))

					return conn
				},
			},
			args: args{
				nv: &driver.NamedValue{Name: "test"},
			},
			wantErr: true,
		},
		{
			name: "pass without CheckNamedValue",
			fields: fields{
				makeConn: func() driver.Conn {
					conn := mocks.NewMockConn(ctrl)

					return conn
				},
			},
			args: args{
				nv: &driver.NamedValue{Name: "test"},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			conn := &Conn{
				Conn: tt.fields.makeConn(),
			}

			err := conn.CheckNamedValue(tt.args.nv)
			switch tt.wantErr {
			case true:
				assert.Error(t, err)
			case false:
				assert.NoError(t, err)
			}
		})
	}
}
