package uffmpeg

import (
	"bytes"
	"context"
	"fmt"
	"git.umu.work/AI/uglib/uhttp"
	"git.umu.work/AI/uglib/ujson"
	"git.umu.work/be/goframework/logger"
	"github.com/pkg/errors"
	"gopkg.in/vansante/go-ffprobe.v2"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var ffprobeBinPath = "ffprobe"

type StreamingMediaType uint8

const (
	AudioStreamingMediaType StreamingMediaType = 1
	VideoStreamingMediaType StreamingMediaType = 2
)

type CodecNameType string

const (
	AudioCodecPCMS16LE CodecNameType = "pcm_s16le"
	AudioCodecAAC      CodecNameType = "aac"
)

type StreamingMediaInfo struct {
	Name       string
	CodecName  CodecNameType
	Duration   time.Duration
	SampleRate int64
	Size       int64
}

func probeURL(ctx context.Context, fileURL string) (data *ffprobe.ProbeData, err error) {
	args := append([]string{
		"-loglevel", "fatal",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
	})
	// Add the file argument
	args = append(args, fileURL)
	cmd := exec.CommandContext(ctx, ffprobeBinPath, args...)
	cmd.SysProcAttr = nil
	// run probe
	var outputBuf bytes.Buffer
	var stdErr bytes.Buffer
	cmd.Stdout = &outputBuf
	cmd.Stderr = &stdErr
	err = cmd.Run()
	if err != nil {
		return nil, err
	}
	if stdErr.Len() > 0 {
		if strings.Contains(stdErr.String(), "Warning:") {
			logger.GetLogger(ctx).Warn(fmt.Sprintf("ffprobe std warning: %s", stdErr.String()))
		} else {
			return nil, errors.New(stdErr.String())
		}
	}
	data = &ffprobe.ProbeData{}
	err = ujson.Unmarshal(outputBuf.Bytes(), data)
	if err != nil {
		return data, fmt.Errorf("error parsing ffprobe output: %w", err)
	}
	if data.Format == nil {
		return data, fmt.Errorf("no format data found in ffprobe output")
	}

	return data, nil
}

func GetStreamingMediaInfo(ctx context.Context, src string, mediaType StreamingMediaType) (*StreamingMediaInfo, error) {
	logger.GetLogger(ctx).Info(fmt.Sprintf("ffprobe uri is %s", src))
	now := time.Now()
	fileName, _, err := uhttp.ExtractFileNameFromURL(src)
	if err != nil {
		logger.GetLogger(ctx).Error(err.Error())
		return nil, err
	}
	data, err := probeURL(ctx, src)
	if err != nil {
		return nil, errors.Errorf("Error getting ffprobe data: %v, uri=%s", err, src)
	}
	duration := data.Format.Duration()
	size, err := strconv.ParseInt(data.Format.Size, 10, 64)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("size not right: %v", data.Format.Size))
	}
	info := &StreamingMediaInfo{
		Name:     fileName,
		Duration: duration,
		Size:     size,
	}
	switch mediaType {
	case AudioStreamingMediaType:
		logger.GetLogger(ctx).Info(fmt.Sprintf("ffprobe deal with video uri %s", src))
		firstStream := data.FirstAudioStream()
		if firstStream == nil {
			return info, nil
		}
		sampleRate, err := strconv.ParseInt(firstStream.SampleRate, 10, 64)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("sample_rate not right: %v", firstStream.SampleRate))
		}
		info.SampleRate = sampleRate
		info.CodecName = CodecNameType(firstStream.CodecName)
	case VideoStreamingMediaType:
		logger.GetLogger(ctx).Info(fmt.Sprintf("ffprobe deal with aduio uri %s", src))
		firstStream := data.FirstVideoStream()
		if firstStream == nil {
			return info, nil
		}
		info.CodecName = CodecNameType(firstStream.CodecName)
	}
	logger.GetLogger(ctx).Info(fmt.Sprintf("remote streaming media info took %+v", time.Since(now)))

	return info, nil
}

func StreamingMediaConvertToWAV(ctx context.Context, src, dst string) error {
	args := []string{
		"-y",
		"-i",
		src,
		"-ac",
		"1",
		"-ar",
		"16000",
		dst,
	}
	cmd := exec.CommandContext(ctx, "ffmpeg", args...)
	_, err := cmd.Output()
	if err != nil {
		logger.GetLogger(ctx).Warn(err.Error())
		return err
	}
	return nil
}
