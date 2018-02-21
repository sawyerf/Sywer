//Version 1.5
package main

import ("fmt"
		"net"
		"log"
		"os"
		"io/ioutil")

const src_path string = "/srv/http/" //the directory where you have the files you want to share 

type Request struct{
	//info request
	ip string
	method string
	host string
	user_agent string
	connection string
	//File To send
	path string
	size int64
	type_path bool //File = True and Directory = False
	err string
}

func (c Request) Header() []byte{ 
	var header string = ""
	if c.method == "GET"{
		header += "HTTP/1.1 "
		switch c.err{
		case "200":
			header += "200 OK\r\n"
		case "404":
			header += "404 Not Found"
		case "301":
			header += "301 Moved Permanently\r\n"
			header += "Location: http://" + c.host + "/" + c.path + "/\r\n"
		}
		if 0 < c.size{
		header += "Accept-Ranges: bytes\r\nContent-Lenght: " + fmt.Sprint(c.size) + "\r\nConnection: close\r\n\r\n"
		} else {
			header += "Accept-Ranges: bytes\r\nonnection: close\r\n\r\n"
		}
	}	
	return []byte(header)
}

func (c Request) Data(conn net.Conn){
	if c.method == "GET"{
		switch c.err{
		case "200":
			if !c.type_path{
				_, _ = conn.Write([]byte("<h1>Index Of " + c.path + "</h1>\n<ul>"))
				files, _ := ioutil.ReadDir(src_path + c.path)
				for _, file := range files{
					_, _ = conn.Write([]byte("<li><a href=\"" + file.Name() + "\">" + file.Name() + "</a></li>\n"))
				}
				_, _ = conn.Write([]byte("</ul>"))
				return
			}
			buffer := make([]byte, 1024)
			file, _ := os.Open(src_path + c.path)
			for{
				size, _ := file.Read(buffer)
				if size == 0{
					return
				} else{
					_, _ = conn.Write(buffer[:size])}
			}
		case "404":
			if c.path == "ip"{
				conn.Write([]byte("<h1>" + c.ip + "</h1>"))
				return
			}
			conn.Write([]byte("<h1>404 Not Found</h1>"))
		case "301":
			conn.Write([]byte("<h1>301 Moved Permanently</h1>"))
		default:
			return
		}
	}
}

func request_analyzer(get string, size int) Request{
	var req Request
	nb := 0
	for i := 3; i < size; i++ {
		if get[i] == 13 && get[i+1] == 10{
			if get[nb:nb+3] == "GET"{
				req.method = "GET"
				req.path = get[nb+5:i-9]
				req.type_path = true
				if req.path == ""{
					req.path = "index.html"
				} else{ if req.path[len(req.path)-1] == 47{
					req.type_path = false}}
			} else {if get[nb:nb+5] == "Host:"{
				req.host = get[nb+6:i]
			} else {if get[nb:nb+11] == "User-Agent:"{
				req.user_agent = get[nb+12:i]}}}
			nb = i + 2
		}
	}
	return req
}

func File_check(name string) (string, int64){
	fi, err := os.Stat(src_path + name)
	if os.IsNotExist(err){
		return "404", 30}
	switch mode := fi.Mode(); {
		case mode.IsDir():
			return "301", 30
		case mode.IsRegular():
			return "200", fi.Size()
		default:
			return "404", 30
	}
}

func Dir_check(name string) string{
	fi, err := os.Stat(src_path + name)
	mode := fi.Mode()
	if os.IsNotExist(err) || mode.IsRegular(){
		return "404"}
	return "200"
}

func Ip(ip net.Addr) string{
	ip_str := fmt.Sprint(ip)
	if ip_str[0] == 91{
		return "127.0.0.1"
	}
	for i := 0; i != len(ip_str); i++{
		if ip_str[i] == 58{
			return ip_str[:i]
		}
	}
	return ip_str
}

func recv(conn net.Conn){
	buffer := make([]byte, 4096)
	size, err := conn.Read(buffer)

	if err != nil {
		fmt.Println(err)
		conn.Close()
		return
	}

	request := request_analyzer(string(buffer[:size+11]), size)
	if request.method == ""{
		conn.Close()
		return}

	request.ip = Ip(conn.RemoteAddr())
	// File and Directory check
	if request.type_path{
		request.err, request.size = File_check(request.path)
	} else{
		request.err = Dir_check(request.path)}

	_, err = conn.Write(request.Header())
	request.Data(conn)
	conn.Close()
	fmt.Println(request.ip, "\t", request.method, "/" + request.path + "\t", request.user_agent, "\t", request.err)
}

func main(){
	server, err := net.Listen("tcp", ":8080")
	if err != nil {
			log.Fatalln(err)
	}
	fmt.Println("[*]Start")
	for{
		conn, err := server.Accept()
		if err != nil {
			fmt.Println("[*]Erreur", err)
		}
		//fmt.Println("[*]Nouvelle ecoute", conn)
		go recv(conn)
	}
}

