package com.example.client_java.session;

import org.springframework.stereotype.Component;
import org.springframework.web.context.annotation.SessionScope;

import java.io.Serializable;

@Component
@SessionScope
public class UserSession implements Serializable {
    private static final long serialVersionUID = 1L;

    private String token;
    private String contaNumero;

    public UserSession() {}

    public String getToken() { return token; }
    public void setToken(String token) { this.token = token; }

    public String getContaNumero() { return contaNumero; }
    public void setContaNumero(String contaNumero) { this.contaNumero = contaNumero; }

    public void limparSessao() {
        this.token = null;
        this.contaNumero = null;
    }
}
