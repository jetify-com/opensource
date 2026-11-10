package ids

import (
	"strings"
	"testing"

	"go.jetify.com/typeid/v2"
)

type idCase struct {
	prefix  string
	newFn   func() (string, error)
	parseFn func(string) (string, error)
}

func wrap[T interface{ String() string }](
	newFn func() (T, error),
	parseFn func(string) (T, error),
	prefix string,
) idCase {
	return idCase{
		prefix: prefix,
		newFn: func() (string, error) {
			id, err := newFn()
			return id.String(), err
		},
		parseFn: func(s string) (string, error) {
			id, err := parseFn(s)
			if err != nil {
				return "", err
			}
			return id.String(), nil
		},
	}
}

func allCases() []idCase {
	return []idCase{
		wrap(NewUserID, ParseUserID, UserPrefix),
		wrap(NewProjectID, ParseProjectID, ProjectPrefix),
		wrap(NewRepoID, ParseRepoID, RepoPrefix),
		wrap(NewOrgID, ParseOrgID, OrgPrefix),
		wrap(NewMemberID, ParseMemberID, MemberPrefix),
		wrap(NewSecretID, ParseSecretID, SecretPrefix),
		wrap(NewDeploymentID, ParseDeploymentID, DeploymentPrefix),
		wrap(NewCustomDomainID, ParseCustomDomainID, CustomDomainPrefix),
	}
}

func TestNewAndParseRoundTrip(t *testing.T) {
	for _, tc := range allCases() {
		t.Run(tc.prefix, func(t *testing.T) {
			s, err := tc.newFn()
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}
			if !strings.HasPrefix(s, tc.prefix+"_") {
				t.Errorf("New() = %q, want prefix %q", s, tc.prefix+"_")
			}

			got, err := tc.parseFn(s)
			if err != nil {
				t.Fatalf("Parse(%q) error = %v", s, err)
			}
			if got != s {
				t.Errorf("Parse(%q).String() = %q", s, got)
			}
		})
	}
}

func TestNewIsUnique(t *testing.T) {
	a, err := NewUserID()
	if err != nil {
		t.Fatal(err)
	}
	b, err := NewUserID()
	if err != nil {
		t.Fatal(err)
	}
	if a.String() == b.String() {
		t.Errorf("two NewUserID() calls returned the same ID %q", a)
	}
}

func TestParseRejectsWrongPrefix(t *testing.T) {
	cases := allCases()
	for i, tc := range cases {
		// Generate an ID belonging to a different type.
		other := cases[(i+1)%len(cases)]
		s, err := other.newFn()
		if err != nil {
			t.Fatal(err)
		}
		t.Run(tc.prefix, func(t *testing.T) {
			_, err := tc.parseFn(s)
			if err == nil {
				t.Fatalf("Parse(%q) succeeded, want error for wrong prefix", s)
			}
			if !strings.Contains(err.Error(), "invalid "+tc.prefix+" ID") {
				t.Errorf("Parse(%q) error = %q, want it to mention the expected prefix", s, err)
			}
		})
	}
}

func TestParseRejectsInvalid(t *testing.T) {
	for _, s := range []string{
		"",
		"user",
		"user_",
		"user_not-a-valid-suffix",
		"USER_01h455vb4pex5vsknk084sn02q",
	} {
		t.Run(s, func(t *testing.T) {
			if _, err := ParseUserID(s); err == nil {
				t.Errorf("ParseUserID(%q) succeeded, want error", s)
			}
		})
	}
}

func TestParseZeroSuffix(t *testing.T) {
	const s = "user_00000000000000000000000000"
	id, err := ParseUserID(s)
	if err != nil {
		t.Fatalf("ParseUserID(%q) error = %v", s, err)
	}
	if id.String() != s {
		t.Errorf("ParseUserID(%q).String() = %q", s, id)
	}
}

func TestShortSuffix(t *testing.T) {
	tid, err := typeid.Parse("user_01h455vb4pex5vsknk084sn02q")
	if err != nil {
		t.Fatal(err)
	}
	if got := ShortSuffix(tid); got != "4sn02q" {
		t.Errorf("ShortSuffix() = %q, want %q", got, "4sn02q")
	}

	// A zero TypeID still has a 26 char suffix, so it should never panic.
	var zero typeid.TypeID
	if got := ShortSuffix(zero); len(got) != 6 {
		t.Errorf("ShortSuffix(zero) = %q, want 6 chars", got)
	}
}
