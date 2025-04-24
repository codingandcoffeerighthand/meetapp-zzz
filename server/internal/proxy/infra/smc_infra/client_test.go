package smc_infra_test

import (
	"context"
	"log"
	"proxy-srv/internal/proxy/configs"
	"proxy-srv/internal/proxy/infra/smc_infra"
	"testing"
)

func Test(t *testing.T) {
	cfg, err := configs.NewConfig(configs.ConfigPath(""))
	if err != nil {
		t.Fatal(err)
	}
	cl, err := smc_infra.NewSMCInfra(&cfg)
	if err != nil {
		t.Fatal(err)
	}

	c, _, err := cl.SubCreateRoomEvent(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	log.Println("subCreateRoomEvent")
	e := <-c
	if e == nil {
		t.Fatal("nil")
	}
	log.Println(e)
	t.Log(e)
}
