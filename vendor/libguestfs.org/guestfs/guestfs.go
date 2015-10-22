/* libguestfs generated file
 * WARNING: THIS FILE IS GENERATED FROM:
 *   generator/ *.ml
 * ANY CHANGES YOU MAKE TO THIS FILE WILL BE LOST.
 *
 * Copyright (C) 2009-2015 Red Hat Inc.
 *
 * This library is free software; you can redistribute it and/or
 * modify it under the terms of the GNU Lesser General Public
 * License as published by the Free Software Foundation; either
 * version 2 of the License, or (at your option) any later version.
 *
 * This library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
 * Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public
 * License along with this library; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA
 */

package guestfs

/*
#cgo CFLAGS:  -DGUESTFS_PRIVATE=1
// #cgo LDFLAGS: -undefined dynamic_lookup -lguestfs
#cgo LDFLAGS: -lguestfs
#include <stdio.h>
#include <stdlib.h>
#include "guestfs.h"

// cgo can't deal with variable argument functions.
static guestfs_h *
_go_guestfs_create_flags (unsigned flags)
{
    return guestfs_create_flags (flags);
}
*/
import "C"

import (
	"fmt"
	"runtime"
	"syscall"
	"unsafe"
)

/* Handle. */
type Guestfs struct {
	g *C.guestfs_h
}

/* Convert handle to string (just for debugging). */
func (g *Guestfs) String() string {
	return "&Guestfs{}"
}

/* Create a new handle with flags. */
type CreateFlags uint

const (
	CREATE_NO_ENVIRONMENT   = CreateFlags(C.GUESTFS_CREATE_NO_ENVIRONMENT)
	CREATE_NO_CLOSE_ON_EXIT = CreateFlags(C.GUESTFS_CREATE_NO_CLOSE_ON_EXIT)
)

func Create_flags(flags CreateFlags) (*Guestfs, error) {
	c_g, err := C._go_guestfs_create_flags(C.uint(flags))
	if c_g == nil {
		return nil, err
	}
	g := &Guestfs{g: c_g}
	// Finalizers aren't guaranteed to run, but try having one anyway ...
	runtime.SetFinalizer(g, (*Guestfs).Close)
	return g, nil
}

/* Create a new handle without flags. */
func Create() (*Guestfs, error) {
	return Create_flags(0)
}

/* Apart from Create() and Create_flags() which return a (handle, error)
 * pair, the other functions return a ([result,] GuestfsError) where
 * GuestfsError is defined here.
 */
type GuestfsError struct {
	Op     string        // operation which failed
	Errmsg string        // string (guestfs_last_error)
	Errno  syscall.Errno // errno (guestfs_last_errno)
}

func (e *GuestfsError) String() string {
	if e.Errno != 0 {
		return fmt.Sprintf("%s: %s", e.Op, e.Errmsg)
	} else {
		return fmt.Sprintf("%s: %s: %s", e.Op, e.Errmsg, e.Errno)
	}
}

func get_error_from_handle(g *Guestfs, op string) *GuestfsError {
	// NB: DO NOT try to free c_errmsg!
	c_errmsg := C.guestfs_last_error(g.g)
	errmsg := C.GoString(c_errmsg)

	errno := syscall.Errno(C.guestfs_last_errno(g.g))

	return &GuestfsError{Op: op, Errmsg: errmsg, Errno: errno}
}

func closed_handle_error(op string) *GuestfsError {
	return &GuestfsError{Op: op, Errmsg: "handle is closed",
		Errno: syscall.Errno(0)}
}

/* Close the handle. */
func (g *Guestfs) Close() *GuestfsError {
	if g.g == nil {
		return closed_handle_error("close")
	}
	C.guestfs_close(g.g)
	g.g = nil
	return nil
}

/* Functions for translating between NULL-terminated lists of
 * C strings and golang []string.
 */
func arg_string_list(xs []string) **C.char {
	r := make([]*C.char, 1+len(xs))
	for i, x := range xs {
		r[i] = C.CString(x)
	}
	r[len(xs)] = nil
	return &r[0]
}

func count_string_list(argv **C.char) int {
	var i int
	for *argv != nil {
		i++
		argv = (**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(argv)) +
			unsafe.Sizeof(*argv)))
	}
	return i
}

func free_string_list(argv **C.char) {
	for *argv != nil {
		//C.free (*argv)
		argv = (**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(argv)) +
			unsafe.Sizeof(*argv)))
	}
}

func return_string_list(argv **C.char) []string {
	r := make([]string, count_string_list(argv))
	var i int
	for *argv != nil {
		r[i] = C.GoString(*argv)
		i++
		argv = (**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(argv)) +
			unsafe.Sizeof(*argv)))
	}
	return r
}

func return_hashtable(argv **C.char) map[string]string {
	r := make(map[string]string)
	for *argv != nil {
		key := C.GoString(*argv)
		argv = (**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(argv)) +
			unsafe.Sizeof(*argv)))
		if *argv == nil {
			panic("odd number of items in hash table")
		}

		r[key] = C.GoString(*argv)
		argv = (**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(argv)) +
			unsafe.Sizeof(*argv)))
	}
	return r
}

/* XXX Events/callbacks not yet implemented. */

type Application struct {
	app_name           string
	app_display_name   string
	app_epoch          int32
	app_version        string
	app_release        string
	app_install_path   string
	app_trans_path     string
	app_publisher      string
	app_url            string
	app_source_package string
	app_summary        string
	app_description    string
}

func return_Application(c *C.struct_guestfs_application) *Application {
	r := Application{}
	r.app_name = C.GoString(c.app_name)
	r.app_display_name = C.GoString(c.app_display_name)
	r.app_epoch = int32(c.app_epoch)
	r.app_version = C.GoString(c.app_version)
	r.app_release = C.GoString(c.app_release)
	r.app_install_path = C.GoString(c.app_install_path)
	r.app_trans_path = C.GoString(c.app_trans_path)
	r.app_publisher = C.GoString(c.app_publisher)
	r.app_url = C.GoString(c.app_url)
	r.app_source_package = C.GoString(c.app_source_package)
	r.app_summary = C.GoString(c.app_summary)
	r.app_description = C.GoString(c.app_description)
	return &r
}

func return_Application_list(c *C.struct_guestfs_application_list) *[]Application {
	nrelems := int(c.len)
	ptr := uintptr(unsafe.Pointer(c.val))
	elemsize := unsafe.Sizeof(*c.val)
	r := make([]Application, nrelems)
	for i := 0; i < nrelems; i++ {
		r[i] = *return_Application((*C.struct_guestfs_application)(unsafe.Pointer(ptr)))
		ptr += elemsize
	}
	return &r
}

type Application2 struct {
	app2_name           string
	app2_display_name   string
	app2_epoch          int32
	app2_version        string
	app2_release        string
	app2_arch           string
	app2_install_path   string
	app2_trans_path     string
	app2_publisher      string
	app2_url            string
	app2_source_package string
	app2_summary        string
	app2_description    string
	app2_spare1         string
	app2_spare2         string
	app2_spare3         string
	app2_spare4         string
}

func return_Application2(c *C.struct_guestfs_application2) *Application2 {
	r := Application2{}
	r.app2_name = C.GoString(c.app2_name)
	r.app2_display_name = C.GoString(c.app2_display_name)
	r.app2_epoch = int32(c.app2_epoch)
	r.app2_version = C.GoString(c.app2_version)
	r.app2_release = C.GoString(c.app2_release)
	r.app2_arch = C.GoString(c.app2_arch)
	r.app2_install_path = C.GoString(c.app2_install_path)
	r.app2_trans_path = C.GoString(c.app2_trans_path)
	r.app2_publisher = C.GoString(c.app2_publisher)
	r.app2_url = C.GoString(c.app2_url)
	r.app2_source_package = C.GoString(c.app2_source_package)
	r.app2_summary = C.GoString(c.app2_summary)
	r.app2_description = C.GoString(c.app2_description)
	r.app2_spare1 = C.GoString(c.app2_spare1)
	r.app2_spare2 = C.GoString(c.app2_spare2)
	r.app2_spare3 = C.GoString(c.app2_spare3)
	r.app2_spare4 = C.GoString(c.app2_spare4)
	return &r
}

func return_Application2_list(c *C.struct_guestfs_application2_list) *[]Application2 {
	nrelems := int(c.len)
	ptr := uintptr(unsafe.Pointer(c.val))
	elemsize := unsafe.Sizeof(*c.val)
	r := make([]Application2, nrelems)
	for i := 0; i < nrelems; i++ {
		r[i] = *return_Application2((*C.struct_guestfs_application2)(unsafe.Pointer(ptr)))
		ptr += elemsize
	}
	return &r
}

type BTRFSBalance struct {
	btrfsbalance_status     string
	btrfsbalance_total      uint64
	btrfsbalance_balanced   uint64
	btrfsbalance_considered uint64
	btrfsbalance_left       uint64
}

func return_BTRFSBalance(c *C.struct_guestfs_btrfsbalance) *BTRFSBalance {
	r := BTRFSBalance{}
	r.btrfsbalance_status = C.GoString(c.btrfsbalance_status)
	r.btrfsbalance_total = uint64(c.btrfsbalance_total)
	r.btrfsbalance_balanced = uint64(c.btrfsbalance_balanced)
	r.btrfsbalance_considered = uint64(c.btrfsbalance_considered)
	r.btrfsbalance_left = uint64(c.btrfsbalance_left)
	return &r
}

func return_BTRFSBalance_list(c *C.struct_guestfs_btrfsbalance_list) *[]BTRFSBalance {
	nrelems := int(c.len)
	ptr := uintptr(unsafe.Pointer(c.val))
	elemsize := unsafe.Sizeof(*c.val)
	r := make([]BTRFSBalance, nrelems)
	for i := 0; i < nrelems; i++ {
		r[i] = *return_BTRFSBalance((*C.struct_guestfs_btrfsbalance)(unsafe.Pointer(ptr)))
		ptr += elemsize
	}
	return &r
}

type BTRFSQgroup struct {
	btrfsqgroup_id   string
	btrfsqgroup_rfer uint64
	btrfsqgroup_excl uint64
}

func return_BTRFSQgroup(c *C.struct_guestfs_btrfsqgroup) *BTRFSQgroup {
	r := BTRFSQgroup{}
	r.btrfsqgroup_id = C.GoString(c.btrfsqgroup_id)
	r.btrfsqgroup_rfer = uint64(c.btrfsqgroup_rfer)
	r.btrfsqgroup_excl = uint64(c.btrfsqgroup_excl)
	return &r
}

func return_BTRFSQgroup_list(c *C.struct_guestfs_btrfsqgroup_list) *[]BTRFSQgroup {
	nrelems := int(c.len)
	ptr := uintptr(unsafe.Pointer(c.val))
	elemsize := unsafe.Sizeof(*c.val)
	r := make([]BTRFSQgroup, nrelems)
	for i := 0; i < nrelems; i++ {
		r[i] = *return_BTRFSQgroup((*C.struct_guestfs_btrfsqgroup)(unsafe.Pointer(ptr)))
		ptr += elemsize
	}
	return &r
}

type BTRFSScrub struct {
	btrfsscrub_data_extents_scrubbed uint64
	btrfsscrub_tree_extents_scrubbed uint64
	btrfsscrub_data_bytes_scrubbed   uint64
	btrfsscrub_tree_bytes_scrubbed   uint64
	btrfsscrub_read_errors           uint64
	btrfsscrub_csum_errors           uint64
	btrfsscrub_verify_errors         uint64
	btrfsscrub_no_csum               uint64
	btrfsscrub_csum_discards         uint64
	btrfsscrub_super_errors          uint64
	btrfsscrub_malloc_errors         uint64
	btrfsscrub_uncorrectable_errors  uint64
	btrfsscrub_unverified_errors     uint64
	btrfsscrub_corrected_errors      uint64
	btrfsscrub_last_physical         uint64
}

func return_BTRFSScrub(c *C.struct_guestfs_btrfsscrub) *BTRFSScrub {
	r := BTRFSScrub{}
	r.btrfsscrub_data_extents_scrubbed = uint64(c.btrfsscrub_data_extents_scrubbed)
	r.btrfsscrub_tree_extents_scrubbed = uint64(c.btrfsscrub_tree_extents_scrubbed)
	r.btrfsscrub_data_bytes_scrubbed = uint64(c.btrfsscrub_data_bytes_scrubbed)
	r.btrfsscrub_tree_bytes_scrubbed = uint64(c.btrfsscrub_tree_bytes_scrubbed)
	r.btrfsscrub_read_errors = uint64(c.btrfsscrub_read_errors)
	r.btrfsscrub_csum_errors = uint64(c.btrfsscrub_csum_errors)
	r.btrfsscrub_verify_errors = uint64(c.btrfsscrub_verify_errors)
	r.btrfsscrub_no_csum = uint64(c.btrfsscrub_no_csum)
	r.btrfsscrub_csum_discards = uint64(c.btrfsscrub_csum_discards)
	r.btrfsscrub_super_errors = uint64(c.btrfsscrub_super_errors)
	r.btrfsscrub_malloc_errors = uint64(c.btrfsscrub_malloc_errors)
	r.btrfsscrub_uncorrectable_errors = uint64(c.btrfsscrub_uncorrectable_errors)
	r.btrfsscrub_unverified_errors = uint64(c.btrfsscrub_unverified_errors)
	r.btrfsscrub_corrected_errors = uint64(c.btrfsscrub_corrected_errors)
	r.btrfsscrub_last_physical = uint64(c.btrfsscrub_last_physical)
	return &r
}

func return_BTRFSScrub_list(c *C.struct_guestfs_btrfsscrub_list) *[]BTRFSScrub {
	nrelems := int(c.len)
	ptr := uintptr(unsafe.Pointer(c.val))
	elemsize := unsafe.Sizeof(*c.val)
	r := make([]BTRFSScrub, nrelems)
	for i := 0; i < nrelems; i++ {
		r[i] = *return_BTRFSScrub((*C.struct_guestfs_btrfsscrub)(unsafe.Pointer(ptr)))
		ptr += elemsize
	}
	return &r
}

type BTRFSSubvolume struct {
	btrfssubvolume_id           uint64
	btrfssubvolume_top_level_id uint64
	btrfssubvolume_path         string
}

func return_BTRFSSubvolume(c *C.struct_guestfs_btrfssubvolume) *BTRFSSubvolume {
	r := BTRFSSubvolume{}
	r.btrfssubvolume_id = uint64(c.btrfssubvolume_id)
	r.btrfssubvolume_top_level_id = uint64(c.btrfssubvolume_top_level_id)
	r.btrfssubvolume_path = C.GoString(c.btrfssubvolume_path)
	return &r
}

func return_BTRFSSubvolume_list(c *C.struct_guestfs_btrfssubvolume_list) *[]BTRFSSubvolume {
	nrelems := int(c.len)
	ptr := uintptr(unsafe.Pointer(c.val))
	elemsize := unsafe.Sizeof(*c.val)
	r := make([]BTRFSSubvolume, nrelems)
	for i := 0; i < nrelems; i++ {
		r[i] = *return_BTRFSSubvolume((*C.struct_guestfs_btrfssubvolume)(unsafe.Pointer(ptr)))
		ptr += elemsize
	}
	return &r
}

type Dirent struct {
	ino  int64
	ftyp byte
	name string
}

func return_Dirent(c *C.struct_guestfs_dirent) *Dirent {
	r := Dirent{}
	r.ino = int64(c.ino)
	r.ftyp = byte(c.ftyp)
	r.name = C.GoString(c.name)
	return &r
}

func return_Dirent_list(c *C.struct_guestfs_dirent_list) *[]Dirent {
	nrelems := int(c.len)
	ptr := uintptr(unsafe.Pointer(c.val))
	elemsize := unsafe.Sizeof(*c.val)
	r := make([]Dirent, nrelems)
	for i := 0; i < nrelems; i++ {
		r[i] = *return_Dirent((*C.struct_guestfs_dirent)(unsafe.Pointer(ptr)))
		ptr += elemsize
	}
	return &r
}

type HivexNode struct {
	hivex_node_h int64
}

func return_HivexNode(c *C.struct_guestfs_hivex_node) *HivexNode {
	r := HivexNode{}
	r.hivex_node_h = int64(c.hivex_node_h)
	return &r
}

func return_HivexNode_list(c *C.struct_guestfs_hivex_node_list) *[]HivexNode {
	nrelems := int(c.len)
	ptr := uintptr(unsafe.Pointer(c.val))
	elemsize := unsafe.Sizeof(*c.val)
	r := make([]HivexNode, nrelems)
	for i := 0; i < nrelems; i++ {
		r[i] = *return_HivexNode((*C.struct_guestfs_hivex_node)(unsafe.Pointer(ptr)))
		ptr += elemsize
	}
	return &r
}

type HivexValue struct {
	hivex_value_h int64
}

func return_HivexValue(c *C.struct_guestfs_hivex_value) *HivexValue {
	r := HivexValue{}
	r.hivex_value_h = int64(c.hivex_value_h)
	return &r
}

func return_HivexValue_list(c *C.struct_guestfs_hivex_value_list) *[]HivexValue {
	nrelems := int(c.len)
	ptr := uintptr(unsafe.Pointer(c.val))
	elemsize := unsafe.Sizeof(*c.val)
	r := make([]HivexValue, nrelems)
	for i := 0; i < nrelems; i++ {
		r[i] = *return_HivexValue((*C.struct_guestfs_hivex_value)(unsafe.Pointer(ptr)))
		ptr += elemsize
	}
	return &r
}

type INotifyEvent struct {
	in_wd     int64
	in_mask   uint32
	in_cookie uint32
	in_name   string
}

func return_INotifyEvent(c *C.struct_guestfs_inotify_event) *INotifyEvent {
	r := INotifyEvent{}
	r.in_wd = int64(c.in_wd)
	r.in_mask = uint32(c.in_mask)
	r.in_cookie = uint32(c.in_cookie)
	r.in_name = C.GoString(c.in_name)
	return &r
}

func return_INotifyEvent_list(c *C.struct_guestfs_inotify_event_list) *[]INotifyEvent {
	nrelems := int(c.len)
	ptr := uintptr(unsafe.Pointer(c.val))
	elemsize := unsafe.Sizeof(*c.val)
	r := make([]INotifyEvent, nrelems)
	for i := 0; i < nrelems; i++ {
		r[i] = *return_INotifyEvent((*C.struct_guestfs_inotify_event)(unsafe.Pointer(ptr)))
		ptr += elemsize
	}
	return &r
}

type IntBool struct {
	i int32
	b int32
}

func return_IntBool(c *C.struct_guestfs_int_bool) *IntBool {
	r := IntBool{}
	r.i = int32(c.i)
	r.b = int32(c.b)
	return &r
}

func return_IntBool_list(c *C.struct_guestfs_int_bool_list) *[]IntBool {
	nrelems := int(c.len)
	ptr := uintptr(unsafe.Pointer(c.val))
	elemsize := unsafe.Sizeof(*c.val)
	r := make([]IntBool, nrelems)
	for i := 0; i < nrelems; i++ {
		r[i] = *return_IntBool((*C.struct_guestfs_int_bool)(unsafe.Pointer(ptr)))
		ptr += elemsize
	}
	return &r
}

type ISOInfo struct {
	iso_system_id              string
	iso_volume_id              string
	iso_volume_space_size      uint32
	iso_volume_set_size        uint32
	iso_volume_sequence_number uint32
	iso_logical_block_size     uint32
	iso_volume_set_id          string
	iso_publisher_id           string
	iso_data_preparer_id       string
	iso_application_id         string
	iso_copyright_file_id      string
	iso_abstract_file_id       string
	iso_bibliographic_file_id  string
	iso_volume_creation_t      int64
	iso_volume_modification_t  int64
	iso_volume_expiration_t    int64
	iso_volume_effective_t     int64
}

func return_ISOInfo(c *C.struct_guestfs_isoinfo) *ISOInfo {
	r := ISOInfo{}
	r.iso_system_id = C.GoString(c.iso_system_id)
	r.iso_volume_id = C.GoString(c.iso_volume_id)
	r.iso_volume_space_size = uint32(c.iso_volume_space_size)
	r.iso_volume_set_size = uint32(c.iso_volume_set_size)
	r.iso_volume_sequence_number = uint32(c.iso_volume_sequence_number)
	r.iso_logical_block_size = uint32(c.iso_logical_block_size)
	r.iso_volume_set_id = C.GoString(c.iso_volume_set_id)
	r.iso_publisher_id = C.GoString(c.iso_publisher_id)
	r.iso_data_preparer_id = C.GoString(c.iso_data_preparer_id)
	r.iso_application_id = C.GoString(c.iso_application_id)
	r.iso_copyright_file_id = C.GoString(c.iso_copyright_file_id)
	r.iso_abstract_file_id = C.GoString(c.iso_abstract_file_id)
	r.iso_bibliographic_file_id = C.GoString(c.iso_bibliographic_file_id)
	r.iso_volume_creation_t = int64(c.iso_volume_creation_t)
	r.iso_volume_modification_t = int64(c.iso_volume_modification_t)
	r.iso_volume_expiration_t = int64(c.iso_volume_expiration_t)
	r.iso_volume_effective_t = int64(c.iso_volume_effective_t)
	return &r
}

func return_ISOInfo_list(c *C.struct_guestfs_isoinfo_list) *[]ISOInfo {
	nrelems := int(c.len)
	ptr := uintptr(unsafe.Pointer(c.val))
	elemsize := unsafe.Sizeof(*c.val)
	r := make([]ISOInfo, nrelems)
	for i := 0; i < nrelems; i++ {
		r[i] = *return_ISOInfo((*C.struct_guestfs_isoinfo)(unsafe.Pointer(ptr)))
		ptr += elemsize
	}
	return &r
}

type LV struct {
	lv_name         string
	lv_uuid         [32]byte
	lv_attr         string
	lv_major        int64
	lv_minor        int64
	lv_kernel_major int64
	lv_kernel_minor int64
	lv_size         uint64
	seg_count       int64
	origin          string
	snap_percent    float32
	copy_percent    float32
	move_pv         string
	lv_tags         string
	mirror_log      string
	modules         string
}

func return_LV(c *C.struct_guestfs_lvm_lv) *LV {
	r := LV{}
	r.lv_name = C.GoString(c.lv_name)
	// XXX doesn't work XXX r.lv_uuid = C.GoBytes (c.lv_uuid, len (c.lv_uuid))
	r.lv_uuid = [32]byte{}
	r.lv_attr = C.GoString(c.lv_attr)
	r.lv_major = int64(c.lv_major)
	r.lv_minor = int64(c.lv_minor)
	r.lv_kernel_major = int64(c.lv_kernel_major)
	r.lv_kernel_minor = int64(c.lv_kernel_minor)
	r.lv_size = uint64(c.lv_size)
	r.seg_count = int64(c.seg_count)
	r.origin = C.GoString(c.origin)
	r.snap_percent = float32(c.snap_percent)
	r.copy_percent = float32(c.copy_percent)
	r.move_pv = C.GoString(c.move_pv)
	r.lv_tags = C.GoString(c.lv_tags)
	r.mirror_log = C.GoString(c.mirror_log)
	r.modules = C.GoString(c.modules)
	return &r
}

func return_LV_list(c *C.struct_guestfs_lvm_lv_list) *[]LV {
	nrelems := int(c.len)
	ptr := uintptr(unsafe.Pointer(c.val))
	elemsize := unsafe.Sizeof(*c.val)
	r := make([]LV, nrelems)
	for i := 0; i < nrelems; i++ {
		r[i] = *return_LV((*C.struct_guestfs_lvm_lv)(unsafe.Pointer(ptr)))
		ptr += elemsize
	}
	return &r
}

type PV struct {
	pv_name           string
	pv_uuid           [32]byte
	pv_fmt            string
	pv_size           uint64
	dev_size          uint64
	pv_free           uint64
	pv_used           uint64
	pv_attr           string
	pv_pe_count       int64
	pv_pe_alloc_count int64
	pv_tags           string
	pe_start          uint64
	pv_mda_count      int64
	pv_mda_free       uint64
}

func return_PV(c *C.struct_guestfs_lvm_pv) *PV {
	r := PV{}
	r.pv_name = C.GoString(c.pv_name)
	// XXX doesn't work XXX r.pv_uuid = C.GoBytes (c.pv_uuid, len (c.pv_uuid))
	r.pv_uuid = [32]byte{}
	r.pv_fmt = C.GoString(c.pv_fmt)
	r.pv_size = uint64(c.pv_size)
	r.dev_size = uint64(c.dev_size)
	r.pv_free = uint64(c.pv_free)
	r.pv_used = uint64(c.pv_used)
	r.pv_attr = C.GoString(c.pv_attr)
	r.pv_pe_count = int64(c.pv_pe_count)
	r.pv_pe_alloc_count = int64(c.pv_pe_alloc_count)
	r.pv_tags = C.GoString(c.pv_tags)
	r.pe_start = uint64(c.pe_start)
	r.pv_mda_count = int64(c.pv_mda_count)
	r.pv_mda_free = uint64(c.pv_mda_free)
	return &r
}

func return_PV_list(c *C.struct_guestfs_lvm_pv_list) *[]PV {
	nrelems := int(c.len)
	ptr := uintptr(unsafe.Pointer(c.val))
	elemsize := unsafe.Sizeof(*c.val)
	r := make([]PV, nrelems)
	for i := 0; i < nrelems; i++ {
		r[i] = *return_PV((*C.struct_guestfs_lvm_pv)(unsafe.Pointer(ptr)))
		ptr += elemsize
	}
	return &r
}

type VG struct {
	vg_name         string
	vg_uuid         [32]byte
	vg_fmt          string
	vg_attr         string
	vg_size         uint64
	vg_free         uint64
	vg_sysid        string
	vg_extent_size  uint64
	vg_extent_count int64
	vg_free_count   int64
	max_lv          int64
	max_pv          int64
	pv_count        int64
	lv_count        int64
	snap_count      int64
	vg_seqno        int64
	vg_tags         string
	vg_mda_count    int64
	vg_mda_free     uint64
}

func return_VG(c *C.struct_guestfs_lvm_vg) *VG {
	r := VG{}
	r.vg_name = C.GoString(c.vg_name)
	// XXX doesn't work XXX r.vg_uuid = C.GoBytes (c.vg_uuid, len (c.vg_uuid))
	r.vg_uuid = [32]byte{}
	r.vg_fmt = C.GoString(c.vg_fmt)
	r.vg_attr = C.GoString(c.vg_attr)
	r.vg_size = uint64(c.vg_size)
	r.vg_free = uint64(c.vg_free)
	r.vg_sysid = C.GoString(c.vg_sysid)
	r.vg_extent_size = uint64(c.vg_extent_size)
	r.vg_extent_count = int64(c.vg_extent_count)
	r.vg_free_count = int64(c.vg_free_count)
	r.max_lv = int64(c.max_lv)
	r.max_pv = int64(c.max_pv)
	r.pv_count = int64(c.pv_count)
	r.lv_count = int64(c.lv_count)
	r.snap_count = int64(c.snap_count)
	r.vg_seqno = int64(c.vg_seqno)
	r.vg_tags = C.GoString(c.vg_tags)
	r.vg_mda_count = int64(c.vg_mda_count)
	r.vg_mda_free = uint64(c.vg_mda_free)
	return &r
}

func return_VG_list(c *C.struct_guestfs_lvm_vg_list) *[]VG {
	nrelems := int(c.len)
	ptr := uintptr(unsafe.Pointer(c.val))
	elemsize := unsafe.Sizeof(*c.val)
	r := make([]VG, nrelems)
	for i := 0; i < nrelems; i++ {
		r[i] = *return_VG((*C.struct_guestfs_lvm_vg)(unsafe.Pointer(ptr)))
		ptr += elemsize
	}
	return &r
}

type MDStat struct {
	mdstat_device string
	mdstat_index  int32
	mdstat_flags  string
}

func return_MDStat(c *C.struct_guestfs_mdstat) *MDStat {
	r := MDStat{}
	r.mdstat_device = C.GoString(c.mdstat_device)
	r.mdstat_index = int32(c.mdstat_index)
	r.mdstat_flags = C.GoString(c.mdstat_flags)
	return &r
}

func return_MDStat_list(c *C.struct_guestfs_mdstat_list) *[]MDStat {
	nrelems := int(c.len)
	ptr := uintptr(unsafe.Pointer(c.val))
	elemsize := unsafe.Sizeof(*c.val)
	r := make([]MDStat, nrelems)
	for i := 0; i < nrelems; i++ {
		r[i] = *return_MDStat((*C.struct_guestfs_mdstat)(unsafe.Pointer(ptr)))
		ptr += elemsize
	}
	return &r
}

type Partition struct {
	part_num   int32
	part_start uint64
	part_end   uint64
	part_size  uint64
}

func return_Partition(c *C.struct_guestfs_partition) *Partition {
	r := Partition{}
	r.part_num = int32(c.part_num)
	r.part_start = uint64(c.part_start)
	r.part_end = uint64(c.part_end)
	r.part_size = uint64(c.part_size)
	return &r
}

func return_Partition_list(c *C.struct_guestfs_partition_list) *[]Partition {
	nrelems := int(c.len)
	ptr := uintptr(unsafe.Pointer(c.val))
	elemsize := unsafe.Sizeof(*c.val)
	r := make([]Partition, nrelems)
	for i := 0; i < nrelems; i++ {
		r[i] = *return_Partition((*C.struct_guestfs_partition)(unsafe.Pointer(ptr)))
		ptr += elemsize
	}
	return &r
}

type Stat struct {
	dev     int64
	ino     int64
	mode    int64
	nlink   int64
	uid     int64
	gid     int64
	rdev    int64
	size    int64
	blksize int64
	blocks  int64
	atime   int64
	mtime   int64
	ctime   int64
}

func return_Stat(c *C.struct_guestfs_stat) *Stat {
	r := Stat{}
	r.dev = int64(c.dev)
	r.ino = int64(c.ino)
	r.mode = int64(c.mode)
	r.nlink = int64(c.nlink)
	r.uid = int64(c.uid)
	r.gid = int64(c.gid)
	r.rdev = int64(c.rdev)
	r.size = int64(c.size)
	r.blksize = int64(c.blksize)
	r.blocks = int64(c.blocks)
	r.atime = int64(c.atime)
	r.mtime = int64(c.mtime)
	r.ctime = int64(c.ctime)
	return &r
}

func return_Stat_list(c *C.struct_guestfs_stat_list) *[]Stat {
	nrelems := int(c.len)
	ptr := uintptr(unsafe.Pointer(c.val))
	elemsize := unsafe.Sizeof(*c.val)
	r := make([]Stat, nrelems)
	for i := 0; i < nrelems; i++ {
		r[i] = *return_Stat((*C.struct_guestfs_stat)(unsafe.Pointer(ptr)))
		ptr += elemsize
	}
	return &r
}

type StatNS struct {
	st_dev        int64
	st_ino        int64
	st_mode       int64
	st_nlink      int64
	st_uid        int64
	st_gid        int64
	st_rdev       int64
	st_size       int64
	st_blksize    int64
	st_blocks     int64
	st_atime_sec  int64
	st_atime_nsec int64
	st_mtime_sec  int64
	st_mtime_nsec int64
	st_ctime_sec  int64
	st_ctime_nsec int64
	st_spare1     int64
	st_spare2     int64
	st_spare3     int64
	st_spare4     int64
	st_spare5     int64
	st_spare6     int64
}

func return_StatNS(c *C.struct_guestfs_statns) *StatNS {
	r := StatNS{}
	r.st_dev = int64(c.st_dev)
	r.st_ino = int64(c.st_ino)
	r.st_mode = int64(c.st_mode)
	r.st_nlink = int64(c.st_nlink)
	r.st_uid = int64(c.st_uid)
	r.st_gid = int64(c.st_gid)
	r.st_rdev = int64(c.st_rdev)
	r.st_size = int64(c.st_size)
	r.st_blksize = int64(c.st_blksize)
	r.st_blocks = int64(c.st_blocks)
	r.st_atime_sec = int64(c.st_atime_sec)
	r.st_atime_nsec = int64(c.st_atime_nsec)
	r.st_mtime_sec = int64(c.st_mtime_sec)
	r.st_mtime_nsec = int64(c.st_mtime_nsec)
	r.st_ctime_sec = int64(c.st_ctime_sec)
	r.st_ctime_nsec = int64(c.st_ctime_nsec)
	r.st_spare1 = int64(c.st_spare1)
	r.st_spare2 = int64(c.st_spare2)
	r.st_spare3 = int64(c.st_spare3)
	r.st_spare4 = int64(c.st_spare4)
	r.st_spare5 = int64(c.st_spare5)
	r.st_spare6 = int64(c.st_spare6)
	return &r
}

func return_StatNS_list(c *C.struct_guestfs_statns_list) *[]StatNS {
	nrelems := int(c.len)
	ptr := uintptr(unsafe.Pointer(c.val))
	elemsize := unsafe.Sizeof(*c.val)
	r := make([]StatNS, nrelems)
	for i := 0; i < nrelems; i++ {
		r[i] = *return_StatNS((*C.struct_guestfs_statns)(unsafe.Pointer(ptr)))
		ptr += elemsize
	}
	return &r
}

type StatVFS struct {
	bsize   int64
	frsize  int64
	blocks  int64
	bfree   int64
	bavail  int64
	files   int64
	ffree   int64
	favail  int64
	fsid    int64
	flag    int64
	namemax int64
}

func return_StatVFS(c *C.struct_guestfs_statvfs) *StatVFS {
	r := StatVFS{}
	r.bsize = int64(c.bsize)
	r.frsize = int64(c.frsize)
	r.blocks = int64(c.blocks)
	r.bfree = int64(c.bfree)
	r.bavail = int64(c.bavail)
	r.files = int64(c.files)
	r.ffree = int64(c.ffree)
	r.favail = int64(c.favail)
	r.fsid = int64(c.fsid)
	r.flag = int64(c.flag)
	r.namemax = int64(c.namemax)
	return &r
}

func return_StatVFS_list(c *C.struct_guestfs_statvfs_list) *[]StatVFS {
	nrelems := int(c.len)
	ptr := uintptr(unsafe.Pointer(c.val))
	elemsize := unsafe.Sizeof(*c.val)
	r := make([]StatVFS, nrelems)
	for i := 0; i < nrelems; i++ {
		r[i] = *return_StatVFS((*C.struct_guestfs_statvfs)(unsafe.Pointer(ptr)))
		ptr += elemsize
	}
	return &r
}

type UTSName struct {
	uts_sysname string
	uts_release string
	uts_version string
	uts_machine string
}

func return_UTSName(c *C.struct_guestfs_utsname) *UTSName {
	r := UTSName{}
	r.uts_sysname = C.GoString(c.uts_sysname)
	r.uts_release = C.GoString(c.uts_release)
	r.uts_version = C.GoString(c.uts_version)
	r.uts_machine = C.GoString(c.uts_machine)
	return &r
}

func return_UTSName_list(c *C.struct_guestfs_utsname_list) *[]UTSName {
	nrelems := int(c.len)
	ptr := uintptr(unsafe.Pointer(c.val))
	elemsize := unsafe.Sizeof(*c.val)
	r := make([]UTSName, nrelems)
	for i := 0; i < nrelems; i++ {
		r[i] = *return_UTSName((*C.struct_guestfs_utsname)(unsafe.Pointer(ptr)))
		ptr += elemsize
	}
	return &r
}

type Version struct {
	major   int64
	minor   int64
	release int64
	extra   string
}

func return_Version(c *C.struct_guestfs_version) *Version {
	r := Version{}
	r.major = int64(c.major)
	r.minor = int64(c.minor)
	r.release = int64(c.release)
	r.extra = C.GoString(c.extra)
	return &r
}

func return_Version_list(c *C.struct_guestfs_version_list) *[]Version {
	nrelems := int(c.len)
	ptr := uintptr(unsafe.Pointer(c.val))
	elemsize := unsafe.Sizeof(*c.val)
	r := make([]Version, nrelems)
	for i := 0; i < nrelems; i++ {
		r[i] = *return_Version((*C.struct_guestfs_version)(unsafe.Pointer(ptr)))
		ptr += elemsize
	}
	return &r
}

type XAttr struct {
	attrname string
	attrval  []byte
}

func return_XAttr(c *C.struct_guestfs_xattr) *XAttr {
	r := XAttr{}
	r.attrname = C.GoString(c.attrname)
	r.attrval = C.GoBytes(unsafe.Pointer(c.attrval), C.int(c.attrval_len))
	return &r
}

func return_XAttr_list(c *C.struct_guestfs_xattr_list) *[]XAttr {
	nrelems := int(c.len)
	ptr := uintptr(unsafe.Pointer(c.val))
	elemsize := unsafe.Sizeof(*c.val)
	r := make([]XAttr, nrelems)
	for i := 0; i < nrelems; i++ {
		r[i] = *return_XAttr((*C.struct_guestfs_xattr)(unsafe.Pointer(ptr)))
		ptr += elemsize
	}
	return &r
}

type XFSInfo struct {
	xfs_mntpoint     string
	xfs_inodesize    uint32
	xfs_agcount      uint32
	xfs_agsize       uint32
	xfs_sectsize     uint32
	xfs_attr         uint32
	xfs_blocksize    uint32
	xfs_datablocks   uint64
	xfs_imaxpct      uint32
	xfs_sunit        uint32
	xfs_swidth       uint32
	xfs_dirversion   uint32
	xfs_dirblocksize uint32
	xfs_cimode       uint32
	xfs_logname      string
	xfs_logblocksize uint32
	xfs_logblocks    uint32
	xfs_logversion   uint32
	xfs_logsectsize  uint32
	xfs_logsunit     uint32
	xfs_lazycount    uint32
	xfs_rtname       string
	xfs_rtextsize    uint32
	xfs_rtblocks     uint64
	xfs_rtextents    uint64
}

func return_XFSInfo(c *C.struct_guestfs_xfsinfo) *XFSInfo {
	r := XFSInfo{}
	r.xfs_mntpoint = C.GoString(c.xfs_mntpoint)
	r.xfs_inodesize = uint32(c.xfs_inodesize)
	r.xfs_agcount = uint32(c.xfs_agcount)
	r.xfs_agsize = uint32(c.xfs_agsize)
	r.xfs_sectsize = uint32(c.xfs_sectsize)
	r.xfs_attr = uint32(c.xfs_attr)
	r.xfs_blocksize = uint32(c.xfs_blocksize)
	r.xfs_datablocks = uint64(c.xfs_datablocks)
	r.xfs_imaxpct = uint32(c.xfs_imaxpct)
	r.xfs_sunit = uint32(c.xfs_sunit)
	r.xfs_swidth = uint32(c.xfs_swidth)
	r.xfs_dirversion = uint32(c.xfs_dirversion)
	r.xfs_dirblocksize = uint32(c.xfs_dirblocksize)
	r.xfs_cimode = uint32(c.xfs_cimode)
	r.xfs_logname = C.GoString(c.xfs_logname)
	r.xfs_logblocksize = uint32(c.xfs_logblocksize)
	r.xfs_logblocks = uint32(c.xfs_logblocks)
	r.xfs_logversion = uint32(c.xfs_logversion)
	r.xfs_logsectsize = uint32(c.xfs_logsectsize)
	r.xfs_logsunit = uint32(c.xfs_logsunit)
	r.xfs_lazycount = uint32(c.xfs_lazycount)
	r.xfs_rtname = C.GoString(c.xfs_rtname)
	r.xfs_rtextsize = uint32(c.xfs_rtextsize)
	r.xfs_rtblocks = uint64(c.xfs_rtblocks)
	r.xfs_rtextents = uint64(c.xfs_rtextents)
	return &r
}

func return_XFSInfo_list(c *C.struct_guestfs_xfsinfo_list) *[]XFSInfo {
	nrelems := int(c.len)
	ptr := uintptr(unsafe.Pointer(c.val))
	elemsize := unsafe.Sizeof(*c.val)
	r := make([]XFSInfo, nrelems)
	for i := 0; i < nrelems; i++ {
		r[i] = *return_XFSInfo((*C.struct_guestfs_xfsinfo)(unsafe.Pointer(ptr)))
		ptr += elemsize
	}
	return &r
}

/* acl_delete_def_file : delete the default POSIX ACL of a directory */
func (g *Guestfs) Acl_delete_def_file(dir string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("acl_delete_def_file")
	}

	c_dir := C.CString(dir)
	defer C.free(unsafe.Pointer(c_dir))

	r := C.guestfs_acl_delete_def_file(g.g, c_dir)

	if r == -1 {
		return get_error_from_handle(g, "acl_delete_def_file")
	}
	return nil
}

/* acl_get_file : get the POSIX ACL attached to a file */
func (g *Guestfs) Acl_get_file(path string, acltype string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("acl_get_file")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	c_acltype := C.CString(acltype)
	defer C.free(unsafe.Pointer(c_acltype))

	r := C.guestfs_acl_get_file(g.g, c_path, c_acltype)

	if r == nil {
		return "", get_error_from_handle(g, "acl_get_file")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* acl_set_file : set the POSIX ACL attached to a file */
func (g *Guestfs) Acl_set_file(path string, acltype string, acl string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("acl_set_file")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	c_acltype := C.CString(acltype)
	defer C.free(unsafe.Pointer(c_acltype))

	c_acl := C.CString(acl)
	defer C.free(unsafe.Pointer(c_acl))

	r := C.guestfs_acl_set_file(g.g, c_path, c_acltype, c_acl)

	if r == -1 {
		return get_error_from_handle(g, "acl_set_file")
	}
	return nil
}

/* add_cdrom : add a CD-ROM disk image to examine */
func (g *Guestfs) Add_cdrom(filename string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("add_cdrom")
	}

	c_filename := C.CString(filename)
	defer C.free(unsafe.Pointer(c_filename))

	r := C.guestfs_add_cdrom(g.g, c_filename)

	if r == -1 {
		return get_error_from_handle(g, "add_cdrom")
	}
	return nil
}

/* Struct carrying optional arguments for Add_domain */
type OptargsAdd_domain struct {
	/* Libvirturi field is ignored unless Libvirturi_is_set == true */
	Libvirturi_is_set bool
	Libvirturi        string
	/* Readonly field is ignored unless Readonly_is_set == true */
	Readonly_is_set bool
	Readonly        bool
	/* Iface field is ignored unless Iface_is_set == true */
	Iface_is_set bool
	Iface        string
	/* Live field is ignored unless Live_is_set == true */
	Live_is_set bool
	Live        bool
	/* Allowuuid field is ignored unless Allowuuid_is_set == true */
	Allowuuid_is_set bool
	Allowuuid        bool
	/* Readonlydisk field is ignored unless Readonlydisk_is_set == true */
	Readonlydisk_is_set bool
	Readonlydisk        string
	/* Cachemode field is ignored unless Cachemode_is_set == true */
	Cachemode_is_set bool
	Cachemode        string
	/* Discard field is ignored unless Discard_is_set == true */
	Discard_is_set bool
	Discard        string
	/* Copyonread field is ignored unless Copyonread_is_set == true */
	Copyonread_is_set bool
	Copyonread        bool
}

/* add_domain : add the disk(s) from a named libvirt domain */
func (g *Guestfs) Add_domain(dom string, optargs *OptargsAdd_domain) (int, *GuestfsError) {
	if g.g == nil {
		return 0, closed_handle_error("add_domain")
	}

	c_dom := C.CString(dom)
	defer C.free(unsafe.Pointer(c_dom))
	c_optargs := C.struct_guestfs_add_domain_argv{}
	if optargs != nil {
		if optargs.Libvirturi_is_set {
			c_optargs.bitmask |= C.GUESTFS_ADD_DOMAIN_LIBVIRTURI_BITMASK
			c_optargs.libvirturi = C.CString(optargs.Libvirturi)
			defer C.free(unsafe.Pointer(c_optargs.libvirturi))
		}
		if optargs.Readonly_is_set {
			c_optargs.bitmask |= C.GUESTFS_ADD_DOMAIN_READONLY_BITMASK
			if optargs.Readonly {
				c_optargs.readonly = 1
			} else {
				c_optargs.readonly = 0
			}
		}
		if optargs.Iface_is_set {
			c_optargs.bitmask |= C.GUESTFS_ADD_DOMAIN_IFACE_BITMASK
			c_optargs.iface = C.CString(optargs.Iface)
			defer C.free(unsafe.Pointer(c_optargs.iface))
		}
		if optargs.Live_is_set {
			c_optargs.bitmask |= C.GUESTFS_ADD_DOMAIN_LIVE_BITMASK
			if optargs.Live {
				c_optargs.live = 1
			} else {
				c_optargs.live = 0
			}
		}
		if optargs.Allowuuid_is_set {
			c_optargs.bitmask |= C.GUESTFS_ADD_DOMAIN_ALLOWUUID_BITMASK
			if optargs.Allowuuid {
				c_optargs.allowuuid = 1
			} else {
				c_optargs.allowuuid = 0
			}
		}
		if optargs.Readonlydisk_is_set {
			c_optargs.bitmask |= C.GUESTFS_ADD_DOMAIN_READONLYDISK_BITMASK
			c_optargs.readonlydisk = C.CString(optargs.Readonlydisk)
			defer C.free(unsafe.Pointer(c_optargs.readonlydisk))
		}
		if optargs.Cachemode_is_set {
			c_optargs.bitmask |= C.GUESTFS_ADD_DOMAIN_CACHEMODE_BITMASK
			c_optargs.cachemode = C.CString(optargs.Cachemode)
			defer C.free(unsafe.Pointer(c_optargs.cachemode))
		}
		if optargs.Discard_is_set {
			c_optargs.bitmask |= C.GUESTFS_ADD_DOMAIN_DISCARD_BITMASK
			c_optargs.discard = C.CString(optargs.Discard)
			defer C.free(unsafe.Pointer(c_optargs.discard))
		}
		if optargs.Copyonread_is_set {
			c_optargs.bitmask |= C.GUESTFS_ADD_DOMAIN_COPYONREAD_BITMASK
			if optargs.Copyonread {
				c_optargs.copyonread = 1
			} else {
				c_optargs.copyonread = 0
			}
		}
	}

	r := C.guestfs_add_domain_argv(g.g, c_dom, &c_optargs)

	if r == -1 {
		return 0, get_error_from_handle(g, "add_domain")
	}
	return int(r), nil
}

/* Struct carrying optional arguments for Add_drive */
type OptargsAdd_drive struct {
	/* Readonly field is ignored unless Readonly_is_set == true */
	Readonly_is_set bool
	Readonly        bool
	/* Format field is ignored unless Format_is_set == true */
	Format_is_set bool
	Format        string
	/* Iface field is ignored unless Iface_is_set == true */
	Iface_is_set bool
	Iface        string
	/* Name field is ignored unless Name_is_set == true */
	Name_is_set bool
	Name        string
	/* Label field is ignored unless Label_is_set == true */
	Label_is_set bool
	Label        string
	/* Protocol field is ignored unless Protocol_is_set == true */
	Protocol_is_set bool
	Protocol        string
	/* Server field is ignored unless Server_is_set == true */
	Server_is_set bool
	Server        []string
	/* Username field is ignored unless Username_is_set == true */
	Username_is_set bool
	Username        string
	/* Secret field is ignored unless Secret_is_set == true */
	Secret_is_set bool
	Secret        string
	/* Cachemode field is ignored unless Cachemode_is_set == true */
	Cachemode_is_set bool
	Cachemode        string
	/* Discard field is ignored unless Discard_is_set == true */
	Discard_is_set bool
	Discard        string
	/* Copyonread field is ignored unless Copyonread_is_set == true */
	Copyonread_is_set bool
	Copyonread        bool
}

/* add_drive : add an image to examine or modify */
func (g *Guestfs) Add_drive(filename string, optargs *OptargsAdd_drive) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("add_drive")
	}

	c_filename := C.CString(filename)
	defer C.free(unsafe.Pointer(c_filename))
	c_optargs := C.struct_guestfs_add_drive_opts_argv{}
	if optargs != nil {
		if optargs.Readonly_is_set {
			c_optargs.bitmask |= C.GUESTFS_ADD_DRIVE_OPTS_READONLY_BITMASK
			if optargs.Readonly {
				c_optargs.readonly = 1
			} else {
				c_optargs.readonly = 0
			}
		}
		if optargs.Format_is_set {
			c_optargs.bitmask |= C.GUESTFS_ADD_DRIVE_OPTS_FORMAT_BITMASK
			c_optargs.format = C.CString(optargs.Format)
			defer C.free(unsafe.Pointer(c_optargs.format))
		}
		if optargs.Iface_is_set {
			c_optargs.bitmask |= C.GUESTFS_ADD_DRIVE_OPTS_IFACE_BITMASK
			c_optargs.iface = C.CString(optargs.Iface)
			defer C.free(unsafe.Pointer(c_optargs.iface))
		}
		if optargs.Name_is_set {
			c_optargs.bitmask |= C.GUESTFS_ADD_DRIVE_OPTS_NAME_BITMASK
			c_optargs.name = C.CString(optargs.Name)
			defer C.free(unsafe.Pointer(c_optargs.name))
		}
		if optargs.Label_is_set {
			c_optargs.bitmask |= C.GUESTFS_ADD_DRIVE_OPTS_LABEL_BITMASK
			c_optargs.label = C.CString(optargs.Label)
			defer C.free(unsafe.Pointer(c_optargs.label))
		}
		if optargs.Protocol_is_set {
			c_optargs.bitmask |= C.GUESTFS_ADD_DRIVE_OPTS_PROTOCOL_BITMASK
			c_optargs.protocol = C.CString(optargs.Protocol)
			defer C.free(unsafe.Pointer(c_optargs.protocol))
		}
		if optargs.Server_is_set {
			c_optargs.bitmask |= C.GUESTFS_ADD_DRIVE_OPTS_SERVER_BITMASK
			c_optargs.server = arg_string_list(optargs.Server)
			defer free_string_list(c_optargs.server)
		}
		if optargs.Username_is_set {
			c_optargs.bitmask |= C.GUESTFS_ADD_DRIVE_OPTS_USERNAME_BITMASK
			c_optargs.username = C.CString(optargs.Username)
			defer C.free(unsafe.Pointer(c_optargs.username))
		}
		if optargs.Secret_is_set {
			c_optargs.bitmask |= C.GUESTFS_ADD_DRIVE_OPTS_SECRET_BITMASK
			c_optargs.secret = C.CString(optargs.Secret)
			defer C.free(unsafe.Pointer(c_optargs.secret))
		}
		if optargs.Cachemode_is_set {
			c_optargs.bitmask |= C.GUESTFS_ADD_DRIVE_OPTS_CACHEMODE_BITMASK
			c_optargs.cachemode = C.CString(optargs.Cachemode)
			defer C.free(unsafe.Pointer(c_optargs.cachemode))
		}
		if optargs.Discard_is_set {
			c_optargs.bitmask |= C.GUESTFS_ADD_DRIVE_OPTS_DISCARD_BITMASK
			c_optargs.discard = C.CString(optargs.Discard)
			defer C.free(unsafe.Pointer(c_optargs.discard))
		}
		if optargs.Copyonread_is_set {
			c_optargs.bitmask |= C.GUESTFS_ADD_DRIVE_OPTS_COPYONREAD_BITMASK
			if optargs.Copyonread {
				c_optargs.copyonread = 1
			} else {
				c_optargs.copyonread = 0
			}
		}
	}

	r := C.guestfs_add_drive_opts_argv(g.g, c_filename, &c_optargs)

	if r == -1 {
		return get_error_from_handle(g, "add_drive")
	}
	return nil
}

/* add_drive_ro : add a drive in snapshot mode (read-only) */
func (g *Guestfs) Add_drive_ro(filename string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("add_drive_ro")
	}

	c_filename := C.CString(filename)
	defer C.free(unsafe.Pointer(c_filename))

	r := C.guestfs_add_drive_ro(g.g, c_filename)

	if r == -1 {
		return get_error_from_handle(g, "add_drive_ro")
	}
	return nil
}

/* add_drive_ro_with_if : add a drive read-only specifying the QEMU block emulation to use */
func (g *Guestfs) Add_drive_ro_with_if(filename string, iface string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("add_drive_ro_with_if")
	}

	c_filename := C.CString(filename)
	defer C.free(unsafe.Pointer(c_filename))

	c_iface := C.CString(iface)
	defer C.free(unsafe.Pointer(c_iface))

	r := C.guestfs_add_drive_ro_with_if(g.g, c_filename, c_iface)

	if r == -1 {
		return get_error_from_handle(g, "add_drive_ro_with_if")
	}
	return nil
}

/* Struct carrying optional arguments for Add_drive_scratch */
type OptargsAdd_drive_scratch struct {
	/* Name field is ignored unless Name_is_set == true */
	Name_is_set bool
	Name        string
	/* Label field is ignored unless Label_is_set == true */
	Label_is_set bool
	Label        string
}

/* add_drive_scratch : add a temporary scratch drive */
func (g *Guestfs) Add_drive_scratch(size int64, optargs *OptargsAdd_drive_scratch) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("add_drive_scratch")
	}
	c_optargs := C.struct_guestfs_add_drive_scratch_argv{}
	if optargs != nil {
		if optargs.Name_is_set {
			c_optargs.bitmask |= C.GUESTFS_ADD_DRIVE_SCRATCH_NAME_BITMASK
			c_optargs.name = C.CString(optargs.Name)
			defer C.free(unsafe.Pointer(c_optargs.name))
		}
		if optargs.Label_is_set {
			c_optargs.bitmask |= C.GUESTFS_ADD_DRIVE_SCRATCH_LABEL_BITMASK
			c_optargs.label = C.CString(optargs.Label)
			defer C.free(unsafe.Pointer(c_optargs.label))
		}
	}

	r := C.guestfs_add_drive_scratch_argv(g.g, C.int64_t(size), &c_optargs)

	if r == -1 {
		return get_error_from_handle(g, "add_drive_scratch")
	}
	return nil
}

/* add_drive_with_if : add a drive specifying the QEMU block emulation to use */
func (g *Guestfs) Add_drive_with_if(filename string, iface string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("add_drive_with_if")
	}

	c_filename := C.CString(filename)
	defer C.free(unsafe.Pointer(c_filename))

	c_iface := C.CString(iface)
	defer C.free(unsafe.Pointer(c_iface))

	r := C.guestfs_add_drive_with_if(g.g, c_filename, c_iface)

	if r == -1 {
		return get_error_from_handle(g, "add_drive_with_if")
	}
	return nil
}

/* Struct carrying optional arguments for Add_libvirt_dom */
type OptargsAdd_libvirt_dom struct {
	/* Readonly field is ignored unless Readonly_is_set == true */
	Readonly_is_set bool
	Readonly        bool
	/* Iface field is ignored unless Iface_is_set == true */
	Iface_is_set bool
	Iface        string
	/* Live field is ignored unless Live_is_set == true */
	Live_is_set bool
	Live        bool
	/* Readonlydisk field is ignored unless Readonlydisk_is_set == true */
	Readonlydisk_is_set bool
	Readonlydisk        string
	/* Cachemode field is ignored unless Cachemode_is_set == true */
	Cachemode_is_set bool
	Cachemode        string
	/* Discard field is ignored unless Discard_is_set == true */
	Discard_is_set bool
	Discard        string
	/* Copyonread field is ignored unless Copyonread_is_set == true */
	Copyonread_is_set bool
	Copyonread        bool
}

/* add_libvirt_dom : add the disk(s) from a libvirt domain */
func (g *Guestfs) Add_libvirt_dom(dom int64, optargs *OptargsAdd_libvirt_dom) (int, *GuestfsError) {
	if g.g == nil {
		return 0, closed_handle_error("add_libvirt_dom")
	}
	c_optargs := C.struct_guestfs_add_libvirt_dom_argv{}
	if optargs != nil {
		if optargs.Readonly_is_set {
			c_optargs.bitmask |= C.GUESTFS_ADD_LIBVIRT_DOM_READONLY_BITMASK
			if optargs.Readonly {
				c_optargs.readonly = 1
			} else {
				c_optargs.readonly = 0
			}
		}
		if optargs.Iface_is_set {
			c_optargs.bitmask |= C.GUESTFS_ADD_LIBVIRT_DOM_IFACE_BITMASK
			c_optargs.iface = C.CString(optargs.Iface)
			defer C.free(unsafe.Pointer(c_optargs.iface))
		}
		if optargs.Live_is_set {
			c_optargs.bitmask |= C.GUESTFS_ADD_LIBVIRT_DOM_LIVE_BITMASK
			if optargs.Live {
				c_optargs.live = 1
			} else {
				c_optargs.live = 0
			}
		}
		if optargs.Readonlydisk_is_set {
			c_optargs.bitmask |= C.GUESTFS_ADD_LIBVIRT_DOM_READONLYDISK_BITMASK
			c_optargs.readonlydisk = C.CString(optargs.Readonlydisk)
			defer C.free(unsafe.Pointer(c_optargs.readonlydisk))
		}
		if optargs.Cachemode_is_set {
			c_optargs.bitmask |= C.GUESTFS_ADD_LIBVIRT_DOM_CACHEMODE_BITMASK
			c_optargs.cachemode = C.CString(optargs.Cachemode)
			defer C.free(unsafe.Pointer(c_optargs.cachemode))
		}
		if optargs.Discard_is_set {
			c_optargs.bitmask |= C.GUESTFS_ADD_LIBVIRT_DOM_DISCARD_BITMASK
			c_optargs.discard = C.CString(optargs.Discard)
			defer C.free(unsafe.Pointer(c_optargs.discard))
		}
		if optargs.Copyonread_is_set {
			c_optargs.bitmask |= C.GUESTFS_ADD_LIBVIRT_DOM_COPYONREAD_BITMASK
			if optargs.Copyonread {
				c_optargs.copyonread = 1
			} else {
				c_optargs.copyonread = 0
			}
		}
	}

	r := C.guestfs_add_libvirt_dom_argv(g.g, nil, &c_optargs)

	if r == -1 {
		return 0, get_error_from_handle(g, "add_libvirt_dom")
	}
	return int(r), nil
}

/* aug_clear : clear Augeas path */
func (g *Guestfs) Aug_clear(augpath string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("aug_clear")
	}

	c_augpath := C.CString(augpath)
	defer C.free(unsafe.Pointer(c_augpath))

	r := C.guestfs_aug_clear(g.g, c_augpath)

	if r == -1 {
		return get_error_from_handle(g, "aug_clear")
	}
	return nil
}

/* aug_close : close the current Augeas handle */
func (g *Guestfs) Aug_close() *GuestfsError {
	if g.g == nil {
		return closed_handle_error("aug_close")
	}

	r := C.guestfs_aug_close(g.g)

	if r == -1 {
		return get_error_from_handle(g, "aug_close")
	}
	return nil
}

/* aug_defnode : define an Augeas node */
func (g *Guestfs) Aug_defnode(name string, expr string, val string) (*IntBool, *GuestfsError) {
	if g.g == nil {
		return &IntBool{}, closed_handle_error("aug_defnode")
	}

	c_name := C.CString(name)
	defer C.free(unsafe.Pointer(c_name))

	c_expr := C.CString(expr)
	defer C.free(unsafe.Pointer(c_expr))

	c_val := C.CString(val)
	defer C.free(unsafe.Pointer(c_val))

	r := C.guestfs_aug_defnode(g.g, c_name, c_expr, c_val)

	if r == nil {
		return &IntBool{}, get_error_from_handle(g, "aug_defnode")
	}
	defer C.guestfs_free_int_bool(r)
	return return_IntBool(r), nil
}

/* aug_defvar : define an Augeas variable */
func (g *Guestfs) Aug_defvar(name string, expr *string) (int, *GuestfsError) {
	if g.g == nil {
		return 0, closed_handle_error("aug_defvar")
	}

	c_name := C.CString(name)
	defer C.free(unsafe.Pointer(c_name))

	var c_expr *C.char = nil
	if expr != nil {
		c_expr = C.CString(*expr)
		defer C.free(unsafe.Pointer(c_expr))
	}

	r := C.guestfs_aug_defvar(g.g, c_name, c_expr)

	if r == -1 {
		return 0, get_error_from_handle(g, "aug_defvar")
	}
	return int(r), nil
}

/* aug_get : look up the value of an Augeas path */
func (g *Guestfs) Aug_get(augpath string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("aug_get")
	}

	c_augpath := C.CString(augpath)
	defer C.free(unsafe.Pointer(c_augpath))

	r := C.guestfs_aug_get(g.g, c_augpath)

	if r == nil {
		return "", get_error_from_handle(g, "aug_get")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* aug_init : create a new Augeas handle */
func (g *Guestfs) Aug_init(root string, flags int) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("aug_init")
	}

	c_root := C.CString(root)
	defer C.free(unsafe.Pointer(c_root))

	r := C.guestfs_aug_init(g.g, c_root, C.int(flags))

	if r == -1 {
		return get_error_from_handle(g, "aug_init")
	}
	return nil
}

/* aug_insert : insert a sibling Augeas node */
func (g *Guestfs) Aug_insert(augpath string, label string, before bool) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("aug_insert")
	}

	c_augpath := C.CString(augpath)
	defer C.free(unsafe.Pointer(c_augpath))

	c_label := C.CString(label)
	defer C.free(unsafe.Pointer(c_label))

	var c_before C.int
	if before {
		c_before = 1
	} else {
		c_before = 0
	}

	r := C.guestfs_aug_insert(g.g, c_augpath, c_label, c_before)

	if r == -1 {
		return get_error_from_handle(g, "aug_insert")
	}
	return nil
}

/* aug_label : return the label from an Augeas path expression */
func (g *Guestfs) Aug_label(augpath string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("aug_label")
	}

	c_augpath := C.CString(augpath)
	defer C.free(unsafe.Pointer(c_augpath))

	r := C.guestfs_aug_label(g.g, c_augpath)

	if r == nil {
		return "", get_error_from_handle(g, "aug_label")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* aug_load : load files into the tree */
func (g *Guestfs) Aug_load() *GuestfsError {
	if g.g == nil {
		return closed_handle_error("aug_load")
	}

	r := C.guestfs_aug_load(g.g)

	if r == -1 {
		return get_error_from_handle(g, "aug_load")
	}
	return nil
}

/* aug_ls : list Augeas nodes under augpath */
func (g *Guestfs) Aug_ls(augpath string) ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("aug_ls")
	}

	c_augpath := C.CString(augpath)
	defer C.free(unsafe.Pointer(c_augpath))

	r := C.guestfs_aug_ls(g.g, c_augpath)

	if r == nil {
		return nil, get_error_from_handle(g, "aug_ls")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* aug_match : return Augeas nodes which match augpath */
func (g *Guestfs) Aug_match(augpath string) ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("aug_match")
	}

	c_augpath := C.CString(augpath)
	defer C.free(unsafe.Pointer(c_augpath))

	r := C.guestfs_aug_match(g.g, c_augpath)

	if r == nil {
		return nil, get_error_from_handle(g, "aug_match")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* aug_mv : move Augeas node */
func (g *Guestfs) Aug_mv(src string, dest string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("aug_mv")
	}

	c_src := C.CString(src)
	defer C.free(unsafe.Pointer(c_src))

	c_dest := C.CString(dest)
	defer C.free(unsafe.Pointer(c_dest))

	r := C.guestfs_aug_mv(g.g, c_src, c_dest)

	if r == -1 {
		return get_error_from_handle(g, "aug_mv")
	}
	return nil
}

/* aug_rm : remove an Augeas path */
func (g *Guestfs) Aug_rm(augpath string) (int, *GuestfsError) {
	if g.g == nil {
		return 0, closed_handle_error("aug_rm")
	}

	c_augpath := C.CString(augpath)
	defer C.free(unsafe.Pointer(c_augpath))

	r := C.guestfs_aug_rm(g.g, c_augpath)

	if r == -1 {
		return 0, get_error_from_handle(g, "aug_rm")
	}
	return int(r), nil
}

/* aug_save : write all pending Augeas changes to disk */
func (g *Guestfs) Aug_save() *GuestfsError {
	if g.g == nil {
		return closed_handle_error("aug_save")
	}

	r := C.guestfs_aug_save(g.g)

	if r == -1 {
		return get_error_from_handle(g, "aug_save")
	}
	return nil
}

/* aug_set : set Augeas path to value */
func (g *Guestfs) Aug_set(augpath string, val string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("aug_set")
	}

	c_augpath := C.CString(augpath)
	defer C.free(unsafe.Pointer(c_augpath))

	c_val := C.CString(val)
	defer C.free(unsafe.Pointer(c_val))

	r := C.guestfs_aug_set(g.g, c_augpath, c_val)

	if r == -1 {
		return get_error_from_handle(g, "aug_set")
	}
	return nil
}

/* aug_setm : set multiple Augeas nodes */
func (g *Guestfs) Aug_setm(base string, sub *string, val string) (int, *GuestfsError) {
	if g.g == nil {
		return 0, closed_handle_error("aug_setm")
	}

	c_base := C.CString(base)
	defer C.free(unsafe.Pointer(c_base))

	var c_sub *C.char = nil
	if sub != nil {
		c_sub = C.CString(*sub)
		defer C.free(unsafe.Pointer(c_sub))
	}

	c_val := C.CString(val)
	defer C.free(unsafe.Pointer(c_val))

	r := C.guestfs_aug_setm(g.g, c_base, c_sub, c_val)

	if r == -1 {
		return 0, get_error_from_handle(g, "aug_setm")
	}
	return int(r), nil
}

/* available : test availability of some parts of the API */
func (g *Guestfs) Available(groups []string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("available")
	}

	c_groups := arg_string_list(groups)
	defer free_string_list(c_groups)

	r := C.guestfs_available(g.g, c_groups)

	if r == -1 {
		return get_error_from_handle(g, "available")
	}
	return nil
}

/* available_all_groups : return a list of all optional groups */
func (g *Guestfs) Available_all_groups() ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("available_all_groups")
	}

	r := C.guestfs_available_all_groups(g.g)

	if r == nil {
		return nil, get_error_from_handle(g, "available_all_groups")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* base64_in : upload base64-encoded data to file */
func (g *Guestfs) Base64_in(base64file string, filename string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("base64_in")
	}

	c_base64file := C.CString(base64file)
	defer C.free(unsafe.Pointer(c_base64file))

	c_filename := C.CString(filename)
	defer C.free(unsafe.Pointer(c_filename))

	r := C.guestfs_base64_in(g.g, c_base64file, c_filename)

	if r == -1 {
		return get_error_from_handle(g, "base64_in")
	}
	return nil
}

/* base64_out : download file and encode as base64 */
func (g *Guestfs) Base64_out(filename string, base64file string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("base64_out")
	}

	c_filename := C.CString(filename)
	defer C.free(unsafe.Pointer(c_filename))

	c_base64file := C.CString(base64file)
	defer C.free(unsafe.Pointer(c_base64file))

	r := C.guestfs_base64_out(g.g, c_filename, c_base64file)

	if r == -1 {
		return get_error_from_handle(g, "base64_out")
	}
	return nil
}

/* blkdiscard : discard all blocks on a device */
func (g *Guestfs) Blkdiscard(device string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("blkdiscard")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_blkdiscard(g.g, c_device)

	if r == -1 {
		return get_error_from_handle(g, "blkdiscard")
	}
	return nil
}

/* blkdiscardzeroes : return true if discarded blocks are read as zeroes */
func (g *Guestfs) Blkdiscardzeroes(device string) (bool, *GuestfsError) {
	if g.g == nil {
		return false, closed_handle_error("blkdiscardzeroes")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_blkdiscardzeroes(g.g, c_device)

	if r == -1 {
		return false, get_error_from_handle(g, "blkdiscardzeroes")
	}
	return r != 0, nil
}

/* blkid : print block device attributes */
func (g *Guestfs) Blkid(device string) (map[string]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("blkid")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_blkid(g.g, c_device)

	if r == nil {
		return nil, get_error_from_handle(g, "blkid")
	}
	defer free_string_list(r)
	return return_hashtable(r), nil
}

/* blockdev_flushbufs : flush device buffers */
func (g *Guestfs) Blockdev_flushbufs(device string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("blockdev_flushbufs")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_blockdev_flushbufs(g.g, c_device)

	if r == -1 {
		return get_error_from_handle(g, "blockdev_flushbufs")
	}
	return nil
}

/* blockdev_getbsz : get blocksize of block device */
func (g *Guestfs) Blockdev_getbsz(device string) (int, *GuestfsError) {
	if g.g == nil {
		return 0, closed_handle_error("blockdev_getbsz")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_blockdev_getbsz(g.g, c_device)

	if r == -1 {
		return 0, get_error_from_handle(g, "blockdev_getbsz")
	}
	return int(r), nil
}

/* blockdev_getro : is block device set to read-only */
func (g *Guestfs) Blockdev_getro(device string) (bool, *GuestfsError) {
	if g.g == nil {
		return false, closed_handle_error("blockdev_getro")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_blockdev_getro(g.g, c_device)

	if r == -1 {
		return false, get_error_from_handle(g, "blockdev_getro")
	}
	return r != 0, nil
}

/* blockdev_getsize64 : get total size of device in bytes */
func (g *Guestfs) Blockdev_getsize64(device string) (int64, *GuestfsError) {
	if g.g == nil {
		return 0, closed_handle_error("blockdev_getsize64")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_blockdev_getsize64(g.g, c_device)

	if r == -1 {
		return 0, get_error_from_handle(g, "blockdev_getsize64")
	}
	return int64(r), nil
}

/* blockdev_getss : get sectorsize of block device */
func (g *Guestfs) Blockdev_getss(device string) (int, *GuestfsError) {
	if g.g == nil {
		return 0, closed_handle_error("blockdev_getss")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_blockdev_getss(g.g, c_device)

	if r == -1 {
		return 0, get_error_from_handle(g, "blockdev_getss")
	}
	return int(r), nil
}

/* blockdev_getsz : get total size of device in 512-byte sectors */
func (g *Guestfs) Blockdev_getsz(device string) (int64, *GuestfsError) {
	if g.g == nil {
		return 0, closed_handle_error("blockdev_getsz")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_blockdev_getsz(g.g, c_device)

	if r == -1 {
		return 0, get_error_from_handle(g, "blockdev_getsz")
	}
	return int64(r), nil
}

/* blockdev_rereadpt : reread partition table */
func (g *Guestfs) Blockdev_rereadpt(device string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("blockdev_rereadpt")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_blockdev_rereadpt(g.g, c_device)

	if r == -1 {
		return get_error_from_handle(g, "blockdev_rereadpt")
	}
	return nil
}

/* blockdev_setbsz : set blocksize of block device */
func (g *Guestfs) Blockdev_setbsz(device string, blocksize int) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("blockdev_setbsz")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_blockdev_setbsz(g.g, c_device, C.int(blocksize))

	if r == -1 {
		return get_error_from_handle(g, "blockdev_setbsz")
	}
	return nil
}

/* blockdev_setra : set readahead */
func (g *Guestfs) Blockdev_setra(device string, sectors int) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("blockdev_setra")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_blockdev_setra(g.g, c_device, C.int(sectors))

	if r == -1 {
		return get_error_from_handle(g, "blockdev_setra")
	}
	return nil
}

/* blockdev_setro : set block device to read-only */
func (g *Guestfs) Blockdev_setro(device string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("blockdev_setro")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_blockdev_setro(g.g, c_device)

	if r == -1 {
		return get_error_from_handle(g, "blockdev_setro")
	}
	return nil
}

/* blockdev_setrw : set block device to read-write */
func (g *Guestfs) Blockdev_setrw(device string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("blockdev_setrw")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_blockdev_setrw(g.g, c_device)

	if r == -1 {
		return get_error_from_handle(g, "blockdev_setrw")
	}
	return nil
}

/* btrfs_balance_cancel : cancel a running or paused balance */
func (g *Guestfs) Btrfs_balance_cancel(path string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("btrfs_balance_cancel")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_btrfs_balance_cancel(g.g, c_path)

	if r == -1 {
		return get_error_from_handle(g, "btrfs_balance_cancel")
	}
	return nil
}

/* btrfs_balance_pause : pause a running balance */
func (g *Guestfs) Btrfs_balance_pause(path string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("btrfs_balance_pause")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_btrfs_balance_pause(g.g, c_path)

	if r == -1 {
		return get_error_from_handle(g, "btrfs_balance_pause")
	}
	return nil
}

/* btrfs_balance_resume : resume a paused balance */
func (g *Guestfs) Btrfs_balance_resume(path string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("btrfs_balance_resume")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_btrfs_balance_resume(g.g, c_path)

	if r == -1 {
		return get_error_from_handle(g, "btrfs_balance_resume")
	}
	return nil
}

/* btrfs_balance_status : show the status of a running or paused balance */
func (g *Guestfs) Btrfs_balance_status(path string) (*BTRFSBalance, *GuestfsError) {
	if g.g == nil {
		return &BTRFSBalance{}, closed_handle_error("btrfs_balance_status")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_btrfs_balance_status(g.g, c_path)

	if r == nil {
		return &BTRFSBalance{}, get_error_from_handle(g, "btrfs_balance_status")
	}
	defer C.guestfs_free_btrfsbalance(r)
	return return_BTRFSBalance(r), nil
}

/* btrfs_device_add : add devices to a btrfs filesystem */
func (g *Guestfs) Btrfs_device_add(devices []string, fs string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("btrfs_device_add")
	}

	c_devices := arg_string_list(devices)
	defer free_string_list(c_devices)

	c_fs := C.CString(fs)
	defer C.free(unsafe.Pointer(c_fs))

	r := C.guestfs_btrfs_device_add(g.g, c_devices, c_fs)

	if r == -1 {
		return get_error_from_handle(g, "btrfs_device_add")
	}
	return nil
}

/* btrfs_device_delete : remove devices from a btrfs filesystem */
func (g *Guestfs) Btrfs_device_delete(devices []string, fs string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("btrfs_device_delete")
	}

	c_devices := arg_string_list(devices)
	defer free_string_list(c_devices)

	c_fs := C.CString(fs)
	defer C.free(unsafe.Pointer(c_fs))

	r := C.guestfs_btrfs_device_delete(g.g, c_devices, c_fs)

	if r == -1 {
		return get_error_from_handle(g, "btrfs_device_delete")
	}
	return nil
}

/* btrfs_filesystem_balance : balance a btrfs filesystem */
func (g *Guestfs) Btrfs_filesystem_balance(fs string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("btrfs_filesystem_balance")
	}

	c_fs := C.CString(fs)
	defer C.free(unsafe.Pointer(c_fs))

	r := C.guestfs_btrfs_filesystem_balance(g.g, c_fs)

	if r == -1 {
		return get_error_from_handle(g, "btrfs_filesystem_balance")
	}
	return nil
}

/* Struct carrying optional arguments for Btrfs_filesystem_defragment */
type OptargsBtrfs_filesystem_defragment struct {
	/* Flush field is ignored unless Flush_is_set == true */
	Flush_is_set bool
	Flush        bool
	/* Compress field is ignored unless Compress_is_set == true */
	Compress_is_set bool
	Compress        string
}

/* btrfs_filesystem_defragment : defragment a file or directory */
func (g *Guestfs) Btrfs_filesystem_defragment(path string, optargs *OptargsBtrfs_filesystem_defragment) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("btrfs_filesystem_defragment")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))
	c_optargs := C.struct_guestfs_btrfs_filesystem_defragment_argv{}
	if optargs != nil {
		if optargs.Flush_is_set {
			c_optargs.bitmask |= C.GUESTFS_BTRFS_FILESYSTEM_DEFRAGMENT_FLUSH_BITMASK
			if optargs.Flush {
				c_optargs.flush = 1
			} else {
				c_optargs.flush = 0
			}
		}
		if optargs.Compress_is_set {
			c_optargs.bitmask |= C.GUESTFS_BTRFS_FILESYSTEM_DEFRAGMENT_COMPRESS_BITMASK
			c_optargs.compress = C.CString(optargs.Compress)
			defer C.free(unsafe.Pointer(c_optargs.compress))
		}
	}

	r := C.guestfs_btrfs_filesystem_defragment_argv(g.g, c_path, &c_optargs)

	if r == -1 {
		return get_error_from_handle(g, "btrfs_filesystem_defragment")
	}
	return nil
}

/* Struct carrying optional arguments for Btrfs_filesystem_resize */
type OptargsBtrfs_filesystem_resize struct {
	/* Size field is ignored unless Size_is_set == true */
	Size_is_set bool
	Size        int64
}

/* btrfs_filesystem_resize : resize a btrfs filesystem */
func (g *Guestfs) Btrfs_filesystem_resize(mountpoint string, optargs *OptargsBtrfs_filesystem_resize) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("btrfs_filesystem_resize")
	}

	c_mountpoint := C.CString(mountpoint)
	defer C.free(unsafe.Pointer(c_mountpoint))
	c_optargs := C.struct_guestfs_btrfs_filesystem_resize_argv{}
	if optargs != nil {
		if optargs.Size_is_set {
			c_optargs.bitmask |= C.GUESTFS_BTRFS_FILESYSTEM_RESIZE_SIZE_BITMASK
			c_optargs.size = C.int64_t(optargs.Size)
		}
	}

	r := C.guestfs_btrfs_filesystem_resize_argv(g.g, c_mountpoint, &c_optargs)

	if r == -1 {
		return get_error_from_handle(g, "btrfs_filesystem_resize")
	}
	return nil
}

/* btrfs_filesystem_sync : sync a btrfs filesystem */
func (g *Guestfs) Btrfs_filesystem_sync(fs string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("btrfs_filesystem_sync")
	}

	c_fs := C.CString(fs)
	defer C.free(unsafe.Pointer(c_fs))

	r := C.guestfs_btrfs_filesystem_sync(g.g, c_fs)

	if r == -1 {
		return get_error_from_handle(g, "btrfs_filesystem_sync")
	}
	return nil
}

/* Struct carrying optional arguments for Btrfs_fsck */
type OptargsBtrfs_fsck struct {
	/* Superblock field is ignored unless Superblock_is_set == true */
	Superblock_is_set bool
	Superblock        int64
	/* Repair field is ignored unless Repair_is_set == true */
	Repair_is_set bool
	Repair        bool
}

/* btrfs_fsck : check a btrfs filesystem */
func (g *Guestfs) Btrfs_fsck(device string, optargs *OptargsBtrfs_fsck) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("btrfs_fsck")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))
	c_optargs := C.struct_guestfs_btrfs_fsck_argv{}
	if optargs != nil {
		if optargs.Superblock_is_set {
			c_optargs.bitmask |= C.GUESTFS_BTRFS_FSCK_SUPERBLOCK_BITMASK
			c_optargs.superblock = C.int64_t(optargs.Superblock)
		}
		if optargs.Repair_is_set {
			c_optargs.bitmask |= C.GUESTFS_BTRFS_FSCK_REPAIR_BITMASK
			if optargs.Repair {
				c_optargs.repair = 1
			} else {
				c_optargs.repair = 0
			}
		}
	}

	r := C.guestfs_btrfs_fsck_argv(g.g, c_device, &c_optargs)

	if r == -1 {
		return get_error_from_handle(g, "btrfs_fsck")
	}
	return nil
}

/* Struct carrying optional arguments for Btrfs_image */
type OptargsBtrfs_image struct {
	/* Compresslevel field is ignored unless Compresslevel_is_set == true */
	Compresslevel_is_set bool
	Compresslevel        int
}

/* btrfs_image : create an image of a btrfs filesystem */
func (g *Guestfs) Btrfs_image(source []string, image string, optargs *OptargsBtrfs_image) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("btrfs_image")
	}

	c_source := arg_string_list(source)
	defer free_string_list(c_source)

	c_image := C.CString(image)
	defer C.free(unsafe.Pointer(c_image))
	c_optargs := C.struct_guestfs_btrfs_image_argv{}
	if optargs != nil {
		if optargs.Compresslevel_is_set {
			c_optargs.bitmask |= C.GUESTFS_BTRFS_IMAGE_COMPRESSLEVEL_BITMASK
			c_optargs.compresslevel = C.int(optargs.Compresslevel)
		}
	}

	r := C.guestfs_btrfs_image_argv(g.g, c_source, c_image, &c_optargs)

	if r == -1 {
		return get_error_from_handle(g, "btrfs_image")
	}
	return nil
}

/* btrfs_qgroup_assign : add a qgroup to a parent qgroup */
func (g *Guestfs) Btrfs_qgroup_assign(src string, dst string, path string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("btrfs_qgroup_assign")
	}

	c_src := C.CString(src)
	defer C.free(unsafe.Pointer(c_src))

	c_dst := C.CString(dst)
	defer C.free(unsafe.Pointer(c_dst))

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_btrfs_qgroup_assign(g.g, c_src, c_dst, c_path)

	if r == -1 {
		return get_error_from_handle(g, "btrfs_qgroup_assign")
	}
	return nil
}

/* btrfs_qgroup_create : create a subvolume quota group */
func (g *Guestfs) Btrfs_qgroup_create(qgroupid string, subvolume string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("btrfs_qgroup_create")
	}

	c_qgroupid := C.CString(qgroupid)
	defer C.free(unsafe.Pointer(c_qgroupid))

	c_subvolume := C.CString(subvolume)
	defer C.free(unsafe.Pointer(c_subvolume))

	r := C.guestfs_btrfs_qgroup_create(g.g, c_qgroupid, c_subvolume)

	if r == -1 {
		return get_error_from_handle(g, "btrfs_qgroup_create")
	}
	return nil
}

/* btrfs_qgroup_destroy : destroy a subvolume quota group */
func (g *Guestfs) Btrfs_qgroup_destroy(qgroupid string, subvolume string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("btrfs_qgroup_destroy")
	}

	c_qgroupid := C.CString(qgroupid)
	defer C.free(unsafe.Pointer(c_qgroupid))

	c_subvolume := C.CString(subvolume)
	defer C.free(unsafe.Pointer(c_subvolume))

	r := C.guestfs_btrfs_qgroup_destroy(g.g, c_qgroupid, c_subvolume)

	if r == -1 {
		return get_error_from_handle(g, "btrfs_qgroup_destroy")
	}
	return nil
}

/* btrfs_qgroup_limit : limit the size of a subvolume */
func (g *Guestfs) Btrfs_qgroup_limit(subvolume string, size int64) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("btrfs_qgroup_limit")
	}

	c_subvolume := C.CString(subvolume)
	defer C.free(unsafe.Pointer(c_subvolume))

	r := C.guestfs_btrfs_qgroup_limit(g.g, c_subvolume, C.int64_t(size))

	if r == -1 {
		return get_error_from_handle(g, "btrfs_qgroup_limit")
	}
	return nil
}

/* btrfs_qgroup_remove : remove a qgroup from its parent qgroup */
func (g *Guestfs) Btrfs_qgroup_remove(src string, dst string, path string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("btrfs_qgroup_remove")
	}

	c_src := C.CString(src)
	defer C.free(unsafe.Pointer(c_src))

	c_dst := C.CString(dst)
	defer C.free(unsafe.Pointer(c_dst))

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_btrfs_qgroup_remove(g.g, c_src, c_dst, c_path)

	if r == -1 {
		return get_error_from_handle(g, "btrfs_qgroup_remove")
	}
	return nil
}

/* btrfs_qgroup_show : show subvolume quota groups */
func (g *Guestfs) Btrfs_qgroup_show(path string) (*[]BTRFSQgroup, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("btrfs_qgroup_show")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_btrfs_qgroup_show(g.g, c_path)

	if r == nil {
		return nil, get_error_from_handle(g, "btrfs_qgroup_show")
	}
	defer C.guestfs_free_btrfsqgroup_list(r)
	return return_BTRFSQgroup_list(r), nil
}

/* btrfs_quota_enable : enable or disable subvolume quota support */
func (g *Guestfs) Btrfs_quota_enable(fs string, enable bool) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("btrfs_quota_enable")
	}

	c_fs := C.CString(fs)
	defer C.free(unsafe.Pointer(c_fs))

	var c_enable C.int
	if enable {
		c_enable = 1
	} else {
		c_enable = 0
	}

	r := C.guestfs_btrfs_quota_enable(g.g, c_fs, c_enable)

	if r == -1 {
		return get_error_from_handle(g, "btrfs_quota_enable")
	}
	return nil
}

/* btrfs_quota_rescan : trash all qgroup numbers and scan the metadata again with the current config */
func (g *Guestfs) Btrfs_quota_rescan(fs string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("btrfs_quota_rescan")
	}

	c_fs := C.CString(fs)
	defer C.free(unsafe.Pointer(c_fs))

	r := C.guestfs_btrfs_quota_rescan(g.g, c_fs)

	if r == -1 {
		return get_error_from_handle(g, "btrfs_quota_rescan")
	}
	return nil
}

/* btrfs_replace : replace a btrfs managed device with another device */
func (g *Guestfs) Btrfs_replace(srcdev string, targetdev string, mntpoint string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("btrfs_replace")
	}

	c_srcdev := C.CString(srcdev)
	defer C.free(unsafe.Pointer(c_srcdev))

	c_targetdev := C.CString(targetdev)
	defer C.free(unsafe.Pointer(c_targetdev))

	c_mntpoint := C.CString(mntpoint)
	defer C.free(unsafe.Pointer(c_mntpoint))

	r := C.guestfs_btrfs_replace(g.g, c_srcdev, c_targetdev, c_mntpoint)

	if r == -1 {
		return get_error_from_handle(g, "btrfs_replace")
	}
	return nil
}

/* btrfs_rescue_chunk_recover : recover the chunk tree of btrfs filesystem */
func (g *Guestfs) Btrfs_rescue_chunk_recover(device string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("btrfs_rescue_chunk_recover")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_btrfs_rescue_chunk_recover(g.g, c_device)

	if r == -1 {
		return get_error_from_handle(g, "btrfs_rescue_chunk_recover")
	}
	return nil
}

/* btrfs_rescue_super_recover : recover bad superblocks from good copies */
func (g *Guestfs) Btrfs_rescue_super_recover(device string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("btrfs_rescue_super_recover")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_btrfs_rescue_super_recover(g.g, c_device)

	if r == -1 {
		return get_error_from_handle(g, "btrfs_rescue_super_recover")
	}
	return nil
}

/* btrfs_scrub_cancel : cancel a running scrub */
func (g *Guestfs) Btrfs_scrub_cancel(path string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("btrfs_scrub_cancel")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_btrfs_scrub_cancel(g.g, c_path)

	if r == -1 {
		return get_error_from_handle(g, "btrfs_scrub_cancel")
	}
	return nil
}

/* btrfs_scrub_resume : resume a previously canceled or interrupted scrub */
func (g *Guestfs) Btrfs_scrub_resume(path string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("btrfs_scrub_resume")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_btrfs_scrub_resume(g.g, c_path)

	if r == -1 {
		return get_error_from_handle(g, "btrfs_scrub_resume")
	}
	return nil
}

/* btrfs_scrub_start : read all data from all disks and verify checksums */
func (g *Guestfs) Btrfs_scrub_start(path string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("btrfs_scrub_start")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_btrfs_scrub_start(g.g, c_path)

	if r == -1 {
		return get_error_from_handle(g, "btrfs_scrub_start")
	}
	return nil
}

/* btrfs_scrub_status : show status of running or finished scrub */
func (g *Guestfs) Btrfs_scrub_status(path string) (*BTRFSScrub, *GuestfsError) {
	if g.g == nil {
		return &BTRFSScrub{}, closed_handle_error("btrfs_scrub_status")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_btrfs_scrub_status(g.g, c_path)

	if r == nil {
		return &BTRFSScrub{}, get_error_from_handle(g, "btrfs_scrub_status")
	}
	defer C.guestfs_free_btrfsscrub(r)
	return return_BTRFSScrub(r), nil
}

/* btrfs_set_seeding : enable or disable the seeding feature of device */
func (g *Guestfs) Btrfs_set_seeding(device string, seeding bool) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("btrfs_set_seeding")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	var c_seeding C.int
	if seeding {
		c_seeding = 1
	} else {
		c_seeding = 0
	}

	r := C.guestfs_btrfs_set_seeding(g.g, c_device, c_seeding)

	if r == -1 {
		return get_error_from_handle(g, "btrfs_set_seeding")
	}
	return nil
}

/* Struct carrying optional arguments for Btrfs_subvolume_create */
type OptargsBtrfs_subvolume_create struct {
	/* Qgroupid field is ignored unless Qgroupid_is_set == true */
	Qgroupid_is_set bool
	Qgroupid        string
}

/* btrfs_subvolume_create : create a btrfs subvolume */
func (g *Guestfs) Btrfs_subvolume_create(dest string, optargs *OptargsBtrfs_subvolume_create) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("btrfs_subvolume_create")
	}

	c_dest := C.CString(dest)
	defer C.free(unsafe.Pointer(c_dest))
	c_optargs := C.struct_guestfs_btrfs_subvolume_create_opts_argv{}
	if optargs != nil {
		if optargs.Qgroupid_is_set {
			c_optargs.bitmask |= C.GUESTFS_BTRFS_SUBVOLUME_CREATE_OPTS_QGROUPID_BITMASK
			c_optargs.qgroupid = C.CString(optargs.Qgroupid)
			defer C.free(unsafe.Pointer(c_optargs.qgroupid))
		}
	}

	r := C.guestfs_btrfs_subvolume_create_opts_argv(g.g, c_dest, &c_optargs)

	if r == -1 {
		return get_error_from_handle(g, "btrfs_subvolume_create")
	}
	return nil
}

/* btrfs_subvolume_delete : delete a btrfs subvolume or snapshot */
func (g *Guestfs) Btrfs_subvolume_delete(subvolume string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("btrfs_subvolume_delete")
	}

	c_subvolume := C.CString(subvolume)
	defer C.free(unsafe.Pointer(c_subvolume))

	r := C.guestfs_btrfs_subvolume_delete(g.g, c_subvolume)

	if r == -1 {
		return get_error_from_handle(g, "btrfs_subvolume_delete")
	}
	return nil
}

/* btrfs_subvolume_get_default : get the default subvolume or snapshot of a filesystem */
func (g *Guestfs) Btrfs_subvolume_get_default(fs string) (int64, *GuestfsError) {
	if g.g == nil {
		return 0, closed_handle_error("btrfs_subvolume_get_default")
	}

	c_fs := C.CString(fs)
	defer C.free(unsafe.Pointer(c_fs))

	r := C.guestfs_btrfs_subvolume_get_default(g.g, c_fs)

	if r == -1 {
		return 0, get_error_from_handle(g, "btrfs_subvolume_get_default")
	}
	return int64(r), nil
}

/* btrfs_subvolume_list : list btrfs snapshots and subvolumes */
func (g *Guestfs) Btrfs_subvolume_list(fs string) (*[]BTRFSSubvolume, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("btrfs_subvolume_list")
	}

	c_fs := C.CString(fs)
	defer C.free(unsafe.Pointer(c_fs))

	r := C.guestfs_btrfs_subvolume_list(g.g, c_fs)

	if r == nil {
		return nil, get_error_from_handle(g, "btrfs_subvolume_list")
	}
	defer C.guestfs_free_btrfssubvolume_list(r)
	return return_BTRFSSubvolume_list(r), nil
}

/* btrfs_subvolume_set_default : set default btrfs subvolume */
func (g *Guestfs) Btrfs_subvolume_set_default(id int64, fs string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("btrfs_subvolume_set_default")
	}

	c_fs := C.CString(fs)
	defer C.free(unsafe.Pointer(c_fs))

	r := C.guestfs_btrfs_subvolume_set_default(g.g, C.int64_t(id), c_fs)

	if r == -1 {
		return get_error_from_handle(g, "btrfs_subvolume_set_default")
	}
	return nil
}

/* btrfs_subvolume_show : return detailed information of the subvolume */
func (g *Guestfs) Btrfs_subvolume_show(subvolume string) (map[string]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("btrfs_subvolume_show")
	}

	c_subvolume := C.CString(subvolume)
	defer C.free(unsafe.Pointer(c_subvolume))

	r := C.guestfs_btrfs_subvolume_show(g.g, c_subvolume)

	if r == nil {
		return nil, get_error_from_handle(g, "btrfs_subvolume_show")
	}
	defer free_string_list(r)
	return return_hashtable(r), nil
}

/* Struct carrying optional arguments for Btrfs_subvolume_snapshot */
type OptargsBtrfs_subvolume_snapshot struct {
	/* Ro field is ignored unless Ro_is_set == true */
	Ro_is_set bool
	Ro        bool
	/* Qgroupid field is ignored unless Qgroupid_is_set == true */
	Qgroupid_is_set bool
	Qgroupid        string
}

/* btrfs_subvolume_snapshot : create a btrfs snapshot */
func (g *Guestfs) Btrfs_subvolume_snapshot(source string, dest string, optargs *OptargsBtrfs_subvolume_snapshot) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("btrfs_subvolume_snapshot")
	}

	c_source := C.CString(source)
	defer C.free(unsafe.Pointer(c_source))

	c_dest := C.CString(dest)
	defer C.free(unsafe.Pointer(c_dest))
	c_optargs := C.struct_guestfs_btrfs_subvolume_snapshot_opts_argv{}
	if optargs != nil {
		if optargs.Ro_is_set {
			c_optargs.bitmask |= C.GUESTFS_BTRFS_SUBVOLUME_SNAPSHOT_OPTS_RO_BITMASK
			if optargs.Ro {
				c_optargs.ro = 1
			} else {
				c_optargs.ro = 0
			}
		}
		if optargs.Qgroupid_is_set {
			c_optargs.bitmask |= C.GUESTFS_BTRFS_SUBVOLUME_SNAPSHOT_OPTS_QGROUPID_BITMASK
			c_optargs.qgroupid = C.CString(optargs.Qgroupid)
			defer C.free(unsafe.Pointer(c_optargs.qgroupid))
		}
	}

	r := C.guestfs_btrfs_subvolume_snapshot_opts_argv(g.g, c_source, c_dest, &c_optargs)

	if r == -1 {
		return get_error_from_handle(g, "btrfs_subvolume_snapshot")
	}
	return nil
}

/* btrfstune_enable_extended_inode_refs : enable extended inode refs */
func (g *Guestfs) Btrfstune_enable_extended_inode_refs(device string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("btrfstune_enable_extended_inode_refs")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_btrfstune_enable_extended_inode_refs(g.g, c_device)

	if r == -1 {
		return get_error_from_handle(g, "btrfstune_enable_extended_inode_refs")
	}
	return nil
}

/* btrfstune_enable_skinny_metadata_extent_refs : enable skinny metadata extent refs */
func (g *Guestfs) Btrfstune_enable_skinny_metadata_extent_refs(device string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("btrfstune_enable_skinny_metadata_extent_refs")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_btrfstune_enable_skinny_metadata_extent_refs(g.g, c_device)

	if r == -1 {
		return get_error_from_handle(g, "btrfstune_enable_skinny_metadata_extent_refs")
	}
	return nil
}

/* btrfstune_seeding : enable or disable seeding of a btrfs device */
func (g *Guestfs) Btrfstune_seeding(device string, seeding bool) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("btrfstune_seeding")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	var c_seeding C.int
	if seeding {
		c_seeding = 1
	} else {
		c_seeding = 0
	}

	r := C.guestfs_btrfstune_seeding(g.g, c_device, c_seeding)

	if r == -1 {
		return get_error_from_handle(g, "btrfstune_seeding")
	}
	return nil
}

/* c_pointer : return the C pointer to the guestfs_h handle */
func (g *Guestfs) C_pointer() (int64, *GuestfsError) {
	if g.g == nil {
		return 0, closed_handle_error("c_pointer")
	}

	r := C.guestfs_c_pointer(g.g)

	if r == -1 {
		return 0, get_error_from_handle(g, "c_pointer")
	}
	return int64(r), nil
}

/* canonical_device_name : return canonical device name */
func (g *Guestfs) Canonical_device_name(device string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("canonical_device_name")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_canonical_device_name(g.g, c_device)

	if r == nil {
		return "", get_error_from_handle(g, "canonical_device_name")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* cap_get_file : get the Linux capabilities attached to a file */
func (g *Guestfs) Cap_get_file(path string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("cap_get_file")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_cap_get_file(g.g, c_path)

	if r == nil {
		return "", get_error_from_handle(g, "cap_get_file")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* cap_set_file : set the Linux capabilities attached to a file */
func (g *Guestfs) Cap_set_file(path string, cap string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("cap_set_file")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	c_cap := C.CString(cap)
	defer C.free(unsafe.Pointer(c_cap))

	r := C.guestfs_cap_set_file(g.g, c_path, c_cap)

	if r == -1 {
		return get_error_from_handle(g, "cap_set_file")
	}
	return nil
}

/* case_sensitive_path : return true path on case-insensitive filesystem */
func (g *Guestfs) Case_sensitive_path(path string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("case_sensitive_path")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_case_sensitive_path(g.g, c_path)

	if r == nil {
		return "", get_error_from_handle(g, "case_sensitive_path")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* cat : list the contents of a file */
func (g *Guestfs) Cat(path string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("cat")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_cat(g.g, c_path)

	if r == nil {
		return "", get_error_from_handle(g, "cat")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* checksum : compute MD5, SHAx or CRC checksum of file */
func (g *Guestfs) Checksum(csumtype string, path string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("checksum")
	}

	c_csumtype := C.CString(csumtype)
	defer C.free(unsafe.Pointer(c_csumtype))

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_checksum(g.g, c_csumtype, c_path)

	if r == nil {
		return "", get_error_from_handle(g, "checksum")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* checksum_device : compute MD5, SHAx or CRC checksum of the contents of a device */
func (g *Guestfs) Checksum_device(csumtype string, device string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("checksum_device")
	}

	c_csumtype := C.CString(csumtype)
	defer C.free(unsafe.Pointer(c_csumtype))

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_checksum_device(g.g, c_csumtype, c_device)

	if r == nil {
		return "", get_error_from_handle(g, "checksum_device")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* checksums_out : compute MD5, SHAx or CRC checksum of files in a directory */
func (g *Guestfs) Checksums_out(csumtype string, directory string, sumsfile string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("checksums_out")
	}

	c_csumtype := C.CString(csumtype)
	defer C.free(unsafe.Pointer(c_csumtype))

	c_directory := C.CString(directory)
	defer C.free(unsafe.Pointer(c_directory))

	c_sumsfile := C.CString(sumsfile)
	defer C.free(unsafe.Pointer(c_sumsfile))

	r := C.guestfs_checksums_out(g.g, c_csumtype, c_directory, c_sumsfile)

	if r == -1 {
		return get_error_from_handle(g, "checksums_out")
	}
	return nil
}

/* chmod : change file mode */
func (g *Guestfs) Chmod(mode int, path string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("chmod")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_chmod(g.g, C.int(mode), c_path)

	if r == -1 {
		return get_error_from_handle(g, "chmod")
	}
	return nil
}

/* chown : change file owner and group */
func (g *Guestfs) Chown(owner int, group int, path string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("chown")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_chown(g.g, C.int(owner), C.int(group), c_path)

	if r == -1 {
		return get_error_from_handle(g, "chown")
	}
	return nil
}

/* clear_backend_setting : remove a single per-backend settings string */
func (g *Guestfs) Clear_backend_setting(name string) (int, *GuestfsError) {
	if g.g == nil {
		return 0, closed_handle_error("clear_backend_setting")
	}

	c_name := C.CString(name)
	defer C.free(unsafe.Pointer(c_name))

	r := C.guestfs_clear_backend_setting(g.g, c_name)

	if r == -1 {
		return 0, get_error_from_handle(g, "clear_backend_setting")
	}
	return int(r), nil
}

/* command : run a command from the guest filesystem */
func (g *Guestfs) Command(arguments []string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("command")
	}

	c_arguments := arg_string_list(arguments)
	defer free_string_list(c_arguments)

	r := C.guestfs_command(g.g, c_arguments)

	if r == nil {
		return "", get_error_from_handle(g, "command")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* command_lines : run a command, returning lines */
func (g *Guestfs) Command_lines(arguments []string) ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("command_lines")
	}

	c_arguments := arg_string_list(arguments)
	defer free_string_list(c_arguments)

	r := C.guestfs_command_lines(g.g, c_arguments)

	if r == nil {
		return nil, get_error_from_handle(g, "command_lines")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* Struct carrying optional arguments for Compress_device_out */
type OptargsCompress_device_out struct {
	/* Level field is ignored unless Level_is_set == true */
	Level_is_set bool
	Level        int
}

/* compress_device_out : output compressed device */
func (g *Guestfs) Compress_device_out(ctype string, device string, zdevice string, optargs *OptargsCompress_device_out) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("compress_device_out")
	}

	c_ctype := C.CString(ctype)
	defer C.free(unsafe.Pointer(c_ctype))

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	c_zdevice := C.CString(zdevice)
	defer C.free(unsafe.Pointer(c_zdevice))
	c_optargs := C.struct_guestfs_compress_device_out_argv{}
	if optargs != nil {
		if optargs.Level_is_set {
			c_optargs.bitmask |= C.GUESTFS_COMPRESS_DEVICE_OUT_LEVEL_BITMASK
			c_optargs.level = C.int(optargs.Level)
		}
	}

	r := C.guestfs_compress_device_out_argv(g.g, c_ctype, c_device, c_zdevice, &c_optargs)

	if r == -1 {
		return get_error_from_handle(g, "compress_device_out")
	}
	return nil
}

/* Struct carrying optional arguments for Compress_out */
type OptargsCompress_out struct {
	/* Level field is ignored unless Level_is_set == true */
	Level_is_set bool
	Level        int
}

/* compress_out : output compressed file */
func (g *Guestfs) Compress_out(ctype string, file string, zfile string, optargs *OptargsCompress_out) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("compress_out")
	}

	c_ctype := C.CString(ctype)
	defer C.free(unsafe.Pointer(c_ctype))

	c_file := C.CString(file)
	defer C.free(unsafe.Pointer(c_file))

	c_zfile := C.CString(zfile)
	defer C.free(unsafe.Pointer(c_zfile))
	c_optargs := C.struct_guestfs_compress_out_argv{}
	if optargs != nil {
		if optargs.Level_is_set {
			c_optargs.bitmask |= C.GUESTFS_COMPRESS_OUT_LEVEL_BITMASK
			c_optargs.level = C.int(optargs.Level)
		}
	}

	r := C.guestfs_compress_out_argv(g.g, c_ctype, c_file, c_zfile, &c_optargs)

	if r == -1 {
		return get_error_from_handle(g, "compress_out")
	}
	return nil
}

/* config : add hypervisor parameters */
func (g *Guestfs) Config(hvparam string, hvvalue *string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("config")
	}

	c_hvparam := C.CString(hvparam)
	defer C.free(unsafe.Pointer(c_hvparam))

	var c_hvvalue *C.char = nil
	if hvvalue != nil {
		c_hvvalue = C.CString(*hvvalue)
		defer C.free(unsafe.Pointer(c_hvvalue))
	}

	r := C.guestfs_config(g.g, c_hvparam, c_hvvalue)

	if r == -1 {
		return get_error_from_handle(g, "config")
	}
	return nil
}

/* Struct carrying optional arguments for Copy_attributes */
type OptargsCopy_attributes struct {
	/* All field is ignored unless All_is_set == true */
	All_is_set bool
	All        bool
	/* Mode field is ignored unless Mode_is_set == true */
	Mode_is_set bool
	Mode        bool
	/* Xattributes field is ignored unless Xattributes_is_set == true */
	Xattributes_is_set bool
	Xattributes        bool
	/* Ownership field is ignored unless Ownership_is_set == true */
	Ownership_is_set bool
	Ownership        bool
}

/* copy_attributes : copy the attributes of a path (file/directory) to another */
func (g *Guestfs) Copy_attributes(src string, dest string, optargs *OptargsCopy_attributes) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("copy_attributes")
	}

	c_src := C.CString(src)
	defer C.free(unsafe.Pointer(c_src))

	c_dest := C.CString(dest)
	defer C.free(unsafe.Pointer(c_dest))
	c_optargs := C.struct_guestfs_copy_attributes_argv{}
	if optargs != nil {
		if optargs.All_is_set {
			c_optargs.bitmask |= C.GUESTFS_COPY_ATTRIBUTES_ALL_BITMASK
			if optargs.All {
				c_optargs.all = 1
			} else {
				c_optargs.all = 0
			}
		}
		if optargs.Mode_is_set {
			c_optargs.bitmask |= C.GUESTFS_COPY_ATTRIBUTES_MODE_BITMASK
			if optargs.Mode {
				c_optargs.mode = 1
			} else {
				c_optargs.mode = 0
			}
		}
		if optargs.Xattributes_is_set {
			c_optargs.bitmask |= C.GUESTFS_COPY_ATTRIBUTES_XATTRIBUTES_BITMASK
			if optargs.Xattributes {
				c_optargs.xattributes = 1
			} else {
				c_optargs.xattributes = 0
			}
		}
		if optargs.Ownership_is_set {
			c_optargs.bitmask |= C.GUESTFS_COPY_ATTRIBUTES_OWNERSHIP_BITMASK
			if optargs.Ownership {
				c_optargs.ownership = 1
			} else {
				c_optargs.ownership = 0
			}
		}
	}

	r := C.guestfs_copy_attributes_argv(g.g, c_src, c_dest, &c_optargs)

	if r == -1 {
		return get_error_from_handle(g, "copy_attributes")
	}
	return nil
}

/* Struct carrying optional arguments for Copy_device_to_device */
type OptargsCopy_device_to_device struct {
	/* Srcoffset field is ignored unless Srcoffset_is_set == true */
	Srcoffset_is_set bool
	Srcoffset        int64
	/* Destoffset field is ignored unless Destoffset_is_set == true */
	Destoffset_is_set bool
	Destoffset        int64
	/* Size field is ignored unless Size_is_set == true */
	Size_is_set bool
	Size        int64
	/* Sparse field is ignored unless Sparse_is_set == true */
	Sparse_is_set bool
	Sparse        bool
	/* Append field is ignored unless Append_is_set == true */
	Append_is_set bool
	Append        bool
}

/* copy_device_to_device : copy from source device to destination device */
func (g *Guestfs) Copy_device_to_device(src string, dest string, optargs *OptargsCopy_device_to_device) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("copy_device_to_device")
	}

	c_src := C.CString(src)
	defer C.free(unsafe.Pointer(c_src))

	c_dest := C.CString(dest)
	defer C.free(unsafe.Pointer(c_dest))
	c_optargs := C.struct_guestfs_copy_device_to_device_argv{}
	if optargs != nil {
		if optargs.Srcoffset_is_set {
			c_optargs.bitmask |= C.GUESTFS_COPY_DEVICE_TO_DEVICE_SRCOFFSET_BITMASK
			c_optargs.srcoffset = C.int64_t(optargs.Srcoffset)
		}
		if optargs.Destoffset_is_set {
			c_optargs.bitmask |= C.GUESTFS_COPY_DEVICE_TO_DEVICE_DESTOFFSET_BITMASK
			c_optargs.destoffset = C.int64_t(optargs.Destoffset)
		}
		if optargs.Size_is_set {
			c_optargs.bitmask |= C.GUESTFS_COPY_DEVICE_TO_DEVICE_SIZE_BITMASK
			c_optargs.size = C.int64_t(optargs.Size)
		}
		if optargs.Sparse_is_set {
			c_optargs.bitmask |= C.GUESTFS_COPY_DEVICE_TO_DEVICE_SPARSE_BITMASK
			if optargs.Sparse {
				c_optargs.sparse = 1
			} else {
				c_optargs.sparse = 0
			}
		}
		if optargs.Append_is_set {
			c_optargs.bitmask |= C.GUESTFS_COPY_DEVICE_TO_DEVICE_APPEND_BITMASK
			if optargs.Append {
				c_optargs.append = 1
			} else {
				c_optargs.append = 0
			}
		}
	}

	r := C.guestfs_copy_device_to_device_argv(g.g, c_src, c_dest, &c_optargs)

	if r == -1 {
		return get_error_from_handle(g, "copy_device_to_device")
	}
	return nil
}

/* Struct carrying optional arguments for Copy_device_to_file */
type OptargsCopy_device_to_file struct {
	/* Srcoffset field is ignored unless Srcoffset_is_set == true */
	Srcoffset_is_set bool
	Srcoffset        int64
	/* Destoffset field is ignored unless Destoffset_is_set == true */
	Destoffset_is_set bool
	Destoffset        int64
	/* Size field is ignored unless Size_is_set == true */
	Size_is_set bool
	Size        int64
	/* Sparse field is ignored unless Sparse_is_set == true */
	Sparse_is_set bool
	Sparse        bool
	/* Append field is ignored unless Append_is_set == true */
	Append_is_set bool
	Append        bool
}

/* copy_device_to_file : copy from source device to destination file */
func (g *Guestfs) Copy_device_to_file(src string, dest string, optargs *OptargsCopy_device_to_file) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("copy_device_to_file")
	}

	c_src := C.CString(src)
	defer C.free(unsafe.Pointer(c_src))

	c_dest := C.CString(dest)
	defer C.free(unsafe.Pointer(c_dest))
	c_optargs := C.struct_guestfs_copy_device_to_file_argv{}
	if optargs != nil {
		if optargs.Srcoffset_is_set {
			c_optargs.bitmask |= C.GUESTFS_COPY_DEVICE_TO_FILE_SRCOFFSET_BITMASK
			c_optargs.srcoffset = C.int64_t(optargs.Srcoffset)
		}
		if optargs.Destoffset_is_set {
			c_optargs.bitmask |= C.GUESTFS_COPY_DEVICE_TO_FILE_DESTOFFSET_BITMASK
			c_optargs.destoffset = C.int64_t(optargs.Destoffset)
		}
		if optargs.Size_is_set {
			c_optargs.bitmask |= C.GUESTFS_COPY_DEVICE_TO_FILE_SIZE_BITMASK
			c_optargs.size = C.int64_t(optargs.Size)
		}
		if optargs.Sparse_is_set {
			c_optargs.bitmask |= C.GUESTFS_COPY_DEVICE_TO_FILE_SPARSE_BITMASK
			if optargs.Sparse {
				c_optargs.sparse = 1
			} else {
				c_optargs.sparse = 0
			}
		}
		if optargs.Append_is_set {
			c_optargs.bitmask |= C.GUESTFS_COPY_DEVICE_TO_FILE_APPEND_BITMASK
			if optargs.Append {
				c_optargs.append = 1
			} else {
				c_optargs.append = 0
			}
		}
	}

	r := C.guestfs_copy_device_to_file_argv(g.g, c_src, c_dest, &c_optargs)

	if r == -1 {
		return get_error_from_handle(g, "copy_device_to_file")
	}
	return nil
}

/* Struct carrying optional arguments for Copy_file_to_device */
type OptargsCopy_file_to_device struct {
	/* Srcoffset field is ignored unless Srcoffset_is_set == true */
	Srcoffset_is_set bool
	Srcoffset        int64
	/* Destoffset field is ignored unless Destoffset_is_set == true */
	Destoffset_is_set bool
	Destoffset        int64
	/* Size field is ignored unless Size_is_set == true */
	Size_is_set bool
	Size        int64
	/* Sparse field is ignored unless Sparse_is_set == true */
	Sparse_is_set bool
	Sparse        bool
	/* Append field is ignored unless Append_is_set == true */
	Append_is_set bool
	Append        bool
}

/* copy_file_to_device : copy from source file to destination device */
func (g *Guestfs) Copy_file_to_device(src string, dest string, optargs *OptargsCopy_file_to_device) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("copy_file_to_device")
	}

	c_src := C.CString(src)
	defer C.free(unsafe.Pointer(c_src))

	c_dest := C.CString(dest)
	defer C.free(unsafe.Pointer(c_dest))
	c_optargs := C.struct_guestfs_copy_file_to_device_argv{}
	if optargs != nil {
		if optargs.Srcoffset_is_set {
			c_optargs.bitmask |= C.GUESTFS_COPY_FILE_TO_DEVICE_SRCOFFSET_BITMASK
			c_optargs.srcoffset = C.int64_t(optargs.Srcoffset)
		}
		if optargs.Destoffset_is_set {
			c_optargs.bitmask |= C.GUESTFS_COPY_FILE_TO_DEVICE_DESTOFFSET_BITMASK
			c_optargs.destoffset = C.int64_t(optargs.Destoffset)
		}
		if optargs.Size_is_set {
			c_optargs.bitmask |= C.GUESTFS_COPY_FILE_TO_DEVICE_SIZE_BITMASK
			c_optargs.size = C.int64_t(optargs.Size)
		}
		if optargs.Sparse_is_set {
			c_optargs.bitmask |= C.GUESTFS_COPY_FILE_TO_DEVICE_SPARSE_BITMASK
			if optargs.Sparse {
				c_optargs.sparse = 1
			} else {
				c_optargs.sparse = 0
			}
		}
		if optargs.Append_is_set {
			c_optargs.bitmask |= C.GUESTFS_COPY_FILE_TO_DEVICE_APPEND_BITMASK
			if optargs.Append {
				c_optargs.append = 1
			} else {
				c_optargs.append = 0
			}
		}
	}

	r := C.guestfs_copy_file_to_device_argv(g.g, c_src, c_dest, &c_optargs)

	if r == -1 {
		return get_error_from_handle(g, "copy_file_to_device")
	}
	return nil
}

/* Struct carrying optional arguments for Copy_file_to_file */
type OptargsCopy_file_to_file struct {
	/* Srcoffset field is ignored unless Srcoffset_is_set == true */
	Srcoffset_is_set bool
	Srcoffset        int64
	/* Destoffset field is ignored unless Destoffset_is_set == true */
	Destoffset_is_set bool
	Destoffset        int64
	/* Size field is ignored unless Size_is_set == true */
	Size_is_set bool
	Size        int64
	/* Sparse field is ignored unless Sparse_is_set == true */
	Sparse_is_set bool
	Sparse        bool
	/* Append field is ignored unless Append_is_set == true */
	Append_is_set bool
	Append        bool
}

/* copy_file_to_file : copy from source file to destination file */
func (g *Guestfs) Copy_file_to_file(src string, dest string, optargs *OptargsCopy_file_to_file) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("copy_file_to_file")
	}

	c_src := C.CString(src)
	defer C.free(unsafe.Pointer(c_src))

	c_dest := C.CString(dest)
	defer C.free(unsafe.Pointer(c_dest))
	c_optargs := C.struct_guestfs_copy_file_to_file_argv{}
	if optargs != nil {
		if optargs.Srcoffset_is_set {
			c_optargs.bitmask |= C.GUESTFS_COPY_FILE_TO_FILE_SRCOFFSET_BITMASK
			c_optargs.srcoffset = C.int64_t(optargs.Srcoffset)
		}
		if optargs.Destoffset_is_set {
			c_optargs.bitmask |= C.GUESTFS_COPY_FILE_TO_FILE_DESTOFFSET_BITMASK
			c_optargs.destoffset = C.int64_t(optargs.Destoffset)
		}
		if optargs.Size_is_set {
			c_optargs.bitmask |= C.GUESTFS_COPY_FILE_TO_FILE_SIZE_BITMASK
			c_optargs.size = C.int64_t(optargs.Size)
		}
		if optargs.Sparse_is_set {
			c_optargs.bitmask |= C.GUESTFS_COPY_FILE_TO_FILE_SPARSE_BITMASK
			if optargs.Sparse {
				c_optargs.sparse = 1
			} else {
				c_optargs.sparse = 0
			}
		}
		if optargs.Append_is_set {
			c_optargs.bitmask |= C.GUESTFS_COPY_FILE_TO_FILE_APPEND_BITMASK
			if optargs.Append {
				c_optargs.append = 1
			} else {
				c_optargs.append = 0
			}
		}
	}

	r := C.guestfs_copy_file_to_file_argv(g.g, c_src, c_dest, &c_optargs)

	if r == -1 {
		return get_error_from_handle(g, "copy_file_to_file")
	}
	return nil
}

/* copy_in : copy local files or directories into an image */
func (g *Guestfs) Copy_in(localpath string, remotedir string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("copy_in")
	}

	c_localpath := C.CString(localpath)
	defer C.free(unsafe.Pointer(c_localpath))

	c_remotedir := C.CString(remotedir)
	defer C.free(unsafe.Pointer(c_remotedir))

	r := C.guestfs_copy_in(g.g, c_localpath, c_remotedir)

	if r == -1 {
		return get_error_from_handle(g, "copy_in")
	}
	return nil
}

/* copy_out : copy remote files or directories out of an image */
func (g *Guestfs) Copy_out(remotepath string, localdir string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("copy_out")
	}

	c_remotepath := C.CString(remotepath)
	defer C.free(unsafe.Pointer(c_remotepath))

	c_localdir := C.CString(localdir)
	defer C.free(unsafe.Pointer(c_localdir))

	r := C.guestfs_copy_out(g.g, c_remotepath, c_localdir)

	if r == -1 {
		return get_error_from_handle(g, "copy_out")
	}
	return nil
}

/* copy_size : copy size bytes from source to destination using dd */
func (g *Guestfs) Copy_size(src string, dest string, size int64) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("copy_size")
	}

	c_src := C.CString(src)
	defer C.free(unsafe.Pointer(c_src))

	c_dest := C.CString(dest)
	defer C.free(unsafe.Pointer(c_dest))

	r := C.guestfs_copy_size(g.g, c_src, c_dest, C.int64_t(size))

	if r == -1 {
		return get_error_from_handle(g, "copy_size")
	}
	return nil
}

/* cp : copy a file */
func (g *Guestfs) Cp(src string, dest string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("cp")
	}

	c_src := C.CString(src)
	defer C.free(unsafe.Pointer(c_src))

	c_dest := C.CString(dest)
	defer C.free(unsafe.Pointer(c_dest))

	r := C.guestfs_cp(g.g, c_src, c_dest)

	if r == -1 {
		return get_error_from_handle(g, "cp")
	}
	return nil
}

/* cp_a : copy a file or directory recursively */
func (g *Guestfs) Cp_a(src string, dest string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("cp_a")
	}

	c_src := C.CString(src)
	defer C.free(unsafe.Pointer(c_src))

	c_dest := C.CString(dest)
	defer C.free(unsafe.Pointer(c_dest))

	r := C.guestfs_cp_a(g.g, c_src, c_dest)

	if r == -1 {
		return get_error_from_handle(g, "cp_a")
	}
	return nil
}

/* cp_r : copy a file or directory recursively */
func (g *Guestfs) Cp_r(src string, dest string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("cp_r")
	}

	c_src := C.CString(src)
	defer C.free(unsafe.Pointer(c_src))

	c_dest := C.CString(dest)
	defer C.free(unsafe.Pointer(c_dest))

	r := C.guestfs_cp_r(g.g, c_src, c_dest)

	if r == -1 {
		return get_error_from_handle(g, "cp_r")
	}
	return nil
}

/* Struct carrying optional arguments for Cpio_out */
type OptargsCpio_out struct {
	/* Format field is ignored unless Format_is_set == true */
	Format_is_set bool
	Format        string
}

/* cpio_out : pack directory into cpio file */
func (g *Guestfs) Cpio_out(directory string, cpiofile string, optargs *OptargsCpio_out) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("cpio_out")
	}

	c_directory := C.CString(directory)
	defer C.free(unsafe.Pointer(c_directory))

	c_cpiofile := C.CString(cpiofile)
	defer C.free(unsafe.Pointer(c_cpiofile))
	c_optargs := C.struct_guestfs_cpio_out_argv{}
	if optargs != nil {
		if optargs.Format_is_set {
			c_optargs.bitmask |= C.GUESTFS_CPIO_OUT_FORMAT_BITMASK
			c_optargs.format = C.CString(optargs.Format)
			defer C.free(unsafe.Pointer(c_optargs.format))
		}
	}

	r := C.guestfs_cpio_out_argv(g.g, c_directory, c_cpiofile, &c_optargs)

	if r == -1 {
		return get_error_from_handle(g, "cpio_out")
	}
	return nil
}

/* dd : copy from source to destination using dd */
func (g *Guestfs) Dd(src string, dest string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("dd")
	}

	c_src := C.CString(src)
	defer C.free(unsafe.Pointer(c_src))

	c_dest := C.CString(dest)
	defer C.free(unsafe.Pointer(c_dest))

	r := C.guestfs_dd(g.g, c_src, c_dest)

	if r == -1 {
		return get_error_from_handle(g, "dd")
	}
	return nil
}

/* debug : debugging and internals */
func (g *Guestfs) Debug(subcmd string, extraargs []string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("debug")
	}

	c_subcmd := C.CString(subcmd)
	defer C.free(unsafe.Pointer(c_subcmd))

	c_extraargs := arg_string_list(extraargs)
	defer free_string_list(c_extraargs)

	r := C.guestfs_debug(g.g, c_subcmd, c_extraargs)

	if r == nil {
		return "", get_error_from_handle(g, "debug")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* debug_drives : debug the drives (internal use only) */
func (g *Guestfs) Debug_drives() ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("debug_drives")
	}

	r := C.guestfs_debug_drives(g.g)

	if r == nil {
		return nil, get_error_from_handle(g, "debug_drives")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* debug_upload : upload a file to the appliance (internal use only) */
func (g *Guestfs) Debug_upload(filename string, tmpname string, mode int) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("debug_upload")
	}

	c_filename := C.CString(filename)
	defer C.free(unsafe.Pointer(c_filename))

	c_tmpname := C.CString(tmpname)
	defer C.free(unsafe.Pointer(c_tmpname))

	r := C.guestfs_debug_upload(g.g, c_filename, c_tmpname, C.int(mode))

	if r == -1 {
		return get_error_from_handle(g, "debug_upload")
	}
	return nil
}

/* device_index : convert device to index */
func (g *Guestfs) Device_index(device string) (int, *GuestfsError) {
	if g.g == nil {
		return 0, closed_handle_error("device_index")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_device_index(g.g, c_device)

	if r == -1 {
		return 0, get_error_from_handle(g, "device_index")
	}
	return int(r), nil
}

/* df : report file system disk space usage */
func (g *Guestfs) Df() (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("df")
	}

	r := C.guestfs_df(g.g)

	if r == nil {
		return "", get_error_from_handle(g, "df")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* df_h : report file system disk space usage (human readable) */
func (g *Guestfs) Df_h() (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("df_h")
	}

	r := C.guestfs_df_h(g.g)

	if r == nil {
		return "", get_error_from_handle(g, "df_h")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* Struct carrying optional arguments for Disk_create */
type OptargsDisk_create struct {
	/* Backingfile field is ignored unless Backingfile_is_set == true */
	Backingfile_is_set bool
	Backingfile        string
	/* Backingformat field is ignored unless Backingformat_is_set == true */
	Backingformat_is_set bool
	Backingformat        string
	/* Preallocation field is ignored unless Preallocation_is_set == true */
	Preallocation_is_set bool
	Preallocation        string
	/* Compat field is ignored unless Compat_is_set == true */
	Compat_is_set bool
	Compat        string
	/* Clustersize field is ignored unless Clustersize_is_set == true */
	Clustersize_is_set bool
	Clustersize        int
}

/* disk_create : create a blank disk image */
func (g *Guestfs) Disk_create(filename string, format string, size int64, optargs *OptargsDisk_create) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("disk_create")
	}

	c_filename := C.CString(filename)
	defer C.free(unsafe.Pointer(c_filename))

	c_format := C.CString(format)
	defer C.free(unsafe.Pointer(c_format))
	c_optargs := C.struct_guestfs_disk_create_argv{}
	if optargs != nil {
		if optargs.Backingfile_is_set {
			c_optargs.bitmask |= C.GUESTFS_DISK_CREATE_BACKINGFILE_BITMASK
			c_optargs.backingfile = C.CString(optargs.Backingfile)
			defer C.free(unsafe.Pointer(c_optargs.backingfile))
		}
		if optargs.Backingformat_is_set {
			c_optargs.bitmask |= C.GUESTFS_DISK_CREATE_BACKINGFORMAT_BITMASK
			c_optargs.backingformat = C.CString(optargs.Backingformat)
			defer C.free(unsafe.Pointer(c_optargs.backingformat))
		}
		if optargs.Preallocation_is_set {
			c_optargs.bitmask |= C.GUESTFS_DISK_CREATE_PREALLOCATION_BITMASK
			c_optargs.preallocation = C.CString(optargs.Preallocation)
			defer C.free(unsafe.Pointer(c_optargs.preallocation))
		}
		if optargs.Compat_is_set {
			c_optargs.bitmask |= C.GUESTFS_DISK_CREATE_COMPAT_BITMASK
			c_optargs.compat = C.CString(optargs.Compat)
			defer C.free(unsafe.Pointer(c_optargs.compat))
		}
		if optargs.Clustersize_is_set {
			c_optargs.bitmask |= C.GUESTFS_DISK_CREATE_CLUSTERSIZE_BITMASK
			c_optargs.clustersize = C.int(optargs.Clustersize)
		}
	}

	r := C.guestfs_disk_create_argv(g.g, c_filename, c_format, C.int64_t(size), &c_optargs)

	if r == -1 {
		return get_error_from_handle(g, "disk_create")
	}
	return nil
}

/* disk_format : detect the disk format of a disk image */
func (g *Guestfs) Disk_format(filename string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("disk_format")
	}

	c_filename := C.CString(filename)
	defer C.free(unsafe.Pointer(c_filename))

	r := C.guestfs_disk_format(g.g, c_filename)

	if r == nil {
		return "", get_error_from_handle(g, "disk_format")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* disk_has_backing_file : return whether disk has a backing file */
func (g *Guestfs) Disk_has_backing_file(filename string) (bool, *GuestfsError) {
	if g.g == nil {
		return false, closed_handle_error("disk_has_backing_file")
	}

	c_filename := C.CString(filename)
	defer C.free(unsafe.Pointer(c_filename))

	r := C.guestfs_disk_has_backing_file(g.g, c_filename)

	if r == -1 {
		return false, get_error_from_handle(g, "disk_has_backing_file")
	}
	return r != 0, nil
}

/* disk_virtual_size : return virtual size of a disk */
func (g *Guestfs) Disk_virtual_size(filename string) (int64, *GuestfsError) {
	if g.g == nil {
		return 0, closed_handle_error("disk_virtual_size")
	}

	c_filename := C.CString(filename)
	defer C.free(unsafe.Pointer(c_filename))

	r := C.guestfs_disk_virtual_size(g.g, c_filename)

	if r == -1 {
		return 0, get_error_from_handle(g, "disk_virtual_size")
	}
	return int64(r), nil
}

/* dmesg : return kernel messages */
func (g *Guestfs) Dmesg() (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("dmesg")
	}

	r := C.guestfs_dmesg(g.g)

	if r == nil {
		return "", get_error_from_handle(g, "dmesg")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* download : download a file to the local machine */
func (g *Guestfs) Download(remotefilename string, filename string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("download")
	}

	c_remotefilename := C.CString(remotefilename)
	defer C.free(unsafe.Pointer(c_remotefilename))

	c_filename := C.CString(filename)
	defer C.free(unsafe.Pointer(c_filename))

	r := C.guestfs_download(g.g, c_remotefilename, c_filename)

	if r == -1 {
		return get_error_from_handle(g, "download")
	}
	return nil
}

/* download_offset : download a file to the local machine with offset and size */
func (g *Guestfs) Download_offset(remotefilename string, filename string, offset int64, size int64) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("download_offset")
	}

	c_remotefilename := C.CString(remotefilename)
	defer C.free(unsafe.Pointer(c_remotefilename))

	c_filename := C.CString(filename)
	defer C.free(unsafe.Pointer(c_filename))

	r := C.guestfs_download_offset(g.g, c_remotefilename, c_filename, C.int64_t(offset), C.int64_t(size))

	if r == -1 {
		return get_error_from_handle(g, "download_offset")
	}
	return nil
}

/* drop_caches : drop kernel page cache, dentries and inodes */
func (g *Guestfs) Drop_caches(whattodrop int) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("drop_caches")
	}

	r := C.guestfs_drop_caches(g.g, C.int(whattodrop))

	if r == -1 {
		return get_error_from_handle(g, "drop_caches")
	}
	return nil
}

/* du : estimate file space usage */
func (g *Guestfs) Du(path string) (int64, *GuestfsError) {
	if g.g == nil {
		return 0, closed_handle_error("du")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_du(g.g, c_path)

	if r == -1 {
		return 0, get_error_from_handle(g, "du")
	}
	return int64(r), nil
}

/* Struct carrying optional arguments for E2fsck */
type OptargsE2fsck struct {
	/* Correct field is ignored unless Correct_is_set == true */
	Correct_is_set bool
	Correct        bool
	/* Forceall field is ignored unless Forceall_is_set == true */
	Forceall_is_set bool
	Forceall        bool
}

/* e2fsck : check an ext2/ext3 filesystem */
func (g *Guestfs) E2fsck(device string, optargs *OptargsE2fsck) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("e2fsck")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))
	c_optargs := C.struct_guestfs_e2fsck_argv{}
	if optargs != nil {
		if optargs.Correct_is_set {
			c_optargs.bitmask |= C.GUESTFS_E2FSCK_CORRECT_BITMASK
			if optargs.Correct {
				c_optargs.correct = 1
			} else {
				c_optargs.correct = 0
			}
		}
		if optargs.Forceall_is_set {
			c_optargs.bitmask |= C.GUESTFS_E2FSCK_FORCEALL_BITMASK
			if optargs.Forceall {
				c_optargs.forceall = 1
			} else {
				c_optargs.forceall = 0
			}
		}
	}

	r := C.guestfs_e2fsck_argv(g.g, c_device, &c_optargs)

	if r == -1 {
		return get_error_from_handle(g, "e2fsck")
	}
	return nil
}

/* e2fsck_f : check an ext2/ext3 filesystem */
func (g *Guestfs) E2fsck_f(device string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("e2fsck_f")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_e2fsck_f(g.g, c_device)

	if r == -1 {
		return get_error_from_handle(g, "e2fsck_f")
	}
	return nil
}

/* echo_daemon : echo arguments back to the client */
func (g *Guestfs) Echo_daemon(words []string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("echo_daemon")
	}

	c_words := arg_string_list(words)
	defer free_string_list(c_words)

	r := C.guestfs_echo_daemon(g.g, c_words)

	if r == nil {
		return "", get_error_from_handle(g, "echo_daemon")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* egrep : return lines matching a pattern */
func (g *Guestfs) Egrep(regex string, path string) ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("egrep")
	}

	c_regex := C.CString(regex)
	defer C.free(unsafe.Pointer(c_regex))

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_egrep(g.g, c_regex, c_path)

	if r == nil {
		return nil, get_error_from_handle(g, "egrep")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* egrepi : return lines matching a pattern */
func (g *Guestfs) Egrepi(regex string, path string) ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("egrepi")
	}

	c_regex := C.CString(regex)
	defer C.free(unsafe.Pointer(c_regex))

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_egrepi(g.g, c_regex, c_path)

	if r == nil {
		return nil, get_error_from_handle(g, "egrepi")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* equal : test if two files have equal contents */
func (g *Guestfs) Equal(file1 string, file2 string) (bool, *GuestfsError) {
	if g.g == nil {
		return false, closed_handle_error("equal")
	}

	c_file1 := C.CString(file1)
	defer C.free(unsafe.Pointer(c_file1))

	c_file2 := C.CString(file2)
	defer C.free(unsafe.Pointer(c_file2))

	r := C.guestfs_equal(g.g, c_file1, c_file2)

	if r == -1 {
		return false, get_error_from_handle(g, "equal")
	}
	return r != 0, nil
}

/* exists : test if file or directory exists */
func (g *Guestfs) Exists(path string) (bool, *GuestfsError) {
	if g.g == nil {
		return false, closed_handle_error("exists")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_exists(g.g, c_path)

	if r == -1 {
		return false, get_error_from_handle(g, "exists")
	}
	return r != 0, nil
}

/* extlinux : install the SYSLINUX bootloader on an ext2/3/4 or btrfs filesystem */
func (g *Guestfs) Extlinux(directory string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("extlinux")
	}

	c_directory := C.CString(directory)
	defer C.free(unsafe.Pointer(c_directory))

	r := C.guestfs_extlinux(g.g, c_directory)

	if r == -1 {
		return get_error_from_handle(g, "extlinux")
	}
	return nil
}

/* fallocate : preallocate a file in the guest filesystem */
func (g *Guestfs) Fallocate(path string, len int) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("fallocate")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_fallocate(g.g, c_path, C.int(len))

	if r == -1 {
		return get_error_from_handle(g, "fallocate")
	}
	return nil
}

/* fallocate64 : preallocate a file in the guest filesystem */
func (g *Guestfs) Fallocate64(path string, len int64) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("fallocate64")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_fallocate64(g.g, c_path, C.int64_t(len))

	if r == -1 {
		return get_error_from_handle(g, "fallocate64")
	}
	return nil
}

/* feature_available : test availability of some parts of the API */
func (g *Guestfs) Feature_available(groups []string) (bool, *GuestfsError) {
	if g.g == nil {
		return false, closed_handle_error("feature_available")
	}

	c_groups := arg_string_list(groups)
	defer free_string_list(c_groups)

	r := C.guestfs_feature_available(g.g, c_groups)

	if r == -1 {
		return false, get_error_from_handle(g, "feature_available")
	}
	return r != 0, nil
}

/* fgrep : return lines matching a pattern */
func (g *Guestfs) Fgrep(pattern string, path string) ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("fgrep")
	}

	c_pattern := C.CString(pattern)
	defer C.free(unsafe.Pointer(c_pattern))

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_fgrep(g.g, c_pattern, c_path)

	if r == nil {
		return nil, get_error_from_handle(g, "fgrep")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* fgrepi : return lines matching a pattern */
func (g *Guestfs) Fgrepi(pattern string, path string) ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("fgrepi")
	}

	c_pattern := C.CString(pattern)
	defer C.free(unsafe.Pointer(c_pattern))

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_fgrepi(g.g, c_pattern, c_path)

	if r == nil {
		return nil, get_error_from_handle(g, "fgrepi")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* file : determine file type */
func (g *Guestfs) File(path string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("file")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_file(g.g, c_path)

	if r == nil {
		return "", get_error_from_handle(g, "file")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* file_architecture : detect the architecture of a binary file */
func (g *Guestfs) File_architecture(filename string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("file_architecture")
	}

	c_filename := C.CString(filename)
	defer C.free(unsafe.Pointer(c_filename))

	r := C.guestfs_file_architecture(g.g, c_filename)

	if r == nil {
		return "", get_error_from_handle(g, "file_architecture")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* filesize : return the size of the file in bytes */
func (g *Guestfs) Filesize(file string) (int64, *GuestfsError) {
	if g.g == nil {
		return 0, closed_handle_error("filesize")
	}

	c_file := C.CString(file)
	defer C.free(unsafe.Pointer(c_file))

	r := C.guestfs_filesize(g.g, c_file)

	if r == -1 {
		return 0, get_error_from_handle(g, "filesize")
	}
	return int64(r), nil
}

/* filesystem_available : check if filesystem is available */
func (g *Guestfs) Filesystem_available(filesystem string) (bool, *GuestfsError) {
	if g.g == nil {
		return false, closed_handle_error("filesystem_available")
	}

	c_filesystem := C.CString(filesystem)
	defer C.free(unsafe.Pointer(c_filesystem))

	r := C.guestfs_filesystem_available(g.g, c_filesystem)

	if r == -1 {
		return false, get_error_from_handle(g, "filesystem_available")
	}
	return r != 0, nil
}

/* fill : fill a file with octets */
func (g *Guestfs) Fill(c int, len int, path string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("fill")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_fill(g.g, C.int(c), C.int(len), c_path)

	if r == -1 {
		return get_error_from_handle(g, "fill")
	}
	return nil
}

/* fill_dir : fill a directory with empty files */
func (g *Guestfs) Fill_dir(dir string, nr int) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("fill_dir")
	}

	c_dir := C.CString(dir)
	defer C.free(unsafe.Pointer(c_dir))

	r := C.guestfs_fill_dir(g.g, c_dir, C.int(nr))

	if r == -1 {
		return get_error_from_handle(g, "fill_dir")
	}
	return nil
}

/* fill_pattern : fill a file with a repeating pattern of bytes */
func (g *Guestfs) Fill_pattern(pattern string, len int, path string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("fill_pattern")
	}

	c_pattern := C.CString(pattern)
	defer C.free(unsafe.Pointer(c_pattern))

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_fill_pattern(g.g, c_pattern, C.int(len), c_path)

	if r == -1 {
		return get_error_from_handle(g, "fill_pattern")
	}
	return nil
}

/* find : find all files and directories */
func (g *Guestfs) Find(directory string) ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("find")
	}

	c_directory := C.CString(directory)
	defer C.free(unsafe.Pointer(c_directory))

	r := C.guestfs_find(g.g, c_directory)

	if r == nil {
		return nil, get_error_from_handle(g, "find")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* find0 : find all files and directories, returning NUL-separated list */
func (g *Guestfs) Find0(directory string, files string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("find0")
	}

	c_directory := C.CString(directory)
	defer C.free(unsafe.Pointer(c_directory))

	c_files := C.CString(files)
	defer C.free(unsafe.Pointer(c_files))

	r := C.guestfs_find0(g.g, c_directory, c_files)

	if r == -1 {
		return get_error_from_handle(g, "find0")
	}
	return nil
}

/* findfs_label : find a filesystem by label */
func (g *Guestfs) Findfs_label(label string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("findfs_label")
	}

	c_label := C.CString(label)
	defer C.free(unsafe.Pointer(c_label))

	r := C.guestfs_findfs_label(g.g, c_label)

	if r == nil {
		return "", get_error_from_handle(g, "findfs_label")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* findfs_uuid : find a filesystem by UUID */
func (g *Guestfs) Findfs_uuid(uuid string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("findfs_uuid")
	}

	c_uuid := C.CString(uuid)
	defer C.free(unsafe.Pointer(c_uuid))

	r := C.guestfs_findfs_uuid(g.g, c_uuid)

	if r == nil {
		return "", get_error_from_handle(g, "findfs_uuid")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* fsck : run the filesystem checker */
func (g *Guestfs) Fsck(fstype string, device string) (int, *GuestfsError) {
	if g.g == nil {
		return 0, closed_handle_error("fsck")
	}

	c_fstype := C.CString(fstype)
	defer C.free(unsafe.Pointer(c_fstype))

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_fsck(g.g, c_fstype, c_device)

	if r == -1 {
		return 0, get_error_from_handle(g, "fsck")
	}
	return int(r), nil
}

/* Struct carrying optional arguments for Fstrim */
type OptargsFstrim struct {
	/* Offset field is ignored unless Offset_is_set == true */
	Offset_is_set bool
	Offset        int64
	/* Length field is ignored unless Length_is_set == true */
	Length_is_set bool
	Length        int64
	/* Minimumfreeextent field is ignored unless Minimumfreeextent_is_set == true */
	Minimumfreeextent_is_set bool
	Minimumfreeextent        int64
}

/* fstrim : trim free space in a filesystem */
func (g *Guestfs) Fstrim(mountpoint string, optargs *OptargsFstrim) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("fstrim")
	}

	c_mountpoint := C.CString(mountpoint)
	defer C.free(unsafe.Pointer(c_mountpoint))
	c_optargs := C.struct_guestfs_fstrim_argv{}
	if optargs != nil {
		if optargs.Offset_is_set {
			c_optargs.bitmask |= C.GUESTFS_FSTRIM_OFFSET_BITMASK
			c_optargs.offset = C.int64_t(optargs.Offset)
		}
		if optargs.Length_is_set {
			c_optargs.bitmask |= C.GUESTFS_FSTRIM_LENGTH_BITMASK
			c_optargs.length = C.int64_t(optargs.Length)
		}
		if optargs.Minimumfreeextent_is_set {
			c_optargs.bitmask |= C.GUESTFS_FSTRIM_MINIMUMFREEEXTENT_BITMASK
			c_optargs.minimumfreeextent = C.int64_t(optargs.Minimumfreeextent)
		}
	}

	r := C.guestfs_fstrim_argv(g.g, c_mountpoint, &c_optargs)

	if r == -1 {
		return get_error_from_handle(g, "fstrim")
	}
	return nil
}

/* get_append : get the additional kernel options */
func (g *Guestfs) Get_append() (*string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("get_append")
	}

	r := C.guestfs_get_append(g.g)
	if r != nil {
		r_s := string(*r)
		return &r_s, nil
	} else {
		return nil, nil
	}
}

/* get_attach_method : get the backend */
func (g *Guestfs) Get_attach_method() (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("get_attach_method")
	}

	r := C.guestfs_get_attach_method(g.g)

	if r == nil {
		return "", get_error_from_handle(g, "get_attach_method")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* get_autosync : get autosync mode */
func (g *Guestfs) Get_autosync() (bool, *GuestfsError) {
	if g.g == nil {
		return false, closed_handle_error("get_autosync")
	}

	r := C.guestfs_get_autosync(g.g)

	if r == -1 {
		return false, get_error_from_handle(g, "get_autosync")
	}
	return r != 0, nil
}

/* get_backend : get the backend */
func (g *Guestfs) Get_backend() (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("get_backend")
	}

	r := C.guestfs_get_backend(g.g)

	if r == nil {
		return "", get_error_from_handle(g, "get_backend")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* get_backend_setting : get a single per-backend settings string */
func (g *Guestfs) Get_backend_setting(name string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("get_backend_setting")
	}

	c_name := C.CString(name)
	defer C.free(unsafe.Pointer(c_name))

	r := C.guestfs_get_backend_setting(g.g, c_name)

	if r == nil {
		return "", get_error_from_handle(g, "get_backend_setting")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* get_backend_settings : get per-backend settings */
func (g *Guestfs) Get_backend_settings() ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("get_backend_settings")
	}

	r := C.guestfs_get_backend_settings(g.g)

	if r == nil {
		return nil, get_error_from_handle(g, "get_backend_settings")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* get_cachedir : get the appliance cache directory */
func (g *Guestfs) Get_cachedir() (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("get_cachedir")
	}

	r := C.guestfs_get_cachedir(g.g)

	if r == nil {
		return "", get_error_from_handle(g, "get_cachedir")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* get_direct : get direct appliance mode flag */
func (g *Guestfs) Get_direct() (bool, *GuestfsError) {
	if g.g == nil {
		return false, closed_handle_error("get_direct")
	}

	r := C.guestfs_get_direct(g.g)

	if r == -1 {
		return false, get_error_from_handle(g, "get_direct")
	}
	return r != 0, nil
}

/* get_e2attrs : get ext2 file attributes of a file */
func (g *Guestfs) Get_e2attrs(file string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("get_e2attrs")
	}

	c_file := C.CString(file)
	defer C.free(unsafe.Pointer(c_file))

	r := C.guestfs_get_e2attrs(g.g, c_file)

	if r == nil {
		return "", get_error_from_handle(g, "get_e2attrs")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* get_e2generation : get ext2 file generation of a file */
func (g *Guestfs) Get_e2generation(file string) (int64, *GuestfsError) {
	if g.g == nil {
		return 0, closed_handle_error("get_e2generation")
	}

	c_file := C.CString(file)
	defer C.free(unsafe.Pointer(c_file))

	r := C.guestfs_get_e2generation(g.g, c_file)

	if r == -1 {
		return 0, get_error_from_handle(g, "get_e2generation")
	}
	return int64(r), nil
}

/* get_e2label : get the ext2/3/4 filesystem label */
func (g *Guestfs) Get_e2label(device string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("get_e2label")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_get_e2label(g.g, c_device)

	if r == nil {
		return "", get_error_from_handle(g, "get_e2label")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* get_e2uuid : get the ext2/3/4 filesystem UUID */
func (g *Guestfs) Get_e2uuid(device string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("get_e2uuid")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_get_e2uuid(g.g, c_device)

	if r == nil {
		return "", get_error_from_handle(g, "get_e2uuid")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* get_hv : get the hypervisor binary */
func (g *Guestfs) Get_hv() (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("get_hv")
	}

	r := C.guestfs_get_hv(g.g)

	if r == nil {
		return "", get_error_from_handle(g, "get_hv")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* get_libvirt_requested_credential_challenge : challenge of i'th requested credential */
func (g *Guestfs) Get_libvirt_requested_credential_challenge(index int) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("get_libvirt_requested_credential_challenge")
	}

	r := C.guestfs_get_libvirt_requested_credential_challenge(g.g, C.int(index))

	if r == nil {
		return "", get_error_from_handle(g, "get_libvirt_requested_credential_challenge")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* get_libvirt_requested_credential_defresult : default result of i'th requested credential */
func (g *Guestfs) Get_libvirt_requested_credential_defresult(index int) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("get_libvirt_requested_credential_defresult")
	}

	r := C.guestfs_get_libvirt_requested_credential_defresult(g.g, C.int(index))

	if r == nil {
		return "", get_error_from_handle(g, "get_libvirt_requested_credential_defresult")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* get_libvirt_requested_credential_prompt : prompt of i'th requested credential */
func (g *Guestfs) Get_libvirt_requested_credential_prompt(index int) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("get_libvirt_requested_credential_prompt")
	}

	r := C.guestfs_get_libvirt_requested_credential_prompt(g.g, C.int(index))

	if r == nil {
		return "", get_error_from_handle(g, "get_libvirt_requested_credential_prompt")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* get_libvirt_requested_credentials : get list of credentials requested by libvirt */
func (g *Guestfs) Get_libvirt_requested_credentials() ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("get_libvirt_requested_credentials")
	}

	r := C.guestfs_get_libvirt_requested_credentials(g.g)

	if r == nil {
		return nil, get_error_from_handle(g, "get_libvirt_requested_credentials")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* get_memsize : get memory allocated to the hypervisor */
func (g *Guestfs) Get_memsize() (int, *GuestfsError) {
	if g.g == nil {
		return 0, closed_handle_error("get_memsize")
	}

	r := C.guestfs_get_memsize(g.g)

	if r == -1 {
		return 0, get_error_from_handle(g, "get_memsize")
	}
	return int(r), nil
}

/* get_network : get enable network flag */
func (g *Guestfs) Get_network() (bool, *GuestfsError) {
	if g.g == nil {
		return false, closed_handle_error("get_network")
	}

	r := C.guestfs_get_network(g.g)

	if r == -1 {
		return false, get_error_from_handle(g, "get_network")
	}
	return r != 0, nil
}

/* get_path : get the search path */
func (g *Guestfs) Get_path() (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("get_path")
	}

	r := C.guestfs_get_path(g.g)

	if r == nil {
		return "", get_error_from_handle(g, "get_path")
	}
	return C.GoString(r), nil
}

/* get_pgroup : get process group flag */
func (g *Guestfs) Get_pgroup() (bool, *GuestfsError) {
	if g.g == nil {
		return false, closed_handle_error("get_pgroup")
	}

	r := C.guestfs_get_pgroup(g.g)

	if r == -1 {
		return false, get_error_from_handle(g, "get_pgroup")
	}
	return r != 0, nil
}

/* get_pid : get PID of hypervisor */
func (g *Guestfs) Get_pid() (int, *GuestfsError) {
	if g.g == nil {
		return 0, closed_handle_error("get_pid")
	}

	r := C.guestfs_get_pid(g.g)

	if r == -1 {
		return 0, get_error_from_handle(g, "get_pid")
	}
	return int(r), nil
}

/* get_program : get the program name */
func (g *Guestfs) Get_program() (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("get_program")
	}

	r := C.guestfs_get_program(g.g)

	if r == nil {
		return "", get_error_from_handle(g, "get_program")
	}
	return C.GoString(r), nil
}

/* get_qemu : get the hypervisor binary (usually qemu) */
func (g *Guestfs) Get_qemu() (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("get_qemu")
	}

	r := C.guestfs_get_qemu(g.g)

	if r == nil {
		return "", get_error_from_handle(g, "get_qemu")
	}
	return C.GoString(r), nil
}

/* get_recovery_proc : get recovery process enabled flag */
func (g *Guestfs) Get_recovery_proc() (bool, *GuestfsError) {
	if g.g == nil {
		return false, closed_handle_error("get_recovery_proc")
	}

	r := C.guestfs_get_recovery_proc(g.g)

	if r == -1 {
		return false, get_error_from_handle(g, "get_recovery_proc")
	}
	return r != 0, nil
}

/* get_selinux : get SELinux enabled flag */
func (g *Guestfs) Get_selinux() (bool, *GuestfsError) {
	if g.g == nil {
		return false, closed_handle_error("get_selinux")
	}

	r := C.guestfs_get_selinux(g.g)

	if r == -1 {
		return false, get_error_from_handle(g, "get_selinux")
	}
	return r != 0, nil
}

/* get_smp : get number of virtual CPUs in appliance */
func (g *Guestfs) Get_smp() (int, *GuestfsError) {
	if g.g == nil {
		return 0, closed_handle_error("get_smp")
	}

	r := C.guestfs_get_smp(g.g)

	if r == -1 {
		return 0, get_error_from_handle(g, "get_smp")
	}
	return int(r), nil
}

/* get_state : get the current state */
func (g *Guestfs) Get_state() (int, *GuestfsError) {
	if g.g == nil {
		return 0, closed_handle_error("get_state")
	}

	r := C.guestfs_get_state(g.g)

	if r == -1 {
		return 0, get_error_from_handle(g, "get_state")
	}
	return int(r), nil
}

/* get_tmpdir : get the temporary directory */
func (g *Guestfs) Get_tmpdir() (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("get_tmpdir")
	}

	r := C.guestfs_get_tmpdir(g.g)

	if r == nil {
		return "", get_error_from_handle(g, "get_tmpdir")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* get_trace : get command trace enabled flag */
func (g *Guestfs) Get_trace() (bool, *GuestfsError) {
	if g.g == nil {
		return false, closed_handle_error("get_trace")
	}

	r := C.guestfs_get_trace(g.g)

	if r == -1 {
		return false, get_error_from_handle(g, "get_trace")
	}
	return r != 0, nil
}

/* get_umask : get the current umask */
func (g *Guestfs) Get_umask() (int, *GuestfsError) {
	if g.g == nil {
		return 0, closed_handle_error("get_umask")
	}

	r := C.guestfs_get_umask(g.g)

	if r == -1 {
		return 0, get_error_from_handle(g, "get_umask")
	}
	return int(r), nil
}

/* get_verbose : get verbose mode */
func (g *Guestfs) Get_verbose() (bool, *GuestfsError) {
	if g.g == nil {
		return false, closed_handle_error("get_verbose")
	}

	r := C.guestfs_get_verbose(g.g)

	if r == -1 {
		return false, get_error_from_handle(g, "get_verbose")
	}
	return r != 0, nil
}

/* getcon : get SELinux security context */
func (g *Guestfs) Getcon() (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("getcon")
	}

	r := C.guestfs_getcon(g.g)

	if r == nil {
		return "", get_error_from_handle(g, "getcon")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* getxattr : get a single extended attribute */
func (g *Guestfs) Getxattr(path string, name string) ([]byte, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("getxattr")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	c_name := C.CString(name)
	defer C.free(unsafe.Pointer(c_name))

	var size C.size_t

	r := C.guestfs_getxattr(g.g, c_path, c_name, &size)

	if r == nil {
		return nil, get_error_from_handle(g, "getxattr")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoBytes(unsafe.Pointer(r), C.int(size)), nil
}

/* getxattrs : list extended attributes of a file or directory */
func (g *Guestfs) Getxattrs(path string) (*[]XAttr, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("getxattrs")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_getxattrs(g.g, c_path)

	if r == nil {
		return nil, get_error_from_handle(g, "getxattrs")
	}
	defer C.guestfs_free_xattr_list(r)
	return return_XAttr_list(r), nil
}

/* glob_expand : expand a wildcard path */
func (g *Guestfs) Glob_expand(pattern string) ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("glob_expand")
	}

	c_pattern := C.CString(pattern)
	defer C.free(unsafe.Pointer(c_pattern))

	r := C.guestfs_glob_expand(g.g, c_pattern)

	if r == nil {
		return nil, get_error_from_handle(g, "glob_expand")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* Struct carrying optional arguments for Grep */
type OptargsGrep struct {
	/* Extended field is ignored unless Extended_is_set == true */
	Extended_is_set bool
	Extended        bool
	/* Fixed field is ignored unless Fixed_is_set == true */
	Fixed_is_set bool
	Fixed        bool
	/* Insensitive field is ignored unless Insensitive_is_set == true */
	Insensitive_is_set bool
	Insensitive        bool
	/* Compressed field is ignored unless Compressed_is_set == true */
	Compressed_is_set bool
	Compressed        bool
}

/* grep : return lines matching a pattern */
func (g *Guestfs) Grep(regex string, path string, optargs *OptargsGrep) ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("grep")
	}

	c_regex := C.CString(regex)
	defer C.free(unsafe.Pointer(c_regex))

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))
	c_optargs := C.struct_guestfs_grep_opts_argv{}
	if optargs != nil {
		if optargs.Extended_is_set {
			c_optargs.bitmask |= C.GUESTFS_GREP_OPTS_EXTENDED_BITMASK
			if optargs.Extended {
				c_optargs.extended = 1
			} else {
				c_optargs.extended = 0
			}
		}
		if optargs.Fixed_is_set {
			c_optargs.bitmask |= C.GUESTFS_GREP_OPTS_FIXED_BITMASK
			if optargs.Fixed {
				c_optargs.fixed = 1
			} else {
				c_optargs.fixed = 0
			}
		}
		if optargs.Insensitive_is_set {
			c_optargs.bitmask |= C.GUESTFS_GREP_OPTS_INSENSITIVE_BITMASK
			if optargs.Insensitive {
				c_optargs.insensitive = 1
			} else {
				c_optargs.insensitive = 0
			}
		}
		if optargs.Compressed_is_set {
			c_optargs.bitmask |= C.GUESTFS_GREP_OPTS_COMPRESSED_BITMASK
			if optargs.Compressed {
				c_optargs.compressed = 1
			} else {
				c_optargs.compressed = 0
			}
		}
	}

	r := C.guestfs_grep_opts_argv(g.g, c_regex, c_path, &c_optargs)

	if r == nil {
		return nil, get_error_from_handle(g, "grep")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* grepi : return lines matching a pattern */
func (g *Guestfs) Grepi(regex string, path string) ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("grepi")
	}

	c_regex := C.CString(regex)
	defer C.free(unsafe.Pointer(c_regex))

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_grepi(g.g, c_regex, c_path)

	if r == nil {
		return nil, get_error_from_handle(g, "grepi")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* grub_install : install GRUB 1 */
func (g *Guestfs) Grub_install(root string, device string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("grub_install")
	}

	c_root := C.CString(root)
	defer C.free(unsafe.Pointer(c_root))

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_grub_install(g.g, c_root, c_device)

	if r == -1 {
		return get_error_from_handle(g, "grub_install")
	}
	return nil
}

/* head : return first 10 lines of a file */
func (g *Guestfs) Head(path string) ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("head")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_head(g.g, c_path)

	if r == nil {
		return nil, get_error_from_handle(g, "head")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* head_n : return first N lines of a file */
func (g *Guestfs) Head_n(nrlines int, path string) ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("head_n")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_head_n(g.g, C.int(nrlines), c_path)

	if r == nil {
		return nil, get_error_from_handle(g, "head_n")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* hexdump : dump a file in hexadecimal */
func (g *Guestfs) Hexdump(path string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("hexdump")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_hexdump(g.g, c_path)

	if r == nil {
		return "", get_error_from_handle(g, "hexdump")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* hivex_close : close the current hivex handle */
func (g *Guestfs) Hivex_close() *GuestfsError {
	if g.g == nil {
		return closed_handle_error("hivex_close")
	}

	r := C.guestfs_hivex_close(g.g)

	if r == -1 {
		return get_error_from_handle(g, "hivex_close")
	}
	return nil
}

/* hivex_commit : commit (write) changes back to the hive */
func (g *Guestfs) Hivex_commit(filename *string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("hivex_commit")
	}

	var c_filename *C.char = nil
	if filename != nil {
		c_filename = C.CString(*filename)
		defer C.free(unsafe.Pointer(c_filename))
	}

	r := C.guestfs_hivex_commit(g.g, c_filename)

	if r == -1 {
		return get_error_from_handle(g, "hivex_commit")
	}
	return nil
}

/* hivex_node_add_child : add a child node */
func (g *Guestfs) Hivex_node_add_child(parent int64, name string) (int64, *GuestfsError) {
	if g.g == nil {
		return 0, closed_handle_error("hivex_node_add_child")
	}

	c_name := C.CString(name)
	defer C.free(unsafe.Pointer(c_name))

	r := C.guestfs_hivex_node_add_child(g.g, C.int64_t(parent), c_name)

	if r == -1 {
		return 0, get_error_from_handle(g, "hivex_node_add_child")
	}
	return int64(r), nil
}

/* hivex_node_children : return list of nodes which are subkeys of node */
func (g *Guestfs) Hivex_node_children(nodeh int64) (*[]HivexNode, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("hivex_node_children")
	}

	r := C.guestfs_hivex_node_children(g.g, C.int64_t(nodeh))

	if r == nil {
		return nil, get_error_from_handle(g, "hivex_node_children")
	}
	defer C.guestfs_free_hivex_node_list(r)
	return return_HivexNode_list(r), nil
}

/* hivex_node_delete_child : delete a node (recursively) */
func (g *Guestfs) Hivex_node_delete_child(nodeh int64) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("hivex_node_delete_child")
	}

	r := C.guestfs_hivex_node_delete_child(g.g, C.int64_t(nodeh))

	if r == -1 {
		return get_error_from_handle(g, "hivex_node_delete_child")
	}
	return nil
}

/* hivex_node_get_child : return the named child of node */
func (g *Guestfs) Hivex_node_get_child(nodeh int64, name string) (int64, *GuestfsError) {
	if g.g == nil {
		return 0, closed_handle_error("hivex_node_get_child")
	}

	c_name := C.CString(name)
	defer C.free(unsafe.Pointer(c_name))

	r := C.guestfs_hivex_node_get_child(g.g, C.int64_t(nodeh), c_name)

	if r == -1 {
		return 0, get_error_from_handle(g, "hivex_node_get_child")
	}
	return int64(r), nil
}

/* hivex_node_get_value : return the named value */
func (g *Guestfs) Hivex_node_get_value(nodeh int64, key string) (int64, *GuestfsError) {
	if g.g == nil {
		return 0, closed_handle_error("hivex_node_get_value")
	}

	c_key := C.CString(key)
	defer C.free(unsafe.Pointer(c_key))

	r := C.guestfs_hivex_node_get_value(g.g, C.int64_t(nodeh), c_key)

	if r == -1 {
		return 0, get_error_from_handle(g, "hivex_node_get_value")
	}
	return int64(r), nil
}

/* hivex_node_name : return the name of the node */
func (g *Guestfs) Hivex_node_name(nodeh int64) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("hivex_node_name")
	}

	r := C.guestfs_hivex_node_name(g.g, C.int64_t(nodeh))

	if r == nil {
		return "", get_error_from_handle(g, "hivex_node_name")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* hivex_node_parent : return the parent of node */
func (g *Guestfs) Hivex_node_parent(nodeh int64) (int64, *GuestfsError) {
	if g.g == nil {
		return 0, closed_handle_error("hivex_node_parent")
	}

	r := C.guestfs_hivex_node_parent(g.g, C.int64_t(nodeh))

	if r == -1 {
		return 0, get_error_from_handle(g, "hivex_node_parent")
	}
	return int64(r), nil
}

/* hivex_node_set_value : set or replace a single value in a node */
func (g *Guestfs) Hivex_node_set_value(nodeh int64, key string, t int64, val []byte) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("hivex_node_set_value")
	}

	c_key := C.CString(key)
	defer C.free(unsafe.Pointer(c_key))

	/* string() cast here is apparently safe because
	 *   "Converting a slice of bytes to a string type yields
	 *   a string whose successive bytes are the elements of
	 *   the slice."
	 */
	c_val := C.CString(string(val))
	defer C.free(unsafe.Pointer(c_val))

	r := C.guestfs_hivex_node_set_value(g.g, C.int64_t(nodeh), c_key, C.int64_t(t), c_val, C.size_t(len(val)))

	if r == -1 {
		return get_error_from_handle(g, "hivex_node_set_value")
	}
	return nil
}

/* hivex_node_values : return list of values attached to node */
func (g *Guestfs) Hivex_node_values(nodeh int64) (*[]HivexValue, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("hivex_node_values")
	}

	r := C.guestfs_hivex_node_values(g.g, C.int64_t(nodeh))

	if r == nil {
		return nil, get_error_from_handle(g, "hivex_node_values")
	}
	defer C.guestfs_free_hivex_value_list(r)
	return return_HivexValue_list(r), nil
}

/* Struct carrying optional arguments for Hivex_open */
type OptargsHivex_open struct {
	/* Verbose field is ignored unless Verbose_is_set == true */
	Verbose_is_set bool
	Verbose        bool
	/* Debug field is ignored unless Debug_is_set == true */
	Debug_is_set bool
	Debug        bool
	/* Write field is ignored unless Write_is_set == true */
	Write_is_set bool
	Write        bool
}

/* hivex_open : open a Windows Registry hive file */
func (g *Guestfs) Hivex_open(filename string, optargs *OptargsHivex_open) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("hivex_open")
	}

	c_filename := C.CString(filename)
	defer C.free(unsafe.Pointer(c_filename))
	c_optargs := C.struct_guestfs_hivex_open_argv{}
	if optargs != nil {
		if optargs.Verbose_is_set {
			c_optargs.bitmask |= C.GUESTFS_HIVEX_OPEN_VERBOSE_BITMASK
			if optargs.Verbose {
				c_optargs.verbose = 1
			} else {
				c_optargs.verbose = 0
			}
		}
		if optargs.Debug_is_set {
			c_optargs.bitmask |= C.GUESTFS_HIVEX_OPEN_DEBUG_BITMASK
			if optargs.Debug {
				c_optargs.debug = 1
			} else {
				c_optargs.debug = 0
			}
		}
		if optargs.Write_is_set {
			c_optargs.bitmask |= C.GUESTFS_HIVEX_OPEN_WRITE_BITMASK
			if optargs.Write {
				c_optargs.write = 1
			} else {
				c_optargs.write = 0
			}
		}
	}

	r := C.guestfs_hivex_open_argv(g.g, c_filename, &c_optargs)

	if r == -1 {
		return get_error_from_handle(g, "hivex_open")
	}
	return nil
}

/* hivex_root : return the root node of the hive */
func (g *Guestfs) Hivex_root() (int64, *GuestfsError) {
	if g.g == nil {
		return 0, closed_handle_error("hivex_root")
	}

	r := C.guestfs_hivex_root(g.g)

	if r == -1 {
		return 0, get_error_from_handle(g, "hivex_root")
	}
	return int64(r), nil
}

/* hivex_value_key : return the key field from the (key, datatype, data) tuple */
func (g *Guestfs) Hivex_value_key(valueh int64) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("hivex_value_key")
	}

	r := C.guestfs_hivex_value_key(g.g, C.int64_t(valueh))

	if r == nil {
		return "", get_error_from_handle(g, "hivex_value_key")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* hivex_value_type : return the data type from the (key, datatype, data) tuple */
func (g *Guestfs) Hivex_value_type(valueh int64) (int64, *GuestfsError) {
	if g.g == nil {
		return 0, closed_handle_error("hivex_value_type")
	}

	r := C.guestfs_hivex_value_type(g.g, C.int64_t(valueh))

	if r == -1 {
		return 0, get_error_from_handle(g, "hivex_value_type")
	}
	return int64(r), nil
}

/* hivex_value_utf8 : return the data field from the (key, datatype, data) tuple */
func (g *Guestfs) Hivex_value_utf8(valueh int64) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("hivex_value_utf8")
	}

	r := C.guestfs_hivex_value_utf8(g.g, C.int64_t(valueh))

	if r == nil {
		return "", get_error_from_handle(g, "hivex_value_utf8")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* hivex_value_value : return the data field from the (key, datatype, data) tuple */
func (g *Guestfs) Hivex_value_value(valueh int64) ([]byte, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("hivex_value_value")
	}

	var size C.size_t

	r := C.guestfs_hivex_value_value(g.g, C.int64_t(valueh), &size)

	if r == nil {
		return nil, get_error_from_handle(g, "hivex_value_value")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoBytes(unsafe.Pointer(r), C.int(size)), nil
}

/* initrd_cat : list the contents of a single file in an initrd */
func (g *Guestfs) Initrd_cat(initrdpath string, filename string) ([]byte, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("initrd_cat")
	}

	c_initrdpath := C.CString(initrdpath)
	defer C.free(unsafe.Pointer(c_initrdpath))

	c_filename := C.CString(filename)
	defer C.free(unsafe.Pointer(c_filename))

	var size C.size_t

	r := C.guestfs_initrd_cat(g.g, c_initrdpath, c_filename, &size)

	if r == nil {
		return nil, get_error_from_handle(g, "initrd_cat")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoBytes(unsafe.Pointer(r), C.int(size)), nil
}

/* initrd_list : list files in an initrd */
func (g *Guestfs) Initrd_list(path string) ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("initrd_list")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_initrd_list(g.g, c_path)

	if r == nil {
		return nil, get_error_from_handle(g, "initrd_list")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* inotify_add_watch : add an inotify watch */
func (g *Guestfs) Inotify_add_watch(path string, mask int) (int64, *GuestfsError) {
	if g.g == nil {
		return 0, closed_handle_error("inotify_add_watch")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_inotify_add_watch(g.g, c_path, C.int(mask))

	if r == -1 {
		return 0, get_error_from_handle(g, "inotify_add_watch")
	}
	return int64(r), nil
}

/* inotify_close : close the inotify handle */
func (g *Guestfs) Inotify_close() *GuestfsError {
	if g.g == nil {
		return closed_handle_error("inotify_close")
	}

	r := C.guestfs_inotify_close(g.g)

	if r == -1 {
		return get_error_from_handle(g, "inotify_close")
	}
	return nil
}

/* inotify_files : return list of watched files that had events */
func (g *Guestfs) Inotify_files() ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("inotify_files")
	}

	r := C.guestfs_inotify_files(g.g)

	if r == nil {
		return nil, get_error_from_handle(g, "inotify_files")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* inotify_init : create an inotify handle */
func (g *Guestfs) Inotify_init(maxevents int) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("inotify_init")
	}

	r := C.guestfs_inotify_init(g.g, C.int(maxevents))

	if r == -1 {
		return get_error_from_handle(g, "inotify_init")
	}
	return nil
}

/* inotify_read : return list of inotify events */
func (g *Guestfs) Inotify_read() (*[]INotifyEvent, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("inotify_read")
	}

	r := C.guestfs_inotify_read(g.g)

	if r == nil {
		return nil, get_error_from_handle(g, "inotify_read")
	}
	defer C.guestfs_free_inotify_event_list(r)
	return return_INotifyEvent_list(r), nil
}

/* inotify_rm_watch : remove an inotify watch */
func (g *Guestfs) Inotify_rm_watch(wd int) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("inotify_rm_watch")
	}

	r := C.guestfs_inotify_rm_watch(g.g, C.int(wd))

	if r == -1 {
		return get_error_from_handle(g, "inotify_rm_watch")
	}
	return nil
}

/* inspect_get_arch : get architecture of inspected operating system */
func (g *Guestfs) Inspect_get_arch(root string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("inspect_get_arch")
	}

	c_root := C.CString(root)
	defer C.free(unsafe.Pointer(c_root))

	r := C.guestfs_inspect_get_arch(g.g, c_root)

	if r == nil {
		return "", get_error_from_handle(g, "inspect_get_arch")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* inspect_get_distro : get distro of inspected operating system */
func (g *Guestfs) Inspect_get_distro(root string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("inspect_get_distro")
	}

	c_root := C.CString(root)
	defer C.free(unsafe.Pointer(c_root))

	r := C.guestfs_inspect_get_distro(g.g, c_root)

	if r == nil {
		return "", get_error_from_handle(g, "inspect_get_distro")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* inspect_get_drive_mappings : get drive letter mappings */
func (g *Guestfs) Inspect_get_drive_mappings(root string) (map[string]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("inspect_get_drive_mappings")
	}

	c_root := C.CString(root)
	defer C.free(unsafe.Pointer(c_root))

	r := C.guestfs_inspect_get_drive_mappings(g.g, c_root)

	if r == nil {
		return nil, get_error_from_handle(g, "inspect_get_drive_mappings")
	}
	defer free_string_list(r)
	return return_hashtable(r), nil
}

/* inspect_get_filesystems : get filesystems associated with inspected operating system */
func (g *Guestfs) Inspect_get_filesystems(root string) ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("inspect_get_filesystems")
	}

	c_root := C.CString(root)
	defer C.free(unsafe.Pointer(c_root))

	r := C.guestfs_inspect_get_filesystems(g.g, c_root)

	if r == nil {
		return nil, get_error_from_handle(g, "inspect_get_filesystems")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* inspect_get_format : get format of inspected operating system */
func (g *Guestfs) Inspect_get_format(root string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("inspect_get_format")
	}

	c_root := C.CString(root)
	defer C.free(unsafe.Pointer(c_root))

	r := C.guestfs_inspect_get_format(g.g, c_root)

	if r == nil {
		return "", get_error_from_handle(g, "inspect_get_format")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* inspect_get_hostname : get hostname of the operating system */
func (g *Guestfs) Inspect_get_hostname(root string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("inspect_get_hostname")
	}

	c_root := C.CString(root)
	defer C.free(unsafe.Pointer(c_root))

	r := C.guestfs_inspect_get_hostname(g.g, c_root)

	if r == nil {
		return "", get_error_from_handle(g, "inspect_get_hostname")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* Struct carrying optional arguments for Inspect_get_icon */
type OptargsInspect_get_icon struct {
	/* Favicon field is ignored unless Favicon_is_set == true */
	Favicon_is_set bool
	Favicon        bool
	/* Highquality field is ignored unless Highquality_is_set == true */
	Highquality_is_set bool
	Highquality        bool
}

/* inspect_get_icon : get the icon corresponding to this operating system */
func (g *Guestfs) Inspect_get_icon(root string, optargs *OptargsInspect_get_icon) ([]byte, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("inspect_get_icon")
	}

	c_root := C.CString(root)
	defer C.free(unsafe.Pointer(c_root))
	c_optargs := C.struct_guestfs_inspect_get_icon_argv{}
	if optargs != nil {
		if optargs.Favicon_is_set {
			c_optargs.bitmask |= C.GUESTFS_INSPECT_GET_ICON_FAVICON_BITMASK
			if optargs.Favicon {
				c_optargs.favicon = 1
			} else {
				c_optargs.favicon = 0
			}
		}
		if optargs.Highquality_is_set {
			c_optargs.bitmask |= C.GUESTFS_INSPECT_GET_ICON_HIGHQUALITY_BITMASK
			if optargs.Highquality {
				c_optargs.highquality = 1
			} else {
				c_optargs.highquality = 0
			}
		}
	}

	var size C.size_t

	r := C.guestfs_inspect_get_icon_argv(g.g, c_root, &size, &c_optargs)

	if r == nil {
		return nil, get_error_from_handle(g, "inspect_get_icon")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoBytes(unsafe.Pointer(r), C.int(size)), nil
}

/* inspect_get_major_version : get major version of inspected operating system */
func (g *Guestfs) Inspect_get_major_version(root string) (int, *GuestfsError) {
	if g.g == nil {
		return 0, closed_handle_error("inspect_get_major_version")
	}

	c_root := C.CString(root)
	defer C.free(unsafe.Pointer(c_root))

	r := C.guestfs_inspect_get_major_version(g.g, c_root)

	if r == -1 {
		return 0, get_error_from_handle(g, "inspect_get_major_version")
	}
	return int(r), nil
}

/* inspect_get_minor_version : get minor version of inspected operating system */
func (g *Guestfs) Inspect_get_minor_version(root string) (int, *GuestfsError) {
	if g.g == nil {
		return 0, closed_handle_error("inspect_get_minor_version")
	}

	c_root := C.CString(root)
	defer C.free(unsafe.Pointer(c_root))

	r := C.guestfs_inspect_get_minor_version(g.g, c_root)

	if r == -1 {
		return 0, get_error_from_handle(g, "inspect_get_minor_version")
	}
	return int(r), nil
}

/* inspect_get_mountpoints : get mountpoints of inspected operating system */
func (g *Guestfs) Inspect_get_mountpoints(root string) (map[string]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("inspect_get_mountpoints")
	}

	c_root := C.CString(root)
	defer C.free(unsafe.Pointer(c_root))

	r := C.guestfs_inspect_get_mountpoints(g.g, c_root)

	if r == nil {
		return nil, get_error_from_handle(g, "inspect_get_mountpoints")
	}
	defer free_string_list(r)
	return return_hashtable(r), nil
}

/* inspect_get_package_format : get package format used by the operating system */
func (g *Guestfs) Inspect_get_package_format(root string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("inspect_get_package_format")
	}

	c_root := C.CString(root)
	defer C.free(unsafe.Pointer(c_root))

	r := C.guestfs_inspect_get_package_format(g.g, c_root)

	if r == nil {
		return "", get_error_from_handle(g, "inspect_get_package_format")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* inspect_get_package_management : get package management tool used by the operating system */
func (g *Guestfs) Inspect_get_package_management(root string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("inspect_get_package_management")
	}

	c_root := C.CString(root)
	defer C.free(unsafe.Pointer(c_root))

	r := C.guestfs_inspect_get_package_management(g.g, c_root)

	if r == nil {
		return "", get_error_from_handle(g, "inspect_get_package_management")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* inspect_get_product_name : get product name of inspected operating system */
func (g *Guestfs) Inspect_get_product_name(root string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("inspect_get_product_name")
	}

	c_root := C.CString(root)
	defer C.free(unsafe.Pointer(c_root))

	r := C.guestfs_inspect_get_product_name(g.g, c_root)

	if r == nil {
		return "", get_error_from_handle(g, "inspect_get_product_name")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* inspect_get_product_variant : get product variant of inspected operating system */
func (g *Guestfs) Inspect_get_product_variant(root string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("inspect_get_product_variant")
	}

	c_root := C.CString(root)
	defer C.free(unsafe.Pointer(c_root))

	r := C.guestfs_inspect_get_product_variant(g.g, c_root)

	if r == nil {
		return "", get_error_from_handle(g, "inspect_get_product_variant")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* inspect_get_roots : return list of operating systems found by last inspection */
func (g *Guestfs) Inspect_get_roots() ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("inspect_get_roots")
	}

	r := C.guestfs_inspect_get_roots(g.g)

	if r == nil {
		return nil, get_error_from_handle(g, "inspect_get_roots")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* inspect_get_type : get type of inspected operating system */
func (g *Guestfs) Inspect_get_type(root string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("inspect_get_type")
	}

	c_root := C.CString(root)
	defer C.free(unsafe.Pointer(c_root))

	r := C.guestfs_inspect_get_type(g.g, c_root)

	if r == nil {
		return "", get_error_from_handle(g, "inspect_get_type")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* inspect_get_windows_current_control_set : get Windows CurrentControlSet of inspected operating system */
func (g *Guestfs) Inspect_get_windows_current_control_set(root string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("inspect_get_windows_current_control_set")
	}

	c_root := C.CString(root)
	defer C.free(unsafe.Pointer(c_root))

	r := C.guestfs_inspect_get_windows_current_control_set(g.g, c_root)

	if r == nil {
		return "", get_error_from_handle(g, "inspect_get_windows_current_control_set")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* inspect_get_windows_systemroot : get Windows systemroot of inspected operating system */
func (g *Guestfs) Inspect_get_windows_systemroot(root string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("inspect_get_windows_systemroot")
	}

	c_root := C.CString(root)
	defer C.free(unsafe.Pointer(c_root))

	r := C.guestfs_inspect_get_windows_systemroot(g.g, c_root)

	if r == nil {
		return "", get_error_from_handle(g, "inspect_get_windows_systemroot")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* inspect_is_live : get live flag for install disk */
func (g *Guestfs) Inspect_is_live(root string) (bool, *GuestfsError) {
	if g.g == nil {
		return false, closed_handle_error("inspect_is_live")
	}

	c_root := C.CString(root)
	defer C.free(unsafe.Pointer(c_root))

	r := C.guestfs_inspect_is_live(g.g, c_root)

	if r == -1 {
		return false, get_error_from_handle(g, "inspect_is_live")
	}
	return r != 0, nil
}

/* inspect_is_multipart : get multipart flag for install disk */
func (g *Guestfs) Inspect_is_multipart(root string) (bool, *GuestfsError) {
	if g.g == nil {
		return false, closed_handle_error("inspect_is_multipart")
	}

	c_root := C.CString(root)
	defer C.free(unsafe.Pointer(c_root))

	r := C.guestfs_inspect_is_multipart(g.g, c_root)

	if r == -1 {
		return false, get_error_from_handle(g, "inspect_is_multipart")
	}
	return r != 0, nil
}

/* inspect_is_netinst : get netinst (network installer) flag for install disk */
func (g *Guestfs) Inspect_is_netinst(root string) (bool, *GuestfsError) {
	if g.g == nil {
		return false, closed_handle_error("inspect_is_netinst")
	}

	c_root := C.CString(root)
	defer C.free(unsafe.Pointer(c_root))

	r := C.guestfs_inspect_is_netinst(g.g, c_root)

	if r == -1 {
		return false, get_error_from_handle(g, "inspect_is_netinst")
	}
	return r != 0, nil
}

/* inspect_list_applications : get list of applications installed in the operating system */
func (g *Guestfs) Inspect_list_applications(root string) (*[]Application, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("inspect_list_applications")
	}

	c_root := C.CString(root)
	defer C.free(unsafe.Pointer(c_root))

	r := C.guestfs_inspect_list_applications(g.g, c_root)

	if r == nil {
		return nil, get_error_from_handle(g, "inspect_list_applications")
	}
	defer C.guestfs_free_application_list(r)
	return return_Application_list(r), nil
}

/* inspect_list_applications2 : get list of applications installed in the operating system */
func (g *Guestfs) Inspect_list_applications2(root string) (*[]Application2, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("inspect_list_applications2")
	}

	c_root := C.CString(root)
	defer C.free(unsafe.Pointer(c_root))

	r := C.guestfs_inspect_list_applications2(g.g, c_root)

	if r == nil {
		return nil, get_error_from_handle(g, "inspect_list_applications2")
	}
	defer C.guestfs_free_application2_list(r)
	return return_Application2_list(r), nil
}

/* inspect_os : inspect disk and return list of operating systems found */
func (g *Guestfs) Inspect_os() ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("inspect_os")
	}

	r := C.guestfs_inspect_os(g.g)

	if r == nil {
		return nil, get_error_from_handle(g, "inspect_os")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* internal_exit : cause the daemon to exit (internal use only) */
func (g *Guestfs) Internal_exit() *GuestfsError {
	if g.g == nil {
		return closed_handle_error("internal_exit")
	}

	r := C.guestfs_internal_exit(g.g)

	if r == -1 {
		return get_error_from_handle(g, "internal_exit")
	}
	return nil
}

/* Struct carrying optional arguments for Internal_test */
type OptargsInternal_test struct {
	/* Obool field is ignored unless Obool_is_set == true */
	Obool_is_set bool
	Obool        bool
	/* Oint field is ignored unless Oint_is_set == true */
	Oint_is_set bool
	Oint        int
	/* Oint64 field is ignored unless Oint64_is_set == true */
	Oint64_is_set bool
	Oint64        int64
	/* Ostring field is ignored unless Ostring_is_set == true */
	Ostring_is_set bool
	Ostring        string
	/* Ostringlist field is ignored unless Ostringlist_is_set == true */
	Ostringlist_is_set bool
	Ostringlist        []string
}

/* internal_test : internal test function - do not use */
func (g *Guestfs) Internal_test(str string, optstr *string, strlist []string, b bool, integer int, integer64 int64, filein string, fileout string, bufferin []byte, optargs *OptargsInternal_test) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("internal_test")
	}

	c_str := C.CString(str)
	defer C.free(unsafe.Pointer(c_str))

	var c_optstr *C.char = nil
	if optstr != nil {
		c_optstr = C.CString(*optstr)
		defer C.free(unsafe.Pointer(c_optstr))
	}

	c_strlist := arg_string_list(strlist)
	defer free_string_list(c_strlist)

	var c_b C.int
	if b {
		c_b = 1
	} else {
		c_b = 0
	}

	c_filein := C.CString(filein)
	defer C.free(unsafe.Pointer(c_filein))

	c_fileout := C.CString(fileout)
	defer C.free(unsafe.Pointer(c_fileout))

	/* string() cast here is apparently safe because
	 *   "Converting a slice of bytes to a string type yields
	 *   a string whose successive bytes are the elements of
	 *   the slice."
	 */
	c_bufferin := C.CString(string(bufferin))
	defer C.free(unsafe.Pointer(c_bufferin))
	c_optargs := C.struct_guestfs_internal_test_argv{}
	if optargs != nil {
		if optargs.Obool_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_OBOOL_BITMASK
			if optargs.Obool {
				c_optargs.obool = 1
			} else {
				c_optargs.obool = 0
			}
		}
		if optargs.Oint_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_OINT_BITMASK
			c_optargs.oint = C.int(optargs.Oint)
		}
		if optargs.Oint64_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_OINT64_BITMASK
			c_optargs.oint64 = C.int64_t(optargs.Oint64)
		}
		if optargs.Ostring_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_OSTRING_BITMASK
			c_optargs.ostring = C.CString(optargs.Ostring)
			defer C.free(unsafe.Pointer(c_optargs.ostring))
		}
		if optargs.Ostringlist_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_OSTRINGLIST_BITMASK
			c_optargs.ostringlist = arg_string_list(optargs.Ostringlist)
			defer free_string_list(c_optargs.ostringlist)
		}
	}

	r := C.guestfs_internal_test_argv(g.g, c_str, c_optstr, c_strlist, c_b, C.int(integer), C.int64_t(integer64), c_filein, c_fileout, c_bufferin, C.size_t(len(bufferin)), &c_optargs)

	if r == -1 {
		return get_error_from_handle(g, "internal_test")
	}
	return nil
}

/* Struct carrying optional arguments for Internal_test_63_optargs */
type OptargsInternal_test_63_optargs struct {
	/* Opt1 field is ignored unless Opt1_is_set == true */
	Opt1_is_set bool
	Opt1        int
	/* Opt2 field is ignored unless Opt2_is_set == true */
	Opt2_is_set bool
	Opt2        int
	/* Opt3 field is ignored unless Opt3_is_set == true */
	Opt3_is_set bool
	Opt3        int
	/* Opt4 field is ignored unless Opt4_is_set == true */
	Opt4_is_set bool
	Opt4        int
	/* Opt5 field is ignored unless Opt5_is_set == true */
	Opt5_is_set bool
	Opt5        int
	/* Opt6 field is ignored unless Opt6_is_set == true */
	Opt6_is_set bool
	Opt6        int
	/* Opt7 field is ignored unless Opt7_is_set == true */
	Opt7_is_set bool
	Opt7        int
	/* Opt8 field is ignored unless Opt8_is_set == true */
	Opt8_is_set bool
	Opt8        int
	/* Opt9 field is ignored unless Opt9_is_set == true */
	Opt9_is_set bool
	Opt9        int
	/* Opt10 field is ignored unless Opt10_is_set == true */
	Opt10_is_set bool
	Opt10        int
	/* Opt11 field is ignored unless Opt11_is_set == true */
	Opt11_is_set bool
	Opt11        int
	/* Opt12 field is ignored unless Opt12_is_set == true */
	Opt12_is_set bool
	Opt12        int
	/* Opt13 field is ignored unless Opt13_is_set == true */
	Opt13_is_set bool
	Opt13        int
	/* Opt14 field is ignored unless Opt14_is_set == true */
	Opt14_is_set bool
	Opt14        int
	/* Opt15 field is ignored unless Opt15_is_set == true */
	Opt15_is_set bool
	Opt15        int
	/* Opt16 field is ignored unless Opt16_is_set == true */
	Opt16_is_set bool
	Opt16        int
	/* Opt17 field is ignored unless Opt17_is_set == true */
	Opt17_is_set bool
	Opt17        int
	/* Opt18 field is ignored unless Opt18_is_set == true */
	Opt18_is_set bool
	Opt18        int
	/* Opt19 field is ignored unless Opt19_is_set == true */
	Opt19_is_set bool
	Opt19        int
	/* Opt20 field is ignored unless Opt20_is_set == true */
	Opt20_is_set bool
	Opt20        int
	/* Opt21 field is ignored unless Opt21_is_set == true */
	Opt21_is_set bool
	Opt21        int
	/* Opt22 field is ignored unless Opt22_is_set == true */
	Opt22_is_set bool
	Opt22        int
	/* Opt23 field is ignored unless Opt23_is_set == true */
	Opt23_is_set bool
	Opt23        int
	/* Opt24 field is ignored unless Opt24_is_set == true */
	Opt24_is_set bool
	Opt24        int
	/* Opt25 field is ignored unless Opt25_is_set == true */
	Opt25_is_set bool
	Opt25        int
	/* Opt26 field is ignored unless Opt26_is_set == true */
	Opt26_is_set bool
	Opt26        int
	/* Opt27 field is ignored unless Opt27_is_set == true */
	Opt27_is_set bool
	Opt27        int
	/* Opt28 field is ignored unless Opt28_is_set == true */
	Opt28_is_set bool
	Opt28        int
	/* Opt29 field is ignored unless Opt29_is_set == true */
	Opt29_is_set bool
	Opt29        int
	/* Opt30 field is ignored unless Opt30_is_set == true */
	Opt30_is_set bool
	Opt30        int
	/* Opt31 field is ignored unless Opt31_is_set == true */
	Opt31_is_set bool
	Opt31        int
	/* Opt32 field is ignored unless Opt32_is_set == true */
	Opt32_is_set bool
	Opt32        int
	/* Opt33 field is ignored unless Opt33_is_set == true */
	Opt33_is_set bool
	Opt33        int
	/* Opt34 field is ignored unless Opt34_is_set == true */
	Opt34_is_set bool
	Opt34        int
	/* Opt35 field is ignored unless Opt35_is_set == true */
	Opt35_is_set bool
	Opt35        int
	/* Opt36 field is ignored unless Opt36_is_set == true */
	Opt36_is_set bool
	Opt36        int
	/* Opt37 field is ignored unless Opt37_is_set == true */
	Opt37_is_set bool
	Opt37        int
	/* Opt38 field is ignored unless Opt38_is_set == true */
	Opt38_is_set bool
	Opt38        int
	/* Opt39 field is ignored unless Opt39_is_set == true */
	Opt39_is_set bool
	Opt39        int
	/* Opt40 field is ignored unless Opt40_is_set == true */
	Opt40_is_set bool
	Opt40        int
	/* Opt41 field is ignored unless Opt41_is_set == true */
	Opt41_is_set bool
	Opt41        int
	/* Opt42 field is ignored unless Opt42_is_set == true */
	Opt42_is_set bool
	Opt42        int
	/* Opt43 field is ignored unless Opt43_is_set == true */
	Opt43_is_set bool
	Opt43        int
	/* Opt44 field is ignored unless Opt44_is_set == true */
	Opt44_is_set bool
	Opt44        int
	/* Opt45 field is ignored unless Opt45_is_set == true */
	Opt45_is_set bool
	Opt45        int
	/* Opt46 field is ignored unless Opt46_is_set == true */
	Opt46_is_set bool
	Opt46        int
	/* Opt47 field is ignored unless Opt47_is_set == true */
	Opt47_is_set bool
	Opt47        int
	/* Opt48 field is ignored unless Opt48_is_set == true */
	Opt48_is_set bool
	Opt48        int
	/* Opt49 field is ignored unless Opt49_is_set == true */
	Opt49_is_set bool
	Opt49        int
	/* Opt50 field is ignored unless Opt50_is_set == true */
	Opt50_is_set bool
	Opt50        int
	/* Opt51 field is ignored unless Opt51_is_set == true */
	Opt51_is_set bool
	Opt51        int
	/* Opt52 field is ignored unless Opt52_is_set == true */
	Opt52_is_set bool
	Opt52        int
	/* Opt53 field is ignored unless Opt53_is_set == true */
	Opt53_is_set bool
	Opt53        int
	/* Opt54 field is ignored unless Opt54_is_set == true */
	Opt54_is_set bool
	Opt54        int
	/* Opt55 field is ignored unless Opt55_is_set == true */
	Opt55_is_set bool
	Opt55        int
	/* Opt56 field is ignored unless Opt56_is_set == true */
	Opt56_is_set bool
	Opt56        int
	/* Opt57 field is ignored unless Opt57_is_set == true */
	Opt57_is_set bool
	Opt57        int
	/* Opt58 field is ignored unless Opt58_is_set == true */
	Opt58_is_set bool
	Opt58        int
	/* Opt59 field is ignored unless Opt59_is_set == true */
	Opt59_is_set bool
	Opt59        int
	/* Opt60 field is ignored unless Opt60_is_set == true */
	Opt60_is_set bool
	Opt60        int
	/* Opt61 field is ignored unless Opt61_is_set == true */
	Opt61_is_set bool
	Opt61        int
	/* Opt62 field is ignored unless Opt62_is_set == true */
	Opt62_is_set bool
	Opt62        int
	/* Opt63 field is ignored unless Opt63_is_set == true */
	Opt63_is_set bool
	Opt63        int
}

/* internal_test_63_optargs : internal test function - do not use */
func (g *Guestfs) Internal_test_63_optargs(optargs *OptargsInternal_test_63_optargs) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("internal_test_63_optargs")
	}
	c_optargs := C.struct_guestfs_internal_test_63_optargs_argv{}
	if optargs != nil {
		if optargs.Opt1_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT1_BITMASK
			c_optargs.opt1 = C.int(optargs.Opt1)
		}
		if optargs.Opt2_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT2_BITMASK
			c_optargs.opt2 = C.int(optargs.Opt2)
		}
		if optargs.Opt3_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT3_BITMASK
			c_optargs.opt3 = C.int(optargs.Opt3)
		}
		if optargs.Opt4_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT4_BITMASK
			c_optargs.opt4 = C.int(optargs.Opt4)
		}
		if optargs.Opt5_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT5_BITMASK
			c_optargs.opt5 = C.int(optargs.Opt5)
		}
		if optargs.Opt6_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT6_BITMASK
			c_optargs.opt6 = C.int(optargs.Opt6)
		}
		if optargs.Opt7_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT7_BITMASK
			c_optargs.opt7 = C.int(optargs.Opt7)
		}
		if optargs.Opt8_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT8_BITMASK
			c_optargs.opt8 = C.int(optargs.Opt8)
		}
		if optargs.Opt9_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT9_BITMASK
			c_optargs.opt9 = C.int(optargs.Opt9)
		}
		if optargs.Opt10_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT10_BITMASK
			c_optargs.opt10 = C.int(optargs.Opt10)
		}
		if optargs.Opt11_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT11_BITMASK
			c_optargs.opt11 = C.int(optargs.Opt11)
		}
		if optargs.Opt12_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT12_BITMASK
			c_optargs.opt12 = C.int(optargs.Opt12)
		}
		if optargs.Opt13_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT13_BITMASK
			c_optargs.opt13 = C.int(optargs.Opt13)
		}
		if optargs.Opt14_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT14_BITMASK
			c_optargs.opt14 = C.int(optargs.Opt14)
		}
		if optargs.Opt15_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT15_BITMASK
			c_optargs.opt15 = C.int(optargs.Opt15)
		}
		if optargs.Opt16_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT16_BITMASK
			c_optargs.opt16 = C.int(optargs.Opt16)
		}
		if optargs.Opt17_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT17_BITMASK
			c_optargs.opt17 = C.int(optargs.Opt17)
		}
		if optargs.Opt18_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT18_BITMASK
			c_optargs.opt18 = C.int(optargs.Opt18)
		}
		if optargs.Opt19_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT19_BITMASK
			c_optargs.opt19 = C.int(optargs.Opt19)
		}
		if optargs.Opt20_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT20_BITMASK
			c_optargs.opt20 = C.int(optargs.Opt20)
		}
		if optargs.Opt21_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT21_BITMASK
			c_optargs.opt21 = C.int(optargs.Opt21)
		}
		if optargs.Opt22_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT22_BITMASK
			c_optargs.opt22 = C.int(optargs.Opt22)
		}
		if optargs.Opt23_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT23_BITMASK
			c_optargs.opt23 = C.int(optargs.Opt23)
		}
		if optargs.Opt24_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT24_BITMASK
			c_optargs.opt24 = C.int(optargs.Opt24)
		}
		if optargs.Opt25_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT25_BITMASK
			c_optargs.opt25 = C.int(optargs.Opt25)
		}
		if optargs.Opt26_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT26_BITMASK
			c_optargs.opt26 = C.int(optargs.Opt26)
		}
		if optargs.Opt27_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT27_BITMASK
			c_optargs.opt27 = C.int(optargs.Opt27)
		}
		if optargs.Opt28_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT28_BITMASK
			c_optargs.opt28 = C.int(optargs.Opt28)
		}
		if optargs.Opt29_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT29_BITMASK
			c_optargs.opt29 = C.int(optargs.Opt29)
		}
		if optargs.Opt30_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT30_BITMASK
			c_optargs.opt30 = C.int(optargs.Opt30)
		}
		if optargs.Opt31_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT31_BITMASK
			c_optargs.opt31 = C.int(optargs.Opt31)
		}
		if optargs.Opt32_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT32_BITMASK
			c_optargs.opt32 = C.int(optargs.Opt32)
		}
		if optargs.Opt33_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT33_BITMASK
			c_optargs.opt33 = C.int(optargs.Opt33)
		}
		if optargs.Opt34_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT34_BITMASK
			c_optargs.opt34 = C.int(optargs.Opt34)
		}
		if optargs.Opt35_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT35_BITMASK
			c_optargs.opt35 = C.int(optargs.Opt35)
		}
		if optargs.Opt36_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT36_BITMASK
			c_optargs.opt36 = C.int(optargs.Opt36)
		}
		if optargs.Opt37_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT37_BITMASK
			c_optargs.opt37 = C.int(optargs.Opt37)
		}
		if optargs.Opt38_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT38_BITMASK
			c_optargs.opt38 = C.int(optargs.Opt38)
		}
		if optargs.Opt39_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT39_BITMASK
			c_optargs.opt39 = C.int(optargs.Opt39)
		}
		if optargs.Opt40_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT40_BITMASK
			c_optargs.opt40 = C.int(optargs.Opt40)
		}
		if optargs.Opt41_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT41_BITMASK
			c_optargs.opt41 = C.int(optargs.Opt41)
		}
		if optargs.Opt42_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT42_BITMASK
			c_optargs.opt42 = C.int(optargs.Opt42)
		}
		if optargs.Opt43_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT43_BITMASK
			c_optargs.opt43 = C.int(optargs.Opt43)
		}
		if optargs.Opt44_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT44_BITMASK
			c_optargs.opt44 = C.int(optargs.Opt44)
		}
		if optargs.Opt45_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT45_BITMASK
			c_optargs.opt45 = C.int(optargs.Opt45)
		}
		if optargs.Opt46_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT46_BITMASK
			c_optargs.opt46 = C.int(optargs.Opt46)
		}
		if optargs.Opt47_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT47_BITMASK
			c_optargs.opt47 = C.int(optargs.Opt47)
		}
		if optargs.Opt48_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT48_BITMASK
			c_optargs.opt48 = C.int(optargs.Opt48)
		}
		if optargs.Opt49_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT49_BITMASK
			c_optargs.opt49 = C.int(optargs.Opt49)
		}
		if optargs.Opt50_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT50_BITMASK
			c_optargs.opt50 = C.int(optargs.Opt50)
		}
		if optargs.Opt51_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT51_BITMASK
			c_optargs.opt51 = C.int(optargs.Opt51)
		}
		if optargs.Opt52_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT52_BITMASK
			c_optargs.opt52 = C.int(optargs.Opt52)
		}
		if optargs.Opt53_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT53_BITMASK
			c_optargs.opt53 = C.int(optargs.Opt53)
		}
		if optargs.Opt54_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT54_BITMASK
			c_optargs.opt54 = C.int(optargs.Opt54)
		}
		if optargs.Opt55_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT55_BITMASK
			c_optargs.opt55 = C.int(optargs.Opt55)
		}
		if optargs.Opt56_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT56_BITMASK
			c_optargs.opt56 = C.int(optargs.Opt56)
		}
		if optargs.Opt57_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT57_BITMASK
			c_optargs.opt57 = C.int(optargs.Opt57)
		}
		if optargs.Opt58_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT58_BITMASK
			c_optargs.opt58 = C.int(optargs.Opt58)
		}
		if optargs.Opt59_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT59_BITMASK
			c_optargs.opt59 = C.int(optargs.Opt59)
		}
		if optargs.Opt60_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT60_BITMASK
			c_optargs.opt60 = C.int(optargs.Opt60)
		}
		if optargs.Opt61_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT61_BITMASK
			c_optargs.opt61 = C.int(optargs.Opt61)
		}
		if optargs.Opt62_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT62_BITMASK
			c_optargs.opt62 = C.int(optargs.Opt62)
		}
		if optargs.Opt63_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_63_OPTARGS_OPT63_BITMASK
			c_optargs.opt63 = C.int(optargs.Opt63)
		}
	}

	r := C.guestfs_internal_test_63_optargs_argv(g.g, &c_optargs)

	if r == -1 {
		return get_error_from_handle(g, "internal_test_63_optargs")
	}
	return nil
}

/* internal_test_close_output : internal test function - do not use */
func (g *Guestfs) Internal_test_close_output() *GuestfsError {
	if g.g == nil {
		return closed_handle_error("internal_test_close_output")
	}

	r := C.guestfs_internal_test_close_output(g.g)

	if r == -1 {
		return get_error_from_handle(g, "internal_test_close_output")
	}
	return nil
}

/* Struct carrying optional arguments for Internal_test_only_optargs */
type OptargsInternal_test_only_optargs struct {
	/* Test field is ignored unless Test_is_set == true */
	Test_is_set bool
	Test        int
}

/* internal_test_only_optargs : internal test function - do not use */
func (g *Guestfs) Internal_test_only_optargs(optargs *OptargsInternal_test_only_optargs) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("internal_test_only_optargs")
	}
	c_optargs := C.struct_guestfs_internal_test_only_optargs_argv{}
	if optargs != nil {
		if optargs.Test_is_set {
			c_optargs.bitmask |= C.GUESTFS_INTERNAL_TEST_ONLY_OPTARGS_TEST_BITMASK
			c_optargs.test = C.int(optargs.Test)
		}
	}

	r := C.guestfs_internal_test_only_optargs_argv(g.g, &c_optargs)

	if r == -1 {
		return get_error_from_handle(g, "internal_test_only_optargs")
	}
	return nil
}

/* internal_test_rbool : internal test function - do not use */
func (g *Guestfs) Internal_test_rbool(val string) (bool, *GuestfsError) {
	if g.g == nil {
		return false, closed_handle_error("internal_test_rbool")
	}

	c_val := C.CString(val)
	defer C.free(unsafe.Pointer(c_val))

	r := C.guestfs_internal_test_rbool(g.g, c_val)

	if r == -1 {
		return false, get_error_from_handle(g, "internal_test_rbool")
	}
	return r != 0, nil
}

/* internal_test_rboolerr : internal test function - do not use */
func (g *Guestfs) Internal_test_rboolerr() (bool, *GuestfsError) {
	if g.g == nil {
		return false, closed_handle_error("internal_test_rboolerr")
	}

	r := C.guestfs_internal_test_rboolerr(g.g)

	if r == -1 {
		return false, get_error_from_handle(g, "internal_test_rboolerr")
	}
	return r != 0, nil
}

/* internal_test_rbufferout : internal test function - do not use */
func (g *Guestfs) Internal_test_rbufferout(val string) ([]byte, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("internal_test_rbufferout")
	}

	c_val := C.CString(val)
	defer C.free(unsafe.Pointer(c_val))

	var size C.size_t

	r := C.guestfs_internal_test_rbufferout(g.g, c_val, &size)

	if r == nil {
		return nil, get_error_from_handle(g, "internal_test_rbufferout")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoBytes(unsafe.Pointer(r), C.int(size)), nil
}

/* internal_test_rbufferouterr : internal test function - do not use */
func (g *Guestfs) Internal_test_rbufferouterr() ([]byte, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("internal_test_rbufferouterr")
	}

	var size C.size_t

	r := C.guestfs_internal_test_rbufferouterr(g.g, &size)

	if r == nil {
		return nil, get_error_from_handle(g, "internal_test_rbufferouterr")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoBytes(unsafe.Pointer(r), C.int(size)), nil
}

/* internal_test_rconstoptstring : internal test function - do not use */
func (g *Guestfs) Internal_test_rconstoptstring(val string) (*string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("internal_test_rconstoptstring")
	}

	c_val := C.CString(val)
	defer C.free(unsafe.Pointer(c_val))

	r := C.guestfs_internal_test_rconstoptstring(g.g, c_val)
	if r != nil {
		r_s := string(*r)
		return &r_s, nil
	} else {
		return nil, nil
	}
}

/* internal_test_rconstoptstringerr : internal test function - do not use */
func (g *Guestfs) Internal_test_rconstoptstringerr() (*string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("internal_test_rconstoptstringerr")
	}

	r := C.guestfs_internal_test_rconstoptstringerr(g.g)
	if r != nil {
		r_s := string(*r)
		return &r_s, nil
	} else {
		return nil, nil
	}
}

/* internal_test_rconststring : internal test function - do not use */
func (g *Guestfs) Internal_test_rconststring(val string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("internal_test_rconststring")
	}

	c_val := C.CString(val)
	defer C.free(unsafe.Pointer(c_val))

	r := C.guestfs_internal_test_rconststring(g.g, c_val)

	if r == nil {
		return "", get_error_from_handle(g, "internal_test_rconststring")
	}
	return C.GoString(r), nil
}

/* internal_test_rconststringerr : internal test function - do not use */
func (g *Guestfs) Internal_test_rconststringerr() (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("internal_test_rconststringerr")
	}

	r := C.guestfs_internal_test_rconststringerr(g.g)

	if r == nil {
		return "", get_error_from_handle(g, "internal_test_rconststringerr")
	}
	return C.GoString(r), nil
}

/* internal_test_rhashtable : internal test function - do not use */
func (g *Guestfs) Internal_test_rhashtable(val string) (map[string]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("internal_test_rhashtable")
	}

	c_val := C.CString(val)
	defer C.free(unsafe.Pointer(c_val))

	r := C.guestfs_internal_test_rhashtable(g.g, c_val)

	if r == nil {
		return nil, get_error_from_handle(g, "internal_test_rhashtable")
	}
	defer free_string_list(r)
	return return_hashtable(r), nil
}

/* internal_test_rhashtableerr : internal test function - do not use */
func (g *Guestfs) Internal_test_rhashtableerr() (map[string]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("internal_test_rhashtableerr")
	}

	r := C.guestfs_internal_test_rhashtableerr(g.g)

	if r == nil {
		return nil, get_error_from_handle(g, "internal_test_rhashtableerr")
	}
	defer free_string_list(r)
	return return_hashtable(r), nil
}

/* internal_test_rint : internal test function - do not use */
func (g *Guestfs) Internal_test_rint(val string) (int, *GuestfsError) {
	if g.g == nil {
		return 0, closed_handle_error("internal_test_rint")
	}

	c_val := C.CString(val)
	defer C.free(unsafe.Pointer(c_val))

	r := C.guestfs_internal_test_rint(g.g, c_val)

	if r == -1 {
		return 0, get_error_from_handle(g, "internal_test_rint")
	}
	return int(r), nil
}

/* internal_test_rint64 : internal test function - do not use */
func (g *Guestfs) Internal_test_rint64(val string) (int64, *GuestfsError) {
	if g.g == nil {
		return 0, closed_handle_error("internal_test_rint64")
	}

	c_val := C.CString(val)
	defer C.free(unsafe.Pointer(c_val))

	r := C.guestfs_internal_test_rint64(g.g, c_val)

	if r == -1 {
		return 0, get_error_from_handle(g, "internal_test_rint64")
	}
	return int64(r), nil
}

/* internal_test_rint64err : internal test function - do not use */
func (g *Guestfs) Internal_test_rint64err() (int64, *GuestfsError) {
	if g.g == nil {
		return 0, closed_handle_error("internal_test_rint64err")
	}

	r := C.guestfs_internal_test_rint64err(g.g)

	if r == -1 {
		return 0, get_error_from_handle(g, "internal_test_rint64err")
	}
	return int64(r), nil
}

/* internal_test_rinterr : internal test function - do not use */
func (g *Guestfs) Internal_test_rinterr() (int, *GuestfsError) {
	if g.g == nil {
		return 0, closed_handle_error("internal_test_rinterr")
	}

	r := C.guestfs_internal_test_rinterr(g.g)

	if r == -1 {
		return 0, get_error_from_handle(g, "internal_test_rinterr")
	}
	return int(r), nil
}

/* internal_test_rstring : internal test function - do not use */
func (g *Guestfs) Internal_test_rstring(val string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("internal_test_rstring")
	}

	c_val := C.CString(val)
	defer C.free(unsafe.Pointer(c_val))

	r := C.guestfs_internal_test_rstring(g.g, c_val)

	if r == nil {
		return "", get_error_from_handle(g, "internal_test_rstring")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* internal_test_rstringerr : internal test function - do not use */
func (g *Guestfs) Internal_test_rstringerr() (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("internal_test_rstringerr")
	}

	r := C.guestfs_internal_test_rstringerr(g.g)

	if r == nil {
		return "", get_error_from_handle(g, "internal_test_rstringerr")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* internal_test_rstringlist : internal test function - do not use */
func (g *Guestfs) Internal_test_rstringlist(val string) ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("internal_test_rstringlist")
	}

	c_val := C.CString(val)
	defer C.free(unsafe.Pointer(c_val))

	r := C.guestfs_internal_test_rstringlist(g.g, c_val)

	if r == nil {
		return nil, get_error_from_handle(g, "internal_test_rstringlist")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* internal_test_rstringlisterr : internal test function - do not use */
func (g *Guestfs) Internal_test_rstringlisterr() ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("internal_test_rstringlisterr")
	}

	r := C.guestfs_internal_test_rstringlisterr(g.g)

	if r == nil {
		return nil, get_error_from_handle(g, "internal_test_rstringlisterr")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* internal_test_rstruct : internal test function - do not use */
func (g *Guestfs) Internal_test_rstruct(val string) (*PV, *GuestfsError) {
	if g.g == nil {
		return &PV{}, closed_handle_error("internal_test_rstruct")
	}

	c_val := C.CString(val)
	defer C.free(unsafe.Pointer(c_val))

	r := C.guestfs_internal_test_rstruct(g.g, c_val)

	if r == nil {
		return &PV{}, get_error_from_handle(g, "internal_test_rstruct")
	}
	defer C.guestfs_free_lvm_pv(r)
	return return_PV(r), nil
}

/* internal_test_rstructerr : internal test function - do not use */
func (g *Guestfs) Internal_test_rstructerr() (*PV, *GuestfsError) {
	if g.g == nil {
		return &PV{}, closed_handle_error("internal_test_rstructerr")
	}

	r := C.guestfs_internal_test_rstructerr(g.g)

	if r == nil {
		return &PV{}, get_error_from_handle(g, "internal_test_rstructerr")
	}
	defer C.guestfs_free_lvm_pv(r)
	return return_PV(r), nil
}

/* internal_test_rstructlist : internal test function - do not use */
func (g *Guestfs) Internal_test_rstructlist(val string) (*[]PV, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("internal_test_rstructlist")
	}

	c_val := C.CString(val)
	defer C.free(unsafe.Pointer(c_val))

	r := C.guestfs_internal_test_rstructlist(g.g, c_val)

	if r == nil {
		return nil, get_error_from_handle(g, "internal_test_rstructlist")
	}
	defer C.guestfs_free_lvm_pv_list(r)
	return return_PV_list(r), nil
}

/* internal_test_rstructlisterr : internal test function - do not use */
func (g *Guestfs) Internal_test_rstructlisterr() (*[]PV, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("internal_test_rstructlisterr")
	}

	r := C.guestfs_internal_test_rstructlisterr(g.g)

	if r == nil {
		return nil, get_error_from_handle(g, "internal_test_rstructlisterr")
	}
	defer C.guestfs_free_lvm_pv_list(r)
	return return_PV_list(r), nil
}

/* internal_test_set_output : internal test function - do not use */
func (g *Guestfs) Internal_test_set_output(filename string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("internal_test_set_output")
	}

	c_filename := C.CString(filename)
	defer C.free(unsafe.Pointer(c_filename))

	r := C.guestfs_internal_test_set_output(g.g, c_filename)

	if r == -1 {
		return get_error_from_handle(g, "internal_test_set_output")
	}
	return nil
}

/* Struct carrying optional arguments for Is_blockdev */
type OptargsIs_blockdev struct {
	/* Followsymlinks field is ignored unless Followsymlinks_is_set == true */
	Followsymlinks_is_set bool
	Followsymlinks        bool
}

/* is_blockdev : test if block device */
func (g *Guestfs) Is_blockdev(path string, optargs *OptargsIs_blockdev) (bool, *GuestfsError) {
	if g.g == nil {
		return false, closed_handle_error("is_blockdev")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))
	c_optargs := C.struct_guestfs_is_blockdev_opts_argv{}
	if optargs != nil {
		if optargs.Followsymlinks_is_set {
			c_optargs.bitmask |= C.GUESTFS_IS_BLOCKDEV_OPTS_FOLLOWSYMLINKS_BITMASK
			if optargs.Followsymlinks {
				c_optargs.followsymlinks = 1
			} else {
				c_optargs.followsymlinks = 0
			}
		}
	}

	r := C.guestfs_is_blockdev_opts_argv(g.g, c_path, &c_optargs)

	if r == -1 {
		return false, get_error_from_handle(g, "is_blockdev")
	}
	return r != 0, nil
}

/* is_busy : is busy processing a command */
func (g *Guestfs) Is_busy() (bool, *GuestfsError) {
	if g.g == nil {
		return false, closed_handle_error("is_busy")
	}

	r := C.guestfs_is_busy(g.g)

	if r == -1 {
		return false, get_error_from_handle(g, "is_busy")
	}
	return r != 0, nil
}

/* Struct carrying optional arguments for Is_chardev */
type OptargsIs_chardev struct {
	/* Followsymlinks field is ignored unless Followsymlinks_is_set == true */
	Followsymlinks_is_set bool
	Followsymlinks        bool
}

/* is_chardev : test if character device */
func (g *Guestfs) Is_chardev(path string, optargs *OptargsIs_chardev) (bool, *GuestfsError) {
	if g.g == nil {
		return false, closed_handle_error("is_chardev")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))
	c_optargs := C.struct_guestfs_is_chardev_opts_argv{}
	if optargs != nil {
		if optargs.Followsymlinks_is_set {
			c_optargs.bitmask |= C.GUESTFS_IS_CHARDEV_OPTS_FOLLOWSYMLINKS_BITMASK
			if optargs.Followsymlinks {
				c_optargs.followsymlinks = 1
			} else {
				c_optargs.followsymlinks = 0
			}
		}
	}

	r := C.guestfs_is_chardev_opts_argv(g.g, c_path, &c_optargs)

	if r == -1 {
		return false, get_error_from_handle(g, "is_chardev")
	}
	return r != 0, nil
}

/* is_config : is in configuration state */
func (g *Guestfs) Is_config() (bool, *GuestfsError) {
	if g.g == nil {
		return false, closed_handle_error("is_config")
	}

	r := C.guestfs_is_config(g.g)

	if r == -1 {
		return false, get_error_from_handle(g, "is_config")
	}
	return r != 0, nil
}

/* Struct carrying optional arguments for Is_dir */
type OptargsIs_dir struct {
	/* Followsymlinks field is ignored unless Followsymlinks_is_set == true */
	Followsymlinks_is_set bool
	Followsymlinks        bool
}

/* is_dir : test if a directory */
func (g *Guestfs) Is_dir(path string, optargs *OptargsIs_dir) (bool, *GuestfsError) {
	if g.g == nil {
		return false, closed_handle_error("is_dir")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))
	c_optargs := C.struct_guestfs_is_dir_opts_argv{}
	if optargs != nil {
		if optargs.Followsymlinks_is_set {
			c_optargs.bitmask |= C.GUESTFS_IS_DIR_OPTS_FOLLOWSYMLINKS_BITMASK
			if optargs.Followsymlinks {
				c_optargs.followsymlinks = 1
			} else {
				c_optargs.followsymlinks = 0
			}
		}
	}

	r := C.guestfs_is_dir_opts_argv(g.g, c_path, &c_optargs)

	if r == -1 {
		return false, get_error_from_handle(g, "is_dir")
	}
	return r != 0, nil
}

/* Struct carrying optional arguments for Is_fifo */
type OptargsIs_fifo struct {
	/* Followsymlinks field is ignored unless Followsymlinks_is_set == true */
	Followsymlinks_is_set bool
	Followsymlinks        bool
}

/* is_fifo : test if FIFO (named pipe) */
func (g *Guestfs) Is_fifo(path string, optargs *OptargsIs_fifo) (bool, *GuestfsError) {
	if g.g == nil {
		return false, closed_handle_error("is_fifo")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))
	c_optargs := C.struct_guestfs_is_fifo_opts_argv{}
	if optargs != nil {
		if optargs.Followsymlinks_is_set {
			c_optargs.bitmask |= C.GUESTFS_IS_FIFO_OPTS_FOLLOWSYMLINKS_BITMASK
			if optargs.Followsymlinks {
				c_optargs.followsymlinks = 1
			} else {
				c_optargs.followsymlinks = 0
			}
		}
	}

	r := C.guestfs_is_fifo_opts_argv(g.g, c_path, &c_optargs)

	if r == -1 {
		return false, get_error_from_handle(g, "is_fifo")
	}
	return r != 0, nil
}

/* Struct carrying optional arguments for Is_file */
type OptargsIs_file struct {
	/* Followsymlinks field is ignored unless Followsymlinks_is_set == true */
	Followsymlinks_is_set bool
	Followsymlinks        bool
}

/* is_file : test if a regular file */
func (g *Guestfs) Is_file(path string, optargs *OptargsIs_file) (bool, *GuestfsError) {
	if g.g == nil {
		return false, closed_handle_error("is_file")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))
	c_optargs := C.struct_guestfs_is_file_opts_argv{}
	if optargs != nil {
		if optargs.Followsymlinks_is_set {
			c_optargs.bitmask |= C.GUESTFS_IS_FILE_OPTS_FOLLOWSYMLINKS_BITMASK
			if optargs.Followsymlinks {
				c_optargs.followsymlinks = 1
			} else {
				c_optargs.followsymlinks = 0
			}
		}
	}

	r := C.guestfs_is_file_opts_argv(g.g, c_path, &c_optargs)

	if r == -1 {
		return false, get_error_from_handle(g, "is_file")
	}
	return r != 0, nil
}

/* is_launching : is launching subprocess */
func (g *Guestfs) Is_launching() (bool, *GuestfsError) {
	if g.g == nil {
		return false, closed_handle_error("is_launching")
	}

	r := C.guestfs_is_launching(g.g)

	if r == -1 {
		return false, get_error_from_handle(g, "is_launching")
	}
	return r != 0, nil
}

/* is_lv : test if device is a logical volume */
func (g *Guestfs) Is_lv(device string) (bool, *GuestfsError) {
	if g.g == nil {
		return false, closed_handle_error("is_lv")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_is_lv(g.g, c_device)

	if r == -1 {
		return false, get_error_from_handle(g, "is_lv")
	}
	return r != 0, nil
}

/* is_ready : is ready to accept commands */
func (g *Guestfs) Is_ready() (bool, *GuestfsError) {
	if g.g == nil {
		return false, closed_handle_error("is_ready")
	}

	r := C.guestfs_is_ready(g.g)

	if r == -1 {
		return false, get_error_from_handle(g, "is_ready")
	}
	return r != 0, nil
}

/* Struct carrying optional arguments for Is_socket */
type OptargsIs_socket struct {
	/* Followsymlinks field is ignored unless Followsymlinks_is_set == true */
	Followsymlinks_is_set bool
	Followsymlinks        bool
}

/* is_socket : test if socket */
func (g *Guestfs) Is_socket(path string, optargs *OptargsIs_socket) (bool, *GuestfsError) {
	if g.g == nil {
		return false, closed_handle_error("is_socket")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))
	c_optargs := C.struct_guestfs_is_socket_opts_argv{}
	if optargs != nil {
		if optargs.Followsymlinks_is_set {
			c_optargs.bitmask |= C.GUESTFS_IS_SOCKET_OPTS_FOLLOWSYMLINKS_BITMASK
			if optargs.Followsymlinks {
				c_optargs.followsymlinks = 1
			} else {
				c_optargs.followsymlinks = 0
			}
		}
	}

	r := C.guestfs_is_socket_opts_argv(g.g, c_path, &c_optargs)

	if r == -1 {
		return false, get_error_from_handle(g, "is_socket")
	}
	return r != 0, nil
}

/* is_symlink : test if symbolic link */
func (g *Guestfs) Is_symlink(path string) (bool, *GuestfsError) {
	if g.g == nil {
		return false, closed_handle_error("is_symlink")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_is_symlink(g.g, c_path)

	if r == -1 {
		return false, get_error_from_handle(g, "is_symlink")
	}
	return r != 0, nil
}

/* is_whole_device : test if a device is a whole device */
func (g *Guestfs) Is_whole_device(device string) (bool, *GuestfsError) {
	if g.g == nil {
		return false, closed_handle_error("is_whole_device")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_is_whole_device(g.g, c_device)

	if r == -1 {
		return false, get_error_from_handle(g, "is_whole_device")
	}
	return r != 0, nil
}

/* is_zero : test if a file contains all zero bytes */
func (g *Guestfs) Is_zero(path string) (bool, *GuestfsError) {
	if g.g == nil {
		return false, closed_handle_error("is_zero")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_is_zero(g.g, c_path)

	if r == -1 {
		return false, get_error_from_handle(g, "is_zero")
	}
	return r != 0, nil
}

/* is_zero_device : test if a device contains all zero bytes */
func (g *Guestfs) Is_zero_device(device string) (bool, *GuestfsError) {
	if g.g == nil {
		return false, closed_handle_error("is_zero_device")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_is_zero_device(g.g, c_device)

	if r == -1 {
		return false, get_error_from_handle(g, "is_zero_device")
	}
	return r != 0, nil
}

/* isoinfo : get ISO information from primary volume descriptor of ISO file */
func (g *Guestfs) Isoinfo(isofile string) (*ISOInfo, *GuestfsError) {
	if g.g == nil {
		return &ISOInfo{}, closed_handle_error("isoinfo")
	}

	c_isofile := C.CString(isofile)
	defer C.free(unsafe.Pointer(c_isofile))

	r := C.guestfs_isoinfo(g.g, c_isofile)

	if r == nil {
		return &ISOInfo{}, get_error_from_handle(g, "isoinfo")
	}
	defer C.guestfs_free_isoinfo(r)
	return return_ISOInfo(r), nil
}

/* isoinfo_device : get ISO information from primary volume descriptor of device */
func (g *Guestfs) Isoinfo_device(device string) (*ISOInfo, *GuestfsError) {
	if g.g == nil {
		return &ISOInfo{}, closed_handle_error("isoinfo_device")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_isoinfo_device(g.g, c_device)

	if r == nil {
		return &ISOInfo{}, get_error_from_handle(g, "isoinfo_device")
	}
	defer C.guestfs_free_isoinfo(r)
	return return_ISOInfo(r), nil
}

/* journal_close : close the systemd journal */
func (g *Guestfs) Journal_close() *GuestfsError {
	if g.g == nil {
		return closed_handle_error("journal_close")
	}

	r := C.guestfs_journal_close(g.g)

	if r == -1 {
		return get_error_from_handle(g, "journal_close")
	}
	return nil
}

/* journal_get : read the current journal entry */
func (g *Guestfs) Journal_get() (*[]XAttr, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("journal_get")
	}

	r := C.guestfs_journal_get(g.g)

	if r == nil {
		return nil, get_error_from_handle(g, "journal_get")
	}
	defer C.guestfs_free_xattr_list(r)
	return return_XAttr_list(r), nil
}

/* journal_get_data_threshold : get the data threshold for reading journal entries */
func (g *Guestfs) Journal_get_data_threshold() (int64, *GuestfsError) {
	if g.g == nil {
		return 0, closed_handle_error("journal_get_data_threshold")
	}

	r := C.guestfs_journal_get_data_threshold(g.g)

	if r == -1 {
		return 0, get_error_from_handle(g, "journal_get_data_threshold")
	}
	return int64(r), nil
}

/* journal_get_realtime_usec : get the timestamp of the current journal entry */
func (g *Guestfs) Journal_get_realtime_usec() (int64, *GuestfsError) {
	if g.g == nil {
		return 0, closed_handle_error("journal_get_realtime_usec")
	}

	r := C.guestfs_journal_get_realtime_usec(g.g)

	if r == -1 {
		return 0, get_error_from_handle(g, "journal_get_realtime_usec")
	}
	return int64(r), nil
}

/* journal_next : move to the next journal entry */
func (g *Guestfs) Journal_next() (bool, *GuestfsError) {
	if g.g == nil {
		return false, closed_handle_error("journal_next")
	}

	r := C.guestfs_journal_next(g.g)

	if r == -1 {
		return false, get_error_from_handle(g, "journal_next")
	}
	return r != 0, nil
}

/* journal_open : open the systemd journal */
func (g *Guestfs) Journal_open(directory string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("journal_open")
	}

	c_directory := C.CString(directory)
	defer C.free(unsafe.Pointer(c_directory))

	r := C.guestfs_journal_open(g.g, c_directory)

	if r == -1 {
		return get_error_from_handle(g, "journal_open")
	}
	return nil
}

/* journal_set_data_threshold : set the data threshold for reading journal entries */
func (g *Guestfs) Journal_set_data_threshold(threshold int64) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("journal_set_data_threshold")
	}

	r := C.guestfs_journal_set_data_threshold(g.g, C.int64_t(threshold))

	if r == -1 {
		return get_error_from_handle(g, "journal_set_data_threshold")
	}
	return nil
}

/* journal_skip : skip forwards or backwards in the journal */
func (g *Guestfs) Journal_skip(skip int64) (int64, *GuestfsError) {
	if g.g == nil {
		return 0, closed_handle_error("journal_skip")
	}

	r := C.guestfs_journal_skip(g.g, C.int64_t(skip))

	if r == -1 {
		return 0, get_error_from_handle(g, "journal_skip")
	}
	return int64(r), nil
}

/* kill_subprocess : kill the hypervisor */
func (g *Guestfs) Kill_subprocess() *GuestfsError {
	if g.g == nil {
		return closed_handle_error("kill_subprocess")
	}

	r := C.guestfs_kill_subprocess(g.g)

	if r == -1 {
		return get_error_from_handle(g, "kill_subprocess")
	}
	return nil
}

/* launch : launch the backend */
func (g *Guestfs) Launch() *GuestfsError {
	if g.g == nil {
		return closed_handle_error("launch")
	}

	r := C.guestfs_launch(g.g)

	if r == -1 {
		return get_error_from_handle(g, "launch")
	}
	return nil
}

/* lchown : change file owner and group */
func (g *Guestfs) Lchown(owner int, group int, path string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("lchown")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_lchown(g.g, C.int(owner), C.int(group), c_path)

	if r == -1 {
		return get_error_from_handle(g, "lchown")
	}
	return nil
}

/* ldmtool_create_all : scan and create Windows dynamic disk volumes */
func (g *Guestfs) Ldmtool_create_all() *GuestfsError {
	if g.g == nil {
		return closed_handle_error("ldmtool_create_all")
	}

	r := C.guestfs_ldmtool_create_all(g.g)

	if r == -1 {
		return get_error_from_handle(g, "ldmtool_create_all")
	}
	return nil
}

/* ldmtool_diskgroup_disks : return the disks in a Windows dynamic disk group */
func (g *Guestfs) Ldmtool_diskgroup_disks(diskgroup string) ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("ldmtool_diskgroup_disks")
	}

	c_diskgroup := C.CString(diskgroup)
	defer C.free(unsafe.Pointer(c_diskgroup))

	r := C.guestfs_ldmtool_diskgroup_disks(g.g, c_diskgroup)

	if r == nil {
		return nil, get_error_from_handle(g, "ldmtool_diskgroup_disks")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* ldmtool_diskgroup_name : return the name of a Windows dynamic disk group */
func (g *Guestfs) Ldmtool_diskgroup_name(diskgroup string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("ldmtool_diskgroup_name")
	}

	c_diskgroup := C.CString(diskgroup)
	defer C.free(unsafe.Pointer(c_diskgroup))

	r := C.guestfs_ldmtool_diskgroup_name(g.g, c_diskgroup)

	if r == nil {
		return "", get_error_from_handle(g, "ldmtool_diskgroup_name")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* ldmtool_diskgroup_volumes : return the volumes in a Windows dynamic disk group */
func (g *Guestfs) Ldmtool_diskgroup_volumes(diskgroup string) ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("ldmtool_diskgroup_volumes")
	}

	c_diskgroup := C.CString(diskgroup)
	defer C.free(unsafe.Pointer(c_diskgroup))

	r := C.guestfs_ldmtool_diskgroup_volumes(g.g, c_diskgroup)

	if r == nil {
		return nil, get_error_from_handle(g, "ldmtool_diskgroup_volumes")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* ldmtool_remove_all : remove all Windows dynamic disk volumes */
func (g *Guestfs) Ldmtool_remove_all() *GuestfsError {
	if g.g == nil {
		return closed_handle_error("ldmtool_remove_all")
	}

	r := C.guestfs_ldmtool_remove_all(g.g)

	if r == -1 {
		return get_error_from_handle(g, "ldmtool_remove_all")
	}
	return nil
}

/* ldmtool_scan : scan for Windows dynamic disks */
func (g *Guestfs) Ldmtool_scan() ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("ldmtool_scan")
	}

	r := C.guestfs_ldmtool_scan(g.g)

	if r == nil {
		return nil, get_error_from_handle(g, "ldmtool_scan")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* ldmtool_scan_devices : scan for Windows dynamic disks */
func (g *Guestfs) Ldmtool_scan_devices(devices []string) ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("ldmtool_scan_devices")
	}

	c_devices := arg_string_list(devices)
	defer free_string_list(c_devices)

	r := C.guestfs_ldmtool_scan_devices(g.g, c_devices)

	if r == nil {
		return nil, get_error_from_handle(g, "ldmtool_scan_devices")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* ldmtool_volume_hint : return the hint field of a Windows dynamic disk volume */
func (g *Guestfs) Ldmtool_volume_hint(diskgroup string, volume string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("ldmtool_volume_hint")
	}

	c_diskgroup := C.CString(diskgroup)
	defer C.free(unsafe.Pointer(c_diskgroup))

	c_volume := C.CString(volume)
	defer C.free(unsafe.Pointer(c_volume))

	r := C.guestfs_ldmtool_volume_hint(g.g, c_diskgroup, c_volume)

	if r == nil {
		return "", get_error_from_handle(g, "ldmtool_volume_hint")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* ldmtool_volume_partitions : return the partitions in a Windows dynamic disk volume */
func (g *Guestfs) Ldmtool_volume_partitions(diskgroup string, volume string) ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("ldmtool_volume_partitions")
	}

	c_diskgroup := C.CString(diskgroup)
	defer C.free(unsafe.Pointer(c_diskgroup))

	c_volume := C.CString(volume)
	defer C.free(unsafe.Pointer(c_volume))

	r := C.guestfs_ldmtool_volume_partitions(g.g, c_diskgroup, c_volume)

	if r == nil {
		return nil, get_error_from_handle(g, "ldmtool_volume_partitions")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* ldmtool_volume_type : return the type of a Windows dynamic disk volume */
func (g *Guestfs) Ldmtool_volume_type(diskgroup string, volume string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("ldmtool_volume_type")
	}

	c_diskgroup := C.CString(diskgroup)
	defer C.free(unsafe.Pointer(c_diskgroup))

	c_volume := C.CString(volume)
	defer C.free(unsafe.Pointer(c_volume))

	r := C.guestfs_ldmtool_volume_type(g.g, c_diskgroup, c_volume)

	if r == nil {
		return "", get_error_from_handle(g, "ldmtool_volume_type")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* lgetxattr : get a single extended attribute */
func (g *Guestfs) Lgetxattr(path string, name string) ([]byte, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("lgetxattr")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	c_name := C.CString(name)
	defer C.free(unsafe.Pointer(c_name))

	var size C.size_t

	r := C.guestfs_lgetxattr(g.g, c_path, c_name, &size)

	if r == nil {
		return nil, get_error_from_handle(g, "lgetxattr")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoBytes(unsafe.Pointer(r), C.int(size)), nil
}

/* lgetxattrs : list extended attributes of a file or directory */
func (g *Guestfs) Lgetxattrs(path string) (*[]XAttr, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("lgetxattrs")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_lgetxattrs(g.g, c_path)

	if r == nil {
		return nil, get_error_from_handle(g, "lgetxattrs")
	}
	defer C.guestfs_free_xattr_list(r)
	return return_XAttr_list(r), nil
}

/* list_9p : list 9p filesystems */
func (g *Guestfs) List_9p() ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("list_9p")
	}

	r := C.guestfs_list_9p(g.g)

	if r == nil {
		return nil, get_error_from_handle(g, "list_9p")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* list_devices : list the block devices */
func (g *Guestfs) List_devices() ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("list_devices")
	}

	r := C.guestfs_list_devices(g.g)

	if r == nil {
		return nil, get_error_from_handle(g, "list_devices")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* list_disk_labels : mapping of disk labels to devices */
func (g *Guestfs) List_disk_labels() (map[string]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("list_disk_labels")
	}

	r := C.guestfs_list_disk_labels(g.g)

	if r == nil {
		return nil, get_error_from_handle(g, "list_disk_labels")
	}
	defer free_string_list(r)
	return return_hashtable(r), nil
}

/* list_dm_devices : list device mapper devices */
func (g *Guestfs) List_dm_devices() ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("list_dm_devices")
	}

	r := C.guestfs_list_dm_devices(g.g)

	if r == nil {
		return nil, get_error_from_handle(g, "list_dm_devices")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* list_filesystems : list filesystems */
func (g *Guestfs) List_filesystems() (map[string]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("list_filesystems")
	}

	r := C.guestfs_list_filesystems(g.g)

	if r == nil {
		return nil, get_error_from_handle(g, "list_filesystems")
	}
	defer free_string_list(r)
	return return_hashtable(r), nil
}

/* list_ldm_partitions : list all Windows dynamic disk partitions */
func (g *Guestfs) List_ldm_partitions() ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("list_ldm_partitions")
	}

	r := C.guestfs_list_ldm_partitions(g.g)

	if r == nil {
		return nil, get_error_from_handle(g, "list_ldm_partitions")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* list_ldm_volumes : list all Windows dynamic disk volumes */
func (g *Guestfs) List_ldm_volumes() ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("list_ldm_volumes")
	}

	r := C.guestfs_list_ldm_volumes(g.g)

	if r == nil {
		return nil, get_error_from_handle(g, "list_ldm_volumes")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* list_md_devices : list Linux md (RAID) devices */
func (g *Guestfs) List_md_devices() ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("list_md_devices")
	}

	r := C.guestfs_list_md_devices(g.g)

	if r == nil {
		return nil, get_error_from_handle(g, "list_md_devices")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* list_partitions : list the partitions */
func (g *Guestfs) List_partitions() ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("list_partitions")
	}

	r := C.guestfs_list_partitions(g.g)

	if r == nil {
		return nil, get_error_from_handle(g, "list_partitions")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* ll : list the files in a directory (long format) */
func (g *Guestfs) Ll(directory string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("ll")
	}

	c_directory := C.CString(directory)
	defer C.free(unsafe.Pointer(c_directory))

	r := C.guestfs_ll(g.g, c_directory)

	if r == nil {
		return "", get_error_from_handle(g, "ll")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* llz : list the files in a directory (long format with SELinux contexts) */
func (g *Guestfs) Llz(directory string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("llz")
	}

	c_directory := C.CString(directory)
	defer C.free(unsafe.Pointer(c_directory))

	r := C.guestfs_llz(g.g, c_directory)

	if r == nil {
		return "", get_error_from_handle(g, "llz")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* ln : create a hard link */
func (g *Guestfs) Ln(target string, linkname string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("ln")
	}

	c_target := C.CString(target)
	defer C.free(unsafe.Pointer(c_target))

	c_linkname := C.CString(linkname)
	defer C.free(unsafe.Pointer(c_linkname))

	r := C.guestfs_ln(g.g, c_target, c_linkname)

	if r == -1 {
		return get_error_from_handle(g, "ln")
	}
	return nil
}

/* ln_f : create a hard link */
func (g *Guestfs) Ln_f(target string, linkname string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("ln_f")
	}

	c_target := C.CString(target)
	defer C.free(unsafe.Pointer(c_target))

	c_linkname := C.CString(linkname)
	defer C.free(unsafe.Pointer(c_linkname))

	r := C.guestfs_ln_f(g.g, c_target, c_linkname)

	if r == -1 {
		return get_error_from_handle(g, "ln_f")
	}
	return nil
}

/* ln_s : create a symbolic link */
func (g *Guestfs) Ln_s(target string, linkname string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("ln_s")
	}

	c_target := C.CString(target)
	defer C.free(unsafe.Pointer(c_target))

	c_linkname := C.CString(linkname)
	defer C.free(unsafe.Pointer(c_linkname))

	r := C.guestfs_ln_s(g.g, c_target, c_linkname)

	if r == -1 {
		return get_error_from_handle(g, "ln_s")
	}
	return nil
}

/* ln_sf : create a symbolic link */
func (g *Guestfs) Ln_sf(target string, linkname string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("ln_sf")
	}

	c_target := C.CString(target)
	defer C.free(unsafe.Pointer(c_target))

	c_linkname := C.CString(linkname)
	defer C.free(unsafe.Pointer(c_linkname))

	r := C.guestfs_ln_sf(g.g, c_target, c_linkname)

	if r == -1 {
		return get_error_from_handle(g, "ln_sf")
	}
	return nil
}

/* lremovexattr : remove extended attribute of a file or directory */
func (g *Guestfs) Lremovexattr(xattr string, path string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("lremovexattr")
	}

	c_xattr := C.CString(xattr)
	defer C.free(unsafe.Pointer(c_xattr))

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_lremovexattr(g.g, c_xattr, c_path)

	if r == -1 {
		return get_error_from_handle(g, "lremovexattr")
	}
	return nil
}

/* ls : list the files in a directory */
func (g *Guestfs) Ls(directory string) ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("ls")
	}

	c_directory := C.CString(directory)
	defer C.free(unsafe.Pointer(c_directory))

	r := C.guestfs_ls(g.g, c_directory)

	if r == nil {
		return nil, get_error_from_handle(g, "ls")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* ls0 : get list of files in a directory */
func (g *Guestfs) Ls0(dir string, filenames string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("ls0")
	}

	c_dir := C.CString(dir)
	defer C.free(unsafe.Pointer(c_dir))

	c_filenames := C.CString(filenames)
	defer C.free(unsafe.Pointer(c_filenames))

	r := C.guestfs_ls0(g.g, c_dir, c_filenames)

	if r == -1 {
		return get_error_from_handle(g, "ls0")
	}
	return nil
}

/* lsetxattr : set extended attribute of a file or directory */
func (g *Guestfs) Lsetxattr(xattr string, val string, vallen int, path string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("lsetxattr")
	}

	c_xattr := C.CString(xattr)
	defer C.free(unsafe.Pointer(c_xattr))

	c_val := C.CString(val)
	defer C.free(unsafe.Pointer(c_val))

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_lsetxattr(g.g, c_xattr, c_val, C.int(vallen), c_path)

	if r == -1 {
		return get_error_from_handle(g, "lsetxattr")
	}
	return nil
}

/* lstat : get file information for a symbolic link */
func (g *Guestfs) Lstat(path string) (*Stat, *GuestfsError) {
	if g.g == nil {
		return &Stat{}, closed_handle_error("lstat")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_lstat(g.g, c_path)

	if r == nil {
		return &Stat{}, get_error_from_handle(g, "lstat")
	}
	defer C.guestfs_free_stat(r)
	return return_Stat(r), nil
}

/* lstatlist : lstat on multiple files */
func (g *Guestfs) Lstatlist(path string, names []string) (*[]Stat, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("lstatlist")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	c_names := arg_string_list(names)
	defer free_string_list(c_names)

	r := C.guestfs_lstatlist(g.g, c_path, c_names)

	if r == nil {
		return nil, get_error_from_handle(g, "lstatlist")
	}
	defer C.guestfs_free_stat_list(r)
	return return_Stat_list(r), nil
}

/* lstatns : get file information for a symbolic link */
func (g *Guestfs) Lstatns(path string) (*StatNS, *GuestfsError) {
	if g.g == nil {
		return &StatNS{}, closed_handle_error("lstatns")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_lstatns(g.g, c_path)

	if r == nil {
		return &StatNS{}, get_error_from_handle(g, "lstatns")
	}
	defer C.guestfs_free_statns(r)
	return return_StatNS(r), nil
}

/* lstatnslist : lstat on multiple files */
func (g *Guestfs) Lstatnslist(path string, names []string) (*[]StatNS, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("lstatnslist")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	c_names := arg_string_list(names)
	defer free_string_list(c_names)

	r := C.guestfs_lstatnslist(g.g, c_path, c_names)

	if r == nil {
		return nil, get_error_from_handle(g, "lstatnslist")
	}
	defer C.guestfs_free_statns_list(r)
	return return_StatNS_list(r), nil
}

/* luks_add_key : add a key on a LUKS encrypted device */
func (g *Guestfs) Luks_add_key(device string, key string, newkey string, keyslot int) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("luks_add_key")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	c_key := C.CString(key)
	defer C.free(unsafe.Pointer(c_key))

	c_newkey := C.CString(newkey)
	defer C.free(unsafe.Pointer(c_newkey))

	r := C.guestfs_luks_add_key(g.g, c_device, c_key, c_newkey, C.int(keyslot))

	if r == -1 {
		return get_error_from_handle(g, "luks_add_key")
	}
	return nil
}

/* luks_close : close a LUKS device */
func (g *Guestfs) Luks_close(device string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("luks_close")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_luks_close(g.g, c_device)

	if r == -1 {
		return get_error_from_handle(g, "luks_close")
	}
	return nil
}

/* luks_format : format a block device as a LUKS encrypted device */
func (g *Guestfs) Luks_format(device string, key string, keyslot int) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("luks_format")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	c_key := C.CString(key)
	defer C.free(unsafe.Pointer(c_key))

	r := C.guestfs_luks_format(g.g, c_device, c_key, C.int(keyslot))

	if r == -1 {
		return get_error_from_handle(g, "luks_format")
	}
	return nil
}

/* luks_format_cipher : format a block device as a LUKS encrypted device */
func (g *Guestfs) Luks_format_cipher(device string, key string, keyslot int, cipher string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("luks_format_cipher")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	c_key := C.CString(key)
	defer C.free(unsafe.Pointer(c_key))

	c_cipher := C.CString(cipher)
	defer C.free(unsafe.Pointer(c_cipher))

	r := C.guestfs_luks_format_cipher(g.g, c_device, c_key, C.int(keyslot), c_cipher)

	if r == -1 {
		return get_error_from_handle(g, "luks_format_cipher")
	}
	return nil
}

/* luks_kill_slot : remove a key from a LUKS encrypted device */
func (g *Guestfs) Luks_kill_slot(device string, key string, keyslot int) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("luks_kill_slot")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	c_key := C.CString(key)
	defer C.free(unsafe.Pointer(c_key))

	r := C.guestfs_luks_kill_slot(g.g, c_device, c_key, C.int(keyslot))

	if r == -1 {
		return get_error_from_handle(g, "luks_kill_slot")
	}
	return nil
}

/* luks_open : open a LUKS-encrypted block device */
func (g *Guestfs) Luks_open(device string, key string, mapname string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("luks_open")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	c_key := C.CString(key)
	defer C.free(unsafe.Pointer(c_key))

	c_mapname := C.CString(mapname)
	defer C.free(unsafe.Pointer(c_mapname))

	r := C.guestfs_luks_open(g.g, c_device, c_key, c_mapname)

	if r == -1 {
		return get_error_from_handle(g, "luks_open")
	}
	return nil
}

/* luks_open_ro : open a LUKS-encrypted block device read-only */
func (g *Guestfs) Luks_open_ro(device string, key string, mapname string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("luks_open_ro")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	c_key := C.CString(key)
	defer C.free(unsafe.Pointer(c_key))

	c_mapname := C.CString(mapname)
	defer C.free(unsafe.Pointer(c_mapname))

	r := C.guestfs_luks_open_ro(g.g, c_device, c_key, c_mapname)

	if r == -1 {
		return get_error_from_handle(g, "luks_open_ro")
	}
	return nil
}

/* lvcreate : create an LVM logical volume */
func (g *Guestfs) Lvcreate(logvol string, volgroup string, mbytes int) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("lvcreate")
	}

	c_logvol := C.CString(logvol)
	defer C.free(unsafe.Pointer(c_logvol))

	c_volgroup := C.CString(volgroup)
	defer C.free(unsafe.Pointer(c_volgroup))

	r := C.guestfs_lvcreate(g.g, c_logvol, c_volgroup, C.int(mbytes))

	if r == -1 {
		return get_error_from_handle(g, "lvcreate")
	}
	return nil
}

/* lvcreate_free : create an LVM logical volume in % remaining free space */
func (g *Guestfs) Lvcreate_free(logvol string, volgroup string, percent int) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("lvcreate_free")
	}

	c_logvol := C.CString(logvol)
	defer C.free(unsafe.Pointer(c_logvol))

	c_volgroup := C.CString(volgroup)
	defer C.free(unsafe.Pointer(c_volgroup))

	r := C.guestfs_lvcreate_free(g.g, c_logvol, c_volgroup, C.int(percent))

	if r == -1 {
		return get_error_from_handle(g, "lvcreate_free")
	}
	return nil
}

/* lvm_canonical_lv_name : get canonical name of an LV */
func (g *Guestfs) Lvm_canonical_lv_name(lvname string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("lvm_canonical_lv_name")
	}

	c_lvname := C.CString(lvname)
	defer C.free(unsafe.Pointer(c_lvname))

	r := C.guestfs_lvm_canonical_lv_name(g.g, c_lvname)

	if r == nil {
		return "", get_error_from_handle(g, "lvm_canonical_lv_name")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* lvm_clear_filter : clear LVM device filter */
func (g *Guestfs) Lvm_clear_filter() *GuestfsError {
	if g.g == nil {
		return closed_handle_error("lvm_clear_filter")
	}

	r := C.guestfs_lvm_clear_filter(g.g)

	if r == -1 {
		return get_error_from_handle(g, "lvm_clear_filter")
	}
	return nil
}

/* lvm_remove_all : remove all LVM LVs, VGs and PVs */
func (g *Guestfs) Lvm_remove_all() *GuestfsError {
	if g.g == nil {
		return closed_handle_error("lvm_remove_all")
	}

	r := C.guestfs_lvm_remove_all(g.g)

	if r == -1 {
		return get_error_from_handle(g, "lvm_remove_all")
	}
	return nil
}

/* lvm_set_filter : set LVM device filter */
func (g *Guestfs) Lvm_set_filter(devices []string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("lvm_set_filter")
	}

	c_devices := arg_string_list(devices)
	defer free_string_list(c_devices)

	r := C.guestfs_lvm_set_filter(g.g, c_devices)

	if r == -1 {
		return get_error_from_handle(g, "lvm_set_filter")
	}
	return nil
}

/* lvremove : remove an LVM logical volume */
func (g *Guestfs) Lvremove(device string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("lvremove")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_lvremove(g.g, c_device)

	if r == -1 {
		return get_error_from_handle(g, "lvremove")
	}
	return nil
}

/* lvrename : rename an LVM logical volume */
func (g *Guestfs) Lvrename(logvol string, newlogvol string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("lvrename")
	}

	c_logvol := C.CString(logvol)
	defer C.free(unsafe.Pointer(c_logvol))

	c_newlogvol := C.CString(newlogvol)
	defer C.free(unsafe.Pointer(c_newlogvol))

	r := C.guestfs_lvrename(g.g, c_logvol, c_newlogvol)

	if r == -1 {
		return get_error_from_handle(g, "lvrename")
	}
	return nil
}

/* lvresize : resize an LVM logical volume */
func (g *Guestfs) Lvresize(device string, mbytes int) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("lvresize")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_lvresize(g.g, c_device, C.int(mbytes))

	if r == -1 {
		return get_error_from_handle(g, "lvresize")
	}
	return nil
}

/* lvresize_free : expand an LV to fill free space */
func (g *Guestfs) Lvresize_free(lv string, percent int) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("lvresize_free")
	}

	c_lv := C.CString(lv)
	defer C.free(unsafe.Pointer(c_lv))

	r := C.guestfs_lvresize_free(g.g, c_lv, C.int(percent))

	if r == -1 {
		return get_error_from_handle(g, "lvresize_free")
	}
	return nil
}

/* lvs : list the LVM logical volumes (LVs) */
func (g *Guestfs) Lvs() ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("lvs")
	}

	r := C.guestfs_lvs(g.g)

	if r == nil {
		return nil, get_error_from_handle(g, "lvs")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* lvs_full : list the LVM logical volumes (LVs) */
func (g *Guestfs) Lvs_full() (*[]LV, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("lvs_full")
	}

	r := C.guestfs_lvs_full(g.g)

	if r == nil {
		return nil, get_error_from_handle(g, "lvs_full")
	}
	defer C.guestfs_free_lvm_lv_list(r)
	return return_LV_list(r), nil
}

/* lvuuid : get the UUID of a logical volume */
func (g *Guestfs) Lvuuid(device string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("lvuuid")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_lvuuid(g.g, c_device)

	if r == nil {
		return "", get_error_from_handle(g, "lvuuid")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* lxattrlist : lgetxattr on multiple files */
func (g *Guestfs) Lxattrlist(path string, names []string) (*[]XAttr, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("lxattrlist")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	c_names := arg_string_list(names)
	defer free_string_list(c_names)

	r := C.guestfs_lxattrlist(g.g, c_path, c_names)

	if r == nil {
		return nil, get_error_from_handle(g, "lxattrlist")
	}
	defer C.guestfs_free_xattr_list(r)
	return return_XAttr_list(r), nil
}

/* max_disks : maximum number of disks that may be added */
func (g *Guestfs) Max_disks() (int, *GuestfsError) {
	if g.g == nil {
		return 0, closed_handle_error("max_disks")
	}

	r := C.guestfs_max_disks(g.g)

	if r == -1 {
		return 0, get_error_from_handle(g, "max_disks")
	}
	return int(r), nil
}

/* Struct carrying optional arguments for Md_create */
type OptargsMd_create struct {
	/* Missingbitmap field is ignored unless Missingbitmap_is_set == true */
	Missingbitmap_is_set bool
	Missingbitmap        int64
	/* Nrdevices field is ignored unless Nrdevices_is_set == true */
	Nrdevices_is_set bool
	Nrdevices        int
	/* Spare field is ignored unless Spare_is_set == true */
	Spare_is_set bool
	Spare        int
	/* Chunk field is ignored unless Chunk_is_set == true */
	Chunk_is_set bool
	Chunk        int64
	/* Level field is ignored unless Level_is_set == true */
	Level_is_set bool
	Level        string
}

/* md_create : create a Linux md (RAID) device */
func (g *Guestfs) Md_create(name string, devices []string, optargs *OptargsMd_create) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("md_create")
	}

	c_name := C.CString(name)
	defer C.free(unsafe.Pointer(c_name))

	c_devices := arg_string_list(devices)
	defer free_string_list(c_devices)
	c_optargs := C.struct_guestfs_md_create_argv{}
	if optargs != nil {
		if optargs.Missingbitmap_is_set {
			c_optargs.bitmask |= C.GUESTFS_MD_CREATE_MISSINGBITMAP_BITMASK
			c_optargs.missingbitmap = C.int64_t(optargs.Missingbitmap)
		}
		if optargs.Nrdevices_is_set {
			c_optargs.bitmask |= C.GUESTFS_MD_CREATE_NRDEVICES_BITMASK
			c_optargs.nrdevices = C.int(optargs.Nrdevices)
		}
		if optargs.Spare_is_set {
			c_optargs.bitmask |= C.GUESTFS_MD_CREATE_SPARE_BITMASK
			c_optargs.spare = C.int(optargs.Spare)
		}
		if optargs.Chunk_is_set {
			c_optargs.bitmask |= C.GUESTFS_MD_CREATE_CHUNK_BITMASK
			c_optargs.chunk = C.int64_t(optargs.Chunk)
		}
		if optargs.Level_is_set {
			c_optargs.bitmask |= C.GUESTFS_MD_CREATE_LEVEL_BITMASK
			c_optargs.level = C.CString(optargs.Level)
			defer C.free(unsafe.Pointer(c_optargs.level))
		}
	}

	r := C.guestfs_md_create_argv(g.g, c_name, c_devices, &c_optargs)

	if r == -1 {
		return get_error_from_handle(g, "md_create")
	}
	return nil
}

/* md_detail : obtain metadata for an MD device */
func (g *Guestfs) Md_detail(md string) (map[string]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("md_detail")
	}

	c_md := C.CString(md)
	defer C.free(unsafe.Pointer(c_md))

	r := C.guestfs_md_detail(g.g, c_md)

	if r == nil {
		return nil, get_error_from_handle(g, "md_detail")
	}
	defer free_string_list(r)
	return return_hashtable(r), nil
}

/* md_stat : get underlying devices from an MD device */
func (g *Guestfs) Md_stat(md string) (*[]MDStat, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("md_stat")
	}

	c_md := C.CString(md)
	defer C.free(unsafe.Pointer(c_md))

	r := C.guestfs_md_stat(g.g, c_md)

	if r == nil {
		return nil, get_error_from_handle(g, "md_stat")
	}
	defer C.guestfs_free_mdstat_list(r)
	return return_MDStat_list(r), nil
}

/* md_stop : stop a Linux md (RAID) device */
func (g *Guestfs) Md_stop(md string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("md_stop")
	}

	c_md := C.CString(md)
	defer C.free(unsafe.Pointer(c_md))

	r := C.guestfs_md_stop(g.g, c_md)

	if r == -1 {
		return get_error_from_handle(g, "md_stop")
	}
	return nil
}

/* mkdir : create a directory */
func (g *Guestfs) Mkdir(path string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("mkdir")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_mkdir(g.g, c_path)

	if r == -1 {
		return get_error_from_handle(g, "mkdir")
	}
	return nil
}

/* mkdir_mode : create a directory with a particular mode */
func (g *Guestfs) Mkdir_mode(path string, mode int) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("mkdir_mode")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_mkdir_mode(g.g, c_path, C.int(mode))

	if r == -1 {
		return get_error_from_handle(g, "mkdir_mode")
	}
	return nil
}

/* mkdir_p : create a directory and parents */
func (g *Guestfs) Mkdir_p(path string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("mkdir_p")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_mkdir_p(g.g, c_path)

	if r == -1 {
		return get_error_from_handle(g, "mkdir_p")
	}
	return nil
}

/* mkdtemp : create a temporary directory */
func (g *Guestfs) Mkdtemp(tmpl string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("mkdtemp")
	}

	c_tmpl := C.CString(tmpl)
	defer C.free(unsafe.Pointer(c_tmpl))

	r := C.guestfs_mkdtemp(g.g, c_tmpl)

	if r == nil {
		return "", get_error_from_handle(g, "mkdtemp")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* Struct carrying optional arguments for Mke2fs */
type OptargsMke2fs struct {
	/* Blockscount field is ignored unless Blockscount_is_set == true */
	Blockscount_is_set bool
	Blockscount        int64
	/* Blocksize field is ignored unless Blocksize_is_set == true */
	Blocksize_is_set bool
	Blocksize        int64
	/* Fragsize field is ignored unless Fragsize_is_set == true */
	Fragsize_is_set bool
	Fragsize        int64
	/* Blockspergroup field is ignored unless Blockspergroup_is_set == true */
	Blockspergroup_is_set bool
	Blockspergroup        int64
	/* Numberofgroups field is ignored unless Numberofgroups_is_set == true */
	Numberofgroups_is_set bool
	Numberofgroups        int64
	/* Bytesperinode field is ignored unless Bytesperinode_is_set == true */
	Bytesperinode_is_set bool
	Bytesperinode        int64
	/* Inodesize field is ignored unless Inodesize_is_set == true */
	Inodesize_is_set bool
	Inodesize        int64
	/* Journalsize field is ignored unless Journalsize_is_set == true */
	Journalsize_is_set bool
	Journalsize        int64
	/* Numberofinodes field is ignored unless Numberofinodes_is_set == true */
	Numberofinodes_is_set bool
	Numberofinodes        int64
	/* Stridesize field is ignored unless Stridesize_is_set == true */
	Stridesize_is_set bool
	Stridesize        int64
	/* Stripewidth field is ignored unless Stripewidth_is_set == true */
	Stripewidth_is_set bool
	Stripewidth        int64
	/* Maxonlineresize field is ignored unless Maxonlineresize_is_set == true */
	Maxonlineresize_is_set bool
	Maxonlineresize        int64
	/* Reservedblockspercentage field is ignored unless Reservedblockspercentage_is_set == true */
	Reservedblockspercentage_is_set bool
	Reservedblockspercentage        int
	/* Mmpupdateinterval field is ignored unless Mmpupdateinterval_is_set == true */
	Mmpupdateinterval_is_set bool
	Mmpupdateinterval        int
	/* Journaldevice field is ignored unless Journaldevice_is_set == true */
	Journaldevice_is_set bool
	Journaldevice        string
	/* Label field is ignored unless Label_is_set == true */
	Label_is_set bool
	Label        string
	/* Lastmounteddir field is ignored unless Lastmounteddir_is_set == true */
	Lastmounteddir_is_set bool
	Lastmounteddir        string
	/* Creatoros field is ignored unless Creatoros_is_set == true */
	Creatoros_is_set bool
	Creatoros        string
	/* Fstype field is ignored unless Fstype_is_set == true */
	Fstype_is_set bool
	Fstype        string
	/* Usagetype field is ignored unless Usagetype_is_set == true */
	Usagetype_is_set bool
	Usagetype        string
	/* Uuid field is ignored unless Uuid_is_set == true */
	Uuid_is_set bool
	Uuid        string
	/* Forcecreate field is ignored unless Forcecreate_is_set == true */
	Forcecreate_is_set bool
	Forcecreate        bool
	/* Writesbandgrouponly field is ignored unless Writesbandgrouponly_is_set == true */
	Writesbandgrouponly_is_set bool
	Writesbandgrouponly        bool
	/* Lazyitableinit field is ignored unless Lazyitableinit_is_set == true */
	Lazyitableinit_is_set bool
	Lazyitableinit        bool
	/* Lazyjournalinit field is ignored unless Lazyjournalinit_is_set == true */
	Lazyjournalinit_is_set bool
	Lazyjournalinit        bool
	/* Testfs field is ignored unless Testfs_is_set == true */
	Testfs_is_set bool
	Testfs        bool
	/* Discard field is ignored unless Discard_is_set == true */
	Discard_is_set bool
	Discard        bool
	/* Quotatype field is ignored unless Quotatype_is_set == true */
	Quotatype_is_set bool
	Quotatype        bool
	/* Extent field is ignored unless Extent_is_set == true */
	Extent_is_set bool
	Extent        bool
	/* Filetype field is ignored unless Filetype_is_set == true */
	Filetype_is_set bool
	Filetype        bool
	/* Flexbg field is ignored unless Flexbg_is_set == true */
	Flexbg_is_set bool
	Flexbg        bool
	/* Hasjournal field is ignored unless Hasjournal_is_set == true */
	Hasjournal_is_set bool
	Hasjournal        bool
	/* Journaldev field is ignored unless Journaldev_is_set == true */
	Journaldev_is_set bool
	Journaldev        bool
	/* Largefile field is ignored unless Largefile_is_set == true */
	Largefile_is_set bool
	Largefile        bool
	/* Quota field is ignored unless Quota_is_set == true */
	Quota_is_set bool
	Quota        bool
	/* Resizeinode field is ignored unless Resizeinode_is_set == true */
	Resizeinode_is_set bool
	Resizeinode        bool
	/* Sparsesuper field is ignored unless Sparsesuper_is_set == true */
	Sparsesuper_is_set bool
	Sparsesuper        bool
	/* Uninitbg field is ignored unless Uninitbg_is_set == true */
	Uninitbg_is_set bool
	Uninitbg        bool
}

/* mke2fs : create an ext2/ext3/ext4 filesystem on device */
func (g *Guestfs) Mke2fs(device string, optargs *OptargsMke2fs) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("mke2fs")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))
	c_optargs := C.struct_guestfs_mke2fs_argv{}
	if optargs != nil {
		if optargs.Blockscount_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKE2FS_BLOCKSCOUNT_BITMASK
			c_optargs.blockscount = C.int64_t(optargs.Blockscount)
		}
		if optargs.Blocksize_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKE2FS_BLOCKSIZE_BITMASK
			c_optargs.blocksize = C.int64_t(optargs.Blocksize)
		}
		if optargs.Fragsize_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKE2FS_FRAGSIZE_BITMASK
			c_optargs.fragsize = C.int64_t(optargs.Fragsize)
		}
		if optargs.Blockspergroup_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKE2FS_BLOCKSPERGROUP_BITMASK
			c_optargs.blockspergroup = C.int64_t(optargs.Blockspergroup)
		}
		if optargs.Numberofgroups_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKE2FS_NUMBEROFGROUPS_BITMASK
			c_optargs.numberofgroups = C.int64_t(optargs.Numberofgroups)
		}
		if optargs.Bytesperinode_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKE2FS_BYTESPERINODE_BITMASK
			c_optargs.bytesperinode = C.int64_t(optargs.Bytesperinode)
		}
		if optargs.Inodesize_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKE2FS_INODESIZE_BITMASK
			c_optargs.inodesize = C.int64_t(optargs.Inodesize)
		}
		if optargs.Journalsize_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKE2FS_JOURNALSIZE_BITMASK
			c_optargs.journalsize = C.int64_t(optargs.Journalsize)
		}
		if optargs.Numberofinodes_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKE2FS_NUMBEROFINODES_BITMASK
			c_optargs.numberofinodes = C.int64_t(optargs.Numberofinodes)
		}
		if optargs.Stridesize_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKE2FS_STRIDESIZE_BITMASK
			c_optargs.stridesize = C.int64_t(optargs.Stridesize)
		}
		if optargs.Stripewidth_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKE2FS_STRIPEWIDTH_BITMASK
			c_optargs.stripewidth = C.int64_t(optargs.Stripewidth)
		}
		if optargs.Maxonlineresize_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKE2FS_MAXONLINERESIZE_BITMASK
			c_optargs.maxonlineresize = C.int64_t(optargs.Maxonlineresize)
		}
		if optargs.Reservedblockspercentage_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKE2FS_RESERVEDBLOCKSPERCENTAGE_BITMASK
			c_optargs.reservedblockspercentage = C.int(optargs.Reservedblockspercentage)
		}
		if optargs.Mmpupdateinterval_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKE2FS_MMPUPDATEINTERVAL_BITMASK
			c_optargs.mmpupdateinterval = C.int(optargs.Mmpupdateinterval)
		}
		if optargs.Journaldevice_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKE2FS_JOURNALDEVICE_BITMASK
			c_optargs.journaldevice = C.CString(optargs.Journaldevice)
			defer C.free(unsafe.Pointer(c_optargs.journaldevice))
		}
		if optargs.Label_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKE2FS_LABEL_BITMASK
			c_optargs.label = C.CString(optargs.Label)
			defer C.free(unsafe.Pointer(c_optargs.label))
		}
		if optargs.Lastmounteddir_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKE2FS_LASTMOUNTEDDIR_BITMASK
			c_optargs.lastmounteddir = C.CString(optargs.Lastmounteddir)
			defer C.free(unsafe.Pointer(c_optargs.lastmounteddir))
		}
		if optargs.Creatoros_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKE2FS_CREATOROS_BITMASK
			c_optargs.creatoros = C.CString(optargs.Creatoros)
			defer C.free(unsafe.Pointer(c_optargs.creatoros))
		}
		if optargs.Fstype_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKE2FS_FSTYPE_BITMASK
			c_optargs.fstype = C.CString(optargs.Fstype)
			defer C.free(unsafe.Pointer(c_optargs.fstype))
		}
		if optargs.Usagetype_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKE2FS_USAGETYPE_BITMASK
			c_optargs.usagetype = C.CString(optargs.Usagetype)
			defer C.free(unsafe.Pointer(c_optargs.usagetype))
		}
		if optargs.Uuid_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKE2FS_UUID_BITMASK
			c_optargs.uuid = C.CString(optargs.Uuid)
			defer C.free(unsafe.Pointer(c_optargs.uuid))
		}
		if optargs.Forcecreate_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKE2FS_FORCECREATE_BITMASK
			if optargs.Forcecreate {
				c_optargs.forcecreate = 1
			} else {
				c_optargs.forcecreate = 0
			}
		}
		if optargs.Writesbandgrouponly_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKE2FS_WRITESBANDGROUPONLY_BITMASK
			if optargs.Writesbandgrouponly {
				c_optargs.writesbandgrouponly = 1
			} else {
				c_optargs.writesbandgrouponly = 0
			}
		}
		if optargs.Lazyitableinit_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKE2FS_LAZYITABLEINIT_BITMASK
			if optargs.Lazyitableinit {
				c_optargs.lazyitableinit = 1
			} else {
				c_optargs.lazyitableinit = 0
			}
		}
		if optargs.Lazyjournalinit_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKE2FS_LAZYJOURNALINIT_BITMASK
			if optargs.Lazyjournalinit {
				c_optargs.lazyjournalinit = 1
			} else {
				c_optargs.lazyjournalinit = 0
			}
		}
		if optargs.Testfs_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKE2FS_TESTFS_BITMASK
			if optargs.Testfs {
				c_optargs.testfs = 1
			} else {
				c_optargs.testfs = 0
			}
		}
		if optargs.Discard_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKE2FS_DISCARD_BITMASK
			if optargs.Discard {
				c_optargs.discard = 1
			} else {
				c_optargs.discard = 0
			}
		}
		if optargs.Quotatype_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKE2FS_QUOTATYPE_BITMASK
			if optargs.Quotatype {
				c_optargs.quotatype = 1
			} else {
				c_optargs.quotatype = 0
			}
		}
		if optargs.Extent_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKE2FS_EXTENT_BITMASK
			if optargs.Extent {
				c_optargs.extent = 1
			} else {
				c_optargs.extent = 0
			}
		}
		if optargs.Filetype_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKE2FS_FILETYPE_BITMASK
			if optargs.Filetype {
				c_optargs.filetype = 1
			} else {
				c_optargs.filetype = 0
			}
		}
		if optargs.Flexbg_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKE2FS_FLEXBG_BITMASK
			if optargs.Flexbg {
				c_optargs.flexbg = 1
			} else {
				c_optargs.flexbg = 0
			}
		}
		if optargs.Hasjournal_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKE2FS_HASJOURNAL_BITMASK
			if optargs.Hasjournal {
				c_optargs.hasjournal = 1
			} else {
				c_optargs.hasjournal = 0
			}
		}
		if optargs.Journaldev_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKE2FS_JOURNALDEV_BITMASK
			if optargs.Journaldev {
				c_optargs.journaldev = 1
			} else {
				c_optargs.journaldev = 0
			}
		}
		if optargs.Largefile_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKE2FS_LARGEFILE_BITMASK
			if optargs.Largefile {
				c_optargs.largefile = 1
			} else {
				c_optargs.largefile = 0
			}
		}
		if optargs.Quota_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKE2FS_QUOTA_BITMASK
			if optargs.Quota {
				c_optargs.quota = 1
			} else {
				c_optargs.quota = 0
			}
		}
		if optargs.Resizeinode_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKE2FS_RESIZEINODE_BITMASK
			if optargs.Resizeinode {
				c_optargs.resizeinode = 1
			} else {
				c_optargs.resizeinode = 0
			}
		}
		if optargs.Sparsesuper_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKE2FS_SPARSESUPER_BITMASK
			if optargs.Sparsesuper {
				c_optargs.sparsesuper = 1
			} else {
				c_optargs.sparsesuper = 0
			}
		}
		if optargs.Uninitbg_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKE2FS_UNINITBG_BITMASK
			if optargs.Uninitbg {
				c_optargs.uninitbg = 1
			} else {
				c_optargs.uninitbg = 0
			}
		}
	}

	r := C.guestfs_mke2fs_argv(g.g, c_device, &c_optargs)

	if r == -1 {
		return get_error_from_handle(g, "mke2fs")
	}
	return nil
}

/* mke2fs_J : make ext2/3/4 filesystem with external journal */
func (g *Guestfs) Mke2fs_J(fstype string, blocksize int, device string, journal string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("mke2fs_J")
	}

	c_fstype := C.CString(fstype)
	defer C.free(unsafe.Pointer(c_fstype))

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	c_journal := C.CString(journal)
	defer C.free(unsafe.Pointer(c_journal))

	r := C.guestfs_mke2fs_J(g.g, c_fstype, C.int(blocksize), c_device, c_journal)

	if r == -1 {
		return get_error_from_handle(g, "mke2fs_J")
	}
	return nil
}

/* mke2fs_JL : make ext2/3/4 filesystem with external journal */
func (g *Guestfs) Mke2fs_JL(fstype string, blocksize int, device string, label string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("mke2fs_JL")
	}

	c_fstype := C.CString(fstype)
	defer C.free(unsafe.Pointer(c_fstype))

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	c_label := C.CString(label)
	defer C.free(unsafe.Pointer(c_label))

	r := C.guestfs_mke2fs_JL(g.g, c_fstype, C.int(blocksize), c_device, c_label)

	if r == -1 {
		return get_error_from_handle(g, "mke2fs_JL")
	}
	return nil
}

/* mke2fs_JU : make ext2/3/4 filesystem with external journal */
func (g *Guestfs) Mke2fs_JU(fstype string, blocksize int, device string, uuid string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("mke2fs_JU")
	}

	c_fstype := C.CString(fstype)
	defer C.free(unsafe.Pointer(c_fstype))

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	c_uuid := C.CString(uuid)
	defer C.free(unsafe.Pointer(c_uuid))

	r := C.guestfs_mke2fs_JU(g.g, c_fstype, C.int(blocksize), c_device, c_uuid)

	if r == -1 {
		return get_error_from_handle(g, "mke2fs_JU")
	}
	return nil
}

/* mke2journal : make ext2/3/4 external journal */
func (g *Guestfs) Mke2journal(blocksize int, device string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("mke2journal")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_mke2journal(g.g, C.int(blocksize), c_device)

	if r == -1 {
		return get_error_from_handle(g, "mke2journal")
	}
	return nil
}

/* mke2journal_L : make ext2/3/4 external journal with label */
func (g *Guestfs) Mke2journal_L(blocksize int, label string, device string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("mke2journal_L")
	}

	c_label := C.CString(label)
	defer C.free(unsafe.Pointer(c_label))

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_mke2journal_L(g.g, C.int(blocksize), c_label, c_device)

	if r == -1 {
		return get_error_from_handle(g, "mke2journal_L")
	}
	return nil
}

/* mke2journal_U : make ext2/3/4 external journal with UUID */
func (g *Guestfs) Mke2journal_U(blocksize int, uuid string, device string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("mke2journal_U")
	}

	c_uuid := C.CString(uuid)
	defer C.free(unsafe.Pointer(c_uuid))

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_mke2journal_U(g.g, C.int(blocksize), c_uuid, c_device)

	if r == -1 {
		return get_error_from_handle(g, "mke2journal_U")
	}
	return nil
}

/* mkfifo : make FIFO (named pipe) */
func (g *Guestfs) Mkfifo(mode int, path string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("mkfifo")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_mkfifo(g.g, C.int(mode), c_path)

	if r == -1 {
		return get_error_from_handle(g, "mkfifo")
	}
	return nil
}

/* Struct carrying optional arguments for Mkfs */
type OptargsMkfs struct {
	/* Blocksize field is ignored unless Blocksize_is_set == true */
	Blocksize_is_set bool
	Blocksize        int
	/* Features field is ignored unless Features_is_set == true */
	Features_is_set bool
	Features        string
	/* Inode field is ignored unless Inode_is_set == true */
	Inode_is_set bool
	Inode        int
	/* Sectorsize field is ignored unless Sectorsize_is_set == true */
	Sectorsize_is_set bool
	Sectorsize        int
	/* Label field is ignored unless Label_is_set == true */
	Label_is_set bool
	Label        string
}

/* mkfs : make a filesystem */
func (g *Guestfs) Mkfs(fstype string, device string, optargs *OptargsMkfs) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("mkfs")
	}

	c_fstype := C.CString(fstype)
	defer C.free(unsafe.Pointer(c_fstype))

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))
	c_optargs := C.struct_guestfs_mkfs_opts_argv{}
	if optargs != nil {
		if optargs.Blocksize_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKFS_OPTS_BLOCKSIZE_BITMASK
			c_optargs.blocksize = C.int(optargs.Blocksize)
		}
		if optargs.Features_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKFS_OPTS_FEATURES_BITMASK
			c_optargs.features = C.CString(optargs.Features)
			defer C.free(unsafe.Pointer(c_optargs.features))
		}
		if optargs.Inode_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKFS_OPTS_INODE_BITMASK
			c_optargs.inode = C.int(optargs.Inode)
		}
		if optargs.Sectorsize_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKFS_OPTS_SECTORSIZE_BITMASK
			c_optargs.sectorsize = C.int(optargs.Sectorsize)
		}
		if optargs.Label_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKFS_OPTS_LABEL_BITMASK
			c_optargs.label = C.CString(optargs.Label)
			defer C.free(unsafe.Pointer(c_optargs.label))
		}
	}

	r := C.guestfs_mkfs_opts_argv(g.g, c_fstype, c_device, &c_optargs)

	if r == -1 {
		return get_error_from_handle(g, "mkfs")
	}
	return nil
}

/* mkfs_b : make a filesystem with block size */
func (g *Guestfs) Mkfs_b(fstype string, blocksize int, device string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("mkfs_b")
	}

	c_fstype := C.CString(fstype)
	defer C.free(unsafe.Pointer(c_fstype))

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_mkfs_b(g.g, c_fstype, C.int(blocksize), c_device)

	if r == -1 {
		return get_error_from_handle(g, "mkfs_b")
	}
	return nil
}

/* Struct carrying optional arguments for Mkfs_btrfs */
type OptargsMkfs_btrfs struct {
	/* Allocstart field is ignored unless Allocstart_is_set == true */
	Allocstart_is_set bool
	Allocstart        int64
	/* Bytecount field is ignored unless Bytecount_is_set == true */
	Bytecount_is_set bool
	Bytecount        int64
	/* Datatype field is ignored unless Datatype_is_set == true */
	Datatype_is_set bool
	Datatype        string
	/* Leafsize field is ignored unless Leafsize_is_set == true */
	Leafsize_is_set bool
	Leafsize        int
	/* Label field is ignored unless Label_is_set == true */
	Label_is_set bool
	Label        string
	/* Metadata field is ignored unless Metadata_is_set == true */
	Metadata_is_set bool
	Metadata        string
	/* Nodesize field is ignored unless Nodesize_is_set == true */
	Nodesize_is_set bool
	Nodesize        int
	/* Sectorsize field is ignored unless Sectorsize_is_set == true */
	Sectorsize_is_set bool
	Sectorsize        int
}

/* mkfs_btrfs : create a btrfs filesystem */
func (g *Guestfs) Mkfs_btrfs(devices []string, optargs *OptargsMkfs_btrfs) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("mkfs_btrfs")
	}

	c_devices := arg_string_list(devices)
	defer free_string_list(c_devices)
	c_optargs := C.struct_guestfs_mkfs_btrfs_argv{}
	if optargs != nil {
		if optargs.Allocstart_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKFS_BTRFS_ALLOCSTART_BITMASK
			c_optargs.allocstart = C.int64_t(optargs.Allocstart)
		}
		if optargs.Bytecount_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKFS_BTRFS_BYTECOUNT_BITMASK
			c_optargs.bytecount = C.int64_t(optargs.Bytecount)
		}
		if optargs.Datatype_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKFS_BTRFS_DATATYPE_BITMASK
			c_optargs.datatype = C.CString(optargs.Datatype)
			defer C.free(unsafe.Pointer(c_optargs.datatype))
		}
		if optargs.Leafsize_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKFS_BTRFS_LEAFSIZE_BITMASK
			c_optargs.leafsize = C.int(optargs.Leafsize)
		}
		if optargs.Label_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKFS_BTRFS_LABEL_BITMASK
			c_optargs.label = C.CString(optargs.Label)
			defer C.free(unsafe.Pointer(c_optargs.label))
		}
		if optargs.Metadata_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKFS_BTRFS_METADATA_BITMASK
			c_optargs.metadata = C.CString(optargs.Metadata)
			defer C.free(unsafe.Pointer(c_optargs.metadata))
		}
		if optargs.Nodesize_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKFS_BTRFS_NODESIZE_BITMASK
			c_optargs.nodesize = C.int(optargs.Nodesize)
		}
		if optargs.Sectorsize_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKFS_BTRFS_SECTORSIZE_BITMASK
			c_optargs.sectorsize = C.int(optargs.Sectorsize)
		}
	}

	r := C.guestfs_mkfs_btrfs_argv(g.g, c_devices, &c_optargs)

	if r == -1 {
		return get_error_from_handle(g, "mkfs_btrfs")
	}
	return nil
}

/* mklost_and_found : make lost+found directory on an ext2/3/4 filesystem */
func (g *Guestfs) Mklost_and_found(mountpoint string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("mklost_and_found")
	}

	c_mountpoint := C.CString(mountpoint)
	defer C.free(unsafe.Pointer(c_mountpoint))

	r := C.guestfs_mklost_and_found(g.g, c_mountpoint)

	if r == -1 {
		return get_error_from_handle(g, "mklost_and_found")
	}
	return nil
}

/* mkmountpoint : create a mountpoint */
func (g *Guestfs) Mkmountpoint(exemptpath string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("mkmountpoint")
	}

	c_exemptpath := C.CString(exemptpath)
	defer C.free(unsafe.Pointer(c_exemptpath))

	r := C.guestfs_mkmountpoint(g.g, c_exemptpath)

	if r == -1 {
		return get_error_from_handle(g, "mkmountpoint")
	}
	return nil
}

/* mknod : make block, character or FIFO devices */
func (g *Guestfs) Mknod(mode int, devmajor int, devminor int, path string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("mknod")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_mknod(g.g, C.int(mode), C.int(devmajor), C.int(devminor), c_path)

	if r == -1 {
		return get_error_from_handle(g, "mknod")
	}
	return nil
}

/* mknod_b : make block device node */
func (g *Guestfs) Mknod_b(mode int, devmajor int, devminor int, path string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("mknod_b")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_mknod_b(g.g, C.int(mode), C.int(devmajor), C.int(devminor), c_path)

	if r == -1 {
		return get_error_from_handle(g, "mknod_b")
	}
	return nil
}

/* mknod_c : make char device node */
func (g *Guestfs) Mknod_c(mode int, devmajor int, devminor int, path string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("mknod_c")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_mknod_c(g.g, C.int(mode), C.int(devmajor), C.int(devminor), c_path)

	if r == -1 {
		return get_error_from_handle(g, "mknod_c")
	}
	return nil
}

/* Struct carrying optional arguments for Mkswap */
type OptargsMkswap struct {
	/* Label field is ignored unless Label_is_set == true */
	Label_is_set bool
	Label        string
	/* Uuid field is ignored unless Uuid_is_set == true */
	Uuid_is_set bool
	Uuid        string
}

/* mkswap : create a swap partition */
func (g *Guestfs) Mkswap(device string, optargs *OptargsMkswap) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("mkswap")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))
	c_optargs := C.struct_guestfs_mkswap_opts_argv{}
	if optargs != nil {
		if optargs.Label_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKSWAP_OPTS_LABEL_BITMASK
			c_optargs.label = C.CString(optargs.Label)
			defer C.free(unsafe.Pointer(c_optargs.label))
		}
		if optargs.Uuid_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKSWAP_OPTS_UUID_BITMASK
			c_optargs.uuid = C.CString(optargs.Uuid)
			defer C.free(unsafe.Pointer(c_optargs.uuid))
		}
	}

	r := C.guestfs_mkswap_opts_argv(g.g, c_device, &c_optargs)

	if r == -1 {
		return get_error_from_handle(g, "mkswap")
	}
	return nil
}

/* mkswap_L : create a swap partition with a label */
func (g *Guestfs) Mkswap_L(label string, device string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("mkswap_L")
	}

	c_label := C.CString(label)
	defer C.free(unsafe.Pointer(c_label))

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_mkswap_L(g.g, c_label, c_device)

	if r == -1 {
		return get_error_from_handle(g, "mkswap_L")
	}
	return nil
}

/* mkswap_U : create a swap partition with an explicit UUID */
func (g *Guestfs) Mkswap_U(uuid string, device string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("mkswap_U")
	}

	c_uuid := C.CString(uuid)
	defer C.free(unsafe.Pointer(c_uuid))

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_mkswap_U(g.g, c_uuid, c_device)

	if r == -1 {
		return get_error_from_handle(g, "mkswap_U")
	}
	return nil
}

/* mkswap_file : create a swap file */
func (g *Guestfs) Mkswap_file(path string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("mkswap_file")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_mkswap_file(g.g, c_path)

	if r == -1 {
		return get_error_from_handle(g, "mkswap_file")
	}
	return nil
}

/* Struct carrying optional arguments for Mktemp */
type OptargsMktemp struct {
	/* Suffix field is ignored unless Suffix_is_set == true */
	Suffix_is_set bool
	Suffix        string
}

/* mktemp : create a temporary file */
func (g *Guestfs) Mktemp(tmpl string, optargs *OptargsMktemp) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("mktemp")
	}

	c_tmpl := C.CString(tmpl)
	defer C.free(unsafe.Pointer(c_tmpl))
	c_optargs := C.struct_guestfs_mktemp_argv{}
	if optargs != nil {
		if optargs.Suffix_is_set {
			c_optargs.bitmask |= C.GUESTFS_MKTEMP_SUFFIX_BITMASK
			c_optargs.suffix = C.CString(optargs.Suffix)
			defer C.free(unsafe.Pointer(c_optargs.suffix))
		}
	}

	r := C.guestfs_mktemp_argv(g.g, c_tmpl, &c_optargs)

	if r == nil {
		return "", get_error_from_handle(g, "mktemp")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* modprobe : load a kernel module */
func (g *Guestfs) Modprobe(modulename string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("modprobe")
	}

	c_modulename := C.CString(modulename)
	defer C.free(unsafe.Pointer(c_modulename))

	r := C.guestfs_modprobe(g.g, c_modulename)

	if r == -1 {
		return get_error_from_handle(g, "modprobe")
	}
	return nil
}

/* mount : mount a guest disk at a position in the filesystem */
func (g *Guestfs) Mount(mountable string, mountpoint string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("mount")
	}

	c_mountable := C.CString(mountable)
	defer C.free(unsafe.Pointer(c_mountable))

	c_mountpoint := C.CString(mountpoint)
	defer C.free(unsafe.Pointer(c_mountpoint))

	r := C.guestfs_mount(g.g, c_mountable, c_mountpoint)

	if r == -1 {
		return get_error_from_handle(g, "mount")
	}
	return nil
}

/* Struct carrying optional arguments for Mount_9p */
type OptargsMount_9p struct {
	/* Options field is ignored unless Options_is_set == true */
	Options_is_set bool
	Options        string
}

/* mount_9p : mount 9p filesystem */
func (g *Guestfs) Mount_9p(mounttag string, mountpoint string, optargs *OptargsMount_9p) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("mount_9p")
	}

	c_mounttag := C.CString(mounttag)
	defer C.free(unsafe.Pointer(c_mounttag))

	c_mountpoint := C.CString(mountpoint)
	defer C.free(unsafe.Pointer(c_mountpoint))
	c_optargs := C.struct_guestfs_mount_9p_argv{}
	if optargs != nil {
		if optargs.Options_is_set {
			c_optargs.bitmask |= C.GUESTFS_MOUNT_9P_OPTIONS_BITMASK
			c_optargs.options = C.CString(optargs.Options)
			defer C.free(unsafe.Pointer(c_optargs.options))
		}
	}

	r := C.guestfs_mount_9p_argv(g.g, c_mounttag, c_mountpoint, &c_optargs)

	if r == -1 {
		return get_error_from_handle(g, "mount_9p")
	}
	return nil
}

/* Struct carrying optional arguments for Mount_local */
type OptargsMount_local struct {
	/* Readonly field is ignored unless Readonly_is_set == true */
	Readonly_is_set bool
	Readonly        bool
	/* Options field is ignored unless Options_is_set == true */
	Options_is_set bool
	Options        string
	/* Cachetimeout field is ignored unless Cachetimeout_is_set == true */
	Cachetimeout_is_set bool
	Cachetimeout        int
	/* Debugcalls field is ignored unless Debugcalls_is_set == true */
	Debugcalls_is_set bool
	Debugcalls        bool
}

/* mount_local : mount on the local filesystem */
func (g *Guestfs) Mount_local(localmountpoint string, optargs *OptargsMount_local) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("mount_local")
	}

	c_localmountpoint := C.CString(localmountpoint)
	defer C.free(unsafe.Pointer(c_localmountpoint))
	c_optargs := C.struct_guestfs_mount_local_argv{}
	if optargs != nil {
		if optargs.Readonly_is_set {
			c_optargs.bitmask |= C.GUESTFS_MOUNT_LOCAL_READONLY_BITMASK
			if optargs.Readonly {
				c_optargs.readonly = 1
			} else {
				c_optargs.readonly = 0
			}
		}
		if optargs.Options_is_set {
			c_optargs.bitmask |= C.GUESTFS_MOUNT_LOCAL_OPTIONS_BITMASK
			c_optargs.options = C.CString(optargs.Options)
			defer C.free(unsafe.Pointer(c_optargs.options))
		}
		if optargs.Cachetimeout_is_set {
			c_optargs.bitmask |= C.GUESTFS_MOUNT_LOCAL_CACHETIMEOUT_BITMASK
			c_optargs.cachetimeout = C.int(optargs.Cachetimeout)
		}
		if optargs.Debugcalls_is_set {
			c_optargs.bitmask |= C.GUESTFS_MOUNT_LOCAL_DEBUGCALLS_BITMASK
			if optargs.Debugcalls {
				c_optargs.debugcalls = 1
			} else {
				c_optargs.debugcalls = 0
			}
		}
	}

	r := C.guestfs_mount_local_argv(g.g, c_localmountpoint, &c_optargs)

	if r == -1 {
		return get_error_from_handle(g, "mount_local")
	}
	return nil
}

/* mount_local_run : run main loop of mount on the local filesystem */
func (g *Guestfs) Mount_local_run() *GuestfsError {
	if g.g == nil {
		return closed_handle_error("mount_local_run")
	}

	r := C.guestfs_mount_local_run(g.g)

	if r == -1 {
		return get_error_from_handle(g, "mount_local_run")
	}
	return nil
}

/* mount_loop : mount a file using the loop device */
func (g *Guestfs) Mount_loop(file string, mountpoint string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("mount_loop")
	}

	c_file := C.CString(file)
	defer C.free(unsafe.Pointer(c_file))

	c_mountpoint := C.CString(mountpoint)
	defer C.free(unsafe.Pointer(c_mountpoint))

	r := C.guestfs_mount_loop(g.g, c_file, c_mountpoint)

	if r == -1 {
		return get_error_from_handle(g, "mount_loop")
	}
	return nil
}

/* mount_options : mount a guest disk with mount options */
func (g *Guestfs) Mount_options(options string, mountable string, mountpoint string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("mount_options")
	}

	c_options := C.CString(options)
	defer C.free(unsafe.Pointer(c_options))

	c_mountable := C.CString(mountable)
	defer C.free(unsafe.Pointer(c_mountable))

	c_mountpoint := C.CString(mountpoint)
	defer C.free(unsafe.Pointer(c_mountpoint))

	r := C.guestfs_mount_options(g.g, c_options, c_mountable, c_mountpoint)

	if r == -1 {
		return get_error_from_handle(g, "mount_options")
	}
	return nil
}

/* mount_ro : mount a guest disk, read-only */
func (g *Guestfs) Mount_ro(mountable string, mountpoint string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("mount_ro")
	}

	c_mountable := C.CString(mountable)
	defer C.free(unsafe.Pointer(c_mountable))

	c_mountpoint := C.CString(mountpoint)
	defer C.free(unsafe.Pointer(c_mountpoint))

	r := C.guestfs_mount_ro(g.g, c_mountable, c_mountpoint)

	if r == -1 {
		return get_error_from_handle(g, "mount_ro")
	}
	return nil
}

/* mount_vfs : mount a guest disk with mount options and vfstype */
func (g *Guestfs) Mount_vfs(options string, vfstype string, mountable string, mountpoint string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("mount_vfs")
	}

	c_options := C.CString(options)
	defer C.free(unsafe.Pointer(c_options))

	c_vfstype := C.CString(vfstype)
	defer C.free(unsafe.Pointer(c_vfstype))

	c_mountable := C.CString(mountable)
	defer C.free(unsafe.Pointer(c_mountable))

	c_mountpoint := C.CString(mountpoint)
	defer C.free(unsafe.Pointer(c_mountpoint))

	r := C.guestfs_mount_vfs(g.g, c_options, c_vfstype, c_mountable, c_mountpoint)

	if r == -1 {
		return get_error_from_handle(g, "mount_vfs")
	}
	return nil
}

/* mountpoints : show mountpoints */
func (g *Guestfs) Mountpoints() (map[string]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("mountpoints")
	}

	r := C.guestfs_mountpoints(g.g)

	if r == nil {
		return nil, get_error_from_handle(g, "mountpoints")
	}
	defer free_string_list(r)
	return return_hashtable(r), nil
}

/* mounts : show mounted filesystems */
func (g *Guestfs) Mounts() ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("mounts")
	}

	r := C.guestfs_mounts(g.g)

	if r == nil {
		return nil, get_error_from_handle(g, "mounts")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* mv : move a file */
func (g *Guestfs) Mv(src string, dest string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("mv")
	}

	c_src := C.CString(src)
	defer C.free(unsafe.Pointer(c_src))

	c_dest := C.CString(dest)
	defer C.free(unsafe.Pointer(c_dest))

	r := C.guestfs_mv(g.g, c_src, c_dest)

	if r == -1 {
		return get_error_from_handle(g, "mv")
	}
	return nil
}

/* nr_devices : return number of whole block devices (disks) added */
func (g *Guestfs) Nr_devices() (int, *GuestfsError) {
	if g.g == nil {
		return 0, closed_handle_error("nr_devices")
	}

	r := C.guestfs_nr_devices(g.g)

	if r == -1 {
		return 0, get_error_from_handle(g, "nr_devices")
	}
	return int(r), nil
}

/* ntfs_3g_probe : probe NTFS volume */
func (g *Guestfs) Ntfs_3g_probe(rw bool, device string) (int, *GuestfsError) {
	if g.g == nil {
		return 0, closed_handle_error("ntfs_3g_probe")
	}

	var c_rw C.int
	if rw {
		c_rw = 1
	} else {
		c_rw = 0
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_ntfs_3g_probe(g.g, c_rw, c_device)

	if r == -1 {
		return 0, get_error_from_handle(g, "ntfs_3g_probe")
	}
	return int(r), nil
}

/* ntfsclone_in : restore NTFS from backup file */
func (g *Guestfs) Ntfsclone_in(backupfile string, device string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("ntfsclone_in")
	}

	c_backupfile := C.CString(backupfile)
	defer C.free(unsafe.Pointer(c_backupfile))

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_ntfsclone_in(g.g, c_backupfile, c_device)

	if r == -1 {
		return get_error_from_handle(g, "ntfsclone_in")
	}
	return nil
}

/* Struct carrying optional arguments for Ntfsclone_out */
type OptargsNtfsclone_out struct {
	/* Metadataonly field is ignored unless Metadataonly_is_set == true */
	Metadataonly_is_set bool
	Metadataonly        bool
	/* Rescue field is ignored unless Rescue_is_set == true */
	Rescue_is_set bool
	Rescue        bool
	/* Ignorefscheck field is ignored unless Ignorefscheck_is_set == true */
	Ignorefscheck_is_set bool
	Ignorefscheck        bool
	/* Preservetimestamps field is ignored unless Preservetimestamps_is_set == true */
	Preservetimestamps_is_set bool
	Preservetimestamps        bool
	/* Force field is ignored unless Force_is_set == true */
	Force_is_set bool
	Force        bool
}

/* ntfsclone_out : save NTFS to backup file */
func (g *Guestfs) Ntfsclone_out(device string, backupfile string, optargs *OptargsNtfsclone_out) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("ntfsclone_out")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	c_backupfile := C.CString(backupfile)
	defer C.free(unsafe.Pointer(c_backupfile))
	c_optargs := C.struct_guestfs_ntfsclone_out_argv{}
	if optargs != nil {
		if optargs.Metadataonly_is_set {
			c_optargs.bitmask |= C.GUESTFS_NTFSCLONE_OUT_METADATAONLY_BITMASK
			if optargs.Metadataonly {
				c_optargs.metadataonly = 1
			} else {
				c_optargs.metadataonly = 0
			}
		}
		if optargs.Rescue_is_set {
			c_optargs.bitmask |= C.GUESTFS_NTFSCLONE_OUT_RESCUE_BITMASK
			if optargs.Rescue {
				c_optargs.rescue = 1
			} else {
				c_optargs.rescue = 0
			}
		}
		if optargs.Ignorefscheck_is_set {
			c_optargs.bitmask |= C.GUESTFS_NTFSCLONE_OUT_IGNOREFSCHECK_BITMASK
			if optargs.Ignorefscheck {
				c_optargs.ignorefscheck = 1
			} else {
				c_optargs.ignorefscheck = 0
			}
		}
		if optargs.Preservetimestamps_is_set {
			c_optargs.bitmask |= C.GUESTFS_NTFSCLONE_OUT_PRESERVETIMESTAMPS_BITMASK
			if optargs.Preservetimestamps {
				c_optargs.preservetimestamps = 1
			} else {
				c_optargs.preservetimestamps = 0
			}
		}
		if optargs.Force_is_set {
			c_optargs.bitmask |= C.GUESTFS_NTFSCLONE_OUT_FORCE_BITMASK
			if optargs.Force {
				c_optargs.force = 1
			} else {
				c_optargs.force = 0
			}
		}
	}

	r := C.guestfs_ntfsclone_out_argv(g.g, c_device, c_backupfile, &c_optargs)

	if r == -1 {
		return get_error_from_handle(g, "ntfsclone_out")
	}
	return nil
}

/* Struct carrying optional arguments for Ntfsfix */
type OptargsNtfsfix struct {
	/* Clearbadsectors field is ignored unless Clearbadsectors_is_set == true */
	Clearbadsectors_is_set bool
	Clearbadsectors        bool
}

/* ntfsfix : fix common errors and force Windows to check NTFS */
func (g *Guestfs) Ntfsfix(device string, optargs *OptargsNtfsfix) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("ntfsfix")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))
	c_optargs := C.struct_guestfs_ntfsfix_argv{}
	if optargs != nil {
		if optargs.Clearbadsectors_is_set {
			c_optargs.bitmask |= C.GUESTFS_NTFSFIX_CLEARBADSECTORS_BITMASK
			if optargs.Clearbadsectors {
				c_optargs.clearbadsectors = 1
			} else {
				c_optargs.clearbadsectors = 0
			}
		}
	}

	r := C.guestfs_ntfsfix_argv(g.g, c_device, &c_optargs)

	if r == -1 {
		return get_error_from_handle(g, "ntfsfix")
	}
	return nil
}

/* Struct carrying optional arguments for Ntfsresize */
type OptargsNtfsresize struct {
	/* Size field is ignored unless Size_is_set == true */
	Size_is_set bool
	Size        int64
	/* Force field is ignored unless Force_is_set == true */
	Force_is_set bool
	Force        bool
}

/* ntfsresize : resize an NTFS filesystem */
func (g *Guestfs) Ntfsresize(device string, optargs *OptargsNtfsresize) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("ntfsresize")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))
	c_optargs := C.struct_guestfs_ntfsresize_opts_argv{}
	if optargs != nil {
		if optargs.Size_is_set {
			c_optargs.bitmask |= C.GUESTFS_NTFSRESIZE_OPTS_SIZE_BITMASK
			c_optargs.size = C.int64_t(optargs.Size)
		}
		if optargs.Force_is_set {
			c_optargs.bitmask |= C.GUESTFS_NTFSRESIZE_OPTS_FORCE_BITMASK
			if optargs.Force {
				c_optargs.force = 1
			} else {
				c_optargs.force = 0
			}
		}
	}

	r := C.guestfs_ntfsresize_opts_argv(g.g, c_device, &c_optargs)

	if r == -1 {
		return get_error_from_handle(g, "ntfsresize")
	}
	return nil
}

/* ntfsresize_size : resize an NTFS filesystem (with size) */
func (g *Guestfs) Ntfsresize_size(device string, size int64) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("ntfsresize_size")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_ntfsresize_size(g.g, c_device, C.int64_t(size))

	if r == -1 {
		return get_error_from_handle(g, "ntfsresize_size")
	}
	return nil
}

/* parse_environment : parse the environment and set handle flags accordingly */
func (g *Guestfs) Parse_environment() *GuestfsError {
	if g.g == nil {
		return closed_handle_error("parse_environment")
	}

	r := C.guestfs_parse_environment(g.g)

	if r == -1 {
		return get_error_from_handle(g, "parse_environment")
	}
	return nil
}

/* parse_environment_list : parse the environment and set handle flags accordingly */
func (g *Guestfs) Parse_environment_list(environment []string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("parse_environment_list")
	}

	c_environment := arg_string_list(environment)
	defer free_string_list(c_environment)

	r := C.guestfs_parse_environment_list(g.g, c_environment)

	if r == -1 {
		return get_error_from_handle(g, "parse_environment_list")
	}
	return nil
}

/* part_add : add a partition to the device */
func (g *Guestfs) Part_add(device string, prlogex string, startsect int64, endsect int64) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("part_add")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	c_prlogex := C.CString(prlogex)
	defer C.free(unsafe.Pointer(c_prlogex))

	r := C.guestfs_part_add(g.g, c_device, c_prlogex, C.int64_t(startsect), C.int64_t(endsect))

	if r == -1 {
		return get_error_from_handle(g, "part_add")
	}
	return nil
}

/* part_del : delete a partition */
func (g *Guestfs) Part_del(device string, partnum int) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("part_del")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_part_del(g.g, c_device, C.int(partnum))

	if r == -1 {
		return get_error_from_handle(g, "part_del")
	}
	return nil
}

/* part_disk : partition whole disk with a single primary partition */
func (g *Guestfs) Part_disk(device string, parttype string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("part_disk")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	c_parttype := C.CString(parttype)
	defer C.free(unsafe.Pointer(c_parttype))

	r := C.guestfs_part_disk(g.g, c_device, c_parttype)

	if r == -1 {
		return get_error_from_handle(g, "part_disk")
	}
	return nil
}

/* part_get_bootable : return true if a partition is bootable */
func (g *Guestfs) Part_get_bootable(device string, partnum int) (bool, *GuestfsError) {
	if g.g == nil {
		return false, closed_handle_error("part_get_bootable")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_part_get_bootable(g.g, c_device, C.int(partnum))

	if r == -1 {
		return false, get_error_from_handle(g, "part_get_bootable")
	}
	return r != 0, nil
}

/* part_get_gpt_guid : get the GUID of a GPT partition */
func (g *Guestfs) Part_get_gpt_guid(device string, partnum int) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("part_get_gpt_guid")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_part_get_gpt_guid(g.g, c_device, C.int(partnum))

	if r == nil {
		return "", get_error_from_handle(g, "part_get_gpt_guid")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* part_get_gpt_type : get the type GUID of a GPT partition */
func (g *Guestfs) Part_get_gpt_type(device string, partnum int) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("part_get_gpt_type")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_part_get_gpt_type(g.g, c_device, C.int(partnum))

	if r == nil {
		return "", get_error_from_handle(g, "part_get_gpt_type")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* part_get_mbr_id : get the MBR type byte (ID byte) from a partition */
func (g *Guestfs) Part_get_mbr_id(device string, partnum int) (int, *GuestfsError) {
	if g.g == nil {
		return 0, closed_handle_error("part_get_mbr_id")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_part_get_mbr_id(g.g, c_device, C.int(partnum))

	if r == -1 {
		return 0, get_error_from_handle(g, "part_get_mbr_id")
	}
	return int(r), nil
}

/* part_get_mbr_part_type : get the MBR partition type */
func (g *Guestfs) Part_get_mbr_part_type(device string, partnum int) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("part_get_mbr_part_type")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_part_get_mbr_part_type(g.g, c_device, C.int(partnum))

	if r == nil {
		return "", get_error_from_handle(g, "part_get_mbr_part_type")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* part_get_name : get partition name */
func (g *Guestfs) Part_get_name(device string, partnum int) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("part_get_name")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_part_get_name(g.g, c_device, C.int(partnum))

	if r == nil {
		return "", get_error_from_handle(g, "part_get_name")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* part_get_parttype : get the partition table type */
func (g *Guestfs) Part_get_parttype(device string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("part_get_parttype")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_part_get_parttype(g.g, c_device)

	if r == nil {
		return "", get_error_from_handle(g, "part_get_parttype")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* part_init : create an empty partition table */
func (g *Guestfs) Part_init(device string, parttype string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("part_init")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	c_parttype := C.CString(parttype)
	defer C.free(unsafe.Pointer(c_parttype))

	r := C.guestfs_part_init(g.g, c_device, c_parttype)

	if r == -1 {
		return get_error_from_handle(g, "part_init")
	}
	return nil
}

/* part_list : list partitions on a device */
func (g *Guestfs) Part_list(device string) (*[]Partition, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("part_list")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_part_list(g.g, c_device)

	if r == nil {
		return nil, get_error_from_handle(g, "part_list")
	}
	defer C.guestfs_free_partition_list(r)
	return return_Partition_list(r), nil
}

/* part_set_bootable : make a partition bootable */
func (g *Guestfs) Part_set_bootable(device string, partnum int, bootable bool) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("part_set_bootable")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	var c_bootable C.int
	if bootable {
		c_bootable = 1
	} else {
		c_bootable = 0
	}

	r := C.guestfs_part_set_bootable(g.g, c_device, C.int(partnum), c_bootable)

	if r == -1 {
		return get_error_from_handle(g, "part_set_bootable")
	}
	return nil
}

/* part_set_gpt_guid : set the GUID of a GPT partition */
func (g *Guestfs) Part_set_gpt_guid(device string, partnum int, guid string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("part_set_gpt_guid")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	c_guid := C.CString(guid)
	defer C.free(unsafe.Pointer(c_guid))

	r := C.guestfs_part_set_gpt_guid(g.g, c_device, C.int(partnum), c_guid)

	if r == -1 {
		return get_error_from_handle(g, "part_set_gpt_guid")
	}
	return nil
}

/* part_set_gpt_type : set the type GUID of a GPT partition */
func (g *Guestfs) Part_set_gpt_type(device string, partnum int, guid string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("part_set_gpt_type")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	c_guid := C.CString(guid)
	defer C.free(unsafe.Pointer(c_guid))

	r := C.guestfs_part_set_gpt_type(g.g, c_device, C.int(partnum), c_guid)

	if r == -1 {
		return get_error_from_handle(g, "part_set_gpt_type")
	}
	return nil
}

/* part_set_mbr_id : set the MBR type byte (ID byte) of a partition */
func (g *Guestfs) Part_set_mbr_id(device string, partnum int, idbyte int) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("part_set_mbr_id")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_part_set_mbr_id(g.g, c_device, C.int(partnum), C.int(idbyte))

	if r == -1 {
		return get_error_from_handle(g, "part_set_mbr_id")
	}
	return nil
}

/* part_set_name : set partition name */
func (g *Guestfs) Part_set_name(device string, partnum int, name string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("part_set_name")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	c_name := C.CString(name)
	defer C.free(unsafe.Pointer(c_name))

	r := C.guestfs_part_set_name(g.g, c_device, C.int(partnum), c_name)

	if r == -1 {
		return get_error_from_handle(g, "part_set_name")
	}
	return nil
}

/* part_to_dev : convert partition name to device name */
func (g *Guestfs) Part_to_dev(partition string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("part_to_dev")
	}

	c_partition := C.CString(partition)
	defer C.free(unsafe.Pointer(c_partition))

	r := C.guestfs_part_to_dev(g.g, c_partition)

	if r == nil {
		return "", get_error_from_handle(g, "part_to_dev")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* part_to_partnum : convert partition name to partition number */
func (g *Guestfs) Part_to_partnum(partition string) (int, *GuestfsError) {
	if g.g == nil {
		return 0, closed_handle_error("part_to_partnum")
	}

	c_partition := C.CString(partition)
	defer C.free(unsafe.Pointer(c_partition))

	r := C.guestfs_part_to_partnum(g.g, c_partition)

	if r == -1 {
		return 0, get_error_from_handle(g, "part_to_partnum")
	}
	return int(r), nil
}

/* ping_daemon : ping the guest daemon */
func (g *Guestfs) Ping_daemon() *GuestfsError {
	if g.g == nil {
		return closed_handle_error("ping_daemon")
	}

	r := C.guestfs_ping_daemon(g.g)

	if r == -1 {
		return get_error_from_handle(g, "ping_daemon")
	}
	return nil
}

/* pread : read part of a file */
func (g *Guestfs) Pread(path string, count int, offset int64) ([]byte, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("pread")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	var size C.size_t

	r := C.guestfs_pread(g.g, c_path, C.int(count), C.int64_t(offset), &size)

	if r == nil {
		return nil, get_error_from_handle(g, "pread")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoBytes(unsafe.Pointer(r), C.int(size)), nil
}

/* pread_device : read part of a device */
func (g *Guestfs) Pread_device(device string, count int, offset int64) ([]byte, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("pread_device")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	var size C.size_t

	r := C.guestfs_pread_device(g.g, c_device, C.int(count), C.int64_t(offset), &size)

	if r == nil {
		return nil, get_error_from_handle(g, "pread_device")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoBytes(unsafe.Pointer(r), C.int(size)), nil
}

/* pvchange_uuid : generate a new random UUID for a physical volume */
func (g *Guestfs) Pvchange_uuid(device string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("pvchange_uuid")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_pvchange_uuid(g.g, c_device)

	if r == -1 {
		return get_error_from_handle(g, "pvchange_uuid")
	}
	return nil
}

/* pvchange_uuid_all : generate new random UUIDs for all physical volumes */
func (g *Guestfs) Pvchange_uuid_all() *GuestfsError {
	if g.g == nil {
		return closed_handle_error("pvchange_uuid_all")
	}

	r := C.guestfs_pvchange_uuid_all(g.g)

	if r == -1 {
		return get_error_from_handle(g, "pvchange_uuid_all")
	}
	return nil
}

/* pvcreate : create an LVM physical volume */
func (g *Guestfs) Pvcreate(device string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("pvcreate")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_pvcreate(g.g, c_device)

	if r == -1 {
		return get_error_from_handle(g, "pvcreate")
	}
	return nil
}

/* pvremove : remove an LVM physical volume */
func (g *Guestfs) Pvremove(device string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("pvremove")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_pvremove(g.g, c_device)

	if r == -1 {
		return get_error_from_handle(g, "pvremove")
	}
	return nil
}

/* pvresize : resize an LVM physical volume */
func (g *Guestfs) Pvresize(device string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("pvresize")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_pvresize(g.g, c_device)

	if r == -1 {
		return get_error_from_handle(g, "pvresize")
	}
	return nil
}

/* pvresize_size : resize an LVM physical volume (with size) */
func (g *Guestfs) Pvresize_size(device string, size int64) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("pvresize_size")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_pvresize_size(g.g, c_device, C.int64_t(size))

	if r == -1 {
		return get_error_from_handle(g, "pvresize_size")
	}
	return nil
}

/* pvs : list the LVM physical volumes (PVs) */
func (g *Guestfs) Pvs() ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("pvs")
	}

	r := C.guestfs_pvs(g.g)

	if r == nil {
		return nil, get_error_from_handle(g, "pvs")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* pvs_full : list the LVM physical volumes (PVs) */
func (g *Guestfs) Pvs_full() (*[]PV, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("pvs_full")
	}

	r := C.guestfs_pvs_full(g.g)

	if r == nil {
		return nil, get_error_from_handle(g, "pvs_full")
	}
	defer C.guestfs_free_lvm_pv_list(r)
	return return_PV_list(r), nil
}

/* pvuuid : get the UUID of a physical volume */
func (g *Guestfs) Pvuuid(device string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("pvuuid")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_pvuuid(g.g, c_device)

	if r == nil {
		return "", get_error_from_handle(g, "pvuuid")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* pwrite : write to part of a file */
func (g *Guestfs) Pwrite(path string, content []byte, offset int64) (int, *GuestfsError) {
	if g.g == nil {
		return 0, closed_handle_error("pwrite")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	/* string() cast here is apparently safe because
	 *   "Converting a slice of bytes to a string type yields
	 *   a string whose successive bytes are the elements of
	 *   the slice."
	 */
	c_content := C.CString(string(content))
	defer C.free(unsafe.Pointer(c_content))

	r := C.guestfs_pwrite(g.g, c_path, c_content, C.size_t(len(content)), C.int64_t(offset))

	if r == -1 {
		return 0, get_error_from_handle(g, "pwrite")
	}
	return int(r), nil
}

/* pwrite_device : write to part of a device */
func (g *Guestfs) Pwrite_device(device string, content []byte, offset int64) (int, *GuestfsError) {
	if g.g == nil {
		return 0, closed_handle_error("pwrite_device")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	/* string() cast here is apparently safe because
	 *   "Converting a slice of bytes to a string type yields
	 *   a string whose successive bytes are the elements of
	 *   the slice."
	 */
	c_content := C.CString(string(content))
	defer C.free(unsafe.Pointer(c_content))

	r := C.guestfs_pwrite_device(g.g, c_device, c_content, C.size_t(len(content)), C.int64_t(offset))

	if r == -1 {
		return 0, get_error_from_handle(g, "pwrite_device")
	}
	return int(r), nil
}

/* read_file : read a file */
func (g *Guestfs) Read_file(path string) ([]byte, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("read_file")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	var size C.size_t

	r := C.guestfs_read_file(g.g, c_path, &size)

	if r == nil {
		return nil, get_error_from_handle(g, "read_file")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoBytes(unsafe.Pointer(r), C.int(size)), nil
}

/* read_lines : read file as lines */
func (g *Guestfs) Read_lines(path string) ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("read_lines")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_read_lines(g.g, c_path)

	if r == nil {
		return nil, get_error_from_handle(g, "read_lines")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* readdir : read directories entries */
func (g *Guestfs) Readdir(dir string) (*[]Dirent, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("readdir")
	}

	c_dir := C.CString(dir)
	defer C.free(unsafe.Pointer(c_dir))

	r := C.guestfs_readdir(g.g, c_dir)

	if r == nil {
		return nil, get_error_from_handle(g, "readdir")
	}
	defer C.guestfs_free_dirent_list(r)
	return return_Dirent_list(r), nil
}

/* readlink : read the target of a symbolic link */
func (g *Guestfs) Readlink(path string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("readlink")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_readlink(g.g, c_path)

	if r == nil {
		return "", get_error_from_handle(g, "readlink")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* readlinklist : readlink on multiple files */
func (g *Guestfs) Readlinklist(path string, names []string) ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("readlinklist")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	c_names := arg_string_list(names)
	defer free_string_list(c_names)

	r := C.guestfs_readlinklist(g.g, c_path, c_names)

	if r == nil {
		return nil, get_error_from_handle(g, "readlinklist")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* realpath : canonicalized absolute pathname */
func (g *Guestfs) Realpath(path string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("realpath")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_realpath(g.g, c_path)

	if r == nil {
		return "", get_error_from_handle(g, "realpath")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* Struct carrying optional arguments for Remount */
type OptargsRemount struct {
	/* Rw field is ignored unless Rw_is_set == true */
	Rw_is_set bool
	Rw        bool
}

/* remount : remount a filesystem with different options */
func (g *Guestfs) Remount(mountpoint string, optargs *OptargsRemount) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("remount")
	}

	c_mountpoint := C.CString(mountpoint)
	defer C.free(unsafe.Pointer(c_mountpoint))
	c_optargs := C.struct_guestfs_remount_argv{}
	if optargs != nil {
		if optargs.Rw_is_set {
			c_optargs.bitmask |= C.GUESTFS_REMOUNT_RW_BITMASK
			if optargs.Rw {
				c_optargs.rw = 1
			} else {
				c_optargs.rw = 0
			}
		}
	}

	r := C.guestfs_remount_argv(g.g, c_mountpoint, &c_optargs)

	if r == -1 {
		return get_error_from_handle(g, "remount")
	}
	return nil
}

/* remove_drive : remove a disk image */
func (g *Guestfs) Remove_drive(label string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("remove_drive")
	}

	c_label := C.CString(label)
	defer C.free(unsafe.Pointer(c_label))

	r := C.guestfs_remove_drive(g.g, c_label)

	if r == -1 {
		return get_error_from_handle(g, "remove_drive")
	}
	return nil
}

/* removexattr : remove extended attribute of a file or directory */
func (g *Guestfs) Removexattr(xattr string, path string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("removexattr")
	}

	c_xattr := C.CString(xattr)
	defer C.free(unsafe.Pointer(c_xattr))

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_removexattr(g.g, c_xattr, c_path)

	if r == -1 {
		return get_error_from_handle(g, "removexattr")
	}
	return nil
}

/* rename : rename a file on the same filesystem */
func (g *Guestfs) Rename(oldpath string, newpath string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("rename")
	}

	c_oldpath := C.CString(oldpath)
	defer C.free(unsafe.Pointer(c_oldpath))

	c_newpath := C.CString(newpath)
	defer C.free(unsafe.Pointer(c_newpath))

	r := C.guestfs_rename(g.g, c_oldpath, c_newpath)

	if r == -1 {
		return get_error_from_handle(g, "rename")
	}
	return nil
}

/* resize2fs : resize an ext2, ext3 or ext4 filesystem */
func (g *Guestfs) Resize2fs(device string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("resize2fs")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_resize2fs(g.g, c_device)

	if r == -1 {
		return get_error_from_handle(g, "resize2fs")
	}
	return nil
}

/* resize2fs_M : resize an ext2, ext3 or ext4 filesystem to the minimum size */
func (g *Guestfs) Resize2fs_M(device string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("resize2fs_M")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_resize2fs_M(g.g, c_device)

	if r == -1 {
		return get_error_from_handle(g, "resize2fs_M")
	}
	return nil
}

/* resize2fs_size : resize an ext2, ext3 or ext4 filesystem (with size) */
func (g *Guestfs) Resize2fs_size(device string, size int64) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("resize2fs_size")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_resize2fs_size(g.g, c_device, C.int64_t(size))

	if r == -1 {
		return get_error_from_handle(g, "resize2fs_size")
	}
	return nil
}

/* rm : remove a file */
func (g *Guestfs) Rm(path string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("rm")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_rm(g.g, c_path)

	if r == -1 {
		return get_error_from_handle(g, "rm")
	}
	return nil
}

/* rm_f : remove a file ignoring errors */
func (g *Guestfs) Rm_f(path string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("rm_f")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_rm_f(g.g, c_path)

	if r == -1 {
		return get_error_from_handle(g, "rm_f")
	}
	return nil
}

/* rm_rf : remove a file or directory recursively */
func (g *Guestfs) Rm_rf(path string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("rm_rf")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_rm_rf(g.g, c_path)

	if r == -1 {
		return get_error_from_handle(g, "rm_rf")
	}
	return nil
}

/* rmdir : remove a directory */
func (g *Guestfs) Rmdir(path string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("rmdir")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_rmdir(g.g, c_path)

	if r == -1 {
		return get_error_from_handle(g, "rmdir")
	}
	return nil
}

/* rmmountpoint : remove a mountpoint */
func (g *Guestfs) Rmmountpoint(exemptpath string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("rmmountpoint")
	}

	c_exemptpath := C.CString(exemptpath)
	defer C.free(unsafe.Pointer(c_exemptpath))

	r := C.guestfs_rmmountpoint(g.g, c_exemptpath)

	if r == -1 {
		return get_error_from_handle(g, "rmmountpoint")
	}
	return nil
}

/* Struct carrying optional arguments for Rsync */
type OptargsRsync struct {
	/* Archive field is ignored unless Archive_is_set == true */
	Archive_is_set bool
	Archive        bool
	/* Deletedest field is ignored unless Deletedest_is_set == true */
	Deletedest_is_set bool
	Deletedest        bool
}

/* rsync : synchronize the contents of two directories */
func (g *Guestfs) Rsync(src string, dest string, optargs *OptargsRsync) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("rsync")
	}

	c_src := C.CString(src)
	defer C.free(unsafe.Pointer(c_src))

	c_dest := C.CString(dest)
	defer C.free(unsafe.Pointer(c_dest))
	c_optargs := C.struct_guestfs_rsync_argv{}
	if optargs != nil {
		if optargs.Archive_is_set {
			c_optargs.bitmask |= C.GUESTFS_RSYNC_ARCHIVE_BITMASK
			if optargs.Archive {
				c_optargs.archive = 1
			} else {
				c_optargs.archive = 0
			}
		}
		if optargs.Deletedest_is_set {
			c_optargs.bitmask |= C.GUESTFS_RSYNC_DELETEDEST_BITMASK
			if optargs.Deletedest {
				c_optargs.deletedest = 1
			} else {
				c_optargs.deletedest = 0
			}
		}
	}

	r := C.guestfs_rsync_argv(g.g, c_src, c_dest, &c_optargs)

	if r == -1 {
		return get_error_from_handle(g, "rsync")
	}
	return nil
}

/* Struct carrying optional arguments for Rsync_in */
type OptargsRsync_in struct {
	/* Archive field is ignored unless Archive_is_set == true */
	Archive_is_set bool
	Archive        bool
	/* Deletedest field is ignored unless Deletedest_is_set == true */
	Deletedest_is_set bool
	Deletedest        bool
}

/* rsync_in : synchronize host or remote filesystem with filesystem */
func (g *Guestfs) Rsync_in(remote string, dest string, optargs *OptargsRsync_in) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("rsync_in")
	}

	c_remote := C.CString(remote)
	defer C.free(unsafe.Pointer(c_remote))

	c_dest := C.CString(dest)
	defer C.free(unsafe.Pointer(c_dest))
	c_optargs := C.struct_guestfs_rsync_in_argv{}
	if optargs != nil {
		if optargs.Archive_is_set {
			c_optargs.bitmask |= C.GUESTFS_RSYNC_IN_ARCHIVE_BITMASK
			if optargs.Archive {
				c_optargs.archive = 1
			} else {
				c_optargs.archive = 0
			}
		}
		if optargs.Deletedest_is_set {
			c_optargs.bitmask |= C.GUESTFS_RSYNC_IN_DELETEDEST_BITMASK
			if optargs.Deletedest {
				c_optargs.deletedest = 1
			} else {
				c_optargs.deletedest = 0
			}
		}
	}

	r := C.guestfs_rsync_in_argv(g.g, c_remote, c_dest, &c_optargs)

	if r == -1 {
		return get_error_from_handle(g, "rsync_in")
	}
	return nil
}

/* Struct carrying optional arguments for Rsync_out */
type OptargsRsync_out struct {
	/* Archive field is ignored unless Archive_is_set == true */
	Archive_is_set bool
	Archive        bool
	/* Deletedest field is ignored unless Deletedest_is_set == true */
	Deletedest_is_set bool
	Deletedest        bool
}

/* rsync_out : synchronize filesystem with host or remote filesystem */
func (g *Guestfs) Rsync_out(src string, remote string, optargs *OptargsRsync_out) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("rsync_out")
	}

	c_src := C.CString(src)
	defer C.free(unsafe.Pointer(c_src))

	c_remote := C.CString(remote)
	defer C.free(unsafe.Pointer(c_remote))
	c_optargs := C.struct_guestfs_rsync_out_argv{}
	if optargs != nil {
		if optargs.Archive_is_set {
			c_optargs.bitmask |= C.GUESTFS_RSYNC_OUT_ARCHIVE_BITMASK
			if optargs.Archive {
				c_optargs.archive = 1
			} else {
				c_optargs.archive = 0
			}
		}
		if optargs.Deletedest_is_set {
			c_optargs.bitmask |= C.GUESTFS_RSYNC_OUT_DELETEDEST_BITMASK
			if optargs.Deletedest {
				c_optargs.deletedest = 1
			} else {
				c_optargs.deletedest = 0
			}
		}
	}

	r := C.guestfs_rsync_out_argv(g.g, c_src, c_remote, &c_optargs)

	if r == -1 {
		return get_error_from_handle(g, "rsync_out")
	}
	return nil
}

/* scrub_device : scrub (securely wipe) a device */
func (g *Guestfs) Scrub_device(device string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("scrub_device")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_scrub_device(g.g, c_device)

	if r == -1 {
		return get_error_from_handle(g, "scrub_device")
	}
	return nil
}

/* scrub_file : scrub (securely wipe) a file */
func (g *Guestfs) Scrub_file(file string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("scrub_file")
	}

	c_file := C.CString(file)
	defer C.free(unsafe.Pointer(c_file))

	r := C.guestfs_scrub_file(g.g, c_file)

	if r == -1 {
		return get_error_from_handle(g, "scrub_file")
	}
	return nil
}

/* scrub_freespace : scrub (securely wipe) free space */
func (g *Guestfs) Scrub_freespace(dir string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("scrub_freespace")
	}

	c_dir := C.CString(dir)
	defer C.free(unsafe.Pointer(c_dir))

	r := C.guestfs_scrub_freespace(g.g, c_dir)

	if r == -1 {
		return get_error_from_handle(g, "scrub_freespace")
	}
	return nil
}

/* set_append : add options to kernel command line */
func (g *Guestfs) Set_append(append *string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("set_append")
	}

	var c_append *C.char = nil
	if append != nil {
		c_append = C.CString(*append)
		defer C.free(unsafe.Pointer(c_append))
	}

	r := C.guestfs_set_append(g.g, c_append)

	if r == -1 {
		return get_error_from_handle(g, "set_append")
	}
	return nil
}

/* set_attach_method : set the backend */
func (g *Guestfs) Set_attach_method(backend string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("set_attach_method")
	}

	c_backend := C.CString(backend)
	defer C.free(unsafe.Pointer(c_backend))

	r := C.guestfs_set_attach_method(g.g, c_backend)

	if r == -1 {
		return get_error_from_handle(g, "set_attach_method")
	}
	return nil
}

/* set_autosync : set autosync mode */
func (g *Guestfs) Set_autosync(autosync bool) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("set_autosync")
	}

	var c_autosync C.int
	if autosync {
		c_autosync = 1
	} else {
		c_autosync = 0
	}

	r := C.guestfs_set_autosync(g.g, c_autosync)

	if r == -1 {
		return get_error_from_handle(g, "set_autosync")
	}
	return nil
}

/* set_backend : set the backend */
func (g *Guestfs) Set_backend(backend string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("set_backend")
	}

	c_backend := C.CString(backend)
	defer C.free(unsafe.Pointer(c_backend))

	r := C.guestfs_set_backend(g.g, c_backend)

	if r == -1 {
		return get_error_from_handle(g, "set_backend")
	}
	return nil
}

/* set_backend_setting : set a single per-backend settings string */
func (g *Guestfs) Set_backend_setting(name string, val string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("set_backend_setting")
	}

	c_name := C.CString(name)
	defer C.free(unsafe.Pointer(c_name))

	c_val := C.CString(val)
	defer C.free(unsafe.Pointer(c_val))

	r := C.guestfs_set_backend_setting(g.g, c_name, c_val)

	if r == -1 {
		return get_error_from_handle(g, "set_backend_setting")
	}
	return nil
}

/* set_backend_settings : replace per-backend settings strings */
func (g *Guestfs) Set_backend_settings(settings []string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("set_backend_settings")
	}

	c_settings := arg_string_list(settings)
	defer free_string_list(c_settings)

	r := C.guestfs_set_backend_settings(g.g, c_settings)

	if r == -1 {
		return get_error_from_handle(g, "set_backend_settings")
	}
	return nil
}

/* set_cachedir : set the appliance cache directory */
func (g *Guestfs) Set_cachedir(cachedir *string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("set_cachedir")
	}

	var c_cachedir *C.char = nil
	if cachedir != nil {
		c_cachedir = C.CString(*cachedir)
		defer C.free(unsafe.Pointer(c_cachedir))
	}

	r := C.guestfs_set_cachedir(g.g, c_cachedir)

	if r == -1 {
		return get_error_from_handle(g, "set_cachedir")
	}
	return nil
}

/* set_direct : enable or disable direct appliance mode */
func (g *Guestfs) Set_direct(direct bool) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("set_direct")
	}

	var c_direct C.int
	if direct {
		c_direct = 1
	} else {
		c_direct = 0
	}

	r := C.guestfs_set_direct(g.g, c_direct)

	if r == -1 {
		return get_error_from_handle(g, "set_direct")
	}
	return nil
}

/* Struct carrying optional arguments for Set_e2attrs */
type OptargsSet_e2attrs struct {
	/* Clear field is ignored unless Clear_is_set == true */
	Clear_is_set bool
	Clear        bool
}

/* set_e2attrs : set ext2 file attributes of a file */
func (g *Guestfs) Set_e2attrs(file string, attrs string, optargs *OptargsSet_e2attrs) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("set_e2attrs")
	}

	c_file := C.CString(file)
	defer C.free(unsafe.Pointer(c_file))

	c_attrs := C.CString(attrs)
	defer C.free(unsafe.Pointer(c_attrs))
	c_optargs := C.struct_guestfs_set_e2attrs_argv{}
	if optargs != nil {
		if optargs.Clear_is_set {
			c_optargs.bitmask |= C.GUESTFS_SET_E2ATTRS_CLEAR_BITMASK
			if optargs.Clear {
				c_optargs.clear = 1
			} else {
				c_optargs.clear = 0
			}
		}
	}

	r := C.guestfs_set_e2attrs_argv(g.g, c_file, c_attrs, &c_optargs)

	if r == -1 {
		return get_error_from_handle(g, "set_e2attrs")
	}
	return nil
}

/* set_e2generation : set ext2 file generation of a file */
func (g *Guestfs) Set_e2generation(file string, generation int64) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("set_e2generation")
	}

	c_file := C.CString(file)
	defer C.free(unsafe.Pointer(c_file))

	r := C.guestfs_set_e2generation(g.g, c_file, C.int64_t(generation))

	if r == -1 {
		return get_error_from_handle(g, "set_e2generation")
	}
	return nil
}

/* set_e2label : set the ext2/3/4 filesystem label */
func (g *Guestfs) Set_e2label(device string, label string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("set_e2label")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	c_label := C.CString(label)
	defer C.free(unsafe.Pointer(c_label))

	r := C.guestfs_set_e2label(g.g, c_device, c_label)

	if r == -1 {
		return get_error_from_handle(g, "set_e2label")
	}
	return nil
}

/* set_e2uuid : set the ext2/3/4 filesystem UUID */
func (g *Guestfs) Set_e2uuid(device string, uuid string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("set_e2uuid")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	c_uuid := C.CString(uuid)
	defer C.free(unsafe.Pointer(c_uuid))

	r := C.guestfs_set_e2uuid(g.g, c_device, c_uuid)

	if r == -1 {
		return get_error_from_handle(g, "set_e2uuid")
	}
	return nil
}

/* set_hv : set the hypervisor binary */
func (g *Guestfs) Set_hv(hv string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("set_hv")
	}

	c_hv := C.CString(hv)
	defer C.free(unsafe.Pointer(c_hv))

	r := C.guestfs_set_hv(g.g, c_hv)

	if r == -1 {
		return get_error_from_handle(g, "set_hv")
	}
	return nil
}

/* set_label : set filesystem label */
func (g *Guestfs) Set_label(mountable string, label string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("set_label")
	}

	c_mountable := C.CString(mountable)
	defer C.free(unsafe.Pointer(c_mountable))

	c_label := C.CString(label)
	defer C.free(unsafe.Pointer(c_label))

	r := C.guestfs_set_label(g.g, c_mountable, c_label)

	if r == -1 {
		return get_error_from_handle(g, "set_label")
	}
	return nil
}

/* set_libvirt_requested_credential : pass requested credential back to libvirt */
func (g *Guestfs) Set_libvirt_requested_credential(index int, cred []byte) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("set_libvirt_requested_credential")
	}

	/* string() cast here is apparently safe because
	 *   "Converting a slice of bytes to a string type yields
	 *   a string whose successive bytes are the elements of
	 *   the slice."
	 */
	c_cred := C.CString(string(cred))
	defer C.free(unsafe.Pointer(c_cred))

	r := C.guestfs_set_libvirt_requested_credential(g.g, C.int(index), c_cred, C.size_t(len(cred)))

	if r == -1 {
		return get_error_from_handle(g, "set_libvirt_requested_credential")
	}
	return nil
}

/* set_libvirt_supported_credentials : set libvirt credentials supported by calling program */
func (g *Guestfs) Set_libvirt_supported_credentials(creds []string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("set_libvirt_supported_credentials")
	}

	c_creds := arg_string_list(creds)
	defer free_string_list(c_creds)

	r := C.guestfs_set_libvirt_supported_credentials(g.g, c_creds)

	if r == -1 {
		return get_error_from_handle(g, "set_libvirt_supported_credentials")
	}
	return nil
}

/* set_memsize : set memory allocated to the hypervisor */
func (g *Guestfs) Set_memsize(memsize int) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("set_memsize")
	}

	r := C.guestfs_set_memsize(g.g, C.int(memsize))

	if r == -1 {
		return get_error_from_handle(g, "set_memsize")
	}
	return nil
}

/* set_network : set enable network flag */
func (g *Guestfs) Set_network(network bool) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("set_network")
	}

	var c_network C.int
	if network {
		c_network = 1
	} else {
		c_network = 0
	}

	r := C.guestfs_set_network(g.g, c_network)

	if r == -1 {
		return get_error_from_handle(g, "set_network")
	}
	return nil
}

/* set_path : set the search path */
func (g *Guestfs) Set_path(searchpath *string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("set_path")
	}

	var c_searchpath *C.char = nil
	if searchpath != nil {
		c_searchpath = C.CString(*searchpath)
		defer C.free(unsafe.Pointer(c_searchpath))
	}

	r := C.guestfs_set_path(g.g, c_searchpath)

	if r == -1 {
		return get_error_from_handle(g, "set_path")
	}
	return nil
}

/* set_pgroup : set process group flag */
func (g *Guestfs) Set_pgroup(pgroup bool) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("set_pgroup")
	}

	var c_pgroup C.int
	if pgroup {
		c_pgroup = 1
	} else {
		c_pgroup = 0
	}

	r := C.guestfs_set_pgroup(g.g, c_pgroup)

	if r == -1 {
		return get_error_from_handle(g, "set_pgroup")
	}
	return nil
}

/* set_program : set the program name */
func (g *Guestfs) Set_program(program string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("set_program")
	}

	c_program := C.CString(program)
	defer C.free(unsafe.Pointer(c_program))

	r := C.guestfs_set_program(g.g, c_program)

	if r == -1 {
		return get_error_from_handle(g, "set_program")
	}
	return nil
}

/* set_qemu : set the hypervisor binary (usually qemu) */
func (g *Guestfs) Set_qemu(hv *string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("set_qemu")
	}

	var c_hv *C.char = nil
	if hv != nil {
		c_hv = C.CString(*hv)
		defer C.free(unsafe.Pointer(c_hv))
	}

	r := C.guestfs_set_qemu(g.g, c_hv)

	if r == -1 {
		return get_error_from_handle(g, "set_qemu")
	}
	return nil
}

/* set_recovery_proc : enable or disable the recovery process */
func (g *Guestfs) Set_recovery_proc(recoveryproc bool) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("set_recovery_proc")
	}

	var c_recoveryproc C.int
	if recoveryproc {
		c_recoveryproc = 1
	} else {
		c_recoveryproc = 0
	}

	r := C.guestfs_set_recovery_proc(g.g, c_recoveryproc)

	if r == -1 {
		return get_error_from_handle(g, "set_recovery_proc")
	}
	return nil
}

/* set_selinux : set SELinux enabled or disabled at appliance boot */
func (g *Guestfs) Set_selinux(selinux bool) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("set_selinux")
	}

	var c_selinux C.int
	if selinux {
		c_selinux = 1
	} else {
		c_selinux = 0
	}

	r := C.guestfs_set_selinux(g.g, c_selinux)

	if r == -1 {
		return get_error_from_handle(g, "set_selinux")
	}
	return nil
}

/* set_smp : set number of virtual CPUs in appliance */
func (g *Guestfs) Set_smp(smp int) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("set_smp")
	}

	r := C.guestfs_set_smp(g.g, C.int(smp))

	if r == -1 {
		return get_error_from_handle(g, "set_smp")
	}
	return nil
}

/* set_tmpdir : set the temporary directory */
func (g *Guestfs) Set_tmpdir(tmpdir *string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("set_tmpdir")
	}

	var c_tmpdir *C.char = nil
	if tmpdir != nil {
		c_tmpdir = C.CString(*tmpdir)
		defer C.free(unsafe.Pointer(c_tmpdir))
	}

	r := C.guestfs_set_tmpdir(g.g, c_tmpdir)

	if r == -1 {
		return get_error_from_handle(g, "set_tmpdir")
	}
	return nil
}

/* set_trace : enable or disable command traces */
func (g *Guestfs) Set_trace(trace bool) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("set_trace")
	}

	var c_trace C.int
	if trace {
		c_trace = 1
	} else {
		c_trace = 0
	}

	r := C.guestfs_set_trace(g.g, c_trace)

	if r == -1 {
		return get_error_from_handle(g, "set_trace")
	}
	return nil
}

/* set_uuid : set the filesystem UUID */
func (g *Guestfs) Set_uuid(device string, uuid string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("set_uuid")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	c_uuid := C.CString(uuid)
	defer C.free(unsafe.Pointer(c_uuid))

	r := C.guestfs_set_uuid(g.g, c_device, c_uuid)

	if r == -1 {
		return get_error_from_handle(g, "set_uuid")
	}
	return nil
}

/* set_uuid_random : set a random UUID for the filesystem */
func (g *Guestfs) Set_uuid_random(device string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("set_uuid_random")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_set_uuid_random(g.g, c_device)

	if r == -1 {
		return get_error_from_handle(g, "set_uuid_random")
	}
	return nil
}

/* set_verbose : set verbose mode */
func (g *Guestfs) Set_verbose(verbose bool) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("set_verbose")
	}

	var c_verbose C.int
	if verbose {
		c_verbose = 1
	} else {
		c_verbose = 0
	}

	r := C.guestfs_set_verbose(g.g, c_verbose)

	if r == -1 {
		return get_error_from_handle(g, "set_verbose")
	}
	return nil
}

/* setcon : set SELinux security context */
func (g *Guestfs) Setcon(context string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("setcon")
	}

	c_context := C.CString(context)
	defer C.free(unsafe.Pointer(c_context))

	r := C.guestfs_setcon(g.g, c_context)

	if r == -1 {
		return get_error_from_handle(g, "setcon")
	}
	return nil
}

/* setxattr : set extended attribute of a file or directory */
func (g *Guestfs) Setxattr(xattr string, val string, vallen int, path string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("setxattr")
	}

	c_xattr := C.CString(xattr)
	defer C.free(unsafe.Pointer(c_xattr))

	c_val := C.CString(val)
	defer C.free(unsafe.Pointer(c_val))

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_setxattr(g.g, c_xattr, c_val, C.int(vallen), c_path)

	if r == -1 {
		return get_error_from_handle(g, "setxattr")
	}
	return nil
}

/* sfdisk : create partitions on a block device */
func (g *Guestfs) Sfdisk(device string, cyls int, heads int, sectors int, lines []string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("sfdisk")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	c_lines := arg_string_list(lines)
	defer free_string_list(c_lines)

	r := C.guestfs_sfdisk(g.g, c_device, C.int(cyls), C.int(heads), C.int(sectors), c_lines)

	if r == -1 {
		return get_error_from_handle(g, "sfdisk")
	}
	return nil
}

/* sfdiskM : create partitions on a block device */
func (g *Guestfs) SfdiskM(device string, lines []string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("sfdiskM")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	c_lines := arg_string_list(lines)
	defer free_string_list(c_lines)

	r := C.guestfs_sfdiskM(g.g, c_device, c_lines)

	if r == -1 {
		return get_error_from_handle(g, "sfdiskM")
	}
	return nil
}

/* sfdisk_N : modify a single partition on a block device */
func (g *Guestfs) Sfdisk_N(device string, partnum int, cyls int, heads int, sectors int, line string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("sfdisk_N")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	c_line := C.CString(line)
	defer C.free(unsafe.Pointer(c_line))

	r := C.guestfs_sfdisk_N(g.g, c_device, C.int(partnum), C.int(cyls), C.int(heads), C.int(sectors), c_line)

	if r == -1 {
		return get_error_from_handle(g, "sfdisk_N")
	}
	return nil
}

/* sfdisk_disk_geometry : display the disk geometry from the partition table */
func (g *Guestfs) Sfdisk_disk_geometry(device string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("sfdisk_disk_geometry")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_sfdisk_disk_geometry(g.g, c_device)

	if r == nil {
		return "", get_error_from_handle(g, "sfdisk_disk_geometry")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* sfdisk_kernel_geometry : display the kernel geometry */
func (g *Guestfs) Sfdisk_kernel_geometry(device string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("sfdisk_kernel_geometry")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_sfdisk_kernel_geometry(g.g, c_device)

	if r == nil {
		return "", get_error_from_handle(g, "sfdisk_kernel_geometry")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* sfdisk_l : display the partition table */
func (g *Guestfs) Sfdisk_l(device string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("sfdisk_l")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_sfdisk_l(g.g, c_device)

	if r == nil {
		return "", get_error_from_handle(g, "sfdisk_l")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* sh : run a command via the shell */
func (g *Guestfs) Sh(command string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("sh")
	}

	c_command := C.CString(command)
	defer C.free(unsafe.Pointer(c_command))

	r := C.guestfs_sh(g.g, c_command)

	if r == nil {
		return "", get_error_from_handle(g, "sh")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* sh_lines : run a command via the shell returning lines */
func (g *Guestfs) Sh_lines(command string) ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("sh_lines")
	}

	c_command := C.CString(command)
	defer C.free(unsafe.Pointer(c_command))

	r := C.guestfs_sh_lines(g.g, c_command)

	if r == nil {
		return nil, get_error_from_handle(g, "sh_lines")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* shutdown : shutdown the hypervisor */
func (g *Guestfs) Shutdown() *GuestfsError {
	if g.g == nil {
		return closed_handle_error("shutdown")
	}

	r := C.guestfs_shutdown(g.g)

	if r == -1 {
		return get_error_from_handle(g, "shutdown")
	}
	return nil
}

/* sleep : sleep for some seconds */
func (g *Guestfs) Sleep(secs int) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("sleep")
	}

	r := C.guestfs_sleep(g.g, C.int(secs))

	if r == -1 {
		return get_error_from_handle(g, "sleep")
	}
	return nil
}

/* stat : get file information */
func (g *Guestfs) Stat(path string) (*Stat, *GuestfsError) {
	if g.g == nil {
		return &Stat{}, closed_handle_error("stat")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_stat(g.g, c_path)

	if r == nil {
		return &Stat{}, get_error_from_handle(g, "stat")
	}
	defer C.guestfs_free_stat(r)
	return return_Stat(r), nil
}

/* statns : get file information */
func (g *Guestfs) Statns(path string) (*StatNS, *GuestfsError) {
	if g.g == nil {
		return &StatNS{}, closed_handle_error("statns")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_statns(g.g, c_path)

	if r == nil {
		return &StatNS{}, get_error_from_handle(g, "statns")
	}
	defer C.guestfs_free_statns(r)
	return return_StatNS(r), nil
}

/* statvfs : get file system statistics */
func (g *Guestfs) Statvfs(path string) (*StatVFS, *GuestfsError) {
	if g.g == nil {
		return &StatVFS{}, closed_handle_error("statvfs")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_statvfs(g.g, c_path)

	if r == nil {
		return &StatVFS{}, get_error_from_handle(g, "statvfs")
	}
	defer C.guestfs_free_statvfs(r)
	return return_StatVFS(r), nil
}

/* strings : print the printable strings in a file */
func (g *Guestfs) Strings(path string) ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("strings")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_strings(g.g, c_path)

	if r == nil {
		return nil, get_error_from_handle(g, "strings")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* strings_e : print the printable strings in a file */
func (g *Guestfs) Strings_e(encoding string, path string) ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("strings_e")
	}

	c_encoding := C.CString(encoding)
	defer C.free(unsafe.Pointer(c_encoding))

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_strings_e(g.g, c_encoding, c_path)

	if r == nil {
		return nil, get_error_from_handle(g, "strings_e")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* swapoff_device : disable swap on device */
func (g *Guestfs) Swapoff_device(device string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("swapoff_device")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_swapoff_device(g.g, c_device)

	if r == -1 {
		return get_error_from_handle(g, "swapoff_device")
	}
	return nil
}

/* swapoff_file : disable swap on file */
func (g *Guestfs) Swapoff_file(file string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("swapoff_file")
	}

	c_file := C.CString(file)
	defer C.free(unsafe.Pointer(c_file))

	r := C.guestfs_swapoff_file(g.g, c_file)

	if r == -1 {
		return get_error_from_handle(g, "swapoff_file")
	}
	return nil
}

/* swapoff_label : disable swap on labeled swap partition */
func (g *Guestfs) Swapoff_label(label string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("swapoff_label")
	}

	c_label := C.CString(label)
	defer C.free(unsafe.Pointer(c_label))

	r := C.guestfs_swapoff_label(g.g, c_label)

	if r == -1 {
		return get_error_from_handle(g, "swapoff_label")
	}
	return nil
}

/* swapoff_uuid : disable swap on swap partition by UUID */
func (g *Guestfs) Swapoff_uuid(uuid string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("swapoff_uuid")
	}

	c_uuid := C.CString(uuid)
	defer C.free(unsafe.Pointer(c_uuid))

	r := C.guestfs_swapoff_uuid(g.g, c_uuid)

	if r == -1 {
		return get_error_from_handle(g, "swapoff_uuid")
	}
	return nil
}

/* swapon_device : enable swap on device */
func (g *Guestfs) Swapon_device(device string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("swapon_device")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_swapon_device(g.g, c_device)

	if r == -1 {
		return get_error_from_handle(g, "swapon_device")
	}
	return nil
}

/* swapon_file : enable swap on file */
func (g *Guestfs) Swapon_file(file string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("swapon_file")
	}

	c_file := C.CString(file)
	defer C.free(unsafe.Pointer(c_file))

	r := C.guestfs_swapon_file(g.g, c_file)

	if r == -1 {
		return get_error_from_handle(g, "swapon_file")
	}
	return nil
}

/* swapon_label : enable swap on labeled swap partition */
func (g *Guestfs) Swapon_label(label string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("swapon_label")
	}

	c_label := C.CString(label)
	defer C.free(unsafe.Pointer(c_label))

	r := C.guestfs_swapon_label(g.g, c_label)

	if r == -1 {
		return get_error_from_handle(g, "swapon_label")
	}
	return nil
}

/* swapon_uuid : enable swap on swap partition by UUID */
func (g *Guestfs) Swapon_uuid(uuid string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("swapon_uuid")
	}

	c_uuid := C.CString(uuid)
	defer C.free(unsafe.Pointer(c_uuid))

	r := C.guestfs_swapon_uuid(g.g, c_uuid)

	if r == -1 {
		return get_error_from_handle(g, "swapon_uuid")
	}
	return nil
}

/* sync : sync disks, writes are flushed through to the disk image */
func (g *Guestfs) Sync() *GuestfsError {
	if g.g == nil {
		return closed_handle_error("sync")
	}

	r := C.guestfs_sync(g.g)

	if r == -1 {
		return get_error_from_handle(g, "sync")
	}
	return nil
}

/* Struct carrying optional arguments for Syslinux */
type OptargsSyslinux struct {
	/* Directory field is ignored unless Directory_is_set == true */
	Directory_is_set bool
	Directory        string
}

/* syslinux : install the SYSLINUX bootloader */
func (g *Guestfs) Syslinux(device string, optargs *OptargsSyslinux) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("syslinux")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))
	c_optargs := C.struct_guestfs_syslinux_argv{}
	if optargs != nil {
		if optargs.Directory_is_set {
			c_optargs.bitmask |= C.GUESTFS_SYSLINUX_DIRECTORY_BITMASK
			c_optargs.directory = C.CString(optargs.Directory)
			defer C.free(unsafe.Pointer(c_optargs.directory))
		}
	}

	r := C.guestfs_syslinux_argv(g.g, c_device, &c_optargs)

	if r == -1 {
		return get_error_from_handle(g, "syslinux")
	}
	return nil
}

/* tail : return last 10 lines of a file */
func (g *Guestfs) Tail(path string) ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("tail")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_tail(g.g, c_path)

	if r == nil {
		return nil, get_error_from_handle(g, "tail")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* tail_n : return last N lines of a file */
func (g *Guestfs) Tail_n(nrlines int, path string) ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("tail_n")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_tail_n(g.g, C.int(nrlines), c_path)

	if r == nil {
		return nil, get_error_from_handle(g, "tail_n")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* Struct carrying optional arguments for Tar_in */
type OptargsTar_in struct {
	/* Compress field is ignored unless Compress_is_set == true */
	Compress_is_set bool
	Compress        string
	/* Xattrs field is ignored unless Xattrs_is_set == true */
	Xattrs_is_set bool
	Xattrs        bool
	/* Selinux field is ignored unless Selinux_is_set == true */
	Selinux_is_set bool
	Selinux        bool
	/* Acls field is ignored unless Acls_is_set == true */
	Acls_is_set bool
	Acls        bool
}

/* tar_in : unpack tarfile to directory */
func (g *Guestfs) Tar_in(tarfile string, directory string, optargs *OptargsTar_in) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("tar_in")
	}

	c_tarfile := C.CString(tarfile)
	defer C.free(unsafe.Pointer(c_tarfile))

	c_directory := C.CString(directory)
	defer C.free(unsafe.Pointer(c_directory))
	c_optargs := C.struct_guestfs_tar_in_opts_argv{}
	if optargs != nil {
		if optargs.Compress_is_set {
			c_optargs.bitmask |= C.GUESTFS_TAR_IN_OPTS_COMPRESS_BITMASK
			c_optargs.compress = C.CString(optargs.Compress)
			defer C.free(unsafe.Pointer(c_optargs.compress))
		}
		if optargs.Xattrs_is_set {
			c_optargs.bitmask |= C.GUESTFS_TAR_IN_OPTS_XATTRS_BITMASK
			if optargs.Xattrs {
				c_optargs.xattrs = 1
			} else {
				c_optargs.xattrs = 0
			}
		}
		if optargs.Selinux_is_set {
			c_optargs.bitmask |= C.GUESTFS_TAR_IN_OPTS_SELINUX_BITMASK
			if optargs.Selinux {
				c_optargs.selinux = 1
			} else {
				c_optargs.selinux = 0
			}
		}
		if optargs.Acls_is_set {
			c_optargs.bitmask |= C.GUESTFS_TAR_IN_OPTS_ACLS_BITMASK
			if optargs.Acls {
				c_optargs.acls = 1
			} else {
				c_optargs.acls = 0
			}
		}
	}

	r := C.guestfs_tar_in_opts_argv(g.g, c_tarfile, c_directory, &c_optargs)

	if r == -1 {
		return get_error_from_handle(g, "tar_in")
	}
	return nil
}

/* Struct carrying optional arguments for Tar_out */
type OptargsTar_out struct {
	/* Compress field is ignored unless Compress_is_set == true */
	Compress_is_set bool
	Compress        string
	/* Numericowner field is ignored unless Numericowner_is_set == true */
	Numericowner_is_set bool
	Numericowner        bool
	/* Excludes field is ignored unless Excludes_is_set == true */
	Excludes_is_set bool
	Excludes        []string
	/* Xattrs field is ignored unless Xattrs_is_set == true */
	Xattrs_is_set bool
	Xattrs        bool
	/* Selinux field is ignored unless Selinux_is_set == true */
	Selinux_is_set bool
	Selinux        bool
	/* Acls field is ignored unless Acls_is_set == true */
	Acls_is_set bool
	Acls        bool
}

/* tar_out : pack directory into tarfile */
func (g *Guestfs) Tar_out(directory string, tarfile string, optargs *OptargsTar_out) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("tar_out")
	}

	c_directory := C.CString(directory)
	defer C.free(unsafe.Pointer(c_directory))

	c_tarfile := C.CString(tarfile)
	defer C.free(unsafe.Pointer(c_tarfile))
	c_optargs := C.struct_guestfs_tar_out_opts_argv{}
	if optargs != nil {
		if optargs.Compress_is_set {
			c_optargs.bitmask |= C.GUESTFS_TAR_OUT_OPTS_COMPRESS_BITMASK
			c_optargs.compress = C.CString(optargs.Compress)
			defer C.free(unsafe.Pointer(c_optargs.compress))
		}
		if optargs.Numericowner_is_set {
			c_optargs.bitmask |= C.GUESTFS_TAR_OUT_OPTS_NUMERICOWNER_BITMASK
			if optargs.Numericowner {
				c_optargs.numericowner = 1
			} else {
				c_optargs.numericowner = 0
			}
		}
		if optargs.Excludes_is_set {
			c_optargs.bitmask |= C.GUESTFS_TAR_OUT_OPTS_EXCLUDES_BITMASK
			c_optargs.excludes = arg_string_list(optargs.Excludes)
			defer free_string_list(c_optargs.excludes)
		}
		if optargs.Xattrs_is_set {
			c_optargs.bitmask |= C.GUESTFS_TAR_OUT_OPTS_XATTRS_BITMASK
			if optargs.Xattrs {
				c_optargs.xattrs = 1
			} else {
				c_optargs.xattrs = 0
			}
		}
		if optargs.Selinux_is_set {
			c_optargs.bitmask |= C.GUESTFS_TAR_OUT_OPTS_SELINUX_BITMASK
			if optargs.Selinux {
				c_optargs.selinux = 1
			} else {
				c_optargs.selinux = 0
			}
		}
		if optargs.Acls_is_set {
			c_optargs.bitmask |= C.GUESTFS_TAR_OUT_OPTS_ACLS_BITMASK
			if optargs.Acls {
				c_optargs.acls = 1
			} else {
				c_optargs.acls = 0
			}
		}
	}

	r := C.guestfs_tar_out_opts_argv(g.g, c_directory, c_tarfile, &c_optargs)

	if r == -1 {
		return get_error_from_handle(g, "tar_out")
	}
	return nil
}

/* tgz_in : unpack compressed tarball to directory */
func (g *Guestfs) Tgz_in(tarball string, directory string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("tgz_in")
	}

	c_tarball := C.CString(tarball)
	defer C.free(unsafe.Pointer(c_tarball))

	c_directory := C.CString(directory)
	defer C.free(unsafe.Pointer(c_directory))

	r := C.guestfs_tgz_in(g.g, c_tarball, c_directory)

	if r == -1 {
		return get_error_from_handle(g, "tgz_in")
	}
	return nil
}

/* tgz_out : pack directory into compressed tarball */
func (g *Guestfs) Tgz_out(directory string, tarball string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("tgz_out")
	}

	c_directory := C.CString(directory)
	defer C.free(unsafe.Pointer(c_directory))

	c_tarball := C.CString(tarball)
	defer C.free(unsafe.Pointer(c_tarball))

	r := C.guestfs_tgz_out(g.g, c_directory, c_tarball)

	if r == -1 {
		return get_error_from_handle(g, "tgz_out")
	}
	return nil
}

/* touch : update file timestamps or create a new file */
func (g *Guestfs) Touch(path string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("touch")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_touch(g.g, c_path)

	if r == -1 {
		return get_error_from_handle(g, "touch")
	}
	return nil
}

/* truncate : truncate a file to zero size */
func (g *Guestfs) Truncate(path string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("truncate")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_truncate(g.g, c_path)

	if r == -1 {
		return get_error_from_handle(g, "truncate")
	}
	return nil
}

/* truncate_size : truncate a file to a particular size */
func (g *Guestfs) Truncate_size(path string, size int64) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("truncate_size")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_truncate_size(g.g, c_path, C.int64_t(size))

	if r == -1 {
		return get_error_from_handle(g, "truncate_size")
	}
	return nil
}

/* Struct carrying optional arguments for Tune2fs */
type OptargsTune2fs struct {
	/* Force field is ignored unless Force_is_set == true */
	Force_is_set bool
	Force        bool
	/* Maxmountcount field is ignored unless Maxmountcount_is_set == true */
	Maxmountcount_is_set bool
	Maxmountcount        int
	/* Mountcount field is ignored unless Mountcount_is_set == true */
	Mountcount_is_set bool
	Mountcount        int
	/* Errorbehavior field is ignored unless Errorbehavior_is_set == true */
	Errorbehavior_is_set bool
	Errorbehavior        string
	/* Group field is ignored unless Group_is_set == true */
	Group_is_set bool
	Group        int64
	/* Intervalbetweenchecks field is ignored unless Intervalbetweenchecks_is_set == true */
	Intervalbetweenchecks_is_set bool
	Intervalbetweenchecks        int
	/* Reservedblockspercentage field is ignored unless Reservedblockspercentage_is_set == true */
	Reservedblockspercentage_is_set bool
	Reservedblockspercentage        int
	/* Lastmounteddirectory field is ignored unless Lastmounteddirectory_is_set == true */
	Lastmounteddirectory_is_set bool
	Lastmounteddirectory        string
	/* Reservedblockscount field is ignored unless Reservedblockscount_is_set == true */
	Reservedblockscount_is_set bool
	Reservedblockscount        int64
	/* User field is ignored unless User_is_set == true */
	User_is_set bool
	User        int64
}

/* tune2fs : adjust ext2/ext3/ext4 filesystem parameters */
func (g *Guestfs) Tune2fs(device string, optargs *OptargsTune2fs) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("tune2fs")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))
	c_optargs := C.struct_guestfs_tune2fs_argv{}
	if optargs != nil {
		if optargs.Force_is_set {
			c_optargs.bitmask |= C.GUESTFS_TUNE2FS_FORCE_BITMASK
			if optargs.Force {
				c_optargs.force = 1
			} else {
				c_optargs.force = 0
			}
		}
		if optargs.Maxmountcount_is_set {
			c_optargs.bitmask |= C.GUESTFS_TUNE2FS_MAXMOUNTCOUNT_BITMASK
			c_optargs.maxmountcount = C.int(optargs.Maxmountcount)
		}
		if optargs.Mountcount_is_set {
			c_optargs.bitmask |= C.GUESTFS_TUNE2FS_MOUNTCOUNT_BITMASK
			c_optargs.mountcount = C.int(optargs.Mountcount)
		}
		if optargs.Errorbehavior_is_set {
			c_optargs.bitmask |= C.GUESTFS_TUNE2FS_ERRORBEHAVIOR_BITMASK
			c_optargs.errorbehavior = C.CString(optargs.Errorbehavior)
			defer C.free(unsafe.Pointer(c_optargs.errorbehavior))
		}
		if optargs.Group_is_set {
			c_optargs.bitmask |= C.GUESTFS_TUNE2FS_GROUP_BITMASK
			c_optargs.group = C.int64_t(optargs.Group)
		}
		if optargs.Intervalbetweenchecks_is_set {
			c_optargs.bitmask |= C.GUESTFS_TUNE2FS_INTERVALBETWEENCHECKS_BITMASK
			c_optargs.intervalbetweenchecks = C.int(optargs.Intervalbetweenchecks)
		}
		if optargs.Reservedblockspercentage_is_set {
			c_optargs.bitmask |= C.GUESTFS_TUNE2FS_RESERVEDBLOCKSPERCENTAGE_BITMASK
			c_optargs.reservedblockspercentage = C.int(optargs.Reservedblockspercentage)
		}
		if optargs.Lastmounteddirectory_is_set {
			c_optargs.bitmask |= C.GUESTFS_TUNE2FS_LASTMOUNTEDDIRECTORY_BITMASK
			c_optargs.lastmounteddirectory = C.CString(optargs.Lastmounteddirectory)
			defer C.free(unsafe.Pointer(c_optargs.lastmounteddirectory))
		}
		if optargs.Reservedblockscount_is_set {
			c_optargs.bitmask |= C.GUESTFS_TUNE2FS_RESERVEDBLOCKSCOUNT_BITMASK
			c_optargs.reservedblockscount = C.int64_t(optargs.Reservedblockscount)
		}
		if optargs.User_is_set {
			c_optargs.bitmask |= C.GUESTFS_TUNE2FS_USER_BITMASK
			c_optargs.user = C.int64_t(optargs.User)
		}
	}

	r := C.guestfs_tune2fs_argv(g.g, c_device, &c_optargs)

	if r == -1 {
		return get_error_from_handle(g, "tune2fs")
	}
	return nil
}

/* tune2fs_l : get ext2/ext3/ext4 superblock details */
func (g *Guestfs) Tune2fs_l(device string) (map[string]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("tune2fs_l")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_tune2fs_l(g.g, c_device)

	if r == nil {
		return nil, get_error_from_handle(g, "tune2fs_l")
	}
	defer free_string_list(r)
	return return_hashtable(r), nil
}

/* txz_in : unpack compressed tarball to directory */
func (g *Guestfs) Txz_in(tarball string, directory string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("txz_in")
	}

	c_tarball := C.CString(tarball)
	defer C.free(unsafe.Pointer(c_tarball))

	c_directory := C.CString(directory)
	defer C.free(unsafe.Pointer(c_directory))

	r := C.guestfs_txz_in(g.g, c_tarball, c_directory)

	if r == -1 {
		return get_error_from_handle(g, "txz_in")
	}
	return nil
}

/* txz_out : pack directory into compressed tarball */
func (g *Guestfs) Txz_out(directory string, tarball string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("txz_out")
	}

	c_directory := C.CString(directory)
	defer C.free(unsafe.Pointer(c_directory))

	c_tarball := C.CString(tarball)
	defer C.free(unsafe.Pointer(c_tarball))

	r := C.guestfs_txz_out(g.g, c_directory, c_tarball)

	if r == -1 {
		return get_error_from_handle(g, "txz_out")
	}
	return nil
}

/* umask : set file mode creation mask (umask) */
func (g *Guestfs) Umask(mask int) (int, *GuestfsError) {
	if g.g == nil {
		return 0, closed_handle_error("umask")
	}

	r := C.guestfs_umask(g.g, C.int(mask))

	if r == -1 {
		return 0, get_error_from_handle(g, "umask")
	}
	return int(r), nil
}

/* Struct carrying optional arguments for Umount */
type OptargsUmount struct {
	/* Force field is ignored unless Force_is_set == true */
	Force_is_set bool
	Force        bool
	/* Lazyunmount field is ignored unless Lazyunmount_is_set == true */
	Lazyunmount_is_set bool
	Lazyunmount        bool
}

/* umount : unmount a filesystem */
func (g *Guestfs) Umount(pathordevice string, optargs *OptargsUmount) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("umount")
	}

	c_pathordevice := C.CString(pathordevice)
	defer C.free(unsafe.Pointer(c_pathordevice))
	c_optargs := C.struct_guestfs_umount_opts_argv{}
	if optargs != nil {
		if optargs.Force_is_set {
			c_optargs.bitmask |= C.GUESTFS_UMOUNT_OPTS_FORCE_BITMASK
			if optargs.Force {
				c_optargs.force = 1
			} else {
				c_optargs.force = 0
			}
		}
		if optargs.Lazyunmount_is_set {
			c_optargs.bitmask |= C.GUESTFS_UMOUNT_OPTS_LAZYUNMOUNT_BITMASK
			if optargs.Lazyunmount {
				c_optargs.lazyunmount = 1
			} else {
				c_optargs.lazyunmount = 0
			}
		}
	}

	r := C.guestfs_umount_opts_argv(g.g, c_pathordevice, &c_optargs)

	if r == -1 {
		return get_error_from_handle(g, "umount")
	}
	return nil
}

/* umount_all : unmount all filesystems */
func (g *Guestfs) Umount_all() *GuestfsError {
	if g.g == nil {
		return closed_handle_error("umount_all")
	}

	r := C.guestfs_umount_all(g.g)

	if r == -1 {
		return get_error_from_handle(g, "umount_all")
	}
	return nil
}

/* Struct carrying optional arguments for Umount_local */
type OptargsUmount_local struct {
	/* Retry field is ignored unless Retry_is_set == true */
	Retry_is_set bool
	Retry        bool
}

/* umount_local : unmount a locally mounted filesystem */
func (g *Guestfs) Umount_local(optargs *OptargsUmount_local) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("umount_local")
	}
	c_optargs := C.struct_guestfs_umount_local_argv{}
	if optargs != nil {
		if optargs.Retry_is_set {
			c_optargs.bitmask |= C.GUESTFS_UMOUNT_LOCAL_RETRY_BITMASK
			if optargs.Retry {
				c_optargs.retry = 1
			} else {
				c_optargs.retry = 0
			}
		}
	}

	r := C.guestfs_umount_local_argv(g.g, &c_optargs)

	if r == -1 {
		return get_error_from_handle(g, "umount_local")
	}
	return nil
}

/* upload : upload a file from the local machine */
func (g *Guestfs) Upload(filename string, remotefilename string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("upload")
	}

	c_filename := C.CString(filename)
	defer C.free(unsafe.Pointer(c_filename))

	c_remotefilename := C.CString(remotefilename)
	defer C.free(unsafe.Pointer(c_remotefilename))

	r := C.guestfs_upload(g.g, c_filename, c_remotefilename)

	if r == -1 {
		return get_error_from_handle(g, "upload")
	}
	return nil
}

/* upload_offset : upload a file from the local machine with offset */
func (g *Guestfs) Upload_offset(filename string, remotefilename string, offset int64) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("upload_offset")
	}

	c_filename := C.CString(filename)
	defer C.free(unsafe.Pointer(c_filename))

	c_remotefilename := C.CString(remotefilename)
	defer C.free(unsafe.Pointer(c_remotefilename))

	r := C.guestfs_upload_offset(g.g, c_filename, c_remotefilename, C.int64_t(offset))

	if r == -1 {
		return get_error_from_handle(g, "upload_offset")
	}
	return nil
}

/* user_cancel : cancel the current upload or download operation */
func (g *Guestfs) User_cancel() *GuestfsError {
	if g.g == nil {
		return closed_handle_error("user_cancel")
	}

	r := C.guestfs_user_cancel(g.g)

	if r == -1 {
		return get_error_from_handle(g, "user_cancel")
	}
	return nil
}

/* utimens : set timestamp of a file with nanosecond precision */
func (g *Guestfs) Utimens(path string, atsecs int64, atnsecs int64, mtsecs int64, mtnsecs int64) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("utimens")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_utimens(g.g, c_path, C.int64_t(atsecs), C.int64_t(atnsecs), C.int64_t(mtsecs), C.int64_t(mtnsecs))

	if r == -1 {
		return get_error_from_handle(g, "utimens")
	}
	return nil
}

/* utsname : appliance kernel version */
func (g *Guestfs) Utsname() (*UTSName, *GuestfsError) {
	if g.g == nil {
		return &UTSName{}, closed_handle_error("utsname")
	}

	r := C.guestfs_utsname(g.g)

	if r == nil {
		return &UTSName{}, get_error_from_handle(g, "utsname")
	}
	defer C.guestfs_free_utsname(r)
	return return_UTSName(r), nil
}

/* version : get the library version number */
func (g *Guestfs) Version() (*Version, *GuestfsError) {
	if g.g == nil {
		return &Version{}, closed_handle_error("version")
	}

	r := C.guestfs_version(g.g)

	if r == nil {
		return &Version{}, get_error_from_handle(g, "version")
	}
	defer C.guestfs_free_version(r)
	return return_Version(r), nil
}

/* vfs_label : get the filesystem label */
func (g *Guestfs) Vfs_label(mountable string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("vfs_label")
	}

	c_mountable := C.CString(mountable)
	defer C.free(unsafe.Pointer(c_mountable))

	r := C.guestfs_vfs_label(g.g, c_mountable)

	if r == nil {
		return "", get_error_from_handle(g, "vfs_label")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* vfs_type : get the Linux VFS type corresponding to a mounted device */
func (g *Guestfs) Vfs_type(mountable string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("vfs_type")
	}

	c_mountable := C.CString(mountable)
	defer C.free(unsafe.Pointer(c_mountable))

	r := C.guestfs_vfs_type(g.g, c_mountable)

	if r == nil {
		return "", get_error_from_handle(g, "vfs_type")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* vfs_uuid : get the filesystem UUID */
func (g *Guestfs) Vfs_uuid(mountable string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("vfs_uuid")
	}

	c_mountable := C.CString(mountable)
	defer C.free(unsafe.Pointer(c_mountable))

	r := C.guestfs_vfs_uuid(g.g, c_mountable)

	if r == nil {
		return "", get_error_from_handle(g, "vfs_uuid")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* vg_activate : activate or deactivate some volume groups */
func (g *Guestfs) Vg_activate(activate bool, volgroups []string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("vg_activate")
	}

	var c_activate C.int
	if activate {
		c_activate = 1
	} else {
		c_activate = 0
	}

	c_volgroups := arg_string_list(volgroups)
	defer free_string_list(c_volgroups)

	r := C.guestfs_vg_activate(g.g, c_activate, c_volgroups)

	if r == -1 {
		return get_error_from_handle(g, "vg_activate")
	}
	return nil
}

/* vg_activate_all : activate or deactivate all volume groups */
func (g *Guestfs) Vg_activate_all(activate bool) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("vg_activate_all")
	}

	var c_activate C.int
	if activate {
		c_activate = 1
	} else {
		c_activate = 0
	}

	r := C.guestfs_vg_activate_all(g.g, c_activate)

	if r == -1 {
		return get_error_from_handle(g, "vg_activate_all")
	}
	return nil
}

/* vgchange_uuid : generate a new random UUID for a volume group */
func (g *Guestfs) Vgchange_uuid(vg string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("vgchange_uuid")
	}

	c_vg := C.CString(vg)
	defer C.free(unsafe.Pointer(c_vg))

	r := C.guestfs_vgchange_uuid(g.g, c_vg)

	if r == -1 {
		return get_error_from_handle(g, "vgchange_uuid")
	}
	return nil
}

/* vgchange_uuid_all : generate new random UUIDs for all volume groups */
func (g *Guestfs) Vgchange_uuid_all() *GuestfsError {
	if g.g == nil {
		return closed_handle_error("vgchange_uuid_all")
	}

	r := C.guestfs_vgchange_uuid_all(g.g)

	if r == -1 {
		return get_error_from_handle(g, "vgchange_uuid_all")
	}
	return nil
}

/* vgcreate : create an LVM volume group */
func (g *Guestfs) Vgcreate(volgroup string, physvols []string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("vgcreate")
	}

	c_volgroup := C.CString(volgroup)
	defer C.free(unsafe.Pointer(c_volgroup))

	c_physvols := arg_string_list(physvols)
	defer free_string_list(c_physvols)

	r := C.guestfs_vgcreate(g.g, c_volgroup, c_physvols)

	if r == -1 {
		return get_error_from_handle(g, "vgcreate")
	}
	return nil
}

/* vglvuuids : get the LV UUIDs of all LVs in the volume group */
func (g *Guestfs) Vglvuuids(vgname string) ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("vglvuuids")
	}

	c_vgname := C.CString(vgname)
	defer C.free(unsafe.Pointer(c_vgname))

	r := C.guestfs_vglvuuids(g.g, c_vgname)

	if r == nil {
		return nil, get_error_from_handle(g, "vglvuuids")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* vgmeta : get volume group metadata */
func (g *Guestfs) Vgmeta(vgname string) ([]byte, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("vgmeta")
	}

	c_vgname := C.CString(vgname)
	defer C.free(unsafe.Pointer(c_vgname))

	var size C.size_t

	r := C.guestfs_vgmeta(g.g, c_vgname, &size)

	if r == nil {
		return nil, get_error_from_handle(g, "vgmeta")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoBytes(unsafe.Pointer(r), C.int(size)), nil
}

/* vgpvuuids : get the PV UUIDs containing the volume group */
func (g *Guestfs) Vgpvuuids(vgname string) ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("vgpvuuids")
	}

	c_vgname := C.CString(vgname)
	defer C.free(unsafe.Pointer(c_vgname))

	r := C.guestfs_vgpvuuids(g.g, c_vgname)

	if r == nil {
		return nil, get_error_from_handle(g, "vgpvuuids")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* vgremove : remove an LVM volume group */
func (g *Guestfs) Vgremove(vgname string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("vgremove")
	}

	c_vgname := C.CString(vgname)
	defer C.free(unsafe.Pointer(c_vgname))

	r := C.guestfs_vgremove(g.g, c_vgname)

	if r == -1 {
		return get_error_from_handle(g, "vgremove")
	}
	return nil
}

/* vgrename : rename an LVM volume group */
func (g *Guestfs) Vgrename(volgroup string, newvolgroup string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("vgrename")
	}

	c_volgroup := C.CString(volgroup)
	defer C.free(unsafe.Pointer(c_volgroup))

	c_newvolgroup := C.CString(newvolgroup)
	defer C.free(unsafe.Pointer(c_newvolgroup))

	r := C.guestfs_vgrename(g.g, c_volgroup, c_newvolgroup)

	if r == -1 {
		return get_error_from_handle(g, "vgrename")
	}
	return nil
}

/* vgs : list the LVM volume groups (VGs) */
func (g *Guestfs) Vgs() ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("vgs")
	}

	r := C.guestfs_vgs(g.g)

	if r == nil {
		return nil, get_error_from_handle(g, "vgs")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* vgs_full : list the LVM volume groups (VGs) */
func (g *Guestfs) Vgs_full() (*[]VG, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("vgs_full")
	}

	r := C.guestfs_vgs_full(g.g)

	if r == nil {
		return nil, get_error_from_handle(g, "vgs_full")
	}
	defer C.guestfs_free_lvm_vg_list(r)
	return return_VG_list(r), nil
}

/* vgscan : rescan for LVM physical volumes, volume groups and logical volumes */
func (g *Guestfs) Vgscan() *GuestfsError {
	if g.g == nil {
		return closed_handle_error("vgscan")
	}

	r := C.guestfs_vgscan(g.g)

	if r == -1 {
		return get_error_from_handle(g, "vgscan")
	}
	return nil
}

/* vguuid : get the UUID of a volume group */
func (g *Guestfs) Vguuid(vgname string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("vguuid")
	}

	c_vgname := C.CString(vgname)
	defer C.free(unsafe.Pointer(c_vgname))

	r := C.guestfs_vguuid(g.g, c_vgname)

	if r == nil {
		return "", get_error_from_handle(g, "vguuid")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* wait_ready : wait until the hypervisor launches (no op) */
func (g *Guestfs) Wait_ready() *GuestfsError {
	if g.g == nil {
		return closed_handle_error("wait_ready")
	}

	r := C.guestfs_wait_ready(g.g)

	if r == -1 {
		return get_error_from_handle(g, "wait_ready")
	}
	return nil
}

/* wc_c : count characters in a file */
func (g *Guestfs) Wc_c(path string) (int, *GuestfsError) {
	if g.g == nil {
		return 0, closed_handle_error("wc_c")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_wc_c(g.g, c_path)

	if r == -1 {
		return 0, get_error_from_handle(g, "wc_c")
	}
	return int(r), nil
}

/* wc_l : count lines in a file */
func (g *Guestfs) Wc_l(path string) (int, *GuestfsError) {
	if g.g == nil {
		return 0, closed_handle_error("wc_l")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_wc_l(g.g, c_path)

	if r == -1 {
		return 0, get_error_from_handle(g, "wc_l")
	}
	return int(r), nil
}

/* wc_w : count words in a file */
func (g *Guestfs) Wc_w(path string) (int, *GuestfsError) {
	if g.g == nil {
		return 0, closed_handle_error("wc_w")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_wc_w(g.g, c_path)

	if r == -1 {
		return 0, get_error_from_handle(g, "wc_w")
	}
	return int(r), nil
}

/* wipefs : wipe a filesystem signature from a device */
func (g *Guestfs) Wipefs(device string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("wipefs")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_wipefs(g.g, c_device)

	if r == -1 {
		return get_error_from_handle(g, "wipefs")
	}
	return nil
}

/* write : create a new file */
func (g *Guestfs) Write(path string, content []byte) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("write")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	/* string() cast here is apparently safe because
	 *   "Converting a slice of bytes to a string type yields
	 *   a string whose successive bytes are the elements of
	 *   the slice."
	 */
	c_content := C.CString(string(content))
	defer C.free(unsafe.Pointer(c_content))

	r := C.guestfs_write(g.g, c_path, c_content, C.size_t(len(content)))

	if r == -1 {
		return get_error_from_handle(g, "write")
	}
	return nil
}

/* write_append : append content to end of file */
func (g *Guestfs) Write_append(path string, content []byte) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("write_append")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	/* string() cast here is apparently safe because
	 *   "Converting a slice of bytes to a string type yields
	 *   a string whose successive bytes are the elements of
	 *   the slice."
	 */
	c_content := C.CString(string(content))
	defer C.free(unsafe.Pointer(c_content))

	r := C.guestfs_write_append(g.g, c_path, c_content, C.size_t(len(content)))

	if r == -1 {
		return get_error_from_handle(g, "write_append")
	}
	return nil
}

/* write_file : create a file */
func (g *Guestfs) Write_file(path string, content string, size int) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("write_file")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	c_content := C.CString(content)
	defer C.free(unsafe.Pointer(c_content))

	r := C.guestfs_write_file(g.g, c_path, c_content, C.int(size))

	if r == -1 {
		return get_error_from_handle(g, "write_file")
	}
	return nil
}

/* Struct carrying optional arguments for Xfs_admin */
type OptargsXfs_admin struct {
	/* Extunwritten field is ignored unless Extunwritten_is_set == true */
	Extunwritten_is_set bool
	Extunwritten        bool
	/* Imgfile field is ignored unless Imgfile_is_set == true */
	Imgfile_is_set bool
	Imgfile        bool
	/* V2log field is ignored unless V2log_is_set == true */
	V2log_is_set bool
	V2log        bool
	/* Projid32bit field is ignored unless Projid32bit_is_set == true */
	Projid32bit_is_set bool
	Projid32bit        bool
	/* Lazycounter field is ignored unless Lazycounter_is_set == true */
	Lazycounter_is_set bool
	Lazycounter        bool
	/* Label field is ignored unless Label_is_set == true */
	Label_is_set bool
	Label        string
	/* Uuid field is ignored unless Uuid_is_set == true */
	Uuid_is_set bool
	Uuid        string
}

/* xfs_admin : change parameters of an XFS filesystem */
func (g *Guestfs) Xfs_admin(device string, optargs *OptargsXfs_admin) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("xfs_admin")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))
	c_optargs := C.struct_guestfs_xfs_admin_argv{}
	if optargs != nil {
		if optargs.Extunwritten_is_set {
			c_optargs.bitmask |= C.GUESTFS_XFS_ADMIN_EXTUNWRITTEN_BITMASK
			if optargs.Extunwritten {
				c_optargs.extunwritten = 1
			} else {
				c_optargs.extunwritten = 0
			}
		}
		if optargs.Imgfile_is_set {
			c_optargs.bitmask |= C.GUESTFS_XFS_ADMIN_IMGFILE_BITMASK
			if optargs.Imgfile {
				c_optargs.imgfile = 1
			} else {
				c_optargs.imgfile = 0
			}
		}
		if optargs.V2log_is_set {
			c_optargs.bitmask |= C.GUESTFS_XFS_ADMIN_V2LOG_BITMASK
			if optargs.V2log {
				c_optargs.v2log = 1
			} else {
				c_optargs.v2log = 0
			}
		}
		if optargs.Projid32bit_is_set {
			c_optargs.bitmask |= C.GUESTFS_XFS_ADMIN_PROJID32BIT_BITMASK
			if optargs.Projid32bit {
				c_optargs.projid32bit = 1
			} else {
				c_optargs.projid32bit = 0
			}
		}
		if optargs.Lazycounter_is_set {
			c_optargs.bitmask |= C.GUESTFS_XFS_ADMIN_LAZYCOUNTER_BITMASK
			if optargs.Lazycounter {
				c_optargs.lazycounter = 1
			} else {
				c_optargs.lazycounter = 0
			}
		}
		if optargs.Label_is_set {
			c_optargs.bitmask |= C.GUESTFS_XFS_ADMIN_LABEL_BITMASK
			c_optargs.label = C.CString(optargs.Label)
			defer C.free(unsafe.Pointer(c_optargs.label))
		}
		if optargs.Uuid_is_set {
			c_optargs.bitmask |= C.GUESTFS_XFS_ADMIN_UUID_BITMASK
			c_optargs.uuid = C.CString(optargs.Uuid)
			defer C.free(unsafe.Pointer(c_optargs.uuid))
		}
	}

	r := C.guestfs_xfs_admin_argv(g.g, c_device, &c_optargs)

	if r == -1 {
		return get_error_from_handle(g, "xfs_admin")
	}
	return nil
}

/* Struct carrying optional arguments for Xfs_growfs */
type OptargsXfs_growfs struct {
	/* Datasec field is ignored unless Datasec_is_set == true */
	Datasec_is_set bool
	Datasec        bool
	/* Logsec field is ignored unless Logsec_is_set == true */
	Logsec_is_set bool
	Logsec        bool
	/* Rtsec field is ignored unless Rtsec_is_set == true */
	Rtsec_is_set bool
	Rtsec        bool
	/* Datasize field is ignored unless Datasize_is_set == true */
	Datasize_is_set bool
	Datasize        int64
	/* Logsize field is ignored unless Logsize_is_set == true */
	Logsize_is_set bool
	Logsize        int64
	/* Rtsize field is ignored unless Rtsize_is_set == true */
	Rtsize_is_set bool
	Rtsize        int64
	/* Rtextsize field is ignored unless Rtextsize_is_set == true */
	Rtextsize_is_set bool
	Rtextsize        int64
	/* Maxpct field is ignored unless Maxpct_is_set == true */
	Maxpct_is_set bool
	Maxpct        int
}

/* xfs_growfs : expand an existing XFS filesystem */
func (g *Guestfs) Xfs_growfs(path string, optargs *OptargsXfs_growfs) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("xfs_growfs")
	}

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))
	c_optargs := C.struct_guestfs_xfs_growfs_argv{}
	if optargs != nil {
		if optargs.Datasec_is_set {
			c_optargs.bitmask |= C.GUESTFS_XFS_GROWFS_DATASEC_BITMASK
			if optargs.Datasec {
				c_optargs.datasec = 1
			} else {
				c_optargs.datasec = 0
			}
		}
		if optargs.Logsec_is_set {
			c_optargs.bitmask |= C.GUESTFS_XFS_GROWFS_LOGSEC_BITMASK
			if optargs.Logsec {
				c_optargs.logsec = 1
			} else {
				c_optargs.logsec = 0
			}
		}
		if optargs.Rtsec_is_set {
			c_optargs.bitmask |= C.GUESTFS_XFS_GROWFS_RTSEC_BITMASK
			if optargs.Rtsec {
				c_optargs.rtsec = 1
			} else {
				c_optargs.rtsec = 0
			}
		}
		if optargs.Datasize_is_set {
			c_optargs.bitmask |= C.GUESTFS_XFS_GROWFS_DATASIZE_BITMASK
			c_optargs.datasize = C.int64_t(optargs.Datasize)
		}
		if optargs.Logsize_is_set {
			c_optargs.bitmask |= C.GUESTFS_XFS_GROWFS_LOGSIZE_BITMASK
			c_optargs.logsize = C.int64_t(optargs.Logsize)
		}
		if optargs.Rtsize_is_set {
			c_optargs.bitmask |= C.GUESTFS_XFS_GROWFS_RTSIZE_BITMASK
			c_optargs.rtsize = C.int64_t(optargs.Rtsize)
		}
		if optargs.Rtextsize_is_set {
			c_optargs.bitmask |= C.GUESTFS_XFS_GROWFS_RTEXTSIZE_BITMASK
			c_optargs.rtextsize = C.int64_t(optargs.Rtextsize)
		}
		if optargs.Maxpct_is_set {
			c_optargs.bitmask |= C.GUESTFS_XFS_GROWFS_MAXPCT_BITMASK
			c_optargs.maxpct = C.int(optargs.Maxpct)
		}
	}

	r := C.guestfs_xfs_growfs_argv(g.g, c_path, &c_optargs)

	if r == -1 {
		return get_error_from_handle(g, "xfs_growfs")
	}
	return nil
}

/* xfs_info : get geometry of XFS filesystem */
func (g *Guestfs) Xfs_info(pathordevice string) (*XFSInfo, *GuestfsError) {
	if g.g == nil {
		return &XFSInfo{}, closed_handle_error("xfs_info")
	}

	c_pathordevice := C.CString(pathordevice)
	defer C.free(unsafe.Pointer(c_pathordevice))

	r := C.guestfs_xfs_info(g.g, c_pathordevice)

	if r == nil {
		return &XFSInfo{}, get_error_from_handle(g, "xfs_info")
	}
	defer C.guestfs_free_xfsinfo(r)
	return return_XFSInfo(r), nil
}

/* Struct carrying optional arguments for Xfs_repair */
type OptargsXfs_repair struct {
	/* Forcelogzero field is ignored unless Forcelogzero_is_set == true */
	Forcelogzero_is_set bool
	Forcelogzero        bool
	/* Nomodify field is ignored unless Nomodify_is_set == true */
	Nomodify_is_set bool
	Nomodify        bool
	/* Noprefetch field is ignored unless Noprefetch_is_set == true */
	Noprefetch_is_set bool
	Noprefetch        bool
	/* Forcegeometry field is ignored unless Forcegeometry_is_set == true */
	Forcegeometry_is_set bool
	Forcegeometry        bool
	/* Maxmem field is ignored unless Maxmem_is_set == true */
	Maxmem_is_set bool
	Maxmem        int64
	/* Ihashsize field is ignored unless Ihashsize_is_set == true */
	Ihashsize_is_set bool
	Ihashsize        int64
	/* Bhashsize field is ignored unless Bhashsize_is_set == true */
	Bhashsize_is_set bool
	Bhashsize        int64
	/* Agstride field is ignored unless Agstride_is_set == true */
	Agstride_is_set bool
	Agstride        int64
	/* Logdev field is ignored unless Logdev_is_set == true */
	Logdev_is_set bool
	Logdev        string
	/* Rtdev field is ignored unless Rtdev_is_set == true */
	Rtdev_is_set bool
	Rtdev        string
}

/* xfs_repair : repair an XFS filesystem */
func (g *Guestfs) Xfs_repair(device string, optargs *OptargsXfs_repair) (int, *GuestfsError) {
	if g.g == nil {
		return 0, closed_handle_error("xfs_repair")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))
	c_optargs := C.struct_guestfs_xfs_repair_argv{}
	if optargs != nil {
		if optargs.Forcelogzero_is_set {
			c_optargs.bitmask |= C.GUESTFS_XFS_REPAIR_FORCELOGZERO_BITMASK
			if optargs.Forcelogzero {
				c_optargs.forcelogzero = 1
			} else {
				c_optargs.forcelogzero = 0
			}
		}
		if optargs.Nomodify_is_set {
			c_optargs.bitmask |= C.GUESTFS_XFS_REPAIR_NOMODIFY_BITMASK
			if optargs.Nomodify {
				c_optargs.nomodify = 1
			} else {
				c_optargs.nomodify = 0
			}
		}
		if optargs.Noprefetch_is_set {
			c_optargs.bitmask |= C.GUESTFS_XFS_REPAIR_NOPREFETCH_BITMASK
			if optargs.Noprefetch {
				c_optargs.noprefetch = 1
			} else {
				c_optargs.noprefetch = 0
			}
		}
		if optargs.Forcegeometry_is_set {
			c_optargs.bitmask |= C.GUESTFS_XFS_REPAIR_FORCEGEOMETRY_BITMASK
			if optargs.Forcegeometry {
				c_optargs.forcegeometry = 1
			} else {
				c_optargs.forcegeometry = 0
			}
		}
		if optargs.Maxmem_is_set {
			c_optargs.bitmask |= C.GUESTFS_XFS_REPAIR_MAXMEM_BITMASK
			c_optargs.maxmem = C.int64_t(optargs.Maxmem)
		}
		if optargs.Ihashsize_is_set {
			c_optargs.bitmask |= C.GUESTFS_XFS_REPAIR_IHASHSIZE_BITMASK
			c_optargs.ihashsize = C.int64_t(optargs.Ihashsize)
		}
		if optargs.Bhashsize_is_set {
			c_optargs.bitmask |= C.GUESTFS_XFS_REPAIR_BHASHSIZE_BITMASK
			c_optargs.bhashsize = C.int64_t(optargs.Bhashsize)
		}
		if optargs.Agstride_is_set {
			c_optargs.bitmask |= C.GUESTFS_XFS_REPAIR_AGSTRIDE_BITMASK
			c_optargs.agstride = C.int64_t(optargs.Agstride)
		}
		if optargs.Logdev_is_set {
			c_optargs.bitmask |= C.GUESTFS_XFS_REPAIR_LOGDEV_BITMASK
			c_optargs.logdev = C.CString(optargs.Logdev)
			defer C.free(unsafe.Pointer(c_optargs.logdev))
		}
		if optargs.Rtdev_is_set {
			c_optargs.bitmask |= C.GUESTFS_XFS_REPAIR_RTDEV_BITMASK
			c_optargs.rtdev = C.CString(optargs.Rtdev)
			defer C.free(unsafe.Pointer(c_optargs.rtdev))
		}
	}

	r := C.guestfs_xfs_repair_argv(g.g, c_device, &c_optargs)

	if r == -1 {
		return 0, get_error_from_handle(g, "xfs_repair")
	}
	return int(r), nil
}

/* zegrep : return lines matching a pattern */
func (g *Guestfs) Zegrep(regex string, path string) ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("zegrep")
	}

	c_regex := C.CString(regex)
	defer C.free(unsafe.Pointer(c_regex))

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_zegrep(g.g, c_regex, c_path)

	if r == nil {
		return nil, get_error_from_handle(g, "zegrep")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* zegrepi : return lines matching a pattern */
func (g *Guestfs) Zegrepi(regex string, path string) ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("zegrepi")
	}

	c_regex := C.CString(regex)
	defer C.free(unsafe.Pointer(c_regex))

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_zegrepi(g.g, c_regex, c_path)

	if r == nil {
		return nil, get_error_from_handle(g, "zegrepi")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* zero : write zeroes to the device */
func (g *Guestfs) Zero(device string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("zero")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_zero(g.g, c_device)

	if r == -1 {
		return get_error_from_handle(g, "zero")
	}
	return nil
}

/* zero_device : write zeroes to an entire device */
func (g *Guestfs) Zero_device(device string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("zero_device")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_zero_device(g.g, c_device)

	if r == -1 {
		return get_error_from_handle(g, "zero_device")
	}
	return nil
}

/* zero_free_space : zero free space in a filesystem */
func (g *Guestfs) Zero_free_space(directory string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("zero_free_space")
	}

	c_directory := C.CString(directory)
	defer C.free(unsafe.Pointer(c_directory))

	r := C.guestfs_zero_free_space(g.g, c_directory)

	if r == -1 {
		return get_error_from_handle(g, "zero_free_space")
	}
	return nil
}

/* zerofree : zero unused inodes and disk blocks on ext2/3 filesystem */
func (g *Guestfs) Zerofree(device string) *GuestfsError {
	if g.g == nil {
		return closed_handle_error("zerofree")
	}

	c_device := C.CString(device)
	defer C.free(unsafe.Pointer(c_device))

	r := C.guestfs_zerofree(g.g, c_device)

	if r == -1 {
		return get_error_from_handle(g, "zerofree")
	}
	return nil
}

/* zfgrep : return lines matching a pattern */
func (g *Guestfs) Zfgrep(pattern string, path string) ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("zfgrep")
	}

	c_pattern := C.CString(pattern)
	defer C.free(unsafe.Pointer(c_pattern))

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_zfgrep(g.g, c_pattern, c_path)

	if r == nil {
		return nil, get_error_from_handle(g, "zfgrep")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* zfgrepi : return lines matching a pattern */
func (g *Guestfs) Zfgrepi(pattern string, path string) ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("zfgrepi")
	}

	c_pattern := C.CString(pattern)
	defer C.free(unsafe.Pointer(c_pattern))

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_zfgrepi(g.g, c_pattern, c_path)

	if r == nil {
		return nil, get_error_from_handle(g, "zfgrepi")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* zfile : determine file type inside a compressed file */
func (g *Guestfs) Zfile(meth string, path string) (string, *GuestfsError) {
	if g.g == nil {
		return "", closed_handle_error("zfile")
	}

	c_meth := C.CString(meth)
	defer C.free(unsafe.Pointer(c_meth))

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_zfile(g.g, c_meth, c_path)

	if r == nil {
		return "", get_error_from_handle(g, "zfile")
	}
	defer C.free(unsafe.Pointer(r))
	return C.GoString(r), nil
}

/* zgrep : return lines matching a pattern */
func (g *Guestfs) Zgrep(regex string, path string) ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("zgrep")
	}

	c_regex := C.CString(regex)
	defer C.free(unsafe.Pointer(c_regex))

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_zgrep(g.g, c_regex, c_path)

	if r == nil {
		return nil, get_error_from_handle(g, "zgrep")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}

/* zgrepi : return lines matching a pattern */
func (g *Guestfs) Zgrepi(regex string, path string) ([]string, *GuestfsError) {
	if g.g == nil {
		return nil, closed_handle_error("zgrepi")
	}

	c_regex := C.CString(regex)
	defer C.free(unsafe.Pointer(c_regex))

	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	r := C.guestfs_zgrepi(g.g, c_regex, c_path)

	if r == nil {
		return nil, get_error_from_handle(g, "zgrepi")
	}
	defer free_string_list(r)
	return return_string_list(r), nil
}
