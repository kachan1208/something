package buffer

type Stream interface {
	Read() (interface{}, error)
	Write(interface{}) error
	Get() chan interface{}
}
