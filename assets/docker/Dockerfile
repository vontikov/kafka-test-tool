FROM scratch

ARG VERSION

LABEL \
  org.opencontainers.image.title="Kafka test tool" \
  org.opencontainers.image.description="\
Accesses Kafka via HTTP API" \
  org.opencontainers.image.version="${VERSION}" \
  org.opencontainers.image.authors="Vladimir Ontikov <vontikov@gmail.com>"

COPY kafka_test_tool /usr/local/bin/
CMD ["kafka_test_tool"]
