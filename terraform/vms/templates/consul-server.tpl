#!/bin/bash

apt-get update && apt-get install -y unzip

# Fetch Consul
cd /tmp
wget https://releases.hashicorp.com/consul/1.6.0/consul_1.6.0_linux_amd64.zip -O ./consul.zip
unzip ./consul.zip
mv ./consul /usr/local/bin

# Create the consul config
mkdir -p /etc/consul
cat << EOF > /etc/consul/config.hcl
data_dir = "/tmp/"
log_level = "DEBUG"

datacenter = "vms"
primary_datacenter = "aks"

server = true

bootstrap_expect = 1
ui = true

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
advertise_addr_wan = "${wan_addr}"
retry_join_wan = ["${primary_cluster_addr}"]
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

systemctl daemon-reload
systemctl start consul.service

## Wait for consul to start and then register a default
## deny all intention

until consul members; do
  echo "Waiting for Consul to start"
  sleep 1
done

consul intention create -deny '*' '*'
