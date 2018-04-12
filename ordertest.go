package main

import (
	"fmt"
	//"net"
	"time"
	"math/rand"
	//"encoding/json"
	//"log"
)

type order struct {
	Name string
	Stamp time.Time
	ID int
}
//order must be placed in q-slice
//order must be removed from q-slice when order is complete

func order_queue(add_order <-chan order, remove_order <-chan order) {

	var q []order
	
	for {
		timecheck_order_queue(q)
		select {
		case newQ := <- add_order:
			q = append(q, newQ)
			fmt.Println("added",newQ.Name)
			
		case removeQ := <- remove_order:
			i := 0
			for _,c := range q {
				if c.ID == removeQ.ID {
					fmt.Println("removing", c.Name)
					q = q[:i+copy(q[i:], q[i+1:])]
					fmt.Println(q)
				}
				i++
			}			
		}
	}
}

func timecheck_order_queue(q []order) {
	for _, c := range q {
		if time.Now().Sub(c.Stamp) > 2*time.Second {
			fmt.Println(c.Name," failed")
			//TODO: Force other elevators to take order
		}
	}
}


func main() {

	add_order := make(chan order)
	remove_order := make(chan order)
	exit := make(chan bool)
	
	go order_queue(add_order, remove_order)
	
	
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	
	//	l1 := order{Name: "Jonas1", Stamp: time.Now(), ID: r1.Intn(1000)}
	
	l1 := order{Name: "Jonas1", Stamp: time.Now(), ID: r1.Intn(1000)}
	add_order <- l1
	time.Sleep(1*time.Second)
	l2 := order{Name: "Jonas2", Stamp: time.Now(), ID: r1.Intn(1000)}
	add_order <- l2
	time.Sleep(1*time.Second)
	l3 := order{Name: "Jonas3", Stamp: time.Now(), ID: r1.Intn(1000)}
	add_order <- l3
	time.Sleep(1*time.Second)
	l4 := order{Name: "Jonas4", Stamp: time.Now(), ID: r1.Intn(1000)}
	add_order <- l4
	remove_order <- l2

	
	<- exit
}