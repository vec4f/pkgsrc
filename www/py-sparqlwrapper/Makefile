# $NetBSD: Makefile,v 1.4 2017/10/07 11:14:42 adam Exp $

DISTNAME=	SPARQLWrapper-1.8.0
PKGNAME=	${PYPKGPREFIX}-${DISTNAME:tl}
CATEGORIES=	www python
MASTER_SITES=	${MASTER_SITE_PYPI:=S/SPARQLWrapper/}

MAINTAINER=	pkgsrc-users@NetBSD.org
HOMEPAGE=	https://rdflib.github.io/sparqlwrapper/
COMMENT=	SPARQL endpoint interface to Python
LICENSE=	w3c

DEPENDS+=	${PYPKGPREFIX}-rdflib>=4.0:../../textproc/py-rdflib

USE_LANGUAGES=	# none

.include "../../lang/python/egg.mk"
.include "../../mk/bsd.pkg.mk"
