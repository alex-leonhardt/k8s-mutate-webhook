#! /bin/sh
set -uo errexit

export CSR_NAME="${APP}.${NAMESPACE}.svc"
openssl req -x509 -newkey rsa:4096 -keyout ssl/kube-admission.key -out ssl/kube-admission.pem -days 3650
