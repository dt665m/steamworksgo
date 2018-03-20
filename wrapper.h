// wrapper.h
int SteamInit();
int SteamIsSteamRunning();
void SteamShutdown();

//authentication
int RgubKeyLength();
int CheckSteamEncryptedAppTicket(const unsigned char *rgubKey, const unsigned char *rgubTicket, const unsigned int cubTicket, const unsigned int appId, const unsigned int timeNow, const unsigned int expectedAccountID, const unsigned int expiryToleranceSeconds);
