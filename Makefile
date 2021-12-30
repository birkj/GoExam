program:
	osascript -e 'tell application "Terminal" to do script "cd $(PWD); go run server/server.go -server_id 0 -server_port 8080 -num_replicas 5"'
	osascript -e 'tell application "Terminal" to do script "cd $(PWD); go run server/server.go -server_id 1 -server_port 8081 -num_replicas 5"'
	osascript -e 'tell application "Terminal" to do script "cd $(PWD); go run server/server.go -server_id 2 -server_port 8082 -num_replicas 5"'
	osascript -e 'tell application "Terminal" to do script "cd $(PWD); go run server/server.go -server_id 3 -server_port 8083 -num_replicas 5"'
	osascript -e 'tell application "Terminal" to do script "cd $(PWD); go run server/server.go -server_id 4 -server_port 8084 -num_replicas 5"'

	#osascript -e 'tell application "Terminal" to do script "cd $(PWD); go run client/client.go -id 1"'
	#osascript -e 'tell application "Terminal" to do script "cd $(PWD); go run client/client.go -id 2"'
	#osascript -e 'tell application "Terminal" to do script "cd $(PWD); go run client/client.go -id 3"'
	#osascript -e 'tell application "Terminal" to do script "cd $(PWD); go run client/client.go -id 4"'


