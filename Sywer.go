//Version 1.5
package main

import ("fmt"
		"net"
		"log"
		"./bin/file"
		"./bin/request"
		"./bin/settings")


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

func recv(conn net.Conn, set settings.Settings){
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
	// File and Directory check
	if req.Type_path{
		req.Err = file.File_check(set.Path + req.Path)
	} else{
		req.Err = file.Dir_check(set.Path + req.Path)}

	_, err = conn.Write(req.Header(set))
	req.Data(conn, set)
	conn.Close()
	fmt.Println(req.Ip, "\t", req.Method, "/" + req.Path + "\t", req.User_agent, "\t", req.Err)
	buffer = buffer[:0]
}

func main(){
	set := settings.Recup("settings.swy")
	server, err := net.Listen("tcp", ":" + set.Port)
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
		go recv(conn, set)
	}
}

