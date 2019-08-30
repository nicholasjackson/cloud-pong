resource "kubernetes_config_map" "central_config" {
  metadata {
    name = "central-config"
  }

  data = {
    "pong-aks-defaults.hcl" = file("${path.module}/consul_config/pong-aks-defaults.hcl")
    "pong-vms-defaults.hcl" = file("${path.module}/consul_config/pong-vms-defaults.hcl")
    "pong-aks-resolver.hcl" = file("${path.module}/consul_config/pong-aks-resolver.hcl")
    "pong-vms-resolver.hcl" = file("${path.module}/consul_config/pong-vms-resolver.hcl")
  }
}

resource "kubernetes_job" "central_config" {
  depends_on = [helm_release.consul]

  metadata {
    name = "central-config"
  }

  spec {
    template {
      metadata {
        labels = {
          version = "v0.0.1"
        }
      }
      spec{
        volume {
          name = kubernetes_config_map.central_config.metadata[0].name
        
          config_map {
            name = kubernetes_config_map.central_config.metadata[0].name
          }
        }

        container {
          image = "nicholasjackson/consul-envoy:v1.6.0-v0.10.0"
      		name = "central-config"

          env {
            name  = "CONSUL_HTTP_ADDR"
            value = "consul-consul-server:8500"
          }

          env {
            name  = "CONSUL_GRPC_ADDR"
            value = "consul-consul-server:8502"
          }

          env {
            name  = "CENTRAL_CONFIG_DIR"
            value = "/config"
          }
          
      		volume_mount {
          	read_only = true  
            mount_path = "/config"
            name = kubernetes_config_map.central_config.metadata[0].name
      		}
        }
      }
    }
  }
}
