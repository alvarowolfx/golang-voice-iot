build:
	GOOS=linux GOARCH=arm go build -ldflags="-s -w" -o robot-arm main.go 
	
copy: 
	rsync -P -a robot-arm pi@orangepizero.local:/home/pi/go

copy-all:
	rsync -P -a orangepizero.key.pem start-robot-arm.sh robot-arm pi@orangepizero.local:/home/pi/go