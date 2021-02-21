module github.com/wolfeidau/lambda-go-extras/example

go 1.16

replace (
	github.com/wolfeidau/lambda-go-extras => ../
	github.com/wolfeidau/lambda-go-extras/middleware/raw => ../middleware/raw
	github.com/wolfeidau/lambda-go-extras/middleware/zerolog => ../middleware/zerolog
)

require (
	github.com/aws/aws-lambda-go v1.22.0
	github.com/rs/zerolog v1.20.0
	github.com/wolfeidau/lambda-go-extras v1.2.2-0.20210221075335-fb39fb29667d
	github.com/wolfeidau/lambda-go-extras/middleware/raw v0.0.0-20210221075335-fb39fb29667d
	github.com/wolfeidau/lambda-go-extras/middleware/zerolog v0.0.0-20210221075335-fb39fb29667d
)
