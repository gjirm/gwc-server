package wg

import (
	//"fmt"

	wg "golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

// GenerateWGKey exported
func GenerateWGKey() (string, string, error){
	
	pKey, err := wg.GenerateKey()
	if err != nil {
		return "","",err
	}
	//fmt.Println("Private key: %v",pKey.String())
	pubKey := pKey.PublicKey()
	//fmt.Println("Public key: %v",pubKey.String())
	return pKey.String(),pubKey.String(),nil
}