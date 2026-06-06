package com.example.client_java.config;

import com.example.client_java.interfaces.TelaContaFactory;
import com.example.client_java.service.BancoApiClient;
import com.example.client_java.service.WebSocketService;
import com.example.client_java.ui.TelaConta;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.web.client.RestClient;

import java.util.Scanner;

@Configuration
public class RestClientConfig {

    @Bean
    public RestClient restClient() {
        return RestClient.builder()
                .baseUrl("http://localhost:8080")
                .build();
    }

    @Bean
    public Scanner scanner() {
        return new Scanner(System.in);
    }

    @Bean
    public TelaContaFactory telaContaFactory(Scanner sc, BancoApiClient bancoApiClient, WebSocketService webSocketService) {
        return conta -> new TelaConta(sc, bancoApiClient, conta, webSocketService);
    }

}
