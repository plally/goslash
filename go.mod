module github.com/plally/goslash

go 1.15

require (
	github.com/aws/aws-lambda-go v1.21.0
	github.com/bwmarrin/discordgo v0.23.3-0.20210327033043-f637c37ba2f0
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/kr/pretty v0.1.0 // indirect
	github.com/sirupsen/logrus v1.7.0
	github.com/stretchr/testify v1.6.1
	golang.org/x/crypto v0.0.0-20190605123033-f99c8df09eb5 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
)

replace github.com/bwmarrin/discordgo => github.com/plally/discordgo v0.23.4
