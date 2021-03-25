package vokernel

import "io"

type multi struct {
	io.Writer
	cs []io.Closer
}

func MultiWriteCloser(ws ...io.Writer) io.WriteCloser {
	m := &multi{Writer: io.MultiWriter(ws...)}
	for _, w := range ws {
		if c, ok := w.(io.Closer); ok {
			m.cs = append(m.cs, c)
		}
	}
	return m
}

func (m *multi) Close() error {
	var first error
	for _, c := range m.cs {
		if err := c.Close(); err != nil && first == nil {
			first = err
		}
	}
	return first
}


func BiWriteCloser(primary io.Writer,second io.Writer)io.WriteCloser{
	return &ForkWriteCloser{
		PrimaryWriter: primary,
		SecondWriter:  second,
	}
}
type ForkWriteCloser struct {
	PrimaryWriter interface{Write (p []byte) (n int, e error)}
	SecondWriter interface{Write (p []byte) (n int, e error)}

}
func (f *ForkWriteCloser) Write (p []byte) (n int, e error){
	go func (){
		if f.SecondWriter!=nil {
			f.SecondWriter.Write(p)
		}
	}()
	return f.PrimaryWriter.Write(p)
}
func (f *ForkWriteCloser) Close () error{
	if c,ok:=f.SecondWriter.(io.Closer);ok{
		go c.Close()
	}
	if c,ok:=f.PrimaryWriter.(io.Closer);ok{
		return c.Close()
	}
	return nil
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
//
//type VolatileReader struct{
//	Source interface{Read (p []byte) (n int, e error)}
//}
//func (vr *VolatileReader) Read (p []byte) (n int, e error){
//	return vr.Source.Read(p)
//}
type Writer interface{Write (p []byte) (n int, e error)}
