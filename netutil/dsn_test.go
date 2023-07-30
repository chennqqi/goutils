package netutil

import (
	"reflect"
	"testing"
)

func TestParseDsn(t *testing.T) {
	type args struct {
		dsn string
	}
	tests := []struct {
		name    string
		args    args
		want    Dsn
		wantErr bool
	}{
		// Add test cases.、
		{name: "test0", args: args{dsn: "mysql://root:z!guwrBhH7p>(127.0.0.1:3306)/cloudscan?charset=utf8mb4&parseTime=True&loc=Local"}, want: Dsn{Scheme: "mysql", Source: "root:z!guwrBhH7p>(127.0.0.1:3306)/cloudscan?charset=utf8mb4&parseTime=True&loc=Local"}, wantErr: false},
		{name: "test1", args: args{dsn: ""}, want: Dsn{Scheme: "", Source: ""}, wantErr: true},
		{name: "test2", args: args{dsn: "mongodb://testuser:psKdqSe@192.168.102.153:27017"}, want: Dsn{Scheme: "mongodb", Source: "testuser:psKdqSe@192.168.102.153:27017"}, wantErr: false},
		{name: "test3", args: args{dsn: "mongodb://192.168.102.153:27017"}, want: Dsn{Scheme: "mongodb", Source: "192.168.102.153:27017"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDsn(tt.args.dsn)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDsn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseDsn() got = %v, want %v", got, tt.want)
			}
		})
	}
}
