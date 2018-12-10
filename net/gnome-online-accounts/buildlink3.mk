# $NetBSD: buildlink3.mk,v 1.7 2018/12/09 18:52:05 adam Exp $

BUILDLINK_TREE+=	gnome-online-accounts

.if !defined(GNOME_ONLINE_ACCOUNTS_BUILDLINK3_MK)
GNOME_ONLINE_ACCOUNTS_BUILDLINK3_MK:=

BUILDLINK_API_DEPENDS.gnome-online-accounts+=	gnome-online-accounts>=3.29.1
BUILDLINK_ABI_DEPENDS.gnome-online-accounts?=	gnome-online-accounts>=3.28.0nb7
BUILDLINK_PKGSRCDIR.gnome-online-accounts?=	../../net/gnome-online-accounts

.include "../../security/libgcrypt/buildlink3.mk"
.include "../../security/libsecret/buildlink3.mk"
.include "../../textproc/json-glib/buildlink3.mk"
.include "../../www/librest07/buildlink3.mk"
.include "../../www/webkit-gtk/buildlink3.mk"
.endif	# GNOME_ONLINE_ACCOUNTS_BUILDLINK3_MK

BUILDLINK_TREE+=	-gnome-online-accounts
