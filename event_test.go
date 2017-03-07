package retraced

import "testing"

func TestHashMismatch(t *testing.T) {
	testEvent := &Event{
		Action: "just.a.test",
		Group: &Group{
			ID: "Customer: XYZ",
		},
		SourceIp:    "1.2.3.4",
		IsAnonymous: true,
		Fields: map[string]string{
			"custom": "123",
			"Custom": "Rate = 50%",
		},
	}

	fakeNew := &NewEventRecord{
		ID:   "0123456789abcdefg",
		Hash: "XXXXXXXXX",
	}
	if err := testEvent.VerifyHash(fakeNew); err != nil {
		// pass
	} else {
		t.Errorf("Hash check should have failed")
	}
}

func TestHashMatch(t *testing.T) {
	testEvent := &Event{
		Action: "even.more.of.a.test",
		Group: &Group{
			ID: "%% :: some %% customer :: %%",
		},
		Actor: &Actor{
			ID: "user@domain.xyz",
		},
		Target: &Target{
			ID: "some_object01234",
		},
		IsAnonymous: false,
		IsFailure:   true,
		Fields: map[string]string{
			";zyx=cba;abc=xyz": "nothing special",
			";Zyx=Cba%Abc=Xyz": "% hi there %",
		},
	}

	fakeNew := &NewEventRecord{
		ID:   "abf053dc4a3042459818833276eec717",
		Hash: "5b570bff4628b35262fb401d2f6c9bb38d29e212f6e0e8ea93445b4e5a253d50",
	}
	if err := testEvent.VerifyHash(fakeNew); err != nil {
		t.Errorf("Hash check should have succeeded")
	} else {
		// pass
	}
}
