package middleware

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/alvarowolfx/golang-voice-iot/mearm"
)

const (
	mqttHost = "mqtt.googleapis.com"
	mqttPort = "8883"
)

type CloudIoTCoreRobotArm struct {
	robotArm   *mearm.RobotArm
	deviceID   string
	projectID  string
	region     string
	registryID string

	privateKey *rsa.PrivateKey
	tlsConfig  *tls.Config
	opts       *mqtt.ClientOptions

	configTopic    string
	telemetryTopic string
	stateTopic     string

	client mqtt.Client
}

func NewCloudIoTCoreRobotArm(robotArm *mearm.RobotArm,
	deviceID,
	projectID,
	region,
	registryID,
	privateKey,
	certsCa string) (*CloudIoTCoreRobotArm, error) {

	keyBytes, err := ioutil.ReadFile(privateKey)
	if err != nil {
		return nil, err
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyBytes)
	if err != nil {
		return nil, err
	}

	certpool := x509.NewCertPool()
	pemCerts, err := ioutil.ReadFile(certsCa)
	if err != nil {
		return nil, err
	}
	certpool.AppendCertsFromPEM(pemCerts)

	config := &tls.Config{
		RootCAs:            certpool,
		ClientAuth:         tls.NoClientCert,
		ClientCAs:          nil,
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{},
		MinVersion:         tls.VersionTLS12,
	}

	return &CloudIoTCoreRobotArm{
		robotArm:   robotArm,
		deviceID:   deviceID,
		projectID:  projectID,
		region:     region,
		registryID: registryID,
		tlsConfig:  config,
		privateKey: key,
	}, nil
}

func (ciot *CloudIoTCoreRobotArm) Init() {
	clientID := fmt.Sprintf("projects/%v/locations/%v/registries/%v/devices/%v",
		ciot.projectID,
		ciot.region,
		ciot.registryID,
		ciot.deviceID,
	)

	ciot.opts = mqtt.NewClientOptions()

	broker := fmt.Sprintf("ssl://%v:%v", mqttHost, mqttPort)

	ciot.opts.
		AddBroker(broker).
		SetClientID(clientID).
		SetTLSConfig(ciot.tlsConfig).
		SetUsername("unused").
		SetProtocolVersion(4)

	ciot.configTopic = fmt.Sprintf("/devices/%v/config", ciot.deviceID)
	ciot.telemetryTopic = fmt.Sprintf("/devices/%v/events", ciot.deviceID)
	ciot.stateTopic = fmt.Sprintf("/devices/%v/state", ciot.deviceID)
}

func (ciot *CloudIoTCoreRobotArm) Start() {

	if ciot.client != nil {
		if ciot.client.IsConnected() {
			ciot.Stop()
		}
	}

	password, err := ciot.generatePassword()
	if err != nil {
		log.Fatal(err)
	}

	ciot.opts.SetPassword(password)

	first := true
	ciot.opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		if first { // Skip first retained config
			first = false
			return
		}
		fmt.Printf("[handler] Topic: %v\n", msg.Topic())
		fmt.Printf("[handler] Payload: %v\n", msg.Payload())

		payload := msg.Payload()
		var config map[string]string
		err := json.Unmarshal(payload, &config)
		if err != nil {
			fmt.Printf("[handler] Error decoding payload")
		}

		fmt.Printf("[handler] Parse Payload: %v\n", config)

		if moveElbow, ok := config["moveelbow"]; ok {
			intValue, err := strconv.Atoi(moveElbow)
			if err == nil {
				ciot.robotArm.MoveElbow(intValue)
			}
		}

		if elbow, ok := config["elbow"]; ok {
			intValue, err := strconv.Atoi(elbow)
			if err == nil {
				ciot.robotArm.SetElbow(intValue)
			}
		}

		if moveShoulder, ok := config["moveshoulder"]; ok {
			intValue, err := strconv.Atoi(moveShoulder)
			if err == nil {
				ciot.robotArm.MoveShoulder(intValue)
			}
		}

		if shoulder, ok := config["shoulder"]; ok {
			intValue, err := strconv.Atoi(shoulder)
			if err == nil {
				ciot.robotArm.SetShoulder(intValue)
			}
		}

		if moveBase, ok := config["movebase"]; ok {
			intValue, err := strconv.Atoi(moveBase)
			if err == nil {
				ciot.robotArm.MoveBase(intValue)
			}
		}

		if base, ok := config["base"]; ok {
			intValue, err := strconv.Atoi(base)
			if err == nil {
				ciot.robotArm.SetBase(intValue)
			}
		}

		if grip, ok := config["grip"]; ok {
			grip = strings.ToLower(grip)
			if grip == "open" {
				ciot.robotArm.OpenGrip()
			} else if grip == "close" {
				ciot.robotArm.CloseGrip()
			}
		}
	})

	ciot.client = mqtt.NewClient(ciot.opts)
	if token := ciot.client.Connect(); token.WaitTimeout(5*time.Second) && token.Error() != nil {
		log.Fatal(token.Error())
	}

	ciot.client.Subscribe(ciot.configTopic, 1, nil)
}

func (ciot *CloudIoTCoreRobotArm) generatePassword() (string, error) {
	token := jwt.New(jwt.SigningMethodRS256)
	token.Claims = jwt.StandardClaims{
		Audience:  ciot.projectID,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	}

	tokenString, err := token.SignedString(ciot.privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (ciot CloudIoTCoreRobotArm) Stop() {
	ciot.client.Disconnect(1000)
}
