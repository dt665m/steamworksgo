// +build linux darwin
// +build cgo

package steamworksgo

/*
#cgo linux,!darwin LDFLAGS: -L${SRCDIR}/linux64 -lsteam_api
#cgo linux,!darwin LDFLAGS: -L${SRCDIR}/linux64 -lsdkencryptedappticket
#cgo darwin LDFLAGS: -L${SRCDIR}/osx32 -lsteam_api
#cgo darwin LDFLAGS: -L${SRCDIR}/osx32 -lsdkencryptedappticket

#include "wrapper.h"
*/
import "C"
import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"
	"unsafe"
)

const (
	EncryptedTicketLength = 1024
)

var (
	rgubExpectedLength int
)

func init() {
	rlength, err := C.RgubKeyLength()
	if err != nil {
		panic("k_nSteamEncryptedAppTicketSymmetricKeyLen retrieval failed: " + err.Error())
	}
	rgubExpectedLength = int(rlength)
}

type SteamWorks struct {
	SteamAPIKey     []byte
	AppID           uint32
	ExpiryTolerance uint32
}

func NewSteamWorks(rgubKeyString string, appID uint32, expiryTolerance uint32) (*SteamWorks, error) {
	//Secret steam key
	rgub, err := hex.DecodeString(rgubKeyString)
	if err != nil {
		return nil, err
	}

	if len(rgub) != rgubExpectedLength {
		return nil, fmt.Errorf("private key length invalid. expected %v got %v", rgubExpectedLength, len(rgub))
	}

	return &SteamWorks{
		SteamAPIKey:     rgub,
		AppID:           appID,
		ExpiryTolerance: expiryTolerance,
	}, nil
}

func (s *SteamWorks) VerifyAppTicket(encryptedTicketB64 string, cubTicket uint32, steamID uint32) error {
	encryptedTicket, err := base64.StdEncoding.DecodeString(encryptedTicketB64)
	if err != nil {
		return err
	}

	if len(encryptedTicket) != EncryptedTicketLength {
		return fmt.Errorf("encrypted ticket length invalid. expected %v got %v", EncryptedTicketLength, len(encryptedTicket))
	}

	keyPtr := (*C.uchar)(unsafe.Pointer(&s.SteamAPIKey[0]))
	ticketPtr := (*C.uchar)(unsafe.Pointer(&encryptedTicket[0]))

	code, err := C.CheckSteamEncryptedAppTicket(
		keyPtr,
		ticketPtr,
		C.uint(cubTicket),
		C.uint(s.AppID),
		C.uint(time.Now().Unix()),
		C.uint(steamID),
		C.uint(s.ExpiryTolerance),
	)
	if err != nil {
		return fmt.Errorf("C.CheckSteamEncryptedAppTicket error: %v", err)
	}

	switch code {
	case 0:
		return nil
	case 1:
		return fmt.Errorf("failed to decrypt ticket, code %v", code)
	case 2:
		return fmt.Errorf("ticket not for specified appid %v, code %v", s.AppID, code)
	case 3:
		return fmt.Errorf("ticket expired, code: %v", code)
	case 4:
		return fmt.Errorf("ticket appid does not match %v, code %v", steamID, code)
	case 5:
		return fmt.Errorf("dlc not verified, code: %v", code)
	}

	return nil
}
