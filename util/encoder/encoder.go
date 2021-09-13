package encoder

type Encoder interface {
	Encode(str string) string
	Decode(str string) ([]byte, error)
}
