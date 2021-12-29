program:
	osascript -e 'tell application "Terminal" to do script "cd $(PWD); go run server/server.go"'
	osascript -e 'tell application "Terminal" to do script "cd $(PWD); go run client/client.go -id 1"'
	osascript -e 'tell application "Terminal" to do script "cd $(PWD); go run client/client.go -id 2"'
	#osascript -e 'tell application "Terminal" to do script "cd $(PWD); go run client/client.go -id 3"'
	#osascript -e 'tell application "Terminal" to do script "cd $(PWD); go run client/client.go -id 4"'


