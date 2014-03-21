package Testing

import(
	"Driver"
	"Network"
	"time"
	."fmt"
)

const (
	up = 150
	down = -150
	stop = 0
)

func Network_Test(){
	net := Network.Init_Net()
	msg := "Hello World!"
	var message string
	net.ToNet <- msg
	time.Sleep(2*time.Second)
	message = <-net.FromNet
	Println(message)
}

func Driver_Test(){
	Driver.Elev_Init()

	Println("Testing motor and floor indicators\n")
	Driver.Set_Speed(up)
	for{
		if(Driver.Get_Floor_Sensor_Signal() == 1){
			Driver.Set_Speed(down)
			Driver.Set_Speed(stop)
			break
		}
	}
	Driver.Set_Floor_Indicator(Driver.Get_Floor_Sensor_Signal())
	time.Sleep(time.Millisecond*2000)
	
	Driver.Set_Speed(up)
	for{
		if(Driver.Get_Floor_Sensor_Signal() == 2){
			Driver.Set_Speed(down)
			Driver.Set_Speed(stop)
			break
		}
	}
	Driver.Set_Floor_Indicator(Driver.Get_Floor_Sensor_Signal())
	time.Sleep(time.Millisecond*2000)

	Driver.Set_Speed(up)
	for{
		if(Driver.Get_Floor_Sensor_Signal() == 3){
			Driver.Set_Speed(down)
			Driver.Set_Speed(stop)
			break
		}
	}
	Driver.Set_Floor_Indicator(Driver.Get_Floor_Sensor_Signal())
	time.Sleep(time.Millisecond*2000)

	Driver.Set_Speed(down)
	for{
		if(Driver.Get_Floor_Sensor_Signal() == 2){
			Driver.Set_Speed(down)
			Driver.Set_Speed(stop)
			break
		}
	}
	Driver.Set_Floor_Indicator(Driver.Get_Floor_Sensor_Signal())
	time.Sleep(time.Millisecond*2000)

	Driver.Set_Speed(down)
	for{
		if(Driver.Get_Floor_Sensor_Signal() == 1){
			Driver.Set_Speed(down)
			Driver.Set_Speed(stop)
			break
		}
	}
	Driver.Set_Floor_Indicator(Driver.Get_Floor_Sensor_Signal())
	time.Sleep(time.Millisecond*2000)

	Driver.Set_Speed(down)
	for{
		if(Driver.Get_Floor_Sensor_Signal() == 0){
			Driver.Set_Speed(down)
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

	Driver.Set_Button_Light(2,1,1)
	time.Sleep(time.Millisecond*500)
	Driver.Set_Button_Light(2,2,1)
	time.Sleep(time.Millisecond*500)
	Driver.Set_Button_Light(2,3,1)
	time.Sleep(time.Millisecond*500)
	Driver.Set_Button_Light(2,1,0)
	time.Sleep(time.Millisecond*500)
	Driver.Set_Button_Light(2,2,0)
	time.Sleep(time.Millisecond*500)
	Driver.Set_Button_Light(2,3,0)

}


