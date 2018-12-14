# $NetBSD: buildlink3.mk,v 1.11 2018/12/13 19:51:31 adam Exp $

BUILDLINK_TREE+=	ktp-common-internals

.if !defined(KTP_COMMON_INTERNALS_BUILDLINK3_MK)
KTP_COMMON_INTERNALS_BUILDLINK3_MK:=

BUILDLINK_API_DEPENDS.ktp-common-internals+=	ktp-common-internals>=0.9.0
BUILDLINK_ABI_DEPENDS.ktp-common-internals?=	ktp-common-internals>=0.9.0nb11
BUILDLINK_PKGSRCDIR.ktp-common-internals?=	../../chat/ktp-common-internals

.include "../../chat/telepathy-qt/buildlink3.mk"
.include "../../x11/kdelibs4/buildlink3.mk"
.endif	# KTP_COMMON_INTERNALS_BUILDLINK3_MK

BUILDLINK_TREE+=	-ktp-common-internals
