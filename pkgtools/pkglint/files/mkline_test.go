package main

import "gopkg.in/check.v1"

func (s *Suite) Test_NewMkLine__varassign(c *check.C) {
	t := s.Init(c)

	mkline := t.NewMkLine("test.mk", 101,
		"VARNAME.param?=value # varassign comment")

	c.Check(mkline.IsVarassign(), equals, true)
	c.Check(mkline.Varname(), equals, "VARNAME.param")
	c.Check(mkline.Varcanon(), equals, "VARNAME.*")
	c.Check(mkline.Varparam(), equals, "param")
	c.Check(mkline.Op(), equals, opAssignDefault)
	c.Check(mkline.Value(), equals, "value")
	c.Check(mkline.VarassignComment(), equals, "# varassign comment")
}

func (s *Suite) Test_NewMkLine__shellcmd(c *check.C) {
	t := s.Init(c)

	mkline := t.NewMkLine("test.mk", 101,
		"\tshell command # shell comment")

	c.Check(mkline.IsShellCommand(), equals, true)
	c.Check(mkline.ShellCommand(), equals, "shell command # shell comment")
}

func (s *Suite) Test_NewMkLine__comment(c *check.C) {
	t := s.Init(c)

	mkline := t.NewMkLine("test.mk", 101,
		"# whole line comment")

	c.Check(mkline.IsComment(), equals, true)
}

func (s *Suite) Test_NewMkLine__empty(c *check.C) {
	t := s.Init(c)

	mkline := t.NewMkLine("test.mk", 101, "")

	c.Check(mkline.IsEmpty(), equals, true)
}

func (s *Suite) Test_NewMkLine__directive(c *check.C) {
	t := s.Init(c)

	mkline := t.NewMkLine("test.mk", 101,
		".  if !empty(PKGNAME:M*-*) && ${RUBY_RAILS_SUPPORTED:[\\#]} == 1 # directive comment")

	c.Check(mkline.IsDirective(), equals, true)
	c.Check(mkline.Indent(), equals, "  ")
	c.Check(mkline.Directive(), equals, "if")
	c.Check(mkline.Args(), equals, "!empty(PKGNAME:M*-*) && ${RUBY_RAILS_SUPPORTED:[#]} == 1")
	c.Check(mkline.DirectiveComment(), equals, "directive comment")
}

func (s *Suite) Test_NewMkLine__include(c *check.C) {
	t := s.Init(c)

	mkline := t.NewMkLine("test.mk", 101,
		".    include \"../../mk/bsd.prefs.mk\" # include comment")

	c.Check(mkline.IsInclude(), equals, true)
	c.Check(mkline.Indent(), equals, "    ")
	c.Check(mkline.MustExist(), equals, true)
	c.Check(mkline.IncludedFile(), equals, "../../mk/bsd.prefs.mk")

	c.Check(mkline.IsSysinclude(), equals, false)
}

func (s *Suite) Test_NewMkLine__sysinclude(c *check.C) {
	t := s.Init(c)

	mkline := t.NewMkLine("test.mk", 101,
		".    include <subdir.mk> # sysinclude comment")

	c.Check(mkline.IsSysinclude(), equals, true)
	c.Check(mkline.Indent(), equals, "    ")
	c.Check(mkline.MustExist(), equals, true)
	c.Check(mkline.IncludedFile(), equals, "subdir.mk")

	c.Check(mkline.IsInclude(), equals, false)
}

func (s *Suite) Test_NewMkLine__dependency(c *check.C) {
	t := s.Init(c)

	mkline := t.NewMkLine("test.mk", 101,
		"target1 target2: source1 source2")

	c.Check(mkline.IsDependency(), equals, true)
	c.Check(mkline.Targets(), equals, "target1 target2")
	c.Check(mkline.Sources(), equals, "source1 source2")
}

func (s *Suite) Test_NewMkLine__dependency_space(c *check.C) {
	t := s.Init(c)

	mkline := t.NewMkLine("test.mk", 101,
		"target : source")

	c.Check(mkline.Targets(), equals, "target")
	c.Check(mkline.Sources(), equals, "source")
	t.CheckOutputLines(
		"NOTE: test.mk:101: Space before colon in dependency line.")
}

func (s *Suite) Test_NewMkLine__varassign_append(c *check.C) {
	t := s.Init(c)

	mkline := t.NewMkLine("test.mk", 101,
		"VARNAME+=value")

	c.Check(mkline.IsVarassign(), equals, true)
	c.Check(mkline.Varname(), equals, "VARNAME")
	c.Check(mkline.Varcanon(), equals, "VARNAME")
	c.Check(mkline.Varparam(), equals, "")
}

func (s *Suite) Test_NewMkLine__merge_conflict(c *check.C) {
	t := s.Init(c)

	mkline := t.NewMkLine("test.mk", 101,
		"<<<<<<<<<<<<<<<<<")

	// Merge conflicts are of neither type.
	c.Check(mkline.IsVarassign(), equals, false)
	c.Check(mkline.IsDirective(), equals, false)
	c.Check(mkline.IsInclude(), equals, false)
	c.Check(mkline.IsEmpty(), equals, false)
	c.Check(mkline.IsComment(), equals, false)
	c.Check(mkline.IsDependency(), equals, false)
	c.Check(mkline.IsShellCommand(), equals, false)
	c.Check(mkline.IsSysinclude(), equals, false)
}

func (s *Suite) Test_NewMkLine__autofix_space_after_varname(c *check.C) {
	t := s.Init(c)

	t.SetupCommandLine("-Wspace")
	filename := t.CreateFileLines("Makefile",
		MkRcsID,
		"VARNAME +=\t${VARNAME}",
		"VARNAME+ =\t${VARNAME+}",
		"VARNAME+ +=\t${VARNAME+}",
		"pkgbase := pkglint")

	CheckfileMk(filename)

	t.CheckOutputLines(
		"NOTE: ~/Makefile:2: Unnecessary space after variable name \"VARNAME\".",
		// FIXME: Don't say anything here because the spaced form is clearer that the compressed form.
		"NOTE: ~/Makefile:4: Unnecessary space after variable name \"VARNAME+\".")

	t.SetupCommandLine("-Wspace", "--autofix")

	CheckfileMk(filename)

	t.CheckOutputLines(
		"AUTOFIX: ~/Makefile:2: Replacing \"VARNAME +=\" with \"VARNAME+=\".",
		// FIXME: Don't fix anything here because the spaced form is clearer that the compressed form.
		"AUTOFIX: ~/Makefile:4: Replacing \"VARNAME+ +=\" with \"VARNAME++=\".")
	t.CheckFileLines("Makefile",
		MkRcsID+"",
		"VARNAME+=\t${VARNAME}",
		"VARNAME+ =\t${VARNAME+}",
		"VARNAME++=\t${VARNAME+}",
		"pkgbase := pkglint")
}

func (s *Suite) Test_NewMkLine__varname_with_hash(c *check.C) {
	t := s.Init(c)

	mkline := t.NewMkLine("Makefile", 123, "VARNAME.#=\tvalue")

	// Parse error because the # starts a comment.
	c.Check(mkline.IsVarassign(), equals, false)

	mkline2 := t.NewMkLine("Makefile", 123, "VARNAME.\\#=\tvalue")

	// FIXME: Varname() should be "VARNAME.#".
	c.Check(mkline2.IsVarassign(), equals, false)

	t.CheckOutputLines(
		"ERROR: Makefile:123: Unknown Makefile line format: \"VARNAME.#=\\tvalue\".",
		"ERROR: Makefile:123: Unknown Makefile line format: \"VARNAME.\\\\#=\\tvalue\".")
}

func (s *Suite) Test_MkLine_Varparam(c *check.C) {
	t := s.Init(c)

	mkline := t.NewMkLine("Makefile", 2, "SUBST_SED.${param}=\tvalue")

	varparam := mkline.Varparam()

	c.Check(varparam, equals, "${param}")
}

func (s *Suite) Test_MkLine_ValueAlign__commented(c *check.C) {
	t := s.Init(c)

	mkline := t.NewMkLine("Makefile", 2, "#SUBST_SED.${param}=\tvalue")

	valueAlign := mkline.ValueAlign()

	c.Check(mkline.IsCommentedVarassign(), equals, true)
	c.Check(valueAlign, equals, "#SUBST_SED.${param}=\t")
}

// Demonstrates how a simple condition is structured internally.
// For most of the checks, using cond.Walk is the simplest way to go.
func (s *Suite) Test_MkLine_Cond(c *check.C) {
	t := s.Init(c)

	mkline := t.NewMkLine("Makefile", 2, ".if ${VAR} == Value")

	cond := mkline.Cond()

	c.Check(cond.CompareVarStr.Var.varname, equals, "VAR")
	c.Check(cond.CompareVarStr.Str, equals, "Value")
	c.Check(mkline.Cond(), equals, cond)
}

func (s *Suite) Test_VarUseContext_String(c *check.C) {
	t := s.Init(c)

	t.SetupVartypes()
	vartype := G.Pkgsrc.VariableType("PKGNAME")
	vuc := &VarUseContext{vartype, vucTimeUnknown, vucQuotBackt, false}

	c.Check(vuc.String(), equals, "(Pkgname time:unknown quoting:backt wordpart:false)")
}

// In variable assignments, a plain '#' introduces a line comment, unless
// it is escaped by a backslash. In shell commands, on the other hand, it
// is interpreted literally.
func (s *Suite) Test_NewMkLine__number_sign(c *check.C) {
	t := s.Init(c)

	mklineVarassignEscaped := t.NewMkLine("filename", 1, "SED_CMD=\t's,\\#,hash,g'")

	c.Check(mklineVarassignEscaped.Varname(), equals, "SED_CMD")
	c.Check(mklineVarassignEscaped.Value(), equals, "'s,#,hash,g'")

	mklineCommandEscaped := t.NewMkLine("filename", 1, "\tsed -e 's,\\#,hash,g'")

	c.Check(mklineCommandEscaped.ShellCommand(), equals, "sed -e 's,\\#,hash,g'")

	// From shells/zsh/Makefile.common, rev. 1.78
	mklineCommandUnescaped := t.NewMkLine("filename", 1, "\t# $ sha1 patches/patch-ac")

	c.Check(mklineCommandUnescaped.ShellCommand(), equals, "# $ sha1 patches/patch-ac")
	t.CheckOutputEmpty() // No warning about parsing the lonely dollar sign.

	mklineVarassignUnescaped := t.NewMkLine("filename", 1, "SED_CMD=\t's,#,hash,'")

	c.Check(mklineVarassignUnescaped.Value(), equals, "'s,")
	t.CheckOutputLines(
		"WARN: filename:1: The # character starts a comment.")
}

func (s *Suite) Test_NewMkLine__varassign_leading_space(c *check.C) {
	t := s.Init(c)

	_ = t.NewMkLine("rubyversion.mk", 427, " _RUBYVER=\t2.15")
	_ = t.NewMkLine("bsd.buildlink3.mk", 132, "   ok:=yes")

	// In mk/buildlink3/bsd.buildlink3.mk, the leading space is really helpful,
	// therefore no warnings for that file.
	t.CheckOutputLines(
		"WARN: rubyversion.mk:427: Makefile lines should not start with space characters.")
}

// Exotic code examples from the pkgsrc infrastructure.
// Hopefully, pkgsrc packages don't need such complicated code.
// Still, pkglint needs to parse them correctly, or it would not
// be able to parse and check the infrastructure files as well.
//
// See Pkgsrc.loadUntypedVars.
func (s *Suite) Test_NewMkLine__infrastructure(c *check.C) {
	t := s.Init(c)

	mklines := t.NewMkLines("infra.mk",
		MkRcsID,
		"         USE_BUILTIN.${_pkg_:S/^-//}:=no",
		".error \"Something went wrong\"",
		".export WRKDIR",
		".export",
		".unexport-env WRKDIR",
		"",
		".ifmake target1",    // Luckily, this is not used in the wild.
		".elifnmake target2", // Neither is this.
		".endif")

	c.Check(mklines.mklines[1].Varcanon(), equals, "USE_BUILTIN.*")
	c.Check(mklines.mklines[2].Directive(), equals, "error")
	c.Check(mklines.mklines[3].Directive(), equals, "export")

	t.CheckOutputLines(
		"WARN: infra.mk:2: Makefile lines should not start with space characters.",
		"ERROR: infra.mk:8: Unknown Makefile line format: \".ifmake target1\".",
		"ERROR: infra.mk:9: Unknown Makefile line format: \".elifnmake target2\".")

	mklines.Check()

	t.CheckOutputLines(
		"WARN: infra.mk:2: USE_BUILTIN.${_pkg_:S/^-//} is defined but not used.",
		"WARN: infra.mk:2: _pkg_ is used but not defined.",
		"ERROR: infra.mk:5: \".export\" requires arguments.",
		"NOTE: infra.mk:2: This variable value should be aligned to column 41.",
		"ERROR: infra.mk:10: Unmatched .endif.")
}

func (s *Suite) Test_MkLine_VariableNeedsQuoting__unknown_rhs(c *check.C) {
	t := s.Init(c)

	mkline := t.NewMkLine("filename", 1, "PKGNAME:= ${UNKNOWN}")
	t.SetupVartypes()

	vuc := &VarUseContext{G.Pkgsrc.VariableType("PKGNAME"), vucTimeParse, vucQuotUnknown, false}
	nq := mkline.VariableNeedsQuoting("UNKNOWN", nil, vuc)

	c.Check(nq, equals, unknown)
}

func (s *Suite) Test_MkLine_VariableNeedsQuoting__append_URL_to_list_of_URLs(c *check.C) {
	t := s.Init(c)

	t.SetupVartypes()
	t.SetupMasterSite("MASTER_SITE_SOURCEFORGE", "http://downloads.sourceforge.net/sourceforge/")
	mkline := t.NewMkLine("Makefile", 95, "MASTER_SITES=\t${HOMEPAGE}")

	vuc := &VarUseContext{G.Pkgsrc.vartypes["MASTER_SITES"], vucTimeRun, vucQuotPlain, false}
	nq := mkline.VariableNeedsQuoting("HOMEPAGE", G.Pkgsrc.vartypes["HOMEPAGE"], vuc)

	c.Check(nq, equals, no)

	MkLineChecker{mkline}.checkVarassign()

	t.CheckOutputEmpty() // Up to version 5.3.6, pkglint warned about a missing :Q here, which was wrong.
}

func (s *Suite) Test_MkLine_VariableNeedsQuoting__append_list_to_list(c *check.C) {
	t := s.Init(c)

	t.SetupVartypes()
	t.SetupMasterSite("MASTER_SITE_SOURCEFORGE", "http://downloads.sourceforge.net/sourceforge/")
	mkline := t.NewMkLine("Makefile", 96, "MASTER_SITES=\t${MASTER_SITE_SOURCEFORGE:=squirrel-sql/}")

	MkLineChecker{mkline}.checkVarassign()

	// Assigning lists to lists is ok.
	t.CheckOutputEmpty()
}

func (s *Suite) Test_MkLine_VariableNeedsQuoting__eval_shell(c *check.C) {
	t := s.Init(c)

	t.SetupVartypes()
	mkline := t.NewMkLine("builtin.mk", 3,
		"USE_BUILTIN.Xfixes!=\t${PKG_ADMIN} pmatch 'pkg-[0-9]*' ${BUILTIN_PKG.Xfixes:Q}")

	MkLineChecker{mkline}.checkVarassign()

	t.CheckOutputLines(
		"WARN: builtin.mk:3: PKG_ADMIN should not be evaluated at load time.",
		"NOTE: builtin.mk:3: The :Q operator isn't necessary for ${BUILTIN_PKG.Xfixes} here.")
}

func (s *Suite) Test_MkLine_VariableNeedsQuoting__command_in_single_quotes(c *check.C) {
	t := s.Init(c)

	t.SetupVartypes()
	mkline := t.NewMkLine("Makefile", 3,
		"SUBST_SED.hpath=\t-e 's|^\\(INSTALL[\t:]*=\\).*|\\1${INSTALL}|'")

	MkLineChecker{mkline}.checkVarassign()

	t.CheckOutputLines(
		"WARN: Makefile:3: Please use ${INSTALL:Q} instead of ${INSTALL} " +
			"and make sure the variable appears outside of any quoting characters.")
}

func (s *Suite) Test_MkLine_VariableNeedsQuoting__command_in_command(c *check.C) {
	t := s.Init(c)

	t.SetupVartypes()
	t.SetupTool("find", "FIND", AtRunTime)
	t.SetupTool("sort", "SORT", AtRunTime)
	G.Pkg = NewPackage(t.File("category/pkgbase"))
	G.Mk = t.NewMkLines("Makefile",
		MkRcsID,
		"GENERATE_PLIST= cd ${DESTDIR}${PREFIX}; ${FIND} * \\( -type f -or -type l \\) | ${SORT};")

	G.Mk.DetermineDefinedVariables()
	MkLineChecker{G.Mk.mklines[1]}.Check()

	t.CheckOutputLines(
		"WARN: Makefile:2: The exitcode of \"${FIND}\" at the left of the | operator is ignored.")
}

func (s *Suite) Test_MkLine_VariableNeedsQuoting__word_as_part_of_word(c *check.C) {
	t := s.Init(c)

	t.SetupVartypes()
	G.Mk = t.NewMkLines("Makefile",
		MkRcsID,
		"EGDIR=\t${EGDIR}/${MACHINE_GNU_PLATFORM}")

	MkLineChecker{G.Mk.mklines[1]}.Check()

	t.CheckOutputEmpty()
}

// As an argument to ${ECHO}, the :Q modifier should be used, but as of
// October 2018, pkglint does not know all shell commands and how they
// handle their arguments. As an argument to xargs(1), the :Q modifier
// would be misplaced, therefore no warning is issued in both these cases.
//
// Based on graphics/circos/Makefile.
func (s *Suite) Test_MkLine_VariableNeedsQuoting__command_as_command_argument(c *check.C) {
	t := s.Init(c)

	t.SetupTool("perl", "PERL5", AtRunTime)
	t.SetupTool("bash", "BASH", AtRunTime)
	t.SetupVartypes()
	mklines := t.NewMkLines("Makefile",
		MkRcsID,
		"\t${RUN} cd ${WRKSRC} && ( ${ECHO} ${PERL5:Q} ; ${ECHO} ) | ${BASH} ./install",
		"\t${RUN} cd ${WRKSRC} && ( ${ECHO} ${PERL5} ; ${ECHO} ) | ${BASH} ./install")

	mklines.Check()

	t.CheckOutputLines(
		"WARN: Makefile:2: The exitcode of the command at the left of the | operator is ignored.",
		"WARN: Makefile:3: The exitcode of the command at the left of the | operator is ignored.")
}

// As seen in mail/mailfront/Makefile.
func (s *Suite) Test_MkLine_VariableNeedsQuoting__URL_as_part_of_word_in_list(c *check.C) {
	t := s.Init(c)

	t.SetupVartypes()
	G.Mk = t.NewMkLines("Makefile",
		MkRcsID,
		"MASTER_SITES=${HOMEPAGE}archive/")

	MkLineChecker{G.Mk.mklines[1]}.Check()

	t.CheckOutputEmpty() // Don't suggest to use ${HOMEPAGE:Q}.
}

// Before November 2018, pkglint did not parse $$(subshell) commands very well.
// As a side effect, it sometimes issued wrong warnings about the :Q modifier.
//
// As seen in www/firefox31/xpi.mk.
func (s *Suite) Test_MkLine_VariableNeedsQuoting__command_in_subshell(c *check.C) {
	t := s.Init(c)

	t.SetupVartypes()
	t.SetupTool("awk", "AWK", AtRunTime)
	t.SetupTool("echo", "ECHO", AtRunTime)
	G.Mk = t.NewMkLines("xpi.mk",
		MkRcsID,
		"\t id=$$(${AWK} '{print}' < ${WRKSRC}/idfile) && echo \"$$id\"",
		"\t id=`${AWK} '{print}' < ${WRKSRC}/idfile` && echo \"$$id\"")

	MkLineChecker{G.Mk.mklines[1]}.Check()
	MkLineChecker{G.Mk.mklines[2]}.Check()

	// Don't suggest to use ${AWK:Q}.
	t.CheckOutputLines(
		"WARN: xpi.mk:2: Invoking subshells via $(...) is not portable enough.")
}

// LDFLAGS (and even more so CPPFLAGS and CFLAGS) may contain special
// shell characters like quotes or backslashes. Therefore, quoting them
// correctly is trickier than with other variables.
func (s *Suite) Test_MkLine_VariableNeedsQuoting__LDFLAGS_in_single_quotes(c *check.C) {
	t := s.Init(c)

	t.SetupVartypes()
	G.Mk = t.NewMkLines("x11/mlterm/Makefile",
		MkRcsID,
		"SUBST_SED.link=-e 's|(LIBTOOL_LINK).*(LIBS)|& ${LDFLAGS:M*:Q}|g'",
		"SUBST_SED.link=-e 's|(LIBTOOL_LINK).*(LIBS)|& '${LDFLAGS:M*:Q}'|g'")

	MkLineChecker{G.Mk.mklines[1]}.Check()
	MkLineChecker{G.Mk.mklines[2]}.Check()

	t.CheckOutputLines(
		"WARN: x11/mlterm/Makefile:2: Please move ${LDFLAGS:M*:Q} outside of any quoting characters.")
}

// No quoting is necessary when lists of options are appended to each other.
// PKG_OPTIONS are declared as "lkShell" although they are processed
// using make's .for loop, which splits them at whitespace and usually
// requires the variable to be declared as "lkSpace".
// In this case it doesn't matter though since each option is an identifier,
// and these do not pose any quoting or escaping problems.
func (s *Suite) Test_MkLine_VariableNeedsQuoting__package_options(c *check.C) {
	t := s.Init(c)

	t.SetupVartypes()
	G.Mk = t.NewMkLines("Makefile",
		MkRcsID,
		"PKG_SUGGESTED_OPTIONS+=\t${PKG_DEFAULT_OPTIONS:Mcdecimal} ${PKG_OPTIONS.py-trytond:Mcdecimal}")

	MkLineChecker{G.Mk.mklines[1]}.Check()

	// No warning about a missing :Q modifier.
	t.CheckOutputEmpty()
}

func (s *Suite) Test_MkLine_VariableNeedsQuoting__tool_in_quotes_in_subshell_in_shellwords(c *check.C) {
	t := s.Init(c)

	t.SetupTool("echo", "ECHO", AtRunTime)
	t.SetupTool("sh", "SH", AtRunTime)
	t.SetupVartypes()
	G.Mk = t.NewMkLines("x11/labltk/Makefile",
		MkRcsID,
		"CONFIGURE_ARGS+=\t-tklibs \"`${SH} -c '${ECHO} $$TK_LD_FLAGS'`\"")

	MkLineChecker{G.Mk.mklines[1]}.Check()

	// Don't suggest ${ECHO:Q} here.
	t.CheckOutputEmpty()
}

func (s *Suite) Test_MkLine_VariableNeedsQuoting__LDADD_in_BUILDLINK_TRANSFORM(c *check.C) {
	t := s.Init(c)

	t.SetupVartypes()
	G.Mk = t.NewMkLines("x11/qt5-qtbase/Makefile.common",
		"BUILDLINK_TRANSFORM+=opt:-ldl:${BUILDLINK_LDADD.dl:M*}")

	MkLineChecker{G.Mk.mklines[0]}.Check()

	// Note: The :M* modifier is not necessary, since this is not a GNU Configure package.
	t.CheckOutputLines(
		"WARN: x11/qt5-qtbase/Makefile.common:1: Please use ${BUILDLINK_LDADD.dl:Q} instead of ${BUILDLINK_LDADD.dl:M*}.")
}

func (s *Suite) Test_MkLine_VariableNeedsQuoting__command_in_message(c *check.C) {
	t := s.Init(c)

	t.SetupVartypes()
	G.Mk = t.NewMkLines("benchmarks/iozone/Makefile",
		"SUBST_MESSAGE.crlf=\tStripping EOL CR in ${REPLACE_PERL}")

	MkLineChecker{G.Mk.mklines[0]}.Check()

	// Don't suggest ${REPLACE_PERL:Q}.
	t.CheckOutputEmpty()
}

func (s *Suite) Test_MkLine_VariableNeedsQuoting__guessed_list_variable_in_quotes(c *check.C) {
	t := s.Init(c)

	t.SetupVartypes()
	G.Mk = t.NewMkLines("audio/jack-rack/Makefile",
		MkRcsID,
		"LADSPA_PLUGIN_PATH=\t${PREFIX}/lib/ladspa",
		"CPPFLAGS+=\t\t-DLADSPA_PATH=\"\\\"${LADSPA_PLUGIN_PATH}\\\"\"")

	G.Mk.Check()

	t.CheckOutputLines(
		"WARN: audio/jack-rack/Makefile:3: The variable LADSPA_PLUGIN_PATH should be quoted as part of a shell word.")
}

func (s *Suite) Test_MkLine_VariableNeedsQuoting__list_in_list(c *check.C) {
	t := s.Init(c)

	t.SetupVartypes()
	G.Mk = t.NewMkLines("x11/eterm/Makefile",
		MkRcsID,
		"DISTFILES=\t${DEFAULT_DISTFILES} ${PIXMAP_FILES}")

	G.Mk.Check()

	// Don't warn about missing :Q modifiers.
	t.CheckOutputLines(
		"WARN: x11/eterm/Makefile:2: PIXMAP_FILES is used but not defined.")
}

func (s *Suite) Test_MkLine_VariableNeedsQuoting__PKGNAME_and_URL_list_in_URL_list(c *check.C) {
	t := s.Init(c)

	t.SetupMasterSite("MASTER_SITE_GNOME", "http://ftp.gnome.org/")
	t.SetupVartypes()
	G.Mk = t.NewMkLines("x11/gtk3/Makefile",
		MkRcsID,
		"MASTER_SITES=\tftp://ftp.gtk.org/${PKGNAME}/ ${MASTER_SITE_GNOME:=subdir/}")

	MkLineChecker{G.Mk.mklines[1]}.checkVarassignVaruse()

	t.CheckOutputEmpty() // Don't warn about missing :Q modifiers.
}

func (s *Suite) Test_MkLine_VariableNeedsQuoting__tool_in_CONFIGURE_ENV(c *check.C) {
	t := s.Init(c)

	t.SetupVartypes()
	t.SetupTool("tar", "TAR", AtRunTime)
	mklines := t.NewMkLines("Makefile",
		MkRcsID,
		"",
		"CONFIGURE_ENV+=\tSYS_TAR_COMMAND_PATH=${TOOLS_TAR:Q}")

	MkLineChecker{mklines.mklines[2]}.checkVarassignVaruse()

	// The TOOLS_* variables only contain the path to the tool,
	// without any additional arguments that might be necessary
	// for invoking the tool properly (e.g. touch -t).
	// Therefore, no quoting is necessary.
	t.CheckOutputLines(
		"NOTE: Makefile:3: The :Q operator isn't necessary for ${TOOLS_TAR} here.")
}

func (s *Suite) Test_MkLine_VariableNeedsQuoting__backticks(c *check.C) {
	t := s.Init(c)

	t.SetupVartypes()
	t.SetupTool("cat", "CAT", AtRunTime)
	mklines := t.NewMkLines("Makefile",
		MkRcsID,
		"",
		"COMPILE_CMD=\tcc `${CAT} ${WRKDIR}/compileflags`",
		"COMMENT_CMD=\techo `echo ${COMMENT}`")

	MkLineChecker{mklines.mklines[2]}.checkVarassignVaruse()
	MkLineChecker{mklines.mklines[3]}.checkVarassignVaruse()

	// Both CAT and WRKDIR are safe from quoting, therefore no warnings.
	// But COMMENT may contain arbitrary characters and therefore must
	// only appear completely unquoted. There is no practical way of
	// using it inside backticks, and luckily there is no need for it.
	t.CheckOutputLines(
		"WARN: Makefile:4: COMMENT may not be used in any file; it is a write-only variable.",
		// TODO: Better suggest that COMMENT should not be used inside backticks or other quotes.
		"WARN: Makefile:4: The variable COMMENT should be quoted as part of a shell word.")
}

// For some well-known directory variables like WRKDIR, PREFIX, LOCALBASE,
// the :Q modifier can be safely removed since pkgsrc will never support
// having special characters in these directory names.
// For guessed variable types be cautious and don't autofix them.
func (s *Suite) Test_MkLine_VariableNeedsQuoting__only_remove_known(c *check.C) {
	t := s.Init(c)

	t.SetupCommandLine("-Wall", "--autofix")
	t.SetupVartypes()

	mklines := t.SetupFileMkLines("Makefile",
		MkRcsID,
		"",
		"demo: .PHONY",
		"\t${ECHO} ${WRKSRC:Q}",
		"\t${ECHO} ${FOODIR:Q}")

	mklines.Check()

	t.CheckOutputLines(
		"AUTOFIX: ~/Makefile:4: Replacing \"${WRKSRC:Q}\" with \"${WRKSRC}\".")
	t.CheckFileLines("Makefile",
		MkRcsID,
		"",
		"demo: .PHONY",
		"\t${ECHO} ${WRKSRC}",
		"\t${ECHO} ${FOODIR:Q}")
}

// TODO: COMPILER_RPATH_FLAG and LINKER_RPATH_FLAG have different types
// defined in vardefs.go; examine why.
func (s *Suite) Test_MkLine_VariableNeedsQuoting__shellword_part(c *check.C) {
	t := s.Init(c)

	t.SetupCommandLine("-Wall,no-space")
	t.SetupVartypes()

	mklines := t.SetupFileMkLines("Makefile",
		MkRcsID,
		"",
		"SUBST_CLASSES+=    class",
		"SUBST_STAGE.class= pre-configure",
		"SUBST_FILES.class= files",
		"SUBST_SED.class=-e s:@LINKER_RPATH_FLAG@:${LINKER_RPATH_FLAG}:g")

	mklines.Check()

	t.CheckOutputLines(
		"NOTE: ~/Makefile:6: The substitution command \"s:@LINKER_RPATH_FLAG@:${LINKER_RPATH_FLAG}:g\" " +
			"can be replaced with \"SUBST_VARS.class+= LINKER_RPATH_FLAG\".")
}

// Tools, when used in a shell command, must not be quoted.
func (s *Suite) Test_MkLine_VariableNeedsQuoting__tool_in_shell_command(c *check.C) {
	t := s.Init(c)

	t.SetupCommandLine("-Wall,no-space")
	t.SetupVartypes()
	t.SetupTool("bash", "BASH", AtRunTime)

	mklines := t.SetupFileMkLines("Makefile",
		MkRcsID,
		"",
		"CONFIG_SHELL= ${BASH}")

	mklines.Check()

	t.CheckOutputEmpty()
}

// As of October 2018, these examples from real pkgsrc end up in the
// final "unknown" case.
func (s *Suite) Test_MkLine_VariableNeedsQuoting__uncovered_cases(c *check.C) {
	t := s.Init(c)

	t.SetupCommandLine("-Wall,no-space")
	t.SetupVartypes()

	mklines := t.SetupFileMkLines("Makefile",
		MkRcsID,
		"",
		"GO_SRCPATH=             ${HOMEPAGE:S,https://,,}",
		"LINKER_RPATH_FLAG:=     ${LINKER_RPATH_FLAG:S/-rpath/& /}",
		"HOMEPAGE=               http://godoc.org/${GO_SRCPATH}",
		"PATH:=                  ${PREFIX}/cross/bin:${PATH}",
		"NO_SRC_ON_FTP=          ${RESTRICTED}")

	mklines.Check()

	t.CheckOutputLines(
		// TODO: Explain why the variable may not be set, by listing the current rules.
		"WARN: ~/Makefile:4: The variable LINKER_RPATH_FLAG may not be set by any package.",
		"WARN: ~/Makefile:4: Please use ${LINKER_RPATH_FLAG:S/-rpath/& /:Q} instead of ${LINKER_RPATH_FLAG:S/-rpath/& /}.",
		"WARN: ~/Makefile:4: LINKER_RPATH_FLAG should not be evaluated at load time.",
		"WARN: ~/Makefile:6: The variable PATH may not be set by any package.",
		"WARN: ~/Makefile:6: PREFIX should not be evaluated at load time.",
		"WARN: ~/Makefile:6: PATH should not be evaluated at load time.")
}

func (s *Suite) Test_MkLine__shell_varuse_in_backt_dquot(c *check.C) {
	t := s.Init(c)

	t.SetupVartypes()
	t.SetupTool("grep", "GREP", AtRunTime)
	mklines := t.NewMkLines("x11/motif/Makefile",
		MkRcsID,
		"post-patch:",
		"\tfiles=`${GREP} -l \".fB$${name}.fP(3)\" *.3`")

	mklines.Check()

	// Just ensure that there are no parse errors.
	t.CheckOutputEmpty()
}

// PR 51696, security/py-pbkdf2/Makefile, r1.2
func (s *Suite) Test_MkLine__comment_in_comment(c *check.C) {
	t := s.Init(c)

	t.SetupVartypes()
	mklines := t.NewMkLines("Makefile",
		MkRcsID,
		"COMMENT=\tPKCS#5 v2.0 PBKDF2 Module")

	mklines.Check()

	t.CheckOutputLines(
		"WARN: Makefile:2: The # character starts a comment.")
}

// Ensures that the conditional variables of a line can be set even
// after initializing the MkLine.
//
// If this test should fail, it is probably because mkLineDirective
// is not a pointer type anymore.
//
// See https://github.com/golang/go/issues/28045.
func (s *Suite) Test_MkLine_ConditionalVars(c *check.C) {
	t := s.Init(c)

	mkline := t.NewMkLine("Makefile", 45, ".include \"../../category/package/buildlink3.mk\"")

	c.Check(mkline.ConditionalVars(), check.HasLen, 0)

	mkline.SetConditionalVars([]string{"OPSYS"})

	c.Check(mkline.ConditionalVars(), deepEquals, []string{"OPSYS"})
}

func (s *Suite) Test_MkLine_ValueSplit(c *check.C) {
	t := s.Init(c)

	test := func(value string, expected ...string) {
		mkline := t.NewMkLine("Makefile", 1, "PATH=\t"+value)
		split := mkline.ValueSplit(value, ":")
		c.Check(split, deepEquals, expected)
	}

	test("#empty",
		[]string(nil)...)

	test("/bin",
		"/bin")

	test("/bin:/sbin",
		"/bin",
		"/sbin")

	test("${DESTDIR}/bin:/bin/${SUBDIR}",
		"${DESTDIR}/bin",
		"/bin/${SUBDIR}")

	test("/bin:${DESTDIR}${PREFIX}:${DESTDIR:S,/,\\:,:S,:,:,}/sbin",
		"/bin",
		"${DESTDIR}${PREFIX}",
		"${DESTDIR:S,/,\\:,:S,:,:,}/sbin")

	test("${VAR:Udefault}::${VAR2}two:words",
		"${VAR:Udefault}",
		"",
		"${VAR2}two",
		"words")
}

func (s *Suite) Test_MkLine_ValueFields(c *check.C) {
	t := s.Init(c)

	test := func(value string, expected ...string) {
		mkline := t.NewMkLine("Makefile", 1, "VAR=\t"+value)
		split := mkline.ValueFields(value)
		c.Check(split, deepEquals, expected)
	}

	test("one   two\t\t${THREE:Uthree:Nsome \tspaces}",
		"one",
		"two",
		"${THREE:Uthree:Nsome \tspaces}")

	test("${VAR:Udefault value} ${VAR2}two words",
		"${VAR:Udefault value}",
		"${VAR2}two",
		"words")
}

// Before 2018-11-26, this test panicked.
func (s *Suite) Test_MkLine_ValueFields__adjacent_vars(c *check.C) {
	t := s.Init(c)

	test := func(value string, expected ...string) {
		mkline := t.NewMkLine("Makefile", 1, "")
		split := mkline.ValueFields(value)
		c.Check(split, deepEquals, expected)
	}

	test("\t; ${RM} ${WRKSRC}",
		";",
		"${RM}",
		"${WRKSRC}")
}

func (s *Suite) Test_MkLine_ValueTokens(c *check.C) {
	t := s.Init(c)

	testTokens := func(value string, expected ...*MkToken) {
		mkline := t.NewMkLine("Makefile", 1, "PATH=\t"+value)
		split := mkline.ValueTokens()
		c.Check(split, deepEquals, expected)
	}

	testTokens("#empty",
		[]*MkToken(nil)...)

	testTokens("value",
		&MkToken{"value", nil})

	testTokens("value ${VAR} rest",
		&MkToken{"value ", nil},
		&MkToken{"${VAR}", NewMkVarUse("VAR")},
		&MkToken{" rest", nil})

	testTokens("value ${UNFINISHED",
		&MkToken{"value ", nil})
}

func (s *Suite) Test_MkLine_ValueTokens__caching(c *check.C) {
	t := s.Init(c)

	mkline := t.NewMkLine("Makefile", 1, "PATH=\tvalue ${UNFINISHED")
	split := mkline.ValueTokens()

	c.Check(split, deepEquals, []*MkToken{{"value ", nil}})

	split2 := mkline.ValueTokens() // This time the slice is taken from the cache.

	// In Go, it's not possible to compare slices for reference equality.
	c.Check(split2, deepEquals, split)
}

func (s *Suite) Test_MkLine_ResolveVarsInRelativePath(c *check.C) {
	t := s.Init(c)

	t.CreateFileLines("lang/lua53/Makefile")
	t.CreateFileLines("lang/php72/Makefile")
	t.CreateFileLines("emulators/suse100_base/Makefile")
	t.CreateFileLines("lang/python36/Makefile")
	mklines := t.SetupFileMkLines("Makefile",
		MkRcsID)
	mkline := mklines.mklines[0]

	checkResolve := func(before string, after string) {
		c.Check(mkline.ResolveVarsInRelativePath(before, false), equals, after)
	}

	checkResolve("", ".")
	checkResolve("${LUA_PKGSRCDIR}", "../../lang/lua53")
	checkResolve("${PHPPKGSRCDIR}", "../../lang/php72")
	checkResolve("${SUSE_DIR_PREFIX}", "suse100")
	checkResolve("${PYPKGSRCDIR}", "../../lang/python36")
	checkResolve("${PYPACKAGE}", "python36")
	checkResolve("${FILESDIR}", "${FILESDIR}")
	checkResolve("${PKGDIR}", "${PKGDIR}")

	G.Pkg = NewPackage(t.File("category/package"))

	checkResolve("${FILESDIR}", "files")
	checkResolve("${PKGDIR}", ".")
}

func (s *Suite) Test_MkLine_ResolveVarsInRelativePath__directory_depth(c *check.C) {
	t := s.Init(c)

	t.SetupVartypes()
	mklines := t.SetupFileMkLines("multimedia/totem/bla.mk",
		MkRcsID,
		"BUILDLINK_PKGSRCDIR.totem?=\t../../multimedia/totem")

	mklines.Check()

	t.CheckOutputLines(
		"ERROR: ~/multimedia/totem/bla.mk:2: There is no package in \"multimedia/totem\".")
}

func (s *Suite) Test_MatchVarassign(c *check.C) {
	s.Init(c)

	checkVarassign := func(text string, commented bool, varname, spaceAfterVarname, op, align, value, spaceAfterValue, comment string) {
		type VarAssign struct {
			commented                  bool
			varname, spaceAfterVarname string
			op, align                  string
			value, spaceAfterValue     string
			comment                    string
		}
		expected := VarAssign{commented, varname, spaceAfterVarname, op, align, value, spaceAfterValue, comment}
		am, acommented, avarname, aspaceAfterVarname, aop, aalign, avalue, aspaceAfterValue, acomment := MatchVarassign(text)
		if !am {
			c.Errorf("Text %q doesn't match variable assignment", text)
			return
		}
		actual := VarAssign{acommented, avarname, aspaceAfterVarname, aop, aalign, avalue, aspaceAfterValue, acomment}
		c.Check(actual, equals, expected)
	}
	checkNotVarassign := func(text string) {
		m, _, _, _, _, _, _, _, _ := MatchVarassign(text)
		if m {
			c.Errorf("Text %q matches variable assignment but shouldn't.", text)
		}
	}

	checkVarassign("C++=c11", false, "C+", "", "+=", "C++=", "c11", "", "")
	checkVarassign("V=v", false, "V", "", "=", "V=", "v", "", "")
	checkVarassign("VAR=#comment", false, "VAR", "", "=", "VAR=", "", "", "#comment")
	checkVarassign("VAR=\\#comment", false, "VAR", "", "=", "VAR=", "#comment", "", "")
	checkVarassign("VAR=\\\\\\##comment", false, "VAR", "", "=", "VAR=", "\\\\#", "", "#comment")
	checkVarassign("VAR=\\", false, "VAR", "", "=", "VAR=", "\\", "", "")
	checkVarassign("VAR += value", false, "VAR", " ", "+=", "VAR += ", "value", "", "")
	checkVarassign(" VAR=value", false, "VAR", "", "=", " VAR=", "value", "", "")
	checkVarassign("VAR=value #comment", false, "VAR", "", "=", "VAR=", "value", " ", "#comment")

	checkNotVarassign("\tVAR=value")
	checkNotVarassign("?=value")
	checkNotVarassign("<=value")
	checkNotVarassign("#")
	checkNotVarassign("VAR.$$=value")

	// A commented variable assignment must start immediately after the comment character.
	// There must be no additional whitespace before the variable name.
	checkVarassign("#VAR=value", true, "VAR", "", "=", "#VAR=", "value", "", "")

	// A single space is typically used for writing documentation, not for commenting out code.
	// Therefore this line doesn't count as commented variable assignment.
	checkNotVarassign("# VAR=value")
}

func (s *Suite) Test_NewMkOperator(c *check.C) {
	c.Check(NewMkOperator(":="), equals, opAssignEval)
	c.Check(NewMkOperator("="), equals, opAssign)

	c.Check(func() { NewMkOperator("???") }, check.Panics, "Invalid operator: ???")
}

func (s *Suite) Test_Indentation(c *check.C) {
	t := s.Init(c)

	ind := NewIndentation()

	mkline := t.NewMkLine("dummy.mk", 5, ".if 0")

	c.Check(ind.Depth("if"), equals, 0)
	c.Check(ind.DependsOn("VARNAME"), equals, false)

	ind.Push(mkline, 2, "")

	c.Check(ind.Depth("if"), equals, 0) // Because "if" is handled in MkLines.TrackBefore.
	c.Check(ind.Depth("endfor"), equals, 0)

	ind.AddVar("LEVEL1.VAR1")

	c.Check(ind.Varnames(), deepEquals, []string{"LEVEL1.VAR1"})

	ind.AddVar("LEVEL1.VAR2")

	c.Check(ind.Varnames(), deepEquals, []string{"LEVEL1.VAR1", "LEVEL1.VAR2"})
	c.Check(ind.DependsOn("LEVEL1.VAR1"), equals, true)
	c.Check(ind.DependsOn("OTHER_VAR"), equals, false)

	ind.Push(mkline, 2, "")

	ind.AddVar("LEVEL2.VAR")

	c.Check(ind.Varnames(), deepEquals, []string{"LEVEL1.VAR1", "LEVEL1.VAR2", "LEVEL2.VAR"})
	c.Check(ind.String(), equals, "[2 (LEVEL1.VAR1 LEVEL1.VAR2) 2 (LEVEL2.VAR)]")

	ind.Pop()

	c.Check(ind.Varnames(), deepEquals, []string{"LEVEL1.VAR1", "LEVEL1.VAR2"})
	c.Check(ind.IsConditional(), equals, true)

	ind.Pop()

	c.Check(ind.Varnames(), check.HasLen, 0)
	c.Check(ind.IsConditional(), equals, false)
	c.Check(ind.String(), equals, "[]")
}

func (s *Suite) Test_Indentation_RememberUsedVariables(c *check.C) {
	t := s.Init(c)

	mkline := t.NewMkLine("Makefile", 123, ".if ${PKGREVISION} > 0")
	ind := NewIndentation()

	ind.RememberUsedVariables(mkline.Cond())

	t.CheckOutputEmpty()
	c.Check(ind.Varnames(), deepEquals, []string{"PKGREVISION"})
}

func (s *Suite) Test_MkLine_DetermineUsedVariables(c *check.C) {
	t := s.Init(c)

	mklines := t.NewMkLines("Makefile",
		MkRcsID,
		"VAR=\t${VALUE} # ${varassign.comment}",
		".if ${OPSYS:M${endianness}} == ${Hello:L} # ${if.comment}",
		".for var in one ${two} three # ${for.comment}",
		"# ${empty.comment}",
		"${TARGETS}: ${SOURCES} # ${dependency.comment}",
		".include \"${OTHER_FILE}\"",
		"",
		"\t${VAR.${param}}",
		"\t${VAR}and${VAR2}",
		"\t${VAR:M${pattern}}",
		"\t$(ROUND_PARENTHESES)",
		"\t$$shellvar",
		"\t$< $@ $x")

	var varnames []string
	for _, mkline := range mklines.mklines {
		varnames = append(varnames, mkline.DetermineUsedVariables()...)
	}

	c.Check(varnames, deepEquals, []string{
		"VALUE",
		"OPSYS",
		"endianness",
		// "Hello" is not a variable name, the :L modifier makes it an expression.
		"two",
		"TARGETS",
		"SOURCES",
		"OTHER_FILE",

		"VAR.${param}",
		"param",
		"VAR",
		"VAR2",
		"VAR",
		"pattern",
		"ROUND_PARENTHESES",
		// Shell variables are ignored here.
		"<",
		"@",
		"x"})
}

func (s *Suite) Test_matchMkDirective(c *check.C) {

	test := func(input, expectedIndent, expectedDirective, expectedArgs, expectedComment string) {
		m, indent, directive, args, comment := matchMkDirective(input)
		c.Check(
			[]interface{}{m, indent, directive, args, comment},
			deepEquals,
			[]interface{}{true, expectedIndent, expectedDirective, expectedArgs, expectedComment})
	}

	test(".if ${VAR} == value", "", "if", "${VAR} == value", "")
	test(".\tendif # comment", "\t", "endif", "", "comment")
	test(".if ${VAR} == \"#\"", "", "if", "${VAR} == \"", "\"")
	test(".if ${VAR:[#]}", "", "if", "${VAR:[#]}", "")
}
