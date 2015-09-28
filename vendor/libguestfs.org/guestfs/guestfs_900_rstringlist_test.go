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
)

func Test900RStringLits(t *testing.T) {
	g, errno := Create()
	if errno != nil {
		t.Errorf("could not create handle: %s", errno)
	}
	defer g.Close()

	actual, err := g.Internal_test_rstringlist("16")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	expected := []string{
		"0", "1", "2", "3", "4", "5", "6", "7",
		"8", "9", "10", "11", "12", "13", "14", "15",
	}
	if !equal(actual, expected) {
		t.Errorf("%s != %s", actual, expected)
	}
}

func equal(xs []string, ys []string) bool {
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
