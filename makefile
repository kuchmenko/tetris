run:
	go run -tags wayland ./src/.

.PHONY: run run-debug

run-debug:
	dlv debug ./src/. --build-flags="-tags=wayland" --headless --listen=127.0.0.1:38697 --accept-multiclient
