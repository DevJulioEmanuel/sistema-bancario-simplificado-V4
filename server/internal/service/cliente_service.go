package service

import (
	"banco-api/internal/model"
	"banco-api/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"errors"
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
	tipo model.TipoConta,
) (*model.Conta, error) {

	cliente, existe := s.clienteRepo.BuscarPorCPF(cpf)

	if !existe {

		hash, err := bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
		if err != nil {
			return nil, errors.New("erro ao processar a segurança da senha")
		}

		cliente = &model.Cliente{
			Nome:  nome,
			CPF:   cpf,
			Senha: string(hash),
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

	conta := &model.Conta{
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
	tipo model.TipoConta,
) (*model.Conta, error) {

	cliente, existeCPF := s.clienteRepo.BuscarPorCPF(cpf)

	if !existeCPF {
		return nil, errors.New("Credenciais inválidas")
	}

	err := bcrypt.CompareHashAndPassword([]byte(cliente.Senha), []byte(senha))
	if err != nil {
		return nil, errors.New("credenciais inválidas")
	}

	conta, existeConta := s.contaRepo.BuscarPorClienteETipo(cpf, tipo)

	if !existeConta {
		return nil, errors.New("cliente não possui esse tipo de conta")
	}

	return conta, nil
}
