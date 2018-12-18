$NetBSD: patch-largefile.h,v 1.2 2018/12/18 14:16:18 hauke Exp $

Fix for Radmind bug #221, accomodating for 64 bit time_t

--- largefile.h.orig	2010-12-13 03:42:49.000000000 +0000
+++ largefile.h
@@ -18,3 +18,19 @@
 #define strtoofft(x,y,z)	(strtol((x),(y),(z)))
 #define PRIofft			"l"
 #endif
+
+#ifndef SIZEOF_TIME_T
+#error "sizeof time_t unknown."
+#endif
+
+#if SIZEOF_TIME_T == 8
+    #ifdef HAVE_STRTOLL
+    #define strtotimet(x,y,z)	(strtoll((x),(y),(z)))
+    #else /* !HAVE_STRTOLL */
+    #define strtotimet(x,y,z)	(strtol((x),(y),(z)))
+    #endif /* HAVE_STRTOLL */
+    #define PRItimet		"ll"
+#else /* SIZEOF_TIME_T != 8 */
+    #define strtotimet(x,y,z)	(strtol((x),(y),(z)))
+    #define PRItimet		"l"
+#endif /* SIZEOF_TIME_T */
