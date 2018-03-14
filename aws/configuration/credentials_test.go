package configuration

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestNewCredentials(t *testing.T) {
	type args struct {
		dir     string
		profile string
	}
	tests := []struct {
		name string
		args args
		want *Credentials
	}{
		{
			name: "NewCredentials",
			args: args{
				dir:     "/tmp",
				profile: "test",
			},
			want: &Credentials{
				file:    "/tmp/credentials",
				profile: "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCredentials(tt.args.dir, tt.args.profile); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCredentials() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCredentials_Save(t *testing.T) {
	type fields struct {
		file    string
		profile string
	}
	type args struct {
		options map[string]string
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantErr     bool
		wantContent string
	}{
		{
			name: "default empty profile",
			fields: fields{
				file:    "/tmp/testconfig",
				profile: "default",
			},
			args: args{
				options: map[string]string{},
			},
			wantErr: false,
			wantContent: `[default]

`,
		},
		{
			name: "default profile",
			fields: fields{
				file:    "/tmp/testconfig",
				profile: "default",
			},
			args: args{
				options: map[string]string{
					"aws_access_key": "12345678",
				},
			},
			wantErr: false,
			wantContent: `[default]
aws_access_key = 12345678

`,
		},
		{
			name: "Unwritable file",
			fields: fields{
				file:    "/unwritable",
				profile: "test",
			},
			args: args{
				options: map[string]string{},
			},
			wantErr:     true,
			wantContent: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := tt.fields.file
			c := &Credentials{
				file:    tt.fields.file,
				profile: tt.fields.profile,
			}
			err := c.Save(tt.args.options)
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
