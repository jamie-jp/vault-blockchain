// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at

//   http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package path

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"

	"github.com/intech/vault-blockchain/internal/model"
	"github.com/intech/vault-blockchain/pkg/utils"
)

type signPathConfig struct {
	basePathConfig
}

func (s *signPathConfig) getPattern() string {
	return fmt.Sprintf("keys/%s/sign", framework.GenericNameRegex("name"))
}

func (s *signPathConfig) getHelpSynopsis() string {
	return "Sign data"
}

func (s *signPathConfig) getFields() map[string]*framework.FieldSchema {
	return map[string]*framework.FieldSchema{
		"name": {
			Type: framework.TypeString,
		},
		"data": {
			Type:        framework.TypeString,
			Description: "The data as hash for sign.",
		},
	}
}

func (s *signPathConfig) getCallbacks() map[logical.Operation]framework.OperationFunc {
	return map[logical.Operation]framework.OperationFunc{
		logical.CreateOperation: s.sign,
	}
}

func (s *signPathConfig) sign(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	dataWrapper := model.NewFieldDataWrapper(data)
	name, err := dataWrapper.MustGetString("name")
	if err != nil {
		return nil, err
	}
	key, err := s.readKey(ctx, req, name)
	if err != nil {
		return nil, fmt.Errorf("error reading key, %v", err)
	}
	// get data to sign
	dataToSign, err := dataWrapper.MustGetString("data")
	if err != nil {
		return nil, err
	}
	dataBytes, err := hexutil.Decode("0x" + dataToSign)
	if err != nil {
		return nil, err
	}
	// get private ecdsa key from key for signing data
	privateKey, err := key.GetPrivateKeyECDSA()
	if err != nil {
		return nil, fmt.Errorf("error reconstructing private key, %v", err)
	}
	defer utils.ZeroKey(privateKey)
	// sign data
	signature, err := crypto.Sign(dataBytes, privateKey)
	if err != nil {
		return nil, err
	}
	return &logical.Response{
		Data: map[string]interface{}{
			"signature": hexutil.Encode(signature),
			"address":   key.PublicKeyStr,
		},
	}, nil
}
