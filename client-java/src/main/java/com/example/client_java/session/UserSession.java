package com.example.client_java.session;

import org.springframework.stereotype.Component;

import java.io.Serializable;

@Component
public class UserSession implements Serializable {
    private static final long serialVersionUID = 1L;

    private String token;
    private String contaNumero;
    private boolean bloquearAtualizacao = false;

    public UserSession() {}

    public String getToken() { return token; }
    public void setToken(String token) { this.token = token; }

    public String getContaNumero() { return contaNumero; }
    public void setContaNumero(String contaNumero) { this.contaNumero = contaNumero; }

    public boolean isBloquearAtualizacao() { return bloquearAtualizacao; }
    public void setBloquearAtualizacao(boolean bloquearAtualizacao) { this.bloquearAtualizacao = bloquearAtualizacao; }

    public void limparSessao() {
        this.token = null;
        this.contaNumero = null;
    }
}
