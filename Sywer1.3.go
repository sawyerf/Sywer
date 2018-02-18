package main

import ("fmt"
		"net"
		"log"
		"os")

const src_path string = "/srv/http/" 

type Request struct{
	ip net.Addr
	method string
	path string
	type_path bool//File = True and Directory = False
	host string
	user_agent string
	connection string
	err string
}


func (c Request) Header(size string) []byte{
	var header string = ""
	if c.method == "GET"{
		header += "HTTP/1.1 " + c.err + "\r\n"
		if c.err == "301 Moved Permanently"{
			header += "Location: http://" + c.host + "/" + c.path + "/\r\n"}
		header += "Accept-Ranges: bytes\r\nContent-Lenght: " + size + "\r\nConnection: close\r\n\r\n"
	}	
	return []byte(header)
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

func file_recup(name string) ([]byte, string){
	if name == ""{
		name = "index.html"}
	fi, err := os.Stat(src_path + name)
	if os.IsNotExist(err){
		return []byte("<h1>404 Not Found</h1>"), "404 Not Found"}
	switch mode := fi.Mode(); {
		case mode.IsDir():
			return []byte("<h1>301 Moved Permanently</h1>"), "301 Moved Permanently"
		case mode.IsRegular():
			file, _ := os.Open(src_path + name)
			buffer := make([]byte, fi.Size())
			size, _ := file.Read(buffer)
	return buffer[:size], "200 OK"
		default:
			return []byte("<h1>404 Not Found</h1>"), "404 Not Found"
	}
}

func directory_recup(name string) ([]byte, string){
	fi, err := os.Stat(src_path + name)
	mode := fi.Mode()
	if os.IsNotExist(err) || mode.IsRegular(){
		return []byte("<h1>404 Not Found</h1>"), "404 Not Found"}
	return []byte("<h1>Index of /" + name + "</h1>"), "200 OK"
}

func recv(conn net.Conn){
	buffer := make([]byte, 4096)
	data := make([]byte, 5*1024)
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
	request.ip = conn.RemoteAddr()
	if request.type_path{
		data, request.err = file_recup(request.path)
	} else{
		data, request.err = directory_recup(request.path)}
	_, err = conn.Write(request.Header(fmt.Sprint(len(data))))
	_, err = conn.Write(data)
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