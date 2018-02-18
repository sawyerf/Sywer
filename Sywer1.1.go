package main

import ("fmt"
		"net"
		"log"
		"os")

type Request struct{
	method string
	path string
	host string
	user_agent string
	connection string
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

func file_recup(name string) ([]byte, bool){
	file, err := os.Open(name)
	if err != nil{
		fmt.Println(err, name)
		return []byte(""), false
	}
	buffer := make([]byte, 5*1024)
	size, err := file.Read(buffer)
	return buffer[:size], true
}

func recv(conn net.Conn){
	buffer := make([]byte, 4096)
	size, err := conn.Read(buffer)
	if err != nil {
		fmt.Println(err)
		return
	}
	request := request_analyzer(string(buffer[:size+11]), size)
	data, eror := file_recup(request.path)
	if eror{
		_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-length: " + string(len(data)) + "\r\nConnection: Close\r\n\r\n"))
		_, err = conn.Write(data)}
	conn.Close()
	fmt.Println(conn, ":", request.method, request.path, request.user_agent)
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