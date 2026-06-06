import requests
from dataclasses import asdict

from models import CadastroResponse, LoginResponse, ContaResponse, PagamentoResponse, MovimentacaoResponse, \
    TransferenciaResponse, ExtratoResponse, RendimentoResponse, CadastroRequest, LoginRequest, PagamentoRequest, \
    ValorRequest, TransferenciaRequest, RendimentoRequest


class BancoApiClient:
    def __init__(self, base_url="http://localhost:8080"):
        self.base_url = base_url
        self.session = requests.Session()

    def _enviar(self, metodo, endpoint, payload=None):
        try:
            url = f"{self.base_url}{endpoint}"

            if payload and hasattr(payload, '__dataclass_fields__'):
                json_data = asdict(payload)
            else:
                json_data = payload

            response = self.session.request(metodo, url, json=json_data)
            response.raise_for_status()
            return response.json()

        except requests.exceptions.HTTPError:
            erro_msg = response.json().get("erro", "Erro desconhecido")
            raise Exception(erro_msg)
        except requests.exceptions.ConnectionError:
            raise Exception("Servidor em Go está offline.")

    def cadastrar(self, req: CadastroRequest) -> CadastroResponse:
        response = self._enviar("POST", "/clientes", req)
        return CadastroResponse(**response)

    def login(self, req: LoginRequest) -> LoginResponse:
        response = self._enviar("POST", "/clientes/login", req)

        token = response.get("token")
        if token:
            self.session.headers.update({"Authorization": f"Bearer {token}"})

        return LoginResponse(**response)

    def obterDados(self, numero: int) -> ContaResponse:
        response = self._enviar("GET", f"/contas/{numero}/")
        return ContaResponse(**response)

    def pagar(self, numero: int, req: PagamentoRequest) -> PagamentoResponse:
        response = self._enviar("POST", f"/contas/{numero}/pagamento", req)
        return PagamentoResponse(**response)

    def depositar(self, numero: int, req: ValorRequest) -> MovimentacaoResponse:
        response = self._enviar("POST", f"/contas/{numero}/depositar", req)
        return MovimentacaoResponse(**response)

    def sacar(self, numero: int, req: ValorRequest) -> MovimentacaoResponse:
        response = self._enviar("POST", f"/contas/{numero}/sacar", req)
        return MovimentacaoResponse(**response)

    def transferir(self, numero: int, req: TransferenciaRequest) -> TransferenciaResponse:
        response = self._enviar("POST", f"/contas/{numero}/transferir", req)
        return TransferenciaResponse(**response)

    def extrato(self, numero: int) -> ExtratoResponse:
        response = self._enviar("GET", f"/contas/{numero}/extrato")
        return ExtratoResponse(**response)

    def rendimento(self, numero: int, req: RendimentoRequest) -> RendimentoResponse:
        response = self._enviar("GET", f"/contas/{numero}/rendimento/{req.meses}")
        return RendimentoResponse(**response)

    def obter_notificacoes(self, numero: int) -> list:
        response = self._enviar("GET", f"/contas/{numero}/notificacoes")
        return response