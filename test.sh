go run server.go &
sleep 1
curl -X POST -d "url=https://people.sc.fsu.edu/~jburkardt/data/csv/addresses.csv" localhost:8000/submit
curl -X GET localhost:8000/check?id=1449244776ba2d5ea034dbc6a4b4b413
curl -X GET localhost:8000/check?id=randomstuff
curl -X POST -d "url=randomstuff" localhost:8000/submit
curl -X GET localhost:8000/check?id=c52f05387cdaf678809f8083ac04a7d5