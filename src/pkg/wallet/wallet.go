package wallet

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/jbenet/go-base58"
)

type Wallet struct {
	PrivateKey    string
	PublicKey     string
	AddressHex    string
	AddressBase58 string
}

func NewWallet() Wallet {
	return Wallet{}
}

func (w *Wallet) Generate() {
	// TODO
	// generate tron private key

}

func (w *Wallet) GetAddress() (addressBase58 string, err error) {

	// private key from bytes
	privateKey, _ := crypto.HexToECDSA(w.PrivateKey)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		errMsg := "cannot assert type: publicKey is not of type *ecdsa.PublicKey"
		fmt.Println(errMsg)
		return "", errors.New(errMsg)
	}

	// try 4
	// hash from hex addr begings with 41
	// encode with SHA256
	// append 0x to addresshex

	addrHexBytes, _ := hexutil.Decode(w.AddressHex)
	// encode with double SHA256
	sha256hash1 := sha256.Sum256(addrHexBytes)
	// do double sha256
	fmt.Println("sha256hash1:", hexutil.Encode(sha256hash1[:]))
	sha245hash1bytes := []byte(sha256hash1[:])
	sha256hash2 := sha256.Sum256(sha245hash1bytes)
	fmt.Println("sha256hash2:", hexutil.Encode(sha256hash2[:]))
	// add to the end first 4 bytes from doublehash (HACK, got it from tron js lib)
	doubleHashBytes := []byte(sha256hash2[0:4])
	eth_address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	eth_address_bytes, _ := hexutil.Decode(eth_address)
	eth4 := append(eth_address_bytes, doubleHashBytes...)
	eth4_bytes := []byte(eth4)
	eth4_bytes = append([]byte{0x41}, eth4_bytes...)
	encoded := base58.Encode(eth4_bytes)
	fmt.Println("eth4 b58 encoded (RIGHT):", encoded)

	// // decode back to hex address
	// decoded_test := base58.Decode(address58)
	// // cut 1 byte from the end
	// decoded_test = decoded_test[:len(decoded_test)-0]
	// fmt.Println("b58 decode test addr: ", hexutil.Encode(decoded_test))
	// fmt.Println("b58 encoded back:", b58.Encode(decoded_test))
	return encoded, nil
}

func (w *Wallet) GetPrivateKey() string {
	return w.PrivateKey
}

func (w *Wallet) GetPublicKey() string {
	return w.PublicKey
}
