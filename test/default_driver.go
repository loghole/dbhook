package test

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"io"

	"github.com/golang/mock/gomock"

	"github.com/loghole/dbhook/mocks"
)

type queryerConn struct {
	driver.Conn
	driver.ConnBeginTx
	driver.QueryerContext
	driver.ExecerContext
}

func MakeDefaultDriver(ctrl *gomock.Controller, name string) driver.Driver {
	drv := mocks.NewMockDriver(ctrl)

	drv.EXPECT().Open(gomock.Any()).AnyTimes().DoAndReturn(func(name string) (driver.Conn, error) {
		conn := mocks.NewMockConn(ctrl)

		queryer := mocks.NewMockQueryerContext(ctrl)
		queryer.EXPECT().
			QueryContext(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().
			DoAndReturn(func(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
				rows := mocks.NewMockRows(ctrl)

				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Columns().AnyTimes().DoAndReturn(func() []string {
					return []string{"id"}
				})

				var rowsCallCounter int

				rows.EXPECT().Next(gomock.Any()).AnyTimes().SetArg(0, []driver.Value{"some"}).
					DoAndReturn(func(args []driver.Value) error {
						if rowsCallCounter >= 2 { // nolint:gomnd // it's test
							return io.EOF
						}

						rowsCallCounter++

						return nil
					})

				return rows, nil
			})

		execer := mocks.NewMockExecerContext(ctrl)
		execer.EXPECT().
			ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().
			DoAndReturn(func(arg0 context.Context, arg1 string, arg2 []driver.NamedValue) (driver.Result, error) {
				result := mocks.NewMockResult(ctrl)

				return result, nil
			})

		return &queryerConn{conn, nil, queryer, execer}, nil
	})

	sql.Register(name, drv)

	return drv
}
