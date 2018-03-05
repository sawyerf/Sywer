package request

import ("fmt"
		"os"
		"io/ioutil"
		"net"
		"../file"
		"../settings")

type Request struct{
	//info request
	Ip string
	Method string
	Host string
	User_agent string
	Connection string
	//File To send
	Path string
	Type_path bool //File = True and Directory = False
	Err string
}

func (c Request) Header(set settings.Settings) []byte{
	var header string = ""
	if c.Method == "GET"{
		var size int64
		header += "HTTP/1.1 "
		switch c.Err{
		case "200":
			header += "200 OK\r\n"
			size = file.File_size(set.Path + c.Path)
		case "404":
			header += "404 Not Found"
			if set.Error404 != ""{
				size = file.File_size(set.Error404)}
		case "301":
			header += "301 Moved Permanently\r\n"
			header += "Location: http://" + c.Host + "/" + c.Path + "/\r\n"
			if set.Error301 != ""{
				size = file.File_size(set.Error301)
			}
		}
		if 0 < size{
			header += "Accept-Ranges: bytes\r\nContent-Lenght: " + fmt.Sprint(size) + "\r\nConnection: close\r\n\r\n"
		} else {
			header += "Accept-Ranges: bytes\r\nonnection: close\r\n\r\n"
		}
	}
	return []byte(header)
}

func (c Request) Data(conn net.Conn, set settings.Settings){
	if c.Method == "GET"{
		switch c.Err{
		case "200":
			if !c.Type_path{
				_, _ = conn.Write([]byte("<h1>Index Of " + c.Path + "</h1>\n<ul>"))
				files, _ := ioutil.ReadDir(set.Path + c.Path)
				for _, file := range files{
					_, _ = conn.Write([]byte("<li><a href=\"" + file.Name() + "\">" + file.Name() + "</a></li>\n"))
				}
				_, _ = conn.Write([]byte("</ul>"))

				return
			}else{
				buffer := make([]byte, 1024)
				file, _ := os.Open(set.Path + c.Path)
				for{
					Size, _ := file.Read(buffer)
					if Size == 0{
						buffer = buffer[:0]
						file.Close()
						return
					} else{
						_, _ = conn.Write(buffer[:Size])}
				}
			}
		case "404":
			if c.Path == "ip"{
					conn.Write([]byte("<h1>" + c.Ip + "</h1>"))
					return}
			if set.Error404 != ""{
				buffer := make([]byte, 1024)
				file, err := os.Open(set.Error404)
				if err != nil{
					conn.Write([]byte("<h1>404 Not Found</h1>"))
					return
				}
				for{
					Size, _ := file.Read(buffer)
					if Size == 0{
						buffer = buffer[:0]
						file.Close()
						return
					} else{
						_, _ = conn.Write(buffer[:Size])}
				}
			}else {
				conn.Write([]byte("<h1>404 Not Found</h1>"))
			}
		case "301":
			if set.Error301 != ""{
				buffer := make([]byte, 1024)
				file, err := os.Open(set.Error301)
				if err != nil {
					conn.Write([]byte("<h1>301 Moved Permanently</h1>"))
				}
				for{
					Size, _ := file.Read(buffer)
					if Size == 0{
						buffer = buffer[:0]
						file.Close()
						return
					} else{
						_, _ = conn.Write(buffer[:Size])}
				}
			}else {
				conn.Write([]byte("<h1>301 Moved Permanently</h1>"))
			}
		default:
			return
		}
	}
}

func Line_request(req Request, data string, size int) Request {
	if size <= 0{
		return req
	}
	if data[0:3] == "GET"{
		req.Method = "GET"
		req.Path = data[5:size-9]
		req.Type_path = true
		if req.Path != ""{
			if req.Path[len(req.Path)-1] == 47{
				req.Type_path = false
			}
		}
		return req
	} else if data[0:5] == "Host:"{
		req.Host = data[6:size]
		return req
	} else if data[0:11] == "User-Agent:"{
		req.User_agent = data[12:size]
		return req
	}
	return req
}

func Request_analyzer(get string, Size int) Request{
	var req Request
	nb := 0
	for i := 3; i < Size; i++{
		if get[i] == 13 && get[i+1] == 10{
			req = Line_request(req, get[nb:i], i-nb)
			nb = i + 2
		}
	}
	return req
}
