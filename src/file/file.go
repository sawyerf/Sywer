package file

import "os"

func File_check(path string) string{
	fi, err := os.Stat(path)
	if err != nil{
		return "404"
	}
	if os.IsNotExist(err){
		return "404"}
	switch mode := fi.Mode(); {
		case mode.IsDir():
			return "301"
		case mode.IsRegular():
			return "200"
		default:
			return "404"
	}
}

func File_size(path string) int64{
	fi, err := os.Stat(path)
	if err != nil{
		return 0
	}
	return fi.Size()
}

func Dir_check(path string) string{
	fi, err := os.Stat(path)
	if err != nil{
		return "404"
	}
	mode := fi.Mode()
	if os.IsNotExist(err) || mode.IsRegular(){
		return "404"}
	return "200"
}

func Content_Type(name string, size int) string{
	var ptype int
	for i:=0; i < size; i++{
		if name[i] == 46{
			ptype = i + 1
		}
	}
	switch name[ptype:]{
	case "html":
		return "text/html"
	case "css":
		return "text/css"
	case "png":
		return "image/png"
	case "jpeg":
		return "image/jpeg"
	case "mp4":
		return "video/mp4"
	case "pdf":
		return "application/pdf"
	case "zip":
		return "application/zip"
	}
	return ""
}

func Name_Decode(name string, len int) string{
	for i:=0; i < len; i++{
		if len - i > 2 && name[i] == 37{
			if name[i+1] == 50 && name[i+2] == 48{ // %20
				name = name[:i] + " " + name[i+3:]
				i = i-1
				len = len - 3
			}
		}
	}
	return name
}
