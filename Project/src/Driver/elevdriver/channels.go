package channels

const (
    PORT4 = 3 
    OBSTRUCTION = (0x300+23)
    STOP = (0x300+22)
    FLOOR_COMMAND1 = (0x300+21)
    FLOOR_COMMAND2 = (0x300+20)
    FLOOR_COMMAND3 = (0x300+19)
    FLOOR_COMMAND4 = (0x300+18)
    FLOOR_UP1 = (0x300+17)
    FLOOR_UP2 = (0x300+16)

    //in port 1
    PORT1 = 2
    FLOOR_DOWN2 = (0x200+0)
    FLOOR_UP3 = (0x200+1)
    FLOOR_DOWN3 = (0x200+2)
    FLOOR_DOWN4 = (0x200+3)
    SENSOR1 = (0x200+4)
    SENSOR2 = (0x200+5)
    SENSOR3 = (0x200+6)
    SENSOR4 = (0x200+7)

    //out port 3
    PORT3 = 3
    MOTORDIR = (0x300+15)
    LIGHT_STOP = (0x300+14)
    LIGHT_COMMAND1 = (0x300+13)
    LIGHT_COMMAND2 = (0x300+12)
    LIGHT_COMMAND3 = (0x300+11)
    LIGHT_COMMAND4 = (0x300+10)
    LIGHT_UP1 = (0x300+9)
    LIGHT_UP2 = (0x300+8)

    //out port 2
    PORT2 = 3
    LIGHT_DOWN2 = (0x300+7)
    LIGHT_UP3 = (0x300+6)
    LIGHT_DOWN3 = (0x300+5)
    LIGHT_DOWN4 = (0x300+4)
    DOOR_OPEN = (0x300+3)
    FLOOR_IND2 = (0x300+1)
    FLOOR_IND1 = (0x300+0)

    //out port 0
    PORT0 = 1
    MOTOR = (0x100+0)

    //non-existing ports (for alignment)
    FLOOR_DOWN1 = -1
    FLOOR_UP4 = -1
    LIGHT_DOWN1 = -1
    LIGHT_UP4 = -1
)

const button_type_t(
    BUTTON_CALL_UP = 0
    BUTTON_CALL_DOWN = 1
    BUTTON_COMMAND = 2
)
