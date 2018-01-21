package commands

import (
    "log"
    "fmt"
)

func FormatCommandEndpoint(cmd Command) string {
    switch cmd.Name {
        case "ADD":
           return fmt.Sprintf("/api/addUser/%s/%f", cmd.Username, cmd.Amount)

        case "QUOTE":
			return fmt.Sprintf("/api/getQuote/%s/%s", cmd.Username, cmd.Symbol)

        case "BUY":
			return fmt.Sprintf("/api/buyOrder/%s/%s/%f", cmd.Username, cmd.Symbol, cmd.Amount)

        case "SELL":
			return fmt.Sprintf("/api/sellOrder/%s/%s/%f", cmd.Username, cmd.Symbol, cmd.Amount)

        case "COMMIT_BUY":
			return fmt.Sprintf("/api/commitBuy/%s", cmd.Username)    

        case "COMMIT_SELL":
			return fmt.Sprintf("/api/commitSell/%s", cmd.Username)

        case "CANCEL_BUY":
            return ""

        case "CANCEL_SELL":
            return ""

        case "DISPLAY_SUMMARY":
            return ""
            
        case "SET_BUY_TRIGGER":
            return ""
            
        case "SET_SELL_TRIGGER":
            return ""

        case "CANCEL_SET_BUY":
            return ""

        case "CANCEL_SET_SELL":
            return ""

        case "SET_BUY_AMOUNT":
            return ""

        case "SET_SELL_AMOUNT":
            return ""

        case "DUMPLOG":
            return ""

        default:
            log.Fatalf("Invalid command: %s", cmd.Name)
    }
    
    return ""
}