
BUILD_NAME = ipinfo-api

docker-build:
	@echo 'Start docker-compose ...'
	@docker-compose up -d --build || (echo "[!] Docker-compose failed $$?"; exit 1)
	@echo '[+] Service started'

start:
	@docker-compose start

stop:
	@docker-compose stop

clean:
	@rm -rf $(BUILD_NAME)

.PHONY: docker-build start stop clean