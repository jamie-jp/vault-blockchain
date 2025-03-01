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

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"

	"github.com/intech/vault-blockchain/internal/model"
)

type basePathConfig struct {
	config
}

func (b basePathConfig) getExistenceFunc() framework.ExistenceFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
		entry, err := req.Storage.Get(ctx, req.Path)
		if err != nil {
			return false, fmt.Errorf("existence check failed, %v", err)
		}
		return entry != nil, nil
	}
}

func (b *basePathConfig) readKey(ctx context.Context, req *logical.Request, name string) (*model.Key, error) {
	path := fmt.Sprintf("keys/%s", name)
	entry, err := req.Storage.Get(ctx, path)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, fmt.Errorf("entry not existed at %v", path)
	}
	var key *model.Key
	err = entry.DecodeJSON(&key)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize key at %s", path)
	}
	if key == nil {
		return nil, fmt.Errorf("key not existed at %s", path)
	}
	return key, nil
}
