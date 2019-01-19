package main

func main() {
	server := CreateApp()
	server.StartServer(":8000")
}
