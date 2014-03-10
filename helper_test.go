package rfb

import (
	"reflect"
	"testing"
)

func TestTypeString(t *testing.T) {
	expected := []int{'H', 'H', 'e', 'e', 'l', 'l', 'l', 'l', 'o', 'o', xkControlL, 'V', xkControlL, 'V', xkAltL, ']', xkAltL, ']'}
	received := make([]int, 0, len(expected))
	cm := &ClientMock{
		SendMessageFunc: func(msg Message) error {
			kem, ok := msg.(*KeyEventMessage)
			if !ok {
				t.Fatalf("Unexpected message type: %s", msg)
			}
			received = append(received, kem.Key)
			return nil
		},
	}

	TypeString(cm, "Hello[Ctrl+V][Alt+\\]]")
	if !reflect.DeepEqual(received, expected) {
		t.Fatalf("Received unexpected key sequence.\n%#v <- Expected\n%#v <- Received", expected, received)
	}
}
