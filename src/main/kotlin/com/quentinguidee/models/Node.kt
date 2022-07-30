package com.quentinguidee.models

import kotlinx.serialization.KSerializer
import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable
import kotlinx.serialization.descriptors.PrimitiveKind
import kotlinx.serialization.descriptors.PrimitiveSerialDescriptor
import kotlinx.serialization.encoding.Decoder
import kotlinx.serialization.encoding.Encoder
import org.jetbrains.exposed.dao.id.UUIDTable
import org.jetbrains.exposed.sql.javatime.datetime
import java.time.LocalDateTime

object Nodes : UUIDTable() {
    val parent = reference("parent_uuid", Nodes).nullable()
    val bucket = reference("bucket_uuid", Buckets)
    val name = varchar("name", 255)
    val type = varchar("type", 255)
    val mime = varchar("mime", 255).nullable()
    val size = integer("size").default(0)

    val createdAt = datetime("created_at").default(LocalDateTime.now())
    val updatedAt = datetime("updated_at").nullable()
    val deletedAt = datetime("deleted_at").nullable()
}

@Serializable
data class Node(
    val uuid: String,
    @SerialName("parent_uuid")
    val parentUUID: String?,
    @SerialName("bucket_uuid")
    val bucketUUID: String,
    val name: String,
    val type: String,
    val mime: String?,
    val size: Int,

    @Serializable(DateSerializer::class) val createdAt: LocalDateTime,
    @Serializable(DateSerializer::class) val updatedAt: LocalDateTime? = null,
    @Serializable(DateSerializer::class) val deletedAt: LocalDateTime? = null,
)

object DateSerializer : KSerializer<LocalDateTime> {
    override val descriptor = PrimitiveSerialDescriptor("LocalDateTime", PrimitiveKind.LONG)
    override fun serialize(encoder: Encoder, value: LocalDateTime) = encoder.encodeString(value.toString())
    override fun deserialize(decoder: Decoder): LocalDateTime = LocalDateTime.parse(decoder.decodeString())
}
