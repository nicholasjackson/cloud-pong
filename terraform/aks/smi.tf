resource "kubernetes_cluster_role" "smi" {
  metadata {
    name = "smi-controller"
  
    labels = {
      "rbac.authorization.k8s.io/aggregate-to-view": "true" 
      "rbac.authorization.k8s.io/aggregate-to-cluster-reader": "true" 
    }
  }


  rule {
    api_groups = ["specs.smi-spec.io", "access.smi-spec.io"] 
    resources = ["tcproutes","htttproutegroups","traffictargets","events"] 
    verbs = ["get", "list", "watch", "create", "update", "patch", "delete"]
  }
}

resource "kubernetes_service_account" "smi" {
  metadata {
    name      = "smi-controller"
    namespace = "default"
  }
}

resource "kubernetes_cluster_role_binding" "smi" {
  metadata {
    name = "smi-controller"
  }

  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "ClusterRole"
    name      = "smi-controller"
  }

  subject {
    kind      = "ServiceAccount"
    name      = kubernetes_service_account.smi.metadata[0].name
    namespace = "default"
  }
}

resource "kubernetes_deployment" "smi" {
  depends_on = [helm_release.consul]

  metadata {
    name = "consul-smi-controller-deployment"
    labels = {
      app = "consul-smi-controller"
    }
  }
  
  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "consul-smi-controller"
      }
    }
    
    template {
      metadata {
        labels = {
          app     = "consul-smi-controller"
        }
      }

      spec {
        service_account_name = kubernetes_service_account.smi.metadata[0].name
        automount_service_account_token = true

        container {
          image = "hashicorp/consul-smi-controller:v0.0.0-alpha.1"
          name  = "consul-smi-controller"
          command = ["/app/smi-controller"]
          args = ["--consul-http-addr=http://$(HOST_IP):8500"]

          env {
            name  = "HOST_IP"
            value_from { 
              field_ref {
                field_path = "status.hostIP"
              }
            }
          }
        }
      }
    }
  }
}
