compile-all-platform: mdserver_linux-x86_64 mdserver_linux-arm mdserver_win-x86_64.exe

mdserver_linux-x86_64:
	GOOS=linux GOARCH=amd64 go build -o out/mdserver_linux-x86_64

mdserver_linux-arm:
	GOOS=linux GOARCH=arm go build -o out/mdserver_linux-arm

mdserver_win-x86_64.exe:
	GOOS=windows GOARCH=amd64 go build -o out/mdserver_win-x86_64.exe

clean:
	rm out/*
