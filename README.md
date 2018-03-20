# SteamWorks Go, a C wrapper for the Steamworks SDK for CGO. 

Included is the Steamworks 1.4.2 SDK.  This isn't a production ready package, but we use it in our backend servers to verify steam game ownership.  Windows is not supported.  Other steam sdk features are currently on hold, as I'm not that proficient at CGO dynamic linking or static linking with Go and any work done would literally be function wrappers rewritten in C. 

Wrapper.cc is the quickest hack I could come up with to access the C++ Steamworks SDK.  If anyone has a better solution for directly linking the Steamworks headers, please ping us.  Otherwise, the package is basically a Go struct that calls a C wrapper that calls the C++ steam dynamic library.

Unit testing with the standard go tools are broken as I have no idea how to make the cgo LDFLAGS look for dynamically linked libraries in either an absolute or relative folder.  The best way to try the package is to build the example and copy your platform specific libraries to the built binary's path.

Contributions welcomed!