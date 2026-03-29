FROM quay.io/keycloak/keycloak:26.5 AS builder

ENV KC_METRICS_ENABLED=true
ENV KC_DB=postgres
ENV KC_HTTP_RELATIVE_PATH=/auth
ENV KC_FEATURES=token-exchange,scripts,authorization,step-up-authentication,client-secret-rotation,client-policies,web-authn,impersonation,admin,admin-fine-grained-authz

RUN /opt/keycloak/bin/kc.sh build

FROM quay.io/keycloak/keycloak:26.5

COPY --from=builder /opt/keycloak/lib/quarkus/ /opt/keycloak/lib/quarkus/

WORKDIR /opt/keycloak

RUN keytool -genkeypair \
  -storepass password \
  -storetype PKCS12 \
  -keyalg RSA \
  -keysize 2048 \
  -dname "CN=server" \
  -alias server \
  -ext "SAN:c=DNS:localhost,IP:127.0.0.1,DNS:keycloak,DNS:kubernetes.lehi-k8s-vanilla.calvarado04.com" \
  -keystore conf/server.keystore

RUN mkdir -p /opt/keycloak/data/import
COPY realm-export.json /opt/keycloak/data/import/realm-export.json

ENV KC_DB=postgres
ENV KC_HTTP_RELATIVE_PATH=/auth
ENV KC_BOOTSTRAP_ADMIN_USERNAME=admin
ENV KC_BOOTSTRAP_ADMIN_PASSWORD=change_me

ENTRYPOINT ["/opt/keycloak/bin/kc.sh", "start", "--optimized", "--import-realm", "--hostname-strict=false", "--proxy-headers=xforwarded"]