# $NetBSD: Makefile,v 1.11 2018/12/19 15:47:11 bsiegert Exp $

DISTNAME=		pup-0.4.0
PKGREVISION=		7
CATEGORIES=		www
MASTER_SITES=		${MASTER_SITE_GITHUB:=ericchiang/}
GITHUB_TAG=		v${PKGVERSION_NOREV}

MAINTAINER=		fhajny@NetBSD.org
HOMEPAGE=		https://github.com/ericchiang/pup/
COMMENT=		Parsing HTML at the command line
LICENSE=		mit

GO_DIST_BASE=		${DISTNAME}
GO_SRCPATH=		github.com/ericchiang/pup

INSTALLATION_DIRS+=	bin share/doc/pup

do-install:
	${INSTALL_PROGRAM} ${WRKDIR}/bin/pup ${DESTDIR}${PREFIX}/bin
	${INSTALL_DATA} ${WRKSRC}/README.md ${DESTDIR}${PREFIX}/share/doc/pup

.include "../../lang/go/go-package.mk"
.include "../../mk/bsd.pkg.mk"
