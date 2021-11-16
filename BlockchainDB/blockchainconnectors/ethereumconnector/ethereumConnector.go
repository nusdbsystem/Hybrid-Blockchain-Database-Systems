package connectors

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	KVStore "hybrid/BlockchainDB/storage/ethereum/contracts/KVStore"
	"hybrid/BlockchainDB/storage/redis"
	TxMgr "hybrid/BlockchainDB/transactionMgr"
)

type EthereumConnector struct {
	Client *ethclient.Client
	KV     *KVStore.Store
	Auth   *bind.TransactOpts
	Hexkey string
	Redis  *redis.RedisKV
	TxMgr  *TxMgr.TransactionMgr
}

func keyToByte32(key string) [32]byte {
	var key32 [32]byte
	copy(key32[:], key)
	return key32
}

func (ethereumConn *EthereumConnector) Read(ctx context.Context, key string) (string, error) {
	//auth, err := ethereumConn.bindTransactOpts(ctx, key) //no bind
	//result, err := ethereumConn.KV.Get(auth, []byte(key))
	result, err := ethereumConn.KV.Items(nil, keyToByte32(key))
	if err != nil {
		log.Println("error EthereumConnector Read ", err)
		return "", err
	}
	//return string(result.Data()), nil
	return string(result), nil
}

func (ethereumConn *EthereumConnector) Write(ctx context.Context, key, value string) (string, error) {

	auth, err := ethereumConn.bindTransactOpts(ctx, key)
	tx, err := ethereumConn.KV.Set(auth, keyToByte32(key), []byte(value))
	//tx, err := ethereumConn.KV.Set(auth, key, value)
	if err != nil {
		log.Println("error EthereumConnector Write ", err)
		return "", err
	}
	txid := tx.Hash().Hex()
	// Enable verify function
	if ethereumConn.Redis != nil {
		err = ethereumConn.Redis.Set([]byte("set_"+key), []byte(txid))
		if err != nil {
			log.Println("error store txid for Set opt verification", err)
		}
		err = ethereumConn.Redis.Set([]byte("get_"+key), []byte(value))
		if err != nil {
			log.Println("error store value for Get opt verification", err)
		}
	}

	return txid, nil
}

func (ethereumConn *EthereumConnector) Verify(ctx context.Context, opt, key string) (bool, error) {

	switch opt {
	case "set": //check transaction status by txid
		txid, err := ethereumConn.Redis.Get([]byte("set_" + key))
		if err != nil {
			return false, fmt.Errorf("txid for set_key not found")
		}
		txhash := common.HexToHash(string(txid))
		log.Println("verifying tx", string(txid))
		receipt, err := ethereumConn.Client.TransactionReceipt(context.Background(), txhash)
		if err != nil {
			return false, fmt.Errorf("TransactionReceipt %v %v", string(txid), err)
		}

		if receipt == nil {
			return false, fmt.Errorf("TransactionReceipt null %v", txhash)
		}
		log.Println("verify receipt status ", receipt.Status)
		_, isPending, err := ethereumConn.Client.TransactionByHash(context.Background(), txhash)
		log.Println("verify tx: isPending ", isPending)
		if receipt.Status == 1 {
			return true, nil
		} else {
			return false, nil
		}
	case "get": //compare value
		getvalue, err := ethereumConn.Redis.Get([]byte("get_" + key))
		if err != nil {
			return false, fmt.Errorf("value for get_key not found")
		}
		//auth, err := ethereumConn.bindTransactOpts(ctx, key)   //no bind
		result, err := ethereumConn.KV.Items(nil, keyToByte32(key)) //Get(auth, []byte(key))
		if err != nil {
			log.Println("error EthereumConnector Read ", err)
			return false, err
		}
		if string(getvalue) == string(result) {
			return true, nil
		} else {
			return false, nil
		}

	default:
		return false, fmt.Errorf("Verify operation only support get/set.")

	}

}

func (ethereumConn *EthereumConnector) bindTransactOpts(ctx context.Context, key string) (*bind.TransactOpts, error) {
	gasPrice, err := ethereumConn.Client.SuggestGasPrice(ctx)
	if err != nil {
		log.Println("error parse a secp256k1 private key.", err)
		return nil, err
	}
	privateKey, err := crypto.HexToECDSA(ethereumConn.Hexkey)
	if err != nil {
		log.Println("error casting public key to ECDSA.", err)
		return nil, err
	}
	auth := bind.NewKeyedTransactor(privateKey)
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(1000000) // in units
	auth.GasPrice = gasPrice

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Println("error casting public key to ECDSA")
		return nil, err
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// nounce unique to avoid transaction failure
	for {
		nonce, err := ethereumConn.Client.PendingNonceAt(ctx, fromAddress)
		if err != nil {
			log.Println("error return the account nonce of the given account in the pending state.", err)
			return nil, err
		}

		if ethereumConn.TxMgr.WriteNounce(int64(nonce), key) {
			auth.Nonce = big.NewInt(int64(nonce))
			//log.Println(auth.Nonce)
			break
		}
	}

	return auth, nil
}
