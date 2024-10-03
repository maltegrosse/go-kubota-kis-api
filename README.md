# Kubota Intelligence Solutions (KIS) API

Golang KIS API bindings.

The Kubota Intelligence Solution (KIS) provides flexible access to machine operation data via their official API.
KIS is available for machines mainly outside of Japan (e.g. South East Asisa). 

## Usage
Fetch the package

```
go get github.com/maltegrosse/go-kubota-kis-api
```

Call the constructor including the API endpoint (e.g. https://someweb-api-kis.net) including private and public key, which can be obtained at the developer app console.

```
k, err := kis.NewKIS("PUBLIC-KEY", "PRIVATE-KEY", "https://someweb-api-kis.net")
	if err != nil {
		panic(err)
	}
```

The application should automatically refresh the access token 


Now all major endpoints can be called by public functions, e.g. the latest position by machine ID:
```
mId := "SOME-MACHINE-UUID"
	pos, err := k.GetLastPositionByMachineUUID(mId, "")
	if err != nil {
		panic(err)
	}
```

Additional examples can be found at `/examples/main.go`

## Limitation
The current status of the KIS API is still under development and can be changed. Not all functions are tested. The API wrapper is based on Kubota API Service. version 1.0.1 [December 07, 2023]


## License
**[MIT license](http://opensource.org/licenses/mit-license.php)**

Copyright 2024 © Malte Grosse.