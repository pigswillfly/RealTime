package Driver

/*
#cgo LDFLAGS: -lcomedi -lm
#include "io.h"
*/
import "C"


func IO_Init(value int){
	C.io_init(C.int(value))
}

func IO_Set_Bit(channel int){
	C.io_set_bit(C.int(channel))
}

func IO_Clear_Bit(channel int){
	C.io_clear_bit(C.int(channel))
}

func IO_Write_Analog(channel int, value int){
	C.io_write_analog(C.int(channel), C.int(value))
}

func IO_Read_Bit(channel int) int{
	return C.io_read_bit(C.int(channel))
}

func IO_Read_Analog(channel int) int{
	return C.io_read_analog(C.int(channel))
}




