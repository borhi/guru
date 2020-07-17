package services

import (
	"context"
	"errors"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"guru/models"
	"guru/repositories"
	"net/http"
	"sync"
	"time"
)

type UserService struct {
	Users                 map[uint64]*models.UserModel
	Statistic             map[uint64]*models.StatisticModel
	UserRepository        repositories.UserRepository
	DepositRepository     repositories.DepositRepository
	TransactionRepository repositories.TransactionRepository
	Ticker                *time.Ticker
	sync.Mutex
}

func (s *UserService) Run(ctx context.Context, r *mux.Router) {
	users := make(map[uint64]*models.UserModel)
	if err := s.UserRepository.FindAll(users); err != nil {
		zap.L().Fatal(err.Error())
	}
	s.Users = users

	statistic := make(map[uint64]*models.StatisticModel)
	if err := s.DepositRepository.FindAllDeposit(statistic); err != nil {
		zap.L().Fatal(err.Error())
	}
	if err := s.TransactionRepository.FindAllBet(statistic); err != nil {
		zap.L().Fatal(err.Error())
	}
	if err := s.TransactionRepository.FindAllWin(statistic); err != nil {
		zap.L().Fatal(err.Error())
	}
	s.Statistic = statistic

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	s.startTicker()

	go func() {
		err := server.ListenAndServe()
		zap.L().Warn("http server terminated", zap.String("error", err.Error()))
	}()
	zap.L().Info("server started")


	<-ctx.Done()
	s.stopTicker()

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()
	zap.L().Info("shutdown initiated")
	if err := server.Shutdown(ctxShutDown); err != nil {
		zap.L().Error("http server shutdown: %v", zap.String("error", err.Error()))
	}
	zap.L().Info("shutdown completed")
}

func (s *UserService) GetUser(id uint64, token string) (*models.GetUserResponseModel, error) {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.Users[id]; !ok {
		return nil, errors.New("not found")
	}

	if s.Users[id].Token != token {
		return nil, errors.New("wrong token")
	}

	return &models.GetUserResponseModel{
		Id:           s.Users[id].Id,
		Balance:      s.Users[id].Balance,
		DepositCount: s.Statistic[id].DepositCount,
		DepositSum:   s.Statistic[id].DepositSum,
		BetCount:     s.Statistic[id].BetCount,
		BetSum:       s.Statistic[id].BetSum,
		WinCount:     s.Statistic[id].WinCount,
		WinSum:       s.Statistic[id].WinSum,
	}, nil
}

func (s *UserService) CreateUser(id uint64, user models.UserModel) {
	s.Lock()
	defer s.Unlock()
	s.Users[id] = &user
	s.Users[id].Status = models.StatusNew
	s.Statistic[id] = &models.StatisticModel{
		Id:           id,
		DepositCount: 0,
		DepositSum:   0,
		BetCount:     0,
		BetSum:       0,
		WinCount:     0,
		WinSum:       0,
	}
}

func (s *UserService) AddDeposit(depositRequest models.DepositRequestModel) (*models.TransactionResponseModel, error) {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.Users[depositRequest.UserId]; !ok {
		return nil, errors.New("not found")
	}

	if s.Users[depositRequest.UserId].Token != depositRequest.Token {
		return nil, errors.New("wrong token")
	}

	if err := s.saveDeposit(depositRequest); err != nil {
		return nil, err
	}

	s.Users[depositRequest.UserId].Balance += depositRequest.Amount
	s.Statistic[depositRequest.UserId].DepositCount += 1
	s.Statistic[depositRequest.UserId].DepositSum += depositRequest.Amount
	s.Users[depositRequest.UserId].Status = models.StatusModified

	return &models.TransactionResponseModel{
		Error:   "",
		Balance: s.Users[depositRequest.UserId].Balance,
	}, nil
}

func (s *UserService) Transaction(transactionRequest models.TransactionRequestModel) (*models.TransactionResponseModel, error) {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.Users[transactionRequest.UserId]; !ok {
		return nil, errors.New("not found")
	}

	if s.Users[transactionRequest.UserId].Token != transactionRequest.Token {
		return nil, errors.New("wrong token")
	}

	if transactionRequest.Type == models.TypeBet && s.Users[transactionRequest.UserId].Balance < transactionRequest.Amount {
		return nil, errors.New("not enough balance")
	}

	if err := s.saveTransaction(transactionRequest); err != nil {
		return nil, err
	}

	if transactionRequest.Type == models.TypeWin {
		s.Users[transactionRequest.UserId].Balance += transactionRequest.Amount
		s.Statistic[transactionRequest.UserId].WinCount += 1
		s.Statistic[transactionRequest.UserId].WinSum += transactionRequest.Amount
	}
	if transactionRequest.Type == models.TypeBet {
		s.Users[transactionRequest.UserId].Balance -= transactionRequest.Amount
		s.Statistic[transactionRequest.UserId].BetCount += 1
		s.Statistic[transactionRequest.UserId].BetSum -= transactionRequest.Amount
	}
	s.Users[transactionRequest.UserId].Status = models.StatusModified

	return &models.TransactionResponseModel{
		Error:   "",
		Balance: s.Users[transactionRequest.UserId].Balance,
	}, nil
}

func (s *UserService) saveDeposit(depositRequest models.DepositRequestModel) error {
	deposit := models.DepositModel{
		Id:            depositRequest.DepositId,
		UserId:        depositRequest.UserId,
		Amount:        depositRequest.Amount,
		BalanceBefore: s.Users[depositRequest.UserId].Balance,
		BalanceAfter:  s.Users[depositRequest.UserId].Balance + depositRequest.Amount,
		CreatedAt:     time.Time{},
	}

	if err := s.DepositRepository.Insert(deposit); err != nil {
		return err
	}

	return nil
}

func (s *UserService) saveTransaction(transactionRequest models.TransactionRequestModel) error {
	balanceAfter := s.Users[transactionRequest.UserId].Balance - transactionRequest.Amount
	if transactionRequest.Type == models.TypeWin {
		balanceAfter = s.Users[transactionRequest.UserId].Balance + transactionRequest.Amount
	}

	transaction := models.TransactionModel{
		Id:            transactionRequest.TransactionId,
		UserId:        transactionRequest.UserId,
		Amount:        transactionRequest.Amount,
		Type:          transactionRequest.Type,
		BalanceBefore: s.Users[transactionRequest.UserId].Balance,
		BalanceAfter:  balanceAfter,
		CreatedAt:     time.Time{},
	}

	if err := s.TransactionRepository.Insert(transaction); err != nil {
		return err
	}

	return nil
}

func (s *UserService) startTicker() {
	go func() {
		for range s.Ticker.C {
			if err := s.saveUser(); err != nil {
				zap.L().Error(err.Error())
			}
		}
	}()
}

func (s *UserService) stopTicker() {
	s.Ticker.Stop()
	if err := s.saveUser(); err != nil {
		zap.L().Error(err.Error())
	}

}

func (s *UserService) saveUser() error {
	s.Lock()
	defer s.Unlock()

	var newUsers []interface{}
	for k := range s.Users {
		if s.Users[k].Status == models.StatusNew {
			newUsers = append(newUsers, *s.Users[k])
		}

		if s.Users[k].Status == models.StatusModified {
			if err := s.UserRepository.Update(s.Users[k]); err != nil {
				return err
			}
		}
	}

	if err := s.UserRepository.Insert(newUsers); err != nil {
		return err
	}

	return nil
}
