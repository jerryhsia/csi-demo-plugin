#!/bin/zsh

make

docker build -t jerry9916/csi-demo-driver:latest .

kubectl delete -f deploy/nginx.yaml
kubectl delete -f deploy/pvc.yaml
kubectl delete -f deploy/storageclass.yaml
kubectl delete -f deploy/csi-demo-driver.yaml

sleep 15

kubectl apply -f deploy/csi-demo-driver.yaml
kubectl apply -f deploy/storageclass.yaml
kubectl apply -f deploy/pvc.yaml
kubectl apply -f deploy/nginx.yaml
