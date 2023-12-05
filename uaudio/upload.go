package uaudio

import (
	"context"
	"fmt"
	"git.umu.work/AI/uglib/ucontext"
	ubytebufferpool "git.umu.work/AI/uglib/udatastructure/byte_buffer_pool"
	ulocal "git.umu.work/AI/uglib/ustorage/local"
	uoss "git.umu.work/AI/uglib/ustorage/oss"
	"git.umu.work/be/goframework/logger"
	"os"
	"path"
	"sync"
	"time"
)

type AudioEncodingType int32

const (
	PCM AudioEncodingType = 0
	WAV AudioEncodingType = 1
)

type AudioInfo struct {
	Name         string
	EncodingType AudioEncodingType
	URI          string
}

type UploadAudio struct {
	fd            *os.File
	tmpDir        string
	audioFilePath string
	oss           uoss.ObjectStorage
	info          *AudioInfo
	infoChan      chan *AudioInfo
	rw            *sync.RWMutex
	wg            *sync.WaitGroup
	bufferPool    chan []byte
	ticker        *time.Ticker
	stopWrite     bool
}

func NewUploadAudio(ctx context.Context, oss uoss.ObjectStorage, info *AudioInfo) (upload *UploadAudio, err error) {
	switch info.EncodingType {
	case PCM:
		info.Name = info.Name + ".pcm"
	case WAV:
		info.Name = info.Name + ".wav"
	}
	tmpDir := os.TempDir()
	tmp := path.Join(tmpDir, info.Name)
	fd, err := ulocal.GetFd(tmp, os.O_WRONLY|os.O_APPEND)
	if err != nil {
		logger.GetLogger(ctx).Error(err.Error())
		return nil, err
	}
	ticker := time.NewTicker(2 * time.Second)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	upload = &UploadAudio{
		fd:            fd,
		wg:            wg,
		tmpDir:        tmpDir,
		audioFilePath: tmp,
		oss:           oss,
		rw:            &sync.RWMutex{},
		info:          info,
		infoChan:      make(chan *AudioInfo),
		bufferPool:    make(chan []byte, 100),
		ticker:        ticker,
	}
	innerCtx := ucontext.NewUValueContext(ctx)
	go upload.flushLoop(innerCtx)

	return upload, nil
}

func (u *UploadAudio) flushLoop(ctx context.Context) {
	defer u.wg.Done()
	pool := ubytebufferpool.NewByteBufferPool(ctx)
	defer ubytebufferpool.ReleaseByteBufferPool(pool)
	var writing bool
	for {
		select {
		case buf, ok := <-u.bufferPool:
			if !ok {
				return
			}
			_, err := pool.Write(buf)
			if err != nil {
				logger.GetLogger(ctx).Warn(err.Error())
				return
			}
		case _, ok := <-u.ticker.C:
			if !ok {
				return
			}
			if writing {
				continue
			}
			data := pool.ByteBuffer
			pool.Reset()
			if len(data) > 0 {
				writing = true
				go func() {
					err := u.write(ctx, data)
					if err != nil {
						logger.GetLogger(ctx).Warn(err.Error())
						return
					}
					writing = false
				}()
			} else {
				if u.stopWrite {
					u.ticker.Stop()
					close(u.bufferPool)
				}
			}
		}
	}
}

func (u *UploadAudio) write(ctx context.Context, content []byte) error {
	_, err := u.fd.Stat()
	if err != nil {
		logger.GetLogger(ctx).Error(err.Error())
		return err
	}
	_, err = u.fd.Write(content)
	if err != nil {
		logger.GetLogger(ctx).Error(err.Error())
		return err
	}

	return nil
}

func (u *UploadAudio) Write(ctx context.Context, content []byte) (err error) {
	u.bufferPool <- content

	return nil
}

func (u *UploadAudio) GetAudioInfo() *AudioInfo {
	u.rw.RLock()
	info := u.info
	u.rw.RUnlock()

	return info
}

func (u *UploadAudio) close(ctx context.Context) error {
	err := u.deleteFileFragments(ctx)
	if err != nil {
		logger.GetLogger(ctx).Error(err.Error())
		return err
	}
	err = u.fd.Close()
	if err != nil {
		logger.GetLogger(ctx).Error(err.Error())
		return err
	}

	return nil
}

func (u *UploadAudio) deleteFileFragments(ctx context.Context) error {
	// 删除临时文件
	err := os.Remove(u.audioFilePath)
	if err != nil {
		logger.GetLogger(ctx).Error(err.Error())
		return err
	}
	logger.GetLogger(ctx).Info(fmt.Sprintf("%s local file is deleted", u.audioFilePath))

	return nil
}

// Upload  生成唯一文件标识并上传文件，返回文件访问地址
func (u *UploadAudio) Upload(ctx context.Context, ossPath string) (audioUri string, err error) {
	defer func() {
		if err != nil {
			e := u.close(ctx)
			if e != nil {
				logger.GetLogger(ctx).Error(e.Error())
			}
		}
	}()
	u.stopWrite = true
	u.wg.Wait()
	logger.GetLogger(ctx).Info("start upload file to oss")
	tmpFileInfo, err := u.fd.Stat()
	if err != nil {
		logger.GetLogger(ctx).Warn(err.Error())
		return "", err
	}
	u.rw.RLock()
	info := u.info
	u.rw.RUnlock()
	switch info.EncodingType {
	case PCM:
		fileNameWithoutExtension := ulocal.RemoveFileExtension(tmpFileInfo.Name())
		pcmFilePath := path.Join(u.tmpDir, tmpFileInfo.Name())
		// pcm加wav头
		wavFileName := fileNameWithoutExtension + ".wav"
		wavFilePath := path.Join(u.tmpDir, wavFileName)
		wavHeader, err := Pcm2Wav([]byte{}, 1, 16000, 16,
			WithAudioLength(int(tmpFileInfo.Size())))
		if err != nil {
			logger.GetLogger(ctx).Error(err.Error())
			return "", err
		}
		// wav 文件写入
		err = ulocal.WriteHead(wavHeader, pcmFilePath, wavFilePath)
		if err != nil {
			logger.GetLogger(ctx).Error(err.Error())
			return "", err
		}
		err = os.Remove(pcmFilePath)
		if err != nil {
			logger.GetLogger(ctx).Error(err.Error())
			return "", err
		}
		u.audioFilePath = wavFilePath
	}
	logger.GetLogger(ctx).Info(fmt.Sprintf("%s audio file  start upload  to %s", u.audioFilePath, ossPath))
	audioUri, err = u.oss.UploadByFile(ctx, u.audioFilePath, ossPath)
	if err != nil {
		logger.GetLogger(ctx).Error(err.Error())
		return "", err
	}
	u.rw.Lock()
	u.info.URI = audioUri
	u.rw.Unlock()

	return audioUri, nil
}
