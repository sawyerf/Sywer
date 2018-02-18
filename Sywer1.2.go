package main

import ("fmt"
		"net"
		"log"
		"os")

type Request struct{
	ip net.Addr
	method string
	path string
	host string
	user_agent string
	connection string
	err string
}
//"HTTP/1.1 200 OK\r\nContent-length: " + fmt.Sprint(len(data)) + "\r\nConnection: Close\r\n\r\n"

func (c Request) Header(size string) []byte{
	var header string = ""
	if c.method == "GET"{
		if c.err == ""{
			header += "HTTP/1.1 200 OK\r\n"
		} else{
			header += "HTTP/1.1 " + c.err + "\r\n"}
		header += "Accept-Ranges: bytes\r\nContent-Lenght: " + size + "\r\nConnection: close\r\n\r\n"
	}	
	return []byte(header)
}

func request_analyzer(get string, size int) Request{
	var req Request
	switch get[:3] {
		case "GET":
			req.method = "GET"
			nb := 0
			for i := 3; i < size; i++ {
				if get[i] == 13 && get[i+1] == 10{
					if get[nb:nb+3] == "GET"{
						req.method = "GET"
						req.path = get[nb+5:i-9]
					} else {if get[nb:nb+5] == "Host:"{
						req.host = get[nb+6:i]
					} else {if get[nb:nb+11] == "User-Agent:"{
						req.user_agent = get[nb+12:i]}}}
					nb = i + 2
				}
			}
		default:
			req.method = ""
	}
	return req
}

func file_recup(name string) ([]byte, string){
	if name == ""{
		name = "index.html"}
	file, err := os.Open(name)
	if err != nil{
		return []byte("<h1>404 Not Found</h1>"), "404 Not Found"
	}
	buffer := make([]byte, 5*1024)
	size, err := file.Read(buffer)
	return buffer[:size], ""
}

func recv(conn net.Conn){
	buffer := make([]byte, 4096)
	data := make([]byte, 5*1024)
	size, err := conn.Read(buffer)
	if err != nil {
		fmt.Println(err)
		return
	}
	request := request_analyzer(string(buffer[:size+11]), size)
	request.ip = conn.RemoteAddr()
	data, request.err = file_recup(request.path)
	_, err = conn.Write(request.Header(fmt.Sprint(len(data))))
	_, err = conn.Write(data)
	conn.Close()
	fmt.Println(request.ip, "\t", request.method, "/" + request.path + "\t", request.user_agent)
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