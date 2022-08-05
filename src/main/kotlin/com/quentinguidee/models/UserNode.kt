package com.quentinguidee.models

import com.quentinguidee.utils.DateSerializer
import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable
import org.jetbrains.exposed.sql.Table
import org.jetbrains.exposed.sql.javatime.datetime
import java.time.LocalDateTime

object UsersNodes : Table("users_nodes") {
    val node = reference("node_uuid", Nodes)
    val user = reference("user_id", Users)
    val seenAt = datetime("seen_at").nullable()
    val editedAt = datetime("edited_at").nullable()

    override val primaryKey = PrimaryKey(user, node)
}

@Serializable
data class UserNode(
    @SerialName("node_uuid")
    val nodeUUID: String,
    @SerialName("user_id")
    val userID: Int,

    @Serializable(DateSerializer::class) @SerialName("seen_at") val seenAt: LocalDateTime?,
    @Serializable(DateSerializer::class) @SerialName("edited_at") val editedAt: LocalDateTime?,
)
