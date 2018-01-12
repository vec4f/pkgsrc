# $NetBSD$

BUILDLINK_TREE+=	libcxx6

.if !defined(LIBCXX6_BUILDLINK3_MK)
LIBCXX6_BUILDLINK3_MK:=

BUILDLINK_API_DEPENDS.libcxx6+=	libcxx6>=6.0.0alpha1
BUILDLINK_PKGSRCDIR.libcxx6?=	../../lang/libcxx6

.include "../../lang/llvm6/buildlink3.mk"
.endif	# LIBCXX6_BUILDLINK3_MK

BUILDLINK_TREE+=	-libcxx6
