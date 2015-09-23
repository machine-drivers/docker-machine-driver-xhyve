/*-
 * Copyright (c) 2011 NetApp, Inc.
 * Copyright (c) 2015 xhyve developers
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions
 * are met:
 * 1. Redistributions of source code must retain the above copyright
 *    notice, this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright
 *    notice, this list of conditions and the following disclaimer in the
 *    documentation and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY NETAPP, INC ``AS IS'' AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED.  IN NO EVENT SHALL NETAPP, INC OR CONTRIBUTORS BE LIABLE
 * FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
 * OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
 * HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
 * LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY
 * OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF
 * SUCH DAMAGE.
 *
 * $FreeBSD$
 */

/*
 * https://github.com/mirage/ocaml-vmnet/blob/master/lib/vmnet_stubs.c
 *
 * Copyright (C) 2014 Anil Madhavapeddy <anil@recoil.org>
 *
 * Permission to use, copy, modify, and distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

#include <stdio.h>

#include <vmnet/vmnet.h>

#include "uuid.h"

static char*
vmnet_get_mac_address_from_uuid(char *guest_uuid_str) {
/*
 * from vmn_create() in https://github.com/mist64/xhyve/blob/master/src/pci_virtio_vmnet.c
 */
  xpc_object_t interface_desc;
  uuid_t uuid;
  __block interface_ref iface;
  __block vmnet_return_t iface_status;
  __block char* mac = malloc(18);
  dispatch_semaphore_t iface_created, iface_stopped;
  dispatch_queue_t if_create_q, if_stop_q;
  uint32_t uuid_status;

  interface_desc = xpc_dictionary_create(NULL, NULL, 0);
  xpc_dictionary_set_uint64(interface_desc, vmnet_operation_mode_key, VMNET_SHARED_MODE);

  uuid_from_string(guest_uuid_str, &uuid, &uuid_status);
  if (uuid_status != uuid_s_ok) {
    fprintf(stderr, "Invalid UUID\n");
    return NULL;
  }

  xpc_dictionary_set_uuid(interface_desc, vmnet_interface_id_key, uuid);
  iface = NULL;
  iface_status = 0;

  if_create_q = dispatch_queue_create("org.xhyve.vmnet.create", DISPATCH_QUEUE_SERIAL);

  iface_created = dispatch_semaphore_create(0);

  iface = vmnet_start_interface(interface_desc, if_create_q,
    ^(vmnet_return_t status, xpc_object_t interface_param)
  {
    iface_status = status;
    if (status != VMNET_SUCCESS || !interface_param) {
      dispatch_semaphore_signal(iface_created);
      return;
    }

    //printf("%s\n", xpc_dictionary_get_string(interface_param, vmnet_mac_address_key));
    const char *macStr = xpc_dictionary_get_string(interface_param, vmnet_mac_address_key);
    strcpy(mac, macStr);

    dispatch_semaphore_signal(iface_created);
  });

  dispatch_semaphore_wait(iface_created, DISPATCH_TIME_FOREVER);
  dispatch_release(if_create_q);

  if (iface == NULL || iface_status != VMNET_SUCCESS) {
    fprintf(stderr, "virtio_net: Could not create vmnet interface, "
      "permission denied or no entitlement?\n");
    return NULL;
  }

  iface_status = 0;

  if_stop_q = dispatch_queue_create("org.xhyve.vmnet.stop", DISPATCH_QUEUE_SERIAL);

  iface_stopped = dispatch_semaphore_create(0);

  iface_status = vmnet_stop_interface(iface, if_stop_q,
    ^(vmnet_return_t status)
  {
    iface_status = status;
    dispatch_semaphore_signal(iface_stopped);
  });

  dispatch_semaphore_wait(iface_stopped, DISPATCH_TIME_FOREVER);
  dispatch_release(if_stop_q);

  if (iface_status != VMNET_SUCCESS) {
    fprintf(stderr, "virtio_net: Could not stop vmnet interface, "
      "permission denied or no entitlement?\n");
    return NULL;
  }

  return mac;
}
