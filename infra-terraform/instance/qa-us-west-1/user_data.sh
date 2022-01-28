#!/usr/bin/env bash

# cat <<EOF > /usr/local/stress_test_loader-node-selfcheck/config.json
# {
#     "VRList": [
#          {
#             "url": "http://qa-p-target.zone2.seconddinnertech.net",
#             "timeout": 2,
#             "urlParameters": "index.html"
#         }
#     ],
#     "waitSeconds": 3000,
#     "listenPort": 9005,
#     "initialWaitForProxyStart": 300,
#     "shutdownThreshHold": 2
# }
# EOF

cat <<EOF > /usr/local/stress_test_loader-node-selfcheck/config.json
{
    "VRList": [],
    "waitSeconds": 3000,
    "listenPort": 9005,
    "initialWaitForProxyStart": 300,
    "shutdownThreshHold": 2,
    "shutdownFlag": false
}
EOF

cat <<EOF > /etc/default/prometheus-stress_test_loader-exporter
ARGS="-stress_test_loader-port ${stress_test_loader_port}"

# prometheus-stress_test_loader-exporter supports the following options:
#  -listen string
#        Address and Port to bind exporter, in host:port format (default ":9301")
#  -metrics-path string
#        Metrics path to expose prometheus metrics (default "/metrics")
#  -stress_test_loader-hostname string
#        Squid hostname (default "localhost")
#  -stress_test_loader-login string
#        Login to stress_test_loader service
#  -stress_test_loader-password string
#        Password to stress_test_loader service
#  -stress_test_loader-port int
#        Squid port to read metrics (default 3128)
EOF

cat <<EOF > /usr/local/stress_test_loader-node-selfcheck/reportbadserver.sh
#!/usr/bin/env bash

PublicIP=$(curl -s http://169.254.169.254/latest/meta-data/public-ipv4)
PrivateIP=$(curl -s http://169.254.169.254/latest/meta-data/local-ipv4)
Zone=$(curl -s http://169.254.169.254/latest/meta-data/placement/availability-zone)
curl -X POST "https://register.zone2.seconddinnertech.net/v1/server/1xCrmez3Xy34qKg4h99RUgqVbQR6sUPVmYSTc0TE1rJYLmgKKpwnOshC1Ejvc2YM/badserver" -H  "accept: application/json" -H  "content-type: application/json" -d "{  \"Cloud\": \"aws\",   \"PrivateIP\": \"\$PrivateIP\",  \"PublicIP\": \"\$PublicIP\",  \"Type\": \"${environment}\",  \"Zone\": \"\$Zone\"}"
EOF

chmod 755 /usr/local/stress_test_loader-node-selfcheck/reportbadserver.sh

cat <<EOF > /root/register.sh
#!/usr/bin/env bash

PublicIP=$(curl -s http://169.254.169.254/latest/meta-data/public-ipv4)
PrivateIP=$(curl -s http://169.254.169.254/latest/meta-data/local-ipv4)
Zone=$(curl -s http://169.254.169.254/latest/meta-data/placement/availability-zone)
curl -X POST "https://register.zone2.seconddinnertech.net/v1/server/1xCrmez3Xy34qKg4h99RUgqVbQR6sUPVmYSTc0TE1rJYLmgKKpwnOshC1Ejvc2YM/" -H  "accept: application/json" -H  "content-type: application/json" -d "{  \"Cloud\": \"aws\",   \"PrivateIP\": \"\$PrivateIP\",  \"PublicIP\": \"\$PublicIP\",  \"Type\": \"${environment}\",  \"Zone\": \"\$Zone\"}"
EOF

source /root/register.sh


cat <<EOF > /etc/stress_test_loader/stress_test_loader.conf
# jc request rules

forwarded_for delete
via off

# add reply header
reply_header_add Proxy-Node-Public-IP ___SED___REPLACEME___

# block domains

acl toblock dstdomain 
http_access deny toblock

# Adapt to list your (internal) IP networks from where browsing
# should be allowed
acl all src all


%{ for addr in stress_test_loader_allowed_cidr ~}
acl localnet src ${addr}
%{ endfor ~}

acl SSL_ports port 443
acl Safe_ports port 80          # http
acl Safe_ports port 443         # https
acl Safe_ports port 1025-65535  # unregistered ports
acl CONNECT method CONNECT

#
# Recommended minimum Access Permission configuration:
#
# Deny requests to certain unsafe ports
http_access deny !Safe_ports

# Deny CONNECT to other than secure SSL ports
http_access deny CONNECT !SSL_ports

# Only allow cachemgr access from localhost
#http_access allow all
cache deny all
http_access allow localhost manager
http_access deny manager

# Allowing access only to AWS sites.
acl allowed_http_sites dstdomain .amazonaws.com
http_access allow allowed_http_sites

# Example rule allowing access from your local networks.
# Adapt localnet in the ACL section to list your (internal) IP networks
# from where browsing should be allowed
#http_access allow all
http_access allow localnet
http_access allow localhost

# And finally deny all other access to this stress_test_loader
http_access deny all

# Listen port
http_port ${stress_test_loader_port}
EOF

export PublicIP=$(curl -s http://169.254.169.254/latest/meta-data/public-ipv4)
sed -i "s/___SED___REPLACEME___/$PublicIP/g" /etc/stress_test_loader/stress_test_loader.conf

cat <<EOF > /etc/systemd/system/stress_test_loader_export.service
[Unit]
Description=Squid exporter

[Service]
ExecStart=/usr/sbin/stress_test_loader-exporter -stress_test_loader-port ${stress_test_loader_port} -listen ":9301"

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl restart stress_test_loader
systemctl restart prometheus-node-exporter
systemctl restart prometheus-stress_test_loader-exporter
systemctl restart stress_test_loader-node-selfcheck.service