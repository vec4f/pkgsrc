$NetBSD$

Use compiler-rt instead of libgcc.
Pull in libcxx correctly.
Specify paths to system objects explicitly.

--- lib/Driver/ToolChains/Solaris.cpp.orig	2018-01-04 07:43:41.000000000 +0000
+++ lib/Driver/ToolChains/Solaris.cpp
@@ -49,8 +49,24 @@ void solaris::Linker::ConstructJob(Compi
                                    const InputInfoList &Inputs,
                                    const ArgList &Args,
                                    const char *LinkingOutput) const {
+  const Driver &D = getToolChain().getDriver();
   ArgStringList CmdArgs;
 
+  std::string LibPath = "/usr/lib/";
+  switch (getToolChain().getArch()) {
+  case llvm::Triple::x86:
+  case llvm::Triple::sparc:
+    break;
+  case llvm::Triple::x86_64:
+    LibPath += "amd64/";
+    break;
+  case llvm::Triple::sparcv9:
+    LibPath += "sparcv9/";
+    break;
+  default:
+    llvm_unreachable("Unsupported architecture");
+  }
+
   // Demangle C++ names in errors
   CmdArgs.push_back("-C");
 
@@ -69,7 +85,7 @@ void solaris::Linker::ConstructJob(Compi
     } else {
       CmdArgs.push_back("--dynamic-linker");
       CmdArgs.push_back(
-          Args.MakeArgString(getToolChain().GetFilePath("ld.so.1")));
+          Args.MakeArgString(LibPath + "ld.so.1"));
     }
   }
 
@@ -83,13 +99,11 @@ void solaris::Linker::ConstructJob(Compi
   if (!Args.hasArg(options::OPT_nostdlib, options::OPT_nostartfiles)) {
     if (!Args.hasArg(options::OPT_shared))
       CmdArgs.push_back(
-          Args.MakeArgString(getToolChain().GetFilePath("crt1.o")));
+          Args.MakeArgString(LibPath + "crt1.o"));
 
-    CmdArgs.push_back(Args.MakeArgString(getToolChain().GetFilePath("crti.o")));
-    CmdArgs.push_back(
-        Args.MakeArgString(getToolChain().GetFilePath("values-Xa.o")));
+    CmdArgs.push_back(Args.MakeArgString(LibPath + "crti.o"));
     CmdArgs.push_back(
-        Args.MakeArgString(getToolChain().GetFilePath("crtbegin.o")));
+        Args.MakeArgString(LibPath + "values-Xa.o"));
   }
 
   getToolChain().AddFilePathLibArgs(Args, CmdArgs);
@@ -100,21 +114,16 @@ void solaris::Linker::ConstructJob(Compi
   AddLinkerInputs(getToolChain(), Inputs, Args, CmdArgs, JA);
 
   if (!Args.hasArg(options::OPT_nostdlib, options::OPT_nodefaultlibs)) {
-    if (getToolChain().ShouldLinkCXXStdlib(Args))
-      getToolChain().AddCXXStdlibLibArgs(Args, CmdArgs);
-    CmdArgs.push_back("-lgcc_s");
-    CmdArgs.push_back("-lc");
-    if (!Args.hasArg(options::OPT_shared)) {
-      CmdArgs.push_back("-lgcc");
+    if (D.CCCIsCXX()) {
+      if (getToolChain().ShouldLinkCXXStdlib(Args))
+        getToolChain().AddCXXStdlibLibArgs(Args, CmdArgs);
       CmdArgs.push_back("-lm");
     }
+    CmdArgs.push_back("-lc");
+    AddRunTimeLibs(getToolChain(), D, CmdArgs, Args);
   }
 
-  if (!Args.hasArg(options::OPT_nostdlib, options::OPT_nostartfiles)) {
-    CmdArgs.push_back(
-        Args.MakeArgString(getToolChain().GetFilePath("crtend.o")));
-  }
-  CmdArgs.push_back(Args.MakeArgString(getToolChain().GetFilePath("crtn.o")));
+  CmdArgs.push_back(Args.MakeArgString(LibPath + "crtn.o"));
 
   getToolChain().addProfileRTLibs(Args, CmdArgs);
 
@@ -127,35 +136,9 @@ void solaris::Linker::ConstructJob(Compi
 Solaris::Solaris(const Driver &D, const llvm::Triple &Triple,
                  const ArgList &Args)
     : Generic_ELF(D, Triple, Args) {
-
-  GCCInstallation.init(Triple, Args);
-
-  path_list &Paths = getFilePaths();
-  if (GCCInstallation.isValid())
-    addPathIfExists(D, GCCInstallation.getInstallPath(), Paths);
-
-  addPathIfExists(D, getDriver().getInstalledDir(), Paths);
-  if (getDriver().getInstalledDir() != getDriver().Dir)
-    addPathIfExists(D, getDriver().Dir, Paths);
-
-  addPathIfExists(D, getDriver().SysRoot + getDriver().Dir + "/../lib", Paths);
-
-  std::string LibPath = "/usr/lib/";
-  switch (Triple.getArch()) {
-  case llvm::Triple::x86:
-  case llvm::Triple::sparc:
-    break;
-  case llvm::Triple::x86_64:
-    LibPath += "amd64/";
-    break;
-  case llvm::Triple::sparcv9:
-    LibPath += "sparcv9/";
-    break;
-  default:
-    llvm_unreachable("Unsupported architecture");
-  }
-
-  addPathIfExists(D, getDriver().SysRoot + LibPath, Paths);
+  // No special handling, the C runtime files are found directly above
+  // and crle handles adding the default system library paths if they
+  // are necessary.
 }
 
 Tool *Solaris::buildAssembler() const {
@@ -164,30 +147,41 @@ Tool *Solaris::buildAssembler() const {
 
 Tool *Solaris::buildLinker() const { return new tools::solaris::Linker(*this); }
 
+void Solaris::AddCXXStdlibLibArgs(const ArgList &Args,
+                                  ArgStringList &CmdArgs) const {
+  CXXStdlibType Type = GetCXXStdlibType(Args);
+
+  // Currently assumes pkgsrc layout where libcxx and libcxxabi are installed
+  // in the same prefixed directory that we are.
+  std::string LibPath;
+  LibPath = llvm::sys::path::parent_path(getDriver().getInstalledDir());
+  LibPath += "/lib";
+
+  switch (Type) {
+  case ToolChain::CST_Libcxx:
+    CmdArgs.push_back(Args.MakeArgString(StringRef("-L") + LibPath));
+    CmdArgs.push_back(Args.MakeArgString(StringRef("-R") + LibPath));
+    CmdArgs.push_back("-lc++");
+    CmdArgs.push_back("-lc++abi");
+    break;
+
+  // XXX: This won't work without being told exactly where libstdc++ is, but
+  // that changes based on distribution.  Maybe return ENOTSUP here?
+  case ToolChain::CST_Libstdcxx:
+    CmdArgs.push_back("-lstdc++");
+    break;
+  }
+}
+
 void Solaris::AddClangCXXStdlibIncludeArgs(const ArgList &DriverArgs,
                                            ArgStringList &CC1Args) const {
   if (DriverArgs.hasArg(options::OPT_nostdlibinc) ||
       DriverArgs.hasArg(options::OPT_nostdincxx))
     return;
 
-  // Include the support directory for things like xlocale and fudged system
-  // headers.
-  // FIXME: This is a weird mix of libc++ and libstdc++. We should also be
-  // checking the value of -stdlib= here and adding the includes for libc++
-  // rather than libstdc++ if it's requested.
-  addSystemInclude(DriverArgs, CC1Args, "/usr/include/c++/v1/support/solaris");
-
-  if (GCCInstallation.isValid()) {
-    GCCVersion Version = GCCInstallation.getVersion();
-    addSystemInclude(DriverArgs, CC1Args,
-                     getDriver().SysRoot + "/usr/gcc/" +
-                     Version.MajorStr + "." +
-                     Version.MinorStr +
-                     "/include/c++/" + Version.Text);
-    addSystemInclude(DriverArgs, CC1Args,
-                     getDriver().SysRoot + "/usr/gcc/" + Version.MajorStr +
-                     "." + Version.MinorStr + "/include/c++/" +
-                     Version.Text + "/" +
-                     GCCInstallation.getTriple().str());
-  }
+  // Currently assumes pkgsrc layout.
+  addSystemInclude(DriverArgs, CC1Args,
+                   llvm::sys::path::parent_path(getDriver().getInstalledDir())
+                   + "/include/c++/v1");
+  return;
 }
