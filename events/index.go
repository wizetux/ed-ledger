package events

import (
  "log"
  "time"
  "ed-ledger/balancesheet"
  "fmt"
)

type Event map[string]interface{}

var expenses map[string]string = map[string]string{"RepairAll":"Cost", "PayFines":"Amount", "RefuelAll":"Cost", "BuyExplorationData":"Cost", "BuyTradeData":"Cost", "MarketBuy":"TotalCost", "BuyAmmo":"Cost", "BuyDrones":"TotalCost", "FetchRemoteModule":"TransferCost", "ModuleBuy":"BuyPrice", "RefuelPartial":"Cost", "Repair":"Cost", "RestockVehicle":"Cost", "ShipyardBuy":"ShipPrice", "ShipyardTransfer":"TransferPrice"}

var income map[string]string = map[string]string{"MissionCompleted":"Reward", "RedeemVoucher":"Amount", "MarketSell":"TotalSale", "ModuleBuy":"SellPrice", "ModuleSell":"SellPrice", "ModuleSellRemote":"SellPrice", "SellDrones":"TotalSale", "ShipyardBuy":"SellPrice", "ShipyardSell":"ShipPrice"}

var ProcessedEvents = 0

func ProccessEvent(balanceSheet *balancesheet.BalanceSheet, event Event)  {
  if balanceSheet.LastEventTime != "" {
    eventTimeStamp, err := time.Parse(time.RFC3339, event["timestamp"].(string))

    if err != nil {
      log.Fatal(err)
    }

    lastEventTimeStamp, err := time.Parse(time.RFC3339, balanceSheet.LastEventTime)
    if err != nil {
      log.Fatal(err)
    }

    if lastEventTimeStamp.After(eventTimeStamp) {
      return
    }
  }

  field, ok := expenses[event["event"].(string)]
  if ok {
    value, ok := event[field].(float64) 
    if ok {
      value := int(value)
      balanceSheet.AddExpense(value)
      ProcessedEvents++
    }
  }

  field, ok = income[event["event"].(string)]
  if ok {
    value, ok := event[field].(float64)
    if ok {
      value := int(value)
      balanceSheet.AddIncome(value)
      ProcessedEvents++
    }
  }

  if event["event"] == "SellOrganicData" {
    for _, bioData := range event["BioData"].([]interface{}) {
      obj := bioData.(map[string]interface{})
      balanceSheet.AddIncome(int(obj["Value"].(float64)))
      ProcessedEvents++
    }
  }

  if event["event"] == "SellExplorationData" {
    baseValue := int(event["BaseValue"].(float64))
    bonusValue := int(event["Bonus"].(float64))
    balanceSheet.AddIncome(baseValue + bonusValue)
    fmt.Printf("SellExplorationData BaseValue: %d  Bonus: %d", baseValue, bonusValue)
    ProcessedEvents++
  }

  balanceSheet.LastEventTime = event["timestamp"].(string)
  return
}
