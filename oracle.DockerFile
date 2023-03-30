FROM ghcr.io/medibloc/panacea-oracle:main AS oracle

# Install prerequisites
RUN apt-get update && apt-get install -y --no-install-recommends git build-essential

COPY ./init_oracle.sh .
COPY ./config.toml .

RUN chmod +x ./init_oracle.sh

EXPOSE 8080 9090

CMD ["sh", "./init_oracle.sh"]
