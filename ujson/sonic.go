package ujson

import (
	"github.com/bytedance/sonic"
	"io"
)

type sonicJsonAPI struct {
	sonic.Config
}

func newSonicJsonAPI(config sonic.Config) uJsonAPI {
	return &sonicJsonAPI{config}
}

func (s *sonicJsonAPI) MarshalToString(v interface{}) (string, error) {
	api := s.Config.Froze()

	return api.MarshalToString(v)
}

func (s *sonicJsonAPI) Marshal(v interface{}) ([]byte, error) {
	api := s.Config.Froze()

	return api.Marshal(v)
}

func (s *sonicJsonAPI) MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	api := s.Config.Froze()

	return api.MarshalIndent(v, prefix, indent)
}

func (s *sonicJsonAPI) UnmarshalFromString(str string, v interface{}) error {
	api := s.Config.Froze()

	return api.UnmarshalFromString(str, v)
}

func (s *sonicJsonAPI) Unmarshal(data []byte, v interface{}) error {
	api := s.Config.Froze()

	return api.Unmarshal(data, v)
}

func (s *sonicJsonAPI) NewEncoder(writer io.Writer) Encoder {
	api := s.Config.Froze()

	return api.NewEncoder(writer)
}

func (s *sonicJsonAPI) NewDecoder(reader io.Reader) Decoder {
	api := s.Config.Froze()

	return api.NewDecoder(reader)
}

func (s *sonicJsonAPI) Valid(data []byte) bool {
	api := s.Config.Froze()

	return api.Valid(data)
}
