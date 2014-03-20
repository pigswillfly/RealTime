package Control

import(
	."fmt"
	."container/list"
	."strconv"
)

const (
)

type Control struct{
	elev *Elevator				// elevator
	elevators *List				// list of other elevators
	backup map[int]*List		// backup list of requests of other elev
	costs map[int]int			// costs of other elevators referenced by id
	
	net *Network				// network
	toCtrlNet chan string		// TCP connection for control messages
	fromCtrlNet chan string
	toAliveNet chan string		// UDP connection for alive messages
	fromAliveNet chan string	
}


func Control(){
	for {

		// ACTIONS WHEN IN RIGHT STATE:
	
		// button push - poll function?
			// add to request list
			// allocate to elevator (wrt direction) - allocate function?
		// alive message
			// confirm in alive list - list management function
			// if not add new
		// dead message
			// remove from alive list - list management function
			// get requests from elevator's request queue and reallocate?
		// reached floor	
			// is it correct? - fulfilled function
			// open/close doors
			// remove request 
		// stop - stop function

		// Process pairs?
		// Primary with backup? - primary dead function
			// primary dead, backup becomes primary
			// new backup


		// states
			// 1. moving up
			// 2. moving down
			// 3. waiting (no requests)
			// 4. opening/closing doors
			// 5. stopped	
	}
}

func (ctrl *Control) Send_Msg(id int, to int, code string, args []string){

	// Make message comma delimited string
	msg := Itoa(id) + "," + Itoa(to) + "," + code
	for i:= range args {
		if args[i] != ""{
			msg += "," + args[i]
		}
	}
	ctrl.toCtrlNet <- msg
}

func (ctrl *Control) Recieve_Msg(){
	for{
		// Message format:
		//	[Sender ID],[Receiver ID],[Code],[Arguments]
		msg := <-fromCtrlNet
		from_id, code, args := Decipher_Msg(msg)
		var send_args []string
		
		switch code{
			case "CostPlease":
				// args -- button, floor
				send_args[0] = Itoa(ctrl.elev.Cost(args[0],args[1])
				Ctrl.Send_Msg(ctrl.elev.id, from_id, "MyCost", send_args)
			case "MyCost":
				// args -- cost
				Ctrl.Add_Cost(id,args[0])
	
			case "MyList":
				// args -- floor requests
				Update_List(from_id, args)

			case "ListPlease":
				// no args
				Send_List(from_id)

			case "AddFloor":
				// args -- floor to add
		}
	}
}

func (ctrl *Control) Decipher_Msg(msg string) (int, string, []string){
	
	substrings := SplitN(msg, ",", -1)		
	from_id,_ := Atoi(substrings[0])
	code := substrings[2]
	args := substrings[3:]
	return id, code, args
}

func (ctrl *Control) Poll_Buttons(){
	
}

func (ctrl *Control) Decide_Elevator(button int, floor int){
	// 
}

func (ctrl *Control) Add_Cost(id int, cost int){

}

func (ctrl *Control) Send_Alive_Msg(){
	for {
		msg := <-ctrl.elev.alive
		ctrl.toAliveNet <-msg
	}
}

func (ctrl *Control) Rec_Alive_Msg(){
	for {
		msg := <-ctrl.fromAliveNet
		id,code,_ := Decipher_Msg(msg)
		if code == "Alive" {
			Update_Elevator_List(id)
		} else {
	}

	// msg := <-ctrl.fromAliveNet
	// get elevator id
	// call Update_List
}

func (ctrl *Control) Update_Elevator_List(id int){
	
}

func (ctrl *Control) Update_List(id int, requests []string){

	if elevators[id] != nil{
		l *List
		l.New()
		req := -1
		for i:= range requests{
			req,_ = Atoi(requests[i])
			if req > -1 {			
				l.PushBack(i)
			}
		}
	}
	// search map of elevators
	// if not there, add
	// reallocate elev.other_id if adding new elevator
}

func (ctrl *Control) Send_List(to int){

	l *List = ctrl.elev.requests
	var send_args []string

	// fill send args will floor requests
	for i:=l.Front(); i!=nil; i=i.Next(){
		send_args[i] = Itoa(i)
	}
	Ctrl.Send_Msg(ctrl.elev.id, to, "MyList", send_args)

}














