# $NetBSD: Makefile,v 1.3 2018/12/14 22:51:32 schmonz Exp $

DISTNAME=		fehQlibs-10
PKGNAME=		fehqlibs-0.9.10
CATEGORIES=		net
MASTER_SITES=		https://www.fehcom.de/ipnet/fehQlibs/
EXTRACT_SUFX=		.tgz

MAINTAINER=		schmonz@NetBSD.org
HOMEPAGE=		https://www.fehcom.de/ipnet/qlibs.html
COMMENT=		State-of-the-art C routines for Internet services
LICENSE=		public-domain

MAKE_JOBS_SAFE=		no

BUILD_TARGET=		libs

SUBST_CLASSES+=		echo
SUBST_STAGE.echo=	pre-configure
SUBST_FILES.echo=	configure install
SUBST_SED.echo=		-e 's|echo -n|${ECHO_N}|g'

SUBST_CLASSES+=		etc
SUBST_STAGE.etc=	do-configure
SUBST_FILES.etc=	dnsstub/dns_rcrw.c
SUBST_SED.etc=		-e 's|/etc/dnsrewrite|${PKG_SYSCONFBASE}/dnsrewrite|g'

BUILD_DEFS+=		PKG_SYSCONFBASE

do-configure:
	cd ${WRKSRC};								\
	${ECHO} "LIBDIR=${DESTDIR}${PREFIX}/lib/qlibs" >> ${WRKSRC}/conf-build;	\
	${ECHO} "HDRDIR=${DESTDIR}${PREFIX}/include/qlibs" >> ${WRKSRC}/conf-build; \
	${MAKE} check

.include "../../mk/bsd.pkg.mk"
