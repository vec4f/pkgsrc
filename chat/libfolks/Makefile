# $NetBSD: Makefile,v 1.12 2017/09/19 08:32:08 wiz Exp $

DISTNAME=	folks-0.4.2
PKGNAME=	lib${DISTNAME}
PKGREVISION=	5
CATEGORIES=	chat
MASTER_SITES=	${MASTER_SITE_GNOME:=sources/folks/0.4/}
EXTRACT_SUFX=	.tar.bz2

MAINTAINER=	pkgsrc-users@NetBSD.org
HOMEPAGE=	https://wiki.gnome.org/action/show/Projects/Folks
COMMENT=	Library to manage contacts
LICENSE=	gnu-lgpl-v2

GNU_CONFIGURE=	yes
USE_LIBTOOL=	yes
USE_TOOLS+=	pkg-config intltool gmake msgfmt
USE_PKGLOCALEDIR= yes

PKGCONFIG_OVERRIDE+=	backends/telepathy/lib/folks-telepathy.pc.in
PKGCONFIG_OVERRIDE+=	folks/folks.pc.in

TEST_TARGET=	check

CPPFLAGS.SunOS+=	-std=c99 -D_XOPEN_SOURCE=600

.include "../../devel/glib2/buildlink3.mk"
.include "../../devel/gobject-introspection/buildlink3.mk"
.include "../../sysutils/dbus-glib/buildlink3.mk"
.include "../../devel/libgee0.6/buildlink3.mk"
.include "../../chat/telepathy-glib/buildlink3.mk"
.include "../../mk/bsd.pkg.mk"
