package Helpers

import(
	."container/list"
)

func Sort_Queue(dir int, l *List){
	temp := l.Front()
	for i:=l.Front(); i!=nil; i=i.Next(){
		for j:=l.Back(); j!=i; j=j.Prev(){
			if dir < 0 {
				// Max first
				if i.Value.(int) < j.Value.(int) {
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
			} else {
				// Min first
				if i.Value.(int) > j.Value.(int) {
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
