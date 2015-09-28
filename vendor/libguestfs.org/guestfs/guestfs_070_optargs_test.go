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

import "testing"

func Test070Optargs(t *testing.T) {
	g, errno := Create()
	if errno != nil {
		t.Errorf("could not create handle: %s", errno)
	}
	defer g.Close()
	err := g.Add_drive("/dev/null", nil)
	if err != nil {
		t.Errorf("%s", err)
	}
	err = g.Add_drive("/dev/null", &OptargsAdd_drive{
		Readonly_is_set: true,
		Readonly:        true,
	})
	if err != nil {
		t.Errorf("%s", err)
	}
	err = g.Add_drive("/dev/null", &OptargsAdd_drive{
		Readonly_is_set: true,
		Readonly:        true,
		Format_is_set:   true,
		Format:          "raw",
	})
	if err != nil {
		t.Errorf("%s", err)
	}
	err = g.Add_drive("/dev/null", &OptargsAdd_drive{
		Readonly_is_set: true,
		Readonly:        true,
		Format_is_set:   true,
		Format:          "raw",
		Iface_is_set:    true,
		Iface:           "virtio",
	})
	if err != nil {
		t.Errorf("%s", err)
	}
}
