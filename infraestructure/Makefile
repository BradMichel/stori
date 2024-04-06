.PHONY: npmi init remove local remove-local refresh

npmi:
	npm install

init: npmi
	@if [ "$(STACK)" = "local" ]; then \
		make local; \
	else \
		pulumi up --stack $(STACK); \
	fi

remove: npmi
	@if [ "$(STACK)" = "local" ]; then \
		make remove-local; \
	else \
		pulumi destroy --stack $(STACK); \
	fi

local:
	docker-compose up --detach
	SLEEP 5
	pulumi login --local
	PULUMI_CONFIG_PASSPHRASE=your_passphrase pulumi refresh --stack local --yes
	PULUMI_CONFIG_PASSPHRASE=your_passphrase pulumi up --stack local --yes

remove-local:
	PULUMI_CONFIG_PASSPHRASE=your_passphrase pulumi destroy --stack local --yes
	docker-compose down

refresh: npmi
	PULUMI_CONFIG_PASSPHRASE=your_passphrase pulumi refresh --stack $(STACK) --yes