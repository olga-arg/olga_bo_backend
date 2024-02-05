module go-lambda-post-confirmation

go 1.18

require (
	commons v0.0.0-00010101000000-000000000000
	github.com/aws/aws-lambda-go v1.45.0
	github.com/jinzhu/gorm v1.9.16
)

require (
	github.com/aws/aws-sdk-go v1.50.3 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/lib/pq v1.1.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/stretchr/testify v1.8.4 // indirect
)

replace commons => ../commons
