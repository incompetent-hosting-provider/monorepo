package com.incompetent.hosting.provider;

import java.io.IOException;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.util.logging.Logger;

import org.keycloak.events.Event;
import org.keycloak.events.EventListenerProvider;
import org.keycloak.events.EventType;
import org.keycloak.events.admin.AdminEvent;
import org.keycloak.events.admin.OperationType;
import org.keycloak.events.admin.ResourceType;
import org.keycloak.models.KeycloakSession;
import org.json.JSONObject;

public class CustomEventListenerProvider implements EventListenerProvider {

    private static final Logger log = Logger.getLogger(CustomEventListenerProvider.class.getName());

    private final String BackendHost = System.getenv("ICHP_BACKEND_HOST");

    private final HttpClient httpClient = HttpClient.newHttpClient();

    public CustomEventListenerProvider(KeycloakSession session) {
        log.info("provider...");
    }

    @Override
    public void onEvent(Event event) {
        if (EventType.REGISTER.equals(event.getType()) || EventType.DELETE_ACCOUNT.equals(event.getType())) {
            callBackendWebhook(buildJSONBody(event.getType(), event.getUserId()));
        }
    }

    private String buildJSONBody(EventType eventType, String userId){
        return new JSONObject().put("keycloakEvent", eventType.toString()).put("userId", userId).toString();
    }

    private void callBackendWebhook(String jsonPayload)  {
        log.info(BackendHost + "/spi-webhooks");
        HttpRequest request = HttpRequest.newBuilder().uri(URI.create(BackendHost + "/spi-webhook")).POST(HttpRequest.BodyPublishers.ofString( jsonPayload)).build();
        log.info("Webhook call");
        try {
            HttpResponse<String> response = httpClient.send(request, HttpResponse.BodyHandlers.ofString());
            log.info(response.toString());
        } catch (Exception exception){
            log.warning("Cannot call backend webhook due to an error" + exception.toString());
        }
    }

    @Override
    public void onEvent(AdminEvent adminEvent, boolean b) {
        if (!adminEvent.getResourceType().equals(ResourceType.USER)){
            log.info("No user edit but: "+adminEvent.getResourceType().toString());
            return;
        }
        if (adminEvent.getOperationType().equals(OperationType.CREATE)){
            callBackendWebhook(buildJSONBody(EventType.REGISTER,adminEvent.getResourcePath().substring("users/".length())));
        }
        else if (adminEvent.getOperationType().equals(OperationType.DELETE)){
            callBackendWebhook(buildJSONBody(EventType.DELETE_ACCOUNT,adminEvent.getResourcePath().substring("users/".length())));
        }
        else {
            log.info("Wrong op type: " + adminEvent.getOperationType().toString());
        }

    }

    @Override
    public void close() {

    }
}
