#!/bin/bash

apt-get update && apt-get install -y unzip

# Fetch Consul
cd /tmp
wget https://releases.hashicorp.com/consul/1.6.0/consul_1.6.0_linux_amd64.zip -O ./consul.zip
unzip ./consul.zip
mv ./consul /usr/local/bin

# Fetch Envoy
wget https://github.com/nicholasjackson/cloud-pong/releases/download/v0.1.1/envoy -O /usr/local/bin/envoy
chmod +x /usr/local/bin/envoy

# Create the consul config
mkdir -p /etc/consul
cat << EOF > /etc/consul/config.hcl
data_dir = "/tmp/"
log_level = "DEBUG"

datacenter = "vms"

bind_addr = "0.0.0.0"
client_addr = "0.0.0.0"

ports {
  grpc = 8502
}

connect {
  enabled = true
}

enable_central_service_config = true

advertise_addr = "${advertise_addr}"
retry_join = ["${consul_cluster_addr}"]
EOF

# Setup system D
cat << EOF > /etc/systemd/system/consul.service
[Unit]
Description=Consul Server
After=syslog.target network.target

[Service]
ExecStart=/usr/local/bin/consul agent -config-file=/etc/consul/config.hcl
ExecStop=/bin/sleep 5
Restart=always

[Install]
WantedBy=multi-user.target
EOF

chmod 644 /etc/systemd/system/consul.service

# Setup system D
cat << EOF > /etc/systemd/system/consul-gateway.service
[Unit]
Description=Consul Gateway
After=syslog.target network.target

[Service]
ExecStart=/usr/local/bin/consul connect envoy -mesh-gateway -register -wan-address ${gateway_addr}:443 -- -l debug
ExecStop=/bin/sleep 5
Restart=always

[Install]
WantedBy=multi-user.target
EOF

chmod 644 /etc/systemd/system/consul-gateway.service

systemctl daemon-reload
systemctl start consul.service
systemctl start consul-gateway.service
