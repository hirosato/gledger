package application

import (
	"strings"

	"github.com/hirosato/gledger/domain"
)

type AccountService struct {
	repository domain.AccountRepository
}

func NewAccountService(repository domain.AccountRepository) *AccountService {
	return &AccountService{
		repository: repository,
	}
}

func (as *AccountService) FindOrCreateAccount(fullName string) *domain.Account {
	return as.repository.FindOrCreateAccount(fullName)
}

func (as *AccountService) GetAccount(fullName string) *domain.Account {
	return as.repository.FindAccount(fullName)
}

func (as *AccountService) GetAllAccounts() []*domain.Account {
	return as.repository.GetAllAccounts()
}

func (as *AccountService) GetRootAccount() *domain.Account {
	return as.repository.GetRootAccount()
}

type InMemoryAccountRepository struct {
	root     *domain.Account
	accounts map[string]*domain.Account
}

func NewInMemoryAccountRepository() *InMemoryAccountRepository {
	root := domain.NewAccount("")
	return &InMemoryAccountRepository{
		root:     root,
		accounts: make(map[string]*domain.Account),
	}
}

func (repo *InMemoryAccountRepository) FindAccount(fullName string) *domain.Account {
	if fullName == "" {
		return repo.root
	}
	return repo.accounts[fullName]
}

func (repo *InMemoryAccountRepository) CreateAccount(fullName string) *domain.Account {
	if fullName == "" {
		return repo.root
	}
	
	parts := strings.Split(fullName, ":")
	current := repo.root
	
	for _, part := range parts {
		current = current.FindOrCreateChild(part)
	}
	
	current.Type = domain.DetermineAccountType(fullName)
	repo.accounts[fullName] = current
	
	return current
}

func (repo *InMemoryAccountRepository) FindOrCreateAccount(fullName string) *domain.Account {
	if existing := repo.FindAccount(fullName); existing != nil {
		return existing
	}
	return repo.CreateAccount(fullName)
}

func (repo *InMemoryAccountRepository) GetRootAccount() *domain.Account {
	return repo.root
}

func (repo *InMemoryAccountRepository) GetAllAccounts() []*domain.Account {
	accounts := make([]*domain.Account, 0, len(repo.accounts))
	for _, account := range repo.accounts {
		accounts = append(accounts, account)
	}
	return accounts
}