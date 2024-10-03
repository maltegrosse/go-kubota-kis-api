package main

import (
	"fmt"
	"time"

	kis "github.com/maltegrosse/go-kubota-kis-api"
)

func main() {

	k, err := kis.NewKIS("publicKEY", "privateKEX", "https://some-KIS-API-Endpoint.net")
	if err != nil {
		panic(err)
	}
	mId := "SOME-MACHINE-UUID"
	pos, err := k.GetLastPositionByMachineUUID(mId, "")
	if err != nil {
		panic(err)
	}
	fmt.Println(pos)

	_, err = k.GetHistoricalPositionByMachineUUID(mId, "", time.Now().Add(-24*time.Hour), time.Now())
	if err != nil {
		panic(err)
	}
	machine, err := k.GetMachineByMachineUUID(mId, "")
	if err != nil {
		panic(err)
	}
	fmt.Println(machine)
	reg, err := k.GetRegistryByMachineUUID(mId, "")
	if err != nil {
		panic(err)
	}
	fmt.Println(reg)

	alarm, err := k.GetHistoricalAlarmByMachineUUID(mId, "", time.Now().Add(-24*time.Hour), time.Now())
	if err != nil {
		panic(err)
	}
	fmt.Println(alarm)

	measure, err := k.GetHistoricalMeasureByMachineUUID(mId, "", time.Now().Add(-24*time.Hour), time.Now())
	if err != nil {
		panic(err)
	}
	fmt.Println(measure)
}
