// Package rfb implements RFC6143 (also known as VNC).
//
// Right now only the client side is supported, but the general
// architecture of the package will allow a futere extension for
// a server side as well.
//
// The RFC is not fully implemenent, a lot of encodings are missing
// and there's no authentication as of now.
//
// https://tools.ietf.org/rfc/rfc6143.txt
package rfb
