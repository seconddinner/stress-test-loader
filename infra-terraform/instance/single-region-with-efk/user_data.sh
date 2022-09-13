#!/usr/bin/env bash

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

cat <<EOF > /etc/rsyslog.d/99-stresstest.conf
if \$programname == 'stress-test-loader-linux' then @127.0.0.1:10514
EOF

# setup filebeat
wget -qO - https://artifacts.elastic.co/GPG-KEY-elasticsearch | sudo apt-key add -
sudo apt-get install apt-transport-https
echo "deb https://artifacts.elastic.co/packages/oss-7.x/apt stable main" | sudo tee -a /etc/apt/sources.list.d/elastic-7.x.list
sudo apt-get update && \
sudo apt-get install -y filebeat='7.10.2'

cat <<EOF > /etc/filebeat/filebeat.yml
scan_frequency: 10s
filebeat.shutdown_timeout: 2s
filebeat.registry.flush: 0s
logging.metrics.enabled: true
logging.level: info
queue.mem:
  events: 65535
  flush.min_events: 2048
  flush.timeout: 3s
setup.ilm:
  enabled: false
  ilm.enabled: false
setup.template:
  enabled: true
  name: snap-nvstresstest-
  pattern: snap-nvstresstest-*
  settings:
    index.number_of_shards: 16
filebeat.inputs:
- type: log
  enabled: true
  tail_files: true
  close_inactive: 5m
  ignore_older: 10m
  close_timeout: 1h
  paths:
    - /var/log/stresstest-client/stresstest-client*.log
  fields:
    index: snap-nvstresstest
  processors:
  - decode_json_fields:
      fields: ["message"]
      target: "stressclient"
  - timestamp:
      field: stressclient.EndTime
      target_field: stressclient.EndTimestamp
      layouts:
        - UNIX_MS
      test:
        - '1661403661761'
output.elasticsearch:
  indices:
  - index: snap-nvstresstest-%%{+yyyy.MM.dd}
    when.contains:
      fields:
        index: snap-nvstresstest
  hosts:
  - https://${es_point}:443
  username: stresstestadmin
  password: "${masterpassword}"
  worker: 2
  bulk_max_size: 4096
EOF


systemctl daemon-reload
systemctl restart stress_test_loader
systemctl restart prometheus-node-exporter
systemctl restart prometheus-stress_test_loader-exporter
systemctl restart stress_test_loader-node-selfcheck.service
systemctl restart filebeat
systemctl restart rsyslog