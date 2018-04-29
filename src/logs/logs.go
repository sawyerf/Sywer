package logs

import ("os"
        "../request"
	"time")

func Log(path string, req request.Request){
    file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
    if err != nil{
        return
    }
    defer file.Close()
    _, err = file.Write([]byte("[" + date() + "] " + req.Ip + "\t" + req.Method + " /" + req.Path + " " + req.Err + "\t" + req.User_agent + "\n"))
    if err != nil{
        return
    }
}

func date() string{
	date := time.RFC822
	return date[:len(date)-4]
}
