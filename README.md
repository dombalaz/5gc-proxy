# 5GC Proxy

5GC Proxy is an simple web server that handles Http requests from the openapi client. The requests are proxied to the 5G core network. The server is bundled into the docker image that can be deployed to the K8s cluster using helmcharts.

## Installation

Build docker image.

```bash
docker image built -t 5gc-proxy .
```

## Usage

You can run the docker image on localhost with docker run
```bash
$ docker container run -it 5gc-proxy
Server is running on port 8080
```

More desirable is to deploy the application, with the 

```bash
helm install -n free5gc 5gc-proxy charts
```

## License

[MIT](https://choosealicense.com/licenses/mit/)