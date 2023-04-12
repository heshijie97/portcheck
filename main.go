package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"portcheck/flaginit"
	"strings"
	"sync"
	"time"
)

// 声明全局等待组变量
var wg sync.WaitGroup

func PortCheck(host, port string, writer *bufio.Writer) {
	defer wg.Done()
	var target string
	if strings.Contains(host, ".") {
		target = fmt.Sprintf("%s:%s", host, port)
	} else if strings.Contains(host, ":") {
		target = fmt.Sprintf("[%s]:%s", host, port)
	}
	conn, err := net.DialTimeout("tcp", target, time.Second*3)
	if err != nil {
		//端口未连接，不需要关闭，否则会报错
		fmt.Printf("%s %s down\n", host, port)
		return
	} else {
		fmt.Printf("%s %s up\n", host, port)
	}
	conn.Close()
}

// 传入参数 -f 指定待ping的ip文件 -d 指定输出结果

func main() {
	//初始化flag
	checkfile, resutlfile := flaginit.InitFlag()
	//打开文件
	cf, err := os.Open(checkfile)
	if err != nil {
		log.Fatal(err)
	}
	defer cf.Close()
	rf, err := os.Create(resutlfile)
	if err != nil {
		log.Fatal(err)
	}
	defer rf.Close()
	//读入带缓冲的io
	reader := bufio.NewReader(cf)
	//写入带缓冲的io
	writer := bufio.NewWriter(rf)
	defer writer.Flush()
	for {
		hostPort, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		//去除换行符
		hostPort = strings.Trim(hostPort, "\n")
		//以空格分割字符串
		arr := strings.Fields(hostPort)
		host := arr[0]
		port := arr[1]
		wg.Add(1) // 登记1个goroutine
		go PortCheck(host, port, writer)
	}
	wg.Wait() // 阻塞等待登记的goroutine完成
}
