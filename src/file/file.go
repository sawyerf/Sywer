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
	fi, _ := os.Stat(path)
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
