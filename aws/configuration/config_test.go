package configuration

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestNewConfig(t *testing.T) {
	type args struct {
		dir     string
		profile string
	}
	tests := []struct {
		name string
		args args
		want *Config
	}{
		{
			name: "NewConfig",
			args: args{
				dir:     "/tmp",
				profile: "test",
			},
			want: &Config{
				file:    "/tmp/config",
				profile: "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewConfig(tt.args.dir, tt.args.profile); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_Save(t *testing.T) {
	type fields struct {
		file    string
		profile string
	}
	type args struct {
		region string
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantErr     bool
		wantContent string
	}{
		{
			name: "default profile",
			fields: fields{
				file:    "/tmp/testconfig",
				profile: "default",
			},
			args: args{
				region: "us-east-1",
			},
			wantErr: false,
			wantContent: `[profile default]
region = us-east-1

`,
		},
		{
			name: "test profile",
			fields: fields{
				file:    "/tmp/testconfig",
				profile: "test",
			},
			args: args{
				region: "ap-northeast-1",
			},
			wantErr: false,
			wantContent: `[profile test]
region = ap-northeast-1

`,
		},
		{
			name: "Unwritable file",
			fields: fields{
				file:    "/unwritable",
				profile: "test",
			},
			args: args{
				region: "ap-northeast-1",
			},
			wantErr:     true,
			wantContent: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := tt.fields.file
			os.Remove(file)
			c := &Config{
				file:    tt.fields.file,
				profile: tt.fields.profile,
			}
			err := c.Save(tt.args.region)
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Save() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil {
				if _, err := os.Stat(file); err != nil {
					t.Errorf("Config.Save() don't save %v", tt.fields.file)
				}
				data, err := ioutil.ReadFile(file)
				if err != nil {
					t.Errorf("%#v", err)
				}
				actual := string(data)
				expected := tt.wantContent
				if actual != expected {
					t.Errorf("'%v' is not equal '%v'", actual, expected)
				}
				os.Remove(tt.fields.file)
			}
		})
	}
}
