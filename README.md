Deployment steps

1. ssh to server

2. go to directory: /home/ubuntu/go/src/authorization

3. run command
    git pull
   to get latest change

4. run command
    go get ./...
   to install all dependecies

5. run command
    go build *.go
   to generate app file

6. run command
    lsof -i :9000
   to find running pid

7. kill running pid by command
    kill <pid>

8. run command
    ./app &
   to start service
    
