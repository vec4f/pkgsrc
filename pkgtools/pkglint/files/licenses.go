package main

import "netbsd.org/pkglint/licenses"

type LicenseChecker struct {
	MkLine MkLine
}

func (lc *LicenseChecker) Check(value string, op MkOperator) {
	expanded := resolveVariableRefs(value) // For ${PERL5_LICENSE}
	cond := licenses.Parse(ifelseStr(op == opAssignAppend, "append-placeholder ", "") + expanded)

	if cond == nil {
		if op == opAssign {
			lc.MkLine.Errorf("Parse error for license condition %q.", value)
		} else {
			lc.MkLine.Errorf("Parse error for appended license condition %q.", value)
		}
		return
	}

	cond.Walk(lc.checkNode)
}

func (lc *LicenseChecker) checkName(license string) {
	licenseFile := ""
	if G.Pkg != nil {
		if mkline := G.Pkg.vars.FirstDefinition("LICENSE_FILE"); mkline != nil {
			licenseFile = G.Pkg.File(mkline.ResolveVarsInRelativePath(mkline.Value(), false))
		}
	}
	if licenseFile == "" {
		licenseFile = G.Pkgsrc.File("licenses/" + license)
		if G.Pkgsrc.UsedLicenses != nil {
			G.Pkgsrc.UsedLicenses[license] = true
		}
	}

	if !fileExists(licenseFile) {
		lc.MkLine.Warnf("License file %s does not exist.", cleanpath(licenseFile))
	}

	switch license {
	case "fee-based-commercial-use",
		"no-commercial-use",
		"no-profit",
		"no-redistribution",
		"shareware":
		lc.MkLine.Errorf("License %q must not be used.", license)
		G.Explain(
			"Instead of using these deprecated licenses, extract the actual",
			"license from the package into the pkgsrc/licenses/ directory",
			"and define LICENSE to that filename.",
			"",
			seeGuide("Handling licenses", "handling-licenses"))
	}
}

func (lc *LicenseChecker) checkNode(cond *licenses.Condition) {
	if name := cond.Name; name != "" && name != "append-placeholder" {
		lc.checkName(name)
		return
	}

	if cond.And && cond.Or {
		lc.MkLine.Errorf("AND and OR operators in license conditions can only be combined using parentheses.")
		G.Explain(
			"Examples for valid license conditions are:",
			"",
			"\tlicense1 AND license2 AND (license3 OR license4)",
			"\t(((license1 OR license2) AND (license3 OR license4)))")
	}
}
