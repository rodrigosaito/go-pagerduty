# go-pagerduty
A simple Pager Duty client for Go

## Installation

```
$ go get github.com/rodrigosaito/go-pagerduty
```

## Getting Started

### Import package

```go
import "github.com/rodrigosaito/go-pagerduty/pagerduty"
```

### Triggering an incident

```go
pd := pagerduty.NewPagerDuty("SERVICE_TOKEN")
pd.Trigger(pagerduty.Trigger{
  Description: message, 
  IncidentKey: "incidentKey", // If you don't send an incidentKey random one will be generated by PagerDuty
})
```
