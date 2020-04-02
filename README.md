# k8s-mutate-webhook

A playground to try build a crude k8s mutating webhook; the goal is to mutate a Pod CREATE request to _always_ use a debian image and by doing this, learning more about
the k8s api, objects, etc. - eventually figure out how scalable this is (could be made) if one had 1000 pods to schedule (concurrently)

This is a companion to the blog post [Writing a very basic kubernetes mutating admission webhook](https://medium.com/ovni/writing-a-very-basic-kubernetes-mutating-admission-webhook-398dbbcb63ec)  

## build 

```
make
```

## test

```
make test
```

## ssl/tls

the `ssl/` dir contains a script to create a self-signed certificate, not sure this will even work when running in k8s but that's part of figuring this out I guess

_NOTE: the app expects the cert/key to be in `ssl/` dir relative to where the app is running/started and currently is hardcoded to `mutateme.{key,pem}`_

```
cd ssl/ 
make 
```

## docker

to create a docker image .. 

```
make docker
```

it'll be tagged with the current git commit (short `ref`) and `:latest`

don't forget to update `IMAGE_PREFIX` in the Makefile or set it when running `make`

### images

[`alexleonhardt/k8s-mutate-webhook`](https://cloud.docker.com/repository/docker/alexleonhardt/k8s-mutate-webhook)


## watcher

useful during devving ... 

```
watcher -watch github.com/alex-leonhardt/k8s-mutate-webhook -run github.com/alex-leonhardt/k8s-mutate-webhook/cmd/
```

## Running in docker-for-mac

```bash
cd ssl && make && cd -
make docker
sed -i '' 's/imagePullPolicy: Always/imagePullPolicy: Never/' deploy/webhook.yaml # use local image
sed -i '' "s/caBundle:.*/caBundle: $(cat ssl/mutateme.pem | base64)/" deploy/webhook.yaml # use local CA 
kubectl label namespace default mutateme=enabled
kubectl apply -f deploy/webhook.yaml

# make sure it's running ...
kubectl get pods
kubectl logs <PDO> --follow

# create example pod to see it working
kubectl apply -f pod.yaml
kubectl get pod c7m -o yaml | grep image: # should be debian
```

## kudos

- [https://github.com/morvencao/kube-mutating-webhook-tutorial](https://github.com/morvencao/kube-mutating-webhook-tutorial)
