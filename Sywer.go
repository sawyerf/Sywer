//Version 0.22
package main

import ("fmt"
		"net"
		"log"
		"./src/file"
		"./src/request"
		"./src/settings"
		"./src/logs")

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

func recv(conn net.Conn, set settings.Settings, wFile *bool){
	buffer := make([]byte, 1024)
	var size_data int = 0
	var req_data string

	//Receive the request
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

	//Analyze the request
	req := request.Request_analyzer(req_data, size_data)
	if req.Method == ""{
		conn.Close()
		return
	}
	switch req.Path{
	case "":
		req.Path = set.Index
	case "favicon.ico":
		if set.Ico != ""{
			req.Path = set.Ico
		}
	}
	req.Ip = Ip(conn.RemoteAddr())
	req.Path = file.Name_Decode(req.Path, len(req.Path))
	// File and Directory check
	if req.Type_path{
		req.Err = file.File_check(set.Path + req.Path)
	} else{
		req.Err = file.Dir_check(set.Path + req.Path)
	}
	if req.Err == "301" && req.Host == ""{
		req.Err = "400"
	}

	//Send the reply
	_, err := conn.Write(req.Header(set))
	if err != nil{
		return
	}
	req.Data(conn, set)
	conn.Close()
	for{
		if *wFile{
			*wFile = false
			logs.Log(set.Logs, req)
			*wFile = true
			break
		}
	}
	buffer = buffer[:0]
}

func main(){
	var wFile bool = true

	//Settings
	set := settings.Recup("settings.swy")
	if !set.Found{
		set = settings.Recup("/data/data/com.termux/files/usr/var/lib/sywer/settings.swy")
		if !set.Found{
			set = settings.Recup("/var/lib/sywer/settings.swy")
			if !set.Found{
				fmt.Println("[!] File Not Found (.../settings.swy)")
			}
		}
	}

	//Init server
	server, err := net.Listen("tcp", ":" + set.Port)
	if err != nil {
			log.Fatalln(err)
	}
	fmt.Println("[*]Start")

	//affect the socket to a thread
	for{
		conn, err := server.Accept()
		if err != nil {
			fmt.Println("[*]Erreur", err)
		}

		go recv(conn, set, &wFile)
	}
}
