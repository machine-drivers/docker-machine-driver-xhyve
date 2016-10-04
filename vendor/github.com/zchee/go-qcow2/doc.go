// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package qcow2 implements manage the QEMU qcow2 image format written in pure Go.
//
// The below contents are qemu-img command line options from the man page.
//
//  "compat"
//   "compat=0.10": uses the traditional image format that can be read by any QEMU since 0.10.
//   "compat=1.1":  enables image format extensions that only QEMU 1.1 and newer understand (this is the default).
//
//  "backing_file"
//   File name of a base image (see create subcommand)
//
//  "backing_fmt"
//   Image format of the base image
//
//  "encryption"
//   If this option is set to "on", the image is encrypted with 128-bit AES-CBC.
//
//  "cluster_size"
//   Changes the qcow2 cluster size (must be between 512 and 2M).
//   Smaller cluster sizes can improve the image file size whereas larger cluster sizes generally provide better performance.
//
//  "preallocation"
//   Preallocation mode (allowed values: "off", "metadata", "falloc", "full").
//   An image with preallocated metadata is initially larger but can improve performance when the image needs to grow.
//   "falloc" and "full" preallocations are like the same options of "raw" format, but sets up metadata also.
//
//  "lazy_refcounts"
//   If this option is set to "on", reference count updates are postponed with the goal of avoiding metadata I/O and improving performance.
//   This is particularly interesting with cache=writethrough which doesn't batch metadata updates.
//   The tradeoff is that after a host crash, the reference count tables must be rebuilt,
//   i.e. on the next open an (automatic) "qemu-img check -r all" is required, which may take some time.
//   This option can only be enabled if "compat=1.1" is specified.
//
//  "nocow"
//   If this option is set to "on", it will turn off COW of the file. It's only valid on btrfs, no effect on other file systems.
//   Btrfs has low performance when hosting a VM image file, even more when the guest on the VM also using btrfs as file system.
//   Turning off COW is a way to mitigate this bad performance. Generally there are two ways to turn off COW on btrfs: a)
//   Disable it by mounting with nodatacow, then all newly created files will be NOCOW. b)
//   For an empty file, add the NOCOW file attribute. That's what this option does.
//   Note: this option is only valid to new or empty files.
//   If there is an existing file which is COW and has data blocks already, it couldn't be changed to NOCOW by setting "nocow=on".
//   One can issue "lsattr filename" to check if the NOCOW flag is set or not (Capital 'C' is NOCOW flag).
package qcow2
