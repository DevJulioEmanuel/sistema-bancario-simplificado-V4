package com.example.client_java.ui;

import com.example.client_java.model.ContaLogada;
import com.example.client_java.model.request.PagamentoRequest;
import com.example.client_java.model.request.RendimentoRequest;
import com.example.client_java.model.request.TransferenciaRequest;
import com.example.client_java.model.request.ValorRequest;
import com.example.client_java.model.response.*;
import com.example.client_java.service.BancoApiClient;
import com.example.client_java.service.WebSocketService;
import org.springframework.web.client.HttpClientErrorException;

import java.text.DecimalFormat;
import java.util.Scanner;

public class TelaConta {

    private final Scanner sc;
    private final BancoApiClient bancoApiClient;
    private final ContaLogada contaLogada;
    private final WebSocketService webSocketService;

    private boolean bloquearAtualizacao = false;

    public TelaConta(Scanner sc, BancoApiClient bancoApiClient, ContaLogada contaLogada, WebSocketService webSocketService) {
        this.sc = sc;
        this.bancoApiClient = bancoApiClient;
        this.contaLogada = contaLogada;
        this.webSocketService = webSocketService;
    }

    public void exibir() {
        boolean logado = true;

        webSocketService.conectar(contaLogada.getNumeroConta(), () -> {
            try {
                if (!bloquearAtualizacao) {
                    ContaResponse dadosNovos = bancoApiClient.obterDados(contaLogada.getNumeroConta());
                    contaLogada.setSaldo(dadosNovos.saldo());
                    renderizarDashboard();
                    exibirMenu();
                }
            } catch (Exception e) {
            }
        });

        while (logado) {
            renderizarDashboard();
            exibirMenu();

            String opcao = sc.nextLine().trim();

            switch (opcao) {
                case "1" -> operacaoDepositar();
                case "2" -> operacaoSacar();
                case "3" -> operacaoTransferir();
                case "4" -> operacaoExtrato();
                case "5" -> {
                    if (contaLogada.isCorrente()) operacaoPagar();
                    else if (contaLogada.isPoupanca()) operacaoProjetarRendimento();
                    else opcaoInvalida();
                }
                case "6" -> operacaoNotificacoes();
                case "0" -> {
                    Banner.limpar();
                    Banner.exibirCabecalho();
                    Banner.info("Sessão encerrada. Até logo!");
                    pausa(1500);
                    logado = false;
                }
                default -> opcaoInvalida();
            }
        }
    }

    private void renderizarDashboard() {
        Banner.limpar();
        Banner.exibirCabecalho("MINHA CONTA");
        Banner.tituloSecao("PAINEL DO CLIENTE");

        Banner.labelSecao("TITULAR");
        Banner.campo("Nome",     Banner.branco(contaLogada.getNomeTitular()));
        Banner.campo("CPF",      Banner.cinza(contaLogada.getCpf()));
        Banner.espaco();

        Banner.labelSecao("CONTA");
        Banner.campo("Número",   Banner.branco(String.valueOf(contaLogada.getNumeroConta())));
        Banner.campo("Tipo",     contaLogada.isCorrente() ? Banner.badgeAzul("CORRENTE") : Banner.badgeVerde("POUPANÇA"));
        Banner.espaco();

        Banner.labelSecao("SALDO DISPONÍVEL");
        System.out.println(Banner.margem() + "  " + Banner.valorPositivo(contaLogada.getSaldo()));
        if (contaLogada.isCorrente()) {
            System.out.println(Banner.margem() + "  " + Banner.cinza(String.format("Limite: R$ %.2f", contaLogada.getLimite())));
        } else {
            System.out.println(Banner.margem() + "  " + Banner.cinza(String.format("Rendimento: %.3f%% a.m.", contaLogada.getRendimento() * 100)));
        }
        Banner.espaco();
        Banner.linhaSeparadora();
        Banner.espaco();
    }

    private void exibirMenu() {
        String m = Banner.margem();
        System.out.println(m + "  " + Banner.cinza("[1]  Depositar"));
        System.out.println(m + "  " + Banner.cinza("[2]  Sacar"));
        System.out.println(m + "  " + Banner.cinza("[3]  Transferir"));
        System.out.println(m + "  " + Banner.cinza("[4]  Ver Extrato"));

        if (contaLogada.isCorrente()) {
            System.out.println(m + "  " + Banner.cinza("[5]  Pagar Boleto / Conta"));
        } else {
            System.out.println(m + "  " + Banner.cinza("[5]  Projetar Rendimento"));
        }

        System.out.println(m + "  " + Banner.cinza("[6]  Ver Notificações"));

        System.out.println(m + "  " + Banner.cinza("[0]  Sair (Logout)"));
        Banner.espaco();
        Banner.prompt("O que deseja fazer?");
    }

    private void operacaoDepositar() {
        bloquearAtualizacao = true;

        Double valor = lerValor("Valor do depósito");
        if (valor == null) return;

        Banner.info("Processando...");
        try {
            ValorRequest valorRequest = new ValorRequest(valor);
            MovimentacaoResponse res = bancoApiClient.depositar(contaLogada.getNumeroConta(), valorRequest);
            Banner.sucesso(res.mensagem());
        } catch (HttpClientErrorException e) {
            String erro = e.getResponseBodyAsString();
            String mensagem = erro
                    .replaceAll(".*\"erro\":\"*", "")
                    .replaceAll("\"}.*", "")
                    .trim();

            Banner.erro("Erro ao depositar: " + mensagem);
            pausa(2000);
        }
        pausa(2000);

        bloquearAtualizacao = false;
    }

    private void operacaoSacar() {
        bloquearAtualizacao = true;

        Double valor = lerValor("Valor do saque");
        if (valor == null) return;

        Banner.info("Processando...");
        try {
            ValorRequest valorRequest = new ValorRequest(valor);
            MovimentacaoResponse res = bancoApiClient.sacar(contaLogada.getNumeroConta(), valorRequest);
            Banner.sucesso(res.mensagem());
        } catch (HttpClientErrorException e) {
            String erro = e.getResponseBodyAsString();
            String mensagem = erro
                    .replaceAll(".*\"erro\":\"*", "")
                    .replaceAll("\"}.*", "")
                    .trim();

            Banner.erro("Erro ao sacar: " + mensagem);
            pausa(2000);
        }
        pausa(2000);

        bloquearAtualizacao = false;
    }

    private void operacaoTransferir() {
        bloquearAtualizacao = true;

        Banner.espaco();
        Banner.labelSecao("TRANSFERÊNCIA");

        Banner.prompt("Conta destino      :");
        String destStr = sc.nextLine().trim();
        int numDestino;
        try {
            numDestino = Integer.parseInt(destStr);
        } catch (NumberFormatException e) {
            Banner.erro("Número de conta inválido.");
            pausa(1500);
            return;
        }

        if (numDestino == contaLogada.getNumeroConta()) {
            Banner.erro("Não é possível transferir para a própria conta.");
            pausa(2000);
            return;
        }

        Double valor = lerValor("Valor da transferência");
        if (valor == null) return;

        Banner.info("Processando...");
        try {
            TransferenciaRequest transferenciaRequest = new TransferenciaRequest(numDestino, valor);
            TransferenciaResponse res = bancoApiClient.transferir(contaLogada.getNumeroConta(), transferenciaRequest);
            Banner.sucesso(res.mensagem());
        } catch (HttpClientErrorException e) {
            String erro = e.getResponseBodyAsString();
            String mensagem = erro
                    .replaceAll(".*\"erro\":\"*", "")
                    .replaceAll("\"}.*", "")
                    .trim();

            Banner.erro("Erro ao transferir: " + mensagem);
            pausa(2000);
        }
        pausa(2500);

        bloquearAtualizacao = false;
    }

    private void operacaoPagar() {
        bloquearAtualizacao = true;

        Banner.espaco();
        Banner.labelSecao("PAGAMENTO DE BOLETO / CONTA");

        Banner.prompt("Descrição          :");
        String descricao = sc.nextLine().trim();
        if (descricao.isBlank()) {
            Banner.erro("Descrição obrigatória.");
            pausa(1500);
            return;
        }

        Double valor = lerValor("Valor do pagamento");
        if (valor == null) return;

        Banner.info("Processando...");
        try {
            PagamentoRequest pagamentoRequest = new PagamentoRequest(descricao, valor);
            PagamentoResponse res = bancoApiClient.pagar(contaLogada.getNumeroConta(), pagamentoRequest);
            Banner.sucesso(res.mensagem() + " — " + descricao);
        } catch (HttpClientErrorException e) {
            String erro = e.getResponseBodyAsString();
            String mensagem = erro
                    .replaceAll(".*\"erro\":\"*", "")
                    .replaceAll("\"}.*", "")
                    .trim();

            Banner.erro("Erro ao pagar: " + mensagem);
            pausa(2000);
        }
        pausa(2000);

        bloquearAtualizacao = false;
    }

    private void operacaoProjetarRendimento() {
        Banner.espaco();
        Banner.labelSecao("PROJEÇÃO DE RENDIMENTO");
        Banner.prompt("Número de meses    :");
        String mesesStr = sc.nextLine().trim();
        int meses;
        try {
            meses = Integer.parseInt(mesesStr);
            if (meses <= 0) throw new NumberFormatException();
        } catch (NumberFormatException e) {
            Banner.erro("Número de meses inválido.");
            pausa(1500);
            return;
        }

        Banner.info("Calculando...");
        try {
            RendimentoRequest rendimentoRequest = new RendimentoRequest(meses);
            RendimentoResponse saldoProjetado = bancoApiClient.rendimento(contaLogada.getNumeroConta(), rendimentoRequest);
            Banner.espaco();
            Banner.labelSecao("RESULTADO DA PROJEÇÃO");
            Banner.campo("Período",         meses + " meses");
            Banner.campo("Saldo atual",     Banner.formatarDinheiro(contaLogada.getSaldo()));
            Banner.campo("Saldo projetado", Banner.valorPositivo(saldoProjetado.saldoProjetado()));
            Banner.espaco();
            Banner.aguardeEnter(sc);
        } catch (Exception e) {
            Banner.erro("Erro ao projetar rendimento: " + e.getMessage());
            pausa(2000);
        }
    }

    private void operacaoExtrato() {
        Banner.limpar();
        Banner.exibirCabecalho("EXTRATO");
        Banner.tituloSecao("MOVIMENTAÇÕES DA CONTA " + contaLogada.getNumeroConta());

        Banner.info("Buscando movimentações...");
        try {
            ExtratoResponse res = bancoApiClient.extrato(contaLogada.getNumeroConta());
            Banner.espaco();
            if (res.transacoes() == null || res.transacoes().isEmpty()) {
                Banner.info("Nenhuma movimentação registrada.");
            } else {
                Banner.labelSecao("HISTÓRICO");
                for (TransacaoResponse t : res.transacoes()) {
                    DecimalFormat df = new java.text.DecimalFormat("+R$ #,##0.00;-R$ #,##0.00");
                    String valorFormatado = df.format(t.valor());
                    System.out.println(Banner.margem() + "  " + Banner.cinza("• ") + t.descricao() + ": " + valorFormatado);
                }
            }
        } catch (HttpClientErrorException e) {
            String erro = e.getResponseBodyAsString();
            String mensagem = erro
                    .replaceAll(".*\"erro\":\"*", "")
                    .replaceAll("\"}.*", "")
                    .trim();

            Banner.erro("Erro ao buscar extrato: " + mensagem);
            pausa(2000);
        }

        Banner.espaco();
        Banner.linhaSeparadora();
        Banner.espaco();
        Banner.aguardeEnter(sc);
    }

    private Double lerValor(String label) {
        Banner.espaco();
        Banner.prompt(label + "  : R$ ");
        String s = sc.nextLine().trim().replace(",", ".");
        try {
            double v = Double.parseDouble(s);
            if (v <= 0) throw new NumberFormatException();
            return v;
        } catch (NumberFormatException e) {
            Banner.erro("Valor inválido. Informe um número positivo.");
            pausa(1500);
            return null;
        }
    }

    private void operacaoNotificacoes() {
        Banner.limpar();
        Banner.exibirCabecalho("NOTIFICAÇÕES");
        Banner.tituloSecao("CAIXA DE ENTRADA");

        Banner.info("Buscando notificações no servidor...");

        try {
            java.util.List<java.util.Map<String, Object>> notificacoes = bancoApiClient.obterNotificacoes(contaLogada.getNumeroConta());

            Banner.espaco();
            if (notificacoes == null || notificacoes.isEmpty()) {
                Banner.info("Sua caixa de notificações está vazia.");
            } else {
                Banner.labelSecao("ALERTAS RECEBIDOS");

                for (java.util.Map<String, Object> notif : notificacoes) {
                    String msg = (String) notif.getOrDefault("mensagem", "Sem mensagem");
                    boolean sucesso = (boolean) notif.getOrDefault("sucesso", false);

                    if (sucesso) {
                        System.out.println(Banner.margem() + "  \u001B[92m• [SUCESSO] " + msg + "\u001B[0m");
                    } else {
                        System.out.println(Banner.margem() + "  \u001B[91m• [FALHA] " + msg + "\u001B[0m");
                    }
                }
            }
        } catch (Exception e) {
            Banner.erro("Erro ao recuperar notificações: " + e.getMessage());
        }

        Banner.espaco();
        Banner.linhaSeparadora();
        Banner.espaco();
        Banner.aguardeEnter(sc);
    }

    private void opcaoInvalida() {
        Banner.erro("Opção inválida.");
        pausa(1000);
    }

    private void pausa(long ms) {
        try { Thread.sleep(ms); } catch (InterruptedException ignored) {}
    }
}