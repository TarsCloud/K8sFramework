#! /bin/bash

if [ -z "$SERVICE_NAME" ]; then
  echo "read \$SERVICE_NAME error"
  exit 255
fi

declare OPENSSL_WORK_DIR=/etc/tarscontrol-cert
declare CA_PEM_FILE=${OPENSSL_WORK_DIR}/ca.pem
declare CA_KEY_FILE=${OPENSSL_WORK_DIR}/ca.key
declare CERT_KEY_FILE=${OPENSSL_WORK_DIR}/cert.key
declare CERT_CSR_FILE=${OPENSSL_WORK_DIR}/cert.csr
declare CERT_PEM_FILE=${OPENSSL_WORK_DIR}/cert.pem
declare NAMESPACE_FILE=/var/run/secrets/kubernetes.io/serviceaccount/namespace

NAMESPACE_VALUE=$(cat ${NAMESPACE_FILE})

TargetCertCN=${SERVICE_NAME}.${NAMESPACE_VALUE}.svc

if [ -z "$NAMESPACE_VALUE" ]; then
  echo "read \$NAMESPACE_FILE error"
  exit 255
fi

if ! cd $OPENSSL_WORK_DIR; then
  echo "cd \$OPENSSL_WORK_DIR error"
  exit 255
fi

#  generate $CERT_KEY_FILE
if ! openssl genrsa -out $CERT_KEY_FILE 2048; then
  echo "generate \$CERT_KEY_FILE error"
  exit 255
fi
#

# generate $CERT_CSR_FILE
if ! openssl req -new -key $CERT_KEY_FILE -out $CERT_CSR_FILE -extensions "v3_req" -subj "/CN=${TargetCertCN}" -reqexts SAN -config <(cat /etc/ssl/openssl.cnf <(printf "[SAN]\nsubjectAltName=DNS:%s" "${TargetCertCN}")); then
  echo "generate \$CERT_CSR_FILE error"
  exit 255
fi
#
# generate $CERT_PEM_FILE
if ! openssl x509 -days 365 -req -in $CERT_CSR_FILE -CA $CA_PEM_FILE -CAkey $CA_KEY_FILE -CAcreateserial -out $CERT_PEM_FILE -extfile <(printf "subjectAltName=DNS:%s" "${TargetCertCN}"); then
  echo "generate \$CERT_PEM_FILE error"
  exit 255
fi
#

# verify $CERT_PEM_FILE
if ! openssl verify -CAfile $CA_PEM_FILE $CERT_PEM_FILE; then
  echo "verify \$CERT_PEM_FILE error"
  exit 255
fi

TARSCONTROL_EXECUTION_FILE="/usr/local/app/tars/tarscontrol/bin/tarscontrol"
exec ${TARSCONTROL_EXECUTION_FILE}
