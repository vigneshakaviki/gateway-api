/*
Copyright The Kubernetes Authors.

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

package main

import (
	"fmt"
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func TestHTTPRouteRequestHeaderModifierNames(t *testing.T) {
	tests := []struct {
		name       string
		add        []gatewayv1.HTTPHeader
		wantErrors []string
	}{
		{
			name: "valid distinct header names in add",
			add: []gatewayv1.HTTPHeader{
				{Name: "my-header-name", Value: "my-header-value"},
				{Name: "other-header-name", Value: "other-value"},
			},
		},
		{
			name: "invalid exact duplicate header name in add",
			add: []gatewayv1.HTTPHeader{
				{Name: "my-header-name", Value: "my-header-value"},
				{Name: "my-header-name", Value: "my-header-value2"},
			},
			wantErrors: []string{`Duplicate value: {"name":"my-header-name"}`},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			route := &gatewayv1.HTTPRoute{
				Name:      fmt.Sprintf("header-names-%v", time.Now().UnixNano()),
				Namespace: metav1.NamespaceDefault,
				Spec: gatewayv1.HTTPRouteSpec{
					Rules: []gatewayv1.HTTPRouteRule{{
						Filters: []gatewayv1.HTTPRouteFilter{{
							Type: gatewayv1.HTTPRouteFilterRequestHeaderModifier,
							RequestHeaderModifier: &gatewayv1.HTTPHeaderFilter{
								Add: tc.add,
							},
						}},
					}},
				},
			}
			validateHTTPRoute(t, route, tc.wantErrors)
		})
	}
}
