$NetBSD$

Support installing older package versions if explicitly requested.

--- impact.c.orig	2016-02-08 16:27:12.000000000 +0000
+++ impact.c
@@ -241,8 +241,8 @@ deps_impact(Plisthead *impacthead, Pkgli
 				 * local pkgname didn't match deps, remote pkg
 				 * has a lesser version than local package.
 				 */
-				if (version_check(plist->full,
-					remotepkg) == 1) {
+				if (pdp->level > 0 &&
+				    version_check(plist->full, remotepkg) == 1) {
 					/*
 					 * proposing a downgrade is definitely
 					 * not useful, not sure what I want to
