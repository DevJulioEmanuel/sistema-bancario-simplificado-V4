package com.example.client_java.service;

import com.example.client_java.model.request.*;
import com.example.client_java.model.response.*;
import com.example.client_java.session.UserSession;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestClient;

@Service
public class BancoApiClient {

    private final RestClient restClient;
    private final UserSession  userSession;

    public BancoApiClient(RestClient restClient, UserSession userSession) {
        this.restClient = restClient;
        this.userSession = userSession;
    }

    public CadastroResponse cadastrar(CadastroRequest cadastroRequest) {
        return restClient.post()
                .uri("/clientes")
                .body(cadastroRequest)
                .retrieve()
                .body(CadastroResponse.class);
    }

    public LoginResponse login(LoginRequest loginRequest) {
        LoginResponse response = restClient.post()
                .uri("/clientes/login")
                .body(loginRequest)
                .retrieve()
                .body(LoginResponse.class);

        if (response != null) {
            userSession.setToken(response.token());
            userSession.setContaNumero(String.valueOf(response.numero()));
        }

        return response;
    }

    public ContaResponse obterDados(int numero) {
        return restClient.get()
                .uri("/contas/" + numero + "/")
                .header("Authorization", "Bearer " + userSession.getToken())
                .retrieve()
                .body(ContaResponse.class);
    }

    public PagamentoResponse pagar(int numero, PagamentoRequest  pagamentoRequest) {
        return restClient.post()
                .uri("/contas/" + numero + "/pagamento")
                .header("Authorization", "Bearer " + userSession.getToken())
                .body(pagamentoRequest)
                .retrieve()
                .body(PagamentoResponse.class);
    }

    public MovimentacaoResponse depositar(int numero, ValorRequest  valorRequest) {
        return restClient.post()
                .uri("/contas/" + numero + "/depositar")
                .header("Authorization", "Bearer " + userSession.getToken())
                .body(valorRequest)
                .retrieve()
                .body(MovimentacaoResponse.class);
    }

    public MovimentacaoResponse sacar(int numero, ValorRequest  valorRequest) {
        return restClient.post()
                .uri("/contas/" + numero + "/sacar")
                .header("Authorization", "Bearer " + userSession.getToken())
                .body(valorRequest)
                .retrieve()
                .body(MovimentacaoResponse.class);
    }

    public TransferenciaResponse transferir(int numero, TransferenciaRequest valorRequest) {
        return restClient.post()
                .uri("/contas/" + numero + "/transferir")
                .header("Authorization", "Bearer " + userSession.getToken())
                .body(valorRequest)
                .retrieve()
                .body(TransferenciaResponse.class);
    }

    public ExtratoResponse extrato(int numero) {
        return restClient.get()
                .uri("/contas/" + numero + "/extrato")
                .header("Authorization", "Bearer " + userSession.getToken())
                .retrieve()
                .body(ExtratoResponse.class);
    }

    public RendimentoResponse rendimento(int numero, RendimentoRequest rendimentoRequest) {
        return restClient.get()
                .uri("/contas/" + numero + "/rendimento/" + rendimentoRequest.meses())
                .header("Authorization", "Bearer " + userSession.getToken())
                .retrieve()
                .body(RendimentoResponse.class);
    }

}
