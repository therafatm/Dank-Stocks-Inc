package commands

import (
    "log"
)


func formatAdd(cmd Command){
	return fmt.Sprintf("/api/addUser/%s/%s", cmd.Username, cmd.Amount)
}

