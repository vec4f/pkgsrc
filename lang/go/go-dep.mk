# $NetBSD$
#
# This file implements dep (https://github.com/golang/dep) support in pkgsrc.
#
# === Package-settable variables ===
#
# GO_DEPS
#	This is a list of dependencies of the format url:sha[:dir].  These
#	dependencies can be generated using the print-go-deps helper target
#	against an extracted WRKSRC, which parses the Godeps file contained in
#	the source that would normally be used by dep.
#
#	Currently only GitHub URLs are supported, so if a Godeps dependency
#	uses a different URL, the equivalent GitHub URL needs to be calculated
#	and the original URL added as the optional :dir argument.  This file
#	contains a list of GODEP_REDIRECTS which translates common redirects,
#	and additions to this list is encouraged.
#
#	The print-go-deps target will reduce SHA strings to the first 8
#	characters to avoid long lines, but the extracted archive will use the
#	full SHA.
#
#	Examples:
#		github.com/golang/protobuf:92554152
#		github.com/collectd/go-collectd:2ce14454:collectd.org
#		github.com/uber-go/atomic:8474b86a:go.uber.org/atomic
#		github.com/golang/time:26559e0f:golang.org/x/time
#

#
# When using GO_DEPS a lot of additional distfiles will be added with just SHAs
# for their filename.  Setting DIST_SUBDIR is required to keep things sensible.
#
DIST_SUBDIR=			${PKGNAME_NOREV}

#
# For now fetching via GitHub is mandatory.
#
DISTFILES?=			${DEFAULT_DISTFILES}
.for dep in ${GO_DEPS}
GODEP_URL.${dep}:=		${dep:C/:/ /g:[1]}
GODEP_SHA.${dep}:=		${dep:C/:/ /g:[2]}
GODEP_DIR.${dep}:=		${dep:C/:/ /g:[3]}
.  if empty(GODEP_DIR.${dep})
GODEP_DIR.${dep}:=		${GODEP_URL.${dep}}
.  endif
GODEP_TGZ.${dep}:=		${GODEP_SHA.${dep}}.tar.gz
SITES.${GODEP_TGZ.${dep}}:=	${MASTER_SITE_GITHUB}${GODEP_URL.${dep}:S,github.com/,,}/archive/
DISTFILES+=			${GODEP_TGZ.${dep}}

post-extract: post-extract-${GODEP_URL.${dep}}
post-extract-${GODEP_URL.${dep}}:
	# We can't use :H here as the GODEP_DIR might just be a host.name
	@${MKDIR} `${DIRNAME} ${WRKDIR}/src/${GODEP_DIR.${dep}}`
	# The * glob is used as we shorten the SHA to 8 characters to keep
	# GO_DEPS lines a reasonable length, but the extracted distfile uses
	# the full SHA.
	@${MV} ${WRKDIR}/${GODEP_URL.${dep}:T}-${GODEP_SHA.${dep}}* \
		${WRKDIR}/src/${GODEP_DIR.${dep}}
.endfor

#
# Non-GitHub sites which are listed in Godeps files but redirect to GitHub.  We
# fetch them from GitHub but move the source to the expected directory name.
#
GODEP_REDIRECTS+=		collectd.org=github.com/collectd/go-collectd
GODEP_REDIRECTS+=		go.uber.org=github.com/uber-go
GODEP_REDIRECTS+=		golang.org/x=github.com/golang

.for url in ${GODEP_REDIRECTS}
GODEP_REDIRECT_FROM.${url}:=	${url:C/=/ /g:[1]:C,/,\/,g}
GODEP_REDIRECT_TO.${url}:=	${url:C/=/ /g:[2]}
.endfor

#
# Add a print-go-deps target to aid GO_DEPS generation.
#
PRINT_GODEPS_AWK+=	${GODEP_REDIRECTS:@url@				\
			/${GODEP_REDIRECT_FROM.${url}}/ {		\
				dir=$$1;				\
				gsub(/${GODEP_REDIRECT_FROM.${url}}/,	\
				    "${GODEP_REDIRECT_TO.${url}}");	\
			}@}
PRINT_GODEPS_AWK+=	{
PRINT_GODEPS_AWK+=		printf("GO_DEPS+=\t%s:%s%s\n", $$1,
PRINT_GODEPS_AWK+=		    substr($$2, 1, 8),
PRINT_GODEPS_AWK+=		    (dir) ? ":" dir : "");
PRINT_GODEPS_AWK+=		dir="";
PRINT_GODEPS_AWK+=	}

.PHONY: print-go-deps
print-go-deps:
	@${AWK} '${PRINT_GODEPS_AWK}' ${WRKSRC}/Godeps
