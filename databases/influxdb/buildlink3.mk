# $NetBSD$

BUILDLINK_TREE+=	influxdb

.if !defined(INFLUXDB_BUILDLINK3_MK)
INFLUXDB_BUILDLINK3_MK:=

BUILDLINK_API_DEPENDS.influxdb+=	influxdb>=0
BUILDLINK_PKGSRCDIR.influxdb?=		../../databases/influxdb
.endif  # INFLUXDB_BUILDLINK3_MK

BUILDLINK_TREE+=	-influxdb
