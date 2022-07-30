package com.quentinguidee.models

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable
import org.jetbrains.exposed.dao.id.IntIdTable

object Users : IntIdTable() {
    val username = varchar("username", 127).uniqueIndex()
    val name = varchar("name", 127)
    val email = varchar("email", 255)
    val profilePicture = varchar("profile_picture", 255)
    val role = varchar("role", 63).nullable()
}

@Serializable
data class User(
    val id: Int,
    val username: String,
    val name: String,
    val email: String,
    @SerialName("profile_picture")
    val profilePicture: String,
    val role: String? = null,
)
