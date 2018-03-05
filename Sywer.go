//Version 0.22
package main

import ("fmt"
		"net"
		"log"
		"./src/file"
		"./src/request"
		"./src/settings")


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
	buffer := make([]byte, 1024)
	var size_data int = 0
	var req_data string
	for{
		size, err := conn.Read(buffer)
		if err != nil {
			fmt.Println(err)
			conn.Close()
			return
		}
		size_data = size_data + size
		req_data = req_data + string(buffer[:size])
		if req_data[size_data-4:size_data] == "\r\n\r\n"{
			break
		}
	}

	req := request.Request_analyzer(req_data, size_data)
	if req.Method == ""{
		conn.Close()
		return
	} else if req.Path == ""{
		req.Path = set.Index
	}

	req.Ip = Ip(conn.RemoteAddr())
	// File and Directory check
	if req.Type_path{
		req.Err = file.File_check(set.Path + req.Path)
	} else{
		req.Err = file.Dir_check(set.Path + req.Path)}

	_, err := conn.Write(req.Header(set))
	if err != nil{
		return
	}
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
