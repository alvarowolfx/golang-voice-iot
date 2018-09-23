package pca9685

type ServoGroup struct {
	*Dev
	minPwm   int
	maxPwm   int
	minAngle int
	maxAngle int
}

type Servo struct {
	group    *ServoGroup
	channel  int
	minAngle int
	maxAngle int
}

func NewServoGroup(dev *Dev, minPwm, maxPwm, minAngle, maxAngle int) *ServoGroup {
	return &ServoGroup{
		Dev:      dev,
		minPwm:   minPwm,
		maxPwm:   maxPwm,
		minAngle: minAngle,
		maxAngle: maxAngle,
	}
}

func (s *ServoGroup) SetMinMaxPwm(minAngle, maxAngle, minPwm, maxPwm int) {
	s.maxPwm = maxPwm
	s.minPwm = minPwm
	s.minAngle = minAngle
	s.maxAngle = maxAngle
}

func (s *ServoGroup) SetAngle(channel, angle int) {
	value := mapValue(angle, s.minAngle, s.maxAngle, s.minPwm, s.maxPwm)
	s.Dev.SetPwm(channel, 0, uint16(value))
}

func (s *ServoGroup) GetServo(channel int) *Servo {
	return &Servo{
		group:    s,
		channel:  channel,
		minAngle: s.minAngle,
		maxAngle: s.maxAngle,
	}
}

func (s *Servo) SetMinMaxAngle(min, max int) {
	s.minAngle = min
	s.maxAngle = max
}

func (s *Servo) SetAngle(angle int) {
	if angle < s.minAngle {
		angle = s.minAngle
	}
	if angle > s.maxAngle {
		angle = s.maxAngle
	}
	s.group.SetAngle(s.channel, angle)
}

func (s *Servo) SetPwm(pwm uint16) {
	s.group.SetPwm(s.channel, 0, pwm)
}

func mapValue(x, inMin, inMax, outMin, outMax int) int {
	return (x-inMin)*(outMax-outMin)/(inMax-inMin) + outMin
}
