package file

import "os"

func File_check(path string) (string, int64){
	fi, err := os.Stat(path)
	if os.IsNotExist(err){
		return "404", 30}
	switch mode := fi.Mode(); {
		case mode.IsDir():
			return "301", 30
		case mode.IsRegular():
			return "200", fi.Size()
		default:
			return "404", 30
	}
}

func Dir_check(path string) string{
	fi, err := os.Stat(path)
	mode := fi.Mode()
	if os.IsNotExist(err) || mode.IsRegular(){
		return "404"}
	return "200"
}