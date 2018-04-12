package network

import (
	"fmt"
	"time"
	"math/rand"

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

