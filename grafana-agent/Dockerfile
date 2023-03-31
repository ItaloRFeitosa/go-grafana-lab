FROM grafana/agent:v0.30.2

COPY ./grafana-agent.yml /etc/agent/agent.yaml

ENTRYPOINT ["/bin/agent", "-config.file=/etc/agent/agent.yaml", "-metrics.wal-directory=/tmp/agent/wal", "-config.expand-env"]