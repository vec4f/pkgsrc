# $NetBSD$

BUILDLINK_TREE+=	compiler-rt6

.if !defined(COMPILER_RT6_BUILDLINK3_MK)
COMPILER_RT6_BUILDLINK3_MK:=

BUILDLINK_API_DEPENDS.compiler-rt6+=	compiler-rt6>=6.0.0alpha1
BUILDLINK_PKGSRCDIR.compiler-rt6?=	../../lang/compiler-rt6

.include "../../lang/llvm6/buildlink3.mk"
.endif	# COMPILER_RT6_BUILDLINK3_MK

BUILDLINK_TREE+=	-compiler-rt6
