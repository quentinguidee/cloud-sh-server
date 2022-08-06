package com.quentinguidee.models

import kotlinx.serialization.Serializable
import org.jetbrains.exposed.dao.id.IntIdTable

object OAuthUsers : IntIdTable("oauth_users") {
    val user = reference("user_id", Users)
    val username = varchar("username", 255)
    val oAuthMethod = reference("oauth_method_id", OAuthMethods)
}

@Serializable
data class OAuthUser(
    val userID: Int,
    val username: String,
    val oAuthMethodID: Int,
)
