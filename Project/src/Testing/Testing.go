package Testing

import(
	"Driver"
	."Network"
	"time"
	."fmt"
	"Control"
)

const (
	up = 300
	down = -300
	stop = 0
	chanSize = 1024
	N_FLOORS = 4
	N_BUTTONS = 3
	BUTTON_CALL_UP = 0
	BUTTON_CALL_DOWN = 1
	BUTTON_COMMAND = 2
	stopWait = 0
	oppWait = 5
)

func Network_Test(){

	// Run on two machines!
	to := make(chan string, chanSize)
	from := make(chan string, chanSize)
	net := Init_Net(to,from)
	msg := "Hello World!"
	var message string
	net.ToNet <- msg
	time.Sleep(2*time.Second)
	message = <-net.FromNet
	Println(message)
}
	
func Elevator_Test(){
	
	msg := make(chan string, chanSize)
	alive := make(chan string, chanSize)
	elev := Control.Init_Elev(msg, alive)
//	go elev.Get_Stats()
	go elev.Update_Floor()
	go elev.No_Friends()
	go elev.Poll_Buttons()
	go elev.Print_Requests()
	go elev.Run()
}

func Driver_Test(){
	Driver.Elev_Init()

	Println("Testing motor and floor indicators\n")
	Driver.Set_Speed(up)
	for{
		if(Driver.Get_Floor_Sensor_Signal() == 1){
			time.Sleep(stopWait)
			Driver.Set_Speed(down)
			time.Sleep(oppWait*time.Millisecond)
			Driver.Set_Speed(stop)
			break
		}
	}
	Driver.Set_Floor_Indicator(Driver.Get_Floor_Sensor_Signal())
	time.Sleep(time.Millisecond*2000)
	
	Driver.Set_Speed(up)
	for{
		if(Driver.Get_Floor_Sensor_Signal() == 2){
			time.Sleep(stopWait)
			Driver.Set_Speed(down)
			time.Sleep(oppWait*time.Millisecond)
			Driver.Set_Speed(stop)
			break
		}
	}
	Driver.Set_Floor_Indicator(Driver.Get_Floor_Sensor_Signal())
	time.Sleep(time.Millisecond*2000)

	Driver.Set_Speed(up)
	for{
		if(Driver.Get_Floor_Sensor_Signal() == 3){
			time.Sleep(stopWait)
			Driver.Set_Speed(down)
			time.Sleep(oppWait*time.Millisecond)
			Driver.Set_Speed(stop)
			break
		}
	}
	Driver.Set_Floor_Indicator(Driver.Get_Floor_Sensor_Signal())
	time.Sleep(time.Millisecond*2000)

	Driver.Set_Speed(down)
	for{
		if(Driver.Get_Floor_Sensor_Signal() == 2){
			time.Sleep(stopWait)
			Driver.Set_Speed(up)
			time.Sleep(oppWait*time.Millisecond)
			Driver.Set_Speed(stop)
			break
		}
	}
	Driver.Set_Floor_Indicator(Driver.Get_Floor_Sensor_Signal())
	time.Sleep(time.Millisecond*2000)

	Driver.Set_Speed(down)
	for{
		if(Driver.Get_Floor_Sensor_Signal() == 1){
			time.Sleep(stopWait)
			Driver.Set_Speed(up)
			time.Sleep(oppWait*time.Millisecond)
			Driver.Set_Speed(stop)
			break
		}
	}
	Driver.Set_Floor_Indicator(Driver.Get_Floor_Sensor_Signal())
	time.Sleep(time.Millisecond*2000)

	Driver.Set_Speed(down)
	for{
		if(Driver.Get_Floor_Sensor_Signal() == 0){
			time.Sleep(stopWait)
			Driver.Set_Speed(up)
			time.Sleep(oppWait*time.Millisecond)
			Driver.Set_Speed(stop)
			break
		}
	}
	Driver.Set_Floor_Indicator(Driver.Get_Floor_Sensor_Signal())
	time.Sleep(time.Millisecond*2000)

	Println("Testing door open light\n")
	Driver.Set_Door_Open_Light(1)
	time.Sleep(time.Millisecond*2000)
	Driver.Set_Door_Open_Light(0)

	Println("Testing stop light\n")
	Driver.Set_Stop_Light(1)
	time.Sleep(time.Millisecond*2000)
	Driver.Set_Stop_Light(0)

	// Set_Button_Light(button int, floor int, value int)
	Println("Testing button lights\n")
	Driver.Set_Button_Light(0,0,1)
	time.Sleep(time.Millisecond*500)
	Driver.Set_Button_Light(0,1,1)
	time.Sleep(time.Millisecond*500)
	Driver.Set_Button_Light(0,2,1)
	time.Sleep(time.Millisecond*500)
	Driver.Set_Button_Light(0,0,0)
	time.Sleep(time.Millisecond*500)
	Driver.Set_Button_Light(0,1,0)
	time.Sleep(time.Millisecond*500)
	Driver.Set_Button_Light(0,2,0)
	time.Sleep(time.Millisecond*500)	

	Driver.Set_Button_Light(1,1,1)
	time.Sleep(time.Millisecond*500)
	Driver.Set_Button_Light(1,2,1)
	time.Sleep(time.Millisecond*500)
	Driver.Set_Button_Light(1,3,1)
	time.Sleep(time.Millisecond*500)
	Driver.Set_Button_Light(1,1,0)
	time.Sleep(time.Millisecond*500)
	Driver.Set_Button_Light(1,2,0)
	time.Sleep(time.Millisecond*500)
	Driver.Set_Button_Light(1,3,0)
	time.Sleep(time.Millisecond*500)

	Driver.Set_Button_Light(2,0,1)
	time.Sleep(time.Millisecond*500)
	Driver.Set_Button_Light(2,1,1)
	time.Sleep(time.Millisecond*500)
	Driver.Set_Button_Light(2,2,1)
	time.Sleep(time.Millisecond*500)
	Driver.Set_Button_Light(2,3,1)
	time.Sleep(time.Millisecond*500)
	Driver.Set_Button_Light(2,0,0)
	time.Sleep(time.Millisecond*500)
	Driver.Set_Button_Light(2,1,0)
	time.Sleep(time.Millisecond*500)
	Driver.Set_Button_Light(2,2,0)
	time.Sleep(time.Millisecond*500)
	Driver.Set_Button_Light(2,3,0)

	Println("Testing buttons")
	for {
		for floor:=0; floor<N_FLOORS; floor++{
			for button:=0; button<N_BUTTONS; button++{
				if Driver.Get_Button_Signal(button, floor) == 1 {
					Println("Button ",button," was pressed at floor ",floor)
					time.Sleep(100*time.Millisecond)
				}
			}
		}
	}

}


