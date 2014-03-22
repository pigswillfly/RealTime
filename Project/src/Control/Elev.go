package Control

import(
	."Driver"
	."fmt"
	."time"
	."container/list"
	."strconv"
	."Helpers"
	."math"
	."strings"
)

const (
	up = 1
	stop = 0
	down = -1
	upSpeed = 300
	downSpeed = -300
	channelSize = 100
	StopWait = 200*Millisecond
)

type Elevator struct {
	ID int
	backup_id int	
	other_id int			
	speed int			
	direction int
	backup_dir int
	floor int
	alignment int
	requests *List
	backup *List
	Msg chan string
	Alive chan string	
}

func Init_Elev(msg chan string, alive chan string) *Elevator {

	// driver init
	if Elev_Init() == 0{
		Println("Elevator initialization failure")
	}

	// declare new elevator
	elev := new(Elevator)

	// set up parameters
	elev.ID = -1
	elev.backup_id = -1
	elev.other_id = -1
	elev.speed = Set_Speed(stop)

	elev.direction = 0
	elev.alignment = Get_Floor_Sensor_Signal()

	if elev.alignment == -1{
		for{
			elev.speed = Set_Speed(downSpeed)
			elev.alignment = Get_Floor_Sensor_Signal()
			if elev.alignment != -1{
				elev.Stop()
				break
			}
		}
	}

	elev.requests = New()
	elev.backup = New()
	elev.Msg = msg
	elev.Alive = alive

	return elev
}

func (elev *Elevator) Pulse(){
	Println("Elevator pulse active")
	for{
		msg := Itoa(elev.ID)+",0,Alive"
		elev.Alive <- msg
		Sleep(500*Millisecond)
	}
}

func (elev *Elevator) Run(){
	Println("Elevator Running")
	l := elev.requests
	dest := -1

	for{
		// If no requests...
		if elev.requests.Len() > 0{
			dest = elev.requests.Front().Value.(int)

			dir := elev.direction

			// if no requests, stop
			if dest == -1{
				elev.direction = stop
			// if elevator above destination
			} else if Get_Floor_Sensor_Signal() != -1 && Get_Floor_Sensor_Signal() > dest {
				// set direction down
				elev.direction = down
				elev.speed = Set_Speed(downSpeed)
			// if elevator below destination
			} else if Get_Floor_Sensor_Signal() != -1 && Get_Floor_Sensor_Signal() < dest {
				// set direction up
				elev.direction = up
				elev.speed = Set_Speed(upSpeed)
			// if stop signal, stop
			} else if Get_Stop_Signal() == 1{
				elev.Stop()
				break
			// otherwise
			} else if Get_Floor_Sensor_Signal() != -1 && Get_Floor_Sensor_Signal() == dest{
				// stop elevator
				elev.Stop()
				elev.direction = 0
				// remove request from queue
				elev.requests.Remove(l.Front())
				// send control message
				msg := Itoa(elev.ID)+",0,Handled"
				elev.Msg <- msg

				// door open light
				Set_Door_Open_Light(1)
				Sleep(2*Second)
				Set_Door_Open_Light(0)

				// button lights
				if elev.requests.Len() > 0{
					// if destination is lower than next request
					if dest < elev.requests.Front().Value.(int) {
						Set_Button_Light(BUTTON_CALL_UP, dest, 0)
					// if destination is higher than next request
					} else if dest > elev.requests.Front().Value.(int) {
						Set_Button_Light(BUTTON_CALL_DOWN, dest, 0)
					} 
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
		} else {
			Sleep(Second)
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
			Set_Floor_Indicator(elev.floor)
		}
	Sleep(10*Millisecond)
	}
}

func (elev *Elevator) Get_Stats(){
	Println("Getting stats...")
	for {
		Println("Elevator stats for elevator \t", elev.ID)
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

	Println("Polling buttons started")
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
					// otherwise send control message
					} else {
						// control message
						msg := Itoa(elev.ID)+",0,Request,"+Itoa(button)+","+Itoa(floor)
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
		Sleep(100*Millisecond)
	}
}

func (elev *Elevator) Add_Request(button int, floor int) int {
	// returns 0 on success, -1 failure
	l := elev.requests

	// search list, if requests already exists, do nothing
	for i:=l.Front(); i!=nil; i=i.Next(){
		if i.Value.(int) == floor{
			return 1
		}
	}

	// if no requests or elevator stopped
	if l.Len() == 0 || elev.direction == 0 {
		// list is empty, insert at front
		_ = l.PushFront(floor)
		return 0
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
		ret = 0
	}
	return ret
}

func (elev *Elevator) Cost(button int, floor int) int{

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
func (elev *Elevator) Stop(){
	Sleep(StopWait)
	if elev.direction == up{
		elev.speed = Set_Speed(downSpeed)
		elev.speed = Set_Speed(stop)
	} else if elev.direction == down{
		elev.speed = Set_Speed(upSpeed)
		elev.speed = Set_Speed(stop)
	}
}

func (elev *Elevator) No_Friends(){

	Println("No Friends established")
	for {
		msg := <-elev.Msg
		substrings := SplitN(msg, ",", -1)		
		code := substrings[2]
		args := substrings[3:]
		if code == "Request"{
			button,_ := Atoi(args[0])
			floor,_ := Atoi(args[1]) 
			elev.Add_Request(button,floor)
		}
	}
}

func (elev *Elevator) Print_Requests(){

	Println("Printing requests...")
	for {
		l:= elev.requests
		if (elev.requests.Len() > 0){
			Printf("Requests: %d", l.Front().Value.(int))
			for i:= l.Front(); i!=nil; i=i.Next(){
				if i==l.Front(){
					continue
				}
				Printf(", %d", i.Value.(int))
			}
			Printf("\n")
		}
		Sleep(2*Second)
	}
}




















