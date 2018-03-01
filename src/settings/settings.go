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
	var data string = string(buffer[:size]) + "       "

	nb := 0
	fmt.Println(data)
	for i:=0; i!=size; i++{
		if data[i] == 13{
			if data[nb:nb+5] == "port	"{
				set.Port = data[nb+5:i]
			} else if data[nb:nb+5] == "path	"{
				set.Path = data[nb+5:i]
			} else if data[nb:nb+6] == "index	"{
				set.Index = data[nb+6:i]
			} else if data[nb:nb+6] == "error_"{
				switch data[nb+6:nb+10] {
				case "404	":
					set.Error404 = data[nb+10:i]
				case "301	":
					set.Error301 = data[nb+10:i]
				default:
					fmt.Println("[!] This type of error is not understand : \"" + data[nb:i])
					os.Exit(0)}
			}
			nb = i + 2
		}
	}
	if set.Port == ""{
		set.Port = "80"
	}
	fmt.Println(set)
	return set
}
