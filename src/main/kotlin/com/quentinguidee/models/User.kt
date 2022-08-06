package com.quentinguidee.models

import com.quentinguidee.utils.DateSerializer
import kotlinx.serialization.Serializable
import org.jetbrains.exposed.dao.id.IntIdTable
import org.jetbrains.exposed.sql.javatime.datetime
import java.time.LocalDateTime

object Users : IntIdTable("users") {
    val username = varchar("username", 127).uniqueIndex()
    val name = varchar("name", 127).nullable()
    val email = varchar("email", 255).nullable()
    val profilePicture = varchar("profile_picture", 255).nullable()
    val role = varchar("role", 63).nullable()

    val createdAt = datetime("created_at").default(LocalDateTime.now())
    val updatedAt = datetime("updated_at").nullable()
}

@Serializable
data class User(
    val id: Int,
    val username: String,
    val name: String? = null,
    val email: String? = null,
    val profilePicture: String? = null,
    val role: String? = null,

    @Serializable(DateSerializer::class) val createdAt: LocalDateTime,
    @Serializable(DateSerializer::class) val updatedAt: LocalDateTime? = null,
)
