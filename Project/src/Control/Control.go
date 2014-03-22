package Control

import(
	."fmt"
	."container/list"
	."strconv"
	."time"
	."strings"
	."Network"
)

const (
	
	// for sorting lists
	minFirst = 1
	maxFirst = -1
	chanSize = 1024
)

type Control struct{
	elev *Elevator				// elevator
	ElevMsg chan string		
	ElevAlive chan string
	elevators *List				// list of other elevators
	friends int					// number of other elevators active
	costs map[int]int			// costs of other elevators referenced by id
	cost_res int				// number of cost responses received
	tie map[int]int				// for tiebreakers, cost referenced by id
	tie_res int					// number of tie responses receieved
	timers map[int]*Timer		// alive msg timers referenced by id
	net *Network				// network
	toNet chan string
	fromNet chan string

}

func Init_Control() *Control{

	// New Control object
	ctrl := new(Control)
	ctrl.elevators = New()
	ctrl.cost_res = -1
	ctrl.tie_res = -1
	ctrl.friends = 0
	ctrl.costs = make(map[int]int)
	ctrl.tie = make(map[int]int)
	ctrl.timers = make(map[int]*Timer)
	
	// Set up network
	ctrl.toNet = make(chan string, chanSize)
	ctrl.fromNet = make(chan string, chanSize)
	ctrl.net = Init_Net(ctrl.toNet, ctrl.fromNet)
	Println("Network established")

	// Find which elevators are present ("Alive" messages receiving)
	Println("Finding friends")
	go ctrl.Recieve_Msg()
	// Wait
	Sleep(Second)
	Println(Itoa(ctrl.friends)+" friends found")

	// Set up elevator
	ctrl.ElevMsg = make(chan string, chanSize)
	ctrl.ElevAlive = make(chan string, chanSize)
	ctrl.elev = Init_Elev(ctrl.ElevMsg, ctrl.ElevAlive)
	ctrl.Set_Elev_ID()

	// Set up alive message pulse
	go ctrl.Send_Alive_Msg()

	// Start receiving messages from elevator
	go ctrl.Rec_Elev_Msg()

	// Start elevator
	go ctrl.elev.Update_Floor()
	go ctrl.elev.Poll_Buttons()
	go ctrl.elev.Run()

	return ctrl
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
	ctrl.net.ToNet <- msg
}

func (ctrl *Control) Decipher_Msg(msg string) (int, int, string, []string){
	
	substrings := SplitN(msg, ",", -1)		
	from_id,_ := Atoi(substrings[0])
	to_id,_ := Atoi(substrings[1])
	code := substrings[2]
	args := substrings[3:]
	return from_id, to_id, code, args
}

func (ctrl *Control) Recieve_Msg(){

	Println("Receiving messages from network")
	for{
		// Message format:
		//	[Sender ID],[Receiver ID],[Code],[Arguments]
		msg := <-ctrl.net.FromNet
		from_id, to_id, code, args := ctrl.Decipher_Msg(msg)
		var send_args []string
		var i,j int

		if code == "Alive"{
			if from_id != ctrl.elev.ID{
				ctrl.Rec_Alive_Msg(from_id)
			}
		} else{
		// check received ID against own
			if to_id == ctrl.elev.ID{
				// action according to code
				// TODO check state?
				switch code{
					case "CostPlease":
						// args -- button, floor
						i,_ = Atoi(args[0])
						j,_ = Atoi(args[1])
						send_args[0] = Itoa(ctrl.elev.Cost(i,j))
						ctrl.Send_Msg(ctrl.elev.ID, from_id, "MyCost", send_args)

					case "MyCost":
						// args -- cost
						i,_ = Atoi(args[0])
						ctrl.costs[from_id] = i
						ctrl.cost_res += 1	
	
					case "TieBreaker":
						// args -- round, floor
						i,_ = Atoi(args[0])
						j,_ = Atoi(args[1])				
						send_args[0] = Itoa(ctrl.elev.Tie_Breaker(i, j))
						ctrl.Send_Msg(ctrl.elev.ID, from_id, "MyTie", send_args)

					case "MyTie":
						// args -- tiebreaker result
						i,_ = Atoi(args[0])
						ctrl.tie[from_id] = i
						ctrl.tie_res += 1
				
					case "ListPlease":
						// no args
						ctrl.Send_List(from_id)
						ctrl.elev.other_id = from_id

					case "MyList":
						// args -- direction, then floor requests
						ctrl.elev.backup_dir,_ = Atoi(args[0])
						ctrl.Update_List(from_id, args[1:])

					case "AddFloor":
						// args -- button, floor for request
						i,_ = Atoi(args[0])
						j,_ = Atoi(args[1])			
						ctrl.elev.Add_Request(i, j)		
				
				}
			}
		}
	}
}	

func (ctrl *Control) Send_Alive_Msg(){

	Println("Begin sending Alive messages to network")
	go ctrl.elev.Pulse()
	// send alive message every time available in the alive channel
	for {
		msg := <-ctrl.ElevAlive
		ctrl.net.ToNet <- msg
	}
}

func (ctrl *Control) Rec_Alive_Msg(id int){
	
	ctrl.Update_Elevator_List(id)
//	go ctrl.Reset_Timer(id)

}

func (ctrl *Control) Rec_Elev_Msg(){
	for{
		msg := <-ctrl.ElevMsg
		_, _, code, request := ctrl.Decipher_Msg(msg)
		if code == "Request"{
			ctrl.New_Request(request)
		} else if code == "Handled"{
			if ctrl.elev.other_id != -1{
				ctrl.Send_List(ctrl.elev.other_id)
			}
		}
	}
}

func (ctrl *Control) New_Request(request []string){
	button,_ := Atoi(request[0])
	floor,_ := Atoi(request[1])
	if ctrl.friends == 0{
		ctrl.elev.Add_Request(button,floor)
	} else {
		ctrl.costs[ctrl.elev.ID] = ctrl.elev.Cost(button,floor)
		ctrl.cost_res = 1
		l := ctrl.elevators
		for i:=l.Front(); i!=nil; i=i.Next(){
			ctrl.Send_Msg(ctrl.elev.ID, i.Value.(int), "CostPlease", request)
		}
	}
}

func (ctrl *Control) Decide_Elevator(button int, floor int){
	l := ctrl.elevators
	best_id := l.Front().Value.(int)
	best := ctrl.costs[best_id]
	best_tie := -1
	var args []string

	for{
		// if all responses received
		if ctrl.cost_res == ctrl.friends+1{
			// find best
			for i:=l.Front(); i!=nil; i=i.Next(){
				if i==l.Front(){
					continue
				} else {
					if ctrl.costs[i.Value.(int)] < best {	
						best_id = i.Value.(int)
						best = ctrl.costs[best_id]
						best_tie = -1
					} else if ctrl.costs[i.Value.(int)] == best{
						best_tie = i.Value.(int)
					}
				}
			}
			// tiebreaker
			if best_tie != -1{
				// round 1
				args[0] = Itoa(1)
				args[1] = Itoa(floor)
				ctrl.tie_res = 0
				ctrl.Send_Msg(ctrl.elev.ID, best_id, "TieBreaker", args)
				ctrl.Send_Msg(ctrl.elev.ID, best_tie, "TieBreaker", args)
				for {
					if ctrl.tie_res == 2 {
						if ctrl.tie[best_id] > ctrl.tie[best_tie]{
							best_id = best_tie
							break
						} else if ctrl.tie[best_id] == ctrl.tie[best_tie]{
							// round 2
							args[0] = Itoa(2)
							ctrl.tie_res = 0
							ctrl.Send_Msg(ctrl.elev.ID, best_id, "TieBreaker", args)
							ctrl.Send_Msg(ctrl.elev.ID, best_tie, "TieBreaker", args)
							for {
								if ctrl.tie_res == 2 {
									if ctrl.tie[best_id] > ctrl.tie[best_tie]{
										best_id = best_tie
										break
									}
									break
								}
							}
							break
						}
					}
				}
			}
			// Send message to selected elevator
			ctrl.Send_Msg(ctrl.elev.ID, best_id, "AddFloor", args)
			break
		}	
	}
}

func (ctrl *Control) Reset_Timer(id int){
	ctrl.timers[id] = NewTimer(5*Second)
	go func(){
		<-ctrl.timers[id].C
		ctrl.Timer_Expired(id)
	}()
	ctrl.timers[id].Stop()
}

func (ctrl *Control) Timer_Expired(id int){

	// Remove elevator from list of other elevators
	l := ctrl.elevators
	for i:=l.Front(); i!=nil; i=i.Next(){
		if i.Value.(int) == id{
			l.Remove(i)
			ctrl.friends -= 1
		}
	}
	// Reset back up elevator
	ctrl.Set_Elev_Backup_ID()

	// Redistribute requests if matches backup ID
	// Assume button = direction of elevator
	if id == ctrl.elev.backup_id{
		button := ctrl.elev.backup_dir
		var request []string
		request[0] = Itoa(button)
		l = ctrl.elev.backup
		for i:=l.Front(); i!=nil; i=i.Next(){
			request[1] = Itoa(i.Value.(int))
			ctrl.New_Request(request)
			Sleep(2*Second)
		}
	}
}

func (ctrl *Control) Set_Elev_ID(){

	l := ctrl.elevators
	if l.Back() == nil{
		ctrl.elev.ID = 1
	} else {
		ctrl.elev.ID = l.Back().Value.(int) + 1
	}
	Println("Elevator ID set to: "+Itoa(ctrl.elev.ID))
}

func (ctrl *Control) Set_Elev_Backup_ID(){
	l := ctrl.elevators
	if l.Back() == nil{
		// no other elevators, don't need backup_id
	} else {
		if ctrl.elev.ID == l.Back().Value.(int){
			if ctrl.elev.backup_id != l.Front().Value.(int){
				// New ID needed
				ctrl.elev.backup_id = l.Front().Value.(int)
			}
		} else {
			for i:=l.Front(); i!=nil; i=i.Next(){
				if ctrl.elev.ID == i.Value.(int) {
					if ctrl.elev.backup_id == i.Next().Value.(int){
						ctrl.elev.backup_id = i.Next().Value.(int)
					}
				}		
			}
		}
	}
}

func (ctrl *Control) Update_Elevator_List(id int){

	// Search list of elevators, if not found, add 
	found := 0
	l := ctrl.elevators
	for i:=l.Front(); i!=nil; i=i.Next(){
		if id==i.Value.(int){
			found = 1
			break
		}
	}
	if found == 0 {
		ctrl.friends += 1
		l.PushBack(id)
		ctrl.Set_Elev_Backup_ID()
	}
}

func (ctrl *Control) Update_List(id int, requests []string){

	if ctrl.elev.backup != nil{
		l := new(List)
		req := -1
		for i:= range requests{
			req,_ = Atoi(requests[i])
			if req > -1 {			
				l.PushBack(i)
			}
		}
		ctrl.elev.backup = l
	}
}

func (ctrl *Control) Request_List(){
	for {	
		if ctrl.elev.backup_id != -1{
			ctrl.Send_Msg(ctrl.elev.ID, ctrl.elev.backup_id, "ListPlease", nil)
		}
		Sleep(5*Second)
	}
}

func (ctrl *Control) Send_List(to int){

	l := ctrl.elev.requests
	var send_args []string

	if ctrl.elev.requests.Len() > 0 {
		// fill send args will floor requests
		send_args[0] = Itoa(ctrl.elev.direction)
		j := 1
		for i:=l.Front(); i!=nil; i=i.Next(){
			send_args[j] = Itoa(i.Value.(int))
			j++
		}

		ctrl.Send_Msg(ctrl.elev.ID, to, "MyList", send_args)
	}
}












