package vokernel

import "io"

type multiW struct {
	io.Writer
	cs []io.Closer
}
func MultiWriteCloser(ws ...io.Writer) io.WriteCloser {
	m := &multiW{Writer: io.MultiWriter(ws...)}
	for _, w := range ws {
		if c, ok := w.(io.Closer); ok {
			m.cs = append(m.cs, c)
		}
	}
	return m
}
func (m *multiW) Close() error {
	var first error
	for _, c := range m.cs {
		if err := c.Close(); err != nil && first == nil {
			first = err
		}
	}
	return first
}

type VolatileWriter struct{
	Destination Writer
}
func (vw *VolatileWriter) Write (p []byte) (n int, e error){
	return vw.Destination.Write(p)
}
func (vw *VolatileWriter) Close () error{
	//return vw.Destination.Close()
	return nil
}
type Writer interface{Write (p []byte) (n int, e error)}
