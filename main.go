package main

func main() {
	server := CreateApp()
	defer server.Close()
	server.StartServer(":8000")
}
