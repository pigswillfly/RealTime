package Control

import(
	."Driver"
	."fmt"
	."time"
	."container/list"
	."strconv"
	."Helpers"
	."math"
)

const (
	up = 1
	stop = 0
	down = -1
	upSpeed = 150
	downSpeed = -150
	channelSize = 100
)

type Elevator struct {
	id int
	other_id int				
	speed int			
	direction int
	other_dir int
	floor int
	alignment int
	requests *List
	Msg chan string
	Alive chan string	
}

func Init_Elev() *Elevator {

	// driver init
	if Elev_Init() == 0{
		Println("Elevator initialization failure")
	}

	// declare new elevator
	elev := new(Elevator)

	// set up parameters
	elev.speed = Set_Speed(stop)
	elev.direction = 0
	elev.alignment = Get_Floor_Sensor_Signal()
	if elev.alignment != 0{
		for{
			elev.speed = Set_Speed(downSpeed)
			elev.alignment = Get_Floor_Sensor_Signal()
			if elev.alignment == 0{
				elev.speed = Set_Speed(stop)
				break
			}
		}
	}
	elev.requests = new(List)
	elev.Msg = make(chan string, channelSize)

	go elev.Pulse()
//	go elev.Get_Stats()
	go elev.Run()
	go elev.Poll_Buttons()

	return elev
}

func (elev *Elevator) Pulse(){
	for{
		msg := Itoa(elev.id)+",0,Alive"
		elev.Alive <- msg
		Sleep(500*Millisecond)
	}
}

func (elev *Elevator) Run(){
	
	l := elev.requests

	for{
		// get next floor as destination
		dest := elev.requests.Front().Value.(int)
		if elev.requests.Front() == nil{
			dest = -1
		}

		// set floor light
		if elev.alignment != -1{
			Set_Floor_Indicator(elev.floor)
		}

		dir := elev.direction

		if dest == -1{
			elev.direction = 0
		// if elevator above destination
		} else if elev.floor > dest {
			// set direction down
			elev.direction = -1 
			elev.speed = Set_Speed(downSpeed)
		// if elevator below destination
		} else if elev.floor < dest {
			// set direction up
			elev.direction = 1
			elev.speed = Set_Speed(upSpeed)
		// if stop signal, stop
		} else if Get_Stop_Signal() == 1 {
			elev.speed = Set_Speed(stop)
			break
		// otherwise
		} else {
			// stop elevator
			elev.speed = Set_Speed(stop)
			// remove request from queue
			elev.requests.Remove(l.Front())

			// button lights
			// if destination is lower than next request
			if dest < elev.requests.Front().Value.(int) {
				Set_Button_Light(BUTTON_CALL_UP, dest, 0)
			// if destination is higher than next request
			} else if dest > elev.requests.Front().Value.(int) {
				Set_Button_Light(BUTTON_CALL_DOWN, dest, 0)
			// otherwise
			} else {
				Set_Button_Light(BUTTON_CALL_UP, dest, 0)
				Set_Button_Light(BUTTON_CALL_DOWN, dest, 0)
			}

			// turn off command light
			Set_Button_Light(BUTTON_COMMAND, dest, 0)

			Sleep(Second)

		}
		
		// sort requests in terms of direction 
		if dir != elev.direction {
			Sort_Queue(elev.direction,elev.requests)
		}
	}
}

func (elev *Elevator) Update_Floor(){
	// update alignment (Get_Floor_Sensor_Signal)
	// control message current floor and direction of travel

	for {
		elev.alignment = Get_Floor_Sensor_Signal()
		if elev.alignment != -1 && elev.floor != elev.alignment{
				elev.floor = elev.alignment
		}
	}
}

func (elev *Elevator) Get_Stats(){
	for {
		Println("Elevator stats for elevator \t", elev.id)
		Println("Speed: \t", elev.speed)
		Println("Direction: \t", elev.direction)
		Println("Floor: \t", elev.floor)
		Println("Number of waiting requests: \t", elev.requests.Len())
		
		Sleep(Second)
	}
}

func (elev *Elevator) Poll_Buttons(){
	// polls buttons for requests
	old_button := -1
	old_floor := -1
	new_button := -1
	new_floor := -1

	for {
		// iterate through all floors & buttons
		new_button = -1
		new_floor = -1
		for floor:=0; floor<N_FLOORS; floor++{
			for button:=0; button<N_BUTTONS; button++{
				// check validity of button call in relation to floor
				if (button == BUTTON_CALL_UP && floor == N_FLOORS-1)|| 
					(button == BUTTON_CALL_DOWN && floor == 0){
					// next iteraion of button loop
					continue 
				}
				// if there is a request
				if Get_Button_Signal(button, floor) == 1 {
					if button == BUTTON_COMMAND{
						// add request to queue if command
						elev.Add_Request(button, floor)
						// TODO check if added
					// otherwise send control message
					} else {
						// control message
						msg := Itoa(elev.id)+",0,Request,"+Itoa(button)+","+Itoa(floor)
						elev.Msg <- msg
					}
					Println("Button ",button," was pressed at floor ",floor)
					new_button = button
					new_floor = floor
				}
			}
		}
	
		if (new_button != -1 && new_floor != -1){
			old_button = new_button
			old_floor = new_floor
			Set_Button_Light(old_button, old_floor, 1)
		}

		if Get_Stop_Signal() == 1{
			return
		}
		Sleep(Millisecond)
	}
}

func (elev *Elevator) Add_Request(button int, floor int) int {
	// returns 0 on success, -1 failure
	l := elev.requests

	// if no requests or elevator stopped
	if l.Len() == 0 || elev.direction == 0 {
		// list is empty, insert at front
		_ = l.PushFront(floor)
		return 0
	}

	// search list, if requests already exists, do nothing
	for i:=l.Front(); i!=nil; i=i.Next(){
		if i.Value.(int) == floor{
			return 1
		}
	}

	// return value
	ret := 1	

	// if direction is up
	if elev.direction == 1 {
		// insert before lowest floor higher than request
		// searching from front to back
		for i:=l.Front(); i!=nil; i=i.Next() {
			// if current value is higher than request
			// & request is higher than current floor
			if i.Value.(int) > floor && floor > elev.floor {
				// if call down break and restart
				if button == BUTTON_CALL_DOWN {
					// next iteration of request queue loop
					continue
				}
				// otherwise, insert before 
				_ = l.InsertBefore(floor, i)
				ret = 0
				break
			}
		}
	// if direction is down
	} else if elev.direction == -1 {
		// insert before highest floor slower than request
		// searching from front to back
		for i:= l.Front(); i!=nil; i=i.Next() {
			// if current value is lower than request
			// & request is lower than current floor
			if (i.Value.(int) < floor) && (floor > elev.floor) {
				// if call up break and restart
				if button == BUTTON_CALL_UP{
					// next iteration of request queue loop
					continue
				}
				// otherwise, insert before
				_ = l.InsertBefore(floor, i)
				ret = 0
				break			
			}
		}
	}
	// otherwise, add floor at back of queue
	if ret < 0 {
		_ = l.PushBack(floor)
	}
	return ret
}

func (elev *Elevator) Cost(button int, floor int) int{

	// TODO check if direction changes when servicing a floor (open/close doors)
	// if no requests in queue, or lift is at floor
	if (elev.direction==stop || elev.requests.Len()==0 || elev.At(floor)){
		return 0

	// if elevator is passing floor and request is in same direction
	} else if (elev.direction==down && elev.Above(floor) && button==down) ||
		(elev.direction==up && elev.Below(floor) && button==up){
		return 1

	// if elevator is passing floor and request is in opposite direction
	} else if (elev.direction==down && elev.Above(floor) && button==up) ||
		(elev.direction==up && elev.Below(floor) && button==down){
		return 2

	// if request floor and direction are opposite to elevator direction
	} else if (elev.direction==down && elev.Below(floor) && button==up) ||
		(elev.direction==up && elev.Above(floor) && button==down){
		return 3

	// if request and elevator direction are same but floor already passed
	} else if (elev.direction==down && elev.Below(floor) && button==down) ||
		(elev.direction==up && elev.Above(floor) && button==up){
		return 4

	// error
	} else {
		return -1
	}
}

func (elev *Elevator) Tie_Breaker(round int, floor int) int{
	if round==1{
		return elev.requests.Len()
	} else if round==2{
		return int(Abs(float64(elev.floor - floor)))
	} else {
		return 100
	}
}

func (elev *Elevator) Above(floor int) bool{
	return floor < elev.floor
}
func (elev *Elevator) Below(floor int) bool{
	return floor > elev.floor
}
func (elev *Elevator) At(floor int) bool{
	return floor == elev.floor
}




















