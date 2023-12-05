package usrt

import (
	"bufio"
	"context"
	"fmt"
	"git.umu.work/AI/uglib/usubtitle"
	"git.umu.work/be/goframework/logger"
	"github.com/asticode/go-astisub"
	"os"
	"strings"
)

func ParserFromFile(ctx context.Context, srtPath string) ([]*usubtitle.Item, error) {
	// 打开SRT字幕文件
	file, err := os.Open(srtPath)
	if err != nil {
		logger.GetLogger(ctx).Warn(err.Error())
		return nil, err
	}
	defer file.Close()
	// 读取SRT文件
	subtitles, err := astisub.ReadFromSRT(file)
	if err != nil {
		logger.GetLogger(ctx).Warn(err.Error())
		return nil, err
	}

	// 遍历字幕并输出内容
	items := make([]*usubtitle.Item, 0, len(subtitles.Items))
	for _, item := range subtitles.Items {
		texts := make([]string, 0)
		for _, line := range item.Lines {
			texts = append(texts, line.String())
		}
		items = append(items, &usubtitle.Item{
			EndAt:   item.EndAt,
			StartAt: item.StartAt,
			Texts:   texts,
		})
	}

	return items, nil
}

func WriteFile(ctx context.Context, srtPath string, items []*usubtitle.Item) error {
	var sb strings.Builder
	lineNum := 0
	for _, item := range items {
		lineNum++
		startTime := formatSRTTime(item.StartAt.Milliseconds())
		endTime := formatSRTTime(item.EndAt.Milliseconds())
		sb.WriteString(fmt.Sprintf("%d\n", lineNum))
		sb.WriteString(fmt.Sprintf("%s --> %s\n", startTime, endTime))
		for _, text := range item.Texts {
			s := strings.TrimSpace(text)
			if len(s) == 0 {
				continue
			}
			sb.WriteString(s)
		}
		sb.WriteString("\n\n")
	}
	// create file
	f, err := os.Create(srtPath)
	if err != nil {
		logger.GetLogger(ctx).Warn(err.Error())
		return err
	}
	// remember to close the file
	defer f.Close()

	// 添加 UTF-8 字节顺序标记 (BOM)
	utf8BOM := []byte{0xEF, 0xBB, 0xBF}
	_, err = f.Write(utf8BOM)
	if err != nil {
		logger.GetLogger(ctx).Warn(err.Error())
		return err
	}

	// create new buffer
	buffer := bufio.NewWriter(f)
	_, err = buffer.WriteString(sb.String())
	if err != nil {
		logger.GetLogger(ctx).Warn(err.Error())
		return err
	}
	// flush buffered data to the file
	if err = buffer.Flush(); err != nil {
		logger.GetLogger(ctx).Warn(err.Error())
		return err
	}

	return nil
}

// formatSRTTime 转化为srt专用的时间表示
func formatSRTTime(ms int64) string {
	h := ms / 3600000
	ms -= h * 3600000
	m := ms / 60000
	ms -= m * 60000
	s := ms / 1000
	ms -= s * 1000
	return fmt.Sprintf("%02d:%02d:%02d,%03d", h, m, s, ms)
}
