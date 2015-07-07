Package qprintable implements quoted-printable encoding as specified
by RFC 2045. It is strict on ouput, generous on input.

Quoting RFC 2045:
> The Quoted-Printable encoding is intended to represent data that
> largely consists of octets that correspond to printable characters in
> the US-ASCII character set.  It encodes the data in such a way that
> the resulting octets are unlikely to be modified by mail transport.
> If the data being encoded are mostly US-ASCII text, the encoded form
> of the data remains largely recognizable by humans.  A body which is
> entirely US-ASCII may also be encoded in Quoted-Printable to ensure
> the integrity of the data should the message pass through a
> character-translating, and/or line-wrapping gateway.


# VARIABLES 

`var BinaryEncoding = &Encoding{false, ""}`

In binary encoding, CR and LF characters are treated like other control
characters sequence and are escaped.

`var MacTextEncoding = &Encoding{true, "\r"}`

`var UnixTextEncoding = &Encoding{true, "\n"}`

`var WindowsTextEncoding = &Encoding{true, "\r\n"}`

A text encoding has to convert its input in the canonical form (as
defined by RFC 2045) : native ends of line (CR for MacTextEncoding,
LF for UnixTextEncoding, CRLF for WindowsTextEncoding) are converted
into CRLF sequences. Non-native EOL sequences (for example, CR on
UnixTextEncoding) are treated as control characters and escaped.

In the decoding process, CRLF sequences are converted to native ends of
line.


# FUNCTIONS 

`func NewDecoder(enc *Encoding, r io.Reader) io.Reader`

Returns a new decoder. Data will be read from r, and decoded
according to enc.

`func NewEncoder(enc *Encoding, w io.Writer) io.WriteCloser`

Returns a new encoder. Any data passed to Write will be encoded
according to enc and then written to w.

Data passed to Write must be in its canonical form. The canonical
form depends on the encoding:

* for binary encoding, anything goes.
* for text encodings, there shouldn't be any CR or LF characters
  other than the one used for end-of-line representation, that is,
  LF on Unix, CR on old Mac, CR+LF on Windows.

It is the responsibility of the caller to ensure that the input stream
is in its canonical form. Any CR of LF character which is not part of
an end-of-line representation will be quoted.

This returns a WriteCloser, but Close has no effect for encoding other
than WindowsEncoding.

For WindowsEncoding, any trailing CR will not be written unless
you call this function. However, note that for a text conforming
to windows canonical form, this should never happen. So this
function is useful only for invalid WindowsEncoding text streams,
you can safely ignore it in all other cases.


# TYPES 

	type Encoding struct {
		// contains unexported fields
	}

`func DetectEncoding(data string) *Encoding`

Try to detect encoding of string:
* strings with no \r will be Unix
* strings with \r and no \n will be Mac
* strings with count(\r\n) == count(\r) == count(\n) will be Windows
* other strings will be binary
