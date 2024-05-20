# 5GC Proxy

5GC Proxy is an simple web server that handles Http requests from the openapi client. The requests are proxied to the 5G core network. The server is bundled into the docker image that can be deployed to the K8s cluster using helmcharts.

## Installation

Build docker image.

```bash
docker image built -t proxy-5gc .
```

It is expected that the `config.json` file is preconfigured before the building the image.

## Usage

You can run the docker image on localhost with docker run
```bash
$ docker container run -it 5gc-proxy
Server is running on port 8080
```

More desirable is to deploy the application using helm. Keep in mind that the image should be available to the kubernetes. It can be sideloaded to the `microk8s`.

```bash
docker image save proxy-5gc | microk8s image import
helm install -n free5gc 5gc-proxy charts
```

Both services will be created using the NodePort, that is specified in the `values.yaml`. The UI is then accessible by default on port `30080`.

## License

[MIT](https://choosealicense.com/licenses/mit/)