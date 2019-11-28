// Copyright © 2019 Michael Gruener
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package inventory

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"reflect"

	"github.com/mitchellh/mapstructure"
)

func NewTillerFromParams(params map[string]interface{}) (*tiller, error) {
	var tiller tiller

	// get the raw tls setting before decoding as the empty string
	// decodes to false and we want the default to be true
	// TODO: Fail if tls is "true" but we do not have certificate data.
	//       Find a way to test the default usecase with the cert data
	//       being encoded in the s3 stored kubeconfig
	if value, ok := params["tls"].(string); !ok || value == "" {
		params["tls"] = true
	}

	if value, ok := params["namespace"].(string); !ok || value == "" {
		params["namespace"] = "kube-system"
	}
	hook := pemDecoderHookFunc()

	decoderConfig := &mapstructure.DecoderConfig{
		DecodeHook:       hook,
		WeaklyTypedInput: true,
		Result:           &tiller,
	}

	decoder, err := mapstructure.NewDecoder(decoderConfig)
	if err != nil {
		return nil, err
	}
	err = decoder.Decode(params)
	if err != nil {
		return nil, err
	}

	return &tiller, nil
}

func tillerDecodeFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if t.Name() == "tiller" {
			var params map[string]interface{}
			err := mapstructure.Decode(data, &params)
			if err != nil {
				return data, err
			}
			settings, err := NewTillerFromParams(params)
			if err != nil {
				return data, err
			}
			return settings, nil
		}
		return data, nil
	}
}

func pemDecoderHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if t.Name() == "Certificate" {
			block, _ := pem.Decode([]byte(data.(string)))
			if block == nil {
				return data, fmt.Errorf("failed to decode pem data")
			}
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return data, fmt.Errorf("failed to parse certificate: %s", err)
			}
			return cert, nil
		}
		if t.Name() == "PrivateKey" {
			block, _ := pem.Decode([]byte(data.(string)))
			if block == nil {
				return data, fmt.Errorf("failed to decode pem data")
			}

			var err error
			var parsedKey interface{}
			if parsedKey, err = x509.ParsePKCS8PrivateKey(block.Bytes); err != nil {
				if parsedKey, err = x509.ParsePKCS1PrivateKey(block.Bytes); err != nil { // note this returns type `interface{}`
					return data, fmt.Errorf("failed to parse rsa key: %s", err)
				}
			}

			privateKey, ok := parsedKey.(*rsa.PrivateKey)
			if !ok {
				return data, fmt.Errorf("failed to parse key as rsa private key")
			}
			return privateKey, nil
		}
		return data, nil
	}
}
