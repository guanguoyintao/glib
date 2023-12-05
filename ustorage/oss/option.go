package uoss

type Options struct {
	ContentType string
}

type Option func(*Options)

func WithContentType(contentType string) Option {
	return func(o *Options) {
		o.ContentType = contentType
	}
}
