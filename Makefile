EXEC_NAME=build/shell

default:
	go build -o ${EXEC_NAME} main.go

password:
	go build -ldflags="-X 'main.password=${PASSWORD}'" -o ${EXEC_NAME} main.go

minify:
	make default
	upx ${EXEC_NAME}

minify_password:
	make password
	upx ${EXEC_NAME}