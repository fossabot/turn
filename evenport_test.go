package turn

import (
	"github.com/ernado/stun"
	"testing"
)

func TestEvenPort(t *testing.T) {
	t.Run("False", func(t *testing.T) {
		m := new(stun.Message)
		p := EvenPort{
			ReservePort: false,
		}
		if err := p.AddTo(m); err != nil {
			t.Error(err)
		}
		m.WriteHeader()
		decoded := new(stun.Message)
		var port EvenPort
		decoded.Write(m.Raw)
		if err := port.GetFrom(m); err != nil {
			t.Fatal(err)
		}
		if port != p {
			t.Fatal("not equal")
		}
	})
	t.Run("AddTo", func(t *testing.T) {
		m := new(stun.Message)
		p := EvenPort{
			ReservePort: true,
		}
		if err := p.AddTo(m); err != nil {
			t.Error(err)
		}
		m.WriteHeader()
		t.Run("GetFrom", func(t *testing.T) {
			decoded := new(stun.Message)
			if _, err := decoded.Write(m.Raw); err != nil {
				t.Fatal("failed to decode message:", err)
			}
			port := EvenPort{}
			if err := port.GetFrom(decoded); err != nil {
				t.Fatal(err)
			}
			if port != p {
				t.Errorf("Decoded %q, expected %q", port, p)
			}
			if wasAllocs(func() {
				port.GetFrom(decoded)
			}) {
				t.Error("Unexpected allocations")
			}
			t.Run("HandleErr", func(t *testing.T) {
				m := new(stun.Message)
				var handle EvenPort
				if err := handle.GetFrom(m); err != stun.ErrAttributeNotFound {
					t.Errorf("%v should be not found", err)
				}
				m.Add(stun.AttrEvenPort, []byte{1, 2, 3})
				if err, ok := handle.GetFrom(m).(*BadAttrLength); !ok {
					t.Errorf("%v should be *BadAttrLength", err)
				}
			})
		})
	})
}
