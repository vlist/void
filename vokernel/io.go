package vokernel

//type BranchWriter struct {
//	PipeWriter interface{Write (p []byte) (n int, e error)}
//	EchoWriter interface{Write (p []byte) (n int, e error)}
//}
//
//func (mw *BranchWriter) Write (p []byte) (n int, e error){
//	if mw.EchoWriter!=nil{
//		go mw.EchoWriter.Write(p)
//	}
//	return mw.PipeWriter.Write(p)
//}

type VolatileWriter struct{
	Destination interface{Write (p []byte) (n int, e error)}
}
func (vw *VolatileWriter) Write (p []byte) (n int, e error){
	return vw.Destination.Write(p)
}
//
//type VolatileReader struct{
//	Source interface{Read (p []byte) (n int, e error)}
//}
//func (vr *VolatileReader) Read (p []byte) (n int, e error){
//	return vr.Source.Read(p)
//}