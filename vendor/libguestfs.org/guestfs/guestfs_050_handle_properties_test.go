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

func Test050HandleProperties(t *testing.T) {
	g, errno := Create()
	if errno != nil {
		t.Errorf("could not create handle: %s", errno)
	}
	defer g.Close()

	v, err := g.Get_verbose()
	if err != nil {
		t.Errorf("%s", err)
	}
	err = g.Set_verbose(v)
	if err != nil {
		t.Errorf("%s", err)
	}

	tr, err := g.Get_trace()
	if err != nil {
		t.Errorf("%s", err)
	}
	err = g.Set_trace(tr)
	if err != nil {
		t.Errorf("%s", err)
	}

	m, err := g.Get_memsize()
	if err != nil {
		t.Errorf("%s", err)
	}
	err = g.Set_memsize(m)
	if err != nil {
		t.Errorf("%s", err)
	}

	p, err := g.Get_path()
	if err != nil {
		t.Errorf("%s", err)
	}
	err = g.Set_path(&p)
	if err != nil {
		t.Errorf("%s", err)
	}
}
