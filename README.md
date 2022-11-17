# Description

A gRCP server which provides endpoints following endpoints:
```sh
  GetAggregatedCategoryScoresForPeriod
  GetTicketScoresForPeriod
  GetOverallScoreForPeriod
```

## Usage

* Clone the frontend repo from [Link](https://github.com/valdrinkuchi/score_web)
* Follow the instructions to start the frontend app.
  ```sh
    npm install
    npm run serve
  ```
* Run Envoy Proxy on local terminal
  ```sh
    docker run -p 8080:8080 valdrinkuchi/envoy_proxy:latest
  ``` 
  
* Start the server via 
    ```go
    go get
    go run ./server/server.go
    ```

Navigate to below address after the frontend app and the server is running
```sh
  http://localhost:3000/
```
