# $NetBSD$

BUILDLINK_TREE+=	libcxxabi6

.if !defined(LIBCXXABI6_BUILDLINK3_MK)
LIBCXXABI6_BUILDLINK3_MK:=

BUILDLINK_API_DEPENDS.libcxxabi6+=	libcxxabi6>=6.0.0alpha1
BUILDLINK_PKGSRCDIR.libcxxabi6?=	../../lang/libcxxabi6

.include "../../lang/llvm6/buildlink3.mk"
.endif	# LIBCXXABI6_BUILDLINK3_MK

BUILDLINK_TREE+=	-libcxxabi6
