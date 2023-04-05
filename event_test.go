package retraced

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashMismatch(t *testing.T) {
	testEvent := &Event{
		Action: "just.a.test",
		Group: &Group{
			ID: "Customer: XYZ",
		},
		SourceIP:    "1.2.3.4",
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
		t.Errorf("Hash check should have succeeded %v", err)
	}
}

func TestHashMatchChannelList(t *testing.T) {
	testEvent := &Event{
		Action: "channel.list",
		Group: &Group{
			ID: "602f21a3fbd3f92302133762808b39af",
		},
		Actor: &Actor{
			ID: "060dbbd5da8c43b57b26179a3bfb7b1a",
		},
		Target: &Target{
			ID: "6da2ecf53d388e107df6e4dbb061b165",
		},
		SourceIP:    "172.19.0.1",
		IsAnonymous: false,
		IsFailure:   false,
	}

	fakeNew := &NewEventRecord{
		ID:   "f59b236a449d43a5b27c8322aadc0503",
		Hash: "2224989b8d83d4b23920f0136f8e3b11ce034d9e0b610ee97c1c198350838a9e",
	}

	hashTarget := string(testEvent.BuildHashTarget(fakeNew))
	expected := "f59b236a449d43a5b27c8322aadc0503:channel.list:6da2ecf53d388e107df6e4dbb061b165:060dbbd5da8c43b57b26179a3bfb7b1a:602f21a3fbd3f92302133762808b39af:172.19.0.1:0:0::"

	assert.New(t).Equal(expected, hashTarget, "Hash targets should be equal")
}

func TestBuildHashNoGroupId(t *testing.T) {
	testEvent := &Event{
		Action: "even.more.of.a.test",
		Actor: &Actor{
			ID: "user@domain.xyz",
		},
		Target: &Target{
			ID: "some_object01234",
		},
		IsAnonymous: false,
		IsFailure:   true,
		Fields: map[string]string{
			"abc=xyz": "nothing special",
		},
	}

	fakeNew := &NewEventRecord{
		ID: "kfbr392",
	}

	hashTarget := string(testEvent.BuildHashTarget(fakeNew))
	expected := "kfbr392:even.more.of.a.test:some_object01234:user@domain.xyz:::1:0:abc%3Dxyz=nothing special;"

	assert.New(t).Equal(expected, hashTarget, "Hash targets should be equal")
}
