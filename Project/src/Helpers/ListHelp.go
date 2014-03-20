package Helpers

import(
	."fmt"
	."container/List"
)

func Sort_Queue(dir int, l *List){
	var temp int
	for i:=l.Front(); i!=nil; i=i.Next(){
		for j:=l.Back(); j!=i; j=j.Prev(){
			if dir < 0 {
				// Max first
				if i.Value < j.Value {
					// Swap		
					if i.Prev() != nil{
						temp = i.Prev()
						l.MoveAfter(i,j)
						l.MoveAfter(j,temp)
					} else if i.Next() != nil {
						temp = i.Next()
						l.MoveAfter(i,j)
						l.MoveBefore(j,temp)f
					}
					break
				}
			} else {
				// Min first
				if i.Value > j.Value {
					// Swap
					if i.Prev() != nil{
						temp = i.Prev()
						l.MoveAfter(i,j)
						l.MoveAfter(j,temp)
					} else if i.Next() != nil {
						temp = i.Next()
						l.MoveAfter(i,j)
						l.MoveBefore(j,temp)
					}
					break
				}
			}
		}
	} 
}
