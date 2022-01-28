output "stress_test_loader_url" {
  description = "Squid URL"

  value = join("", ["http://", module.autoscale-ap-southeast-1.dns, ":", var.stress_test_loader_port])
}

output "PNS_version" {
  value = local.PNS_version
}

