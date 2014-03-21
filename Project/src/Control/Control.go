package Control

import(
//	."fmt"
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
)

type Control struct{
	elev *Elevator				// elevator
	elevators *List				// list of other elevators
	backup map[int]*List		// backup list of requests of other elev
	costs map[int]int			// costs of other elevators referenced by id
	cost_res int				// number of cost responses received
	tie map[int]int				// for tiebreakers, cost referenced by id
	tie_res int					// number of tie responses receieved
	timers map[int]*Timer		// alive msg timers referenced by id
	net *Network				// network
}

func Init_Control() *Control{

	// New Control object
	ctrl := new(Control)

	// Set up network
	ctrl.net = Init_Net()

	// Find which elevators are present ("Alive" messages receiving)
	go ctrl.Recieve_Msg()
	// Wait
	Sleep(Second*5)

	// Set up elevator
	ctrl.elev = Init_Elev()
	ctrl.Set_Elev_ID()

	// Set up alive message pulse
	go ctrl.Send_Pulse()

	// Start polling buttons
	go ctrl.Poll_Buttons()

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
	for{
		// Message format:
		//	[Sender ID],[Receiver ID],[Code],[Arguments]
		msg := <-ctrl.net.FromNet
		from_id, to_id, code, args := ctrl.Decipher_Msg(msg)
		var send_args []string

		if code == "Alive"{
			ctrl.Rec_Alive_Msg(from_id)
		} else{
		// check received ID against own
			if to_id == ctrl.elev.id{
				// action according to code
				// TODO check state?
				switch code{
					case "CostPlease":
						// args -- button, floor
						send_args[0] = Itoa(ctrl.elev.Cost(Atoi(args[0]),Atoi(args[1])))
						ctrl.Send_Msg(ctrl.elev.id, from_id, "MyCost", send_args)

					case "MyCost":
						// args -- cost
						ctrl.costs[from_id] = Atoi(args[0])
						ctrl.cost_res += 1	
	
					case "TieBreaker":
						// args -- round, floor
						send_args[0] = Itoa(ctrl.elev.Tie_Breaker(Atoi(args[0]), Atoi(args[1])))
						ctrl.Send_Msg(ctrl.elev.id, from_id, "MyTie", send_args)

					case "MyTie":
						// args -- tiebreaker result
						ctrl.tie[from_id] = Atoi(args[0])
						ctrl.tie_res += 1
				
					case "ListPlease":
						// no args
						ctrl.Send_List(from_id)

					case "MyList":
						// args -- direction, then floor requests
						ctrl.elev.other_dir,_ = Atoi(args[0])
						ctrl.Update_List(from_id, args[1:])

					case "AddFloor":
						// args -- button, floor for request
						ctrl.elev.Add_Request(Atoi(args[0]), Atoi(args[1]))		
				
				}
			}
		}
	}
}	

func (ctrl *Control) Send_Pulse(){

	// send alive message every time available in the alive channel
	for {
		msg := <-ctrl.elev.alive
		ctrl.net.ToNet <- msg
	}
}

func (ctrl *Control) Rec_Alive_Msg(id int){

	ctrl.Update_Elevator_List(id)
	// TODO reset timer for that elevator
	go ctrl.Reset_Timer(id)

}

func (ctrl *Control) Poll_Buttons(){
	for{
		msg := <-ctrl.elev.msg
		id, _, code, request := ctrl.Decipher_Msg(msg)
		if code == "Request"{
			ctrl.cost_res = 0
			l := ctrl.elevators
			for i:=l.Front(); i!=nil; i=i.Next(){
				ctrl.Send_Msg(id, i.Value.(int), "CostPlease", request)
			}
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
		if ctrl.cost_res == l.Len(){
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
				ctrl.Send_Msg(ctrl.elev.id, best_id, "TieBreaker", args)
				ctrl.Send_Msg(ctrl.elev.id, best_tie, "TieBreaker", args)
				for {
					if ctrl.tie_res == 2 {
						if ctrl.tie[best_id] > ctrl.tie[best_tie]{
							best_id = best_tie
							break
						} else if ctrl.tie[best_id] == ctrl.tie[best_tie]{
							// round 2
							args[0] = Itoa(2)
							ctrl.tie_res = 0
							ctrl.Send_Msg(ctrl.elev.id, best_id, "TieBreaker", args)
							ctrl.Send_Msg(ctrl.elev.id, best_tie, "TieBreaker", args)
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
			ctrl.Send_Msg(ctrl.elev.id, best_id, "AddFloor", args)
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
		stop:= ctrl.timers[id].Stop()
}

func (ctrl *Control) Timer_Expired(id int){
		// TODO
		// Timers for detecting dead elevators
		// 
}

func (ctrl *Control) Set_Elev_ID(){
	l := ctrl.elevators
	ctrl.elev.id = l.Back().Value.(int) + 1
}

func (ctrl *Control) Set_Elev_Other_ID(){
	l := ctrl.elevators
	if ctrl.elev.id == l.Back().Value.(int){
		ctrl.elev.other_id = l.Front().Value.(int)
	} else {
		for i:=l.Front(); i!=nil; i=i.Next(){
			if ctrl.elev.id == i.Value.(int) {
				ctrl.elev.other_id = i.Next().Value.(int)
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
		l.PushBack(id)
		ctrl.Set_Elev_Other_ID()
	}
}

func (ctrl *Control) Update_List(id int, requests []string){

	if ctrl.backup[id] != nil{
		l := new(List)
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

	l := ctrl.elev.requests
	var send_args []string

	// fill send args will floor requests
	j := 0
	for i:=l.Front(); i!=nil; i=i.Next(){
		send_args[j] = Itoa(i.Value.(int))
		j++
	}
	ctrl.Send_Msg(ctrl.elev.id, to, "MyList", send_args)

}












