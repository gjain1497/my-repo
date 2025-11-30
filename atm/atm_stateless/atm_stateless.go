package main

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"
)

var (
	ErrCardNotValid             = errors.New("card is not valid")
	ErrCardNotActive            = errors.New("card is not active")
	ErrCardExpired              = errors.New("card has expired")
	ErrCardDoesNotExist         = errors.New("card does not exist")
	ErrInvalidPIN               = errors.New("invalid PIN")
	ErrInsufficientCash         = errors.New("ATM has insufficient cash")
	ErrCannotDispense           = errors.New("cannot dispense exact amount")
	ErrInsufficientFunds        = errors.New("insufficient funds")
	ErrDailyLimitExceeded       = errors.New("daily limit exceeded")
	ErrAccountNotFound          = errors.New("account not found")
	ErrTransactionsDoesNotExist = errors.New("No transactions found for this account")
	ErrTransactionDoesNotExist  = errors.New("No transactions found for this transaction id")
)

type User struct {
	Id       string
	Name     string
	Phone    string
	Location Location
}

type BankService interface {
	ValidatePin(cardNumber, pin string) (bool, error)
	DebitAccount(accountId string, amount float64) error
	CreditAccount(accountId string, amount float64) error
}

type BankServiceV1 struct {
}

func (s *BankServiceV1) ValidatePin(cardNumber, pin string) (bool, error) {
	// Validate PIN with bank's system
	//maybe call Bank's API
	return true, nil
}

func (s *BankServiceV1) DebitAccount(accountId string, amount float64) error {
	return nil
}

func (s *BankServiceV1) CreditAccount(accountId string, amount float64) error {
	return nil
}

type Bank struct {
	Id       string
	Name     string
	Location Location
}

type Location struct {
	City    string
	Street  string
	Pincode string
}

type ATMService interface {
	Withdraw(cardNumber, pin string, amount float64) error
	Deposit(cardNumber, pin string, amount float64, denominations map[float64]int) error
	CheckBalance(cardNumber, pin string) (float64, error)
}

type ATMServiceV1 struct {
	//ATMs               map[string]*ATM
	ATM                *ATM //each ATM machine just manages one ATM
	BankService        BankService
	CardService        CardService
	AccountService     AccountService
	TransactionService TransactionService
	ReceiptService     ReceiptService
}

func NewATMServiceV1(atm *ATM, txnServ TransactionService, bankServ BankService, acctServ AccountService, cardServ CardService, receiptServ ReceiptService) (*ATMServiceV1, error) {
	return &ATMServiceV1{
		ATM:                atm,
		TransactionService: txnServ,
		BankService:        bankServ,
		AccountService:     acctServ,
		CardService:        cardServ,
		ReceiptService:     receiptServ,
	}, nil
}

func (s *ATMServiceV1) depositCash(amount float64, denominationsComing map[float64]int) error {
	// 19200 / 500
	originalAmount := amount
	denominations := s.ATM.CashInventory

	//500
	for denominationComing, count := range denominationsComing { //(500->2)
		denominations[denominationComing] += count
	}

	s.ATM.CurrBalance += originalAmount

	return nil
}
func (s *ATMServiceV1) dispenseCash(amount float64) error {
	// Check if ATM has enough cash
	if s.ATM.CurrBalance < amount {
		return ErrInsufficientCash
	}

	// TODO: Handle cash denominations from CashInventory

	// 19200 / 500
	originalAmount := amount
	denominations := s.ATM.CashInventory

	// Sort denominations in descending order
	var keys []float64
	for denomination := range denominations {
		keys = append(keys, denomination)
	}
	sort.Sort(sort.Reverse(sort.Float64Slice(keys)))

	for _, denomination := range keys {
		count := denominations[denomination]
		if count > 0 {
			countDenominationUsed := int(amount / denomination)
			if countDenominationUsed > count { //10000/500 = 20
				countDenominationUsed = count
			}
			s.ATM.CashInventory[denomination] -= countDenominationUsed
			amount -= float64(countDenominationUsed) * denomination
		}
	}

	//Check if we could dispense exact amount
	if amount > 0 {
		// Rollback the inventory changes
		// (In production, you'd track changes and rollback)
		return errors.New("cannot dispense exact amount with available denominations")
	}

	s.ATM.CurrBalance -= originalAmount

	return nil
}

func (s *ATMServiceV1) Withdraw(cardNumber, pin string, amount float64) error {
	//1. validate card
	err := s.CardService.ValidateCard(cardNumber)
	if err != nil {
		return err
	}

	//2. validate pin
	isValidPIN, err := s.BankService.ValidatePin(cardNumber, pin)
	if err != nil {
		return err
	}
	if !isValidPIN {
		return errors.New("invalid PIN")
	}
	accountId, err := s.CardService.GetAccountDetails(cardNumber)
	if err != nil {
		return err
	}

	err = s.AccountService.CanWithdraw(accountId, amount)
	if err != nil {
		return err
	}

	if s.ATM.CurrBalance < amount {
		return errors.New("ATM has insufficient cash")
	}

	err = s.AccountService.DebitAccount(accountId, amount)
	if err != nil {
		return err
	}

	err = s.dispenseCash(amount)
	if err != nil {
		return err
	}

	txn, err := s.TransactionService.CreateTransaction(
		accountId,
		s.ATM.Id, // Current ATM's ID
		amount,
		Withdraw, // Transaction type
	)

	// Generate receipt (optional - don't fail withdrawal if receipt fails)
	_, _ = s.ReceiptService.GenerateReceipt(txn.Id)

	return nil
}
func (s *ATMServiceV1) Deposit(cardNumber, pin string, amount float64, denominations map[float64]int) error {
	//1. validate card
	err := s.CardService.ValidateCard(cardNumber)
	if err != nil {
		return err
	}

	//2. validate pin
	isValidPIN, err := s.BankService.ValidatePin(cardNumber, pin)
	if err != nil {
		return err
	}
	if !isValidPIN {
		return errors.New("invalid PIN")
	}
	accountId, err := s.CardService.GetAccountDetails(cardNumber)
	if err != nil {
		return err
	}

	err = s.AccountService.CreditAccount(accountId, amount)
	if err != nil {
		return err
	}

	err = s.depositCash(amount, denominations)
	if err != nil {
		return err
	}

	_, err = s.TransactionService.CreateTransaction(
		accountId,
		s.ATM.Id, // Current ATM's ID
		amount,
		Deposit, // Transaction type
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *ATMServiceV1) CheckBalance(cardNumber, pin string) (float64, error) {
	// 1. Validate card
	err := s.CardService.ValidateCard(cardNumber)
	if err != nil {
		return 0, err
	}

	// 2. Validate PIN
	isValidPIN, err := s.BankService.ValidatePin(cardNumber, pin)
	if err != nil {
		return 0, err
	}
	if !isValidPIN {
		return 0, ErrInvalidPIN
	}

	// 3. Get account
	accountId, err := s.CardService.GetAccountDetails(cardNumber)
	if err != nil {
		return 0, err
	}

	// 4. Get account details
	account, err := s.AccountService.GetAccount(accountId)
	if err != nil {
		return 0, err
	}

	// 5. Create balance inquiry transaction (optional)
	_, _ = s.TransactionService.CreateTransaction(
		accountId,
		s.ATM.Id,
		0, // No amount for balance inquiry
		BalanceInquiry,
	)

	return account.CurrBalance, nil

}

type ATM struct {
	Id            string
	BankId        string
	Location      *Location
	CurrBalance   float64
	CashInventory map[float64]int //(denomination-> count) (500->5), (200->3), (100->23)
}

type AccountService interface {
	GetAccount(accountId string) (*Account, error)
	CanWithdraw(accountId string, amount float64) error
	DebitAccount(accountId string, amount float64) error
	CreditAccount(accountId string, amount float64) error
}

type AccountServiceV1 struct {
	accounts map[string]*Account //(account_id -> account object)
	mu       sync.RWMutex
}

// constructor
func NewAccountServiceV1() (*AccountServiceV1, error) {
	return &AccountServiceV1{
		accounts: make(map[string]*Account),
	}, nil
}

func (s *AccountServiceV1) GetAccount(accountId string) (*Account, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	account, exists := s.accounts[accountId]
	if !exists {
		return nil, ErrAccountNotFound
	}

	return account, nil
}
func (s *AccountServiceV1) CanWithdraw(accountId string, amount float64) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	account, exists := s.accounts[accountId]
	if !exists {
		return ErrAccountNotFound
	}
	// Check if sufficient balance
	if account.CurrBalance < amount {
		return ErrInsufficientFunds
	}
	// Check daily limit
	// TODO: Track today's withdrawals
	// For now, simple check:
	if amount > account.DailyLimit {
		return ErrDailyLimitExceeded
	}
	return nil

}
func (s *AccountServiceV1) DebitAccount(accountId string, amount float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	account, exists := s.accounts[accountId]
	if !exists {
		return ErrAccountNotFound
	}

	// Double-check balance (defensive programming)
	if account.CurrBalance < amount {
		return ErrInsufficientFunds
	}
	// Debit the account
	account.CurrBalance -= amount
	return nil
}
func (s *AccountServiceV1) CreditAccount(accountId string, amount float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	account, exists := s.accounts[accountId]
	if !exists {
		return ErrAccountNotFound
	}

	// Credit the account
	account.CurrBalance += amount
	return nil
}

type Account struct {
	Id          string
	UserId      string
	BankId      string
	CurrBalance float64
	AccountType AccountType
	DailyLimit  float64
}
type AccountType string

const (
	Savings AccountType = "SAVINGS"
	Current AccountType = "CURRENT"
)

type CardStatus string

const (
	Active  CardStatus = "ACTIVE"
	Blocked CardStatus = "BLOCKED"
	Expired CardStatus = "EXPIRED"
)

type CardService interface {
	ValidateCard(cardNumber string) error
	GetCard(cardNumber string) (*Card, error)
	BlockCard(cardNumber string) error
	GetAccountDetails(cardNumber string) (string, error)
}

type CardServiceV1 struct {
	Cards map[string]*Card //(cardNumber, card)
	mu    sync.RWMutex
}

func NewCardServiceV1() (*CardServiceV1, error) {
	return &CardServiceV1{
		Cards: make(map[string]*Card),
	}, nil
}

func (s *CardServiceV1) ValidateCard(cardNumber string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	card, exists := s.Cards[cardNumber]
	if !exists {
		return ErrCardDoesNotExist
	}

	//check card is active
	if card.Status != Active {
		return ErrCardNotActive
	}

	//check if card is expired
	if time.Now().After(card.ExpiryDate) {
		return ErrCardExpired
	}

	return nil
}

func (s *CardServiceV1) GetCard(cardNumber string) (*Card, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	card, exists := s.Cards[cardNumber]
	if !exists {
		return nil, ErrCardDoesNotExist
	}
	return card, nil
}

func (s *CardServiceV1) BlockCard(cardNumber string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	card, exists := s.Cards[cardNumber]
	if !exists {
		return ErrCardDoesNotExist
	}

	card.Status = Blocked
	return nil
}
func (s *CardServiceV1) GetAccountDetails(cardNumber string) (string, error) {
	card, ok := s.Cards[cardNumber]
	if !ok {
		return "", ErrCardNotValid
	}
	return card.AccountId, nil
}

type Card struct {
	CardNumber string
	UserId     string
	AccountId  string
	Name       string
	ExpiryDate time.Time
	Status     CardStatus
	// ❌ NO CVV - Illegal to store per PCI-DSS
	//    CVV is only for "card not present" online transactions

	// ❌ NO PIN - Stored securely in Bank's system (hashed)
	//    ATM sends PIN to Bank for validation
	//    Bank never shares PIN with ATM
}

type TransactionService interface {
	CreateTransaction(accountId, atmId string, amount float64, txnType TransactionType) (*Transaction, error)
	GetTransaction(transactionId string) (*Transaction, error)
	GetTransactionHistory(accountId string) ([]*Transaction, error)
}

type TransactionServiceV1 struct {
	totalTransactions   map[string]*Transaction   // txnId → transaction object
	accountTransactions map[string][]*Transaction // accountId → list of transactions
	mu                  sync.RWMutex              // For thread safety
}

func NewTransactionServiceV1() (*TransactionServiceV1, error) {
	return &TransactionServiceV1{
		totalTransactions:   make(map[string]*Transaction),
		accountTransactions: make(map[string][]*Transaction),
	}, nil
}
func (s *TransactionServiceV1) CreateTransaction(accountId, atmId string, amount float64, txnType TransactionType) (*Transaction, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	txnId := generateTransactionId()

	transaction := &Transaction{
		Id:        txnId,
		AccountId: accountId,
		Amount:    amount,
		Type:      txnType,
		Status:    Success, // Assume success at creation
		ATMId:     atmId,
		CreatedAt: time.Now(),
	}
	// Store transaction //kinda inserting in db table
	s.totalTransactions[txnId] = transaction
	s.accountTransactions[accountId] = append(s.accountTransactions[accountId], transaction)

	return transaction, nil
}

func (s *TransactionServiceV1) GetTransaction(transactionId string) (*Transaction, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	transaction, ok := s.totalTransactions[transactionId]
	if !ok {
		return nil, ErrTransactionDoesNotExist
	}

	return transaction, nil
}

func (s *TransactionServiceV1) GetTransactionHistory(accountId string) ([]*Transaction, error) {

	transactions, ok := s.accountTransactions[accountId]
	if !ok {
		return nil, ErrTransactionsDoesNotExist
	}

	return transactions, nil
}

type Transaction struct {
	Id        string
	AccountId string //1 account can have many transaction
	Amount    float64
	Type      TransactionType
	Status    TransactionStatus
	ATMId     string
	CreatedAt time.Time
}
type TransactionType string

const (
	Withdraw       TransactionType = "WITHDRAW"
	Deposit        TransactionType = "DEPOSIT"
	BalanceInquiry TransactionType = "BALANCE_INQUIRY"
)

type TransactionStatus string

const (
	Success TransactionStatus = "SUCCESS"
	Failed  TransactionStatus = "FAILED"
	Pending TransactionStatus = "PENDING"
)

type ReceiptService interface {
	GenerateReceipt(transactionId string) (*Receipt, error)
}

type ReceiptServiceV1 struct {
	TransactionService TransactionService
	AccountService     AccountService
	Receipts           map[string]*Receipt //reciept_id, receipt object
}

func NewReceiptServiceV1(txnServ TransactionService, acctServ AccountService) (*ReceiptServiceV1, error) {
	return &ReceiptServiceV1{
		TransactionService: txnServ,
		AccountService:     acctServ,
		Receipts:           make(map[string]*Receipt),
	}, nil
}

func (s *ReceiptServiceV1) GenerateReceipt(transactionId string) (*Receipt, error) {
	// 1. Get transaction details
	transaction, err := s.TransactionService.GetTransaction(transactionId)
	if err != nil {
		return nil, err
	}
	// 2. Get account details (to get balance)
	account, err := s.AccountService.GetAccount(transaction.AccountId)
	if err != nil {
		return nil, err
	}
	// 3. Generate receipt
	receipt := &Receipt{
		TransactionId: transactionId,
		Amount:        transaction.Amount,
		Balance:       account.CurrBalance, // Balance AFTER transaction
		CreatedAt:     time.Now(),
		Summary:       fmt.Sprintf("%s of %.2f completed", transaction.Type, transaction.Amount),
	}
	// 4. Store receipt (optional)
	s.Receipts[transactionId] = receipt

	return receipt, nil
}

type Receipt struct {
	TransactionId string
	Amount        float64
	Balance       float64
	CreatedAt     time.Time
	Summary       string
}

// Helper function to generate transaction ID
func generateTransactionId() string {
	return fmt.Sprintf("TXN-%d", time.Now().UnixNano())
}

// ATM                *ATM //each ATM machine just manages one ATM
//
//	BankService        BankService
//	CardService        CardService
//	AccountService     AccountService
//	TransactionService TransactionService
//	ReceiptService     ReceiptService
func main() {
	fmt.Println("=== ATM System Starting ===\n")

	// 1. Initialize all services
	bankServ := &BankServiceV1{}

	cardServ, err := NewCardServiceV1()
	if err != nil {
		fmt.Println("Error creating card service:", err)
		return
	}

	accountServ, err := NewAccountServiceV1()
	if err != nil {
		fmt.Println("Error creating account service:", err)
		return
	}

	transactionServ, err := NewTransactionServiceV1()
	if err != nil {
		fmt.Println("Error creating transaction service:", err)
		return
	}

	receiptServ, err := NewReceiptServiceV1(transactionServ, accountServ)
	if err != nil {
		fmt.Println("Error creating receipt service:", err)
		return
	}

	// 2. Create ATM with cash inventory
	cashInv := make(map[float64]int)
	cashInv[500] = 10 // 5,000
	cashInv[200] = 20 // 4,000
	cashInv[100] = 30 // 3,000
	// Total: 12,000

	atm := &ATM{
		Id:     "ATM-001",
		BankId: "BANK-001",
		Location: &Location{
			City:    "Hyderabad",
			Street:  "Lanco Hills",
			Pincode: "500089",
		},
		CurrBalance:   12000, // Sum of all denominations
		CashInventory: cashInv,
	}

	// 3. Initialize ATM Service
	atmServ, err := NewATMServiceV1(atm, transactionServ, bankServ, accountServ, cardServ, receiptServ)
	if err != nil {
		fmt.Println("Error creating ATM service:", err)
		return
	}

	// 4. CREATE TEST DATA (CRITICAL!)
	fmt.Println("Creating test user, account, and card...\n")

	// Create test account
	testAccount := &Account{
		Id:          "ACC-001",
		UserId:      "USER-001",
		BankId:      "BANK-001",
		CurrBalance: 50000, // Starting balance: ₹50,000
		AccountType: Savings,
		DailyLimit:  20000, // Daily limit: ₹20,000
	}
	accountServ.accounts["ACC-001"] = testAccount

	// Create test card
	testCard := &Card{
		CardNumber: "CARD-001",
		UserId:     "USER-001",
		AccountId:  "ACC-001",
		Name:       "John Doe",
		ExpiryDate: time.Now().AddDate(2, 0, 0), // Valid for 2 years
		Status:     Active,
	}
	cardServ.Cards["CARD-001"] = testCard

	fmt.Println("✅ Test data created!")
	fmt.Printf("   Card: %s\n", testCard.CardNumber)
	fmt.Printf("   Account: %s\n", testAccount.Id)
	fmt.Printf("   Initial Balance: ₹%.2f\n\n", testAccount.CurrBalance)

	// 5. TEST DEPOSIT
	fmt.Println("=== Testing Deposit ===")
	depositCashInv := make(map[float64]int)
	depositCashInv[500] = 2  // 1,000
	depositCashInv[200] = 10 // 2,000
	depositCashInv[100] = 30 // 3,000
	// Total deposit: ₹6,000

	err = atmServ.Deposit("CARD-001", "1234", 6000, depositCashInv)
	if err != nil {
		fmt.Println("❌ Error in deposit:", err)
	} else {
		fmt.Println("✅ Deposit successful!")
		fmt.Printf("   Amount: ₹6,000\n")
		fmt.Printf("   New Balance: ₹%.2f\n\n", testAccount.CurrBalance)
	}

	// 6. TEST WITHDRAW
	fmt.Println("=== Testing Withdrawal ===")
	err = atmServ.Withdraw("CARD-001", "1234", 5000)
	if err != nil {
		fmt.Println("❌ Error in withdrawal:", err)
	} else {
		fmt.Println("✅ Withdrawal successful!")
		fmt.Printf("   Amount: ₹5,000\n")
		fmt.Printf("   New Balance: ₹%.2f\n\n", testAccount.CurrBalance)
	}

	// 7. TEST CHECK BALANCE
	fmt.Println("=== Testing Check Balance ===")
	balance, err := atmServ.CheckBalance("CARD-001", "1234")
	if err != nil {
		fmt.Println("❌ Error checking balance:", err)
	} else {
		fmt.Println("✅ Balance check successful!")
		fmt.Printf("   Current Balance: ₹%.2f\n\n", balance)
	}

	// 8. TEST TRANSACTION HISTORY
	fmt.Println("=== Transaction History ===")
	transactions, err := transactionServ.GetTransactionHistory("ACC-001")
	if err != nil {
		fmt.Println("❌ Error getting transaction history:", err)
	} else {
		fmt.Printf("Total Transactions: %d\n", len(transactions))
		for i, txn := range transactions {
			fmt.Printf("%d. %s - ₹%.2f - %s - %s\n",
				i+1, txn.Type, txn.Amount, txn.Status, txn.CreatedAt.Format("2006-01-02 15:04:05"))
		}
	}

	fmt.Println("\n=== ATM System Test Complete ===")
}
