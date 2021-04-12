package wg

import (
	//"fmt"

	wg "golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

// GenerateWGKey exported
func GenerateWGKey() (string, string, error) {

	pKey, err := wg.GenerateKey()
	if err != nil {
		return "", "", err
	}
	pubKey := pKey.PublicKey()
	return pKey.String(), pubKey.String(), nil
}
