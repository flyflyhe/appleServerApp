package config

import "fyne.io/fyne/v2"

func GetYamlConfig() []byte {
	return resourceConfigYaml.StaticContent
}

func GetApplePrivateKey() []byte {
	return resourcePrivatePem.StaticContent
}

func GetAppleRootCaPem() []byte {
	return resourceApplePem.StaticContent
}

func GetFontTTF() fyne.Resource {
	return resourceSimkaiTtf
}
