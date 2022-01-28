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

func TestTx_Commit(t *testing.T) {
	t.Parallel()

	var (
		ctx  = context.Background()
		ctrl = gomock.NewController(t)
	)

	type fields struct {
		makeTx    func() driver.Tx
		makeHooks func() Hook
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "with hooks",
			fields: fields{
				makeTx: func() driver.Tx {
					tx := mocks.NewMockTx(ctrl)
					tx.EXPECT().
						Commit().
						Times(1).
						DoAndReturn(func() error {
							return nil
						})

					return tx
				},
				makeHooks: func() Hook {
					hook := NewHooks(WithHooksError(&testFixErrorHookWithError{}))

					return hook
				},
			},
		},
		{
			name: "with before error",
			fields: fields{
				makeTx: func() driver.Tx {
					return nil
				},
				makeHooks: func() Hook {
					hook := NewHooks(WithHooksBefore(&testBeforeHookWithError{}))

					return hook
				},
			},
			wantErr: true,
		},
		{
			name: "with after error",
			fields: fields{
				makeTx: func() driver.Tx {
					tx := mocks.NewMockTx(ctrl)
					tx.EXPECT().
						Commit().
						Times(1).
						DoAndReturn(func() error {
							return nil
						})

					return tx
				},
				makeHooks: func() Hook {
					hook := NewHooks(WithHooksAfter(&testAfterHookWithError{}))

					return hook
				},
			},
			wantErr: true,
		},
		{
			name: "with error without hooks",
			fields: fields{
				makeTx: func() driver.Tx {
					tx := mocks.NewMockTx(ctrl)
					tx.EXPECT().
						Commit().
						Times(1).
						DoAndReturn(func() error {
							return errors.New("some error")
						})

					return tx
				},
				makeHooks: func() Hook {
					return nil
				},
			},
			wantErr: true,
		},
		{
			name: "with error",
			fields: fields{
				makeTx: func() driver.Tx {
					tx := mocks.NewMockTx(ctrl)
					tx.EXPECT().
						Commit().
						Times(1).
						DoAndReturn(func() error {
							return errors.New("some error")
						})

					return tx
				},
				makeHooks: func() Hook {
					hook := NewHooks(WithHooksError(&testErrorHookWithError{}))

					return hook
				},
			},
			wantErr: true,
		},
		{
			name: "with fix error",
			fields: fields{
				makeTx: func() driver.Tx {
					tx := mocks.NewMockTx(ctrl)
					tx.EXPECT().
						Commit().
						Times(1).
						DoAndReturn(func() error {
							return errors.New("some error")
						})

					return tx
				},
				makeHooks: func() Hook {
					hook := NewHooks(WithHooksError(&testFixErrorHookWithError{}))

					return hook
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tx := &Tx{
				Tx:    tt.fields.makeTx(),
				hooks: tt.fields.makeHooks(),
				ctx:   ctx,
			}

			err := tx.Commit()
			switch tt.wantErr {
			case true:
				assert.Error(t, err)
			case false:
				assert.NoError(t, err)
			}
		})
	}
}

func TestTx_Rollback(t *testing.T) {
	t.Parallel()

	var (
		ctx  = context.Background()
		ctrl = gomock.NewController(t)
	)

	type fields struct {
		makeTx    func() driver.Tx
		makeHooks func() Hook
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "with hooks",
			fields: fields{
				makeTx: func() driver.Tx {
					tx := mocks.NewMockTx(ctrl)
					tx.EXPECT().
						Rollback().
						Times(1).
						DoAndReturn(func() error {
							return nil
						})

					return tx
				},
				makeHooks: func() Hook {
					hook := NewHooks(WithHooksError(&testFixErrorHookWithError{}))

					return hook
				},
			},
		},
		{
			name: "with before error",
			fields: fields{
				makeTx: func() driver.Tx {
					return nil
				},
				makeHooks: func() Hook {
					hook := NewHooks(WithHooksBefore(&testBeforeHookWithError{}))

					return hook
				},
			},
			wantErr: true,
		},
		{
			name: "with after error",
			fields: fields{
				makeTx: func() driver.Tx {
					tx := mocks.NewMockTx(ctrl)
					tx.EXPECT().
						Rollback().
						Times(1).
						DoAndReturn(func() error {
							return nil
						})

					return tx
				},
				makeHooks: func() Hook {
					hook := NewHooks(WithHooksAfter(&testAfterHookWithError{}))

					return hook
				},
			},
			wantErr: true,
		},
		{
			name: "with error without hooks",
			fields: fields{
				makeTx: func() driver.Tx {
					tx := mocks.NewMockTx(ctrl)
					tx.EXPECT().
						Rollback().
						Times(1).
						DoAndReturn(func() error {
							return errors.New("some error")
						})

					return tx
				},
				makeHooks: func() Hook {
					return nil
				},
			},
			wantErr: true,
		},
		{
			name: "with error",
			fields: fields{
				makeTx: func() driver.Tx {
					tx := mocks.NewMockTx(ctrl)
					tx.EXPECT().
						Rollback().
						Times(1).
						DoAndReturn(func() error {
							return errors.New("some error")
						})

					return tx
				},
				makeHooks: func() Hook {
					hook := NewHooks(WithHooksError(&testErrorHookWithError{}))

					return hook
				},
			},
			wantErr: true,
		},
		{
			name: "with fix error",
			fields: fields{
				makeTx: func() driver.Tx {
					tx := mocks.NewMockTx(ctrl)
					tx.EXPECT().
						Rollback().
						Times(1).
						DoAndReturn(func() error {
							return errors.New("some error")
						})

					return tx
				},
				makeHooks: func() Hook {
					hook := NewHooks(WithHooksError(&testFixErrorHookWithError{}))

					return hook
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tx := &Tx{
				Tx:    tt.fields.makeTx(),
				hooks: tt.fields.makeHooks(),
				ctx:   ctx,
			}

			err := tx.Rollback()
			switch tt.wantErr {
			case true:
				assert.Error(t, err)
			case false:
				assert.NoError(t, err)
			}
		})
	}
}
