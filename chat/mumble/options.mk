# $NetBSD: options.mk,v 1.1 2018/12/22 18:17:39 nia Exp $

PKG_OPTIONS_VAR=		PKG_OPTIONS.mumble
PKG_SUPPORTED_OPTIONS=		alsa dbus pulseaudio portaudio
PKG_SUGGESTED_OPTIONS=		dbus
PKG_SUGGESTED_OPTIONS.Linux=	alsa

.include "../../mk/bsd.options.mk"

.if !empty(PKG_OPTIONS:Malsa)
CONFIG_OPTIONS+=	CONFIG+=alsa
.include "../../audio/alsa-lib/buildlink3.mk"
.else
CONFIG_OPTIONS+=	CONFIG+=no-alsa
.endif

.if !empty(PKG_OPTIONS:Mdbus)
.include "../../x11/qt4-qdbus/buildlink3.mk"
.else
CONFIG_OPTIONS+=	CONFIG+=no-dbus
.endif

.if !empty(PKG_OPTIONS:Mpulseaudio)
.include "../../audio/pulseaudio/buildlink3.mk"
.else
CONFIG_OPTIONS+=	CONFIG+=no-pulseaudio
.endif

.if !empty(PKG_OPTIONS:Mportaudio)
.include "../../audio/portaudio/buildlink3.mk"
.else
CONFIG_OPTIONS+=	CONFIG+=no-portaudio
.endif
