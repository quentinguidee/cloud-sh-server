package com.quentinguidee.models

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable
import org.jetbrains.exposed.dao.id.IntIdTable

object Sessions : IntIdTable() {
    val user = reference("user_id", Users)
    val token = varchar("token", 63)
}

@Serializable
class Session(
    val id: Int?,
    @SerialName("user_id")
    val userID: Int,
    val token: String,
)
