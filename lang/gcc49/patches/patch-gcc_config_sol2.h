$NetBSD$

Don't apply solaris_cxx_decl_mangling_context, it breaks clang which
expects the correct std::tm et al.  Breaks libstdc++ ABI, who cares.

--- gcc/config/sol2.h.orig	2014-05-28 11:37:50.000000000 +0000
+++ gcc/config/sol2.h
@@ -242,7 +242,6 @@ along with GCC; see the file COPYING3.
 /* Allow macro expansion in #pragma pack.  */
 #define HANDLE_PRAGMA_PACK_WITH_EXPANSION
 
-#define TARGET_CXX_DECL_MANGLING_CONTEXT solaris_cxx_decl_mangling_context
 
 /* Solaris/x86 as and gas support unquoted section names.  */
 #define SECTION_NAME_FORMAT	"%s"
