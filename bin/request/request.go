package request

import ("fmt"
		"os"
		"io/ioutil"
		"net")

type Request struct{
	//info request
	Ip string
	Method string
	Host string
	User_agent string
	Connection string
	//File To send
	Src_path string
	Path string
	Size int64
	Type_path bool //File = True and Directory = False
	Err string
}

func (c Request) Header() []byte{ 
	var header string = ""
	if c.Method == "GET"{
		header += "HTTP/1.1 "
		switch c.Err{
		case "200":
			header += "200 OK\r\n"
		case "404":
			header += "404 Not Found"
		case "301":
			header += "301 Moved Permanently\r\n"
			header += "Location: http://" + c.Host + "/" + c.Path + "/\r\n"
		}
		if 0 < c.Size{
		header += "Accept-Ranges: bytes\r\nContent-Lenght: " + fmt.Sprint(c.Size) + "\r\nConnection: close\r\n\r\n"
		} else {
			header += "Accept-Ranges: bytes\r\nonnection: close\r\n\r\n"
		}
	}	
	return []byte(header)
}

func (c Request) Data(conn net.Conn){
	if c.Method == "GET"{
		switch c.Err{
		case "200":
			if !c.Type_path{
				_, _ = conn.Write([]byte("<h1>Index Of " + c.Path + "</h1>\n<ul>"))
				files, _ := ioutil.ReadDir(c.Src_path + c.Path)
				for _, file := range files{
					_, _ = conn.Write([]byte("<li><a href=\"" + file.Name() + "\">" + file.Name() + "</a></li>\n"))
				}
				_, _ = conn.Write([]byte("</ul>"))
				return
			}
			buffer := make([]byte, 1024)
			file, _ := os.Open(c.Src_path + c.Path)
			for{
				Size, _ := file.Read(buffer)
				if Size == 0{
					return
				} else{
					_, _ = conn.Write(buffer[:Size])}
			}
		case "404":
			if c.Path == "ip"{
				conn.Write([]byte("<h1>" + c.Ip + "</h1>"))
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

func Request_analyzer(get string, Size int) Request{
	var req Request
	nb := 0
	for i := 3; i < Size; i++ {
		if get[i] == 13 && get[i+1] == 10{
			if get[nb:nb+3] == "GET"{
				req.Method = "GET"
				req.Path = get[nb+5:i-9]
				req.Type_path = true
				if req.Path == ""{
					req.Path = "index.html"
				} else{ if req.Path[len(req.Path)-1] == 47{
					req.Type_path = false}}
			} else {if get[nb:nb+5] == "Host:"{
				req.Host = get[nb+6:i]
			} else {if get[nb:nb+11] == "User-Agent:"{
				req.User_agent = get[nb+12:i]}}}
			nb = i + 2
		}
	}
	return req
}
