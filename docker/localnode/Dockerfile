FROM ubuntu:latest
ENV HOME="/root" PATH="/root/go/bin:/bluelink-lab/blk-chain/integration_test/upgrade_module/scripts/:$PATH"
RUN apt-get update && apt-get install -y make build-essential git jq python3 curl vim uuid-runtime
RUN curl -L https://go.dev/dl/go1.21.12.linux-amd64.tar.gz | tar xvzf - -C /usr/local/
RUN curl -L https://foundry.paradigm.xyz | bash
RUN curl -sL https://deb.nodesource.com/setup_16.x | bash
RUN /root/.foundry/bin/foundryup
RUN mkdir -p /root/go/pkg/mod && \
    mkdir -p /root/.cache && \
    chmod -R a+rwX /root
SHELL ["/bin/bash", "-c"]

WORKDIR /bluelink-lab/blk-chain

EXPOSE 26656 26657 26658 9090 9091 7171

CMD ["/usr/bin/deploy.sh"]

STOPSIGNAL SIGTERM

COPY scripts/deploy.sh /usr/bin/deploy.sh
COPY scripts/step0_build.sh /usr/bin/build.sh
COPY scripts/step1_configure_init.sh /usr/bin/configure_init.sh
COPY scripts/step2_genesis.sh /usr/bin/genesis.sh
COPY scripts/step3_add_validator_to_genesis.sh /usr/bin/add_validator_to_gensis.sh
COPY scripts/step4_config_override.sh /usr/bin/config_override.sh
COPY scripts/step5_start_she.sh /usr/bin/start_she.sh
COPY scripts/step6_start_price_feeder.sh /usr/bin/start_price_feeder.sh

