package pca9685

import (
	"math"
	"time"

	"periph.io/x/periph/conn/i2c"
)

const (
	PCA9685_ADDRESS uint16 = 0x40

	MODE1         byte = 0x00
	MODE2         byte = 0x01
	SUBADR1       byte = 0x02 // NOSONAR
	SUBADR2       byte = 0x03 // NOSONAR
	SUBADR3       byte = 0x04 // NOSONAR
	PRESCALE      byte = 0xFE
	LED0_ON_L     byte = 0x06
	LED0_ON_H     byte = 0x07
	LED0_OFF_L    byte = 0x08
	LED0_OFF_H    byte = 0x09
	ALL_LED_ON_L  byte = 0xFA
	ALL_LED_ON_H  byte = 0xFB
	ALL_LED_OFF_L byte = 0xFC
	ALL_LED_OFF_H byte = 0xFD

	// Bits
	RESTART byte = 0x80 // NOSONAR
	SLEEP   byte = 0x10
	ALLCALL byte = 0x01
	INVRT   byte = 0x10 // NOSONAR
	OUTDRV  byte = 0x04
)

type Dev struct {
	dev *i2c.Dev
}

func NewI2CAddress(bus i2c.Bus, address uint16) (*Dev, error) {
	dev := &Dev{
		dev: &i2c.Dev{Bus: bus, Addr: address},
	}
	err := dev.init()
	return dev, err
}

func NewI2C(bus i2c.Bus) (*Dev, error) {
	return NewI2CAddress(bus, PCA9685_ADDRESS)
}

func (d *Dev) init() error {
	d.SetAllPwm(0, 0)
	d.dev.Write([]byte{MODE2, OUTDRV})
	d.dev.Write([]byte{MODE1, ALLCALL})

	time.Sleep(100 * time.Millisecond)

	var mode1 byte
	err := d.dev.Tx([]byte{MODE1}, []byte{mode1})

	if err != nil {
		return err
	}

	mode1 = mode1 & ^SLEEP
	d.dev.Write([]byte{MODE1, mode1 & 0xFF})

	time.Sleep(5 * time.Millisecond)

	err = d.SetPwmFreq(50)
	return err
}

func (d *Dev) SetPwmFreq(freqHz float32) error {
	var prescaleval float32 = 25000000.0 //# 25MHz
	prescaleval /= 4096.0                //# 12-bit
	prescaleval /= freqHz
	prescaleval -= 1.0

	prescale := int(math.Floor(float64(prescaleval + 0.5)))

	var oldmode byte
	err := d.dev.Tx([]byte{MODE1}, []byte{oldmode})

	if err != nil {
		return err
	}

	newmode := (byte)((oldmode & 0x7F) | 0x10) // sleep
	d.dev.Write([]byte{MODE1, newmode})        // go to sleep
	d.dev.Write([]byte{PRESCALE, byte(prescale)})
	d.dev.Write([]byte{MODE1, oldmode})

	time.Sleep(100 * time.Millisecond)

	d.dev.Write([]byte{MODE1, (byte)(oldmode | 0x80)})
	return nil
}

func (d *Dev) SetAllPwm(on uint16, off uint16) {
	d.dev.Write([]byte{ALL_LED_ON_L, byte(on) & 0xFF})
	d.dev.Write([]byte{ALL_LED_ON_H, byte(on >> 8)})
	d.dev.Write([]byte{ALL_LED_OFF_L, byte(off) & 0xFF})
	d.dev.Write([]byte{ALL_LED_OFF_H, byte(off >> 8)})
}

func (d *Dev) SetPwm(channel int, on uint16, off uint16) {
	d.dev.Write([]byte{LED0_ON_L + byte(4*channel), byte(on) & 0xFF})
	d.dev.Write([]byte{LED0_ON_H + byte(4*channel), byte(on >> 8)})
	d.dev.Write([]byte{LED0_OFF_L + byte(4*channel), byte(off) & 0xFF})
	d.dev.Write([]byte{LED0_OFF_H + byte(4*channel), byte(off >> 8)})
}
