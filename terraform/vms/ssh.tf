resource "tls_private_key" "pong" {
  algorithm   = "RSA"
  rsa_bits = "4096"
}
