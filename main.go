package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/alvarowolfx/golang-voice-iot/mearm"
	"github.com/alvarowolfx/golang-voice-iot/middleware"

	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/experimental/devices/pca9685"
	"periph.io/x/periph/host"
)

func main() {
	_, err := host.Init()
	if err != nil {
		log.Fatal(err)
	}

	bus, err := i2creg.Open("0")
	if err != nil {
		log.Fatal(err)
	}

	pca, err := pca9685.NewI2C(bus, pca9685.I2CAddr)
	if err != nil {
		log.Fatal(err)
	}

	pca.SetPwmFreq(50 * physic.Hertz)
	pca.SetAllPwm(0, 0)
	servos := pca9685.NewServoGroup(pca, 50, 650, 0, 180)

	gripServo := servos.GetServo(0)
	baseServo := servos.GetServo(1)
	elbowServo := servos.GetServo(2)
	shoulderServo := servos.GetServo(3)

	gripServo.SetMinMaxAngle(15, 120)
	elbowServo.SetMinMaxAngle(50, 110)    // Set limit of the robot arm
	shoulderServo.SetMinMaxAngle(60, 140) // Set limit of the robot arm

	robotArm := mearm.NewRobotArm(gripServo, baseServo, elbowServo, shoulderServo)
	robotArm.Init()

	port := "8080"
	httpRobotArm := middleware.NewHttpRobotArm(robotArm, port)
	go httpRobotArm.Start() // Configure http handlers

	deviceID, err := os.Hostname()
	if err != nil {
		fmt.Println("Raw")
		log.Fatal(err)
	}

	projectID := os.Getenv("PROJECT_ID")
	region := os.Getenv("REGION")
	registryID := os.Getenv("REGISTRY_ID")
	privateKey := os.Getenv("DEVICE_KEY")
	certsCa := os.Getenv("CA_CERTS")

	ciotRobotArm, err := middleware.NewCloudIoTCoreRobotArm(robotArm, deviceID, projectID, region, registryID, privateKey, certsCa)
	if err != nil {
		log.Fatal(err)
	}
	ciotRobotArm.Init()  // Setup
	ciotRobotArm.Start() // Configure mqtt handlers

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-sigs // Wait for signal
		log.Println(sig)

		// Shutdown all services

		httpRobotArm.Stop()
		log.Println("Http server stopped")

		ciotRobotArm.Stop()
		log.Println("Mqtt handler stopped")

		done <- true

	}()

	log.Println("Press ctrl+c to stop...")
	<-done // Wait

}
