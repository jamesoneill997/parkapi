package logs

import (
	"fmt"
	"os"
	"time"
)

//LogError function takes an error as a parameter and writes the error to the log file that is contained within this directory
func LogError(err error) {
	//open file with write only permissions
	f, er := os.OpenFile("./logs/log.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

	//check err
	if er != nil {
		fmt.Println(er)
	}

	//defer closing file until end of function
	defer f.Close()

	//check err
	if _, e := f.WriteString(fmt.Sprintf("%s Error: %s\n", time.Now(), err.Error())); err != nil {
		fmt.Println(e)
	}
}
