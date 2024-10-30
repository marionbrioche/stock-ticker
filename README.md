# Stock Ticker

## Go

Go version `1.23.2`

Build go binary `go build .`

## Docker

Build docker image `docker build .`

Pull the image here `docker pull mabriot/stock-ticker:1.0.0`

## Kubernetes

Deploy k8s manifest using `kubectl apply -f k8s/manifest.yaml`

To access the app either `kubectl port-forward service/stock-ticker-service 8081:8080` open you browser and `http://localhost:8081/api/closing-prices`

Or get the ingress address and add it to your hostfile and go to `http://stockticker.example.com/api/closing-prices`

