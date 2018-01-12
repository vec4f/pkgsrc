# $NetBSD$

DISTNAME=	libcxx
PKGNAME=	${DISTNAME}6-6.0.0alpha1
CATEGORIES=	lang devel

MASTER_SITES=	${MASTER_SITE_GITHUB:=llvm-mirror/}
GITHUB_PROJECT=	${DISTNAME}
GITHUB_TAG=	390d70d797bde90fe31216b2833edec30ca483fe

MAINTAINER=	pkgsrc-users@NetBSD.org
HOMEPAGE=	http://libcxx.llvm.org/
COMMENT=	C++ Standard Library
LICENSE=	modified-bsd OR mit

# libcxx6 requires libcxxabi6's cxxabi.h to build, but libcxxabi6 requires
# libcxx6's includes to build, so we have this unholy extra distfiles thing
# in both of them to get them to build against each other without causing
# circular dependencies.  This must be kept in sync with libcxxabi6/Makefile.
#
LIBCXXABI_SHA=		9245d481eb3e890f708ff2d7dadf2a10c04748ba
LIBCXXABI=		libcxxabi-${LIBCXXABI_SHA}.tar.gz
SITES.${LIBCXXABI}=	-${MASTER_SITES}/libcxxabi/archive/${LIBCXXABI_SHA}.tar.gz

DISTFILES=	${DEFAULT_DISTFILES} ${LIBCXXABI}

CONFIGURE_DIRS=		${WRKDIR}/build
CMAKE_ARG_PATH=		${WRKSRC}

USE_LANGUAGES=		c c++
USE_CMAKE=		yes
GCC_REQD+=		4.8

CMAKE_INSTALL_PREFIX=	${PREFIX}/llvm6
BUILDLINK_PASSTHRU_RPATHDIRS+=	${CMAKE_INSTALL_PREFIX}/lib
LDFLAGS+=		-L${CMAKE_INSTALL_PREFIX}/lib
LDFLAGS+=		${COMPILER_RPATH_FLAG}${CMAKE_INSTALL_PREFIX}/lib

CMAKE_ARGS+=	-DCMAKE_BUILD_TYPE=Release
CMAKE_ARGS+=	-DCMAKE_C_COMPILER=${CC:Q}
CMAKE_ARGS+=	-DCMAKE_CXX_COMPILER=${CXX:Q}
CMAKE_ARGS+=	-DCMAKE_CXX_LINK_FLAGS=${LDFLAGS:Q}
CMAKE_ARGS+=	-DLIBCXX_CXX_ABI=libcxxabi
CMAKE_ARGS+=	-DLIBCXX_CXX_ABI_INCLUDE_PATHS=${WRKDIR}/libcxxabi-${LIBCXXABI_SHA}/include
CMAKE_ARGS+=	-DLIBCXX_ENABLE_ABI_LINKER_SCRIPT:BOOL=OFF

post-extract:
	${RUN} mkdir -p ${WRKDIR}/build

.include "../../lang/libcxxabi6/buildlink3.mk"
.include "../../lang/llvm6/buildlink3.mk"
.include "../../mk/bsd.pkg.mk"
