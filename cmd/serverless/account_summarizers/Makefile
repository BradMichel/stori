.PHONY: init remove logs-initiator logs-summarizer logs-processor info invoke-initiator invoke-summarizer invoke-notificator

init:
	../../../node_modules/.bin/serverless --stage $(STACK) deploy

remove:
	../../../node_modules/.bin/serverless --stage $(STACK) remove

info:
	../../../node_modules/.bin/serverless info --stage $(STACK)

invoke-initiator:
	../../../node_modules/.bin/serverless invoke -f initiator-v1 --stage $(STACK)

invoke-summarizer:
	../../../node_modules/.bin/serverless invoke -f summarizer-v1 --stage $(STACK)

invoke-notificator:
	../../../node_modules/.bin/serverless invoke -f email-notificator-v1 --stage $(STACK)

logs-initiator:
	../../../node_modules/.bin/serverless logs -f initiator-v1 -t --stage $(STACK)

logs-summarizer:
	../../../node_modules/.bin/serverless logs -f summarizer-v1 -t --stage $(STACK)

logs-notificator:
	../../../node_modules/.bin/serverless logs -f email-notificator-v1 -t --stage $(STACK)
