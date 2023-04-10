**Usage**

1. install kind
```
$ brew install kind
```
2. create a kind cluster
```
$ kind create cluster --config kind-config.yaml
Creating cluster "kind" ...
 ✓ Ensuring node image (kindest/node:v1.16.3) 🖼
 ✓ Preparing nodes 📦 
 ✓ Writing configuration 📜 
 ✓ Starting control-plane 🕹️ 
 ✓ Installing CNI 🔌 
 ✓ Installing StorageClass 💾 
 ✓ Joining worker nodes 🚜 
Set kubectl context to "kind-kind"
You can now use your cluster with:

kubectl cluster-info --context kind-kind

Thanks for using kind! 😊
```

3. confirm `kind-kind` context
```
$ kube cluster-info --context kind-kind
```
4. docker build for app image
```
$ docker build -t envweb:v1 .
Sending build context to Docker daemon  3.072kB
Step 1/8 : FROM golang:alpine AS builder
 ---> 69cf534c966a
Step 2/8 : WORKDIR /work
 ---> Using cache
 ---> ac2f9cedbec1
Step 3/8 : COPY main.go .
 ---> Using cache
 ---> 1c8da2825756
Step 4/8 : RUN go build -o envweb .
 ---> Using cache
 ---> 6dcb50b4c927
Step 5/8 : FROM alpine
 ---> 965ea09ff2eb
Step 6/8 : WORKDIR /exec
 ---> Using cache
 ---> 59303d271366
Step 7/8 : COPY --from=builder /work/envweb .
 ---> Using cache
 ---> 6aae52df31e3
Step 8/8 : CMD ["./envweb"]
 ---> Using cache
 ---> 7b34cbc36c39
Successfully built 7b34cbc36c39
Successfully tagged envweb:v1

$ docker images
REPOSITORY                              TAG                 IMAGE ID            CREATED             SIZE
envweb                                  v1                  7b34cbc36c39        14 minutes ago      13MB
golang                                  alpine              69cf534c966a        5 days ago          359MB
kindest/node                            v1.16.3             14809a9a48fc        6 days ago          1.26GB
alpine                                  latest              965ea09ff2eb        7 weeks ago         5.55MB
```
5. check images
```
$ docker images
REPOSITORY                              TAG                 IMAGE ID            CREATED             SIZE
envweb                                  v1                  7b34cbc36c39        14 minutes ago      13MB
golang                                  alpine              69cf534c966a        5 days ago          359MB
kindest/node                            v1.16.3             14809a9a48fc        6 days ago          1.26GB
alpine       
```

6. load the iamge to kind cluster so it will be loaded to `kind-worker` and `kind-control-plane`
```
$ kind load docker-image envweb:v1
Image: "envweb:v1" with ID "sha256:7b34cbc36c3902eb8093668638c4259b4ebd56e5bf8e81a3f0ea08b3c40fab7f" not present on node "kind-worker"
Image: "envweb:v1" with ID "sha256:7b34cbc36c3902eb8093668638c4259b4ebd56e5bf8e81a3f0ea08b3c40fab7f" not present on node "kind-control-plane"
```

7. install Istio
```
$ curl -L https://istio.io/downloadIstio | sh -
```

8. deploy Istio to kind cluster with istioctl
```
$ istioctl manifest apply --set profile=default
Preparing manifests for these components:
- Tracing
- Injector
- Kiali
- Cni
- Telemetry
- Pilot
- Base
- PrometheusOperator
- IngressGateway
- EgressGateway
- Policy
- Grafana
- Citadel
- CertManager
- NodeAgent
- Prometheus
- CoreDNS
- Galley

Applying manifest for component Base
Finished applying manifest for component Base
Applying manifest for component Tracing
Applying manifest for component Galley
Applying manifest for component IngressGateway
Applying manifest for component Kiali
Applying manifest for component Policy
Applying manifest for component Citadel
Applying manifest for component Prometheus
Applying manifest for component EgressGateway
Applying manifest for component Pilot
Applying manifest for component Injector
Applying manifest for component Telemetry
Applying manifest for component Grafana
Finished applying manifest for component Tracing
Finished applying manifest for component Prometheus
Finished applying manifest for component Citadel
Finished applying manifest for component Galley
Finished applying manifest for component Kiali
Finished applying manifest for component Injector
Finished applying manifest for component EgressGateway
Finished applying manifest for component Policy
Finished applying manifest for component IngressGateway
Finished applying manifest for component Pilot
Finished applying manifest for component Grafana
Finished applying manifest for component Telemetry
```
9. change the service type of `istio ingress gatway` to `type: NodePort` from `type: LoadBalancer`
```
# Please edit the object below. Lines beginning with a '#' will be ignored,
# and an empty file will abort the edit. If an error occurs while saving this file will be
# reopened with the relevant failures.
#
apiVersion: v1
kind: Service
metadata:
  labels:
    app: istio-ingressgateway
    istio: ingressgateway
  name: istio-ingressgateway
  namespace: istio-system
spec:
  type: NodePort # <- change here from LoadBalancer
  selector:
    app: istio-ingressgateway
  ports:
  - name: http2
    nodePort: 30080 # <- it should be matched with the config on the kind side
    port: 80
    protocol: TCP
    targetPort: 80
  - name: https
...
```

10. check services
```
$ kube -n istio-system get svc
NAME                     TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)                                                                                                                      AGE
grafana                  ClusterIP   10.96.74.124     <none>        3000/TCP                                                                                                                     25h
istio-citadel            ClusterIP   10.97.100.43     <none>        8060/TCP,15014/TCP                                                                                                           25h
istio-egressgateway      ClusterIP   10.102.103.38    <none>        80/TCP,443/TCP,15443/TCP                                                                                                     25h
istio-galley             ClusterIP   10.105.177.225   <none>        443/TCP,15014/TCP,9901/TCP,15019/TCP                                                                                         25h
istio-ingressgateway     NodePort    10.108.210.154   <none>        15020:30693/TCP,80:30080/TCP,443:31800/TCP,15029:32323/TCP,15030:30564/TCP,15031:32108/TCP,15032:31978/TCP,15443:30129/TCP   25h
istio-pilot              ClusterIP   10.107.90.39     <none>        15010/TCP,15011/TCP,8080/TCP,15014/TCP                                                                                       25h
istio-policy             ClusterIP   10.97.222.97     <none>        9091/TCP,15004/TCP,15014/TCP                                                                                                 25h
istio-sidecar-injector   ClusterIP   10.101.148.158   <none>        443/TCP                                                                                                                      25h
istio-telemetry          ClusterIP   10.96.219.241    <none>        9091/TCP,15004/TCP,15014/TCP,42422/TCP                                                                                       25h
jaeger-agent             ClusterIP   None             <none>        5775/UDP,6831/UDP,6832/UDP                                                                                                   25h
jaeger-collector         ClusterIP   10.111.95.238    <none>        14267/TCP,14268/TCP,14250/TCP                                                                                                25h
jaeger-query             ClusterIP   10.107.163.49    <none>        16686/TCP                                                                                                                    25h
kiali                    ClusterIP   10.100.80.83     <none>        20001/TCP                                                                                                                    25h
prometheus               ClusterIP   10.100.148.98    <none>        9090/TCP                                                                                                                     25h
tracing                  ClusterIP   10.108.193.129   <none>        80/TCP                                                                                                                       25h
zipkin                   ClusterIP   10.109.29.40     <none>        9411/TCP                                                                                                                     25h
```

11. deploy the golang web app. note that you need to label the namespace you deploy the app to.
```
$ kube label namespace default istio-injection=enabled
namespace/default labeled
$ kube get ns --show-labels
NAME              STATUS   AGE   LABELS
default           Active   41m   istio-injection=enabled
istio-system      Active   36m   istio-injection=disabled,istio-operator-managed=Reconcile,operator.istio.io/component=Base,operator.istio.io/managed=Reconcile,operator.istio.io/version=1.4.0
kube-node-lease   Active   41m   <none>
kube-public       Active   41m   <none>
kube-system       Active   41m   <none>

$ kube apply -f manifests/web.yml
service/web created
deployment.apps/web-v1 created
deployment.apps/web-v2 created

```
if you don't want to label namespaces you can apply inject manifesta with istioctl
```
$ istioctl kube-inject -f manifests/web.yml | kubectl apply -f -
```



12. check pod
```
$ kube get pods
NAME                      READY   STATUS    RESTARTS   AGE
web-v1-6f4f6c8c7f-nh2cl   2/2     Running   0          9s
web-v2-66bfbb9dd7-dkp7x   2/2     Running   0          9s
```
13. deploy the following API resources of istio :　
- `DestinationRule`
- `VirtualService`
- `Gateway`
```
$ kube apply -f manifests/istio.yml
destinationrule.networking.istio.io/app created
virtualservice.networking.istio.io/app created
gateway.networking.istio.io/web-gateway created
```

14. deploy the rate limiting service
```

$ kube apply -f ratelimit/manifest.yaml 
```
15. apply configMap for setting up descriptors
you can modify the rate limiting settings from here : 
```
apiVersion: v1
kind: ConfigMap
metadata:
  name: ratelimit-config
data:
  config.yaml: |
    domain: ratelimit
    descriptors:
      - key: USER_AGENT
        value: "hogehoge"
        rate_limit:
          unit: minute
          requests_per_unit: 5
      - key: PATH
        rate_limit:
          unit: minute
          requests_per_unit: 5
```
```
$ kube apply -f ratelimit/config.yaml
```
16. apply envoyfilters for global rate limiting
```
$ kube apply -f envoyfileters/manifest.yaml
$ kube apply -f envoyfileters/descriptor.yaml
```


15. check if the rate limiting is working
```
$ for i in `seq 10`; do curl localhost:30080/hoge; echo; done
The app 1 is responding
The app 1 is responding
The app 1 is responding
The app 1 is responding
The app 1 is responding




```

    