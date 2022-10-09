FROM quay.io/keycloak/keycloak:19.0 as builder

ENV KC_METRICS_ENABLED=true
ENV KC_FEATURES=token-exchange,scripts,authorization,step-up-authentication,client-secret-rotation,client-policies,step-up-authentication,web-authn,impersonation,admin2,admin-fine-grained-authz
RUN mkdir -p /opt/keycloak/data/import
ADD realm-export.json /opt/keycloak/data/import/realm-export.json
RUN /opt/keycloak/bin/kc.sh import --file /opt/keycloak/data/import/realm-export.json
ENV KC_DB=postgres
RUN /opt/keycloak/bin/kc.sh build --db=postgres

FROM quay.io/keycloak/keycloak:19.0
COPY --from=builder /opt/keycloak/lib/quarkus/ /opt/keycloak/lib/quarkus/
WORKDIR /opt/keycloak
# for demonstration purposes only, please make sure to use proper certificates in production instead
RUN keytool -genkeypair -storepass password -storetype PKCS12 -keyalg RSA -keysize 2048 -dname "CN=server" -alias server -ext "SAN:c=DNS:localhost,IP:127.0.0.1,DNS:keycloak,DNS:kubernetes.lehi-k8s-vanilla.calvarado04.com" -keystore conf/server.keystore
RUN mkdir -p /opt/keycloak/data/import
ADD realm-export.json /opt/keycloak/data/import/realm-export.json
ENV KEYCLOAK_ADMIN=admin
ENV KEYCLOAK_ADMIN_PASSWORD=change_me
ENTRYPOINT ["/opt/keycloak/bin/kc.sh", "start", "--import-realm", "--db=postgres", "--http-relative-path=/auth", "--hostname-strict-https=false", "--hostname-strict=false"]