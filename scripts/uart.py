import time
import serial

ser = serial.Serial("/dev/ttyTHS2", 9600)

while 1:
    ser.write('a')
    time.sleep(1)
