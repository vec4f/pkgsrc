# $NetBSD: Makefile,v 1.3 2018/11/15 09:22:51 adam Exp $

DISTNAME=	django-filter-2.0.0
PKGNAME=	${PYPKGPREFIX}-${DISTNAME}
CATEGORIES=	www python
MASTER_SITES=	${MASTER_SITE_PYPI:=d/django-filter/}

MAINTAINER=	joerg@NetBSD.org
HOMEPAGE=	https://github.com/carltongibson/django-filter
COMMENT=	Generic dynamic filtering support for Django
LICENSE=	modified-bsd

DEPENDS+=	${PYPKGPREFIX}-django>=1.11:../../www/py-djangorestframework
DEPENDS+=	${PYPKGPREFIX}-djangorestframework-[0-9]*:../../www/py-djangorestframework

USE_LANGUAGES=	# none

.include "../../lang/python/egg.mk"
.include "../../mk/bsd.pkg.mk"
