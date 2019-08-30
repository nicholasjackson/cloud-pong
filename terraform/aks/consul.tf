# Create the public loadbalancer for Consul Gateways
resource "kubernetes_service" "gateways" {
  metadata {
    name = "gateways"
  }
  spec {
    selector = {
      app       = "consul"
      component = "mesh-gateway"
    }

    port {
      port        = 443
      target_port = 443
    }

    type = "LoadBalancer"
  }
}

# Create public loadbalancer for Consul
# WARNING: this is not prodution config
# Consul should not be exposed publically
resource "kubernetes_service" "consul" {
  metadata {
    name = "consul-lb"
  }

  spec {
    selector = {
      app       = "consul"
      component = "server"
      release   = "consul"
    }

    port {
      name        = "api-ui"
      port        = 80
      target_port = 8500
    }

    port {
      name        = "api-api"
      port        = 8500
      target_port = 8500
    }

    port {
      name        = "serf-wan-tcp"
      port        = 8302
      target_port = 8302
      protocol    = "TCP"
    }

    port {
      name        = "consul-wan-rpc"
      port        = 8300
      target_port = 8300
      protocol    = "TCP"
    }

    type = "LoadBalancer"
  }
}

provider "helm" {
  kubernetes {
    host                   = azurerm_kubernetes_cluster.aks.kube_config.0.host
    client_certificate     = base64decode(azurerm_kubernetes_cluster.aks.kube_config.0.client_certificate)
    client_key             = base64decode(azurerm_kubernetes_cluster.aks.kube_config.0.client_key)
    cluster_ca_certificate = base64decode(azurerm_kubernetes_cluster.aks.kube_config.0.cluster_ca_certificate)
  }

  service_account = "tiller"
  tiller_image    = "gcr.io/kubernetes-helm/tiller:v2.13.1"
}

resource "helm_release" "consul" {
  depends_on = [kubernetes_cluster_role_binding.tiller]

  name      = "consul"
  chart     = "${path.module}/helm/consul"
  namespace = "default"

  set {
    name  = "global.image"
    value = "consul:1.6.0"
  }

  set {
    name  = "global.datacenter"
    value = "aks"
  }

  set {
    name  = "server.replicas"
    value = 1
  }

  set {
    name  = "server.bootstrapExpect"
    value = 1
  }

  set {
    name  = "client.grpc"
    value = true
  }

  set {
    name  = "connectInject.enabled"
    value = true
  }

  set {
    name  = "centralConfig.enabled"
    value = true
  }

  set_string {
    name  = "server.extraConfig"
    value = "\"{\\\"advertise_addr_wan\\\": \\\"${kubernetes_service.consul.load_balancer_ingress.0.ip}\\\"\\, \\\"primary_datacenter\\\": \\\"aks\\\"}\""
  }

  set {
    name  = "connectInject.centralConfig.enabled"
    value = true
  }

  set {
    name  = "connectInject.centralConfig.proxyDefaults"
    value = <<EOF
      {
        "envoy_prometheus_bind_addr": "0.0.0.0:9102"
      }
    EOF
  }

  set {
    name  = "connectInject.imageEnvoy"
    value = "envoyproxy/envoy:v1.10.0"
  }

  set {
    name  = "meshGateway.enabled"
    value = true
  }
  
  set {
    name  = "meshGateway.enableHealthChecks"
    value = false
  }

  set {
    name  = "meshGateway.mode"
    value = "local"
  }

  set {
    name  = "meshGateway.wanAddress.useNodeIP"
    value = false
  }

  set {
    name  = "meshGateway.wanAddress.host"
    value = "${kubernetes_service.gateways.load_balancer_ingress.0.ip}"
  }

  set {
    name  = "meshGateway.imageEnvoy"
    value = "envoyproxy/envoy:v1.10.0"
  }
}
