
#include <stdlib.h>
#include <stdio.h>
#include <string.h>

int i;
FILE *fd;
char filename[] = "alive.txt";

void primary();
void backup();

void die(char *s)
{
    perror(s);
    exit(1);
}

int
main(int argc, char *argv[]){

	int pid;

	do{
		pid = fork();
		if (pid > 0){
			// parent
			printf("I'm the parent with %d, and created child with  %d\n", getpid(), pid);

			primary();
		}
		else if (pid == 0){
			// child
			printf("I'm the child process %d, and my parent is %d\n", getpid(), getppid());
		
			backup();
		}
		else{	// pid == -1
			printf("fork() failed\n");
		}

	}while(1);

	return 0;
}

void
primary(){

	// Open status file for writing
	fd = fopen(filename, "w");
	if (fd == NULL){
		fprintf(stderr, "Can't open file to write\n");
		die("primary fopen()");
	}

	int n; 

	while(1){

		// Print number to screen
		if(printf("%d \n", i) > -1){
			// Successful print
			// Print number to status file
			fprintf(fd, "%d\n", i);
			i++;
		}

		sleep(1);
	}

	fclose(fd);
}

void
backup(){

	int dead = 0;
	int prev, read;
	do{
		printf("Backup alive\n");

		// Open status file for reading
		fd = fopen(filename, "r");
		if(fd == NULL){
			fprintf(stderr, "Can't open file to read\n");
			die("backup fopen()");
		}
	
		// Read file until end of file is set
		do{
			fscanf(fd, "%d", &read);
			// store current number
			i = read;
		}while(feof(fd) != 0);

		dead++;

		fclose(fd);

	}while(dead < 2); // when same value is read twice, become primary

}





