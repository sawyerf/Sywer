package main

import ("fmt"
		"net"
		"log"
		"os")

type Request struct{
	method string
	raw string
	path string
	host string
	user_agent string
	connection string
}

func request_analyzer(get string, size int) Request{
	var req Request = Request{raw: get}
	switch get[:3] {
	case "GET":
		req.method = "GET"
		for i := 3; i <= size; i++ {
			if get[i] == 10{
				req.path = get[5:i-10]
				return req
			}
		}
	default:
		req.method = ""
	}
	return req
}

func file_recup(name string) ([]byte, bool){
	file, err := os.Open(name[0:])
	if err != nil{
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
	request := request_analyzer(string(buffer[:size]), size)
	data, eror := file_recup(request.path)				
	if eror{
		_, err = conn.Write(data)}
	conn.Close()
	fmt.Println(conn, ":", request.method, request.path)
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
		fmt.Println("[*]Nouvelle ecoute", conn)
		go recv(conn)
	}
}