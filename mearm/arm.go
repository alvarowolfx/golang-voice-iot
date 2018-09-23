package mearm

import (
	"math"
	"time"

	"github.com/alvarowolfx/golang-voice-iot/pca9685"
)

const (
	delayBetweenAngle = 15 * time.Millisecond
)

type RobotArm struct {
	grip     *pca9685.Servo
	elbow    *pca9685.Servo
	base     *pca9685.Servo
	shoulder *pca9685.Servo

	gripState                      bool
	elbowPos, shoulderPos, basePos int
}

type RobotArmPos struct {
	GripState   string `json:"grip"`
	ElbowPos    int    `json:"elbow"`
	ShoulderPos int    `json:"shoulder"`
	BasePos     int    `json:"base"`
}

func NewRobotArm(grip, base, elbow, shoulder *pca9685.Servo) *RobotArm {
	return &RobotArm{
		grip:        grip,
		base:        base,
		elbow:       elbow,
		shoulder:    shoulder,
		gripState:   true,
		elbowPos:    0,
		shoulderPos: 0,
		basePos:     0,
	}
}

func (ra *RobotArm) Init() {
	ra.SetArm(90, 90, 90)
	ra.CloseGrip()
}

/*
 * https://github.com/tobiastoft/ArduinoEasing/blob/master/Easing.cpp
 */
func easeInOutSine(t, b, c, d float64) float64 {
	return -c/2*(math.Cos(math.Pi*t/d)-1) + b
}

func (ra *RobotArm) updateServo(s *pca9685.Servo, oldValue, newValue int) {
	if oldValue == newValue {
		return
	}

	diff := newValue - oldValue
	duration := 53
	go func() {
		for pos := 0; pos < duration; pos++ {
			angle := easeInOutSine(float64(pos), float64(oldValue), float64(diff), float64(duration))
			intAngle := int(math.Ceil(angle))
			s.SetAngle(intAngle)
			time.Sleep(delayBetweenAngle)
		}
	}()
}

func (ra *RobotArm) SetBase(value int) {
	ra.updateServo(ra.base, ra.basePos, value)
	ra.basePos = value
}

func (ra *RobotArm) MoveBase(value int) {
	ra.updateServo(ra.base, ra.basePos, ra.basePos+value)
	ra.basePos += value
}

func (ra *RobotArm) SetShoulder(value int) {
	ra.updateServo(ra.shoulder, ra.shoulderPos, value)
	ra.shoulderPos = value
}

func (ra *RobotArm) MoveShoulder(value int) {
	ra.updateServo(ra.shoulder, ra.shoulderPos, ra.shoulderPos+value)
	ra.shoulderPos += value
}

func (ra *RobotArm) SetElbow(value int) {
	ra.updateServo(ra.elbow, ra.elbowPos, value)
	ra.elbowPos = value
}

func (ra *RobotArm) MoveElbow(value int) {
	ra.updateServo(ra.elbow, ra.elbowPos, ra.elbowPos+value)
	ra.elbowPos += value
}

func (ra *RobotArm) SetArm(base, shoulder, elbow int) {

	ra.updateServo(ra.base, ra.basePos, base)
	ra.basePos = base

	ra.updateServo(ra.shoulder, ra.shoulderPos, shoulder)
	ra.shoulderPos = shoulder

	ra.updateServo(ra.elbow, ra.elbowPos, elbow)
	ra.elbowPos = elbow
}

func (ra *RobotArm) CloseGrip() {
	if ra.gripState {
		ra.updateServo(ra.grip, 15, 120)
		ra.gripState = false
	}
}

func (ra *RobotArm) OpenGrip() {
	if !ra.gripState {
		ra.updateServo(ra.grip, 120, 15)
		ra.gripState = true
	}
}

func (ra *RobotArm) RobotArmPos() *RobotArmPos {
	grip := "OPEN"
	if !ra.gripState {
		grip = "CLOSED"
	}
	return &RobotArmPos{
		GripState:   grip,
		ElbowPos:    ra.elbowPos,
		ShoulderPos: ra.shoulderPos,
		BasePos:     ra.basePos,
	}
}
