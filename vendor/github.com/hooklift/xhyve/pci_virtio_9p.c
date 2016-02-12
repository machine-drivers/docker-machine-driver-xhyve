/*-
 * Copyright (c) 2016 Jakub Klama <jceel@FreeBSD.org>
 * Copyright (c) 2016 iXsystems Inc.
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions
 * are met:
 * 1. Redistributions of source code must retain the above copyright
 *    notice, this list of conditions and the following disclaimer
 *    in this position and unchanged.
 * 2. Redistributions in binary form must reproduce the above copyright
 *    notice, this list of conditions and the following disclaimer in the
 *    documentation and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE AUTHOR AND CONTRIBUTORS ``AS IS'' AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED.  IN NO EVENT SHALL THE AUTHOR OR CONTRIBUTORS BE LIABLE
 * FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
 * OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
 * HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
 * LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY
 * OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF
 * SUCH DAMAGE.
 */

/*
 * virtio filesystem passtrough using 9p protocol.
 */

#include <sys/cdefs.h>

#include <sys/param.h>
#include <sys/uio.h>

#include <errno.h>
#include <fcntl.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <assert.h>
#include <pthread.h>

#pragma clang diagnostic push
#pragma clang diagnostic ignored "-Wpadded"
#include <lib9p.h>
#pragma clang diagnostic pop

#include <xhyve/xhyve.h>
#include <xhyve/pci_emul.h>
#include <xhyve/virtio.h>


#define VT9P_RINGSZ	64

static int pci_vt9p_debug = 0;
#define DPRINTF(params) if (pci_vt9p_debug) printf params

/*
 * Per-device softc
 */
struct pci_vt9p_softc {
	struct virtio_softc      vsc_vs;
	struct vqueue_info       vsc_vq;
	pthread_mutex_t          vsc_mtx;
	uint64_t                 vsc_cfg;
	uint64_t                 vsc_features;
	char *                   vsc_rootpath;
	struct pci_vt9p_config * vsc_config;
	struct l9p_backend *     vsc_fs_backend;
	struct l9p_server *      vsc_server;
        struct l9p_connection *  vsc_conn;
};

struct pci_vt9p_request {
	struct iovec *		vsr_iov;
	size_t			vsr_niov;
	size_t			vsr_respidx;
	size_t			vsr_iolen;
};

#pragma clang diagnostic push
#pragma clang diagnostic ignored "-Wzero-length-array"
struct pci_vt9p_config {
	uint16_t tag_len;
	char tag[0];
};
#pragma clang diagnostic pop

static void pci_vt9p_reset(void *);
static void pci_vt9p_notify(void *, struct vqueue_info *);
static int pci_vt9p_cfgread(void *, int, int, uint32_t *);
static void pci_vt9p_neg_features(void *, uint64_t);

static struct virtio_consts vt9p_vi_consts = {
	"vt9p",			/* our name */
	1,			/* we support 1 virtqueue */
	0,			/* config reg size */
	pci_vt9p_reset,		/* reset */
	pci_vt9p_notify,	/* device-wide qnotify */
	pci_vt9p_cfgread,	/* read virtio config */
	NULL,			/* write virtio config */
	pci_vt9p_neg_features,	/* apply negotiated features */
	(1 << 0),		/* our capabilities */
};


static void
pci_vt9p_reset(void *vsc)
{
	struct pci_vt9p_softc *sc;

	sc = vsc;

	DPRINTF(("vt9p: device reset requested !\n"));
	vi_reset_dev(&sc->vsc_vs);
}

static void
pci_vt9p_neg_features(void *vsc, uint64_t negotiated_features)
{
	struct pci_vt9p_softc *sc = vsc;

	sc->vsc_features = negotiated_features;
}

static int
pci_vt9p_cfgread(void *vsc, int offset, int size, uint32_t *retval)
{
	struct pci_vt9p_softc *sc = vsc;
	void *ptr;

	ptr = (uint8_t *)sc->vsc_config + offset;
	memcpy(retval, ptr, size);
	return (0);
}

static int
pci_vt9p_get_buffer(struct l9p_request *req, struct iovec *iov, size_t *niov,
    void *arg __unused)
{
	struct pci_vt9p_request *preq = req->lr_aux;
	size_t n = preq->vsr_niov - preq->vsr_respidx;
	
	memcpy(iov, preq->vsr_iov + preq->vsr_respidx, n * sizeof(struct iovec));
	*niov = n;
	return (0);
}

static int
pci_vt9p_send(struct l9p_request *req, const struct iovec *iov __unused,
    const size_t niov __unused, const size_t iolen, void *arg __unused)
{
	struct pci_vt9p_request *preq = req->lr_aux;

	preq->vsr_iolen = iolen;
	return (0);
}

static void
pci_vt9p_notify(void *vsc, struct vqueue_info *vq)
{
	struct iovec iov[8];
	struct pci_vt9p_softc *sc;
	struct pci_vt9p_request preq;
	uint16_t idx, i;
	uint16_t flags[8];
	int n;

	sc = vsc;

	while (vq_has_descs(vq)) {
		n = vq_getchain(vq, &idx, iov, 8, flags);
		preq.vsr_iov = iov;
		preq.vsr_niov = (size_t)n;
		preq.vsr_respidx = 0;

		/* Count readable descriptors */
		for (i = 0; i < n; i++) {
			if (flags[i] & VRING_DESC_F_WRITE)
				break;

			preq.vsr_respidx++;
		}

		for (i = 0; i < n; i++)
			DPRINTF(("vt9p: vt9p_notify(): desc%d base=%p, len=%zu, flags=0x%04x\r\n", i, iov[i].iov_base, iov[i].iov_len, flags[i]));

		l9p_connection_recv(sc->vsc_conn, iov, preq.vsr_respidx, &preq);

		/*
		 * Release this chain and handle more
		 */
		vq_relchain(vq, idx, (uint32_t)preq.vsr_iolen);
	}
	vq_endchains(vq, 1);	/* Generate interrupt if appropriate. */
}


static int
pci_vt9p_init(struct pci_devinst *pi, char *opts)
{
	struct pci_vt9p_softc *sc;
	char *opt;
	char *sharename = NULL;
	char *rootpath = NULL;

	if (opts == NULL) {
		printf("virtio-9p: share name and path required\n");
		return (1);
	}

	sc = calloc(1, sizeof(struct pci_vt9p_softc));
	sc->vsc_config = calloc(1, sizeof(struct pci_vt9p_config) + 128);

	while ((opt = strsep(&opts, ",")) != NULL) {
		if (sharename == NULL) {
			sharename = strsep(&opt, "=");
			rootpath = strdup(opt);
			continue;
		}

		if (strcmp(opt, "ro") == 0)
			DPRINTF(("read-only mount requested\r\n"));
	}

	sc->vsc_config->tag_len = (uint16_t)strlen(sharename);
	strncpy(sc->vsc_config->tag, sharename, strlen(sharename));
	
	if (l9p_backend_fs_init(&sc->vsc_fs_backend, rootpath) != 0) {
		errno = ENXIO;
		return (1);
	}

	if (l9p_server_init(&sc->vsc_server, sc->vsc_fs_backend) != 0) {
		errno = ENXIO;
		return (1);
	}

	if (l9p_connection_init(sc->vsc_server, &sc->vsc_conn) != 0) {
		errno = EIO;
		return (1);
	}

	l9p_connection_on_send_response(sc->vsc_conn, pci_vt9p_send, NULL);
	l9p_connection_on_get_response_buffer(sc->vsc_conn, pci_vt9p_get_buffer, NULL);

	vi_softc_linkup(&sc->vsc_vs, &vt9p_vi_consts, sc, pi, &sc->vsc_vq);
	sc->vsc_vs.vs_mtx = &sc->vsc_mtx;

	sc->vsc_vq.vq_qsize = VT9P_RINGSZ;

	/* initialize config space */
	pci_set_cfgdata16(pi, PCIR_DEVICE, VIRTIO_DEV_9P);
	pci_set_cfgdata16(pi, PCIR_VENDOR, VIRTIO_VENDOR);
	pci_set_cfgdata8(pi, PCIR_CLASS, PCIC_STORAGE);
	pci_set_cfgdata16(pi, PCIR_SUBDEV_0, VIRTIO_TYPE_9P);
	pci_set_cfgdata16(pi, PCIR_SUBVEND_0, VIRTIO_VENDOR);

	if (vi_intr_init(&sc->vsc_vs, 1, fbsdrun_virtio_msix()))
		return (1);
	vi_set_io_bar(&sc->vsc_vs, 0);

	return (0);
}

static struct pci_devemu pci_dev_9p = {
	.pe_emu =	"virtio-9p",
	.pe_init =	pci_vt9p_init,
	.pe_barwrite =	vi_pci_write,
	.pe_barread =	vi_pci_read
};
PCI_EMUL_SET(pci_dev_9p);
