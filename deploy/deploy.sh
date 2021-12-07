#!/bin/zsh

make

docker build -t jerry9916/demo-csi-driver:latest .

kubectl delete -f deploy/nginx.yaml
kubectl delete -f deploy/pvc.yaml
kubectl delete -f deploy/storageclass.yaml
kubectl delete -f deploy/demo-csi-driver.yaml
kubectl delete -f deploy/rbac.yaml

sleep 15

kubectl apply -f deploy/rbac.yaml
kubectl apply -f deploy/demo-csi-driver.yaml
kubectl apply -f deploy/storageclass.yaml
kubectl apply -f deploy/pvc.yaml
kubectl apply -f deploy/nginx.yaml
