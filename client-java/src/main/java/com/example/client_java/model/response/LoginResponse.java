package com.example.client_java.model.response;

public record LoginResponse(
        String mensagem,
        String nome,
        int numero,
        String token
        ) {}
