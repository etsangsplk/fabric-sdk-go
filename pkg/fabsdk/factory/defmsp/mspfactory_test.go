/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package defmsp

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/test/mockmsp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	cryptosuiteImpl "github.com/hyperledger/fabric-sdk-go/pkg/core/cryptosuite"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/mocks"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk/factory/defcore"
	mspimpl "github.com/hyperledger/fabric-sdk-go/pkg/msp"
)

func TestCreateUserStore(t *testing.T) {
	factory := NewProviderFactory()

	config := mocks.NewMockIdentityConfig()

	userStore, err := factory.CreateUserStore(config)
	if err != nil {
		t.Fatalf("Unexpected error creating state store %v", err)
	}

	_, ok := userStore.(*mspimpl.CertFileUserStore)
	if !ok {
		t.Fatalf("Unexpected state store created")
	}
}

func newMockUserStore(t *testing.T) msp.UserStore {
	factory := NewProviderFactory()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockConfig := mockmsp.NewMockIdentityConfig(mockCtrl)

	mockClientConfig := msp.ClientConfig{
		CredentialStore: msp.CredentialStoreType{
			Path: "/tmp/fabsdkgo_test/store",
		},
	}
	mockConfig.EXPECT().Client().Return(&mockClientConfig, nil)

	userStore, err := factory.CreateUserStore(mockConfig)
	if err != nil {
		t.Fatalf("Unexpected error creating user store %v", err)
	}
	return userStore
}
func TestCreateUserStoreByConfig(t *testing.T) {
	userStore := newMockUserStore(t)

	_, ok := userStore.(*mspimpl.CertFileUserStore)
	if !ok {
		t.Fatalf("Unexpected user store created")
	}
}

func TestCreateUserStoreEmptyConfig(t *testing.T) {
	factory := NewProviderFactory()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockConfig := mockmsp.NewMockIdentityConfig(mockCtrl)

	mockClientConfig := msp.ClientConfig{}
	mockConfig.EXPECT().Client().Return(&mockClientConfig, nil)

	_, err := factory.CreateUserStore(mockConfig)
	if err == nil {
		t.Fatal("Expected error creating user store")
	}
}

func TestCreateUserStoreFailConfig(t *testing.T) {
	factory := NewProviderFactory()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockConfig := mockmsp.NewMockIdentityConfig(mockCtrl)

	mockConfig.EXPECT().Client().Return(nil, errors.New("error"))

	_, err := factory.CreateUserStore(mockConfig)
	if err == nil {
		t.Fatal("Expected error creating user store")
	}
}

func TestCreateIdentityManager(t *testing.T) {

	coreFactory := defcore.NewProviderFactory()

	configBackend, err := config.FromFile("../../../../test/fixtures/config/config_test.yaml")()
	if err != nil {
		t.Fatalf(err.Error())
	}

	cryptoCfg := cryptosuiteImpl.ConfigFromBackend(configBackend)
	if err != nil {
		t.Fatalf(err.Error())
	}

	endpointCfg, err := fab.ConfigFromBackend(configBackend)
	if err != nil {
		t.Fatalf(err.Error())
	}

	identityCfg, err := mspimpl.ConfigFromBackend(configBackend)
	if err != nil {
		t.Fatalf(err.Error())
	}

	cryptosuite, err := coreFactory.CreateCryptoSuiteProvider(cryptoCfg)
	if err != nil {
		t.Fatalf("Unexpected error creating cryptosuite provider %v", err)
	}

	factory := NewProviderFactory()
	userStore, err := factory.CreateUserStore(identityCfg)
	if err != nil {
		t.Fatalf("Unexpected error creating user store %v", err)
	}

	provider, err := factory.CreateIdentityManagerProvider(endpointCfg, cryptosuite, userStore)
	if err != nil {
		t.Fatalf("Unexpected error creating provider %v", err)
	}

	mgr, ok := provider.IdentityManager("Org1")
	if !ok {
		t.Fatalf("Unexpected error creating identity manager %v", err)
	}

	_, ok = mgr.(msp.IdentityManager)
	if !ok {
		t.Fatalf("Unexpected signing manager created")
	}
}
