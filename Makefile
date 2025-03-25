main:
	gcc `pkg-config --cflags gtk+-3.0` main.c -o main.o `pkg-config --libs gtk+-3.0`

clean:
	rm main.o
