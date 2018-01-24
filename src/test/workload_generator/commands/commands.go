package commands

import (
    "strconv"
    "log"
)

type Command struct {
    Name string     `json:"name,omitempty"`
    Username string `json:"username,omitempty"`
    Symbol string   `json:"symbol,omitempty"`
    Amount int      `json:"amount,omitempty"`
    Filename string `json:"filename,omitempty"`
}

func parseCommandUserSymbolAmount(data []string) Command {
	amount, err := strconv.Atoi(data[4])
    if err != nil{
        log.Fatalf("Could not parse Amount: %s \n %s" , data[4], data)
    }
    return Command{Name: data[1], Username: data[2], Symbol: data[3], Amount: amount}
}

func parseCommandUserAmount(data []string) Command {
	amount, err := strconv.Atoi(data[3])
    if err != nil{
        log.Fatalf("Could not parse Amount: %s \n %s" , data[3], data)
    }
    return Command{Name: data[1], Username: data[2], Amount: amount}
}

func parseCommandUserSymbol(data []string) Command {
    return Command{Name: data[1], Username: data[2], Symbol: data[3]}
}
 
func parseCommandUser(data []string) Command {
    return Command{Name: data[1], Username: data[2]}
}

func ParseAdd(data []string) Command {
   	return parseCommandUserAmount(data)
}

func ParseQuote(data []string) Command {
    return parseCommandUserSymbol(data)
}

func ParseBuy(data []string) Command {
   	return parseCommandUserSymbolAmount(data)
}

func ParseSell(data []string) Command {
   	return parseCommandUserSymbolAmount(data)
}

func ParseCommitBuy(data []string) Command {
  	return parseCommandUser(data)
}

func ParseCommitSell(data []string) Command {
   	return parseCommandUser(data)
}

func ParseCancelSell(data []string) Command {
	return parseCommandUser(data)
}

func ParseCancelBuy(data []string) Command {
	return parseCommandUser(data)
}

func ParseDisplaySummary(data []string) Command {
	return parseCommandUser(data)
}

func ParseSetSellTrigger(data []string) Command {
   	return parseCommandUserSymbolAmount(data)
}

func ParseSetBuyTrigger(data []string) Command {
   	return parseCommandUserSymbolAmount(data)
}

func ParseCancelSetSell(data []string) Command {
    return parseCommandUserSymbol(data)
}

func ParseCancelSetBuy(data []string) Command {
    return parseCommandUserSymbol(data)
}


func ParseSetSellAmount(data []string) Command {
    return parseCommandUserSymbolAmount(data)
}

func ParseSetBuyAmount(data []string) Command {
    return parseCommandUserSymbolAmount(data)
}

func ParseDumplog(data []string) Command {
	return Command{Name: data[1], Filename: data[2]}
}

func ParseData(data []string) Command {
    if len(data) < 2{
        log.Fatal("Length of entry too short")
    }
    var command Command

    switch cmdName := data[1]; cmdName {
        case "ADD":
            command = ParseAdd(data)

        case "QUOTE":
            command = ParseQuote(data)

        case "BUY":
            command = ParseBuy(data)

        case "SELL":
            command = ParseSell(data)

        case "COMMIT_BUY":
            command = ParseCommitBuy(data)

        case "COMMIT_SELL":
            command = ParseCommitSell(data)

        case "CANCEL_BUY":
            command = ParseCommitBuy(data)

        case "CANCEL_SELL":
            command = ParseCommitSell(data)

        case "DISPLAY_SUMMARY":
            command = ParseCommitSell(data)

        case "SET_BUY_TRIGGER":
            command = ParseSetBuyTrigger(data)

        case "SET_SELL_TRIGGER":
            command = ParseSetSellTrigger(data)

        case "CANCEL_SET_BUY":
            command = ParseCancelSetBuy(data)

        case "CANCEL_SET_SELL":
            command = ParseCancelSetSell(data)

        case "SET_BUY_AMOUNT":
            command = ParseCancelSetBuy(data)

        case "SET_SELL_AMOUNT":
            command = ParseCancelSetSell(data)

        case "DUMPLOG":
            command = ParseDumplog(data)

        default:
            log.Fatalf("Invalid command: %s", data[1])
    }
    return command
}