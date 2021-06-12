# Fault Tolerance Library for Golang
### Table of Contents
#### 1. [Introduction](#Introduction)
#### 2. [Usage Example](#Usage-Example)
#### 3. [Configuration Guidelines](#Configuration-Guidelines)


## Introduction
Resilix-go is a fault tolerance library designed to be flexible on its failure handling with configurable parameters to suit any fault tolerance strategy.

## Usage Example
You can take a look at this [demo project](https://github.com/alfian853/resilix-demo)

```go 
// Resilix execution will recover any panic occured in runtime.
executed, result, err = resilix.Go("thirdparty-1").ExecuteSupplier(func() interface{} {
    return CallThirdPartyApi(request)
})

// if there is result and error returned from CallThirdPartyApiWithErrorHandling(request), the result and error will be returned
// but if there is a unhandled panic occured, it returns nil(result) and error(util.UnhandledError)
executed, result, err = resilix.Go("thirdparty-1").ExecuteCheckedSupplier(func() (interface{}, error) {
    return CallThirdPartyApiWithErrorHandling(request)
})

// if panic occurred Somefunction(), it will be counted as failure
executed, err = resilix.Go("some-process-1").Execute(SomeFunction)

// if error returned SomefunctionWithErrorHandling(), it also will be counted as failure
executed, err = resilix.Go("some-process-2").ExecuteChecked(SomeFunctionWithErrorHandling)
```

if you want to customize the configuration please read the [Configuration Guidelines](##Configuration)

## Configuration Guidelines
Configuration example:
```go
cfg := config.NewConfiguration()
cfg.ErrorThreshold = 0.2
cfg.SlidingWindowStrategy = config.SwStrategy_CountBased
cfg.RetryStrategy = config.RetryStrategy_Pessimistic
cfg.WaitDurationInOpenState = 5000
cfg.SlidingWindowMaxSize = 10
cfg.NumberOfRetryInHalfOpenState = 5
cfg.MinimumCallToEvaluate = 2

// register the "foo" contextKey to resilix
resilix.Register("foo", cfg)

// execution on behalf the "foo" contextKey
resilix.Go("foo").Execute(SomeFunc)
```

### Config.ErrorThreshold
Configures the error threshold in percentage in close state and half-open state (for retry process).
when the error rate is greater or equal to the threshold, the circuit will trip to open state.

|default value|`0.5`|
|:---:|:---|
|possible value|from `0.0` to `1.0`|


### Config.SlidingWindowStrategy
Configures the sliding-window aggregation type.

- SwStrategy_CountBased: aggregate error rate by the last ***n*** records.
- SwStrategy_TimeBased: aggregate error rate from the last ***t*** milliseconds.

|default value|`SwStrategy_CountBased`|
|:---:|:---|
|possible value|`SwStrategy_CountBased` or `SwStrategy_TimeBased`|

### Config.SlidingWindowMaxSize
Configures the maximum size for `SlidingWindowStrategy=SwStrategy_CountBased`

|default value|`20`|
|:---:|:---|

### Config.SlidingWindowTimeRange
Configures the window time range for `slidingWindowStrategy=SwStrategy_TimeBased`

|default value|`15000` ms|
|:---:|:---|


### Config.RetryStrategy
Configures the retry strategy.
<br>

- RetryStrategy_Pessimistic: Do retry once at a time as long as the error threshold hasn't reached, and the retry limit (`retryCountInHalfOpenState`) hasn't exceeded.
- RetryStrategy_Optimistic: Do retry as long as the error threshold hasn't reached, and the retry limit (`retryCountInHalfOpenState`) hasn't exceeded.
it will result in more errors if multiple retries happened at the same time and get errors.

|default value|`RetryStrategy_Pessimistic`|
|:---:|:---|
|possible value|`RetryStrategy_Optimistic` or `RetryStrategy_Pessimistic`|

### Config.NumberOfRetryInHalfOpenState
Configures the number of minimum retry in the half-open state before deciding which state it belongs to. 

|default value|`10`|
|:---:|:------------------------|

### Config.MinimumCallToEvaluate
Configures the minimum records to be evaluated to prevent 100% error rate if the
first call is failed.

|default value|`10`|
|:---:|:------------------------|

