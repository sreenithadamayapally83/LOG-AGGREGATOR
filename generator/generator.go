package main

import (
	"fmt"
	"math/rand"
	"os"
)

func main() {
	file, err := os.Create("transactions.log")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	statuses := []string{"SUCCESS", "FAILURE", "PENDING", "REJECTED"}
	fmt.Println("Generating 100,000 logs...")

	for i := 0; i < 100000; i++ {
		line := fmt.Sprintf("TXNID:%d | STATUS:%s | LATENCY:%dms\n",
			rand.Intn(1000000), statuses[rand.Intn(4)], rand.Intn(200))
		file.WriteString(line)
	}
	fmt.Println("Done! 'transactions.log' is ready.")
}
