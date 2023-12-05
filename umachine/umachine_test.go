package umachine

import (
	"context"
	"testing"
)

func TestGetMachineIP(t *testing.T) {
	// 创建一个带取消功能的上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ip, err := GetMachineIP(ctx)
	if err != nil {
		t.Errorf("GetMachineIP() error: %v", err)
		return
	}

	// 确认 IP 地址不为空
	if ip == "" {
		t.Errorf("GetMachineIP() returned empty IP address")
	}
}

func TestIPv4ToInt64(t *testing.T) {
	tests := []struct {
		name    string
		ipv4    string
		want    int64
		wantErr bool
	}{
		{
			name:    "Valid IPv4",
			ipv4:    "192.168.0.1",
			want:    3232235521,
			wantErr: false,
		},
		{
			name:    "Invalid IPv4",
			ipv4:    "256.0.0.1",
			want:    0,
			wantErr: true,
		},
		{
			name:    "Loopback IPv4",
			ipv4:    "127.0.0.1",
			want:    2130706433,
			wantErr: false,
		},
	}
	// 用来存储已经转换过的 int64 值
	convertedValues := make(map[int64]bool)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IPv4ToInt64(context.TODO(), tt.ipv4)
			if (err != nil) != tt.wantErr {
				t.Errorf("IPv4ToInt64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IPv4ToInt64() got = %v, want %v", got, tt.want)
			}

			// 验证转换后的 int64 值的唯一性
			if _, ok := convertedValues[got]; ok {
				t.Errorf("Duplicate int64 value: %v", got)
			} else {
				convertedValues[got] = true
			}
		})
	}
}
