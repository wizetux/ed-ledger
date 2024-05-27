package main

import (
  "fmt"
  "os"
  "io"
  "log"
  "encoding/json"
  
  "ed-ledger/events"
  "ed-ledger/balancesheet"
)

func main() {
    args := os.Args[1:]
    if len(args) < 1 {
      fmt.Println("Please provide a journal to process")
      return
    }

    logFile, err := os.Open(args[0])
    if err != nil {
      fmt.Println("Error opening the log file")
      return
    }
    defer logFile.Close()

    balanceSheet, err := balancesheet.Load("./balance.json")
    if err != nil {
      log.Fatal(err)
    }

    dec := json.NewDecoder(logFile)
    for {
      var event events.Event
      if err := dec.Decode(&event); err == io.EOF {
        break
      } else if err != nil {
        log.Fatal(err)
      }
      events.ProccessEvent(&balanceSheet, event)
    }
    fmt.Printf("%d events processed\n", events.ProcessedEvents)
    balanceSheet.DebugPrint()
    balanceSheet.Save("./balance.json")
}

