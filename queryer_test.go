package dbhook

import (
	"context"
	"database/sql/driver"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/loghole/dbhook/mocks"
)

// TODO fix
func TestQueryerContext_QueryContext(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish() // remove on gomock version >= 1.5.0

	connCtx := mocks.NewMockQueryerContext(ctrl)

	ctx := context.Background()

	connCtx.EXPECT().QueryContext(ctx, "any query", []driver.NamedValue{}).DoAndReturn(func(ctx context.Context, q string, args []driver.NamedValue) (driver.Rows, error) {
		return nil, nil
	})

	type args struct {
		ctx   context.Context
		query string
		args  []driver.NamedValue
	}
	tests := []struct {
		name    string
		args    args
		want    driver.Rows
		wantErr bool
	}{
		{
			name: "pass",
			args: args{
				ctx:   ctx,
				query: "any query",
				args:  []driver.NamedValue{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := connCtx.QueryContext(tt.args.ctx, tt.args.query, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryerContext.QueryContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QueryerContext.QueryContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueryerContext_queryContext(t *testing.T) {
	type fields struct {
		Conn *Conn
	}
	type args struct {
		ctx   context.Context
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn := &QueryerContext{
				Conn: tt.fields.Conn,
			}
			got, err := conn.queryContext(tt.args.ctx, tt.args.query, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryerContext.queryContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QueryerContext.queryContext() = %v, want %v", got, tt.want)
			}
		})
	}
}
