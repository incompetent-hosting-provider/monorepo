FROM maven:3.8.7-openjdk-18-slim as BUILD

WORKDIR /app/build

COPY event-listener .

RUN mvn clean compile assembly:single

FROM quay.io/keycloak/keycloak:22.0.5 AS SETUP

COPY --from=BUILD /app/build/target/custom-event-listener-jar-with-dependencies.jar /opt/keycloak/providers/custom-event-listener.jar

WORKDIR /opt/keycloak

RUN /opt/keycloak/bin/kc.sh build

FROM quay.io/keycloak/keycloak:22.0.5 AS MAIN

COPY --from=SETUP /opt/keycloak/ /opt/keycloak/
