# $NetBSD$

DISTNAME=	libcxxabi
PKGNAME=	${DISTNAME}6-6.0.0alpha1
CATEGORIES=	lang devel

MASTER_SITES=	${MASTER_SITE_GITHUB:=llvm-mirror/}
GITHUB_PROJECT=	${DISTNAME}
GITHUB_TAG=	9245d481eb3e890f708ff2d7dadf2a10c04748ba

MAINTAINER=	pkgsrc-users@NetBSD.org
HOMEPAGE=	http://libcxxabi.llvm.org/
COMMENT=	C++ Standard Library
LICENSE=	modified-bsd OR mit

# libcxx6 requires libcxxabi6's cxxabi.h to build, but libcxxabi6 requires
# libcxx6's includes to build, so we have this unholy extra distfiles thing
# in both of them to get them to build against each other without causing
# circular dependencies.  This must be kept in sync with libcxxabi6/Makefile.
#
LIBCXX_SHA=		390d70d797bde90fe31216b2833edec30ca483fe
LIBCXX=			libcxx-${LIBCXX_SHA}.tar.gz
SITES.${LIBCXX}=	-${MASTER_SITES}/libcxx/archive/${LIBCXX_SHA}.tar.gz

DISTFILES=	${DEFAULT_DISTFILES} ${LIBCXX}

CONFIGURE_DIRS=		${WRKDIR}/build
CMAKE_ARG_PATH=		${WRKSRC}

USE_LANGUAGES=		c c++
USE_CMAKE=		yes
GCC_REQD+=		4.8

CMAKE_INSTALL_PREFIX=	${PREFIX}/llvm6

CMAKE_ARGS+=	-DCMAKE_BUILD_TYPE=Release
CMAKE_ARGS+=	-DCMAKE_C_COMPILER=${CC:Q}
CMAKE_ARGS+=	-DCMAKE_CXX_COMPILER=${CXX:Q}
CMAKE_ARGS+=	-DLIBCXXABI_LIBCXX_INCLUDES=${WRKDIR}/libcxx-${LIBCXX_SHA}/include

post-extract:
	${RUN} mkdir -p ${WRKDIR}/build

.include "../../lang/llvm6/buildlink3.mk"
.include "../../mk/bsd.pkg.mk"
