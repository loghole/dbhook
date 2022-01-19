package dbhook

import (
	"database/sql/driver"
	"reflect"
	"testing"

	sqlite "github.com/mattn/go-sqlite3"
)

func TestWrap(t *testing.T) {
	t.Parallel()

	type args struct {
		drv driver.Driver
		hks Hook
	}

	tests := []struct {
		name string
		args args
		want driver.Driver
	}{
		{
			name: "pass",
			args: args{
				drv: &sqlite.SQLiteDriver{},
				hks: &testHook{},
			},
			want: &Driver{Driver: &sqlite.SQLiteDriver{}, hooks: &testHook{}},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := Wrap(tt.args.drv, tt.args.hks); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Wrap() = %v, want %v", got, tt.want)
			}
		})
	}
}
