package uaudio

import (
	"fmt"
	"io/ioutil"
	"path"
	"testing"
)

func TestPcm2Wav(t *testing.T) {
	type args struct {
		fileName      string
		wavFileName   string
		numChannels   uint16
		samplingRate  uint32
		bitsPerSample uint16
		opts          []ConfigOption
	}
	tests := []struct {
		name       string
		args       args
		wantResult []byte
	}{
		{
			name: "16k16bit",
			args: args{
				fileName:      "16k16bit.pcm",
				wavFileName:   "16k16bit.wav",
				numChannels:   1,
				samplingRate:  16000,
				bitsPerSample: 16,
			},
			wantResult: nil,
		},
		{
			name: "16k16bit 指定pcm length10",
			args: args{
				fileName:      "16k16bit.pcm",
				wavFileName:   "16k16bit-100g.wav",
				numChannels:   1,
				samplingRate:  16000,
				bitsPerSample: 16,
				opts: []ConfigOption{
					WithAudioLength(100 * (1 << (3 * 10))),
				},
			},
			wantResult: nil,
		},
		{
			name: "16k16bit 指定pcm length10",
			args: args{
				fileName:      "16k16bit.pcm",
				wavFileName:   "16k16bit-100m.wav",
				numChannels:   1,
				samplingRate:  16000,
				bitsPerSample: 16,
				opts: []ConfigOption{
					WithAudioLength(100 * (1 << (2 * 10))),
				},
			},
			wantResult: nil,
		},
		{
			name: "16k16bit 指定pcm length10",
			args: args{
				fileName:      "16k16bit.pcm",
				wavFileName:   "16k16bit-1g.wav",
				numChannels:   1,
				samplingRate:  16000,
				bitsPerSample: 16,
				opts: []ConfigOption{
					WithAudioLength(1 * (1 << (3 * 10))),
				},
			},
			wantResult: nil,
		},
		{
			name: "16k16bit 指定pcm length10000",
			args: args{
				fileName:      "16k16bit.pcm",
				wavFileName:   "16k16bit-100k.wav",
				numChannels:   1,
				samplingRate:  16000,
				bitsPerSample: 16,
				opts: []ConfigOption{
					WithAudioLength(100 * (1 << (1 * 10))),
				},
			},
			wantResult: nil,
		},
	}
	audioFilePrefixPath := "./test"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			audioFilePath := path.Join(audioFilePrefixPath, tt.args.fileName)
			audioBytes, err := ioutil.ReadFile(audioFilePath)
			if err != nil {
				t.Error(err)
			}
			gotResult, err := Pcm2Wav(audioBytes, tt.args.numChannels, tt.args.samplingRate,
				tt.args.bitsPerSample, tt.args.opts...)
			if err != nil {
				t.Error(err)
			}
			saveAudioFile := path.Join("test", "tmp_"+tt.args.wavFileName+".wav")
			fmt.Println(saveAudioFile)
			err = ioutil.WriteFile(saveAudioFile, gotResult, 0644)
			if err != nil {
				t.Error(err)
			}
		})
	}
}
