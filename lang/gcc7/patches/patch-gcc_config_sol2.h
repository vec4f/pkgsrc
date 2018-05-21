$NetBSD$

Enable full __cxa_atexit support.

--- gcc/config/sol2.h.orig	2017-11-21 09:31:12.135035000 +0000
+++ gcc/config/sol2.h
@@ -178,7 +178,7 @@ along with GCC; see the file COPYING3.
 				   shared|" PIE_SPEC ":crtbeginS.o%s; \
 				   :crtbegin.o%s}"
 #else
-#define STARTFILE_CRTBEGIN_SPEC	"crtbegin.o%s"
+#define STARTFILE_CRTBEGIN_SPEC	"%{shared:crtbeginS.o%s;:crtbegin.o%s}"
 #endif
 
 #if ENABLE_VTABLE_VERIFY
@@ -228,7 +228,7 @@ along with GCC; see the file COPYING3.
 			       shared|" PIE_SPEC ":crtendS.o%s; \
 			       :crtend.o%s}"
 #else
-#define ENDFILE_CRTEND_SPEC "crtend.o%s"
+#define ENDFILE_CRTEND_SPEC "%{shared:crtendS.o%s;:crtend.o%s}"
 #endif
 
 #undef  ENDFILE_SPEC
