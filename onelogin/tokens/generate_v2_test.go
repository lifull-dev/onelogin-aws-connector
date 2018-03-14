package tokens

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"
)

func TestTokens_Generate(t *testing.T) {
	type fields struct {
		Endpoint     string
		ClientToken  string
		ClientSecret string
	}
	type response struct {
		code int
		body string
	}
	f := fields{
		Endpoint:     "",
		ClientToken:  "client-token",
		ClientSecret: "client-secret,",
	}
	now := time.Now().UTC()
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	tests := []struct {
		name    string
		fields  fields
		req     *GenerateRequest
		res     response
		want    *GenerateResponse
		wantErr bool
	}{
		{
			name:   "success",
			fields: f,
			req: &GenerateRequest{
				GrantType: "client_credentials",
			},
			res: response{
				code: 200,
				body: fmt.Sprintf(`{
					"access_token": "access-token",
					"created_at": "%s",
					"expires_in": 3600,
					"refresh_token": "refresh-token",
					"token_type": "token-type",
					"account_id": 1234567
				}`, now.Format("2006-01-02T15:04:05.000Z")),
			},
			want: &GenerateResponse{
				AccessToken:  "access-token",
				CreatedAt:    now.Format("2006-01-02T15:04:05.000Z"),
				ExpiresIn:    3600,
				RefreshToken: "refresh-token",
				TokenType:    "token-type",
				AccountID:    1234567,
			},
			wantErr: false,
		},
		{
			name:   "failed",
			fields: f,
			req: &GenerateRequest{
				GrantType: "client_credentials",
			},
			res: response{
				code: 200,
				body: `{
					"status": {
						"error": true,
						"code": 400,
						"type": "bad request",
						"message": "grant_type is incorrect/absent"
					}
				}`,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				defer r.Body.Close()
				body, err := ioutil.ReadAll(r.Body)
				if err != nil {
					t.Errorf("%v", err)
				}
				var input GenerateRequest
				if err := json.Unmarshal(body, &input); err != nil {
					t.Errorf("%v", err)
				}
				if !reflect.DeepEqual(&input, tt.req) {
					t.Errorf("Tokens.Generate() = %#v, want %#v", &input, tt.req)
				}
				w.Header().Set("Content-Type", "text/plain; charset=utf-8")
				w.Header().Set("X-Content-Type-Options", "nosniff")
				w.WriteHeader(tt.res.code)
				fmt.Fprintln(w, bytes.NewBuffer([]byte(tt.res.body)))
			}))
			defer ts.Close()
			u, _ := url.Parse(ts.URL)
			endpoint := fmt.Sprintf("%s:%s", u.Hostname(), u.Port())
			tt.fields.Endpoint = endpoint
			g := &Tokens{
				Endpoint:     tt.fields.Endpoint,
				ClientToken:  tt.fields.ClientToken,
				ClientSecret: tt.fields.ClientSecret,
				HTTPClient:   httpClient,
			}
			got, err := g.Generate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Tokens.Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tokens.Generate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokens_Refresh(t *testing.T) {
	type fields struct {
		Endpoint     string
		ClientToken  string
		ClientSecret string
	}
	type args struct {
		input *RefreshRequest
	}
	type response struct {
		code int
		body string
	}
	f := fields{
		Endpoint:     "",
		ClientToken:  "client-token",
		ClientSecret: "client-secret,",
	}
	now := time.Now().UTC()
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		req     *RefreshRequest
		res     response
		want    *GenerateResponse
		wantErr bool
	}{
		{
			name:   "success",
			fields: f,
			args: args{
				input: &RefreshRequest{
					AccessToken:  "access-token",
					RefreshToken: "refresh-token",
				},
			},
			req: &RefreshRequest{
				GrantType:    "refresh_token",
				AccessToken:  "access-token",
				RefreshToken: "refresh-token",
			},
			res: response{
				code: 200,
				body: fmt.Sprintf(`{
					"access_token": "access-token",
					"created_at": "%s",
					"expires_in": 3600,
					"refresh_token": "refresh-token",
					"token_type": "token-type",
					"account_id": 1234567
				}`, now.Format("2006-01-02T15:04:05.000Z")),
			},
			want: &GenerateResponse{
				AccessToken:  "access-token",
				CreatedAt:    now.Format("2006-01-02T15:04:05.000Z"),
				ExpiresIn:    3600,
				RefreshToken: "refresh-token",
				TokenType:    "token-type",
				AccountID:    1234567,
			},
			wantErr: false,
		},
		{
			name:   "failed",
			fields: f,
			args: args{
				input: &RefreshRequest{
					AccessToken:  "access-token",
					RefreshToken: "refresh-token",
				},
			},
			req: &RefreshRequest{
				GrantType:    "refresh_token",
				AccessToken:  "access-token",
				RefreshToken: "refresh-token",
			},
			res: response{
				code: 200,
				body: `{
					"status": {
						"error": true,
						"code": 400,
						"type": "bad request",
						"message": "grant_type is incorrect/absent"
					}
				}`,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				defer r.Body.Close()
				body, err := ioutil.ReadAll(r.Body)
				if err != nil {
					t.Errorf("%v", err)
				}
				var input RefreshRequest
				if err := json.Unmarshal(body, &input); err != nil {
					t.Errorf("%v", err)
				}
				if !reflect.DeepEqual(&input, tt.req) {
					t.Errorf("Tokens.Generate() = %#v, want %#v", &input, tt.req)
				}
				w.Header().Set("Content-Type", "text/plain; charset=utf-8")
				w.Header().Set("X-Content-Type-Options", "nosniff")
				w.WriteHeader(tt.res.code)
				fmt.Fprintln(w, bytes.NewBuffer([]byte(tt.res.body)))
			}))
			defer ts.Close()
			u, _ := url.Parse(ts.URL)
			endpoint := fmt.Sprintf("%s:%s", u.Hostname(), u.Port())
			tt.fields.Endpoint = endpoint
			g := &Tokens{
				Endpoint:     tt.fields.Endpoint,
				ClientToken:  tt.fields.ClientToken,
				ClientSecret: tt.fields.ClientSecret,
				HTTPClient:   httpClient,
			}
			got, err := g.Refresh(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Tokens.Refresh() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tokens.Refresh() = %v, want %v", got, tt.want)
			}
		})
	}
}
