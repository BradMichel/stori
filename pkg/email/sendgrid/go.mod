module github.com/BradMichel/stori/pkg/email/sendgrid

go 1.22

replace github.com/BradMichel/stori => ../../../

require (
	github.com/BradMichel/stori v0.0.0-00010101000000-000000000000
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/sendgrid/sendgrid-go v3.14.0+incompatible

)

require github.com/sendgrid/rest v2.6.9+incompatible // indirect
