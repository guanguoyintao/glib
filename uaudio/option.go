package uaudio

type audioConf struct {
	length *uint32
}

type ConfigOption func(o *audioConf)

func WithAudioLength(length int) ConfigOption {
	l := uint32(length)
	return func(o *audioConf) {
		o.length = &l
	}
}
