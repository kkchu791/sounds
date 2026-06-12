package main

import (
	"github.com/kkchu791/sounds/internal/app/broker"
)

func main() {
	broker.Run()
}

// func handleError(err error) {
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// }

// p := model.NewPartition()
// b := model.NewBroker("broker-0", p)
// partition_key := "outdoor"
// m := model.NewMessage("chirp", partition_key)
// b.Append(m)
// m = model.NewMessage("cool", partition_key)
// b.Append(m)
// m = model.NewMessage("bom", partition_key)
// b.Append(m)

// m, err := b.Read(0)
// fmt.Println(m.Sound, m.Key, m.Timestamp)
// m, err = b.Read(1)
// fmt.Println(m.Sound, m.Key, m.Timestamp)

// m, err = b.Read(2)
// fmt.Println(m.Sound, m.Key, m.Timestamp)

// handleError(err)
