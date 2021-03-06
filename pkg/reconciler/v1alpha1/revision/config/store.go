/*
Copyright 2018 The Knative Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

import (
	"context"

	pkglogging "github.com/knative/pkg/logging"
	"github.com/knative/serving/pkg/autoscaler"
	"github.com/knative/serving/pkg/logging"
	"github.com/knative/serving/pkg/reconciler/config"
)

type cfgKey struct{}

// +k8s:deepcopy-gen=false
type Config struct {
	Controller    *Controller
	Network       *Network
	Observability *Observability
	Logging       *pkglogging.Config
	Autoscaler    *autoscaler.Config
}

func FromContext(ctx context.Context) *Config {
	return ctx.Value(cfgKey{}).(*Config)
}

func ToContext(ctx context.Context, c *Config) context.Context {
	return context.WithValue(ctx, cfgKey{}, c)
}

// +k8s:deepcopy-gen=false
type Store struct {
	*config.UntypedStore
}

func NewStore(logger config.Logger) *Store {
	store := &Store{
		UntypedStore: config.NewUntypedStore(
			"revision",
			logger,
			config.Constructors{
				ControllerConfigName:    NewControllerConfigFromConfigMap,
				NetworkConfigName:       NewNetworkFromConfigMap,
				ObservabilityConfigName: NewObservabilityFromConfigMap,
				autoscaler.ConfigName:   autoscaler.NewConfigFromConfigMap,
				logging.ConfigName:      logging.NewConfigFromConfigMap,
			},
		),
	}

	return store
}

func (s *Store) ToContext(ctx context.Context) context.Context {
	return ToContext(ctx, s.Load())
}

func (s *Store) Load() *Config {
	return &Config{
		Controller:    s.UntypedLoad(ControllerConfigName).(*Controller).DeepCopy(),
		Network:       s.UntypedLoad(NetworkConfigName).(*Network).DeepCopy(),
		Observability: s.UntypedLoad(ObservabilityConfigName).(*Observability).DeepCopy(),
		Logging:       s.UntypedLoad(logging.ConfigName).(*pkglogging.Config).DeepCopy(),
		Autoscaler:    s.UntypedLoad(autoscaler.ConfigName).(*autoscaler.Config).DeepCopy(),
	}
}
