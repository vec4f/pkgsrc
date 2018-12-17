package pkglint

import "gopkg.in/check.v1"

func (s *Suite) Test_MkLineChecker_Check__url2pkg(c *check.C) {
	t := s.Init(c)

	t.SetupVartypes()

	mkline := t.NewMkLine("filename.mk", 1, "# url2pkg-marker")

	MkLineChecker{mkline}.Check()

	t.CheckOutputLines(
		"ERROR: filename.mk:1: This comment indicates unfinished work (url2pkg).")
}

func (s *Suite) Test_MkLineChecker_Check__buildlink3_include_prefs(c *check.C) {
	t := s.Init(c)

	t.SetupVartypes()

	t.CreateFileLines("mk/bsd.prefs.mk")
	mklines := t.SetupFileMkLines("category/package/buildlink3.mk",
		".include \"../../mk/bsd.prefs.mk\"")
	// If the buildlink3.mk file doesn't actually exist, resolving the
	// relative path fails since that depends on the actual file system,
	// not on syntactical paths; see os.Stat in CheckRelativePath.
	//
	// TODO: Refactor relpath to be independent of a filesystem.

	MkLineChecker{mklines.mklines[0]}.Check()

	t.CheckOutputLines(
		"NOTE: ~/category/package/buildlink3.mk:1: For efficiency reasons, " +
			"please include bsd.fast.prefs.mk instead of bsd.prefs.mk.")
}

func (s *Suite) Test_MkLineChecker_checkInclude(c *check.C) {
	t := s.Init(c)

	t.SetupVartypes()

	t.CreateFileLines("pkgtools/x11-links/buildlink3.mk")
	t.CreateFileLines("graphics/jpeg/buildlink3.mk")
	t.CreateFileLines("devel/intltool/buildlink3.mk")
	t.CreateFileLines("devel/intltool/builtin.mk")
	mklines := t.SetupFileMkLines("category/package/filename.mk",
		MkRcsID,
		"",
		".include \"../../pkgtools/x11-links/buildlink3.mk\"",
		".include \"../../graphics/jpeg/buildlink3.mk\"",
		".include \"../../devel/intltool/buildlink3.mk\"",
		".include \"../../devel/intltool/builtin.mk\"")

	mklines.Check()

	t.CheckOutputLines(
		"ERROR: ~/category/package/filename.mk:3: "+
			"../../pkgtools/x11-links/buildlink3.mk must not be included directly. "+
			"Include \"../../mk/x11.buildlink3.mk\" instead.",
		"ERROR: ~/category/package/filename.mk:4: "+
			"../../graphics/jpeg/buildlink3.mk must not be included directly. "+
			"Include \"../../mk/jpeg.buildlink3.mk\" instead.",
		"WARN: ~/category/package/filename.mk:5: "+
			"Please write \"USE_TOOLS+= intltool\" instead of this line.",
		"ERROR: ~/category/package/filename.mk:6: "+
			"../../devel/intltool/builtin.mk must not be included directly. "+
			"Include \"../../devel/intltool/buildlink3.mk\" instead.")
}

func (s *Suite) Test_MkLineChecker_checkInclude__Makefile(c *check.C) {
	t := s.Init(c)

	mkline := t.NewMkLine(t.File("Makefile"), 2, ".include \"../../other/package/Makefile\"")

	MkLineChecker{mkline}.checkInclude()

	t.CheckOutputLines(
		"ERROR: ~/Makefile:2: Relative path \"../../other/package/Makefile\" does not exist.",
		"ERROR: ~/Makefile:2: Other Makefiles must not be included directly.")
}

func (s *Suite) Test_MkLineChecker_checkInclude__Makefile_exists(c *check.C) {
	t := s.Init(c)

	t.CreateFileLines("other/existing/Makefile")
	t.SetupPackage("category/package",
		".include \"../../other/existing/Makefile\"",
		".include \"../../other/not-found/Makefile\"")

	G.checkdirPackage(t.File("category/package"))

	t.CheckOutputLines(
		"ERROR: ~/category/package/Makefile:20: Cannot read \"../../other/existing/Makefile\".")
}

func (s *Suite) Test_MkLineChecker_checkDirective(c *check.C) {
	t := s.Init(c)

	t.SetupVartypes()

	mklines := t.NewMkLines("category/package/filename.mk",
		MkRcsID,
		"",
		".for",
		".endfor",
		"",
		".if",
		".else don't",
		".endif invalid-arg",
		"",
		".ifdef FNAME_MK",
		".endif",
		".ifndef FNAME_MK",
		".endif",
		"",
		".for var in a b c",
		".endfor",
		".undef var")

	mklines.Check()

	t.CheckOutputLines(
		"ERROR: category/package/filename.mk:3: \".for\" requires arguments.",
		"ERROR: category/package/filename.mk:6: \".if\" requires arguments.",
		"ERROR: category/package/filename.mk:7: \".else\" does not take arguments. "+
			"If you meant \"else if\", use \".elif\".",
		"ERROR: category/package/filename.mk:8: \".endif\" does not take arguments.",
		"WARN: category/package/filename.mk:10: The \".ifdef\" directive is deprecated. "+
			"Please use \".if defined(FNAME_MK)\" instead.",
		"WARN: category/package/filename.mk:12: The \".ifndef\" directive is deprecated. "+
			"Please use \".if !defined(FNAME_MK)\" instead.",
		"NOTE: category/package/filename.mk:17: Using \".undef\" after a \".for\" loop is unnecessary.")
}

func (s *Suite) Test_MkLineChecker_checkDirective__for_loop_varname(c *check.C) {
	t := s.Init(c)

	t.SetupVartypes()

	mklines := t.NewMkLines("filename.mk",
		MkRcsID,
		"",
		".for VAR in a b c", // Should be lowercase.
		".endfor",
		"",
		".for _var_ in a b c", // Should be written without underscores.
		".endfor",
		"",
		".for .var. in a b c", // Should be written without dots.
		".endfor",
		"",
		".for ${VAR} in a b c", // The variable name really must be an identifier.
		".endfor")

	mklines.Check()

	t.CheckOutputLines(
		"WARN: filename.mk:3: The variable name \"VAR\" in the .for loop should not contain uppercase letters.",
		"WARN: filename.mk:6: Variable names starting with an underscore (_var_) are reserved for internal pkgsrc use.",
		"ERROR: filename.mk:9: Invalid variable name \".var.\".",
		"ERROR: filename.mk:12: Invalid variable name \"${VAR}\".")
}

func (s *Suite) Test_MkLineChecker_checkDependencyRule(c *check.C) {
	t := s.Init(c)

	t.SetupVartypes()

	mklines := t.NewMkLines("category/package/filename.mk",
		MkRcsID,
		"",
		".PHONY: target-1",
		"target-2: .PHONY",
		".ORDER: target-1 target-2",
		"target-1:",
		"target-2:",
		"target-3:")

	mklines.Check()

	t.CheckOutputLines(
		"WARN: category/package/filename.mk:8: Undeclared target \"target-3\".")
}

func (s *Suite) Test_MkLineChecker_checkVartype__simple_type(c *check.C) {
	t := s.Init(c)

	t.SetupCommandLine("-Wtypes")
	t.SetupVartypes()

	// Since COMMENT is defined in vardefs.go its type is certain instead of guessed.
	vartype := G.Pkgsrc.VariableType("COMMENT")

	c.Assert(vartype, check.NotNil)
	c.Check(vartype.basicType.name, equals, "Comment")
	c.Check(vartype.guessed, equals, false)
	c.Check(vartype.kindOfList, equals, lkNone)

	mkline := t.NewMkLine("Makefile", 123, "COMMENT=\tA nice package")
	MkLineChecker{mkline}.checkVartype(mkline.Varname(), mkline.Op(), mkline.Value(), mkline.VarassignComment())

	t.CheckOutputLines(
		"WARN: Makefile:123: COMMENT should not begin with \"A\".")
}

func (s *Suite) Test_MkLineChecker_checkVartype(c *check.C) {
	t := s.Init(c)

	t.SetupVartypes()
	mkline := t.NewMkLine("filename", 1, "DISTNAME=gcc-${GCC_VERSION}")

	MkLineChecker{mkline}.checkVartype("DISTNAME", opAssign, "gcc-${GCC_VERSION}", "")

	t.CheckOutputEmpty()
}

// The command line option -Wno-types can be used to suppress the type checks.
// Suppressing it is rarely needed and comes from Feb 12 2005 when this feature was introduced.
// Since then the type system has matured and proven effective.
func (s *Suite) Test_MkLineChecker_checkVartype__skip(c *check.C) {
	t := s.Init(c)

	t.SetupCommandLine("-Wno-types")
	t.SetupVartypes()
	mkline := t.NewMkLine("filename", 1, "DISTNAME=invalid:::distname")

	MkLineChecker{mkline}.Check()

	t.CheckOutputEmpty()
}

func (s *Suite) Test_MkLineChecker_checkVartype__append_to_non_list(c *check.C) {
	t := s.Init(c)

	t.SetupVartypes()
	mklines := t.NewMkLines("filename.mk",
		MkRcsID,
		"DISTNAME+=\tsuffix",
		"COMMENT=\tComment for",
		"COMMENT+=\tthe package")

	mklines.Check()

	t.CheckOutputLines(
		"WARN: filename.mk:2: The variable DISTNAME may not be appended to (only set, given a default value) in this file.",
		"WARN: filename.mk:2: The \"+=\" operator should only be used with lists, not with DISTNAME.")
}

// Pkglint once interpreted all lists as consisting of shell tokens,
// splitting this URL at the ampersand.
func (s *Suite) Test_MkLineChecker_checkVarassign__URL_with_shell_special_characters(c *check.C) {
	t := s.Init(c)

	G.Pkg = NewPackage(t.File("graphics/gimp-fix-ca"))
	t.SetupVartypes()
	mkline := t.NewMkLine("filename", 10, "MASTER_SITES=http://registry.gimp.org/file/fix-ca.c?action=download&id=9884&file=")

	MkLineChecker{mkline}.checkVarassign()

	t.CheckOutputEmpty()
}

func (s *Suite) Test_MkLineChecker_checkDirectiveCond(c *check.C) {
	t := s.Init(c)

	t.SetupCommandLine("-Wtypes")
	t.SetupVartypes()

	testCond := func(cond string, output ...string) {
		MkLineChecker{t.NewMkLine("filename", 1, cond)}.checkDirectiveCond()
		if len(output) > 0 {
			t.CheckOutputLines(output...)
		} else {
			t.CheckOutputEmpty()
		}
	}

	testCond(".if !empty(PKGSRC_COMPILER:Mmycc)",
		"WARN: filename:1: The pattern \"mycc\" cannot match any of "+
			"{ ccache ccc clang distcc f2c gcc hp icc ido "+
			"mipspro mipspro-ucode pcc sunpro xlc } for PKGSRC_COMPILER.")

	testCond(".elif ${A} != ${B}")

	testCond(".if ${HOMEPAGE} == \"mailto:someone@example.org\"",
		"WARN: filename:1: \"mailto:someone@example.org\" is not a valid URL.")

	testCond(".if !empty(PKGSRC_RUN_TEST:M[Y][eE][sS])",
		"WARN: filename:1: PKGSRC_RUN_TEST should be matched "+
			"against \"[yY][eE][sS]\" or \"[nN][oO]\", not \"[Y][eE][sS]\".")

	testCond(".if !empty(IS_BUILTIN.Xfixes:M[yY][eE][sS])")

	testCond(".if !empty(${IS_BUILTIN.Xfixes:M[yY][eE][sS]})",
		"WARN: filename:1: The empty() function takes a variable name as parameter, "+
			"not a variable expression.")

	testCond(".if ${PKGSRC_COMPILER} == \"msvc\"",
		"WARN: filename:1: \"msvc\" is not valid for PKGSRC_COMPILER. "+
			"Use one of { ccache ccc clang distcc f2c gcc hp icc ido mipspro mipspro-ucode pcc sunpro xlc } instead.",
		"WARN: filename:1: Use ${PKGSRC_COMPILER:Mmsvc} instead of the == operator.")

	testCond(".if ${PKG_LIBTOOL:Mlibtool}",
		"NOTE: filename:1: PKG_LIBTOOL should be compared using == instead of matching against \":Mlibtool\".")

	testCond(".if ${MACHINE_PLATFORM:MUnknownOS-*-*} || ${MACHINE_ARCH:Mx86}",
		"WARN: filename:1: "+
			"The pattern \"UnknownOS\" cannot match any of "+
			"{ AIX BSDOS Bitrig Cygwin Darwin DragonFly FreeBSD FreeMiNT GNUkFreeBSD HPUX Haiku "+
			"IRIX Interix Linux Minix MirBSD NetBSD OSF1 OpenBSD QNX SCO_SV SunOS UnixWare "+
			"} for the operating system part of MACHINE_PLATFORM.",
		"WARN: filename:1: "+
			"The pattern \"x86\" cannot match any of "+
			"{ aarch64 aarch64eb alpha amd64 arc arm arm26 arm32 cobalt coldfire convex dreamcast earm "+
			"earmeb earmhf earmhfeb earmv4 earmv4eb earmv5 earmv5eb earmv6 earmv6eb earmv6hf earmv6hfeb "+
			"earmv7 earmv7eb earmv7hf earmv7hfeb evbarm hpcmips hpcsh hppa hppa64 i386 i586 i686 ia64 "+
			"m68000 m68k m88k mips mips64 mips64eb mips64el mipseb mipsel mipsn32 mlrisc ns32k pc532 pmax "+
			"powerpc powerpc64 rs6000 s390 sh3eb sh3el sparc sparc64 vax x86_64 "+
			"} for MACHINE_ARCH.",
		"NOTE: filename:1: MACHINE_ARCH should be compared using == instead of matching against \":Mx86\".")

	testCond(".if ${MASTER_SITES:Mftp://*} == \"ftp://netbsd.org/\"")
}

func (s *Suite) Test_MkLineChecker_checkVarassign(c *check.C) {
	t := s.Init(c)

	t.SetupVartypes()

	G.Mk = t.NewMkLines("Makefile",
		MkRcsID,
		"ac_cv_libpari_libs+=\t-L${BUILDLINK_PREFIX.pari}/lib") // From math/clisp-pari/Makefile, rev. 1.8

	MkLineChecker{G.Mk.mklines[1]}.checkVarassign()

	t.CheckOutputLines(
		"WARN: Makefile:2: ac_cv_libpari_libs is defined but not used.")
}

func (s *Suite) Test_MkLineChecker_checkVarassignPermissions(c *check.C) {
	t := s.Init(c)

	t.SetupCommandLine("-Wall,no-space")
	t.SetupVartypes()
	mklines := t.NewMkLines("options.mk",
		MkRcsID,
		"PKG_DEVELOPER?= yes",
		"BUILD_DEFS?=    VARBASE")

	mklines.Check()

	t.CheckOutputLines(
		"WARN: options.mk:2: The variable PKG_DEVELOPER may not be given a default value by any package.",
		"WARN: options.mk:2: Please include \"../../mk/bsd.prefs.mk\" before using \"?=\".",
		"WARN: options.mk:3: The variable BUILD_DEFS may not be given a default value (only appended to) in this file.")
}

// Don't check the permissions for infrastructure files since they have their own rules.
func (s *Suite) Test_MkLineChecker_checkVarassignPermissions__infrastructure(c *check.C) {
	t := s.Init(c)

	t.SetupVartypes()
	t.CreateFileLines("mk/infra.mk",
		MkRcsID,
		"",
		"PKG_DEVELOPER?=\tyes")
	t.CreateFileLines("mk/bsd.pkg.mk")

	G.CheckDirent(t.File("mk/infra.mk"))

	t.CheckOutputEmpty()
}

func (s *Suite) Test_MkLineChecker_checkVarassignVaruse(c *check.C) {
	t := s.Init(c)

	t.SetupVartypes()

	mkline := t.NewMkLine("module.mk", 123, "PLIST_SUBST+=\tLOCALBASE=${LOCALBASE:Q}")

	MkLineChecker{mkline}.checkVarassignVaruse()

	t.CheckOutputLines(
		"WARN: module.mk:123: Please use PREFIX instead of LOCALBASE.",
		"NOTE: module.mk:123: The :Q operator isn't necessary for ${LOCALBASE} here.")
}

func (s *Suite) Test_MkLineChecker_checkVarusePermissions(c *check.C) {
	t := s.Init(c)

	t.SetupVartypes()
	mklines := t.NewMkLines("options.mk",
		MkRcsID,
		"COMMENT=\t${GAMES_USER}",
		"COMMENT:=\t${PKGBASE}",
		"PYPKGPREFIX=${PKGBASE}")
	G.Pkgsrc.UserDefinedVars.Define("GAMES_USER", mklines.mklines[0])

	mklines.Check()

	t.CheckOutputLines(
		"WARN: options.mk:3: PKGBASE should not be evaluated at load time.",
		"WARN: options.mk:4: The variable PYPKGPREFIX may not be set in this file; it would be ok in pyversion.mk.",
		"WARN: options.mk:4: PKGBASE should not be evaluated indirectly at load time.",
		"NOTE: options.mk:4: This variable value should be aligned to column 17.")
}

func (s *Suite) Test_MkLineChecker_checkVarusePermissions__load_time(c *check.C) {
	t := s.Init(c)

	t.SetupVartypes()
	mklines := t.NewMkLines("options.mk",
		MkRcsID,
		"WRKSRC:=${.CURDIR}",
		".if ${PKG_SYSCONFDIR.gdm} != \"etc\"",
		".endif")

	mklines.Check()

	// Evaluating PKG_SYSCONFDIR.* at load time is probably ok,
	// though pkglint cannot prove anything here.
	//
	// Evaluating .CURDIR at load time is definitely ok since it is defined from the beginning.
	t.CheckOutputLines(
		"NOTE: options.mk:2: This variable value should be aligned to column 17.")
}

func (s *Suite) Test_MkLineChecker_checkVarusePermissions__load_time_guessed(c *check.C) {
	t := s.Init(c)

	t.SetupVartypes()
	t.SetupTool("install", "", AtRunTime)
	mklines := t.NewMkLines("install-docfiles.mk",
		MkRcsID,
		"DOCFILES=\ta b c",
		"do-install:",
		".for f in ${DOCFILES}",
		"\tinstall -c ${WRKSRC}/${f} ${DESTDIR}${PREFIX}/${f}",
		".endfor")

	mklines.Check()

	// No warning for using DOCFILES at compile-time. Since the variable
	// name is not one of the predefined names from vardefs.go, the
	// variable's type is guessed based on the name (see
	// Pkgsrc.VariableType).
	//
	// These guessed variables are typically defined and used only in
	// a single file, and in this context, mistakes are usually found
	// quickly.
	t.CheckOutputEmpty()
}

func (s *Suite) Test_MkLineChecker_checkVarusePermissions__PKGREVISION(c *check.C) {
	t := s.Init(c)

	t.SetupVartypes()
	mklines := t.NewMkLines("any.mk",
		MkRcsID,
		// PKGREVISION may only be set in Makefile, not used at load time; see vardefs.go.
		".if defined(PKGREVISION)",
		".endif")

	mklines.Check()

	t.CheckOutputLines(
		"WARN: any.mk:2: PKGREVISION should not be evaluated at load time.",
		"WARN: any.mk:2: PKGREVISION may not be used in any file; it is a write-only variable.")
}

func (s *Suite) Test_MkLineChecker_Check__warn_varuse_LOCALBASE(c *check.C) {
	t := s.Init(c)

	t.SetupVartypes()
	mkline := t.NewMkLine("options.mk", 56, "PKGNAME=${LOCALBASE}")

	MkLineChecker{mkline}.Check()

	t.CheckOutputLines(
		"WARN: options.mk:56: Please use PREFIX instead of LOCALBASE.")
}

func (s *Suite) Test_MkLineChecker_CheckRelativePkgdir(c *check.C) {
	t := s.Init(c)

	t.CreateFileLines("other/package/Makefile")
	// Must be in the filesystem because of directory references.
	mklines := t.SetupFileMkLines("category/package/Makefile",
		"# dummy")
	ck := MkLineChecker{mklines.mklines[0]}

	ck.CheckRelativePkgdir("../pkgbase")
	ck.CheckRelativePkgdir("../../other/package")
	ck.CheckRelativePkgdir("../../other/does-not-exist")

	// FIXME: The diagnostics for does-not-exist are redundant.
	t.CheckOutputLines(
		"ERROR: ~/category/package/Makefile:1: Relative path \"../pkgbase\" does not exist.",
		"WARN: ~/category/package/Makefile:1: \"../pkgbase\" is not a valid relative package directory.",
		"ERROR: ~/category/package/Makefile:1: Relative path \"../../other/does-not-exist\" does not exist.",
		"ERROR: ~/category/package/Makefile:1: There is no package in \"other/does-not-exist\".")
}

// PR pkg/46570, item 2
func (s *Suite) Test_MkLineChecker__unclosed_varuse(c *check.C) {
	t := s.Init(c)

	mklines := t.NewMkLines("Makefile",
		MkRcsID,
		"EGDIRS=\t${EGDIR/apparmor.d ${EGDIR/dbus-1/system.d ${EGDIR/pam.d")

	mklines.Check()

	t.CheckOutputLines(
		"WARN: Makefile:2: Unclosed Make variable starting at \"${EGDIR/apparmor.d $...\".",
		"WARN: Makefile:2: EGDIRS is defined but not used.",

		// XXX: This warning is redundant because of the "Unclosed" warning above.
		"WARN: Makefile:2: Internal pkglint error in MkLine.Tokenize at "+
			"\"${EGDIR/apparmor.d ${EGDIR/dbus-1/system.d ${EGDIR/pam.d\".")
}

func (s *Suite) Test_MkLineChecker_Check__varuse_modifier_L(c *check.C) {
	t := s.Init(c)

	t.SetupVartypes()
	mklines := t.NewMkLines("x11/xkeyboard-config/Makefile",
		"FILES_SUBST+=XKBCOMP_SYMLINK=${${XKBBASE}/xkbcomp:L:Q}",
		"FILES_SUBST+=XKBCOMP_SYMLINK=${${XKBBASE}/xkbcomp:Q}")

	MkLineChecker{mklines.mklines[0]}.Check()
	MkLineChecker{mklines.mklines[1]}.Check()

	// In line 1, don't warn that ${XKBBASE}/xkbcomp is used but not defined.
	// This is because the :L modifier interprets everything before as an expression
	// instead of a variable name.
	//
	// In line 2 the :L modifier is missing, therefore ${XKBBASE}/xkbcomp is the
	// name of another variable, and that variable is not known. Only XKBBASE is known.
	//
	// FIXME: The below warnings are wrong because the MkParser does not recognize the
	// slash as part of a variable name. Because of that, parsing stops before the $.
	// The warning "Unclosed Make variable" wrongly assumes that any parse error from
	// a variable use is because of unclosed braces, which it isn't in this case.h
	t.CheckOutputLines(
		"WARN: x11/xkeyboard-config/Makefile:2: Unclosed Make variable starting at \"${${XKBBASE}/xkbcomp...\".",
		"WARN: x11/xkeyboard-config/Makefile:2: Unclosed Make variable starting at \"${${XKBBASE}/xkbcomp...\".")
}

func (s *Suite) Test_MkLineChecker_checkDirectiveCond__comparison_with_shell_command(c *check.C) {
	t := s.Init(c)

	t.SetupVartypes()
	mklines := t.NewMkLines("security/openssl/Makefile",
		MkRcsID,
		".if ${PKGSRC_COMPILER} == \"gcc\" && ${CC} == \"cc\"",
		".endif")

	mklines.Check()

	// Don't warn about unknown shell command "cc".
	t.CheckOutputLines(
		"WARN: security/openssl/Makefile:2: Use ${PKGSRC_COMPILER:Mgcc} instead of the == operator.")
}

func (s *Suite) Test_MkLineChecker_checkDirectiveCondEmpty(c *check.C) {
	t := s.Init(c)

	t.SetupVartypes()
	mkline := t.NewMkLine("module.mk", 123, ".if ${PKGPATH} == \"category/package\"")
	ck := MkLineChecker{mkline}

	// FIXME: checkDirectiveCondEmpty cannot know whether it is empty(...) or !empty(...).
	// It must know that to generate the proper diagnostics.

	ck.checkDirectiveCondEmpty(NewMkVarUse("PKGPATH", "Mpattern"))

	// When the pattern contains placeholders, it cannot be converted to == or !=.
	ck.checkDirectiveCondEmpty(NewMkVarUse("PKGPATH", "Mpa*n"))

	ck.checkDirectiveCondEmpty(NewMkVarUse("PKGPATH", "tl", "Mpattern"))

	ck.checkDirectiveCondEmpty(NewMkVarUse("PKGPATH", "Ncategory/package"))

	// ${PKGPATH:None:Ntwo} is a short variant of ${PKGPATH} != "one" && ${PKGPATH} != "two",
	// therefore no note is logged in this case.
	ck.checkDirectiveCondEmpty(NewMkVarUse("PKGPATH", "None", "Ntwo"))

	// Note: this combination doesn't make sense since the patterns "one" and "two" don't overlap.
	ck.checkDirectiveCondEmpty(NewMkVarUse("PKGPATH", "Mone", "Mtwo"))

	t.CheckOutputLines(
		"NOTE: module.mk:123: PKGPATH should be compared using == instead of matching against \":Mpattern\".",
		"NOTE: module.mk:123: PKGPATH should be compared using == instead of matching against \":Mpattern\".",
		"NOTE: module.mk:123: PKGPATH should be compared using != instead of matching against \":Ncategory/package\".",
		"NOTE: module.mk:123: PKGPATH should be compared using == instead of matching against \":Mone\".",
		"NOTE: module.mk:123: PKGPATH should be compared using == instead of matching against \":Mtwo\".")
}

func (s *Suite) Test_MkLineChecker_checkDirectiveCond__comparing_PKGSRC_COMPILER_with_eqeq(c *check.C) {
	t := s.Init(c)

	t.SetupVartypes()
	G.Mk = t.NewMkLines("audio/pulseaudio/Makefile",
		MkRcsID,
		".if ${OPSYS} == \"Darwin\" && ${PKGSRC_COMPILER} == \"clang\"",
		".endif")

	G.Mk.Check()

	t.CheckOutputLines(
		"WARN: audio/pulseaudio/Makefile:2: Use ${PKGSRC_COMPILER:Mclang} instead of the == operator.")
}

func (s *Suite) Test_MkLineChecker_checkVartype__CFLAGS_with_backticks(c *check.C) {
	t := s.Init(c)

	t.SetupVartypes()
	G.Mk = t.NewMkLines("chat/pidgin-icb/Makefile",
		MkRcsID,
		"CFLAGS+=\t`pkg-config pidgin --cflags`")
	mkline := G.Mk.mklines[1]

	words, rest := splitIntoMkWords(mkline.Line, mkline.Value())

	c.Check(words, deepEquals, []string{"`pkg-config pidgin --cflags`"})
	c.Check(rest, equals, "")

	ck := MkLineChecker{G.Mk.mklines[1]}
	ck.checkVartype("CFLAGS", opAssignAppend, "`pkg-config pidgin --cflags`", "")

	// No warning about "`pkg-config" being an unknown CFlag.
	t.CheckOutputEmpty()
}

// See PR 46570, Ctrl+F "4. Shell quoting".
// Pkglint is correct, since the shell sees this definition for
// CPPFLAGS as three words, not one word.
func (s *Suite) Test_MkLineChecker_checkVartype__CFLAGS(c *check.C) {
	t := s.Init(c)

	t.SetupVartypes()
	mklines := t.NewMkLines("Makefile",
		MkRcsID,
		"CPPFLAGS.SunOS+=\t-DPIPECOMMAND=\\\"/usr/sbin/sendmail -bs %s\\\"")

	mklines.Check()

	t.CheckOutputLines(
		"WARN: Makefile:2: Unknown compiler flag \"-bs\".",
		"WARN: Makefile:2: Compiler flag \"%s\\\\\\\"\" should start with a hyphen.")
}

func (s *Suite) Test_MkLineChecker_checkDirectiveIndentation__autofix(c *check.C) {
	t := s.Init(c)

	t.SetupCommandLine("--autofix", "-Wspace")
	lines := t.SetupFileLines("filename.mk",
		MkRcsID,
		".if defined(A)",
		".for a in ${A}",
		".if defined(C)",
		".endif",
		".endfor",
		".endif")
	mklines := NewMkLines(lines)

	mklines.Check()

	t.CheckOutputLines(
		"AUTOFIX: ~/filename.mk:3: Replacing \".\" with \".  \".",
		"AUTOFIX: ~/filename.mk:4: Replacing \".\" with \".    \".",
		"AUTOFIX: ~/filename.mk:5: Replacing \".\" with \".    \".",
		"AUTOFIX: ~/filename.mk:6: Replacing \".\" with \".  \".")
	t.CheckFileLines("filename.mk",
		"# $"+"NetBSD$",
		".if defined(A)",
		".  for a in ${A}",
		".    if defined(C)",
		".    endif",
		".  endfor",
		".endif")
}

// Up to 2018-01-28, pkglint applied the autofix also to the continuation
// lines, which is incorrect. It replaced the dot in "4.*" with spaces.
func (s *Suite) Test_MkLineChecker_checkDirectiveIndentation__autofix_multiline(c *check.C) {
	t := s.Init(c)

	t.SetupCommandLine("-Wall", "--autofix")
	t.SetupVartypes()
	mklines := t.SetupFileMkLines("options.mk",
		MkRcsID,
		".if ${PKGNAME} == pkgname",
		".if \\",
		"   ${PLATFORM:MNetBSD-4.*}",
		".endif",
		".endif")

	mklines.Check()

	t.CheckOutputLines(
		"AUTOFIX: ~/options.mk:3: Replacing \".\" with \".  \".",
		"AUTOFIX: ~/options.mk:5: Replacing \".\" with \".  \".")

	t.CheckFileLines("options.mk",
		MkRcsID,
		".if ${PKGNAME} == pkgname",
		".  if \\",
		"   ${PLATFORM:MNetBSD-4.*}",
		".  endif",
		".endif")
}

func (s *Suite) Test_MkLineChecker_CheckVaruseShellword(c *check.C) {
	t := s.Init(c)

	t.SetupVartypes()
	mklines := t.SetupFileMkLines("options.mk",
		MkRcsID,
		"GOPATH=\t${WRKDIR}",
		"do-build:",
		"\tcd ${WRKSRC} && GOPATH=${GOPATH} PATH=${PATH} :")

	mklines.Check()

	// For WRKSRC and GOPATH, no quoting is necessary since pkgsrc directories by
	// definition don't contain special characters. Therefore they don't need the
	// :Q, not even when used as part of a shell word.

	// For PATH, the quoting is necessary because it may contain directories outside
	// of pkgsrc, and these may contain special characters.

	t.CheckOutputLines(
		"WARN: ~/options.mk:4: The variable PATH should be quoted as part of a shell word.")
}

func (s *Suite) Test_MkLineChecker_CheckVaruseShellword__mstar(c *check.C) {
	t := s.Init(c)

	t.SetupCommandLine("-Wall,no-space")
	t.SetupVartypes()
	mklines := t.SetupFileMkLines("options.mk",
		MkRcsID,
		"CONFIGURE_ARGS+=        ${CFLAGS:Q}",
		"CONFIGURE_ARGS+=        ${CFLAGS:M*:Q}",
		"CONFIGURE_ARGS+=        ${ADA_FLAGS:Q}",
		"CONFIGURE_ARGS+=        ${ADA_FLAGS:M*:Q}",
		"CONFIGURE_ENV+=         ${CFLAGS:Q}",
		"CONFIGURE_ENV+=         ${CFLAGS:M*:Q}",
		"CONFIGURE_ENV+=         ${ADA_FLAGS:Q}",
		"CONFIGURE_ENV+=         ${ADA_FLAGS:M*:Q}")

	mklines.Check()

	// FIXME: There should be some notes and warnings about missing :M*;
	// these are currently prevented by the PERL5 case in VariableNeedsQuoting.
	t.CheckOutputLines(
		"WARN: ~/options.mk:4: ADA_FLAGS is used but not defined.")
}

func (s *Suite) Test_MkLineChecker_CheckVaruseShellword__mstar_not_needed(c *check.C) {
	t := s.Init(c)

	t.SetupCommandLine("-Wall,no-space")
	pkg := t.SetupPackage("category/package",
		"MAKE_FLAGS+=\tCFLAGS=${CFLAGS:M*:Q}",
		"MAKE_FLAGS+=\tLFLAGS=${LDFLAGS:M*:Q}")
	G.Pkgsrc.LoadInfrastructure()
	// FIXME: It is too easy to forget this important call.

	// This package is guaranteed to not use GNU_CONFIGURE.
	// Since the :M* hack is only needed for GNU_CONFIGURE, it is not necessary here.
	G.CheckDirent(pkg)

	// FIXME: Duplicate diagnostics.
	t.CheckOutputLines(
		"NOTE: ~/category/package/Makefile:20: The :M* modifier is not needed here.",
		"NOTE: ~/category/package/Makefile:20: The :M* modifier is not needed here.",
		"NOTE: ~/category/package/Makefile:21: The :M* modifier is not needed here.",
		"NOTE: ~/category/package/Makefile:21: The :M* modifier is not needed here.")
}

func (s *Suite) Test_MkLineChecker_CheckVaruseShellword__q_not_needed(c *check.C) {
	t := s.Init(c)

	pkg := t.SetupPackage("category/package",
		"MASTER_SITES=\t${HOMEPAGE:Q}")
	G.Pkgsrc.LoadInfrastructure()

	G.CheckDirent(pkg)

	t.CheckOutputLines(
		"NOTE: ~/category/package/Makefile:5: The :Q operator isn't necessary for ${HOMEPAGE} here.")
}

// The ${VARNAME:=suffix} expression should only be used with lists.
// It typically appears in MASTER_SITE definitions.
func (s *Suite) Test_MkLineChecker_CheckVaruse__eq_nonlist(c *check.C) {
	t := s.Init(c)

	t.SetupVartypes()
	t.SetupMasterSite("MASTER_SITE_GITHUB", "https://github.com/")
	mklines := t.SetupFileMkLines("options.mk",
		MkRcsID,
		"WRKSRC=\t\t${WRKDIR:=/subdir}",
		"MASTER_SITES=\t${MASTER_SITE_GITHUB:=organization/}")

	mklines.Check()

	t.CheckOutputLines(
		"WARN: ~/options.mk:2: The :from=to modifier should only be used with lists, not with WRKDIR.")
}

func (s *Suite) Test_MkLineChecker_CheckVaruse__for(c *check.C) {
	t := s.Init(c)

	t.SetupVartypes()
	t.SetupMasterSite("MASTER_SITE_GITHUB", "https://github.com/")
	mklines := t.SetupFileMkLines("options.mk",
		MkRcsID,
		".for var in a b c",
		"\t: ${var}",
		".endfor")

	mklines.Check()

	t.CheckOutputEmpty()
}

// When a parameterized variable is defined in the pkgsrc infrastructure,
// it does not generate a warning about being "used but not defined".
// Even if the variable parameter differs, like .Linux and .SunOS in this
// case. This pattern is typical for pkgsrc, therefore pkglint doesn't
// check that the variable names match exactly.
func (s *Suite) Test_MkLineChecker_CheckVaruse__varcanon(c *check.C) {
	t := s.Init(c)

	t.SetupVartypes()
	t.SetupPkgsrc()
	t.CreateFileLines("mk/sys-vars.mk",
		MkRcsID,
		"CPPPATH.Linux=\t/usr/bin/cpp")
	G.Pkgsrc.LoadInfrastructure()

	ck := MkLineChecker{t.NewMkLine("module.mk", 101, "COMMENT=\t${CPPPATH.SunOS}")}

	ck.CheckVaruse(NewMkVarUse("CPPPATH.SunOS"), &VarUseContext{
		vartype: &Vartype{
			kindOfList: lkNone,
			basicType:  BtPathname,
			aclEntries: nil,
			guessed:    true,
		},
		time:       vucTimeRun,
		quoting:    vucQuotPlain,
		IsWordPart: false,
	})

	t.CheckOutputEmpty()
}

// Any variable that is defined in the pkgsrc infrastructure in mk/**/*.mk is
// considered defined, and no "used but not defined" warning is logged for it.
//
// See Pkgsrc.loadUntypedVars.
func (s *Suite) Test_MkLineChecker_CheckVaruse__defined_in_infrastructure(c *check.C) {
	t := s.Init(c)

	t.SetupPkgsrc()
	t.SetupVartypes()
	t.CreateFileLines("mk/deeply/nested/infra.mk",
		MkRcsID,
		"INFRA_VAR?=\tvalue")
	G.Pkgsrc.LoadInfrastructure()
	mklines := t.SetupFileMkLines("category/package/module.mk",
		MkRcsID,
		"do-fetch:",
		"\t: ${INFRA_VAR} ${UNDEFINED}")

	mklines.Check()

	t.CheckOutputLines(
		"WARN: ~/category/package/module.mk:3: UNDEFINED is used but not defined.")
}

func (s *Suite) Test_MkLineChecker_CheckVaruse__build_defs(c *check.C) {
	t := s.Init(c)

	// XXX: This paragraph should not be necessary since VARBASE and X11_TYPE
	// are also defined in vardefs.go.
	t.SetupPkgsrc()
	t.CreateFileLines("mk/defaults/mk.conf",
		"VARBASE?= /usr/pkg/var")
	G.Pkgsrc.LoadInfrastructure()

	t.SetupCommandLine("-Wall,no-space")
	t.SetupVartypes()
	mklines := t.SetupFileMkLines("options.mk",
		MkRcsID,
		"COMMENT=                ${VARBASE} ${X11_TYPE}",
		"PKG_FAIL_REASON+=       ${VARBASE} ${X11_TYPE}",
		"BUILD_DEFS+=            X11_TYPE")

	mklines.Check()

	t.CheckOutputLines(
		"WARN: ~/options.mk:2: The user-defined variable VARBASE is used but not added to BUILD_DEFS.")
}

func (s *Suite) Test_MkLineChecker_CheckVaruse__complicated_range(c *check.C) {
	t := s.Init(c)

	t.SetupCommandLine("--show-autofix", "--source")
	t.SetupVartypes()
	mkline := t.NewMkLine("mk/compiler/gcc.mk", 150,
		"CC:=\t${CC:C/^/_asdf_/1:M_asdf_*:S/^_asdf_//}")

	MkLineChecker{mkline}.Check()

	// FIXME: The check is called two times, even though it only produces a single NOTE.
	t.CheckOutputLines(
		"NOTE: mk/compiler/gcc.mk:150: "+
			"The modifier \":C/^/_asdf_/1:M_asdf_*:S/^_asdf_//\" can be written as \":[1]\".",
		"AUTOFIX: mk/compiler/gcc.mk:150: "+
			"Replacing \":C/^/_asdf_/1:M_asdf_*:S/^_asdf_//\" with \":[1]\".",
		"-\tCC:=\t${CC:C/^/_asdf_/1:M_asdf_*:S/^_asdf_//}",
		"+\tCC:=\t${CC:[1]}")
}

func (s *Suite) Test_MkLineChecker_CheckVaruse__deprecated_PKG_DEBUG(c *check.C) {
	t := s.Init(c)

	t.SetupVartypes()
	G.Pkgsrc.initDeprecatedVars()

	mkline := t.NewMkLine("module.mk", 123,
		"\t${_PKG_SILENT}${_PKG_DEBUG} :")

	MkLineChecker{mkline}.Check()

	t.CheckOutputLines(
		"WARN: module.mk:123: Use of \"_PKG_SILENT\" is deprecated. Use RUN (with more error checking) instead.",
		"WARN: module.mk:123: Use of \"_PKG_DEBUG\" is deprecated. Use RUN (with more error checking) instead.")
}

// PR 46570, item "15. net/uucp/Makefile has a make loop"
func (s *Suite) Test_MkLineChecker_checkVaruseUndefined__indirect_variables(c *check.C) {
	t := s.Init(c)

	t.SetupTool("echo", "ECHO", AfterPrefsMk)
	mkline := t.NewMkLine("net/uucp/Makefile", 123, "\techo ${UUCP_${var}}")

	MkLineChecker{mkline}.Check()

	// No warning about UUCP_${var} being used but not defined.
	//
	// Normally, parameterized variables use a dot instead of an underscore as separator.
	// This is one of the few other cases. Pkglint doesn't warn about dynamic variable
	// names like UUCP_${var} or SITES_${distfile}.
	//
	// It does warn about simple variable names though, like ${var} in this example.
	t.CheckOutputLines(
		"WARN: net/uucp/Makefile:123: var is used but not defined.")
}

func (s *Suite) Test_MkLineChecker_checkVarassignSpecific(c *check.C) {
	t := s.Init(c)

	t.SetupPkgsrc()
	G.Pkgsrc.LoadInfrastructure()
	t.SetupCommandLine("-Wall,no-space")
	t.SetupVartypes()
	mklines := t.SetupFileMkLines("module.mk",
		MkRcsID,
		"EGDIR=                  ${PREFIX}/etc/rc.d",
		"_TOOLS_VARNAME.sed=     SED",
		"DIST_SUBDIR=            ${PKGNAME}",
		"WRKSRC=                 ${PKGNAME}",
		"SITES_distfile.tar.gz=  ${MASTER_SITE_GITHUB:=user/}",
		// TODO: The first of the below assignments should be flagged as redundant by RedundantScope;
		// that check is currently only implemented for package Makefiles, not for other files.
		"PYTHON_VERSIONS_ACCEPTED= -13",
		"PYTHON_VERSIONS_ACCEPTED= 27 36")

	mklines.Check()

	// TODO: Split this test into several, one for each topic.
	t.CheckOutputLines(
		"WARN: ~/module.mk:2: Please use the RCD_SCRIPTS mechanism to install rc.d scripts automatically to ${RCD_SCRIPTS_EXAMPLEDIR}.",
		"WARN: ~/module.mk:3: _TOOLS_VARNAME.sed is defined but not used.",
		"WARN: ~/module.mk:3: Variable names starting with an underscore (_TOOLS_VARNAME.sed) are reserved for internal pkgsrc use.",
		"WARN: ~/module.mk:4: PKGNAME should not be used in DIST_SUBDIR as it includes the PKGREVISION. Please use PKGNAME_NOREV instead.",
		"WARN: ~/module.mk:5: PKGNAME should not be used in WRKSRC as it includes the PKGREVISION. Please use PKGNAME_NOREV instead.",
		"WARN: ~/module.mk:6: SITES_distfile.tar.gz is defined but not used.",
		"WARN: ~/module.mk:6: SITES_* is deprecated. Please use SITES.* instead.",
		"WARN: ~/module.mk:7: The variable PYTHON_VERSIONS_ACCEPTED may not be set "+
			"(only given a default value, appended to) in this file; it would be ok in Makefile, Makefile.common, options.mk.",
		"WARN: ~/module.mk:7: Invalid version number \"-13\".",
		"ERROR: ~/module.mk:7: All values for PYTHON_VERSIONS_ACCEPTED must be positive integers.",
		"WARN: ~/module.mk:8: The variable PYTHON_VERSIONS_ACCEPTED may not be set "+
			"(only given a default value, appended to) in this file; it would be ok in Makefile, Makefile.common, options.mk.",
		"WARN: ~/module.mk:8: The values for PYTHON_VERSIONS_ACCEPTED should be in decreasing order.")
}

func (s *Suite) Test_MkLineChecker_checkText(c *check.C) {
	t := s.Init(c)

	t.SetupPkgsrc()
	G.Pkgsrc.LoadInfrastructure()

	t.SetupCommandLine("-Wall,no-space")
	mklines := t.SetupFileMkLines("module.mk",
		MkRcsID,
		"CFLAGS+=                -Wl,--rpath,${PREFIX}/lib",
		"PKG_FAIL_REASON+=       \"Group ${GAMEGRP} doesn't exist.\"")

	mklines.Check()

	t.CheckOutputLines(
		"WARN: ~/module.mk:2: Please use ${COMPILER_RPATH_FLAG} instead of \"-Wl,--rpath,\".",
		"WARN: ~/module.mk:3: Use of \"GAMEGRP\" is deprecated. Use GAMES_GROUP instead.")
}

func (s *Suite) Test_MkLineChecker_checkText__WRKSRC(c *check.C) {
	t := s.Init(c)

	t.SetupCommandLine("-Wall", "--explain")
	mklines := t.SetupFileMkLines("module.mk",
		MkRcsID,
		"pre-configure:",
		"\tcd ${WRKSRC}/..")

	mklines.Check()

	t.CheckOutputLines(
		"WARN: ~/module.mk:3: Building the package should take place entirely inside ${WRKSRC}, not \"${WRKSRC}/..\".",
		"",
		"\tWRKSRC should be defined so that there is no need to do anything",
		"\toutside of this directory.",
		"",
		"\tExample:",
		"",
		"\t\tWRKSRC=\t${WRKDIR}",
		"\t\tCONFIGURE_DIRS=\t${WRKSRC}/lib ${WRKSRC}/src",
		"\t\tBUILD_DIRS=\t${WRKSRC}/lib ${WRKSRC}/src ${WRKSRC}/cmd",
		"",
		"\tSee the pkgsrc guide, section \"Directories used during the build",
		"\tprocess\":",
		"\thttps://www.NetBSD.org/docs/pkgsrc/pkgsrc.html#build.builddirs",
		"",
		"WARN: ~/module.mk:3: WRKSRC is used but not defined.")
}

func (s *Suite) Test_MkLineChecker_CheckRelativePath(c *check.C) {
	t := s.Init(c)

	t.SetupPkgsrc()
	G.Pkgsrc.LoadInfrastructure()
	t.CreateFileLines("wip/package/Makefile")
	t.CreateFileLines("wip/package/module.mk")
	mklines := t.SetupFileMkLines("category/package/module.mk",
		MkRcsID,
		"DEPENDS+=       wip-package-[0-9]*:../../wip/package",
		".include \"../../wip/package/module.mk\"",
		"",
		"DEPENDS+=       unresolvable-[0-9]*:../../lang/${LATEST_PYTHON}",
		".include \"../../lang/${LATEST_PYTHON}/module.mk\"",
		"",
		".include \"module.mk\"",
		".include \"../../category/../category/package/module.mk\"", // Oops
		".include \"../../mk/bsd.prefs.mk\"",
		".include \"../package/module.mk\"")

	mklines.Check()

	t.CheckOutputLines(
		"ERROR: ~/category/package/module.mk:2: A main pkgsrc package must not depend on a pkgsrc-wip package.",
		"ERROR: ~/category/package/module.mk:3: A main pkgsrc package must not depend on a pkgsrc-wip package.",
		"WARN: ~/category/package/module.mk:5: LATEST_PYTHON is used but not defined.",
		// TODO: This warning is unspecific, there is also a pkglint warning "should be ../../category/package".
		"WARN: ~/category/package/module.mk:11: Invalid relative path \"../package/module.mk\".")
}
