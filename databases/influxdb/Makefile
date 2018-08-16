# $NetBSD$

DISTNAME=	influxdb-1.6.1
MASTER_SITES=	${MASTER_SITE_GITHUB:=influxdata/}
CATEGORIES=	databases
GITHUB_PROJECT=	influxdb
GITHUB_TAG=	v${PKGVERSION_NOREV}

MAINTAINER=	pkgsrc-users@NetBSD.org
HOMEPAGE=	https://github.com/influxdata/influxdb
COMMENT=	Scalable datastore for metrics, events, and real-time analytics
LICENSE=	mit

BUILD_DEPENDS+=	asciidoc-[0-9]*:../../textproc/asciidoc
BUILD_DEPENDS+=	xmlto-[0-9]*:../../textproc/xmlto

USE_TOOLS+=	gmake

GO_DIST_BASE=	${DISTNAME}
GO_SRCPATH=	github.com/influxdata/influxdb

GO_DEPS+=	github.com/collectd/go-collectd:2ce14454:collectd.org
GO_DEPS+=	github.com/BurntSushi/toml:a368813c
GO_DEPS+=	github.com/RoaringBitmap/roaring:d6540aab
GO_DEPS+=	github.com/beorn7/perks:4c0e8459
GO_DEPS+=	github.com/bmizerany/pat:6226ea59
GO_DEPS+=	github.com/boltdb/bolt:2f1ce7a8
GO_DEPS+=	github.com/cespare/xxhash:5c37fe37
GO_DEPS+=	github.com/davecgh/go-spew:346938d6
GO_DEPS+=	github.com/dgrijalva/jwt-go:06ea1031
GO_DEPS+=	github.com/dgryski/go-bitstream:9f22ccc2
GO_DEPS+=	github.com/glycerine/go-unsnap-stream:62a9a9eb
GO_DEPS+=	github.com/gogo/protobuf:1adfc126
GO_DEPS+=	github.com/golang/protobuf:92554152
GO_DEPS+=	github.com/golang/snappy:d9eb7a3d
GO_DEPS+=	github.com/google/go-cmp:3af367b6
GO_DEPS+=	github.com/influxdata/influxql:c661ab7d
GO_DEPS+=	github.com/influxdata/usage-client:6d389537
GO_DEPS+=	github.com/influxdata/yamux:1f58ded5
GO_DEPS+=	github.com/influxdata/yarpc:f0da2db1
GO_DEPS+=	github.com/jsternberg/zap-logfmt:ac4bd917
GO_DEPS+=	github.com/jwilder/encoding:b4e1701a
GO_DEPS+=	github.com/klauspost/compress:6c8db69c
GO_DEPS+=	github.com/klauspost/cpuid:ae7887de
GO_DEPS+=	github.com/klauspost/crc32:cb6bfca9
GO_DEPS+=	github.com/klauspost/pgzip:0bf5dcad
GO_DEPS+=	github.com/mattn/go-isatty:6ca4dbf5
GO_DEPS+=	github.com/matttproud/golang_protobuf_extensions:3247c845
GO_DEPS+=	github.com/mschoch/smat:90eadee7
GO_DEPS+=	github.com/opentracing/opentracing-go:328fceb7
GO_DEPS+=	github.com/paulbellamy/ratecounter:524851a9
GO_DEPS+=	github.com/peterh/liner:6106ee4f
GO_DEPS+=	github.com/philhofer/fwd:bb6d471d
GO_DEPS+=	github.com/prometheus/client_golang:661e31bf
GO_DEPS+=	github.com/prometheus/client_model:99fa1f4b
GO_DEPS+=	github.com/prometheus/common:e4aa40a9
GO_DEPS+=	github.com/prometheus/procfs:54d17b57
GO_DEPS+=	github.com/retailnext/hllpp:101a6d2f
GO_DEPS+=	github.com/tinylib/msgp:b2b6a672
GO_DEPS+=	github.com/willf/bitset:d860f346
GO_DEPS+=	github.com/xlab/treeprint:f3a15cfd
GO_DEPS+=	github.com/uber-go/atomic:8474b86a:go.uber.org/atomic
GO_DEPS+=	github.com/uber-go/multierr:3c493748:go.uber.org/multierr
GO_DEPS+=	github.com/uber-go/zap:35aad584:go.uber.org/zap
GO_DEPS+=	github.com/golang/crypto:c3a3ad6d:golang.org/x/crypto
GO_DEPS+=	github.com/golang/net:92b859f3:golang.org/x/net
GO_DEPS+=	github.com/golang/sync:1d60e460:golang.org/x/sync
GO_DEPS+=	github.com/golang/sys:d8e400bc:golang.org/x/sys
GO_DEPS+=	github.com/golang/text:f21a4dfb:golang.org/x/text
GO_DEPS+=	github.com/golang/time:26559e0f:golang.org/x/time

INSTALLATION_DIRS+=	bin sbin ${PKGMANDIR}/man1

post-build:
	cd ${WRKSRC}/man && ${SETENV} ${MAKE_ENV} ${GMAKE}
	${SED} -e 's,/var/lib,${VARBASE}/db,g' ${WRKSRC}/etc/config.sample.toml \
		>${WRKDIR}/.config.sample.toml

do-install:
	${INSTALL_PROGRAM} ${WRKDIR}/bin/influxd ${DESTDIR}${PREFIX}/sbin/${prog}
.for prog in influx influx_inspect influx_tsm influx_stress
	${INSTALL_PROGRAM} ${WRKDIR}/bin/${prog} ${DESTDIR}${PREFIX}/bin/${prog}
.endfor
	${INSTALL_MAN} ${WRKSRC}/man/*.1 ${DESTDIR}${PREFIX}/${PKGMANDIR}/man1/

.include "../../lang/go/go-package.mk"
.include "../../mk/bsd.pkg.mk"
