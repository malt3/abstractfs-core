package validate_test

import (
	"testing"

	"github.com/malt3/abstractfs-core/api"
	"github.com/malt3/abstractfs-core/validate"
	"github.com/stretchr/testify/assert"
)

func TestValidateStat(t *testing.T) {
	testCases := map[string]struct {
		dir            string
		stat           api.Stat
		wantFlatErrors []validate.StatValidationError
		wantTreeErrors []validate.StatValidationError
	}{
		"empty root node": {
			stat: api.Stat{
				Kind: api.KindDirectory,
			},
			wantFlatErrors: []validate.StatValidationError{
				{Stat: api.Stat{Kind: "directory"}, Reason: validate.NodeNameEmpty},
				{Stat: api.Stat{Kind: "directory"}, Reason: validate.FlatNameNotRooted},
			},
		},
		"rooted root node": {
			stat: api.Stat{
				Name: "/",
				Kind: api.KindDirectory,
			},
			wantTreeErrors: []validate.StatValidationError{
				{Stat: api.Stat{Name: "/", Kind: "directory"}, Reason: validate.NodeNameContainsSlash},
				{Stat: api.Stat{Name: "/", Kind: "directory"}, Reason: validate.TreeNameRooted},
			},
		},
		"regular node": {
			dir: "foo",
			stat: api.Stat{
				Name:    "bar",
				Kind:    api.KindRegular,
				Payload: "sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
			},
			wantFlatErrors: []validate.StatValidationError{
				{
					Dir: "foo",
					Stat: api.Stat{
						Name: "bar", Kind: "regular", Payload: "sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
					},
					Reason: validate.FlatNameNotRooted,
				},
			},
		},
		"regular node missing payload": {
			dir: "foo",
			stat: api.Stat{
				Name: "bar",
				Kind: api.KindRegular,
			},
			wantFlatErrors: []validate.StatValidationError{
				{
					Dir: "foo",
					Stat: api.Stat{
						Name: "bar", Kind: "regular",
					},
					Reason: validate.FlatNameNotRooted,
				},
				{
					Dir: "foo",
					Stat: api.Stat{
						Name: "bar", Kind: "regular",
					},
					Reason: validate.PayloadEmpty,
				},
				{
					Dir: "foo",
					Stat: api.Stat{
						Name: "bar", Kind: "regular",
					},
					Reason: validate.PayloadInvalidSRI,
				},
			},
			wantTreeErrors: []validate.StatValidationError{
				{
					Dir: "foo",
					Stat: api.Stat{
						Name: "bar", Kind: "regular",
					},
					Reason: validate.PayloadEmpty,
				},
				{
					Dir: "foo",
					Stat: api.Stat{
						Name: "bar", Kind: "regular",
					},
					Reason: validate.PayloadInvalidSRI,
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			gotFlat := validate.ValidateStat(tc.dir, tc.stat, validate.ValidationModeFlat)
			gotTree := validate.ValidateStat(tc.dir, tc.stat, validate.ValidationModeTree)
			assert.ElementsMatch(tc.wantFlatErrors, gotFlat)
			assert.ElementsMatch(tc.wantTreeErrors, gotTree)
		})
	}
}
