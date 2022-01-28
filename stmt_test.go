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

type testHook struct{}

func (h *testHook) Before(ctx context.Context, input *HookInput) (context.Context, error) {
	return ctx, nil
}

func (h *testHook) After(ctx context.Context, input *HookInput) (context.Context, error) {
	return ctx, nil
}

func (h *testHook) Error(ctx context.Context, input *HookInput) (context.Context, error) {
	return ctx, input.Error
}

type testBeforeHookWithError struct{}

func (h *testBeforeHookWithError) Before(ctx context.Context, input *HookInput) (context.Context, error) {
	return ctx, errors.New("some error")
}

type testAfterHookWithError struct{}

func (h *testAfterHookWithError) After(ctx context.Context, input *HookInput) (context.Context, error) {
	return ctx, errors.New("some error")
}

type testErrorHookWithError struct{}

func (h *testErrorHookWithError) Error(ctx context.Context, input *HookInput) (context.Context, error) {
	return ctx, input.Error
}

type testFixErrorHookWithError struct{}

func (h *testFixErrorHookWithError) Error(ctx context.Context, input *HookInput) (context.Context, error) {
	return ctx, nil
}

func TestStmt_ExecContext(t *testing.T) {
	t.Parallel()

	var (
		ctx  = context.Background()
		ctrl = gomock.NewController(t)
	)

	type fields struct {
		makeStmt  func() driver.Stmt
		makeHooks func() Hook
	}

	type args struct {
		args []driver.NamedValue
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    driver.Result
		wantErr bool
	}{
		{
			name: "with driver.Stmt",
			fields: fields{
				makeStmt: func() driver.Stmt {
					stmt := mocks.NewMockStmt(ctrl)
					stmt.EXPECT().
						Exec([]driver.Value{}).
						Times(1).
						DoAndReturn(func(arg []driver.Value) (driver.Result, error) {
							result := mocks.NewMockResult(ctrl)

							return result, nil
						})

					return stmt
				},
				makeHooks: func() Hook {
					return nil
				},
			},
			want: mocks.NewMockResult(ctrl),
		},
		{
			name: "with driver.StmtExecContext",
			fields: fields{
				makeStmt: func() driver.Stmt {
					stmt := mocks.NewMockStmtWithContext(ctrl)
					stmt.EXPECT().
						ExecContext(ctx, nil).
						Times(1).
						DoAndReturn(func(ctx context.Context, arg []driver.NamedValue) (driver.Result, error) {
							result := mocks.NewMockResult(ctrl)

							return result, nil
						})

					return stmt
				},
				makeHooks: func() Hook {
					return nil
				},
			},
			want: mocks.NewMockResult(ctrl),
		},
		{
			name: "with error",
			fields: fields{
				makeStmt: func() driver.Stmt {
					stmt := mocks.NewMockStmtWithContext(ctrl)
					stmt.EXPECT().
						ExecContext(ctx, nil).
						Times(1).
						DoAndReturn(func(ctx context.Context, arg []driver.NamedValue) (driver.Result, error) {
							return nil, errors.New("some error")
						})

					return stmt
				},
				makeHooks: func() Hook {
					return nil
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "with hook",
			fields: fields{
				makeStmt: func() driver.Stmt {
					stmt := mocks.NewMockStmtWithContext(ctrl)
					stmt.EXPECT().
						ExecContext(ctx, nil).
						Times(1).
						DoAndReturn(func(ctx context.Context, arg []driver.NamedValue) (driver.Result, error) {
							result := mocks.NewMockResult(ctrl)

							return result, nil
						})

					return stmt
				},
				makeHooks: func() Hook {
					hook := NewHooks(WithHook(&testHook{}))

					return hook
				},
			},
			want: mocks.NewMockResult(ctrl),
		},
		{
			name: "with before error",
			fields: fields{
				makeStmt: func() driver.Stmt {
					return nil
				},
				makeHooks: func() Hook {
					hook := NewHooks(WithHooksBefore(&testBeforeHookWithError{}))

					return hook
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "with after error",
			fields: fields{
				makeStmt: func() driver.Stmt {
					stmt := mocks.NewMockStmtWithContext(ctrl)
					stmt.EXPECT().
						ExecContext(ctx, nil).
						Times(1).
						DoAndReturn(func(ctx context.Context, arg []driver.NamedValue) (driver.Result, error) {
							result := mocks.NewMockResult(ctrl)

							return result, nil
						})

					return stmt
				},
				makeHooks: func() Hook {
					hook := NewHooks(WithHooksAfter(&testAfterHookWithError{}))

					return hook
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "with error and after hook",
			fields: fields{
				makeStmt: func() driver.Stmt {
					stmt := mocks.NewMockStmtWithContext(ctrl)
					stmt.EXPECT().
						ExecContext(ctx, nil).
						Times(1).
						DoAndReturn(func(ctx context.Context, arg []driver.NamedValue) (driver.Result, error) {
							return nil, errors.New("some error")
						})

					return stmt
				},
				makeHooks: func() Hook {
					hook := NewHooks(WithHooksAfter(&testAfterHookWithError{}))

					return hook
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "with fix error",
			fields: fields{
				makeStmt: func() driver.Stmt {
					stmt := mocks.NewMockStmtWithContext(ctrl)
					stmt.EXPECT().
						ExecContext(ctx, nil).
						Times(1).
						DoAndReturn(func(ctx context.Context, arg []driver.NamedValue) (driver.Result, error) {
							return nil, errors.New("some error")
						})

					return stmt
				},
				makeHooks: func() Hook {
					hook := NewHooks(WithHooksError(&testFixErrorHookWithError{}))

					return hook
				},
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			stmt := &Stmt{
				Stmt:  tt.fields.makeStmt(),
				hooks: tt.fields.makeHooks(),
			}

			got, err := stmt.ExecContext(ctx, tt.args.args)
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

func TestStmt_Close(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	type fields struct {
		makeStmt func() driver.Stmt
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "pass",
			fields: fields{
				makeStmt: func() driver.Stmt {
					stmt := mocks.NewMockStmtWithContext(ctrl)
					stmt.EXPECT().
						Close().
						DoAndReturn(func() error {
							return nil
						})

					return stmt
				},
			},
		},
		{
			name: "some error",
			fields: fields{
				makeStmt: func() driver.Stmt {
					stmt := mocks.NewMockStmtWithContext(ctrl)
					stmt.EXPECT().
						Close().
						DoAndReturn(func() error {
							return errors.New("some error")
						})

					return stmt
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			stmt := &Stmt{Stmt: tt.fields.makeStmt()}

			err := stmt.Close()
			switch tt.wantErr {
			case true:
				assert.Error(t, err)
			case false:
				assert.NoError(t, err)
			}
		})
	}
}

func TestStmt_NumInput(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	type fields struct {
		makeStmt func() driver.Stmt
	}

	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "pass",
			fields: fields{
				makeStmt: func() driver.Stmt {
					stmt := mocks.NewMockStmtWithContext(ctrl)
					stmt.EXPECT().
						NumInput().
						DoAndReturn(func() int {
							return 5
						})

					return stmt
				},
			},
			want: 5,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			stmt := &Stmt{Stmt: tt.fields.makeStmt()}

			got := stmt.NumInput()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestStmt_QueryContext(t *testing.T) {
	t.Parallel()

	var (
		ctx  = context.Background()
		ctrl = gomock.NewController(t)
	)

	type fields struct {
		makeStmt  func() driver.Stmt
		makeHooks func() Hook
	}

	type args struct {
		args []driver.NamedValue
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
				makeStmt: func() driver.Stmt {
					stmt := mocks.NewMockStmtWithContext(ctrl)
					stmt.EXPECT().
						QueryContext(ctx, nil).
						Times(1).
						DoAndReturn(func(ctx context.Context, arg []driver.NamedValue) (driver.Rows, error) {
							rows := mocks.NewMockRows(ctrl)

							return rows, nil
						})

					return stmt
				},
				makeHooks: func() Hook {
					hook := NewHooks(WithHook(&testHook{}))

					return hook
				},
			},
			want: mocks.NewMockRows(ctrl),
		},
		{
			name: "pass with Query",
			fields: fields{
				makeStmt: func() driver.Stmt {
					stmt := mocks.NewMockStmt(ctrl)
					stmt.EXPECT().
						Query([]driver.Value{}).
						Times(1).
						DoAndReturn(func(arg []driver.Value) (driver.Rows, error) {
							rows := mocks.NewMockRows(ctrl)

							return rows, nil
						})

					return stmt
				},
				makeHooks: func() Hook {
					hook := NewHooks(WithHook(&testHook{}))

					return hook
				},
			},
			want: mocks.NewMockRows(ctrl),
		},
		{
			name: "with before error",
			fields: fields{
				makeStmt: func() driver.Stmt {
					return nil
				},
				makeHooks: func() Hook {
					hook := NewHooks(WithHooksBefore(&testBeforeHookWithError{}))

					return hook
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "with after error",
			fields: fields{
				makeStmt: func() driver.Stmt {
					stmt := mocks.NewMockStmtWithContext(ctrl)
					stmt.EXPECT().
						QueryContext(ctx, nil).
						Times(1).
						DoAndReturn(func(ctx context.Context, arg []driver.NamedValue) (driver.Rows, error) {
							rows := mocks.NewMockRows(ctrl)

							return rows, nil
						})

					return stmt
				},
				makeHooks: func() Hook {
					hook := NewHooks(WithHooksAfter(&testAfterHookWithError{}))

					return hook
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "with error",
			fields: fields{
				makeStmt: func() driver.Stmt {
					stmt := mocks.NewMockStmtWithContext(ctrl)
					stmt.EXPECT().
						QueryContext(ctx, nil).
						Times(1).
						DoAndReturn(func(ctx context.Context, arg []driver.NamedValue) (driver.Rows, error) {
							return nil, errors.New("some error")
						})

					return stmt
				},
				makeHooks: func() Hook {
					hook := NewHooks(WithHooksError(&testErrorHookWithError{}))

					return hook
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "with error without hooks",
			fields: fields{
				makeStmt: func() driver.Stmt {
					stmt := mocks.NewMockStmtWithContext(ctrl)
					stmt.EXPECT().
						QueryContext(ctx, nil).
						Times(1).
						DoAndReturn(func(ctx context.Context, arg []driver.NamedValue) (driver.Rows, error) {
							return nil, errors.New("some error")
						})

					return stmt
				},
				makeHooks: func() Hook {
					return nil
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "with fix error",
			fields: fields{
				makeStmt: func() driver.Stmt {
					stmt := mocks.NewMockStmtWithContext(ctrl)
					stmt.EXPECT().
						QueryContext(ctx, nil).
						Times(1).
						DoAndReturn(func(ctx context.Context, arg []driver.NamedValue) (driver.Rows, error) {
							return nil, errors.New("some error")
						})

					return stmt
				},
				makeHooks: func() Hook {
					hook := NewHooks(WithHooksError(&testFixErrorHookWithError{}))

					return hook
				},
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			stmt := &Stmt{
				Stmt:  tt.fields.makeStmt(),
				hooks: tt.fields.makeHooks(),
			}

			got, err := stmt.QueryContext(ctx, tt.args.args)
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
