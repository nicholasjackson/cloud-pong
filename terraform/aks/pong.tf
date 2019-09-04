resource "kubernetes_deployment" "pong" {
  depends_on = [helm_release.consul]

  metadata {
    name = "pong-api"
    labels = {
      app = "pong-api"
    }
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "pong-api"
      }
    }

    template {
      metadata {
        labels = {
          app     = "pong-api"
          version = "v0.2.0"
        }

        annotations = {
          "consul.hashicorp.com/connect-inject"            = "true"
          "consul.hashicorp.com/connect-service-upstreams" = "pong-vms:6001"
          "consul.hashicorp.com/connect-service-name"      = "pong-aks"
        }
      }

      spec {
        container {
          image = "nicholasjackson/cloud-pong-api:go-v0.3.0"
          name  = "pong-aks"

          port {
            name           = "http"
            container_port = 6000
          }

          resources {
            limits {
              cpu    = "0.5"
              memory = "512Mi"
            }
            requests {
              cpu    = "0.1"
              memory = "50Mi"
            }
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "pong_service" {
  metadata {
    name = "pong-lb"
  }

  spec {
    selector = {
      app = "pong-api"
    }

    port {
      name        = "http"
      port        = 6000
      target_port = 6000
    }

    type = "LoadBalancer"
  }
}
