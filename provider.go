package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

type Provider interface {
	StreamData(ctx context.Context, index int) ([]byte, error)
}

type MockProvider1 struct{}

func (p *MockProvider1) StreamData(ctx context.Context, index int) ([]byte, error) {
	time.Sleep(time.Second)
	// mimicking stream error
	if rand.Float32() < 0.1 {
		return nil, fmt.Errorf("provider 1 error")
	}
	chunk := DataChunk{
		Index:     index,
		Message:   fmt.Sprintf("Chunk %d from provider 1", index),
		Timestamp: time.Now(),
	}
	data, err := json.Marshal(chunk)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type MockProvider2 struct{}

func (p *MockProvider2) StreamData(ctx context.Context, index int) ([]byte, error) {
	time.Sleep(time.Second)
	chunk := DataChunk{
		Index:     index,
		Message:   fmt.Sprintf("Chunk %d from provider 2", index),
		Timestamp: time.Now(),
	}
	data, err := json.Marshal(chunk)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type ProviderManager struct {
	providers    []Provider
	reachable    map[Provider]int
	nonreachable map[Provider]string
}

func NewProviderManager() *ProviderManager {
	return &ProviderManager{
		providers: []Provider{
			&MockProvider1{},
			&MockProvider2{},
		},
		reachable: map[Provider]int{
			&MockProvider1{}: 0,
			&MockProvider2{}: 0,
		},
		nonreachable: make(map[Provider]string),
	}
}

func (pm *ProviderManager) GetProvider() Provider {
	// return pm.providers[rand.Intn(len(pm.providers))]
	for provider, _ := range pm.reachable {
		return provider
	}
	return nil
}

func (pm *ProviderManager) SwitchProvider(ctx context.Context, index int, provider Provider) ([]byte, error) {
	fmt.Printf("[*] switching router for index %d chunk\n", index)
	if _, ok := pm.reachable[provider]; ok {
		if pm.reachable[provider]+1 >= 3 {
			delete(pm.reachable, provider)
			pm.nonreachable[provider] = ""
		} else {
			pm.reachable[provider] += 1
			//  you can try with the same provider to fetch data
			// data, err := provider.StreamData(ctx, index)
			// if err != nil {
			// 	return pm.switchProvider(ctx, index, provider)
			// }
			// return data, nil
		}
	}
	newprovider := pm.GetProvider()
	if newprovider != nil {
		data, err := newprovider.StreamData(ctx, index)
		if err != nil {
			return pm.SwitchProvider(ctx, index, newprovider)
		}
		return data, nil
	}
	return nil, fmt.Errorf("error switching the provider, no reachable providers")
}

func (pm *ProviderManager) StreamData(ctx context.Context, index int) ([]byte, error) {
	provider := pm.GetProvider()
	data, err := provider.StreamData(ctx, index)
	if err != nil {
		return pm.SwitchProvider(ctx, index, provider)
	}
	return data, nil
}
