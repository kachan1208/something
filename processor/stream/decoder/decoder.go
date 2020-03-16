package decoder

import (
	"io"

	"github.com/kachan1208/something/processor/stream/buffer"
)

type StreamDecoder interface {
	Decode(io.Reader, interface{}, buffer.Stream) (uint64, error)
}
