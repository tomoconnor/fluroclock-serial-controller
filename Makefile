build:
	go build -o serial-controller .

run:
	go build -o serial-controller .
	./serial-controller

clean:
	go clean
	rm serial-controller

install:
	install serial-controller /usr/local/bin/serial-controller
	install -m 664 serial-controller.service /etc/systemd/system/serial-controller.service
	install -m 644 tz.json /etc/fclock-tz.json
	systemctl enable serial-controller.service
	systemctl start serial-controller.service
	systemctl daemon-reload
