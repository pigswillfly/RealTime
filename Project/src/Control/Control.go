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
	tie map[int]int				// for tiebreakers, cost referened by id
	timers map[int]int			// alive msg timers referenced by id
	
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
	//	[Sender ID],[Receiver ID],[Code],[Arguments]
	msg := Itoa(id) + "," + Itoa(to) + "," + code
	for i:= range args {
		if args[i] != ""{
			msg += "," + args[i]
		}
	}
	// Send to network
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
				ctrl.costs[id] = args[0]	
	
			case "TieBreaker":
				// args -- round, floor
				send_args[0] = Itoa(ctrl.elev.TieBreaker(args[0], args[1]))
				Ctrl.Send_Msg(ctrl.elev.id, from_id, "MyTie", send_args)

			case "MyTie":
				// args -- tiebreaker result
				ctrl.tie[id] = args[0]
				
			case "ListPlease":
				// no args
				Send_List(from_id)

			case "MyList":
				// args -- direction, then floor requests
				ctrl.elev.other_dir,_ = Atoi(args[0])
				Update_List(from_id, args[1:])

			case "AddFloor":
				// args -- button, floor for request
				ctrl.elev.Add_Request(args[0], args[1])
				
		}
	}
}

func (ctrl *Control) Set_Elev_ID(){
	l *List = ctrl.elevators
	ctrl.elev.id = l.Back()+1

	// TODO possibly use socket ids for elevator ids somehow?
}

func (ctrl *Control) Set_Elev_Other_ID(){
	l *List = ctrl.elevators
	if ctrl.elev.id == l.Back(){
		ctrl.elev.other_id = l.Front()
	} else {
		for i:=l.Front(); i!=nil; i=i.Next(){
			if ctrl.elev.id == i {
				ctrl.elev.other_id = i.Next()
			}		
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
	for{
		msg := <-ctrl.elev.msg
		id, request := Decipher_Request(msg)
		l *List = ctrl.elevators
		for i:=l.Front(); i!=nil; i=i.Next(){
			Send_Msg(id, i, "CostPlease", request)
		}
	}
}

func (ctrl *Control) Decipher_Request(msg string) (int, []string){
	substrings := SplitN(msg, ",", -1)
	id,_ := Atoi(substrings[0])
	args := substrings[1:]

	return id, args
}

func (ctrl *Control) Decide_Elevator(button int, floor int){
	// TODO
	// Send request for costs
	// wait until all responses recieved
	// use list of elevators to find lowest cost
	// use tiebreaker if necessary
	// find lowest cost
	// send "AddFloor" message
}

func (ctrl *Control) Send_Alive_Msg(){

	// send alive message every time available in the alive channel
	for {
		msg := <-ctrl.elev.alive
		ctrl.toAliveNet <-msg
	}
}

func (ctrl *Control) Rec_Alive_Msg(){

	// if message available, receive message
	for {
		msg := <-ctrl.fromAliveNet

		// decode
		id,code,_ := Decipher_Msg(msg)
		if code == "Alive" {
			Update_Elevator_List(id)
			// reset timer for that elevator
		} else {
			// error
		}

		// TODO
		// Timers for detecting dead elevators
}
	}
}

func (ctrl *Control) Update_Elevator_List(id int){

	// Search list of elevators, if not found, add 
	found := 0
	l *List = ctrl.elevators
	for i:=l.Front(); i!=nil; i=i.Next(){
		if i==id{
			found = 1
			break
		}
	}
	if found == 0 {
		l.PushBack(id)
		Set_Elev_Other_ID()
	}
}

func (ctrl *Control) Update_List(id int, requests []string){

	if backup[id] != nil{
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














