package uhttp

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJoinPath(t *testing.T) {
	type args struct {
		baseUrl string
		elem    []string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "with space",
			args: args{
				baseUrl: "",
				elem:    []string{"a", "b", "c"},
			},
			want:    "a/b/c",
			wantErr: false,
		},
		{
			name: "with /",
			args: args{
				baseUrl: "http://foo/",
				elem:    []string{"a", "b", "c"},
			},
			want:    "http://foo/a/b/c",
			wantErr: false,
		},
		{
			name: "without /",
			args: args{
				baseUrl: "http://foo",
				elem:    []string{"a", "b", "c"},
			},
			want:    "http://foo/a/b/c",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := JoinPath(tt.args.baseUrl, tt.args.elem...)
			if (err != nil) != tt.wantErr {
				t.Errorf("JoinPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("JoinPath() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTLDParser(t *testing.T) {
	type args struct {
		urlStr string
	}
	tests := []struct {
		name       string
		args       args
		wantDomain string
		wantErr    bool
	}{
		{
			name:       "empty url",
			args:       args{urlStr: ""},
			wantDomain: "",
			wantErr:    false,
		},
		{
			name:       "invalid host",
			args:       args{urlStr: "https://mktsocom/api/uc?k=v"},
			wantDomain: "",
			wantErr:    true,
		},
		{
			name:       "domain",
			args:       args{urlStr: "mktso.com"},
			wantDomain: "mktso.com",
			wantErr:    false,
		},
		{
			name:       "not a url",
			args:       args{urlStr: "mktso"},
			wantDomain: "mktso",
			wantErr:    false,
		},
		{
			name:       "domain",
			args:       args{urlStr: "https://mktso.com/api/uc?k=v"},
			wantDomain: "mktso.com",
			wantErr:    false,
		},
		{
			name:       "sub domain",
			args:       args{urlStr: "https://api.mktso.com/api/uc?k=v"},
			wantDomain: "mktso.com",
			wantErr:    false,
		},
		{
			name:       "port",
			args:       args{urlStr: "https://api.mktso.com:8080/api/uc?k=v"},
			wantDomain: "mktso.com",
			wantErr:    false,
		},
		{
			name:       "null",
			args:       args{urlStr: ""},
			wantDomain: "",
			wantErr:    false,
		},
		{
			name:       "not found www",
			args:       args{urlStr: "https://vodafone.com.au/"},
			wantDomain: "vodafone.com.au",
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tld, err := TLDParser(tt.args.urlStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("TLDParser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			ok := assert.Equal(t, tld.Domain, tt.wantDomain)
			if !ok {
				t.Errorf("TLDParser() gotDomain = %v, want %v", tld.Domain, tt.wantDomain)
			}
		})
	}
}
