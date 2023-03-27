FROM ubuntu:focal-20220801 AS ego-base

ARG PSW_VERSION=2.17.100.3-focal1
ARG DCAP_VERSION=1.14.100.3-focal1
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates gnupg libcurl4 wget jq \
  && wget -qO- https://download.01.org/intel-sgx/sgx_repo/ubuntu/intel-sgx-deb.key | apt-key add \
  && echo 'deb [arch=amd64] https://download.01.org/intel-sgx/sgx_repo/ubuntu focal main' >> /etc/apt/sources.list \
  && wget -qO- https://packages.microsoft.com/keys/microsoft.asc | apt-key add \
  && echo 'deb [arch=amd64] https://packages.microsoft.com/ubuntu/20.04/prod focal main' >> /etc/apt/sources.list \
  && apt-get update && apt-get install -y --no-install-recommends \
  libsgx-ae-id-enclave=$DCAP_VERSION \
  libsgx-ae-pce=$PSW_VERSION \
  libsgx-ae-qe3=$DCAP_VERSION \
  libsgx-dcap-ql=$DCAP_VERSION \
  libsgx-enclave-common=$PSW_VERSION \
  libsgx-launch=$PSW_VERSION \
  libsgx-pce-logic=$DCAP_VERSION \
  libsgx-qe3-logic=$DCAP_VERSION \
  libsgx-urts=$PSW_VERSION \
  && apt-get install -y az-dcap-client

# Install EGo
# https://docs.edgeless.systems/ego/#/getting-started/install?id=install-the-deb-package
RUN wget https://github.com/edgelesssys/ego/releases/download/v1.0.1/ego_1.0.1_amd64.deb \
  && apt install ./ego_1.0.1_amd64.deb

####################################################

FROM ego-base AS build

# Install prerequisites
RUN apt-get update && apt-get install -y --no-install-recommends git build-essential

RUN git clone "https://github.com/medibloc/panacea-oracle.git"

# Build oracled
WORKDIR /go/src/github.com/medibloc/panacea-oracle

RUN make build

RUN ego sign ./scripts/enclave-prod.json

COPY ./init_oracle.sh .

RUN chmod +x init_oracle.sh

####################################################

FROM ego-base

COPY --from=build /go/src/github.com/medibloc/panacea-oracle/build/oracled /usr/bin/oracled
RUN chmod +x /usr/bin/oracled

EXPOSE 8080 9090

CMD ["sh", "./init_oracle.sh"]
