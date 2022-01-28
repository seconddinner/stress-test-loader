resource "aws_route53_record" "california" {
  zone_id = "${var.route53_zone_id}"
  name    = "${var.domainname}"
  type    = "A"

  alias {
    name                   = "${aws_s3_bucket.b.website_domain}"
    zone_id                = "${aws_s3_bucket.b.hosted_zone_id}"
    evaluate_target_health = true
  }

  geolocation_routing_policy {
    country     = "US"
    subdivision = "CA"
  }

  set_identifier = "California"
}
