@echo off
if "%1%"=="" (
    set cmd=container
) else (
    set cmd=%1%
)
if "%cmd%"=="container" (
    docker-compose build test
    goto :eof
)
if "%cmd%"=="test" (
    docker-compose run --rm test
    goto :eof
)
if "%cmd%"=="coverage" (
    docker-compose run --rm test sh -c "go test -coverprofile=coverage.out ./... && go tool cover -html coverage.out -o coverage.html"
    coverage.html
    goto :eof
)
if "%cmd%"=="clean" (
	docker image rm go-shopify
	del coverage.html
	del coverage.out
    goto :eof
)
