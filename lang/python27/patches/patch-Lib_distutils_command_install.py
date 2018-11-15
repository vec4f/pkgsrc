$NetBSD: patch-Lib_distutils_command_install.py,v 1.1 2018/06/17 19:21:21 adam Exp $

Add a knob (enviroment variable) for disabling installation of egg metadata
in extensions until we have infrastructure in place for dealing w/ it.

Support multiarch.

--- Lib/distutils/command/install.py.orig	2014-12-10 15:59:34.000000000 +0000
+++ Lib/distutils/command/install.py
@@ -41,8 +41,8 @@
 
 INSTALL_SCHEMES = {
     'unix_prefix': {
-        'purelib': '$base/lib/python$py_version_short/site-packages',
-        'platlib': '$platbase/lib/python$py_version_short/site-packages',
+        'purelib': '$base/lib/python$py_version_short/site-packages$libarchsuffix',
+        'platlib': '$platbase/lib/python$py_version_short/site-packages$libarchsuffix',
         'headers': '$base/include/python$py_version_short/$dist_name',
         'scripts': '$base/bin',
         'data'   : '$base',
@@ -298,6 +298,10 @@
 
         py_version = (string.split(sys.version))[0]
         (prefix, exec_prefix) = get_config_vars('prefix', 'exec_prefix')
+        if sys.maxsize > 2**32:
+            self.libarchsuffix = "@LIBARCHSUFFIX.64@"
+        else:
+            self.libarchsuffix = "@LIBARCHSUFFIX.32@"
         self.config_vars = {'dist_name': self.distribution.get_name(),
                             'dist_version': self.distribution.get_version(),
                             'dist_fullname': self.distribution.get_fullname(),
@@ -310,6 +314,7 @@
                             'exec_prefix': exec_prefix,
                             'userbase': self.install_userbase,
                             'usersite': self.install_usersite,
+                            'libarchsuffix': self.libarchsuffix,
                            }
         self.expand_basedirs()
 
@@ -666,7 +666,8 @@ class install (Command):
                     ('install_headers', has_headers),
                     ('install_scripts', has_scripts),
                     ('install_data',    has_data),
-                    ('install_egg_info', lambda self:True),
                    ]
+    if not os.environ.has_key('PKGSRC_PYTHON_NO_EGG'):
+        sub_commands += [('install_egg_info', lambda self:True),]
 
 # class install
