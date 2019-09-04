#!/bin/bash

apt-get update && apt-get install -y unzip

cd /tmp

# Fetch Pong
wget https://github.com/nicholasjackson/cloud-pong/releases/download/v0.3.0/pong-api -O /usr/local/bin/pong-api
chmod +x /usr/local/bin/pong-api

# Fetch Envoy
wget https://github.com/nicholasjackson/cloud-pong/releases/download/v0.3.0/envoy -O /usr/local/bin/envoy
chmod +x /usr/local/bin/envoy

# Fetch Consul
wget https://releases.hashicorp.com/consul/1.6.0/consul_1.6.0_linux_amd64.zip -O ./consul.zip
unzip ./consul.zip
mv ./consul /usr/local/bin

# Create the consul config
mkdir -p /etc/consul/config

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

# Create config and register service
cat << EOF > /etc/consul/config/pong.json
{
  "service": {
    "name": "pong-vms",
    "id":"pong-vms",
    "port": 6000,
    "connect": { 
      "sidecar_service": {
        "port": 20000,
        "proxy": {
          "upstreams": [
            {
              "destination_name": "pong-aks",
              "local_bind_address": "127.0.0.1",
              "local_bind_port": 6001
            }
          ]
        }
      }
    }  
  }
}
EOF

# Setup systemd Consul Agent
cat << EOF > /etc/systemd/system/consul.service
[Unit]
Description=Consul Server
After=syslog.target network.target

[Service]
ExecStart=/usr/local/bin/consul agent -config-file=/etc/consul/config.hcl -config-dir=/etc/consul/config
ExecStop=/bin/sleep 5
Restart=always

[Install]
WantedBy=multi-user.target
EOF

chmod 644 /etc/systemd/system/consul.service

# Setup systemd Envoy Sidecar
cat << EOF > /etc/systemd/system/consul-envoy.service
[Unit]
Description=Consul Envoy
After=syslog.target network.target

[Service]
ExecStart=/usr/local/bin/consul connect envoy -sidecar-for pong-vms
ExecStop=/bin/sleep 5
Restart=always

[Install]
WantedBy=multi-user.target
EOF

chmod 644 /etc/systemd/system/consul-envoy.service

# Setup systemd Pong
cat << EOF > /etc/systemd/system/pong.service
[Unit]
Description=Pong
After=syslog.target network.target

[Service]
Environment=PLAYER=2
ExecStart=/usr/local/bin/pong-api
ExecStop=/bin/sleep 5
Restart=always

[Install]
WantedBy=multi-user.target
EOF

chmod 644 /etc/systemd/system/pong.service

systemctl daemon-reload
systemctl start consul.service
systemctl start consul-envoy.service
systemctl start pong.service
