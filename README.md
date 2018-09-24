# Golang + MeArm + IoT + Google Assistant 

Demo project on how to run a Golang program on an embbeded hardware like Raspberry/Orange Pi. In this case this project works as a controller for a MeArm robot arm that can be controlled by voice using Google Assistant.

## BOM - Bill of Materials

* OrangePi Zero 
* MeArm - Project can be downloaded [here](https://www.thingiverse.com/thing:993759)
* 4 Mini Servo motors
* PCA9685 16 Channel PWM module

## Schematic 

Work in Progress

# Pre Requisites 

* Non Root Access to GPIO: https://opi-gpio.readthedocs.io/en/latest/install.html#non-root-access

## How to build for OrangePi 

This command will generate a binary file compatible with Arm architecture.

`make build`

## Copy to OrangePi

Change your Pi address on the Makefile, then run the command: 

`make copy`