export CURR_DIR=$(pwd)
export PORTWORX_GRPC_URL=localhost:9020
export BROKER_URL=http://localhost:8080
export KUBECONFIG=/Users/camartinez/.kube/vanilla-lehi


# Forward K8s 

kubectl port-forward svc/portworx-api 9020:9020 -n kube-system &

# Start Broker

cd ${CURR_DIR}/broker 

go run cmd/api/*.go &

cd ${CURR_DIR}

# Start Frontend
sleep 5
cd ${CURR_DIR}/frontend

go run cmd/web/*.go &

cd ${CURR_DIR}

