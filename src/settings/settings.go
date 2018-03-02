package settings

import ("fmt"
		"os")

type Settings struct{
	Port string
	Path string
	Index string
	Error404 string
	Error301 string
}

func Line_analyzer(set Settings, data string, i int) Settings{
	if len(data) == 0{
		return set
	}
	if data[0:5] == "port	"{
		set.Port = data[5:i]
	} else if data[0:5] == "path	"{
		set.Path = data[5:i]
	} else if data[0:6] == "index	"{
		set.Index = data[6:i]
	} else if data[0:6] == "error_"{
		switch data[6:10] {
		case "404	":
			set.Error404 = data[10:i]
		case "301	":
			set.Error301 = data[10:i]
		default:
			fmt.Println("[!] This type of error is not understand : \"" + data[0:i])
			os.Exit(0)
		}
	}
	return set
}

func Recup(path string) Settings{
	var set Settings
	//Open file
	fi, err := os.Stat(path)
	if err != nil{
		set.Port = "80"
		return set
	}
	siz := fi.Size()
	buffer := make([]byte, siz)
	file, _ := os.Open(path)
	size, _ := file.Read(buffer)
	var data string = string(buffer[:size])
	nb := 0
	for i := 0; i!=size; i++{
		if data[i] == 13 {
			set = Line_analyzer(set, data[nb:i], i-nb)
			nb = i + 2
			i++
		} else if data[i] == 10{
			set = Line_analyzer(set, data[nb:i], i-nb)
			nb = i + 1
		}
	}
	if set.Port == ""{
		set.Port = "80"
	}
	if set.Index == ""{
		set.Index =	 "index.html"
	}
	fmt.Println(set)
	return set
}
