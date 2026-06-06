package com.example.client_java.service;

import com.example.client_java.session.UserSession;
import org.springframework.stereotype.Service;
import org.springframework.web.socket.TextMessage;
import org.springframework.web.socket.WebSocketSession;
import org.springframework.web.socket.client.standard.StandardWebSocketClient;
import org.springframework.web.socket.handler.TextWebSocketHandler;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;

import java.net.URI;

@Service
public class WebSocketService {
    private final UserSession userSession;
    private final ObjectMapper objectMapper;

    public WebSocketService(UserSession userSession) {
        this.userSession = userSession;
        this.objectMapper = new ObjectMapper();
    }

    public void conectar(int numeroConta, Runnable onUpdate){
        Thread thread = new Thread(() -> {
            try {
                StandardWebSocketClient client = new StandardWebSocketClient();
                String url = "ws://localhost:8080/ws/" + numeroConta;

                client.execute(new TextWebSocketHandler() {
                    @Override
                    protected void handleTextMessage(WebSocketSession session, TextMessage message) throws Exception {
                        JsonNode json = objectMapper.readTree(message.getPayload());
                        if ("saldo_atualizado".equals(json.get("evento").asText())) {
                            onUpdate.run();
                        }
                    }
                }, null, URI.create(url)).get();
            } catch (Exception e){
                e.printStackTrace();
            }
        });

        thread.setDaemon(true);
        thread.start();
    }
}
