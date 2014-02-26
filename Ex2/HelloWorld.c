
#include <pthread.h>
#include <stdio.h>
#include <stdlib.h>

int i = 0;
pthread_mutex_t m;

void *adder()
{
	int x;
	for (x=0; x<1000000; x++)
	{
		pthread_mutex_lock(&m);
		i++;
		pthread_mutex_unlock(&m);
	}
	return NULL;
}

void *subber()
{
	int x;
	for (x=0; x<1000000; x++)
	{
		pthread_mutex_lock(&m);
		i--;
		pthread_mutex_unlock(&m);
	}
	return NULL;
}

int main (int argc, char* argv[])
{

	pthread_t pid1, pid2;
	int err;

	pthread_mutex_init(&m, NULL);
	
	err = pthread_create(&pid1, NULL, adder, NULL);
	if (err < 0)
		printf("Error creating thread 1\n");

	err = pthread_create(&pid2, NULL, subber, NULL);
	if (err < 0)
		printf("Error creating thread 2\n");

	pthread_join(pid1, NULL);
	pthread_join(pid2, NULL);

	pthread_mutex_destroy(&m);

	printf("Done: %i\n", i);
	
	return 0;

}


