## Usage

Refer to our [docs](https://boxyhq.com/docs/retraced/sdks/golang) for details on the usage of this SDK.

## Local development

## Install the dependencies to run the tests

```
go get
```

## Run tests

```
go test -v
```

## Run SDK client test

Ensure you are running Retraced, then run:

```
go test -timeout 30s -run ^TestClientQuery$ github.com/retracedhq/retraced-go/tests -count=1
```
