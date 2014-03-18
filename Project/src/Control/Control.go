package Control

type Control struct{

}


func (ctrl *Control) Control(){
	for {

		// ACTIONS WHEN IN RIGHT STATE:
	
		// button push - poll function?
			// add to request list
			// allocate to elevator (wrt direction) - allocate function?
		// alive message
			// confirm in alive list - list management function
			// if not add new
		// dead message
			// remove from alive list - list management function
			// get requests from elevator's request queue and reallocate?
		// reached floor	
			// is it correct? - fulfilled function
			// open/close doors
			// remove request 
		// stop - stop function

		// Process pairs?
		// Primary with backup? - primary dead function
			// primary dead, backup becomes primary
			// new backup


		// states
			// 1. moving up
			// 2. moving down
			// 3. waiting (no requests)
			// 4. opening/closing doors
			// 5. stopped	
	}
}
