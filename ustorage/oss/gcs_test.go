package uoss

import (
	"context"
	"fmt"
	"git.umu.work/AI/uglib/uhash"
	"path"
	"testing"
)

func TestGCSClient_UploadByWavFile(t *testing.T) {

	type args struct {
		ctx           context.Context
		localFilePath string
		audioName     string
		gcsFilePath   string
		credentials   string
		bucket        string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "en",
			args: args{
				ctx:           context.Background(),
				localFilePath: "/Users/guanguoyintao/data/en.wav",
				audioName:     "en.wav",
				gcsFilePath:   "asr/en",
				credentials: `{
				"type": "service_account",
				"project_id": "umu-gapi",
				"private_key_id": "6a71f886c114fe21f5d1bc7861197ab58e694675",
				"private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDAEXM/dnJBg6eV\nkDJSnewiIeSlENUbq1iMwlZD+f/ZDlfVPqknJ1cmcH9KOzZJOlgwUI6XFWBgq2b4\nDCog9VU8bAfuhiaCHmGV5+EJXqA8B1J7uZg4EEKyh2Y1cZsK4wO39rY5+J0Z40Pb\ndA1z59w/SYAb7VJf0N26/a0XhO1jSetaZxYRKoerilyo8dd0oTfGDgbDpzjCAveo\ntF5LcwPTVnZFQsKVYxILFx11un1XczXNl15ofR4G4bkMW3bLEylYULvR2Q3HPEAj\nRh6Xwx9Qqki5Uhg5f+PhIq2x3aWiv3KkD0d6/s+FZEIcl8rgMoCRehB07dJYjs8V\nD0g9bMJ5AgMBAAECggEAAUCVeXsnUo465P3Y/iXoC+7s1tCtx9QaeMx5k7biq3Xt\nYUdBaUjXeuyUpkjAv8IUIIVRI+LW5O2pKoJCmHVqovRQ9uxD/jXLGbDQP4EXVvZ0\n2+t0kFHQa+nEr+63ypHhWW02RSZKuLtjgV/koilUWJgWydfTYV5ZjvMORY8V0L+Z\npBG2JGgyz2wb1wrcDqtKqvxzZWIwPl6btyHwCvKB1Xe7tEq0tm1qq4yGHjE8aJWu\n7++vI50tOSfq+jokcQbxkwXxZlUlS0WgNKdZYq+oQgwrjvSezcqqBVMKsJ7sqjvs\ndtlZJY1VejTWa03CgkC96h7Yj+u4bgAHhk41m1LRawKBgQDjW1WfdRsWWPSi37p6\nM3/wTd3YOgBCH5W4Ne+c1TFREoOJih0YLSkKNxvHr7iW2FOWTtclZ3friOu4MBbV\nDuuMdu7sYwOIr9ulsi4AlrWf3EYF+lfLx/ODWiQBMC85vMkJ3YHyJbHVCkJmoXOQ\nrPThD/zqOr+utfYMGBCjgG5RmwKBgQDYQ/6VoehAH4d9ThloriOImypC4wx/E79V\nxUtUrCWNJOUnjHhmcrGN/VLx5FmcExvV4fNmxvT/ulOf1mKL02i8O9ZiK3eA4xAn\n9shpeIy492Z5T0Gd+a96PG2PiDP4xigTmTEh/T45AQZIbztYcpilwlYZAcEhd043\nrJLbLX33ewKBgCyDMkVQ6Se77NGCmgDY2mCS7i8qU+ieRHLXZH1BJDGqPUSNXNrh\n5JoSZgb3eV1XJy5Taz3wfwMHQJdEGwRFsopCss5nKEb3nzpWozkMSKzutGrxM4U7\nNYru+AOfim90bavXyVw+Uw3Rv2RWkciuhBcYST0WmXa8O1rszz0jpPWnAoGBAJvJ\nz1km+tFGZSnE2tTnrL05WWY5fMUGCqzUpQFnyN02GMU2kXzrXjA0rk2F+29M7J5U\nWCxPYindpWPc4bxsHGSMjlDFkx6LcxX21kP0DRspTO5SQ6hzNw9M/HeK3DV3OWN0\n1UnwzaPuswH27r82PfQaXD9DXkanVPcSH/A3gyiFAoGAMiUNrhDYs7oOnc0LyRv0\nJrlRkyxO2BTu/hYmgG9WC1oHavsC/5dRLe/nQ1rdaaCeFhfZX+FgNS9/0qdS8yvF\nMZ+aTdd+zOfrHHpVZjALllx6oOKKXDFgVJefpXTujr9pzGiaESv73EsoT+rmK/0W\nMn7FhjeFQw9NyE0OZSgEjPs=\n-----END PRIVATE KEY-----\n",
				"client_email": "asr-poly@umu-gapi.iam.gserviceaccount.com",
				"client_id": "101150886029617617119",
				"auth_uri": "https://accounts.google.com/o/oauth2/auth",
				"token_uri": "https://oauth2.googleapis.com/token",
				"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
				"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/asr-poly%40umu-gapi.iam.gserviceaccount.com"
			}`,
				bucket: "asr-poly",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := NewGCSClient(tt.args.ctx, tt.args.credentials, tt.args.bucket)
			if err != nil {
				t.Error(err)
			}
			gucUrl, err := g.UploadByFile(tt.args.ctx, tt.args.localFilePath, path.Join(tt.args.gcsFilePath, uhash.HashMD532(tt.args.audioName)))
			if (err != nil) != tt.wantErr {
				t.Errorf("UploadByWavFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Printf("gcs url is %s\n", gucUrl)
		})
	}
}
