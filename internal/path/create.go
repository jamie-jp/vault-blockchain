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
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"

	"github.com/intech/vault-blockchain/internal/model"
	"github.com/intech/vault-blockchain/pkg/utils"
)

type createPathConfig struct {
	basePathConfig
}

func (c *createPathConfig) getPattern() string {
	return "create"
}

func (c *createPathConfig) getHelpSynopsis() string {
	return "Create a private key"
}

func (c *createPathConfig) getCallbacks() map[logical.Operation]framework.OperationFunc {
	return map[logical.Operation]framework.OperationFunc{
		logical.CreateOperation: c.create,
	}
}

func (c *createPathConfig) create(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// create private key
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}
	defer utils.ZeroKey(privateKey)
	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyString := hexutil.Encode(privateKeyBytes)[2:]
	// create public key
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("error casting public key to ECDSA")
	}
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	publicKeyString := hexutil.Encode(publicKeyBytes)[4:]
	// store key info
	key := model.NewKey(privateKeyString, publicKeyString)
	entry, err := logical.StorageEntryJSON(fmt.Sprintf("keys/%s", key.PublicKeyStr), key)
	if err != nil {
		return nil, err
	}
	err = req.Storage.Put(ctx, entry)
	if err != nil {
		return nil, err
	}
	return &logical.Response{
		Data: map[string]interface{}{
			"publicKey": key.PublicKeyStr,
			"path": fmt.Sprintf("keys/%s", key.PublicKeyStr),
		},
	}, nil
}
