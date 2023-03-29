#!/bin/bash

SCRIPT_DIR=$(cd `dirname $0` && pwd)

rm -rf /oracle/.oracle
unset SGX_AESM_ADDR
OE_SIMULATION=1 ego run /usr/bin/oracled init

#CHAIN_ID="testing"
#ORACLE_MNEMONIC="gasp shy describe man hello blossom motor monkey seven mule shallow almost bunker hello wife clarify tissue best actress hub wisdom crane ridge heavy"
#PANACEA_HOST="127.0.0.1"


cat ${SCRIPT_DIR}/config.toml | \
  sed "s|__CHAIN_ID__|${CHAIN_ID}|g" | \
  sed "s|__ORACLE_MNEMONIC__|${ORACLE_MNEMONIC}|g" | \
  sed "s|__PANACEA_HOST__|${PANACEA_HOST}|g" \
  > /oracle/.oracle/config.toml

cat /oracle/.oracle/config.toml

sleep time