//+build unit

package decoder

import (
	"strings"
	"testing"

	"github.com/kachan1208/something/processor/model"
)

//Func init that will init everything

//Few simple mocks

//Test decode string stream success(check err and num of elems)
//Test decode broken string stream unsuccess(check err)
//Test decode file stream success
//Test decode wrong type stream
//Test when output stream buffer passed

var (
	BenchData = `
	{
		"AEAJM": {
		"name": "Ajman",
		"city": "Ajman",
		"country": "United Arab Emirates",
		"alias": [],
		"regions": [],
		"coordinates": [
			55.5136433,
			25.4052165
		],
		"province": "Ajman",
		"timezone": "Asia/Dubai",
		"unlocs": [
			"AEAJM"
		],
		"code": "52000"
		}
	}
	`
)

type DiscardStream struct{}

func (d *DiscardStream) Read() (interface{}, error)   { return nil, nil }
func (d *DiscardStream) Write(data interface{}) error { return nil }
func (d *DiscardStream) Get() chan interface{}        { return nil }

func BenchmarkJSONStreamDecoderOnOneElement(b *testing.B) {
	dec := NewJSONStreamDecoder()
	reader := strings.NewReader(BenchData)
	var port model.Port
	discardStream := DiscardStream{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		dec.Decode(reader, &port, &discardStream)
	}
}
