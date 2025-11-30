package main

import (
	"errors"
	"fmt"
	"sort"
	"strings"
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

type CashDispenser interface {
	Dispense(amount float64) error
	Deposit(denoms map[float64]int) error
	GetCurrentBalance() float64
}

type CashDispenserV1 struct {
	CurrBalance   float64
	CashInventory map[float64]int
	mu            sync.Mutex
}

func NewCashDispenserV1(currBalance float64, cashInventory map[float64]int) (*CashDispenserV1, error) {
	return &CashDispenserV1{
		CurrBalance:   currBalance,
		CashInventory: cashInventory,
	}, nil
}

func (d *CashDispenserV1) GetCurrentBalance() float64 {
	return d.CurrBalance
}

func (d *CashDispenserV1) Deposit(denoms map[float64]int) error {
	total := 0.0

	for denom, count := range denoms {
		total += denom * float64(count)
		d.CashInventory[denom] += count
	}

	d.CurrBalance += total
	return nil
}

func (d *CashDispenserV1) Dispense(amount float64) error {

	if d.CurrBalance < amount {
		return ErrInsufficientCash
	}

	remaining := amount
	used := map[float64]int{}
	inv := d.CashInventory

	// Sort denominations
	var keys []float64
	for denom := range inv {
		keys = append(keys, denom)
	}
	sort.Sort(sort.Reverse(sort.Float64Slice(keys)))

	// Greedy dispensing
	for _, denom := range keys {
		count := inv[denom]
		if count <= 0 {
			continue
		}

		need := int(remaining / denom)
		if need > count {
			need = count
		}

		if need > 0 {
			used[denom] = need
			remaining -= float64(need) * denom
		}

		if remaining <= 0 {
			break
		}
	}

	// Check if exact amount dispensed
	if remaining > 0 {
		return ErrCannotDispense
	}

	// Apply changes
	for denom, cnt := range used {
		inv[denom] -= cnt
	}
	d.CurrBalance -= amount

	return nil
}

type ATMService interface {
	Withdraw(cardNumber, pin string, amount float64) error
	Deposit(cardNumber, pin string, amount float64, denominations map[float64]int) error
	CheckBalance(cardNumber, pin string) (float64, error)
}

type ATMServiceV1 struct {
	ATMid                string
	BankService          BankService
	CardService          CardService
	AccountService       AccountService
	TransactionService   TransactionService
	ReceiptService       ReceiptService
	CashDispenserService CashDispenser
}

func NewATMServiceV1(atmId string, txnServ TransactionService, bankServ BankService, acctServ AccountService, cardServ CardService, receiptServ ReceiptService, dispenser CashDispenser,
) (*ATMServiceV1, error) {
	return &ATMServiceV1{
		ATMid:                atmId,
		TransactionService:   txnServ,
		BankService:          bankServ,
		AccountService:       acctServ,
		CardService:          cardServ,
		ReceiptService:       receiptServ,
		CashDispenserService: dispenser,
	}, nil
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

	err = s.CashDispenserService.Dispense(amount)
	if err != nil {
		return err
	}

	err = s.AccountService.DebitAccount(accountId, amount)
	if err != nil {
		return err
	}

	txn, err := s.TransactionService.CreateTransaction(
		accountId,
		s.ATMid, // Current ATM's ID
		amount,
		Withdraw, // Transaction type
	)

	// Generate receipt (optional - don't fail withdrawal if receipt fails)
	s.ReceiptService.GenerateReceipt(txn.Id)

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

	// 4. Calculate total from denominations
	calculatedAmount := 0.0
	for denom, count := range denominations {
		calculatedAmount += denom * float64(count)
	}

	// 5. Validate user-passed amount matches actual cash
	if amount != calculatedAmount {
		return fmt.Errorf("amount mismatch: declared %.2f but counted %.2f", amount, calculatedAmount)
	}

	//add to atm
	err = s.CashDispenserService.Deposit(denominations)
	if err != nil {
		return err
	}

	//only if adding to atm was successful then only add in account
	err = s.AccountService.CreditAccount(accountId, amount)
	if err != nil {
		return err
	}

	_, err = s.TransactionService.CreateTransaction(
		accountId,
		s.ATMid, // Current ATM's ID
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
		s.ATMid,
		0, // No amount for balance inquiry
		BalanceInquiry,
	)

	return account.CurrBalance, nil

}

type ATM struct {
	Id       string
	BankId   string
	Location *Location
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

// State pattern
type OperationType string

const (
	OpWithdraw OperationType = "WITHDRAW"
	OpDeposit  OperationType = "DEPOSIT"
	OpBalance  OperationType = "BALANCE"
)

type ATMState interface {
	InsertCard(ctx *ATMController, card string) error
	EnterPIN(ctx *ATMController, pin string) error
	SelectOperation(ctx *ATMController, op OperationType) error
	EnterAmount(ctx *ATMController, amount float64) error
	Execute(ctx *ATMController) error
	Cancel(ctx *ATMController) error
	EnterDenominations(ctx *ATMController, denominations map[float64]int) error
}

type ATMController struct {
	currentState ATMState

	//session data
	cardNumber    string
	pin           string
	operation     OperationType
	amount        float64
	denominations map[float64]int

	// Your existing service
	atmService ATMService
}

func NewATMController(atmService ATMService) *ATMController {
	return &ATMController{
		currentState: &IdleState{},
		atmService:   atmService,
	}
}
func (ctx *ATMController) InsertCard(cardNumber string) error {
	return ctx.currentState.InsertCard(ctx, cardNumber)
}

func (ctx *ATMController) EnterPIN(pin string) error {
	return ctx.currentState.EnterPIN(ctx, pin)
}

func (ctx *ATMController) SelectOperation(op OperationType) error {
	return ctx.currentState.SelectOperation(ctx, op)
}
func (ctx *ATMController) EnterAmount(amount float64) error {
	return ctx.currentState.EnterAmount(ctx, amount)
}

func (ctx *ATMController) Execute() error {
	return ctx.currentState.Execute(ctx)
}
func (ctx *ATMController) Cancel() error {
	return ctx.currentState.Cancel(ctx)
}
func (ctx *ATMController) EnterDenominations(denominations map[float64]int) error {
	return ctx.currentState.EnterDenominations(ctx, denominations)
}

func (ctx *ATMController) reset() {
	ctx.cardNumber = ""
	ctx.pin = ""
	ctx.operation = ""
	ctx.amount = 0
	ctx.denominations = nil
}

// Step 4: Base State (Default Implementations)
type BaseATMState struct{}

func (s *BaseATMState) InsertCard(ctx *ATMController, cardNumber string) error {
	return errors.New("❌ cannot insert card in this state")
}

func (s *BaseATMState) EnterPIN(ctx *ATMController, pin string) error {
	return errors.New("❌ cannot enter PIN in this state")
}

func (s *BaseATMState) SelectOperation(ctx *ATMController, op OperationType) error {
	return errors.New("❌ cannot select operation in this state")
}

func (s *BaseATMState) EnterAmount(ctx *ATMController, amount float64) error {
	return errors.New("❌ cannot enter amount in this state")
}

func (s *BaseATMState) Execute(ctx *ATMController) error {
	return errors.New("❌ cannot execute in this state")
}

func (s *BaseATMState) EnterDenominations(ctx *ATMController, denominations map[float64]int) error {
	return errors.New("❌ cannot enter denominations in this state")
}

func (s *BaseATMState) Cancel(ctx *ATMController) error {
	fmt.Println("❌ Transaction cancelled")
	fmt.Println("💳 Card ejected")
	ctx.reset()
	ctx.currentState = &IdleState{}
	return nil
}

//Concrete states

// 1. Idle state
type IdleState struct {
	BaseATMState
}

func (s *IdleState) InsertCard(ctx *ATMController, cardNumber string) error {
	fmt.Println("Card inserted:", cardNumber)
	ctx.cardNumber = cardNumber
	ctx.currentState = &CardInsertState{}
	fmt.Println("📌 Please enter your PIN")
	return nil
}
func (s *IdleState) Cancel(ctx *ATMController) error {
	fmt.Println("Nothing to cancel")
	return nil
}

// 2. Card Inserted State
type CardInsertState struct {
	BaseATMState
}

func (s *CardInsertState) EnterPIN(ctx *ATMController, pin string) error {
	fmt.Println("🔐 Validating PIN...")
	ctx.pin = pin
	ctx.currentState = &PINValidatedState{}
	fmt.Println("✅ PIN validated!")
	fmt.Println("📌 Select operation:")
	fmt.Println("   1. Withdraw")
	fmt.Println("   2. Deposit")
	fmt.Println("   3. Balance Inquiry")
	return nil
}

// 3. PIN Validated State
type PINValidatedState struct {
	BaseATMState
}

func (s *PINValidatedState) SelectOperation(ctx *ATMController, op OperationType) error {
	fmt.Printf("✅ Operation selected: %s\n", op)
	ctx.operation = op
	if op == OpDeposit {
		ctx.currentState = &DenomiantionAndAmountEntryState{}
		fmt.Println("Please insert cash into the deposit slot")
	} else if op == OpBalance {
		ctx.currentState = &ReadyToExecuteState{}
		fmt.Println("📌 Press Execute to confirm")
	} else {
		ctx.currentState = &AmountEntryState{}
		fmt.Println("📌 Enter amount:")
	}
	return nil
}

type DenomiantionAndAmountEntryState struct {
	BaseATMState
}

// assuming we get the denominations map from ATM's hardware
// which is actually counting the notes of each type,
// preparing this map and sending it to our code
func (s *DenomiantionAndAmountEntryState) EnterDenominations(ctx *ATMController, denominations map[float64]int) error {
	total := 0.0
	for denom, count := range denominations {
		total += denom * float64(count)
	}

	ctx.amount = total
	ctx.denominations = denominations
	ctx.currentState = &ReadyToExecuteState{}
	fmt.Printf("✅ Cash counted: ₹%.2f\n", total)
	fmt.Println("   Denominations detected:")
	for denom, count := range denominations {
		fmt.Printf("   ₹%.0f × %d = ₹%.0f\n", denom, count, denom*float64(count))
	}
	fmt.Println("📌 Press Execute to confirm deposit")

	return nil
}

// 4. Operation Selected State
type AmountEntryState struct {
	BaseATMState
}

func (s *AmountEntryState) EnterAmount(ctx *ATMController, amount float64) error {
	fmt.Printf("✅ Amount entered: ₹%.2f\n", amount)
	ctx.amount = amount
	ctx.currentState = &ReadyToExecuteState{}
	fmt.Printf("📌 Confirm %s of ₹%.2f? Press Execute\n", ctx.operation, ctx.amount)
	return nil
}

// 5. Ready To Execute State
type ReadyToExecuteState struct {
	BaseATMState
}

func (s *ReadyToExecuteState) Execute(ctx *ATMController) error {
	fmt.Println("⏳ Processing transaction...")

	var err error

	switch ctx.operation {
	case OpWithdraw:
		err = ctx.atmService.Withdraw(ctx.cardNumber, ctx.pin, ctx.amount)
		if err == nil {
			fmt.Printf("✅ Withdrawal successful! ₹%.2f dispensed\n", ctx.amount)
		}

	case OpDeposit:
		err = ctx.atmService.Deposit(ctx.cardNumber, ctx.pin, ctx.amount, ctx.denominations)
		if err == nil {
			fmt.Printf("✅ Deposit successful! ₹%.2f deposited\n", ctx.amount)
		}

	case OpBalance:
		balance, balErr := ctx.atmService.CheckBalance(ctx.cardNumber, ctx.pin)
		if balErr != nil {
			err = balErr
		} else {
			fmt.Printf("✅ Your current balance: ₹%.2f\n", balance)
		}
	}

	if err != nil {
		fmt.Println("❌ Transaction failed:", err)
		ctx.reset()
		ctx.currentState = &IdleState{}
		return err
	}

	fmt.Println("💳 Please take your card")
	ctx.reset()
	ctx.currentState = &IdleState{}
	return nil
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
	printHeader("ATM SYSTEM - COMPLETE DEMONSTRATION")

	// ============================================
	// SETUP: Initialize All Services
	// ============================================
	bankServ := &BankServiceV1{}
	cardServ, _ := NewCardServiceV1()
	accountServ, _ := NewAccountServiceV1()
	transactionServ, _ := NewTransactionServiceV1()
	receiptServ, _ := NewReceiptServiceV1(transactionServ, accountServ)

	// Create ATM with cash inventory
	cashInv := map[float64]int{
		500: 10, // ₹5,000
		200: 20, // ₹4,000
		100: 30, // ₹3,000
	}

	atm := &ATM{
		Id:       "ATM-001",
		BankId:   "BANK-001",
		Location: &Location{City: "Hyderabad", Street: "Lanco Hills", Pincode: "500089"},
	}
	dispenserServ, _ := NewCashDispenserV1(12000, cashInv)
	atmServ, _ := NewATMServiceV1(atm.Id, transactionServ, bankServ, accountServ, cardServ, receiptServ, dispenserServ)

	// Create test data
	testAccount := &Account{
		Id:          "ACC-001",
		UserId:      "USER-001",
		BankId:      "BANK-001",
		CurrBalance: 50000,
		AccountType: Savings,
		DailyLimit:  20000,
	}
	accountServ.accounts["ACC-001"] = testAccount

	testCard := &Card{
		CardNumber: "CARD-001",
		UserId:     "USER-001",
		AccountId:  "ACC-001",
		Name:       "John Doe",
		ExpiryDate: time.Now().AddDate(2, 0, 0),
		Status:     Active,
	}
	cardServ.Cards["CARD-001"] = testCard

	fmt.Println("\n✅ System Initialized")
	fmt.Printf("   💳 Card: %s (John Doe)\n", testCard.CardNumber)
	fmt.Printf("   🏦 Account: %s\n", testAccount.Id)
	fmt.Printf("   💰 Initial Balance: ₹%.2f\n", testAccount.CurrBalance)

	// ============================================
	// PART 1: Direct Service Calls (Backend API)
	// ============================================
	printSectionHeader("PART 1: Direct Service Calls (Backend API Style)")

	// Test 1: Deposit
	fmt.Println("📥 Test 1: Deposit ₹6,000")
	depositDenom := map[float64]int{500: 2, 200: 10, 100: 30}
	if err := atmServ.Deposit("CARD-001", "1234", 6000, depositDenom); err != nil {
		fmt.Println("   ❌ Error:", err)
	} else {
		fmt.Printf("   ✅ Success! Balance: ₹%.2f\n", testAccount.CurrBalance)
	}
	fmt.Println()

	// Test 2: Withdraw
	fmt.Println("📤 Test 2: Withdraw ₹5,000")
	if err := atmServ.Withdraw("CARD-001", "1234", 5000); err != nil {
		fmt.Println("   ❌ Error:", err)
	} else {
		fmt.Printf("   ✅ Success! Balance: ₹%.2f\n", testAccount.CurrBalance)
	}
	fmt.Println()

	// Test 3: Check Balance
	fmt.Println("💵 Test 3: Check Balance")
	if balance, err := atmServ.CheckBalance("CARD-001", "1234"); err != nil {
		fmt.Println("   ❌ Error:", err)
	} else {
		fmt.Printf("   ✅ Current Balance: ₹%.2f\n", balance)
	}
	fmt.Println()

	// Test 4: Transaction History
	fmt.Println("📜 Test 4: Transaction History")
	if txns, err := transactionServ.GetTransactionHistory("ACC-001"); err != nil {
		fmt.Println("   ❌ Error:", err)
	} else {
		fmt.Printf("   ✅ Total Transactions: %d\n", len(txns))
		for i, txn := range txns {
			fmt.Printf("      %d. %s | ₹%.2f | %s\n", i+1, txn.Type, txn.Amount, txn.Status)
		}
	}
	fmt.Println()

	// ============================================
	// PART 2: State Pattern (Physical ATM)
	// ============================================
	printSectionHeader("PART 2: State Pattern (Physical ATM Machine)")

	fmt.Println("🏧 ATM Machine Ready\n")

	// ==========================================
	// Scenario 1: Withdraw
	// ==========================================
	printScenarioHeader("SCENARIO 1: Withdraw ₹2,000")

	atmController := NewATMController(atmServ)

	fmt.Println("▶️  Step 1: Insert Card")
	atmController.InsertCard("CARD-001")
	fmt.Println()

	fmt.Println("▶️  Step 2: Enter PIN")
	atmController.EnterPIN("1234")
	fmt.Println()

	fmt.Println("▶️  Step 3: Select Withdraw")
	atmController.SelectOperation(OpWithdraw)
	fmt.Println()

	fmt.Println("▶️  Step 4: Enter Amount")
	atmController.EnterAmount(2000)
	fmt.Println()

	fmt.Println("▶️  Step 5: Execute Transaction")
	atmController.Execute()
	fmt.Printf("   📊 Balance After: ₹%.2f\n\n", testAccount.CurrBalance)

	// ==========================================
	// Scenario 2: Deposit
	// ==========================================
	printScenarioHeader("SCENARIO 2: Deposit ₹3,000")

	atmController2 := NewATMController(atmServ)

	fmt.Println("▶️  Step 1: Insert Card")
	atmController2.InsertCard("CARD-001")
	fmt.Println()

	fmt.Println("▶️  Step 2: Enter PIN")
	atmController2.EnterPIN("1234")
	fmt.Println()

	fmt.Println("▶️  Step 3: Select Deposit")
	atmController2.SelectOperation(OpDeposit)
	fmt.Println()

	fmt.Println("▶️  Step 4: Insert Cash (Hardware Simulation)")
	fmt.Println("   [User inserts physical notes into deposit slot]")
	fmt.Println("   [ATM Hardware: Scanning notes...]")
	fmt.Println("   [ATM Hardware: Counting denominations...]")

	hardwareDetected := map[float64]int{500: 6}
	atmController2.EnterDenominations(hardwareDetected)
	fmt.Println()

	fmt.Println("▶️  Step 5: Execute Transaction")
	atmController2.Execute()
	fmt.Printf("   📊 Balance After: ₹%.2f\n\n", testAccount.CurrBalance)

	// ==========================================
	// Scenario 3: Balance Inquiry
	// ==========================================
	printScenarioHeader("SCENARIO 3: Balance Inquiry")

	atmController3 := NewATMController(atmServ)

	fmt.Println("▶️  Step 1: Insert Card")
	atmController3.InsertCard("CARD-001")
	fmt.Println()

	fmt.Println("▶️  Step 2: Enter PIN")
	atmController3.EnterPIN("1234")
	fmt.Println()

	fmt.Println("▶️  Step 3: Select Balance Inquiry")
	atmController3.SelectOperation(OpBalance)
	fmt.Println()

	fmt.Println("▶️  Step 4: Execute")
	atmController3.Execute()
	fmt.Println()

	// ==========================================
	// Scenario 4: State Validation
	// ==========================================
	printScenarioHeader("SCENARIO 4: State Validation Test")

	atmController4 := NewATMController(atmServ)

	fmt.Println("❌ Test 1: Enter PIN without card")
	if err := atmController4.EnterPIN("1234"); err != nil {
		fmt.Println("   ✅ Correctly rejected:", err)
	}
	fmt.Println()

	fmt.Println("❌ Test 2: Execute without any input")
	atmController5 := NewATMController(atmServ)
	if err := atmController5.Execute(); err != nil {
		fmt.Println("   ✅ Correctly rejected:", err)
	}
	fmt.Println()

	fmt.Println("❌ Test 3: Cancel in Idle state")
	atmController6 := NewATMController(atmServ)
	atmController6.Cancel()
	fmt.Println()

	// ==========================================
	// Scenario 5: Multi-step with Cancel
	// ==========================================
	printScenarioHeader("SCENARIO 5: Cancel Transaction Mid-Flow")

	atmController7 := NewATMController(atmServ)

	fmt.Println("▶️  User inserts card and enters PIN")
	atmController7.InsertCard("CARD-001")
	atmController7.EnterPIN("1234")
	fmt.Println()

	fmt.Println("▶️  User selects Withdraw")
	atmController7.SelectOperation(OpWithdraw)
	fmt.Println()

	fmt.Println("▶️  User changes mind and cancels")
	atmController7.Cancel()
	fmt.Println()

	// ==========================================
	// Final Summary
	// ==========================================
	printSectionHeader("FINAL SUMMARY")

	fmt.Printf("💰 Final Account Balance: ₹%.2f\n", testAccount.CurrBalance)

	if txns, _ := transactionServ.GetTransactionHistory("ACC-001"); txns != nil {
		fmt.Printf("📊 Total Transactions Processed: %d\n", len(txns))

		var totalWithdrawn, totalDeposited float64
		for _, txn := range txns {
			if txn.Type == Withdraw {
				totalWithdrawn += txn.Amount
			} else if txn.Type == Deposit {
				totalDeposited += txn.Amount
			}
		}
		fmt.Printf("📤 Total Withdrawn: ₹%.2f\n", totalWithdrawn)
		fmt.Printf("📥 Total Deposited: ₹%.2f\n", totalDeposited)
	}

	fmt.Println("\n" + strings.Repeat("═", 62))
	fmt.Println("✅ ATM SYSTEM DEMONSTRATION COMPLETE")
	fmt.Println(strings.Repeat("═", 62))
}

// ==========================================
// Helper Functions
// ==========================================

func printHeader(title string) {
	fmt.Println("\n" + strings.Repeat("═", 62))
	fmt.Println(centerText(title, 62))
	fmt.Println(strings.Repeat("═", 62))
}

func printSectionHeader(title string) {
	fmt.Println("\n" + strings.Repeat("─", 62))
	fmt.Println(centerText(title, 62))
	fmt.Println(strings.Repeat("─", 62) + "\n")
}

func printScenarioHeader(title string) {
	fmt.Println("┌" + strings.Repeat("─", 60) + "┐")
	fmt.Println("│" + centerText(title, 60) + "│")
	fmt.Println("└" + strings.Repeat("─", 60) + "┘\n")
}

func centerText(text string, width int) string {
	if len(text) >= width {
		return text[:width]
	}
	padding := (width - len(text)) / 2
	return strings.Repeat(" ", padding) + text + strings.Repeat(" ", width-padding-len(text))
}
