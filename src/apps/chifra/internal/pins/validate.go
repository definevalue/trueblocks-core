// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.

package pinsPkg

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/validate"
)

func (opts *PinsOptions) ValidatePins() error {
	opts.TestLog()

	if opts.BadFlag != nil {
		return opts.BadFlag
	}

	if opts.InitAll {
		return validate.Deprecated("--init_all", "--init --all")

	} else if opts.Freshen {
		return validate.Deprecated("--freshen", "--init")

	} else if opts.Remote {
		return validate.Deprecated("--remote", "")

	} else if opts.Share {
		return validate.Usage("The --share flag has been deprecated. See the chifra chunks command instead.")
	}

	if opts.List && opts.Init {
		return validate.Usage("Please choose only one of {0}.", "--list or --init")
	}

	if !opts.List && !opts.Init {
		return validate.Usage("Please choose at least one of {0}.", "--list or --init")
	}

	if opts.All && !opts.Init {
		return validate.Usage("The {0} option is available only with {1}.", "--all", "--init")
	}

	return opts.Globals.ValidateGlobals()
}
