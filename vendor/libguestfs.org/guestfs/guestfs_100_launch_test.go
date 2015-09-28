// +build appliance

/* libguestfs Go tests
 * Copyright (C) 2013 Red Hat Inc.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.
 */

package guestfs

import (
	"testing"
	//	"sort"
)

func Test100Launch(t *testing.T) {
	g, errno := Create()
	if errno != nil {
		t.Errorf("could not create handle: %s", errno)
	}
	defer g.Close()

	err := g.Add_drive_scratch(500*1024*1024, nil)
	if err != nil {
		t.Errorf("%s", err)
	}

	err = g.Launch()
	if err != nil {
		t.Errorf("%s", err)
	}

	err = g.Pvcreate("/dev/sda")
	if err != nil {
		t.Errorf("%s", err)
	}

	err = g.Vgcreate("VG", []string{"/dev/sda"})
	if err != nil {
		t.Errorf("%s", err)
	}

	err = g.Lvcreate("LV1", "VG", 200)
	if err != nil {
		t.Errorf("%s", err)
	}
	err = g.Lvcreate("LV2", "VG", 200)
	if err != nil {
		t.Errorf("%s", err)
	}

	lvs, err := g.Lvs()
	if err != nil {
		t.Errorf("%s", err)
	}
	expected := []string{"/dev/VG/LV1", "/dev/VG/LV2"}
	if !equal(lvs, expected) {
		t.Errorf("g.Lvs: %s != %s", lvs, expected)
	}

	err = g.Mkfs("ext2", "/dev/VG/LV1", nil)
	if err != nil {
		t.Errorf("%s", err)
	}
	err = g.Mount("/dev/VG/LV1", "/")
	if err != nil {
		t.Errorf("%s", err)
	}
	err = g.Mkdir("/p")
	if err != nil {
		t.Errorf("%s", err)
	}
	err = g.Touch("/q")
	if err != nil {
		t.Errorf("%s", err)
	}

	_ /*dirs*/, err = g.Readdir("/")
	if err != nil {
		t.Errorf("%s", err)
	}
	//sort.Sort (byName (dirs))
	// XXX Sort interface is needlessly complicated

	err = g.Shutdown()
	if err != nil {
		t.Errorf("%s", err)
	}
}

/* - declared in guestfs_900_rstringlist_test.go
func equal (xs []string, ys []string) bool {
	if len(xs) != len(ys) {
		return false
	}
	for i, x := range xs {
		if x != ys[i] {
			return false
		}
	}
	return true
}
*/
