To Run Docker:

docker-compose up --build (for the first run)

docker-compose up (for future runs without rebuilding)

docker-compose up --build --remove-orphans (to rebuild)

Curl function:

curl --location --request POST 'http://127.0.0.1:8080/hello' --header 'Content-Type: application/json' --data-raw '{"message": "Arbitrary Name or Value"}'        

For windows:

Invoke-WebRequest -Uri 'http://127.0.0.1:8080/hello' -Method POST -Headers @{'Content-Type'='application/json'} -Body '{"message": "hohohohoho"}'