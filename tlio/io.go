package tlio


import (
	"errors"
	"sync"
)

type StringWriter interface{
	WriterString(s string)(n int,err error)
}

var ErrShortBuffer = errors.New("short buffer")

var ErrUnexpectedEOF = errors.New("unexpected EOF")


/*
如果w实现了StringWriter接口,调用StringWriter实现的WriterString方法，否则调用Write([]byte(s))
*/
func WriteString(w Writer,s string)(n int,err error){
	if ww,ok := w.(StringWriter);ok {
		return w.WriterString(s)
	}
	return w.Write([]byte(s))
}

/*
Buf长度小于min,报ErrShortBuffer错误
使用r.Read循环读取数据至buf,如果已读取字节数大于等于min或者err不为nil则退出循环。
退出循环后查看已读取字节数大于等于min则设置err为nil然后返回，如果已读取字节数大于0小于min且err==EOF，则err设置为ErrUnexpectedEOF然后返回。
否则直接返回已读取字节数和err
*/
func ReadAtLeast(r Reader,buf []byte,min int)(n int,err error){
	if min <= 0{
		return 0,nil
	}
	if len(fub) < min{
		return 0,ErrShortBuffer
	}
	for n < min && err == nil{		
		var nn int
		nn,err = r.Read(buf)
		n+= nn
	}
	if n < min && err == EOF{
		return n,ErrUnexpectedEOF
	}else if n >= min{
		err = nil
	}
	return 	
}

/*
使用r.Read读取数据到buf,直到buf满了或者中途遇到err不为nil。如果buf未读满且err为EOF，则设置err为ErrUnexpectedEOF
*/
func ReadFull(r Reader,buf []byte)(n int,err error){
	if len(buf) <= 0{
		return 0,ErrShortBuffer
	}
	var nn int
	for n < len(buf) && err == nil{
		nn,err = r.Read(buf[n:])
		n+=nn
	}
	if n <len(buf) && err == EOF{
		return n,ErrUnexpectedEOF
	}	
	return 
}


var ErrBufferLenght = errors.New("n must big than 0")

/*
使用src.Read读取长度为n的内容，并把读取到得内容使用dst.Write进行写入,返回实际写入的数据，如果err==nil的话
*/
func CopyN(dst Writer, src Reader, n int64) (written int64, err error){
	if n <= 0{
		return 0,nil
	}
	buf := make([]byte,n)
	written,err = ReadFull(src,buf)
	if  err != nil{
		return 
	}
	written,err = w.Writer(buf)
	return
}

/*
和CopyN类似，只不过把src内容全部写到dst
*/
func Copy(dst Writer, src Reader) (written int64, err error){
	buf := make([]byte,1024)//默认开1k空间
	buffer := make([]byte,1)	
	for err == nil{
		written,err = src.Read(buf])
		buffer = append(buffer,buf)
	}
	if err != EOF{
		return
	}
	return dst.Writer(buffer)
}

/*
Src的内容读到buf中，buf的内容写到dst中，返回实际写入的字节数
*/
func CopyBuffer(dst Writer, src Reader, buf []byte) (written int64, err error){
	return CopyN(dst,src,len(buf))
}

/*
上面Copy和CopyBuffer的实际实现，如果src实现了WriterTo接口,调用src的WriteTo来进行写入。
*/
func copyBuffer(dst Writer, src Reader, buf []byte) (written int64, err error){
	return CopyN(dst,src,len(buf))
}

type LimitedReader struct {
	R Reader // underlying reader
	N int64  // max bytes remaining
}

func LimitReader(r Reader, n int64) Reader { return &LimitedReader{r, n} }

/*
结构的方法,读取，最多只能读取结构中的N。如果N小于等于0返回EOF。可重复调用，每次调用N减去实际读取的长度。
*/
func (l *LimitedReader) Read(p []byte) (n int, err error){
	if l.N <= 0{
		return 0,EOF
	}
	n,err = l.R.Read(p)
	if n > l.N {
		n = l.N
		l.N = 0
	}else{
		l.N-=n
	}	
	return
}

type SectionReader struct {
	r     ReaderAt
	base  int64
	off   int64
	limit int64
}

func NewSectionReader(r ReaderAt, off int64, n int64) *SectionReader{
	return &SectionReader{r,off,off,off + n}
}

/*
结构体方法,读取到p，从偏移量off开始读取，到limit结束。意思就是最多度limit-off。
*/
func (s *SectionReader) Read(p []byte) (n int, err error){
	
}


