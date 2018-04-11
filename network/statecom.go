package network

import (//"bytes"
	"encoding/json"
	"fmt"
	"net"
	"time"
	"../state"
	"log"
)

var remote_elev1_alive bool = false
var remote_elev2_alive bool = false

const stasjon13 string = "129.241.187.152:10001"
const stasjon14 string = "129.241.187.142:10001"
const stasjon17 string = "129.241.187.145:10001"
const stasjon20 string = "129.241.187.155:10001"
const stasjon22 string = "129.241.187.56:10001"
const stasjon23 string = "129.241.187.57:10001"
const stasjon10 string = "129.241.187.158:10001"
const stasjon11 string = "129.241.187.159:10001"


const targetIP string = stasjon11
const (
	REMOTE_1   int	= 1
	REMOTE_2	= 2
)

//const targetIP string = stasjon11
const remote_elev_IP1 = stasjon 10
const remote_elev_IP2 = stasjon 11


//TODO: make remote button event listener

func Broadcast_state(bcast <- chan state.Elevator) {
	localip := get_localip()

	local_addr, err := net.ResolveUDPAddr("udp", localip + ":0")
	state.Check(err)
	target_addr,err := net.ResolveUDPAddr("udp", remote_elev_IP1)
	state.Check(err)
	target_addr2,err := net.ResolveUDPAddr("udp", remote_elev_IP2)
	state.Check(err)
	out_connection, err := net.DialUDP("udp", local_addr, target_addr)
	state.Check(err)
	out_connection2, err := net.DialUDP("udp", local_addr, target_addr)
	state.Check(err)
	defer out_connection.Close()
	defer out_connection2.Close()
	
	for {
		select {
		case data := <- bcast:
			send_state(data, out_connection)
			send_state(data, out_connection2)
		}
	}
}


func Poll_remote_state(output *state.Elevator) {
	localip := get_localip()

	listen_addr, err := net.ResolveUDPAddr("udp", localip + ":10001")
	state.Check(err)
	input, _ := net.ListenUDP("udp", listen_addr)
	state.Check(err)
	defer input.Close()
	
	wd_reset := make(chan bool)

	for {
		if (is_alive() == false) {
			go watchdog(wd_reset)
			fmt.Println("Connection established!")
		}
		wd_reset <- true
		
		*output = read_state(input)
		fmt.Println("Received: ", *output)
	}
}

func Poll_remote_state2(output *state.Elevator) {
	localip := get_localip()

	listen_addr, err := net.ResolveUDPAddr("udp", localip + ":10002")
	state.Check(err)
	input, _ := net.ListenUDP("udp", listen_addr)
	state.Check(err)
	defer input.Close()
	
	wd_reset := make(chan bool)

	for {
		if (is_alive() == false) {
			go watchdog(wd_reset)
			fmt.Println("Connection established!")
		}
		wd_reset <- true
		
		*output = read_state(input)
		fmt.Println("Received: ", *output)
	}
}
//**unfinished**
func Broadcast_order() {

	//broadcast order to other elevators
	localip := get_localip()

	local_addr, err := net.ResolveUDPAddr("udp", localip + ":0")
	state.Check(err)
	target_addr,err := net.ResolveUDPAddr("udp", remote_elev_IP1)
	state.Check(err)
	target_addr2,err := net.ResolveUDPAddr("udp", remote_elev_IP2)
	state.Check(err)
	out_connection, err := net.DialUDP("udp", local_addr, target_addr)
	state.Check(err)
	out_connection2, err := net.DialUDP("udp", local_addr, target_addr)
	state.Check(err)
	local_connection, err := net.DialUDP("udp", local_addr, localip + ":10003")
	state.Check(err)
	defer out_connection.Close()
	defer out_connection2.Close()
	defer local_connection.Close()
	
	for {
		select {
		case data := <- bcast: //change to something free
			send_order(data, out_connection) //send_order needs to be made
			send_order(data, out_connection2)
			send_order(data, local_connection)
		}
	}

}
//**unfinished**
func Poll_order()

	//receive order
	//Triggers on receiving order from port
	//polls order and sends order through channel to evaluate function
	network.Broadcast_state()
	network.Poll remote_state()
	network.Poll_remote_state2()
	<- states_complete
				
	cost_local := driver.Evaluate(e, order)
	cost_e2 := driver.Evaluate(e2, order)
	cost_e3 := driver.Evaluate(e3, order)
	Choose_elevator(cost_local, cost_e2. cost_e3)
	<- order_taken
	//Order_accept()

}

func send_state(state state.Elevator, connection *net.UDPConn) {
	jsonRequest, err := json.Marshal(state)
		if err != nil {
			log.Print("Marshal Register information failed.")
			log.Fatal(err)
		}
		connection.Write(jsonRequest)
	fmt.Println("Broadcasting: ", state)
	
}

func read_state(connection *net.UDPConn) state.Elevator {
	var message state.Elevator
	inputBytes := make([]byte, 4096)
	fmt.Println("Starts listening....")
	length, _, _ := connection.ReadFromUDP(inputBytes)
                err := json.Unmarshal(inputBytes[:length], &message)
		if err != nil {
			log.Print(err)
			//continue
		}
	
	return message
}

func watchdog(wd_reset <- chan bool) {
	//fmt.Println("Watchdog activated!\n")
	set_alive(true)
	for i := 0; i < 10; i++ {
		time.Sleep(500*time.Millisecond)
		select {
		case <- wd_reset:
			i = 0

		default:
		}
	}
	set_alive(false)
	fmt.Println("Connection lost with elevator 2.")
}


func get_localip() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	state.Check(err)
	defer conn.Close()

	ip_with_port := conn.LocalAddr().String()

	var ip string = ""
	for _, char := range ip_with_port {
		if (char == ':') {
			break
		}
		ip += string(char)
	}
	return ip
}

func set_alive(b bool) {
	_alive = b
}

func is_alive() bool {
	return _alive
}
