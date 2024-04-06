.PHONY: build init remove

build:
	make -C serverless/account_summarizers build

init:
	make -C serverless/account_summarizers init STACK=$(STACK)

remove:
	make -C serverless/account_summarizers remove STACK=$(STACK)

