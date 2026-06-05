package service

import (
	"banco-api/internal/repository"
	"errors"
	"shared"
)

type ClienteService struct {
	clienteRepo *repository.ClienteRepository
	contaRepo   *repository.ContaRepository
}

func NewClienteService(
	clienteRepo *repository.ClienteRepository,
	contaRepo *repository.ContaRepository,
) *ClienteService {
	return &ClienteService{
		clienteRepo: clienteRepo,
		contaRepo:   contaRepo,
	}
}

func (s *ClienteService) Cadastrar(
	nome,
	cpf,
	senha string,
	tipo shared.TipoConta,
) (*shared.Conta, error) {

	cliente, existe := s.clienteRepo.BuscarPorCPF(cpf)

	if !existe {

		cliente = &shared.Cliente{
			Nome:  nome,
			CPF:   cpf,
			Senha: senha,
		}

		s.clienteRepo.Salvar(cliente)
	}

	_, existeConta := s.contaRepo.BuscarPorClienteETipo(
		cpf,
		tipo,
	)

	if existeConta {
		return nil, errors.New(
			"cliente já possui esse tipo de conta",
		)
	}

	conta := &shared.Conta{
		Numero:  s.contaRepo.GerarNumeroConta(tipo),
		Saldo:   0,
		Titular: cliente,
		Tipo:    tipo,
	}

	s.contaRepo.Salvar(conta)

	return conta, nil
}

func (s *ClienteService) Login(
	cpf,
	senha string,
	tipo shared.TipoConta,
) (*shared.Conta, error) {

	cliente, existeCPF := s.clienteRepo.BuscarPorCPF(cpf)

	if !existeCPF {
		return nil, errors.New("cliente não encontrado")
	}

	if cliente.Senha != senha {
		return nil, errors.New("senha incorreta")
	}

	conta, existeConta := s.contaRepo.BuscarPorClienteETipo(cpf, tipo)

	if !existeConta {
		return nil, errors.New("cliente não possui esse tipo de conta")
	}

	return conta, nil
}
