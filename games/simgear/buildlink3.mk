# $NetBSD: buildlink3.mk,v 1.20 2018/11/23 08:06:28 ryoon Exp $

BUILDLINK_TREE+=	simgear

.if !defined(SIMGEAR_BUILDLINK3_MK)
SIMGEAR_BUILDLINK3_MK:=

BUILDLINK_API_DEPENDS.simgear+=	simgear>=2018.2.2
BUILDLINK_ABI_DEPENDS.simgear?=	simgear>=2018.2.2nb6
BUILDLINK_PKGSRCDIR.simgear?=	../../games/simgear

.include "../../audio/openal-soft/buildlink3.mk"
.include "../../devel/boost-headers/buildlink3.mk"
.include "../../devel/zlib/buildlink3.mk"
.include "../../graphics/MesaLib/buildlink3.mk"
.include "../../graphics/osg/buildlink3.mk"
.include "../../www/curl/buildlink3.mk"
.include "../../mk/pthread.buildlink3.mk"
.endif	# SIMGEAR_BUILDLINK3_MK

BUILDLINK_TREE+=	-simgear
