# $NetBSD: Makefile,v 1.32 2018/06/15 11:49:17 tsutsui Exp $
#

DISTNAME=	smpeg-0.4.4
PKGREVISION=	17
CATEGORIES=	devel multimedia
MASTER_SITES=	#

MAINTAINER=	pkgsrc-users@NetBSD.org
HOMEPAGE=	http://www.lokigames.com/development/smpeg.php3
COMMENT=	SDL MPEG Player Library
LICENSE=	gnu-lgpl-v2

USE_TOOLS+=	gmake
GNU_CONFIGURE=	yes
USE_LANGUAGES=	c c++
USE_LIBTOOL=	yes

CPPFLAGS+=	-Dunix=1

# smpeg sources (glmovie.c & glmovie-tile.c) assume <malloc.h> exists.
# Create a fake one for platforms that don't have it.
#
post-wrapper:
	if ${TEST} ! -f /usr/include/malloc.h; then			\
		${ECHO} "#include <stdlib.h>"				\
			> ${BUILDLINK_DIR}/include/malloc.h;		\
	fi

.include "../../mk/bsd.prefs.mk"

.if ${MACHINE_ARCH} == "i386"
. if ${OPSYS} != "Darwin" && \
     empty(MACHINE_PLATFORM:MNetBSD-[89].*-i386)
CONFIGURE_ARGS+=	--enable-mmx
. endif
.endif

MAKE_ENV.SunOS+=	SunOS_LDFLAGS="-shared -mimpure-text"

.include "options.mk"
.include "../../devel/SDL/buildlink3.mk"
.include "../../mk/bsd.pkg.mk"
