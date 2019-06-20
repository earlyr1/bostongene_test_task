# bostongene_test_task
Golang server which counts md5 hash of file by its url.

Usage example:

Server console:
```bash
go run server.go
Job started, url: https://people.sc.fsu.edu/~jburkardt/data/csv/addresses.csv , id: 1449244776ba2d5ea034dbc6a4b4b413
Job with id 1449244776ba2d5ea034dbc6a4b4b413 requested
Job with id randomstuff requested
Job started, url: randomstuff , id: c52f05387cdaf678809f8083ac04a7d5
Job with id c52f05387cdaf678809f8083ac04a7d5 requested
```

Client console:
```bash
curl -X POST -d "url=https://people.sc.fsu.edu/~jburkardt/data/csv/addresses.csv" localhost:8000/submit
>{"id":"1449244776ba2d5ea034dbc6a4b4b413"}
curl -X GET localhost:8000/check?id=1449244776ba2d5ea034dbc6a4b4b413
>{"md5":"32078264d936c895907f1de187734274","status":finished"}
curl -X GET localhost:8000/check?id=randomstuff
>{"md5":"","status":not exists"}
curl -X POST -d "url=randomstuff" localhost:8000/submit
>{"id":"c52f05387cdaf678809f8083ac04a7d5"}
curl -X GET localhost:8000/check?id=c52f05387cdaf678809f8083ac04a7d5
>{"md5":"","status":error downloading a file"}
```
