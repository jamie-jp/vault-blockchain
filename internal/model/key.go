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

package model

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"

	"github.com/bsostech/vault-blockchain/pkg/utils"
)

// Key
type Key struct {
	PrivateKeyStr string `json:"private_key"`
	PublicKeyStr  string `json:"public_key"` // public key derived from the private key
}

// NewKey returns Key Json
func NewKey(privateKeyStr string, publicKeyStr string) *Key {
	return &Key{
		PrivateKeyStr: privateKeyStr,
		PublicKeyStr:  publicKeyStr,
	}
}

// GetPrivateKeyECDSA key for signing data
func (a *Key) GetPrivateKeyECDSA() (*ecdsa.PrivateKey, error) {
	// Get private key from key
	return crypto.HexToECDSA(a.PrivateKeyStr)
}

// GetPublicKeyECDSA key for validating signature
func (a *Key) GetPublicKeyECDSA() (*ecdsa.PublicKey, error) {
	privateKeyECDSA, err := a.GetPrivateKeyECDSA()
	if err != nil {
		return nil, err
	}
	defer utils.ZeroKey(privateKeyECDSA)
	publicKey := privateKeyECDSA.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("error casting public key to ECDSA")
	}
	return publicKeyECDSA, nil
}

