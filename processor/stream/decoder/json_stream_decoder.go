package decoder

import (
	"encoding/json"
	"io"

	"github.com/kachan1208/something/processor/stream/buffer"
)

type JSONStreamDecoder struct{}

func NewJSONStreamDecoder() *JSONStreamDecoder {
	return &JSONStreamDecoder{}
}

func (j *JSONStreamDecoder) Decode(
	input io.Reader,
	obj interface{},
	output buffer.Stream,
) (n uint64, err error) {
	dec := json.NewDecoder(input)
	_, err = dec.Token()
	if err != nil {
		return
	}

	for dec.More() {
		_, err = dec.Token()
		if err != nil {
			return
		}

		err = dec.Decode(&obj)
		if err != nil {
			return
		}

		err = output.Write(obj)
		if err != nil {
			return
		}

		n++
	}

	return
}
