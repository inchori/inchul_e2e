#!/bin/bash

SCRIPT_DIR=$(cd `dirname $0` && pwd)

rm -rf /oracle/.oracle
unset SGX_AESM_ADDR
OE_SIMULATION=1 ego run /usr/bin/oracled init

cat ${SCRIPT_DIR}/config.toml | \
  sed "s|__CHAIN_ID__|${CHAIN_ID}|g" | \
  sed "s|__ORACLE_MNEMONIC__|${ORACLE_MNEMONIC}|g" \
  > /oracle/.oracle/config.toml


OE_SIMULATION=1 ego run /usr/bin/oracled gen-oracle-key \
 --trusted-block-height 1 --trusted-block-hash "${TRUSTED_BLOCK_HASH}"

sleep 100