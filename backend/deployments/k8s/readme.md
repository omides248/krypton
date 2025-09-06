# --------------------- Install k3s ---------------------

# 1. Install
curl -sfL https://get.k3s.io | sh -s - --disable=traefik


# 2. kubectl config k3s 
mkdir -p ~/.kube
sudo cp /etc/rancher/k3s/k3s.yaml ~/.kube/config
sudo chown $(id -u):$(id -g) ~/.kube/config

# 3. Test nodes
kubectl get nodes


# 4. Uninstall k3s
sudo /usr/local/bin/k3s-uninstall.sh
sudo rm -rf /etc/rancher/k3s
sudo rm -rf /var/lib/rancher/k3s
sudo rm -rf /var/lib/kubelet
sudo rm -rf ~/.kube


# --------------------- Install Helm ---------------------
curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
chmod 700 get_helm.sh
./get_helm.sh
helm version


# --------------------- Install Ingress Controller ---------------------
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update

helm -n ingress-nginx install ingress-nginx ingress-nginx/ingress-nginx --create-namespace


# --------------------- Config Mongo ---------------------

# 1. Connect to mongo shell
mongosh

# 2. Use admin database
use admin

# 3. Create new root user for admin database
db.createUser({
user: "omides248",
pwd: "123123",
roles: [ { role: "root", db: "admin" } ]
})

# 4. Enable authentication for mongodb
nano /etc/mongod.conf
security:
  authorization: enabled

# 5. Restart for set config
sudo systemctl restart mongod

# 6. Test authentication
mongosh -u omides248 -p 123123 --authenticationDatabase admin


# --------------------- Apply config ---------------------
kubectl apply -k k8s/overlays/dev
kubectl apply -k k8s/overlays/stg
kubectl apply -k k8s/overlays/prod

# --------------------- Goland connect to kubectl ---------------------
nano .kube/config

C:\Users\omide\.kube\config


# --------------------- Install kubernetes dashboard ---------------------

# 1. Install
helm repo add kubernetes-dashboard https://kubernetes.github.io/dashboard/
helm upgrade --install kubernetes-dashboard kubernetes-dashboard/kubernetes-dashboard --create-namespace --namespace kubernetes-dashboard

# 2. Config
kubectl -n kubernetes-dashboard create serviceaccount admin-user
kubectl create clusterrolebinding admin-user --clusterrole=cluster-admin --serviceaccount=kubernetes-dashboard:admin-user
kubectl -n kubernetes-dashboard create token admin-user

# 3. Run
kubectl -n kubernetes-dashboard port-forward svc/kubernetes-dashboard-kong-proxy 8443:443 --address 0.0.0.0

# --------------------- Install MinIO ---------------------

# 1. Download the MinIO RPM
wget https://dl.min.io/server/minio/release/linux-amd64/minio_20250723155402.0.0_amd64.deb -O minio.deb
sudo dpkg -i minio.deb

mkdir -p /mnt/disk1/minio
chown -R minio-user:minio-user /mnt/disk1/

nano /etc/default/minio
--->
## Volume to be used for MinIO server.
MINIO_VOLUMES="/mnt/disk1/minio"

## Use if you want to run MinIO on a custom port.
#MINIO_OPTS="--address :9198 --console-address :9199"

## Root user for the server.
MINIO_ROOT_USER=minioadmin

## Root secret for the server.
MINIO_ROOT_PASSWORD=minioadmin123123

## set this for MinIO to reload entries with 'mc admin service restart'
#MINIO_CONFIG_ENV_FILE=/etc/default/minio
<---

systemctl start minio

# --------------------- Connect to pods cli ---------------------
kubectl exec -it catalog-deployment-7b46c7664d-6gb88 -- /bin/sh



















