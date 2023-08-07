// Copyright 2018 GRAIL, Inc.  All rights reserved.
// Use of this source code is governed by the Apache-2.0
// license that can be found in the LICENSE file.

//go:build !cgo
// +build !cgo

package libdeflate

func init() {
	need_build_tags_use_slow_gzip()
}
