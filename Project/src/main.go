package main

import (
    "Driver"
)


func main() {
	Driver.Elev_Init()
	Driver.Set_Speed(300)
	Sleep(Second)
	Driver.Set_Speed(-300)
	
}
