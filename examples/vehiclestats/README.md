### Build against live, local roadtrip package from this checked-out repo
```
go mod edit -replace=github.com/nugget/roadtrip-go/roadtrip="../../roadtrip"
```

### Build against current public release of the roadtrip package
```
go mod edit -dropreplace=github.com/nugget/roadtrip-go/roadtrip
```

### Convenient run command for testing
```
clearbuffer && go run . --file ../CSV/Example\ Vehicle.csv
```
