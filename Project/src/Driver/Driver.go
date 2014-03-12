package Driver

import (
    "fmt"
    "math"   
    "C"
)

var lamp_channel_matrix [][]int = [4][3]int{
    [3]int{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
    [3]int{LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
    [3]int{LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
    [3]int{LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4}}

var button_channel_matrix [][]int = [4][3]int{
    [3]int{FLOOR_UP1, FLOOR_DOWN1, FLOOR_COMMAND1},
    [3]int{FLOOR_UP2, FLOOR_DOWN2, FLOOR_COMMAND2},
    [3]int{FLOOR_UP3, FLOOR_DOWN3, FLOOR_COMMAND3},
    [3]int{FLOOR_UP4, FLOOR_DOWN4, FLOOR_COMMAND4}}
    
const (
    BUTTON_CALL_UP int = 0
    BUTTON_CALL_DOWN int = 1
    BUTTON_COMMAND int = 2
)
   

func elev_init() int{
    // Initialize hardware
    if !IO_Init() {
        return 0
    }
    
    // Zero all floor button lights
    i := 0
    for i=0; i< N_FLOORS; i++ {
        if i != 0 {
            set_button_light(BUTTON_CALL_DOWN, i, 0)
        } 
        if i != N_FLOORS-1 {
            set_button_light(BUTTON_CALL_UP, i, 0)
        }
        set_button_light(BUTTON_COMMAND, i, 0)
    }
    
    // Clear stop light, door open light, set floor indicator to ground floor
    set_stop_light(0)
    set_door_open_light(0)
    set_floor_indicator(0)
    
    return 1
    
} 

func set_speed(speed int){
    // To sharply stop elevator, direction bit is toggled before setting speed to 0
    last_speed := 0;
    
    // If to start (speed > 0)
    // If to stop (speed == 0)
    
    if speed > 0 {                  
        IO_Clear_Bit(MOTORDIR)
    } else if speed < 0 {
        IO_Set_Bit(MOTODIR)
    } else if (last_speed < 0){    
        IO_Clear_Bit(MOTORDIR)
    } else if (last_speed > 0){
        IO_Set_Bit(MOTORDIR)
    }
    
    last_speed = speed
    
    // Write new setting to motor
    IO_Write_Analog(MOTOR, 2048 + 4*Abs(speed))
}

func get_floor_sensor_signal() int{
    if IO_Read_Bit(SENSOR1){
        return 0
    } else if IO_Read_Bit(SENSOR2){
        return 1
    } else if IO_Read_Bit(SENSOR3){
        return 2
    } else if IO_Read_Bit(SENSOR4){
        return 3
    } else{
        return -1
    }
}

func get_button_signal(button, floor int) int{
    // Make sure floor is 0 or greater
    if floor < 0{
        return -1
    }
    // Make sure floor doesn't exceed top floor
    if floor > N_FLOORS{
        return -1
    }
    // Make sure not call up and floor top floor
    if button == BUTTON_CALL_UP && floor == N_FLOORS-1{
        return -1
    }
    // Make sure not call down and floor 0
    if button == BUTTON_CALL_DOWN && floor == 0{
        return -1
    }
    // Make sure a button has been pressed
    if button != BUTTON_CALL_UP && button != BUTTON_CALL_DOWN && button != BUTTON_COMMAND{
        return -1
    }
    
    if IO_Read_Bit(button_channel_matrix[floor][button]){
        return 1
    } else{
        return 0
    }
    
}

func get_stop_signal()int{
    return IO_Read_Bit(STOP)
}

func get_obstruction_signal()int{
    return IO_Read_Bit(OBSTRUCTION)
}

func set_floor_indicator(floor int){
    // Make sure floor is 0 or greater
    if floor < 0{
        return -1
    }
    // Make sure floor doesn't exceed top floor
    if floor > N_FLOORS{
        return -1
    }
    
    // Binary encoding. One light must always be on
    if floor & 0x02{
        IO_Set_Bit(FLOOR_IND1)
    } else{
        IO_Clear_Bit(FLOOR_IND1)
    }
    
    if floor & 0x01{
        IO_Set_Bit(FLOOR_IND2)
    } else{
        IO_Clear_Bit(FLOOR_IND2)
    }
}

func set_button_light(button, floor int, value int){
    // Make sure floor is 0 or greater
    if floor < 0{
        return -1
    }
    // Make sure floor doesn't exceed top floor
    if floor > N_FLOORS{
        return -1
    }
    // Make sure not call up and floor top floor
    if button == BUTTON_CALL_UP && floor == N_FLOORS-1{
        return -1
    }
    // Make sure not call down and floor 0
    if button == BUTTON_CALL_DOWN && floor == 0{
        return -1
    }
    // Make sure a button has been pressed
    if button != BUTTON_CALL_UP && button != BUTTON_CALL_DOWN && button != BUTTON_COMMAND{
        return -1
    }
    
    if value == 1{
        IO_Set_Bit(lamp_channel_matrix[floor][button])
    } else{
        IO_Clear_Bit(lamp_channel_matrix[floor][button])    
    }
}

func set_stop_light(value int){
    if value{
        IO_Set_Bit(LIGHT_STOP)
    } else{
        IO_Clear_Bit(LIGHT_STOP)
    }
}

func set_door_open_light(value int){
    if value{
        IO_Set_Bit(DOOR_OPEN)
    } else{
        IO_Clear_Bit(DOOR_OPEN)
    }
}


