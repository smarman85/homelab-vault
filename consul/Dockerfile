FROM smarman/ubuntu18.04:latest

MAINTAINER Consul Team <consul@hashicorp.com>

ENV CONSUL_VERSION=1.6.2

ENV HASHICORP_RELEASES=https://releases.hashicorp.com

RUN addgroup consul && \
    useradd -m -g consul consul

RUN set -eux && \
    apt-get update && \
    apt-get install -y curl \
    dumb-init \
    net-tools \
    openssl \
    jq && \
    wget ${HASHICORP_RELEASES}/consul/${CONSUL_VERSION}/consul_${CONSUL_VERSION}_linux_amd64.zip && \
    unzip -d /bin consul_${CONSUL_VERSION}_linux_amd64.zip && \
    cd /tmp && \
    rm -rf /tmp/build && \
# tiny smoke test to ensure the binary we downloaded runs
    consul version

COPY ./config/consul-config.json /consul/config/consul-config.json

RUN mkdir -p /consul/data && \
    mkdir -p /consul/config && \
    chown -R consul:consul /consul && \
    sed -i "s/ADVERTISE_ADDR/$(ifconfig eth0 | grep inet | awk '{print $2}')/g" "/consul/config/consul-config.json"

# set up nsswitch.conf for Go's "netgo" implementation which is used by Consul,
# otherwise DNS supercedes the container's hosts file, which we don't want.
RUN test -e /etc/nsswitch.conf || echo 'hosts: files dns' > /etc/nsswitch.conf

# Expose the consul data directory as a volume since there's mutable state in there.
VOLUME /consul/data

# Server RPC is used for communication between Consul clients and servers for internal
# request forwarding.
EXPOSE 8300

# Serf LAN and WAN (WAN is used only by Consul servers) are used for gossip between
# Consul agents. LAN is within the datacenter and WAN is between just the Consul
# servers in all datacenters.
EXPOSE 8301 8301/udp 8302 8302/udp

# HTTP and DNS (both TCP and UDP) are the primary interfaces that applications
# use to interact with Consul.
EXPOSE 8500 8600 8600/udp

# Consul doesn't need root privileges so we run it as the consul user from the
# entry point script. The entry point script also uses dumb-init as the top-level
# process to reap any zombie processes created by Consul sub-processes.
ENTRYPOINT ["/usr/bin/dumb-init", "--"]

# By default you'll get an insecure single-node development server that stores
# everything in RAM, exposes a web UI and HTTP endpoints, and bootstraps itself.
# Don't use this configuration for production.
CMD ["consul", "agent", "-config-file", "/consul/config/consul-config.json"]
