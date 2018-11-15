$NetBSD: patch-Lib_sysconfig.py,v 1.1 2018/06/17 19:21:21 adam Exp $

chop "-RELEASE" or "-DEVELOPMENT" from release for DragonFly.

--- Lib/sysconfig.py.orig	2014-12-10 15:59:40.000000000 +0000
+++ Lib/sysconfig.py
@@ -9,8 +9,8 @@ _INSTALL_SCHEMES = {
     'posix_prefix': {
         'stdlib': '{base}/lib/python{py_version_short}',
         'platstdlib': '{platbase}/lib/python{py_version_short}',
-        'purelib': '{base}/lib/python{py_version_short}/site-packages',
-        'platlib': '{platbase}/lib/python{py_version_short}/site-packages',
+        'purelib': '{base}/lib/python{py_version_short}/site-packages{libarchsuffix}',
+        'platlib': '{platbase}/lib/python{py_version_short}/site-packages{libarchsuffix}',
         'include': '{base}/include/python{py_version_short}',
         'platinclude': '{platbase}/include/python{py_version_short}',
         'scripts': '{base}/bin',
@@ -277,7 +277,10 @@ def get_makefile_filename():
     """Return the path of the Makefile."""
     if _PYTHON_BUILD:
         return os.path.join(_PROJECT_BASE, "Makefile")
-    return os.path.join(get_path('platstdlib'), "config", "Makefile")
+    if sys.maxsize > 2**32:
+        return os.path.join(get_path('platstdlib'), "config", "@LIBARCHSUFFIX.64@".lstrip('/'), "Makefile")
+    else:
+        return os.path.join(get_path('platstdlib'), "config", "@LIBARCHSUFFIX.32@".lstrip('/'), "Makefile")
 
 # Issue #22199: retain undocumented private name for compatibility
 _get_makefile_filename = get_makefile_filename
@@ -465,6 +468,10 @@ def get_config_vars(*args):
         _CONFIG_VARS['base'] = _PREFIX
         _CONFIG_VARS['platbase'] = _EXEC_PREFIX
         _CONFIG_VARS['projectbase'] = _PROJECT_BASE
+        if sys.maxsize > 2**32:
+            _CONFIG_VARS['libarchsuffix'] = "@LIBARCHSUFFIX.64@"
+        else:
+            _CONFIG_VARS['libarchsuffix'] = "@LIBARCHSUFFIX.32@"
 
         if os.name in ('nt', 'os2'):
             _init_non_posix(_CONFIG_VARS)
@@ -607,6 +607,8 @@ def get_platform():
         osname, release, machine = _osx_support.get_platform_osx(
                                             get_config_vars(),
                                             osname, release, machine)
+    elif osname[:9] == "dragonfly":
+        release = str.split(release, '-')[0]
 
     return "%s-%s-%s" % (osname, release, machine)
 
