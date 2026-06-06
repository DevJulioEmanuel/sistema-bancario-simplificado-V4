package com.example.client_java;

import com.example.client_java.interfaces.TelaContaFactory;
import com.example.client_java.model.ContaLogada;
import com.example.client_java.ui.Banner;
import com.example.client_java.ui.TelaCadastro;
import com.example.client_java.ui.TelaConta;
import com.example.client_java.ui.TelaLogin;
import org.springframework.boot.CommandLineRunner;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.ApplicationContext;

import java.util.Scanner;

@SpringBootApplication
public class BancoCliente implements CommandLineRunner {

	private final TelaLogin telaLogin;
	private final TelaCadastro telaCadastro;
	private final TelaContaFactory telaContaFactory;

	public BancoCliente(TelaLogin telaLogin, TelaCadastro telaCadastro, ApplicationContext context, TelaContaFactory telaContaFactory) {
		this.telaLogin = telaLogin;
		this.telaCadastro = telaCadastro;
		this.telaContaFactory = telaContaFactory;
	}

	public static void main(String[] args) {
		SpringApplication app = new SpringApplication(BancoCliente.class);

		app.setBannerMode(org.springframework.boot.Banner.Mode.CONSOLE);

		app.setLogStartupInfo(true);

		app.run(args);
	}

	@Override
	public void run(String... args) {
		Banner.limpar();
		Banner.exibirCabecalho();
		pausa(1000);

		Scanner sc = new Scanner(System.in);
		boolean rodando = true;

		while (rodando) {
			Banner.limpar();
			Banner.exibirCabecalho("MENU PRINCIPAL");
			exibirMenuPrincipal();

			String opcao = sc.nextLine().trim();
			switch (opcao) {
				case "1" -> {
					ContaLogada sessao = telaLogin.exibir();
					if (sessao != null) {
						TelaConta telaConta = telaContaFactory.criar(sessao);
						telaConta.exibir();
					}
				}
				case "2" -> {
					telaCadastro.exibir();
				}
				case "0" -> {
					Banner.limpar();
					Banner.exibirCabecalho();
					Banner.info("Obrigado por usar o Banco API. Até logo!");
					pausa(1500);
					rodando = false;
				}
				default -> {
					Banner.erro("Opção inválida. Tente novamente.");
					pausa(1000);
				}
			}
		}

		sc.close();
		System.exit(0);
	}

	private static void exibirMenuPrincipal() {
		String m = Banner.margem();
		System.out.println(m + Banner.SEPARADOR);
		System.out.println(m + "  [1]  Entrar na minha conta");
		System.out.println(m + "  [2]  Abrir nova conta");
		System.out.println(m + "  [0]  Sair");
		System.out.println(m + Banner.SEPARADOR);
		System.out.print(m + "  Escolha: ");
	}

	private static void pausa(long ms) {
		try { Thread.sleep(ms); } catch (InterruptedException ignored) {}
	}
}