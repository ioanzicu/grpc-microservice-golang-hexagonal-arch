SAMPLE Microservice with gRPC
Hexagonal Architecture



# Testing

Testing gRPC using [grpcurl](https://github.com/fullstorydev/grpcurl)



## Mysql Docker image:

    docker run -p 3306:3306 \
        -e MYSQL_ROOT_PASSWORD=verysecretpass \
        -e MYSQL_DATABASE=payment mysql

## Check port usage
    
    lsof -i tcp:3000
    
    netstat -vanp tcp | grep 3000
    
    kill -9 <PID>
    

## Update Go Pakcages / Modules

    go get -u ./... && go mod tidy


## Tag the commit


### Create the tag

    git tag v1.0

### Push tags

    git push origin v1.0

    OR 

    git push --tags

### List tags

    git tag

