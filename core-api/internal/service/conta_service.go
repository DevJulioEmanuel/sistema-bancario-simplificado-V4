package service

import (
	"banco-api/internal/repository"
	"errors"
	"math"
	"shared"
)

type ContaService struct {
	repo *repository.ContaRepository
}

func NewContaService(repo *repository.ContaRepository) *ContaService {
	return &ContaService{repo: repo}
}

func (s *ContaService) BuscarConta(numero int) (*shared.Conta, error) {

	conta, ok := s.repo.BuscarPorNumero(numero)

	if !ok {
		return nil, errors.New("conta não encontrada")
	}

	return conta, nil
}

func (s *ContaService) ObterExtrato(numero int) (*shared.Conta, error) {
	conta, err := s.BuscarConta(numero)

	if err != nil {
		return nil, err
	}

	return conta, nil
}

func (s *ContaService) CalcularRendimento(numero int, meses int) (float64, float64, *shared.Conta, error) {

	conta, err := s.BuscarConta(numero)

	if err != nil {
		return 0, 0, nil, err
	}

	if conta.Tipo != shared.TipoPoupanca {
		return 0, 0, nil, errors.New("apenas contas poupanca possuem rendimento")
	}

	saldoProjetado := conta.Saldo * math.Pow(1+shared.TaxaMensal, float64(meses))

	rendimento := saldoProjetado - conta.Saldo

	saldoProjetado = math.Round(saldoProjetado*100) / 100

	rendimento = math.Round(rendimento*100) / 100

	return saldoProjetado, rendimento, conta, nil
}
