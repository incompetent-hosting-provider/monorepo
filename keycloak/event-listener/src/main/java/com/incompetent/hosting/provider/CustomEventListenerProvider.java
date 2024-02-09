package com.incompetent.hosting.provider;

import java.util.logging.Logger;

import org.keycloak.events.Event;
import org.keycloak.events.EventListenerProvider;
import org.keycloak.events.EventType;
import org.keycloak.events.admin.AdminEvent;
import org.keycloak.models.KeycloakSession;

public class CustomEventListenerProvider implements EventListenerProvider {

    private static final Logger log = Logger.getLogger(CustomEventListenerProvider.class.getName());

    public CustomEventListenerProvider(KeycloakSession session) {
    }

    @Override
    public void onEvent(Event event) {

        if (EventType.REGISTER.equals(event.getType())) {
            log.info("--------\nReceived register event\n-------------");
            return;
        }

        if (EventType.DELETE_ACCOUNT.equals(event.getType())){
            log.info("------\nReceived account delete event\n-----");
            return;
        }

    }

    @Override
    public void onEvent(AdminEvent adminEvent, boolean b) {

    }

    @Override
    public void close() {

    }
}