package uaudio

import (
	"encoding/binary"
	"errors"
)

type WavHeaderRiffChunk struct {
	ChunkID   [4]byte // 内容为"RIFF"
	ChunkSize [4]byte // wav文件的字节数, 不包含ChunkID和ChunkSize这8个字节）
	Format    [4]byte // 内容为WAVE
}

func (rc *WavHeaderRiffChunk) setChunkID() {
	copy(rc.ChunkID[:], "RIFF")
}

func (rc *WavHeaderRiffChunk) setChunkSize(wavLength uint32) {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, 36+wavLength)
	copy(rc.ChunkSize[:], bs[:4])
}

func (rc *WavHeaderRiffChunk) setFormat() {
	copy(rc.Format[:], "WAVE")
}

type WavHeaderFmtSubChunk struct {
	Subchunk1ID   [4]byte // 内容为"fmt "
	Subchunk1Size [4]byte // Fmt所占字节数，为16
	AudioFormat   [2]byte // 存储音频的编码格式，pcm为1
	NumChannels   [2]byte // 通道数, 单通道为1,双通道为2
	SampleRate    [4]byte // 采样率，如8k, 44.1k等
	ByteRate      [4]byte // 每秒存储的byte数，其值=SampleRate * NumChannels * BitsPerSample/8
	BlockAlign    [2]byte // 块对齐大小，其值=NumChannels * BitsPerSample/8
	BitsPerSample [2]byte // 每个采样点的bit数，一般为8,16,32等。
}

func (fsc *WavHeaderFmtSubChunk) setSubchunk1ID() {
	copy(fsc.Subchunk1ID[:], "fmt ")
}

func (fsc *WavHeaderFmtSubChunk) setSubchunk1Size() {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(16))
	copy(fsc.Subchunk1Size[:], bs[:4])
}

func (fsc *WavHeaderFmtSubChunk) setAudioFormat(audioFormat uint16) {
	bs := make([]byte, 2)
	binary.LittleEndian.PutUint16(bs, audioFormat)
	copy(fsc.AudioFormat[:], bs[:2])
}

func (fsc *WavHeaderFmtSubChunk) setNumChannels(numChannels uint16) {
	bs := make([]byte, 2)
	binary.LittleEndian.PutUint16(bs, numChannels)
	copy(fsc.NumChannels[:], bs[:2])
}

func (fsc *WavHeaderFmtSubChunk) setSampleRate(samplingRate uint32) {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, samplingRate)
	copy(fsc.SampleRate[:], bs[:4])
}

func (fsc *WavHeaderFmtSubChunk) setByteRate(samplingRate uint32, numChannels, bitsPerSample uint16) {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, samplingRate*uint32(numChannels)*uint32(bitsPerSample)/8)
	copy(fsc.ByteRate[:], bs[:4])
}

func (fsc *WavHeaderFmtSubChunk) setBlockAlign(numChannels, bitsPerSample uint16) {
	bs := make([]byte, 2)
	binary.LittleEndian.PutUint16(bs, numChannels*bitsPerSample/8)
	copy(fsc.BlockAlign[:], bs[:2])
}

func (fsc *WavHeaderFmtSubChunk) setBitsPerSample(bitsPerSample uint16) {
	bs := make([]byte, 2)
	binary.LittleEndian.PutUint16(bs, bitsPerSample)
	copy(fsc.BitsPerSample[:], bs[:2])
}

type WavHeaderDataSubChunk struct {
	Subchunk2ID   [4]byte // 内容为"data"
	Subchunk2Size [4]byte // 内容为接下来的正式的数据部分的字节数，其值=NumSamples * NumChannels * BitsPerSample/8
}

func (dsc *WavHeaderDataSubChunk) setSubchunk2ID() {
	copy(dsc.Subchunk2ID[:], "data")
}

func (dsc *WavHeaderDataSubChunk) setSubchunk2Size(wavLength uint32) {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, wavLength)
	copy(dsc.Subchunk2Size[:], bs[:4])
}

type WavHeader struct {
	WavHeaderRiffChunk
	WavHeaderFmtSubChunk
	WavHeaderDataSubChunk
}

func (wh *WavHeader) Marshal() []byte {
	header := make([]byte, 0, 44)
	// riff chunk
	header = append(header, wh.ChunkID[:]...)
	header = append(header, wh.ChunkSize[:]...)
	header = append(header, wh.Format[:]...)
	// format subchunk
	header = append(header, wh.Subchunk1ID[:]...)
	header = append(header, wh.Subchunk1Size[:]...)
	header = append(header, wh.AudioFormat[:]...)
	header = append(header, wh.NumChannels[:]...)
	header = append(header, wh.SampleRate[:]...)
	header = append(header, wh.ByteRate[:]...)
	header = append(header, wh.BlockAlign[:]...)
	header = append(header, wh.BitsPerSample[:]...)
	// data subchunk
	header = append(header, wh.Subchunk2ID[:]...)
	header = append(header, wh.Subchunk2Size[:]...)

	return header
}

func (wh *WavHeader) Unmarshal() string {
	// todo
	return ""
}

func Pcm2Wav(pcm []byte, numChannels uint16, samplingRate uint32, bitsPerSample uint16, opts ...ConfigOption) ([]byte, error) {
	// 加载选项
	audioConf := &audioConf{}
	for _, opt := range opts {
		opt(audioConf)
	}
	if numChannels != 1 && numChannels != 2 {
		return nil, errors.New("invalid_channels_value")
	}
	if samplingRate != 8000 && samplingRate != 16000 {
		return nil, errors.New("invalid_sample_rate_value")
	}
	if bitsPerSample != 8 && bitsPerSample != 16 {
		return nil, errors.New("invalid_bits_per_sample_value")
	}
	var wavLength uint32
	if audioConf.length != nil {
		wavLength = *audioConf.length
	} else {
		wavLength = uint32(len(pcm))
	}
	// create header
	// riff chunk
	rc := WavHeaderRiffChunk{}
	rc.setChunkID()
	rc.setFormat()
	rc.setChunkSize(wavLength)
	// format subchunk
	fsc := WavHeaderFmtSubChunk{}
	fsc.setSubchunk1ID()
	fsc.setSubchunk1Size()
	fsc.setAudioFormat(1)
	fsc.setNumChannels(numChannels)
	fsc.setSampleRate(samplingRate)
	fsc.setByteRate(samplingRate, numChannels, bitsPerSample)
	fsc.setBlockAlign(numChannels, bitsPerSample)
	fsc.setBitsPerSample(bitsPerSample)
	// data subchunk
	dsc := WavHeaderDataSubChunk{}
	dsc.setSubchunk2ID()
	dsc.setSubchunk2Size(wavLength)
	// wav header
	header := &WavHeader{
		WavHeaderRiffChunk:    rc,
		WavHeaderFmtSubChunk:  fsc,
		WavHeaderDataSubChunk: dsc,
	}
	// 拼接wav header和原始pcm
	wav := append(header.Marshal(), pcm...)

	return wav, nil
}
