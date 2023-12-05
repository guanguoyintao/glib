package umachine

import (
	"bufio"
	"context"
	"fmt"
	"git.umu.work/AI/uglib/uerrors"
	"git.umu.work/AI/uglib/uregexp"
	"git.umu.work/be/goframework/logger"
	"net"
	"os"
	"strings"
)

func GetMachineIP(ctx context.Context) (string, error) {
	file, err := os.Open("/etc/hosts")
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[0] != "127.0.0.1" && len(uregexp.IPv4s(fields[0])) > 0 {
			logger.GetLogger(ctx).Info(fmt.Sprintf("machine ip %+v", fields[0]))
			return fields[0], nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", nil
}

func IPv4ToInt64(ctx context.Context, ipv4 string) (int64, error) {
	parsedIP := net.ParseIP(ipv4)
	if parsedIP == nil {
		return 0, uerrors.UErrorInvalidIP
	}

	ipv4Bytes := parsedIP.To4()
	if ipv4Bytes == nil {
		return 0, uerrors.UErrorInvalidIP
	}

	ipv4Int := int64(ipv4Bytes[0])<<24 +
		int64(ipv4Bytes[1])<<16 +
		int64(ipv4Bytes[2])<<8 +
		int64(ipv4Bytes[3])

	return ipv4Int, nil
}
