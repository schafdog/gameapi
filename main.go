package main

import ()

func main() {
	a := App{}
	a.Initialize("127.0.0.1", "gameapi", "user?")
	defer a.Close()
	a.Run(":8000")
}
