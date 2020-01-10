FROM smarman/ubuntu18.04:latest

# This is the release of Vault to pull in.
ARG CONSUL_VERSION=1.6.2
ARG VAULT_VERSION=1.3.1
ENV VAULT_ADDR=http://127.0.0.1:8200

# Create a vault user and group first so the IDs get set the same way,
# even as the rest of this may change over time.
RUN addgroup vault && \
    useradd -m -g vault vault && \
    addgroup consul && \
    useradd -m -g consul consul

# Set up certificates, our base tools, and Vault.
RUN set -eux; \
    apt-get -y update && \
    apt-get install -y openssl \
    dumb-init \
    net-tools \
    libcap2-bin && \
    mkdir -p /tmp/build && \
    cd /tmp/build && \
    wget https://releases.hashicorp.com/vault/${VAULT_VERSION}/vault_${VAULT_VERSION}_linux_amd64.zip && \
    unzip -d /bin vault_${VAULT_VERSION}_linux_amd64.zip && \
    wget https://releases.hashicorp.com/consul/${CONSUL_VERSION}/consul_${CONSUL_VERSION}_linux_amd64.zip && \
    unzip -d /bin consul_${CONSUL_VERSION}_linux_amd64.zip && \
    cd /tmp && \
    rm -rf /tmp/build && \
    rm -rf /root/.gnupg

COPY ./config/vault-config.json /vault/config/vault-config.json
COPY ./config/consul-config.json /var/consul.json

# /vault/logs is made available to use as a location to store audit logs, if
# desired; /vault/file is made available to use as a location with the file
# storage backend, if desired; the server will be started with /vault/config as
# the configuration directory so you can add additional config files in that
# location.

RUN mkdir -p /vault/logs && \
    mkdir -p /vault/file && \
    mkdir -p /vault/config && \
    chown -R vault:vault /vault && \
    mkdir -p /consul/data && \
    chown -R consul:consul /consul && \
    setcap cap_ipc_lock=+ep $(readlink -f $(which vault))  && \
    chown consul:consul /var/consul.json && \
    sed -i "s/BIND_ADDR/$(ifconfig eth0 | grep inet | awk '{print $2}')/g" "/var/consul.json"
    
    

# Expose the logs directory as a volume since there's potentially long-running
# state in there
VOLUME /vault/logs

# Expose the file directory as a volume since there's potentially long-running
# state in there
VOLUME /vault/file

# 8200/tcp is the primary interface that applications use to interact with
# Vault.
EXPOSE 8200

# The entry point script uses dumb-init as the top-level process to reap any
# zombie processes created by Vault sub-processes.
#
# For production derivatives of this container, you shoud add the IPC_LOCK
# capability so that Vault can mlock memory.
#COPY docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
#ENTRYPOINT ["docker-entrypoint.sh"]
USER vault
ENTRYPOINT ["/usr/bin/dumb-init", "--"]

# By default you'll get a single-node development server that stores everything
# in RAM and bootstraps itself. Don't use this configuration for production.

# !!!!! remove disable_mlock: false in config file before going prod
CMD ["vault", "server", "-config", "/vault/config/vault-config.json"]
#USER consul
#CMD ["consul", "agent", "-config-file", "/var/consul.json"]
#CMD ["/bin/bash"]

