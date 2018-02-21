//Version 1.5
package main

import ("fmt"
		"net"
		"log"
		"./bin/file"
		"./bin/request")

const src_path string = "/srv/http/" //the directory where you have the files you want to share 

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

	req := request.Request_analyzer(string(buffer[:size+11]), size)
	fmt.Println(req)
	if req.Method == ""{
		conn.Close()
		return}

	req.Ip = Ip(conn.RemoteAddr())
	req.Src_path = src_path
	// File and Directory check
	if req.Type_path{
		req.Err, req.Size = file.File_check(req.Src_path + req.Path)
	} else{
		req.Err = file.Dir_check(req.Src_path + req.Path)}

	_, err = conn.Write(req.Header())
	req.Data(conn)
	conn.Close()
	fmt.Println(req.Ip, "\t", req.Method, "/" + req.Path + "\t", req.User_agent, "\t", req.Err)
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

