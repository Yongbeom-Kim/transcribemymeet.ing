#!/bin/bash
BASEDIR=$(git rev-parse --show-toplevel)

FILES=(
    $BASEDIR/secrets/gcloud_service_key.json
    $BASEDIR/tf_backend/terraform.tfstate
    $BASEDIR/tf_backend/terraform.tfstate.backup
)

ENCRYPTED_SUFFIX=".enc"

encrypt() {
    for FILE in ${FILES[@]}; do
        echo "Encrypting $FILE"
        sops -e --input-type json --output-type json $FILE > $FILE$ENCRYPTED_SUFFIX
    done
}

decrypt() {
    for FILE in ${FILES[@]}; do
        echo "Decrypting $FILE"
        sops -d --input-type json --output-type json $FILE$ENCRYPTED_SUFFIX > $FILE
    done
}

if [ $# -lt 1 ]; then
  echo "Usage: $0 [encrypt|decrypt]"
  exit 1
fi

CMD=$1
shift

case $CMD in
  encrypt)
    encrypt $@
    ;;
  decrypt)
    decrypt $@
    ;;
  *)
    echo "Usage: $0 [encrypt|decrypt]"
    exit 1
    ;;
esac