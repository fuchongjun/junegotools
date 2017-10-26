package study

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

)

func spinn(delay time.Duration) {
	for {
		for _, r := range `-/|\` {
			fmt.Printf("\r%c", r)
			time.Sleep(delay)
		}
	}
}
func fib(x int) int {
	if x < 2 {
		return x
	}
	return fib(x-1) + fib(x-2)
}
func GoroutineTest1() {
	go spinn(100 * time.Millisecond)
	const n = 40
	fibN := fib(n)
	fmt.Printf("\rFibonacci(%d)=%d", n, fibN)
}

/*并发小例子
func TimeSendServer(port int)  {
	listenner,err:=net.Listen("tcp","localhost:"+strconv.Itoa(port))

	if err!=nil {
		log.Fatal(err)
	}
	for {
		conn,nil:=listenner.Accept()
		if err!=nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}
func handleConn(c net.Conn)  {
	defer c.Close()
	for{
		_,err:=io.WriteString(c,time.Now().Format("15:04:05\n"))
		if err!=nil {
			return
		}
		time.Sleep(1*time.Second)
	}
}

func TimeClient()  {
	conn,err:=net.Dial("tcp","localhost:8000")
	if err!=nil {
		log.Fatal(err)
	}
	defer conn.Close()
	mustCopy(os.Stdout,conn)
}
func mustCopy(dst io.Writer,src io.Reader)  {
	if _,err:=io.Copy(dst,src);err!=nil {
		log.Fatal(err)
	}
}
*/
//===========回声服务器==========
func SoundServer(port int) {
	listenner, err := net.Listen("tcp", "localhost:"+strconv.Itoa(port))

	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, nil := listenner.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		go func(c net.Conn) {
			input := bufio.NewScanner(c)
			for input.Scan() {
				go echo(c, input.Text(), 1*time.Second)
			}
			c.Close()
		}(conn)
	}
}
func echo(c net.Conn, shout string, delay time.Duration) {
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToLower(shout))
}

func SoundClient() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	go mustCopy(os.Stdout, conn)
	//s := bufio.NewScanner(os.Stdin)//从输入读取一行，换行为结束
	//for s.Scan() {
	//	fmt.Fprintf(conn, "%s\n", s.Text())
	//}
	in:=bufio.NewReader(os.Stdin)
	for  {
		r,_,err:= in.ReadLine()
		if err==io.EOF {
			break
		}
		if err!=nil {
		}
		mustCopy(conn,strings.NewReader(string(r)+"\n"))
	}


}
func SoundClient2() {
	pRemoteTCPAddr, err := net.ResolveTCPAddr("tcp", "localhost:8000")
	conn, err := net.DialTCP("tcp",nil, pRemoteTCPAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	done:=make(chan struct{})
	go func() {
		io.Copy(os.Stdout, conn)//阻塞，读conn中的内容，停止写之后才能读
		log.Println("done")
		done<- struct {}{} //发送通道信号
	}()
	mustCopy(conn,strings.NewReader("Hello\n"))//向conn中写内容
	//mustCopy(conn,os.Stdin)
	conn.CloseWrite()
	<-done//收到通道信号
}
func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
