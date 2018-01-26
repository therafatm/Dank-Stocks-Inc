package commands

import (
	"fmt"
	"log"
)

func FormatCommandEndpoint(cmd Command) string {
	switch cmd.Name {
	case "BALANCE":
		return fmt.Sprintf("/api/availableBalance/%s/%s", cmd.Username, cmd.Tnum)

	case "SHARES":
		return fmt.Sprintf("/api/availableShares/%s/%s/%s", cmd.Username, cmd.Symbol, cmd.Tnum)

	case "ADD":
		return fmt.Sprintf("/api/add/%s/%d/%s", cmd.Username, cmd.Amount, cmd.Tnum)

	case "QUOTE":
		return fmt.Sprintf("/api/getQuote/%s/%s/%s", cmd.Username, cmd.Symbol, cmd.Tnum)

	case "BUY":
		return fmt.Sprintf("/api/buy/%s/%s/%d/%s", cmd.Username, cmd.Symbol, cmd.Amount, cmd.Tnum)

	case "SELL":
		return fmt.Sprintf("/api/sell/%s/%s/%d/%s", cmd.Username, cmd.Symbol, cmd.Amount, cmd.Tnum)

	case "COMMIT_BUY":
		return fmt.Sprintf("/api/commitBuy/%s/%s", cmd.Username, cmd.Tnum)

	case "COMMIT_SELL":
		return fmt.Sprintf("/api/commitSell/%s/%s", cmd.Username, cmd.Tnum)

	case "CANCEL_BUY":
		return fmt.Sprintf("/api/cancelBuy/%s/%s", cmd.Username, cmd.Tnum)

	case "CANCEL_SELL":
		return fmt.Sprintf("/api/cancelSell/%s/%s", cmd.Username, cmd.Tnum)

	case "DISPLAY_SUMMARY":
		return ""

	case "SET_BUY_TRIGGER":
		return fmt.Sprintf("/api/setBuyTrigger/%s/%s/%d/%s", cmd.Username, cmd.Symbol, cmd.Amount, cmd.Tnum)

	case "SET_SELL_TRIGGER":
		return fmt.Sprintf("/api/setSellTrigger/%s/%s/%d/%s", cmd.Username, cmd.Symbol, cmd.Amount, cmd.Tnum)

	case "CANCEL_SET_BUY":
		return fmt.Sprintf("/api/cancelSetBuy/%s/%s/%s", cmd.Username, cmd.Symbol, cmd.Tnum)

	case "CANCEL_SET_SELL":
		return fmt.Sprintf("/api/cancelSetSell/%s/%s/%s", cmd.Username, cmd.Symbol, cmd.Tnum)

	case "SET_BUY_AMOUNT":
		return fmt.Sprintf("/api/setBuyAmount/%s/%s/%d/%s", cmd.Username, cmd.Symbol, cmd.Amount, cmd.Tnum)

	case "SET_SELL_AMOUNT":
		return fmt.Sprintf("/api/setSellAmount/%s/%s/%d/%s", cmd.Username, cmd.Symbol, cmd.Amount, cmd.Tnum)

	case "EXECUTE_TRIGGERS":
		return fmt.Sprintf("/api/executeTriggers/%s/%s", cmd.Username, cmd.Tnum)

	case "DUMPLOG":
		return ""

	default:
		log.Fatalf("Invalid command: %s", cmd.Name)
	}

	return ""
}
