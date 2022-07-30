package com.quentinguidee.models

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable
import org.jetbrains.exposed.dao.id.IntIdTable

object GitHubUsers : IntIdTable("github_users") {
    val user = reference("user_id", Users)
    val username = varchar("username", 255)
}

@Serializable
data class GitHubUser(
    @SerialName("user_id")
    val userID: Int,
    val username: String,
)
