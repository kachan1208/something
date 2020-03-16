package buffer

type InterfaceStream struct {
	buf chan interface{}
}

func NewInterfaceStream(size int) *InterfaceStream {
	return &InterfaceStream{
		buf: make(chan interface{}, size),
	}
}

func (p *InterfaceStream) Read() (interface{}, error) {
	return <-p.buf, nil
}

func (p *InterfaceStream) Get() chan interface{} {
	return p.buf
}

func (p *InterfaceStream) Write(obj interface{}) error {
	p.buf <- obj
	return nil
}
