# $NetBSD: Makefile,v 1.38 2018/03/24 07:51:15 adam Exp $

DISTNAME=	virtualenv-15.2.0
PKGNAME=	${PYPKGPREFIX}-${DISTNAME}
CATEGORIES=	devel python
MASTER_SITES=	${MASTER_SITE_PYPI:=v/virtualenv/}

MAINTAINER=	pkgsrc-users@NetBSD.org
HOMEPAGE=	https://pypi.python.org/pypi/virtualenv
COMMENT=	Virtual Python Environment builder
LICENSE=	mit

TEST_DEPENDS+=	${PYPKGPREFIX}-mock-[0-9]*:../../devel/py-mock
TEST_DEPENDS+=	${PYPKGPREFIX}-test-[0-9]*:../../devel/py-test

USE_LANGUAGES=	# none

REPLACE_PYTHON=	virtualenv.py

.include "../../mk/bsd.prefs.mk"

SUBST_CLASSES+=		arch
SUBST_STAGE.arch=	pre-configure
SUBST_FILES.arch=	virtualenv.py
.if !empty(MULTIARCH:M[Yy][Ee][Ss])
SUBST_SED.arch=		-e 's|@BINARCHSUFFIX32@|${BINARCHSUFFIX.32:S/\///}|'
SUBST_SED.arch+=	-e 's|@BINARCHSUFFIX64@|${BINARCHSUFFIX.64:S/\///}|'
.else
SUBST_SED.arch=		-e 's|@BINARCHSUFFIX32@||'
SUBST_SED.arch+=	-e 's|@BINARCHSUFFIX64@||'
.endif

post-install:
	cd ${DESTDIR}${PREFIX}/bin && ${MV} virtualenv virtualenv-${PYVERSSUFFIX} || ${TRUE}

.include "../../lang/python/application.mk"
.include "../../lang/python/egg.mk"
.include "../../mk/bsd.pkg.mk"
