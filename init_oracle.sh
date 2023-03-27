#!/bin/bash

CONFIG_PATH="/home_mnt/.oracle"

rm -rf ~/.${CONFIG_PATH}
ego run oracled init

cat ${CONFIG_PATH}/config.toml | \
  sed "s|__CHAIN_ID__|${CHAIN_ID}|g" | \
  sed "s|__ORACLE_MNEMONIC__|${ORACLE_MNEMONIC}|g" | \
  sed "s|__PANACEA_HOST__|${PANANCE_HOST}|g" \
  > ~/.${CONFIG_PATH}/config.toml