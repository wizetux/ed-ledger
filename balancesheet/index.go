package balancesheet

import (
  "fmt"
  "io"
  "os"
  "encoding/json"
)

type BalanceSheet struct {
  Balance int
  TotalExpenses int
  TotalIncome int
  LastEventTime string
}

func (balancesheet BalanceSheet) Save(filePath string) (err error) {
  byteData, err := json.Marshal(balancesheet)
  if err != nil {
  }

  err = os.WriteFile(filePath, byteData, 0666)
  return err
}

func Load(filePath string) (balanceSheet BalanceSheet, err error) {
  _, err = os.Stat(filePath)

  if err != nil {
    fmt.Printf("File %s did not exist. Returning empty balance sheet\n", filePath)
    return BalanceSheet{
      Balance: 0,
      TotalExpenses: 0,
      TotalIncome: 0,
      LastEventTime: "",
    }, nil
  }

  file, err := os.Open(filePath)
  if err != nil {
    return
  }
  defer file.Close()

  bytes, err := io.ReadAll(file)
  if err != nil {
    return
  }

  json.Unmarshal(bytes, &balanceSheet)

  return
}

func (balanceSheet *BalanceSheet) AddExpense(value int) {
  balanceSheet.TotalExpenses += value
  balanceSheet.Balance -= value
}

func (balanceSheet *BalanceSheet) AddIncome(value int) {
  balanceSheet.TotalIncome += value
  balanceSheet.Balance += value
}

func (balanceSheet *BalanceSheet) DebugPrint() {
  byteData, err := json.MarshalIndent(balanceSheet, "", "   ")
  if err != nil {
  }

  fmt.Println(string(byteData))
}
