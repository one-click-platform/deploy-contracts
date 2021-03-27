package config

import (
	"fmt"
	"net/url"
	"reflect"

	"github.com/ethereum/go-ethereum/accounts"

	"github.com/ethereum/go-ethereum/accounts/keystore"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type Node interface {
	Client() *ethclient.Client
	Account() accounts.Account
	KeyStore() *keystore.KeyStore
}

func NewNode(getter kv.Getter) Node {
	return &node{getter: getter}
}

type node struct {
	getter kv.Getter
	cfg    comfig.Once
}

func (n *node) Client() *ethclient.Client {
	return n.config().Client
}

func (n *node) Account() accounts.Account {
	return n.config().Account
}

func (n *node) KeyStore() *keystore.KeyStore {
	return n.config().KeyStore
}

type ethConfig struct {
	KeyStore *keystore.KeyStore
	Account  accounts.Account
	Client   *ethclient.Client
}

func (n *node) config() ethConfig {
	return n.cfg.Do(func() interface{} {
		var cfg struct {
			Endpoint *url.URL `fig:"url,required"`
			KeyDir   string   `fig:"keydir,required"`
			Address  string   `fig:"address"`
			Password string   `fig:"password,required"`
		}
		err := figure.
			Out(&cfg).
			With(figure.BaseHooks, hooks).
			From(kv.MustGetStringMap(n.getter, "eth")).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to parse eth config entry"))
		}

		var result ethConfig
		result.Client, err = ethclient.Dial(cfg.Endpoint.String())
		if err != nil {
			panic(errors.Wrap(err, "failed to connect to eth client"))
		}

		result.KeyStore = keystore.NewKeyStore(cfg.KeyDir, keystore.StandardScryptN, keystore.StandardScryptP)
		result.Account = result.KeyStore.Accounts()[0]
		if cfg.Address != "" {
			result.Account = accounts.Account{Address: common.HexToAddress(cfg.Address)}
		}

		err = result.KeyStore.Unlock(result.Account, cfg.Password)
		if err != nil {
			panic(errors.Wrap(err, "failed to unlock account"))
		}

		return result
	}).(ethConfig)
}

var hooks = figure.Hooks{
	"common.Address": func(value interface{}) (reflect.Value, error) {
		switch v := value.(type) {
		case string:
			if !common.IsHexAddress(v) {
				// provide value does not look like valid address
				return reflect.Value{}, errors.New("invalid address")
			}
			return reflect.ValueOf(common.HexToAddress(v)), nil
		default:
			return reflect.Value{}, fmt.Errorf("unsupported conversion from %T", value)
		}
	},
	"*ecdsa.PrivateKey": func(raw interface{}) (reflect.Value, error) {
		switch value := raw.(type) {
		case string:
			kp, err := crypto.HexToECDSA(value)
			if err != nil {
				return reflect.Value{}, errors.Wrap(err, "failed to init keypair")
			}
			return reflect.ValueOf(kp), nil
		default:
			return reflect.Value{}, fmt.Errorf("cant init keypair from type: %T", value)
		}
	},
}
