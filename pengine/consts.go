// Copyright (C) 2021 SeanTolstoyevski -  mailto:seantolstoyevski@protonmail.com
// The source code of this project is licensed under the MIT license.
// You can find the license on the repo's main folder.
// Provided without warranty of any kind.

package pengine

type MODE uint8

// Encryption / decryption modes
const (

	//
	MODECBC MODE = 1

	//
	MODECFB MODE = 2

	//
	MODECTR MODE = 3

	//
	MODEOFB MODE = 4

	//
	MODEGCM MODE = 5
)
