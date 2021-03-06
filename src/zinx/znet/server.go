/*
   @Author: Jason Chen
   @Date: 2020/11/12 10:18
   @Code: UTF-8
*/
package znet

import (
	"errors"
	"fmt"
	"net"
	"zinx/ziface"
)

//iServer的接口实现，定义一个Server的服务器模块
type  Server struct {
	//服务器的名称
	Name string
	//版本
	IPversion string
	//监听
	IP string
	//端口
	Port int

}

//定义当前客户端链接多绑定的handle api(目前这个handle是写死的，以后优化应该有用户自定义handle)
func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error  {
	//会先业务
	fmt.Println("[Conn Handle] callbackToClient..")
	if _,err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("write back buf err",err)
		return errors.New("callBackToClient error")
	}
	return nil
}

//启动服务器
func (s *Server)Start()  {
	fmt.Printf("[start] Server Listenner at IP: %s, Port: %d, is starting\n", s.IP, s.Port)

	go func() {
		//1 获取一个TCP的addr
		addr, err := net.ResolveTCPAddr(s.IPversion, fmt.Sprintf("%s:%d", s.IP,s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error: ",err)
			return
		}
		//2 监听服务器地址
		listenner, err := net.ListenTCP(s.IPversion, addr)
		if err != nil {
			fmt.Println("listen",s.IPversion, "err", err)
			return
		}

		fmt.Println("start Zinx server succ", s.Name, "succ listenning...")
		var cid uint32
		cid = 0

		//3 阻塞的等待客户端连接， 处理客户端连接业务
		for  {
			//如果有客户端链接过来，阻塞会返回
			conn, err := listenner.AcceptTCP()
			if err != nil{
				fmt.Println("Accept err", err)
				continue
			}

			//go func() {
			//	for  {
			//		buf := make([]byte, 512)
			//		cnt, err := conn.Read(buf)
			//		if err != nil {
			//			fmt.Println("recv buf err",err)
			//			continue
			//		}
			//
			//		fmt.Printf("recv client buf %s cnt: %d", buf, cnt)
			//
			//		//回显功能
			//		if _, err := conn.Write(buf[:cnt]); err != nil {
			//			fmt.Println("buf Write err", err)
			//			continue
			//		}
			//	}
			//}()


			//将处理新链接的业务方法和conn进行绑定， 得到我们的链接模块
			dealConn := NewConnection(conn, cid, CallBackToClient)
			cid++

			//启动当前的链接业务处理
			go dealConn.Start()
		}
	}()



}
//停止服务器
func (s *Server)Stop()  {
	//TODO 将一些服务器的资源，状态或者一些开启的服务停止
}
//运行服务器
func (s *Server)Serve()  {
	//启动server的服务功能
	s.Start()

	//TODO 做一些启动服务器之后的额外业务

	//阻塞状态
	select {

	}
}

//初始化服务器的方法
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name : name,
		IPversion: "tcp4",
		IP: "0.0.0.0",
		Port: 8999,

	}
	return s
}