#! /bin/sh
set -o errexit

export APP="${1:-mutateme}"
export NAMESPACE="${2:-default}"
export CSR_NAME="${APP}.${NAMESPACE}.svc"

echo "... creating ${app}.key"
openssl genrsa -out ${APP}.key 2048

echo "... creating ${app}.csr"
cat >csr.conf<<EOF
[req]
req_extensions = v3_req
distinguished_name = req_distinguished_name
[req_distinguished_name]
[ v3_req ]
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names
[alt_names]
DNS.1 = ${APP}
DNS.2 = ${APP}.${NAMESPACE}
DNS.3 = ${CSR_NAME}
DNS.4 = ${CSR_NAME}.cluster.local
EOF
openssl req -new -key ${APP}.key -subj "/CN=${CSR_NAME}" -out ${APP}.csr -config csr.conf

echo "... deleting existing csr, if any"
kubectl delete csr ${CSR_NAME} 2>/dev/null || true
	
echo "... creating kubernetes CSR object"
kubectl create -f - <<EOF
apiVersion: certificates.k8s.io/v1beta1
kind: CertificateSigningRequest
metadata:
  name: ${CSR_NAME}
spec:
  groups:
  - system:authenticated
  request: $(cat ${APP}.csr | base64 | tr -d '\n')
  usages:
  - digital signature
  - key encipherment
  - server auth
EOF

SECONDS=0
while true; do
  echo "... waiting for csr to be present in kubernetes"
  kubectl get csr ${CSR_NAME} > /dev/null 2>&1
  if [ "$?" -eq 0 ]; then
      break
  fi
  if [[ $SECONDS -ge 60 ]]; then
    echo "[!] timed out waiting for csr"
    exit 1
  fi
  sleep 2
done

kubectl certificate approve ${CSR_NAME}

SECONDS=0
while true; do
  echo "... waiting for serverCert to be present in kubernetes"
  serverCert=$(kubectl get csr ${CSR_NAME} -o jsonpath='{.status.certificate}')
  if [[ $serverCert != "" ]]; then 
    break
  fi
  if [[ $SECONDS -ge 60 ]]; then
    echo "[!] timed out waiting for serverCert"
    exit 1
  fi
  sleep 2
done

echo "... creating ${app}.pem cert file"
echo ${serverCert} | openssl base64 -d -A -out ${APP}.pem
