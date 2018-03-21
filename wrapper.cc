// wrapper
// +build !windows
//#include "steam_api.h"
#include "public/steamencryptedappticket.h"
#include "stdio.h"
#include "stdlib.h"

extern "C" {
	#include "wrapper.h"
	#include "_cgo_export.h"
}

/*
The server side check for SteamUser()->GetEncryptedAppTicket( rgubTicket, sizeof( rgubTicket), &cubTicket );
Required Params:

*uchar rgubEncryptedTicket Ticket, (length not needed by the steam lib?)
uint32 cubTicket, 
uint32 AppId_t, and 
uint32 Time-expiration serialized to steam specs (#TODO)
uint32 expectedSteamID

*/
int CheckSteamEncryptedAppTicket(
		const uint8 *rgubKey, 
		const uint8 *rgubTicket, 
		const uint32 cubTicket, 
		const AppId_t appID, 
		const RTime32 timeNow, 
		const AccountID_t expectedAccountID,
		const unsigned int expiryToleranceSeconds) {	
	
	int keylen = 32;
	uint8 rgubDecrypted[1024];
	uint32 cubDecrypted = sizeof( rgubDecrypted );
	if ( !SteamEncryptedAppTicket_BDecryptTicket( rgubTicket, cubTicket, rgubDecrypted, &cubDecrypted, rgubKey, keylen ) )  //sizeof( rgubKey )
		return 1;
	

	//printf("Decrypt Success!\n");
	AppId_t ticketAppId = SteamEncryptedAppTicket_GetTicketAppID( rgubDecrypted, cubDecrypted );
	//printf("AppId in Ticket: %d\n", ticketAppId);
	if ( !SteamEncryptedAppTicket_BIsTicketForApp( rgubDecrypted, cubDecrypted, appID ) ) {
		return 2;
	}

	//check time
	RTime32 issueTime = SteamEncryptedAppTicket_GetTicketIssueTime( rgubDecrypted, cubDecrypted );
	if ( timeNow - issueTime > expiryToleranceSeconds ) {
		printf("ticket issue time: %d\n", issueTime);
		printf("ticket elapsed seconds: %d\n", timeNow - issueTime);
		printf("ticket tolerance %d\n", expiryToleranceSeconds);
		return 3;
	}
	
	//get and check accountID
	CSteamID steamIDFromTicket;
	SteamEncryptedAppTicket_GetTicketSteamID( rgubDecrypted, cubDecrypted, &steamIDFromTicket );
	if ( expectedAccountID != steamIDFromTicket.GetAccountID()  ) {
		printf("steamid in ticket: %d\n", (int)steamIDFromTicket.ConvertToUint64());
		return 4;
	}
	
	//check ownership (supposedly for dlc)
	if ( !SteamEncryptedAppTicket_BUserOwnsAppInTicket( rgubDecrypted, cubDecrypted, appID ) ){
		return 5; 
	}

	return 0;	
}

int RgubKeyLength() {
	return k_nSteamEncryptedAppTicketSymmetricKeyLen;
}

//CSteamAPIContext myContext;

int SteamInit() {
	//myContext.Init();
		return 1;
		/*
	if (SteamAPI_InitSafe())
	{		
		myContext.Init();
		return 1;
	}
	else
	{
		const char b[] = "SteamAPI_InitSafe failed."
		GoLog(b);
	}
	return 0;*/
}

int SteamIsSteamRunning() {
	return 0;
	/*
	if (SteamAPI_IsSteamRunning())
		return 1;
	else
		return 0;
		*/
}

void SteamShutdown() {
	//SteamAPI_Shutdown();
}