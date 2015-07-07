package qprintable

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

type encDetectionTestData struct {
	data   string
	result *Encoding
}

type universalTestData struct {
	decoded, encoded string
}

type eolTestData struct {
	decoded, unixEncoded, winEncoded, macEncoded, binEncoded string
}

var encDetectionTests = []encDetectionTestData{
	{"Hello\nworld", UnixTextEncoding},
	{"Hello\rworld", MacTextEncoding},
	{"Hello\r\nworld", WindowsTextEncoding},
	{"Hello\rworld\n", BinaryEncoding},
	{"Hello\rworld\r\n", BinaryEncoding},
	{"Hello\nworld\r\n", BinaryEncoding},
	{"\r\nHello\rworld\n", BinaryEncoding},
	{"Hello world", UnixTextEncoding},
}

var universalTests = []universalTestData{ // Will be tested with all encodings
	{"Touché !", "Touch=C3=A9 !"},
	{"Pi = 4", "Pi =3D 4"},
	{"Ascii 27 (Escape) = \033", "Ascii 27 (Escape) =3D =1B"},
}

var eolTests = []eolTestData{
	// Since non-eol CR or LF characters are invalid in text encodings, these strings are invalid.
	// In consequence, we lose idempotence (decode(encode(x)) can be ≠ x), so we will only test
	// that encode(decodedValue) = encodedValue. decode(encodedValue) = decodedValue will be tested
	// only for binary encoding
	{"Hello\nworld", "Hello\r\nworld", "Hello=0Aworld", "Hello=0Aworld", "Hello=0Aworld"},
	{"Hello\r\nworld", "Hello=0D\r\nworld", "Hello\r\nworld", "Hello\r\n=0Aworld", "Hello=0D=0Aworld"},
	{"Hello\rworld", "Hello=0Dworld", "Hello=0Dworld", "Hello\r\nworld", "Hello=0Dworld"},
}

var wrapTests = []universalTestData{ // Will be tested with Unix encoding
	{ // Should wrap at 75 characters
		"123456789 123456789 123456789 123456789 123456789 123456789 123456789 12345" +
			"123456789 123456789 123456789 123456789 123456789 123456789 123456789 12345" +
			"12345",
		"123456789 123456789 123456789 123456789 123456789 123456789 123456789 12345" + "=\r\n" +
			"123456789 123456789 123456789 123456789 123456789 123456789 123456789 12345" + "=\r\n" +
			"12345",
	},
	{ // Should wrap at 75 characters even when a character is expanded
		"123456789 123456789 1é89 123456789 123456789 123456789 123456789 12345" +
			"123456789 123456789 123456789 1é89 123456789 123456789 123456789 12345" +
			"12345",
		"123456789 123456789 1=C3=A989 123456789 123456789 123456789 123456789 12345" + "=\r\n" +
			"123456789 123456789 123456789 1=C3=A989 123456789 123456789 123456789 12345" + "=\r\n" +
			"12345",
	},
	{ // =xx sequences are atomic
		"123456789 123456789 123456789 123456789 123456789 123456789 123456789 1234é" +
			"789 123456789 123456789 123456789 123456789 123456789 123456789 12345" +
			"12345",
		"123456789 123456789 123456789 123456789 123456789 123456789 123456789 1234" + "=\r\n" +
			"=C3=A9789 123456789 123456789 123456789 123456789 123456789 123456789 12345" + "=\r\n" +
			"12345",
	},
	{ // Hard line breaks should reset counter
		"1234\n" +
			"123456789 123456789 123456789 123456789 123456789 123456789 123456789 12345" +
			"12345",
		"1234\r\n" +
			"123456789 123456789 123456789 123456789 123456789 123456789 123456789 12345" + "=\r\n" +
			"12345",
	},
}

func testEqual(t *testing.T, testName string, expected, actual []byte) bool {
	if bytes.Compare(expected, actual) != 0 {
		t.Logf("Test %s: result is not what was expected !", testName)
		t.Logf(" Expected result: %#v", string(expected))
		t.Logf("   Actual result: %#v", string(actual))
		t.Fail()

		return false
	}
	return true
}

/**
 * Encoding detection tests
 */
func TestEncodingDetection(t *testing.T) {
	for _, data := range encDetectionTests {
		if DetectEncoding(data.data) != data.result {
			t.Errorf("DetectEncoding(%#v): expected %#v, got %#v", data.data, data.result.nativeEol, DetectEncoding(data.data).nativeEol)
		}
	}
}

/**
 * Encoder tests
 */

func testEncodeChunked(t *testing.T, chunkSize int, testName string, enc *Encoding, decoded, encoded []byte) {
	var part []byte
	encBuf := bytes.NewBuffer(nil)
	encoder := NewEncoder(enc, encBuf)
	for decoded != nil {
		if len(decoded) > chunkSize {
			part = decoded[:chunkSize]
			decoded = decoded[chunkSize:]
		} else {
			part = decoded
			decoded = nil
		}
		_, err := encoder.Write(part)
		if err != nil {
			t.Errorf("Test %s: encoder error: %s", testName, err.Error())
			return
		}
	}
	testEqual(t, testName, encoded, encBuf.Bytes())
}

func testEncode(t *testing.T, testName string, enc *Encoding, decoded, encoded []byte) {
	testEncodeChunked(t, len(decoded), fmt.Sprintf("%s/c*", testName), enc, decoded, encoded)
	testEncodeChunked(t, 1, fmt.Sprintf("%s/c1", testName), enc, decoded, encoded)
	testEncodeChunked(t, 3, fmt.Sprintf("%s/c3", testName), enc, decoded, encoded)
	testEncodeChunked(t, 4, fmt.Sprintf("%s/c4", testName), enc, decoded, encoded)
	testEncodeChunked(t, 16, fmt.Sprintf("%s/c16", testName), enc, decoded, encoded)
}

func TestUniversalEncode(t *testing.T) {
	for i, testData := range universalTests {
		testEncode(t, fmt.Sprintf("%d/Binary", i+1), BinaryEncoding, []byte(testData.decoded), []byte(testData.encoded))
		testEncode(t, fmt.Sprintf("%d/Unix", i+1), UnixTextEncoding, []byte(testData.decoded), []byte(testData.encoded))
		testEncode(t, fmt.Sprintf("%d/Windows", i+1), WindowsTextEncoding, []byte(testData.decoded), []byte(testData.encoded))
		testEncode(t, fmt.Sprintf("%d/Mac", i+1), MacTextEncoding, []byte(testData.decoded), []byte(testData.encoded))
	}
}

func TestEOLEncode(t *testing.T) {
	for i, testData := range eolTests {
		testEncode(t, fmt.Sprintf("%d/Binary", i+1), BinaryEncoding, []byte(testData.decoded), []byte(testData.binEncoded))
		testEncode(t, fmt.Sprintf("%d/Unix", i+1), UnixTextEncoding, []byte(testData.decoded), []byte(testData.unixEncoded))
		testEncode(t, fmt.Sprintf("%d/Windows", i+1), WindowsTextEncoding, []byte(testData.decoded), []byte(testData.winEncoded))
		testEncode(t, fmt.Sprintf("%d/Mac", i+1), MacTextEncoding, []byte(testData.decoded), []byte(testData.macEncoded))
	}
}

func TestWrapEncode(t *testing.T) {
	for i, testData := range wrapTests {
		testEncode(t, fmt.Sprintf("%d/Unix", i+1), UnixTextEncoding, []byte(testData.decoded), []byte(testData.encoded))
	}
}

/**
 * Decoder tests
 */

func testDecodeChunked(t *testing.T, chunkSize int, testName string, enc *Encoding, decoded, encoded []byte) {
	decoder := NewDecoder(enc, bytes.NewBuffer(encoded))
	decBuf := bytes.NewBuffer(nil)
	chunk := make([]byte, chunkSize)

	var err error
	var n int

	for err != io.EOF {
		n, err = decoder.Read(chunk)
		decBuf.Write(chunk[:n])
	}

	testEqual(t, testName, decoded, decBuf.Bytes())
}

func testDecode(t *testing.T, testName string, enc *Encoding, decoded, encoded []byte) {
	testDecodeChunked(t, len(decoded), fmt.Sprintf("%s/c*", testName), enc, decoded, encoded)
	testDecodeChunked(t, 1, fmt.Sprintf("%s/c1", testName), enc, decoded, encoded)
	testDecodeChunked(t, 3, fmt.Sprintf("%s/c3", testName), enc, decoded, encoded)
	testDecodeChunked(t, 4, fmt.Sprintf("%s/c4", testName), enc, decoded, encoded)
	testDecodeChunked(t, 16, fmt.Sprintf("%s/c16", testName), enc, decoded, encoded)
}

func TestUniversalDecode(t *testing.T) {
	for i, testData := range universalTests {
		testDecode(t, fmt.Sprintf("%d/Binary", i+1), BinaryEncoding, []byte(testData.decoded), []byte(testData.encoded))
		testDecode(t, fmt.Sprintf("%d/Unix", i+1), UnixTextEncoding, []byte(testData.decoded), []byte(testData.encoded))
		testDecode(t, fmt.Sprintf("%d/Windows", i+1), WindowsTextEncoding, []byte(testData.decoded), []byte(testData.encoded))
		testDecode(t, fmt.Sprintf("%d/Mac", i+1), MacTextEncoding, []byte(testData.decoded), []byte(testData.encoded))
	}
}

func TestEOLDecode(t *testing.T) {
	testDecode(t, "0/Unix", UnixTextEncoding, []byte(eolTests[0].decoded), []byte(eolTests[0].unixEncoded))
	testDecode(t, "0/Binary", BinaryEncoding, []byte(eolTests[0].decoded), []byte(eolTests[0].binEncoded))
	testDecode(t, "1/Windows", WindowsTextEncoding, []byte(eolTests[1].decoded), []byte(eolTests[1].winEncoded))
	testDecode(t, "1/Binary", BinaryEncoding, []byte(eolTests[1].decoded), []byte(eolTests[1].binEncoded))
	testDecode(t, "2/Unix", MacTextEncoding, []byte(eolTests[2].decoded), []byte(eolTests[2].macEncoded))
	testDecode(t, "2/Binary", BinaryEncoding, []byte(eolTests[2].decoded), []byte(eolTests[2].binEncoded))
}

func TestWrapDecode(t *testing.T) {
	for i, testData := range wrapTests {
		testDecode(t, fmt.Sprintf("%d/Unix", i+1), UnixTextEncoding, []byte(testData.decoded), []byte(testData.encoded))
	}
}
