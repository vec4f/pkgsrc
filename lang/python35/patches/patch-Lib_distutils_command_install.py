$NetBSD: patch-Lib_distutils_command_install.py,v 1.1 2015/12/05 17:12:13 adam Exp $

Support multiarch.

--- Lib/distutils/command/install.py.orig	2012-02-23 20:22:44.000000000 +0000
+++ Lib/distutils/command/install.py
@@ -29,8 +29,8 @@ WINDOWS_SCHEME = {
 
 INSTALL_SCHEMES = {
     'unix_prefix': {
-        'purelib': '$base/lib/python$py_version_short/site-packages',
-        'platlib': '$platbase/lib/python$py_version_short/site-packages',
+        'purelib': '$base/lib/python$py_version_short/site-packages$libarchsuffix',
+        'platlib': '$platbase/lib/python$py_version_short/site-packages$libarchsuffix',
         'headers': '$base/include/python$py_version_short$abiflags/$dist_name',
         'scripts': '$base/bin',
         'data'   : '$base',
@@ -281,6 +281,10 @@ class install(Command):
 
         py_version = sys.version.split()[0]
         (prefix, exec_prefix) = get_config_vars('prefix', 'exec_prefix')
+        if sys.maxsize > 2**32:
+            self.libarchsuffix = "@LIBARCHSUFFIX.64@"
+        else:
+            self.libarchsuffix = "@LIBARCHSUFFIX.32@"
         try:
             abiflags = sys.abiflags
         except AttributeError:
@@ -297,6 +301,7 @@ class install(Command):
                             'sys_exec_prefix': exec_prefix,
                             'exec_prefix': exec_prefix,
                             'abiflags': abiflags,
+                            'libarchsuffix': self.libarchsuffix,
                            }
 
         if HAS_USER_SITE:
@@ -646,5 +651,6 @@ class install(Command):
                     ('install_headers', has_headers),
                     ('install_scripts', has_scripts),
                     ('install_data',    has_data),
-                    ('install_egg_info', lambda self:True),
                    ]
+    if not os.environ.get('PKGSRC_PYTHON_NO_EGG'):
+        sub_commands += [('install_egg_info', lambda self:True),]
